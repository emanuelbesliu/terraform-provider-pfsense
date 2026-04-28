package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallVirtualIPsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallVirtualIPsDataSource)(nil)
)

type FirewallVirtualIPsModel struct {
	VirtualIPs types.List `tfsdk:"virtual_ips"`
}

func NewFirewallVirtualIPsDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallVirtualIPsDataSource{}
}

type FirewallVirtualIPsDataSource struct {
	client *pfsense.Client
}

func (m *FirewallVirtualIPsModel) Set(ctx context.Context, vips pfsense.VirtualIPs) diag.Diagnostics {
	var diags diag.Diagnostics

	vipModels := []FirewallVirtualIPModel{}
	for _, v := range vips {
		var vipModel FirewallVirtualIPModel
		diags.Append(vipModel.Set(ctx, v)...)
		vipModels = append(vipModels, vipModel)
	}

	vipsValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: FirewallVirtualIPModel{}.AttrTypes()}, vipModels)
	diags.Append(newDiags...)
	m.VirtualIPs = vipsValue

	return diags
}

func (d *FirewallVirtualIPsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_virtual_ips", req.ProviderTypeName)
}

func (d *FirewallVirtualIPsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all firewall virtual IPs.",
		MarkdownDescription: "Retrieves all [firewall virtual IPs](https://docs.netgate.com/pfsense/en/latest/firewall/virtual-ip-addresses.html).",
		Attributes: map[string]schema.Attribute{
			"virtual_ips": schema.ListNestedAttribute{
				Description: "List of virtual IPs.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *FirewallVirtualIPsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallVirtualIPsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallVirtualIPsModel

	vips, err := d.client.GetVirtualIPs(ctx)
	if addError(&resp.Diagnostics, "Unable to get Virtual IPs", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *vips)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
