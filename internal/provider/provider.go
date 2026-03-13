package provider

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

const (
	envURL      = "TF_PFSENSE_URL"
	envUsername = "TF_PFSENSE_USERNAME"
	envPassword = "TF_PFSENSE_PASSWORD"
)

var _ provider.Provider = (*pfSenseProvider)(nil)

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &pfSenseProvider{
			version: version,
		}
	}
}

type pfSenseProvider struct {
	version string
}

type pfSenseProviderModel struct {
	URL              types.String `tfsdk:"url"`
	Username         types.String `tfsdk:"username"`
	Password         types.String `tfsdk:"password"`
	TLSSkipVerify    types.Bool   `tfsdk:"tls_skip_verify"`
	MaxAttempts      types.Int64  `tfsdk:"max_attempts"`
	ConcurrentWrites types.Bool   `tfsdk:"concurrent_writes"`
}

func (p *pfSenseProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "pfsense"
	resp.Version = p.version
}

func (p *pfSenseProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Interact with pfSense firewall/router.",
		MarkdownDescription: "Interact with [pfSense](https://www.pfsense.org/) firewall/router.\n\nCredentials can be provided via provider configuration or environment variables. Environment variables are used as fallback when the corresponding attribute is not set in the provider block. This is useful for CI/CD pipelines where sensitive values should not be stored in source control.\n\n| Attribute | Environment Variable | Default |\n|---|---|---|\n| `url` | `TF_PFSENSE_URL` | `" + pfsense.DefaultURL + "` |\n| `username` | `TF_PFSENSE_USERNAME` | `" + pfsense.DefaultUsername + "` |\n| `password` | `TF_PFSENSE_PASSWORD` | _(required)_ |",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				Description:         fmt.Sprintf("pfSense administration URL. Can also be set with the %s environment variable. Defaults to '%s'.", envURL, pfsense.DefaultURL),
				MarkdownDescription: fmt.Sprintf("pfSense administration URL. Can also be set with the `%s` environment variable. Defaults to `%s`.", envURL, pfsense.DefaultURL),
				Optional:            true,
			},
			"username": schema.StringAttribute{
				Description:         fmt.Sprintf("pfSense administration username. Can also be set with the %s environment variable. Defaults to '%s'.", envUsername, pfsense.DefaultUsername),
				MarkdownDescription: fmt.Sprintf("pfSense administration username. Can also be set with the `%s` environment variable. Defaults to `%s`.", envUsername, pfsense.DefaultUsername),
				Optional:            true,
			},
			"password": schema.StringAttribute{
				Description:         fmt.Sprintf("pfSense administration password. Can also be set with the %s environment variable.", envPassword),
				MarkdownDescription: fmt.Sprintf("pfSense administration password. Can also be set with the `%s` environment variable.", envPassword),
				Optional:            true,
				Sensitive:           true,
			},
			"tls_skip_verify": schema.BoolAttribute{
				Description:         fmt.Sprintf("Skip verification of TLS certificates, defaults to '%t'.", pfsense.DefaultTLSSkipVerify),
				MarkdownDescription: fmt.Sprintf("Skip verification of TLS certificates, defaults to `%t`.", pfsense.DefaultTLSSkipVerify),
				Optional:            true,
			},
			"max_attempts": schema.Int64Attribute{
				Description:         fmt.Sprintf("Maximum number of attempts (only applicable for retryable errors), defaults to '%d'.", pfsense.DefaultMaxAttempts),
				MarkdownDescription: fmt.Sprintf("Maximum number of attempts (only applicable for retryable errors), defaults to `%d`.", pfsense.DefaultMaxAttempts),
				Optional:            true,
			},
			"concurrent_writes": schema.BoolAttribute{
				Description:         fmt.Sprintf("Enable concurrent pfSense configuration writes. Be aware that pfSense's XML configuration system does not support write operations at scale, which can lead to overwrites and unexpected behavior. Defaults to '%t'.", pfsense.DefaultConcurrentWrites),
				MarkdownDescription: fmt.Sprintf("Enable concurrent pfSense configuration writes. Be aware that pfSense's XML configuration system does not support write operations at scale, which can lead to overwrites and unexpected behavior. Defaults to `%t`.", pfsense.DefaultConcurrentWrites),
				Optional:            true,
			},
		},
	}
}

