package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*DNSResolverConfigFilesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*DNSResolverConfigFilesDataSource)(nil)
)

func NewDNSResolverConfigFilesDataSource() datasource.DataSource { //nolint:ireturn
	return &DNSResolverConfigFilesDataSource{}
}

type DNSResolverConfigFilesDataSource struct {
	client *pfsense.Client
}

func (d *DNSResolverConfigFilesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_dnsresolver_configfiles", req.ProviderTypeName)
}

func (d *DNSResolverConfigFilesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all DNS resolver config files. Config files contain custom Unbound configuration.",
		MarkdownDescription: "Retrieves all DNS resolver [config files](https://docs.netgate.com/pfsense/en/latest/services/dns/resolver.html). Config files contain custom Unbound configuration.",
		Attributes: map[string]schema.Attribute{
			"all": schema.ListNestedAttribute{
				Description: "All config files.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: DNSResolverConfigFileModel{}.descriptions()["name"].Description,
							Computed:    true,
						},
						"content": schema.StringAttribute{
							Description:         DNSResolverConfigFileModel{}.descriptions()["content"].Description,
							MarkdownDescription: DNSResolverConfigFileModel{}.descriptions()["content"].MarkdownDescription,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DNSResolverConfigFilesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *DNSResolverConfigFilesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DNSResolverConfigFilesModel

	configFiles, err := d.client.GetDNSResolverConfigFiles(ctx)
	if addError(&resp.Diagnostics, "Unable to get config files", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *configFiles)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
