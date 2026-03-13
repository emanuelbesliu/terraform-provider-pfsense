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
	_ datasource.DataSource              = (*GatewayGroupsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*GatewayGroupsDataSource)(nil)
)

type GatewayGroupsModel struct {
	GatewayGroups types.List `tfsdk:"gateway_groups"`
}

func NewGatewayGroupsDataSource() datasource.DataSource { //nolint:ireturn
	return &GatewayGroupsDataSource{}
}

type GatewayGroupsDataSource struct {
	client *pfsense.Client
}

func (m *GatewayGroupsModel) Set(ctx context.Context, groups pfsense.GatewayGroups) diag.Diagnostics {
	var diags diag.Diagnostics

	groupModels := []GatewayGroupModel{}
	for _, g := range groups {
		var groupModel GatewayGroupModel
		diags.Append(groupModel.Set(ctx, g)...)
		groupModels = append(groupModels, groupModel)
	}

	groupsValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: GatewayGroupModel{}.AttrTypes()}, groupModels)
	diags.Append(newDiags...)
	m.GatewayGroups = groupsValue

	return diags
}

func (d *GatewayGroupsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_gateway_groups", req.ProviderTypeName)
}

func (d *GatewayGroupsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all gateway groups. Gateway groups organize multiple gateways into tiers for failover and load balancing.",
		MarkdownDescription: "Retrieves all [gateway groups](https://docs.netgate.com/pfsense/en/latest/routing/gateway-groups.html). Gateway groups organize multiple gateways into tiers for failover and load balancing.",
		Attributes: map[string]schema.Attribute{
			"gateway_groups": schema.ListNestedAttribute{
				Description: "List of gateway groups.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: GatewayGroupModel{}.descriptions()["name"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: GatewayGroupModel{}.descriptions()["description"].Description,
							Computed:    true,
						},
						"trigger": schema.StringAttribute{
							Description:         GatewayGroupModel{}.descriptions()["trigger"].Description,
							MarkdownDescription: GatewayGroupModel{}.descriptions()["trigger"].MarkdownDescription,
							Computed:            true,
						},
						"keep_failover_states": schema.StringAttribute{
							Description: GatewayGroupModel{}.descriptions()["keep_failover_states"].Description,
							Computed:    true,
						},
						"members": schema.ListNestedAttribute{
							Description: GatewayGroupModel{}.descriptions()["members"].Description,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"gateway": schema.StringAttribute{
										Description: GatewayGroupMemberModel{}.descriptions()["gateway"].Description,
										Computed:    true,
									},
									"tier": schema.Int64Attribute{
										Description: GatewayGroupMemberModel{}.descriptions()["tier"].Description,
										Computed:    true,
									},
									"virtual_ip": schema.StringAttribute{
										Description: GatewayGroupMemberModel{}.descriptions()["virtual_ip"].Description,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *GatewayGroupsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *GatewayGroupsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GatewayGroupsModel

	groups, err := d.client.GetGatewayGroups(ctx)
	if addError(&resp.Diagnostics, "Unable to get gateway groups", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *groups)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
