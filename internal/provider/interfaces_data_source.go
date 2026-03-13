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
	_ datasource.DataSource              = (*InterfacesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*InterfacesDataSource)(nil)
)

type InterfacesModel struct {
	Interfaces types.List `tfsdk:"interfaces"`
}

func NewInterfacesDataSource() datasource.DataSource { //nolint:ireturn
	return &InterfacesDataSource{}
}

type InterfacesDataSource struct {
	client *pfsense.Client
}

func (m *InterfacesModel) Set(ctx context.Context, ifaces pfsense.Interfaces) diag.Diagnostics {
	var diags diag.Diagnostics

	ifaceModels := []InterfaceModel{}
	for _, i := range ifaces {
		var ifaceModel InterfaceModel
		diags.Append(ifaceModel.Set(ctx, i)...)
		ifaceModels = append(ifaceModels, ifaceModel)
	}

	ifacesValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: InterfaceModel{}.AttrTypes()}, ifaceModels)
	diags.Append(newDiags...)
	m.Interfaces = ifacesValue

	return diags
}

func (d *InterfacesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_interfaces", req.ProviderTypeName)
}

func (d *InterfacesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all interface assignments.",
		MarkdownDescription: "Retrieves all [interface assignments](https://docs.netgate.com/pfsense/en/latest/interfaces/index.html).",
		Attributes: map[string]schema.Attribute{
			"interfaces": schema.ListNestedAttribute{
				Description: "List of interface assignments.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"logical_name": schema.StringAttribute{
							Description: InterfaceModel{}.descriptions()["logical_name"].Description,
							Computed:    true,
						},
						"port": schema.StringAttribute{
							Description: InterfaceModel{}.descriptions()["port"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: InterfaceModel{}.descriptions()["description"].Description,
							Computed:    true,
						},
						"enabled": schema.BoolAttribute{
							Description: InterfaceModel{}.descriptions()["enabled"].Description,
							Computed:    true,
						},
						"ipv4_type": schema.StringAttribute{
							Description:         InterfaceModel{}.descriptions()["ipv4_type"].Description,
							MarkdownDescription: InterfaceModel{}.descriptions()["ipv4_type"].MarkdownDescription,
							Computed:            true,
						},
						"ipv4_address": schema.StringAttribute{
							Description: InterfaceModel{}.descriptions()["ipv4_address"].Description,
							Computed:    true,
						},
						"ipv4_subnet": schema.StringAttribute{
							Description: InterfaceModel{}.descriptions()["ipv4_subnet"].Description,
							Computed:    true,
						},
						"ipv4_gateway": schema.StringAttribute{
							Description: InterfaceModel{}.descriptions()["ipv4_gateway"].Description,
							Computed:    true,
						},
						"ipv6_type": schema.StringAttribute{
							Description:         InterfaceModel{}.descriptions()["ipv6_type"].Description,
							MarkdownDescription: InterfaceModel{}.descriptions()["ipv6_type"].MarkdownDescription,
							Computed:            true,
						},
						"ipv6_address": schema.StringAttribute{
							Description: InterfaceModel{}.descriptions()["ipv6_address"].Description,
							Computed:    true,
						},
						"ipv6_subnet": schema.StringAttribute{
							Description: InterfaceModel{}.descriptions()["ipv6_subnet"].Description,
							Computed:    true,
						},
						"ipv6_gateway": schema.StringAttribute{
							Description: InterfaceModel{}.descriptions()["ipv6_gateway"].Description,
							Computed:    true,
						},
						"spoof_mac": schema.StringAttribute{
							Description: InterfaceModel{}.descriptions()["spoof_mac"].Description,
							Computed:    true,
						},
						"mtu": schema.Int64Attribute{
							Description: InterfaceModel{}.descriptions()["mtu"].Description,
							Computed:    true,
						},
						"mss": schema.Int64Attribute{
							Description: InterfaceModel{}.descriptions()["mss"].Description,
							Computed:    true,
						},
						"block_private": schema.BoolAttribute{
							Description: InterfaceModel{}.descriptions()["block_private"].Description,
							Computed:    true,
						},
						"block_bogons": schema.BoolAttribute{
							Description: InterfaceModel{}.descriptions()["block_bogons"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *InterfacesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *InterfacesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InterfacesModel

	ifaces, err := d.client.GetInterfaces(ctx)
	if addError(&resp.Diagnostics, "Unable to get interfaces", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *ifaces)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
