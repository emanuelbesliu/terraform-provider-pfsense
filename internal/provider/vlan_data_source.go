package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*VLANDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*VLANDataSource)(nil)
)

type VLANDataSourceModel struct {
	VLANModel
}

func NewVLANDataSource() datasource.DataSource { //nolint:ireturn
	return &VLANDataSource{}
}

type VLANDataSource struct {
	client *pfsense.Client
}

func (d *VLANDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_vlan", req.ProviderTypeName)
}

func (d *VLANDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single VLAN by its interface name (e.g. 'vmx0.100'). VLANs allow segmenting a physical network into virtual networks using 802.1Q tagging.",
		MarkdownDescription: "Retrieves a single [VLAN](https://docs.netgate.com/pfsense/en/latest/interfaces/vlan.html) by its interface name (e.g. `vmx0.100`). VLANs allow segmenting a physical network into virtual networks using 802.1Q tagging.",
		Attributes: map[string]schema.Attribute{
			"parent_interface": schema.StringAttribute{
				Description: VLANModel{}.descriptions()["parent_interface"].Description,
				Computed:    true,
			},
			"tag": schema.Int64Attribute{
				Description: VLANModel{}.descriptions()["tag"].Description,
				Computed:    true,
			},
			"pcp": schema.Int64Attribute{
				Description: VLANModel{}.descriptions()["pcp"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: VLANModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
			"vlan_interface": schema.StringAttribute{
				Description: VLANModel{}.descriptions()["vlan_interface"].Description,
				Required:    true,
			},
		},
	}
}

func (d *VLANDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *VLANDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VLANDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vlan, err := d.client.GetVLAN(ctx, data.VLANInterface.ValueString())
	if addError(&resp.Diagnostics, "Unable to get VLAN", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *vlan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
