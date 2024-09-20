package provider

import (
	"context"
	"os"
	"terraform-provider-kuma/internal/kuma"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/spf13/viper"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &kumaProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &kumaProvider{
			version: version,
		}
	}
}

type kumaProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
}

type KumaConfiguration struct {
	Host     string
	Username string
	Password string
}

// Uptime KumaProvider is the provider implementation.
type kumaProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *kumaProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kuma"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *kumaProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "URL for Uptime Kuma API Server. May also be provided via KUMA_API_HOST environment variable.",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username for Uptime Kuma API. May also be provided via KUMA_API_USERNAME environment variable.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password for Uptime Kuma API. May also be provided via KUMA_API_PASSWORD environment variable.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *kumaProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config kumaProviderModel

	tflog.Info(ctx, "Configuring Kuma API client")

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var c struct {
		Host     string
		Username string
		Password string
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/kuma/")

	if err := viper.ReadInConfig(); err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root(""),
			"Unknown Uptime Kuma API Host",
			"The provdier can read config file from home path. "+err.Error(),
		)
	}

	if err := viper.Unmarshal(&c); err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root(""),
			"Unknown Uptime Kuma API Host",
			"The provdier can marshal config file from home path. "+err.Error(),
		)
	}

	host := os.Getenv("KUMA_API_HOST")
	username := os.Getenv("KUMA_API_USERNAME")
	password := os.Getenv("KUMA_API_PASSWORD")

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	} else if host == "" {
		host = c.Host
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	} else if username == "" {
		username = c.Username
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	} else if password == "" {
		password = c.Password
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Uptime Kuma API Host",
			"The provider cannot create the Uptime Kuma API client as there is a missing or empty value for the Uptime Kuma API host. "+
				"Set the host value in the configuration or use the Uptime KUMA_API_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Uptime Kuma API Username",
			"The provider cannot create the Uptime Kuma API client as there is a missing or empty value for the Uptime Kuma API username. "+
				"Set the username value in the configuration or use the Uptime KUMA_API_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Uptime Kuma API Password",
			"The provider cannot create the Uptime Kuma API client as there is a missing or empty value for the Uptime Kuma API password. "+
				"Set the password value in the configuration or use the Uptime KUMA_API_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "kuma_api_host", host)
	ctx = tflog.SetField(ctx, "kuma_api_username", username)
	ctx = tflog.SetField(ctx, "kuma_api_password", password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "kuma_api_password")

	tflog.Debug(ctx, "Creating Kuma API client")

	// Create a new Uptime Kuma client using the configuration values
	client, err := kuma.NewClient(&host, &username, &password)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Uptime Kuma API Client",
			"An unexpected error occurred when creating the Uptime Kuma API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Uptime Kuma Client Error: "+err.Error(),
		)
		return
	}

	// Make the Uptime Kuma client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Kuma API client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *kumaProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTagsDataSource,
		NewMonitorsDataSource,
		NewNotificationsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *kumaProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTagResource,
		NewHttpMonitorResource,
		NewGroupResource,
	}
}
