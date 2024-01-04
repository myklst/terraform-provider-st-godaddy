package godaddy_provider

import (
	"context"

	"github.com/myklst/terraform-provider-st-godaddy/godaddy/api"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NewGodaddyNameServerResource() resource.Resource {
	return &godaddyNameServerResource{}
}

type godaddyNameServerResource struct {
	client *api.Client
}

type godaddyNameServerResourceModel struct {
	Domain      types.String `tfsdk:"domain"`
	NameServers types.List   `tfsdk:"nameservers"`
}

// Metadata returns the resource godaddy_domain type name.
func (r *godaddyNameServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nameserver_attachment"
}

// Configure adds the provider configured client to the resource.
func (r *godaddyNameServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*api.Client)
}

func (r *godaddyNameServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage the nameservers of a domain in GoDaddy",
		Attributes: map[string]schema.Attribute{
			"domain": &schema.StringAttribute{
				Description: "Domain name to manage in GoDaddy",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"nameservers": &schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "The authoritative name server for this domain",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(2),
				},
			},
		},
	}
}

// Create a new godaddy_nameserver resource
func (r *godaddyNameServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	fmtlog(ctx, "[resourceNameServerCreate!]")
	var plan *godaddyNameServerResourceModel
	getPlanDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getPlanDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameservers = []string{}
	for _, x := range plan.NameServers.Elements() {
		nameservers = append(nameservers, x.String())
	}

	domain := plan.Domain.ValueString()
	fmtlog(ctx, plan.NameServers.String())
	diag1 := r.createNameServer(ctx, domain, nameservers)

	resp.Diagnostics.Append(diag1)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state items
	state := &godaddyNameServerResourceModel{
		Domain:      plan.Domain,
		NameServers: plan.NameServers,
	}

	setStateDiags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *godaddyNameServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	fmtlog(ctx, "[resourceNameServerRead!]")

	var state *godaddyNameServerResourceModel
	getStateDiags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domain := state.Domain.ValueString()
	response, err := r.client.GetDomainNameServers(domain)

	var ns []attr.Value

	for _, element := range response {
		ns = append(ns, types.StringValue(element))
	}

	if err == nil {
		state := &godaddyNameServerResourceModel{}
		state.Domain = types.StringValue(domain)
		nameServersList, err := types.ListValue(types.StringType, ns)
		if err != nil {
			return
		}

		state.NameServers = nameServersList
		setStateDiags := resp.State.Set(ctx, state)
		resp.Diagnostics.Append(setStateDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	} else {
		resp.Diagnostics.AddError("Get nameservers of domain error ", err.Error())
	}
}

func (r *godaddyNameServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	fmtlog(ctx, "[resourceNameServerUpdate!]")

	var plan *godaddyNameServerResourceModel

	getPlanDiags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(getPlanDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var nameservers []string
	for _, x := range plan.NameServers.Elements() {
		nameservers = append(nameservers, x.String())
	}

	domain := plan.Domain.ValueString()

	diag1 := r.createNameServer(ctx, domain, nameservers)

	resp.Diagnostics.Append(diag1)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set state items
	state := &godaddyNameServerResourceModel{
		Domain:      plan.Domain,
		NameServers: plan.NameServers,
	}
	setStateDiags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(setStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *godaddyNameServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	fmtlog(ctx, "[resourceNameServerUpdate!]")

	var state *godaddyNameServerResourceModel
	getStateDiags := req.State.Get(ctx, &state)

	resp.Diagnostics.Append(getStateDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	domainName := state.Domain.ValueString()

	diag := r.deleteNameServer(ctx, domainName)
	resp.Diagnostics.Append(diag)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *godaddyNameServerResource) createNameServer(ctx context.Context, domainName string, nameservers []string) diag.Diagnostic {
	client := r.client

	var ns api.NameServers
	ns.NameServers = nameservers

	err := client.UpdateNameServers(domainName, ns)
	if err != nil {
		fmtlog(ctx, "Setting nameservers for [%s] failed!", domainName)
		return DiagnosticErrorOf(err, "Setting nameservers for [%s] failed!!", domainName)
	}

	fmtlog(ctx, "Setting nameservers for [%s] success!", domainName)
	return nil
}

func (r *godaddyNameServerResource) deleteNameServer(ctx context.Context, domainName string) diag.Diagnostic {
	client := r.client

	var ns api.NameServers
	ns.NameServers = []string{"NS07.domaincontrol.com", "NS08.domaincontrol.com"}

	err := client.UpdateNameServers(domainName, ns)

	if err != nil {
		fmtlog(ctx, "Setting nameservers for [%s] failed!", domainName)
		return DiagnosticErrorOf(err, "Setting nameservers for [%s] failed!!", domainName)
	}

	fmtlog(ctx, "Setting nameservers for [%s] success!", domainName)
	return nil
}
