package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*GatewayDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*GatewayDataSource)(nil)
)

type GatewayDataSourceModel struct {
	GatewayModel
}

func NewGatewayDataSource() datasource.DataSource { //nolint:ireturn
	return &GatewayDataSource{}
}

type GatewayDataSource struct {
	client *pfsense.Client
}

func (d *GatewayDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_gateway", req.ProviderTypeName)
}

func (d *GatewayDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single gateway by name. Gateways are used by static routes and can be organized into gateway groups for failover and load balancing.",
		MarkdownDescription: "Retrieves a single [gateway](https://docs.netgate.com/pfsense/en/latest/routing/gateways.html) by name. Gateways are used by static routes and can be organized into gateway groups for failover and load balancing.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: GatewayModel{}.descriptions()["name"].Description,
				Required:    true,
			},
			"interface": schema.StringAttribute{
				Description: GatewayModel{}.descriptions()["interface"].Description,
				Computed:    true,
			},
			"ipprotocol": schema.StringAttribute{
				Description:         GatewayModel{}.descriptions()["ipprotocol"].Description,
				MarkdownDescription: GatewayModel{}.descriptions()["ipprotocol"].MarkdownDescription,
				Computed:            true,
			},
			"gateway": schema.StringAttribute{
				Description: GatewayModel{}.descriptions()["gateway"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: GatewayModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: GatewayModel{}.descriptions()["disabled"].Description,
				Computed:    true,
			},
			"default_gw": schema.BoolAttribute{
				Description: GatewayModel{}.descriptions()["default_gw"].Description,
				Computed:    true,
			},
			"monitor": schema.StringAttribute{
				Description: GatewayModel{}.descriptions()["monitor"].Description,
				Computed:    true,
			},
			"monitor_disable": schema.BoolAttribute{
				Description: GatewayModel{}.descriptions()["monitor_disable"].Description,
				Computed:    true,
			},
			"action_disable": schema.BoolAttribute{
				Description: GatewayModel{}.descriptions()["action_disable"].Description,
				Computed:    true,
			},
			"force_down": schema.BoolAttribute{
				Description: GatewayModel{}.descriptions()["force_down"].Description,
				Computed:    true,
			},
			"weight": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["weight"].Description,
				Computed:    true,
			},
			"non_local_gateway": schema.BoolAttribute{
				Description: GatewayModel{}.descriptions()["non_local_gateway"].Description,
				Computed:    true,
			},
			"latency_low": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["latency_low"].Description,
				Computed:    true,
			},
			"latency_high": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["latency_high"].Description,
				Computed:    true,
			},
			"loss_low": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["loss_low"].Description,
				Computed:    true,
			},
			"loss_high": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["loss_high"].Description,
				Computed:    true,
			},
			"interval": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["interval"].Description,
				Computed:    true,
			},
			"loss_interval": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["loss_interval"].Description,
				Computed:    true,
			},
			"time_period": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["time_period"].Description,
				Computed:    true,
			},
			"alert_interval": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["alert_interval"].Description,
				Computed:    true,
			},
			"data_payload": schema.Int64Attribute{
				Description: GatewayModel{}.descriptions()["data_payload"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *GatewayDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *GatewayDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data GatewayDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	gw, err := d.client.GetGateway(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Unable to get gateway", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *gw)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
