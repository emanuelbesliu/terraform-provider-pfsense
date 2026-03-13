package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*GatewayGroupDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*GatewayGroupDataSource)(nil)
)

type GatewayGroupDataSourceModel struct {
	GatewayGroupModel
}

func NewGatewayGroupDataSource() datasource.DataSource { //nolint:ireturn
	return &GatewayGroupDataSource{}
}

type GatewayGroupDataSource struct {
	client *pfsense.Client
}

func (d *GatewayGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_gateway_group", req.ProviderTypeName)
}

func (d *GatewayGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single gateway group by name. Gateway groups organize multiple gateways into tiers for failover and load balancing.",
		MarkdownDescription: "Retrieves a single [gateway group](https://docs.netgate.com/pfsense/en/latest/routing/gateway-groups.html) by name. Gateway groups organize multiple gateways into tiers for failover and load balancing.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: GatewayGroupModel{}.descriptions()["name"].Description,
				Required:    true,
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
	}
}

func (d *GatewayGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *GatewayGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GatewayGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.client.GetGatewayGroup(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Unable to get gateway group", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *group)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
