package godaddy_provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type godaddyProviderModel struct {
	Key     types.String `tfsdk:"key"`
	Secret  types.String `tfsdk:"secret"`
	Baseurl types.String `tfsdk:"baseurl"`
}

// Metadata returns the provider type name.
func (p *godaddyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "st-godaddy"
}

// Schema defines the provider-level schema for configuration data.
func (p *godaddyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The godaddy domain provider is used to interact with the godaddy to manage domains from it. " +
			"The provider needs to be configured with the proper credentials before it can be used.",
		Attributes: map[string]schema.Attribute{
			"key": schema.StringAttribute{
				Description: "GoDaddy API Key.",
				Required:    true,
			},
			"secret": schema.StringAttribute{
				Description: "GoDaddy API Secret.",
				Required:    true,
			},
			"baseurl": schema.StringAttribute{
				Description: "GoDaddy Base Url(defaults to production).",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *godaddyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

	var config godaddyProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	key := config.Key.ValueString()
	secret := config.Secret.ValueString()
	baseUrl := config.Baseurl.ValueString()

	cfg := Config{
		Key:     key,
		Secret:  secret,
		BaseURL: baseUrl,
	}

	cli, err := cfg.Client()
	if err != nil {
		resp.Diagnostics.AddError("create godaddy client error", err.Error())
		return
	}

	resp.ResourceData = cli
}

func (p *godaddyProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGodaddyDomainResource,
	}
}

func (p *godaddyProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// New is a helper function to simplify provider server
func New() provider.Provider {
	return &godaddyProvider{}
}

type godaddyProvider struct{}
