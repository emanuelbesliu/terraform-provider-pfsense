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
	_ datasource.DataSource              = (*SystemTunablesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*SystemTunablesDataSource)(nil)
)

type SystemTunablesModel struct {
	SystemTunables types.List `tfsdk:"system_tunables"`
}

func NewSystemTunablesDataSource() datasource.DataSource { //nolint:ireturn
	return &SystemTunablesDataSource{}
}

type SystemTunablesDataSource struct {
	client *pfsense.Client
}

func (m *SystemTunablesModel) Set(ctx context.Context, tunables pfsense.SystemTunables) diag.Diagnostics {
	var diags diag.Diagnostics

	tunableModels := []SystemTunableModel{}
	for _, t := range tunables {
		var tunableModel SystemTunableModel
		diags.Append(tunableModel.Set(ctx, t)...)
		tunableModels = append(tunableModels, tunableModel)
	}

	tunablesValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: SystemTunableModel{}.AttrTypes()}, tunableModels)
	diags.Append(newDiags...)
	m.SystemTunables = tunablesValue

	return diags
}

func (d *SystemTunablesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_tunables", req.ProviderTypeName)
}

func (d *SystemTunablesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all system tunables (sysctls). Tunables allow adjusting FreeBSD kernel parameters that control networking, memory, and other system behaviors.",
		MarkdownDescription: "Retrieves all system [tunables](https://docs.netgate.com/pfsense/en/latest/system/advanced-tunables.html) (sysctls). Tunables allow adjusting FreeBSD kernel parameters that control networking, memory, and other system behaviors.",
		Attributes: map[string]schema.Attribute{
			"system_tunables": schema.ListNestedAttribute{
				Description: "List of system tunables.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tunable": schema.StringAttribute{
							Description: SystemTunableModel{}.descriptions()["tunable"].Description,
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: SystemTunableModel{}.descriptions()["value"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: SystemTunableModel{}.descriptions()["description"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *SystemTunablesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *SystemTunablesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemTunablesModel

	tunables, err := d.client.GetSystemTunables(ctx)
	if addError(&resp.Diagnostics, "Unable to get system tunables", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *tunables)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
