package godaddy_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"log"
	"strconv"
	api2 "terraform-provider-st-godaddy/godaddy/api"
)

const (
	attrMode  = "mode"
	attrYears = "years"
)
const MODE_CREATE = "create"
const MODE_RENEW = "renew"

func NewGodaddyDomainResource() resource.Resource {
	return &godaddyDomainResource{}
}

type godaddyDomainResource struct {
	client *api2.Client
}

type godaddyDomainResourceModel struct {
	Domain  types.String `tfsdk:"domain"`
	Mode    types.String `tfsdk:"mode"`
	Years   types.Int64  `tfsdk:"years"`
	Contact types.String `tfsdk:"contact"`
}

// Metadata returns the resource godaddy_domain type name.
func (r *godaddyDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "godaddy_domain"
}

// Configure adds the provider configured client to the resource.
func (r *godaddyDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api2.Client)
}

func (r *godaddyDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import RecordId and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
}

// Schema defines the schema for the godaddy_domain resource.
func (r *godaddyDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides a godaddy_domain resource.",
		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Description: "Purchased available domain name on your account",
				Required:    true,
			},
			"mode": schema.StringAttribute{
				Description: "domain operation type, include create, renew.",
				Required:    true,
			},
			"years": schema.Int64Attribute{
				Description: "Number of years to register",
				Required:    true,
			},
			"contact": schema.StringAttribute{
				Description: "Contact info in json format",
				Required:    true,
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
	mode := plan.Mode.ValueString()
	years := plan.Years.ValueInt64()
	contact := plan.Contact.ValueString()

	switch mode {
	case MODE_CREATE:
		var contactInfo api2.RegisterDomainInfo
		diag1 := readContactInfo(contact, &contactInfo)
		resp.Diagnostics.Append(diag1)
		if resp.Diagnostics.HasError() {
			return
		}
		diag2 := createDomain(ctx, r.client, domain, years, contactInfo)
		resp.Diagnostics.Append(diag2)
		if resp.Diagnostics.HasError() {
			return
		}
	case MODE_RENEW:
		diag := renewDomain(ctx, r.client, domain, years)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	default:
		resp.Diagnostics.AddError("invalid mode value", mode)
	}

	// Set state items
	state := &godaddyDomainResourceModel{}
	state.Mode = plan.Mode
	state.Domain = plan.Domain
	state.Years = plan.Years

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

	_, err := r.client.GetDomain(domain)

	if err == nil {
		setStateDiags := resp.State.Set(ctx, state)
		resp.Diagnostics.Append(setStateDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
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

	newDomain := plan.Domain.ValueString()
	newMode := plan.Mode.ValueString()
	newYear := plan.Years.ValueInt64()
	contact := plan.Contact.ValueString()

	var state *godaddyDomainResourceModel
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
	oldDomain := state.Domain.ValueString()

	switch newMode {
	case MODE_CREATE:
		//delete old domain first
		if oldDomain != newDomain {
			diag1 := deleteDomain(ctx, r.client, oldDomain)
			resp.Diagnostics.Append(diag1)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		var contactInfo api2.RegisterDomainInfo
		diag2 := readContactInfo(contact, &contactInfo)
		resp.Diagnostics.Append(diag2)
		if resp.Diagnostics.HasError() {
			return
		}
		//create new domain then
		diag3 := createDomain(ctx, r.client, newDomain, newYear, contactInfo)
		resp.Diagnostics.Append(diag3)
		if resp.Diagnostics.HasError() {
			return
		}

	case MODE_RENEW:
		//can't do anything about old domain
		diag := renewDomain(ctx, r.client, newDomain, newYear)
		resp.Diagnostics.Append(diag)
		if resp.Diagnostics.HasError() {
			return
		}
	default:
		resp.Diagnostics.AddError("invalid mode value", newMode)
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

	diag := deleteDomain(ctx, r.client, domainName)
	resp.Diagnostics.Append(diag)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func createDomain(cxt context.Context, client *api2.Client, domainName string, year int64, _domainInfo api2.RegisterDomainInfo) diag.Diagnostic {

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

func renewDomain(cxt context.Context, client *api2.Client, domainName string, year int64) diag.Diagnostic {

	err := client.DomainRenew(domainName, strconv.FormatInt(year, 10))
	if err != nil {
		fmtlog(cxt, "Renew [%s] failed!", domainName)
		return DiagnosticErrorOf(err, "Renew [%s] failed!!", domainName)
	}
	fmtlog(cxt, "Renew [%s] success!", domainName)
	return nil
}

func deleteDomain(cxt context.Context, client *api2.Client, domainName string) diag.Diagnostic {

	err := client.DomainCancel(domainName)
	if err != nil {
		fmtlog(cxt, "Delete [%s] failed!", domainName)
		return DiagnosticErrorOf(err, "Delete [%s] failed!", domainName)
	}
	fmtlog(cxt, "Delete [%s] success!", domainName)
	return nil
}

func fmtlog(ctx context.Context, format string, a ...any) {
	msg := fmt.Sprintf(format, a)
	tflog.Info(ctx, msg)
}

func DiagnosticErrorOf(err error, format string, a ...any) diag.Diagnostic {
	msg := fmt.Sprintf(format, a)
	if err != nil {
		return diag.NewErrorDiagnostic(msg, err.Error())
	} else {
		return diag.NewErrorDiagnostic(msg, "")
	}
}

func readContactInfo(contact string, domainInfo *api2.RegisterDomainInfo) diag.Diagnostic {

	var contactInfo api2.Contact

	err := json.Unmarshal([]byte(contact), &contactInfo)

	if err != nil {
		return DiagnosticErrorOf(nil, "parse contact info failed!, json: ", contact)
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