func (p *pfSenseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data pfSenseProviderModel

	tflog.Info(ctx, "Configuring pfSense client")

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.URL.IsUnknown() {
		path := path.Root("url")
		summary, detail := unknownProviderValue(path)
		resp.Diagnostics.AddAttributeError(path, summary, detail)
	}

	if data.Username.IsUnknown() {
		path := path.Root("username")
		summary, detail := unknownProviderValue(path)
		resp.Diagnostics.AddAttributeError(path, summary, detail)
	}

	if data.Password.IsUnknown() {
		path := path.Root("password")
		summary, detail := unknownProviderValue(path)
		resp.Diagnostics.AddAttributeError(path, summary, detail)
	}

	if data.TLSSkipVerify.IsUnknown() {
		path := path.Root("tls_skip_verify")
		summary, detail := unknownProviderValue(path)
		resp.Diagnostics.AddAttributeError(path, summary, detail)
	}

	if data.MaxAttempts.IsUnknown() {
		path := path.Root("max_attempts")
		summary, detail := unknownProviderValue(path)
		resp.Diagnostics.AddAttributeError(path, summary, detail)
	}

	if data.ConcurrentWrites.IsUnknown() {
		path := path.Root("concurrent_writes")
		summary, detail := unknownProviderValue(path)
		resp.Diagnostics.AddAttributeError(path, summary, detail)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve URL: HCL config > env var > client default
	urlValue := data.URL.ValueString()
	if data.URL.IsNull() || data.URL.ValueString() == "" {
		if v, ok := os.LookupEnv(envURL); ok {
			urlValue = v
			tflog.Debug(ctx, fmt.Sprintf("Using %s environment variable for URL", envURL))
		}
	}

	// Resolve username: HCL config > env var > client default
	usernameValue := data.Username.ValueString()
	if data.Username.IsNull() || data.Username.ValueString() == "" {
		if v, ok := os.LookupEnv(envUsername); ok {
			usernameValue = v
			tflog.Debug(ctx, fmt.Sprintf("Using %s environment variable for username", envUsername))
		}
	}

	// Resolve password: HCL config > env var (no default — required)
	passwordValue := data.Password.ValueString()
	if data.Password.IsNull() || data.Password.ValueString() == "" {
		if v, ok := os.LookupEnv(envPassword); ok {
			passwordValue = v
			tflog.Debug(ctx, fmt.Sprintf("Using %s environment variable for password", envPassword))
		}
	}

	if passwordValue == "" {
		resp.Diagnostics.AddError(
			"Missing pfSense password",
			fmt.Sprintf("The provider cannot find a password for pfSense. "+
				"Set the password in the provider configuration or use the %s environment variable.", envPassword),
		)

		return
	}

	var opts pfsense.Options

	if urlValue != "" {
		url, err := url.Parse(urlValue)
		if err != nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("url"),
				"pfSense URL cannot be parsed",
				err.Error(),
			)
		}

		opts.URL = url
	}

	if usernameValue != "" {
		opts.Username = usernameValue
	}

	opts.Password = passwordValue

	if !data.TLSSkipVerify.IsNull() {
		opts.TLSSkipVerify = data.TLSSkipVerify.ValueBoolPointer()
	}

	if !data.MaxAttempts.IsNull() {
		i := int(data.MaxAttempts.ValueInt64())
		opts.MaxAttempts = &i
	}

	if !data.ConcurrentWrites.IsNull() {
		opts.ConcurrentWrites = data.ConcurrentWrites.ValueBoolPointer()
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating pfSense client")

	client, err := pfsense.NewClient(ctx, &opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create pfSense client",
			"An unexpected error occurred when creating the pfSense client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"pfSense client URL: "+opts.URL.String()+"\n"+
				"pfSense client Error: "+err.Error(),
		)

		return
	}

	ctx = tflog.SetField(ctx, "pfsense_url", client.Options.URL.String())
	ctx = tflog.SetField(ctx, "pfsense_username", client.Options.Username)
	ctx = tflog.SetField(ctx, "pfsense_password", client.Options.Password)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "pfsense_password")

	resp.DataSourceData = client
	resp.ResourceData = client
	resp.EphemeralResourceData = client

	tflog.Info(ctx, "Configured pfSense client", map[string]any{"success": true})
}

func (p *pfSenseProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDHCPv4StaticMappingDataSource,
		NewDHCPv4StaticMappingsDataSource,
		NewDNSResolverConfigFileDataSource,
		NewDNSResolverConfigFilesDataSource,
		NewDNSResolverDomainOverrideDataSource,
		NewDNSResolverDomainOverridesDataSource,
		NewDNSResolverHostOverrideDataSource,
		NewDNSResolverHostOverridesDataSource,
		NewExecutePHPCommandDataSource,
		NewFirewallAliasesDataSource,
		NewFirewallIPAliasDataSource,
		NewFirewallPortAliasDataSource,
		NewGatewayDataSource,
		NewGatewayGroupDataSource,
		NewGatewayGroupsDataSource,
		NewGatewaysDataSource,
		NewInterfaceDataSource,
		NewInterfaceGroupDataSource,
		NewInterfaceGroupsDataSource,
		NewInterfacesDataSource,
		NewRouteDataSource,
		NewRoutesDataSource,
		NewSystemVersionDataSource,
		NewVLANDataSource,
		NewVLANsDataSource,
	}
}

func (p *pfSenseProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDNSResolverConfigFileResource,
		NewDNSResolverDomainOverrideResource,
		NewDNSResolverHostOverrideResource,
		NewFirewallIPAliasResource,
		NewFirewallPortAliasResource,
		NewGatewayResource,
		NewGatewayGroupResource,
		NewInterfaceResource,
		NewInterfaceGroupResource,
		NewRouteResource,
		NewVLANResource,
		NewDHCPv4StaticMappingResource,
		NewExecutePHPCommandResource,
	}
}
