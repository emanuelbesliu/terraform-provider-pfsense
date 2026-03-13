package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*DNSResolverConfigFileDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*DNSResolverConfigFileDataSource)(nil)
)

type DNSResolverConfigFileDataSourceModel struct {
	DNSResolverConfigFileModel
}

func NewDNSResolverConfigFileDataSource() datasource.DataSource { //nolint:ireturn
	return &DNSResolverConfigFileDataSource{}
}

type DNSResolverConfigFileDataSource struct {
	client *pfsense.Client
}

func (d *DNSResolverConfigFileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dnsresolver_configfile", req.ProviderTypeName)
}

func (d *DNSResolverConfigFileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single DNS resolver config file by name. Config files contain custom Unbound configuration.",
		MarkdownDescription: "Retrieves a single DNS resolver [config file](https://docs.netgate.com/pfsense/en/latest/services/dns/resolver.html) by name. Config files contain custom Unbound configuration.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: DNSResolverConfigFileModel{}.descriptions()["name"].Description,
				Required:    true,
			},
			"content": schema.StringAttribute{
				Description:         DNSResolverConfigFileModel{}.descriptions()["content"].Description,
				MarkdownDescription: DNSResolverConfigFileModel{}.descriptions()["content"].MarkdownDescription,
				Computed:            true,
			},
		},
	}
}

func (d *DNSResolverConfigFileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *DNSResolverConfigFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DNSResolverConfigFileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	configFile, err := d.client.GetDNSResolverConfigFile(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Unable to get config file", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *configFile)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
