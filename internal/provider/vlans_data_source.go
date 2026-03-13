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
	_ datasource.DataSource              = (*VLANsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*VLANsDataSource)(nil)
)

type VLANsModel struct {
	VLANs types.List `tfsdk:"vlans"`
}

func NewVLANsDataSource() datasource.DataSource { //nolint:ireturn
	return &VLANsDataSource{}
}

type VLANsDataSource struct {
	client *pfsense.Client
}

func (m *VLANsModel) Set(ctx context.Context, vlans pfsense.VLANs) diag.Diagnostics {
	var diags diag.Diagnostics

	vlanModels := []VLANModel{}
	for _, v := range vlans {
		var vlanModel VLANModel
		diags.Append(vlanModel.Set(ctx, v)...)
		vlanModels = append(vlanModels, vlanModel)
	}

	vlansValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: VLANModel{}.AttrTypes()}, vlanModels)
	diags.Append(newDiags...)
	m.VLANs = vlansValue

	return diags
}

func (d *VLANsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_vlans", req.ProviderTypeName)
}

func (d *VLANsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all VLANs. VLANs allow segmenting a physical network into virtual networks using 802.1Q tagging.",
		MarkdownDescription: "Retrieves all [VLANs](https://docs.netgate.com/pfsense/en/latest/interfaces/vlan.html). VLANs allow segmenting a physical network into virtual networks using 802.1Q tagging.",
		Attributes: map[string]schema.Attribute{
			"vlans": schema.ListNestedAttribute{
				Description: "List of VLANs.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
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
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *VLANsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *VLANsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VLANsModel

	vlans, err := d.client.GetVLANs(ctx)
	if addError(&resp.Diagnostics, "Unable to get VLANs", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *vlans)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
