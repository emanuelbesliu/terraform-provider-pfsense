package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallVirtualIPDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallVirtualIPDataSource)(nil)
)

type FirewallVirtualIPDataSourceModel struct {
	FirewallVirtualIPModel
}

func NewFirewallVirtualIPDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallVirtualIPDataSource{}
}

type FirewallVirtualIPDataSource struct {
	client *pfsense.Client
}

func (d *FirewallVirtualIPDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_virtual_ip", req.ProviderTypeName)
}

func (d *FirewallVirtualIPDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single firewall virtual IP by its unique ID.",
		MarkdownDescription: "Retrieves a single [firewall virtual IP](https://docs.netgate.com/pfsense/en/latest/firewall/virtual-ip-addresses.html) by its unique ID.",
		Attributes: map[string]schema.Attribute{
			"mode": schema.StringAttribute{
				Description: FirewallVirtualIPModel{}.descriptions()["mode"].Description,
				Computed:    true,
			},
			"interface": schema.StringAttribute{
				Description: FirewallVirtualIPModel{}.descriptions()["interface"].Description,
				Computed:    true,
			},
			"vhid": schema.Int64Attribute{
				Description: FirewallVirtualIPModel{}.descriptions()["vhid"].Description,
				Computed:    true,
			},
			"advskew": schema.Int64Attribute{
				Description: FirewallVirtualIPModel{}.descriptions()["advskew"].Description,
				Computed:    true,
			},
			"advbase": schema.Int64Attribute{
				Description: FirewallVirtualIPModel{}.descriptions()["advbase"].Description,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: FirewallVirtualIPModel{}.descriptions()["password"].Description,
				Computed:    true,
				Sensitive:   true,
			},
			"subnet": schema.StringAttribute{
				Description: FirewallVirtualIPModel{}.descriptions()["subnet"].Description,
				Computed:    true,
			},
			"subnet_bits": schema.Int64Attribute{
				Description: FirewallVirtualIPModel{}.descriptions()["subnet_bits"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: FirewallVirtualIPModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
			"unique_id": schema.StringAttribute{
				Description: FirewallVirtualIPModel{}.descriptions()["unique_id"].Description,
				Required:    true,
			},
		},
	}
}

func (d *FirewallVirtualIPDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallVirtualIPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallVirtualIPDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vip, err := d.client.GetVirtualIP(ctx, data.UniqueID.ValueString())
	if addError(&resp.Diagnostics, "Unable to get Virtual IP", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *vip)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
