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
	_ datasource.DataSource              = (*InterfaceBridgesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*InterfaceBridgesDataSource)(nil)
)

type InterfaceBridgesModel struct {
	Bridges types.List `tfsdk:"bridges"`
}

func NewInterfaceBridgesDataSource() datasource.DataSource { //nolint:ireturn
	return &InterfaceBridgesDataSource{}
}

type InterfaceBridgesDataSource struct {
	client *pfsense.Client
}

func (m *InterfaceBridgesModel) Set(ctx context.Context, bridges pfsense.Bridges) diag.Diagnostics {
	var diags diag.Diagnostics

	models := []InterfaceBridgeModel{}
	for _, b := range bridges {
		var model InterfaceBridgeModel
		diags.Append(model.Set(ctx, b)...)
		models = append(models, model)
	}

	listValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: InterfaceBridgeModel{}.AttrTypes()}, models)
	diags.Append(newDiags...)
	m.Bridges = listValue

	return diags
}

func (d *InterfaceBridgesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_interface_bridges", req.ProviderTypeName)
}

func (d *InterfaceBridgesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all bridge interfaces.",
		MarkdownDescription: "Retrieves all [bridge interfaces](https://docs.netgate.com/pfsense/en/latest/interfaces/bridges.html).",
		Attributes: map[string]schema.Attribute{
			"bridges": schema.ListNestedAttribute{
				Description: "List of bridge interfaces.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: interfaceBridgeDataSourceAttributes(false),
				},
			},
		},
	}
}

func (d *InterfaceBridgesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *InterfaceBridgesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InterfaceBridgesModel

	bridges, err := d.client.GetBridges(ctx)
	if addError(&resp.Diagnostics, "Unable to get bridges", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *bridges)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
