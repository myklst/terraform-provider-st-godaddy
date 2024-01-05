package godaddy_provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type godaddyProvider struct{}

type godaddyProviderModel struct {
	Key     types.String `tfsdk:"key"`
	Secret  types.String `tfsdk:"secret"`
	Baseurl types.String `tfsdk:"baseurl"`
}

// New is a helper function to simplify provider server
func New() provider.Provider {
	return &godaddyProvider{}
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
				Sensitive:   true,
			},
			"secret": schema.StringAttribute{
				Description: "GoDaddy API Secret.",
				Required:    true,
				Sensitive:   true,
			},
			"baseurl": schema.StringAttribute{
				Description: "GoDaddy Base Url(defaults to production).",
				Required:    true,
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

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if config.Key.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			"Unknown API key",
			"The provider cannot create the Godaddy API client as there is an unknown configuration value for the"+
				"Godaddy API key. Set the value statically in the configuration, or use the GODADDY_API_KEY "+
				"environment variable.",
		)
	}
	if config.Secret.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret"),
			"Unknown API secret",
			"The provider cannot create the Godaddy API client as there is an unknown configuration value for the"+
				"Godaddy API secret. Set the value statically in the configuration, or use the GODADDY_API_SECRET "+
				"environment variable.",
		)
	}
	if config.Baseurl.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("baseurl"),
			"Unknown baseurl",
			"The provider cannot create the Godaddy API client as there is an unknown configuration value for the"+
				"Godaddy baseurl. Set the value statically in the configuration, or use the GODADDY_BASE_URL "+
				"environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	var (
		key,
		secret,
		baseUrl string
	)

	if !config.Key.IsNull() {
		key = config.Key.ValueString()
	} else {
		key = os.Getenv("GODADDY_API_KEY")
	}
	if key == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("key"),
			"Missing Godaddy API Key",
			"The provider cannot create the Godaddy API client as there is a "+
				"missing or empty value for the Godaddy API Key. Set the "+
				"API Key value in the configuration or use the GODADDY_API_KEY "+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if !config.Secret.IsNull() {
		secret = config.Secret.ValueString()
	} else {
		secret = os.Getenv("GODADDY_API_SECRET")
	}
	if secret == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("secret"),
			"Missing Godaddy API Secret",
			"The provider cannot create the Godaddy API client as there is a "+
				"missing or empty value for the Godaddy API Secret. Set the "+
				"API Secret value in the configuration or use the GODADDY_API_SECRET "+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if !config.Baseurl.IsNull() {
		baseUrl = config.Baseurl.ValueString()
	} else {
		baseUrl = os.Getenv("GODADDY_BASEURL")
	}
	if baseUrl == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("baseUrl"),
			"Missing Godaddy Base Url",
			"The provider cannot create the Godaddy API client as there is a "+
				"missing or empty value for the Godaddy base url. Set the "+
				"base url value in the configuration or use the GODADDY_BASEURL "+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	// If any of the expected configuration are missing, return
	// errors with provider-specific guidance.
	if resp.Diagnostics.HasError() {
		return
	}

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
		NewGodaddyNameServerResource,
	}
}

func (p *godaddyProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}
