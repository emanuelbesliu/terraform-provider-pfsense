package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*OpenVPNServersDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*OpenVPNServersDataSource)(nil)
)

func NewOpenVPNServersDataSource() datasource.DataSource { //nolint:ireturn
	return &OpenVPNServersDataSource{}
}

type OpenVPNServersDataSource struct {
	client *pfsense.Client
}

func (d *OpenVPNServersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_openvpn_servers", req.ProviderTypeName)
}

func (d *OpenVPNServersDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all OpenVPN server instances.",
		MarkdownDescription: "Retrieves all OpenVPN [server](https://docs.netgate.com/pfsense/en/latest/vpn/openvpn/index.html) instances.",
		Attributes: map[string]schema.Attribute{
			"all": schema.ListNestedAttribute{
				Description: "All OpenVPN server instances.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: openVPNServerComputedAttributes(false),
				},
			},
		},
	}
}

func (d *OpenVPNServersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *OpenVPNServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OpenVPNServersModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	servers, err := d.client.GetOpenVPNServers(ctx)
	if addError(&resp.Diagnostics, "Unable to get OpenVPN servers", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *servers)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
