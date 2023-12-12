package godaddy_provider

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"terraform-provider-st-godaddy/godaddy/api"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

const (
	MODE_CREATE = "create"
	MODE_RENEW  = "renew"
	MODE_SKIP   = "skip"
)

func NewGodaddyDomainResource() resource.Resource {
	return &godaddyDomainResource{}
}

type godaddyDomainResource struct {
	client *api.Client
}

type godaddyDomainResourceModel struct {
	Domain           types.String `tfsdk:"domain"`
	MinDaysRemaining types.Int64  `tfsdk:"min_days_remaining"`
	Years            types.Int64  `tfsdk:"purchase_years"`
	Contact          types.String `tfsdk:"contact"`
	Expires          types.String `tfsdk:"expires"`
	Renew            types.Bool   `tfsdk:"renew"`
}

// Metadata returns the resource godaddy_domain type name.
func (r *godaddyDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_domain"
}

// Configure adds the provider configured client to the resource.
func (r *godaddyDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

func (r *godaddyDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import RecordId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
}

// Schema defines the schema for the godaddy_domain resource.
func (r *godaddyDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage a domain in Godaddy",
		Attributes: map[string]schema.Attribute{
			"domain": &schema.StringAttribute{
				Description: "Domain name to manage in NameCheap",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"min_days_remaining": &schema.Int64Attribute{
				MarkdownDescription: "The minimum amount of days remaining on the expiration of a domain before a " +
					"renewal is attempted. The default is `30`. A negative value means that the domain will " +
					"never be renewed. Zero value is not allowed",
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(30),
				Validators: []validator.Int64{
					int64validator.NoneOf(int64(0)),
				},
			},
			"purchase_years": &schema.Int64Attribute{
				MarkdownDescription: "Number of years to purchase and renew. The default is `1`.",
				Optional:            true,
				Computed:            true,
				Default:             int64default.StaticInt64(1),
			},
			"contact": schema.StringAttribute{
				Description: "Contact info in json format",
				Required:    true,
			},
			"expires": schema.StringAttribute{
				Description: "The ISO 8601 string representing the expiry date of the domain",
				Computed:    true,
			},
			"renew": schema.BoolAttribute{
				Description: "Whether to renew the domain. This is a special schema attribute used by the custom provider. Practitioners must not touch this value.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

// Create a new godaddy_domain resource
func (r *godaddyDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	fmtlog(ctx, "[resourceDomainCreate!]")
	var plan *godaddyDomainResourceModel
	getPlanDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getPlanDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := plan.Domain.ValueString()
	years := plan.Years.ValueInt64()
	contact := plan.Contact.ValueString()

	var contactInfo api.RegisterDomainInfo
	diag1 := r.readContactInfo(contact, &contactInfo)
	resp.Diagnostics.Append(diag1)
	if resp.Diagnostics.HasError() {
		return
	}
	diag2 := r.createDomain(ctx, domain, years, contactInfo)
	resp.Diagnostics.Append(diag2)
	if resp.Diagnostics.HasError() {
		return
	}

	var res api.Domain
	var diag3 diag.Diagnostic

	operation := func() error {
		res, diag3 = r.getDomain(ctx, domain)

		if  res.Status == "ACTIVE" {
			return nil
		} else {
			log.Println("domain expiry time is not yet in ACTIVE state")
			return errors.New("domain expiry time is not yet in ACTIVE state")
		}
	}

	err := backoff.Retry(operation, backoff.NewExponentialBackOff())

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	resp.Diagnostics.Append(diag3)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state items
	state := &godaddyDomainResourceModel{
		Domain:           plan.Domain,
		Years:            plan.Years,
		MinDaysRemaining: plan.MinDaysRemaining,
		Contact:          plan.Contact,
		Expires:          types.StringValue(res.Expires),
		Renew:            types.BoolValue(false),
	}

	setStateDiags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read godaddy_domain resource information
func (r *godaddyDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	fmtlog(ctx, "[resourceDomainRead!]")

	var state *godaddyDomainResourceModel
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := state.Domain.ValueString()

	res, err := r.client.GetDomain(domain)

	state = &godaddyDomainResourceModel{
		Domain:           state.Domain,
		Years:            state.Years,
		MinDaysRemaining: state.MinDaysRemaining,
		Contact:          state.Contact,
		Expires:          types.StringValue(res.Expires),
		Renew:            basetypes.NewBoolValue(false),
	}

	if err == nil {
		setStateDiags := resp.State.Set(ctx, state)
		resp.Diagnostics.Append(setStateDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		if strings.Contains(err.Error(), "Domain is invalid") {
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError("Get domain info error ", err.Error())
		}
	}

	newMode, diag := r.calculateMode(state)
	resp.Diagnostics.Append(diag)
	if resp.Diagnostics.HasError() {
		return
	}
	var renew bool

	switch newMode {

	case MODE_RENEW:
		renew = true
		state.Renew = types.BoolValue(renew)
		setStateDiags := resp.State.Set(ctx, state)
		resp.Diagnostics.Append(setStateDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	case MODE_SKIP:
		renew = false

	default:
		resp.Diagnostics.AddError("invalid mode value", newMode)
		return
	}
}

// Update godaddy_domain resource and sets the updated Terraform state on success.
func (r *godaddyDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	fmtlog(ctx, "[resourceRecordUpdate!]")

	var plan *godaddyDomainResourceModel
	getPlanDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getPlanDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *godaddyDomainResourceModel
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Renew.ValueBool() {
		fmtlog(ctx, "CalculateMode Complete,Renew = %s", state.Renew.String())
		diag := r.renewDomain(ctx, r.client, plan.Domain.ValueString(), plan.Years.ValueInt64())
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var res api.Domain
	var diag3 diag.Diagnostic

	operation := func() error {
		res, diag3 = r.getDomain(ctx, state.Domain.ValueString())

		const layout string = "2006-01-02T15:04:05.000Z"
		timeFromAPI, apiTimeErr := time.Parse(layout, res.Expires)
		timeFromStateFile, stateFileTimeErr := time.Parse(layout, state.Expires.ValueString())

		if apiTimeErr != nil {
			return apiTimeErr
		}

		if stateFileTimeErr != nil {
			return stateFileTimeErr
		}

		if timeFromAPI.After(timeFromStateFile) {
			return nil
		} else {
			log.Println("domain expiry time is not yet updated")
			return errors.New("domain expiry time is not yet updated")
		}
	}

	err := backoff.Retry(operation, backoff.NewExponentialBackOff())

	if err != nil {
		resp.Diagnostics.AddError(err.Error(), err.Error())
		return
	}

	resp.Diagnostics.Append(diag3)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state items
	state = &godaddyDomainResourceModel{}
	state.Domain = plan.Domain
	state.Years = plan.Years
	state.MinDaysRemaining = plan.MinDaysRemaining
	state.Contact = plan.Contact
	state.Expires = types.StringValue(res.Expires)
	state.Renew = types.BoolValue(false)

	setStateDiags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete godaddy_domain resource and removes the Terraform state on success.
func (r *godaddyDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	fmtlog(ctx, "[resourceRecordDelete!]")

	var state *godaddyDomainResourceModel

	// Retrieve values from plan
	getStateDiags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainName := state.Domain.ValueString()

	diag := r.deleteDomain(ctx, domainName)
	resp.Diagnostics.Append(diag)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *godaddyDomainResource) createDomain(cxt context.Context, domainName string, year int64, _domainInfo api.RegisterDomainInfo) diag.Diagnostic {

	client := r.client
	var domains []string
	domains = append(domains, domainName)
	log.Println("domain", domainName, "do not exist, check whether it's available to purchase....")
	available, err := client.DomainAvailable(domains)
	if err != nil {
		return DiagnosticErrorOf(err, "DomainAvailable for [%s] failed!!", domainName)
	}
	if !available {
		return DiagnosticErrorOf(nil, "[%s] is not available!", domainName)
	}

	//extract tld
	/*
		tld, _, err := gotld.GetTld(domainName)
		agreement, err := client.GetAgreement(tld.Tld, false)
		if err != nil {
			return DiagnosticErrorOf(err, "GetAgreement for  [%s] failed!!", domainName)
		}
		//construct agreement keys
		var agreementKeys []string
		for _, v := range agreement {
			agreementKeys = append(agreementKeys, v.AgreementKey)
		}*/

	err = client.Purchase(domainName, _domainInfo, strconv.FormatInt(year, 10))

	if err != nil {
		fmtlog(cxt, "Creating [%s] failed!", domainName)
		return DiagnosticErrorOf(err, "Creating [%s] failed!!", domainName)
	}
	fmtlog(cxt, "Creating [%s] success!", domainName)
	return nil
}

func (r *godaddyDomainResource) getDomain(ctx context.Context, domainName string) (api.Domain, diag.Diagnostic) {
	client := r.client
	res, err := client.GetDomain(domainName)

	if err != nil {
		return api.Domain{}, DiagnosticErrorOf(err, "Unable to get domain [%s]", domainName)
	} else {
		return *res, nil
	}
}

func (r *godaddyDomainResource) renewDomain(cxt context.Context, client *api.Client, domainName string, year int64) diag.Diagnostic {

	err := client.DomainRenew(domainName, strconv.FormatInt(year, 10))
	if err != nil {
		fmtlog(cxt, "Renew [%s] failed!", domainName)
		return DiagnosticErrorOf(err, "Renew [%s] failed!!", domainName)
	}
	fmtlog(cxt, "Renew [%s] success!", domainName)
	return nil
}

func (r *godaddyDomainResource) deleteDomain(cxt context.Context, domainName string) diag.Diagnostic {

	err := r.client.DomainCancel(domainName)
	if err != nil {
		fmtlog(cxt, "Delete [%s] failed!", domainName)
		return DiagnosticErrorOf(err, "Delete [%s] failed!", domainName)
	}
	fmtlog(cxt, "Delete [%s] success!", domainName)
	return nil
}

func (r *godaddyDomainResource) calculateMode(state *godaddyDomainResourceModel) (string, diag.Diagnostic) {
	minDaysRemain := state.MinDaysRemaining.ValueInt64()
	expires := state.Expires.ValueString()

	return ParseTimeAndCalculateMode(expires, minDaysRemain)
}

func ParseTimeAndCalculateMode(expires string, minDaysRemain int64) (string, diag.Diagnostic) {
	const layout = "2006-01-02T15:04:05.000Z"

	exp, err := time.Parse(layout, expires)
	if err != nil {
		return "", DiagnosticErrorOf(err, "Time string [%s] cannot be parsed", expires)
	}

	diff := time.Until(exp)
	minDays, err := time.ParseDuration(fmt.Sprintf("%dh", minDaysRemain*24))

	if err != nil {
		return "", DiagnosticErrorOf(err, "Value of Min Days Remain [%d] is invalid", minDaysRemain)
	}

	if diff < minDays {
		return MODE_RENEW, nil
	}

	return MODE_SKIP, nil
}

func fmtlog(ctx context.Context, format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	tflog.Info(ctx, msg)
}

func DiagnosticErrorOf(err error, format string, a ...any) diag.Diagnostic {
	msg := fmt.Sprintf(format, a...)
	if err != nil {
		return diag.NewErrorDiagnostic(msg, err.Error())
	} else {
		return diag.NewErrorDiagnostic(msg, "")
	}
}

func (r *godaddyDomainResource) readContactInfo(contact string, domainInfo *api.RegisterDomainInfo) diag.Diagnostic {

	var contactInfo api.Contact

	err := json.Unmarshal([]byte(contact), &contactInfo)

	if err != nil {
		return DiagnosticErrorOf(err, "parse contact info failed!, json: %s", contact)
	}

	//admin
	domainInfo.ContactAdmin = contactInfo
	//ContactBilling
	domainInfo.ContactBilling = contactInfo
	//reg
	domainInfo.ContactRegistrant = contactInfo
	//tech
	domainInfo.ContactTech = contactInfo

	/*  for debug
	log.Println(domainInfo.ContactAdmin.NameLast)
	log.Println(domainInfo.ContactAdmin.NameFirst)
	log.Println(domainInfo.ContactAdmin.Phone)
	log.Println(domainInfo.ContactAdmin.Fax)
	log.Println(domainInfo.ContactAdmin.Organization)
	log.Println(domainInfo.ContactAdmin.NameMiddle)
	log.Println(domainInfo.ContactAdmin.JobTitle)
	log.Println(domainInfo.ContactAdmin.Email)
	log.Println(domainInfo.ContactAdmin.AddressMailing.Address1)
	log.Println(domainInfo.ContactAdmin.AddressMailing.Country)
	*/

	return nil
}
