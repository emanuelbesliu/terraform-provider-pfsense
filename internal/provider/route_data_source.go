package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*RouteDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*RouteDataSource)(nil)
)

type RouteDataSourceModel struct {
	RouteModel
}

func NewRouteDataSource() datasource.DataSource { //nolint:ireturn
	return &RouteDataSource{}
}

type RouteDataSource struct {
	client *pfsense.Client
}

func (d *RouteDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_static_route", req.ProviderTypeName)
}

func (d *RouteDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single static route by network. Static routes direct traffic to a specific network through a gateway.",
		MarkdownDescription: "Retrieves a single [static route](https://docs.netgate.com/pfsense/en/latest/routing/static.html) by network. Static routes direct traffic to a specific network through a gateway.",
		Attributes: map[string]schema.Attribute{
			"network": schema.StringAttribute{
				Description: RouteModel{}.descriptions()["network"].Description,
				Required:    true,
			},
			"gateway": schema.StringAttribute{
				Description: RouteModel{}.descriptions()["gateway"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: RouteModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: RouteModel{}.descriptions()["disabled"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *RouteDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *RouteDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RouteDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	routes, err := d.client.GetRoutes(ctx)
	if addError(&resp.Diagnostics, "Unable to get static routes", err) {
		return
	}

	route, err := routes.GetByNetwork(data.Network.ValueString())
	if addError(&resp.Diagnostics, "Unable to find static route", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *route)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
