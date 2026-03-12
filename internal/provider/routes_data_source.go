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
	_ datasource.DataSource              = (*RoutesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*RoutesDataSource)(nil)
)

type RoutesModel struct {
	Routes types.List `tfsdk:"routes"`
}

func NewRoutesDataSource() datasource.DataSource { //nolint:ireturn
	return &RoutesDataSource{}
}

type RoutesDataSource struct {
	client *pfsense.Client
}

func (m *RoutesModel) Set(ctx context.Context, routes pfsense.Routes) diag.Diagnostics {
	var diags diag.Diagnostics

	routeModels := []RouteModel{}
	for _, route := range routes {
		var routeModel RouteModel
		diags.Append(routeModel.Set(ctx, route)...)
		routeModels = append(routeModels, routeModel)
	}

	routesValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: RouteModel{}.AttrTypes()}, routeModels)
	diags.Append(newDiags...)
	m.Routes = routesValue

	return diags
}

func (d *RoutesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_static_routes", req.ProviderTypeName)
}

func (d *RoutesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all static routes. Static routes direct traffic to a specific network through a gateway.",
		MarkdownDescription: "Retrieves all [static routes](https://docs.netgate.com/pfsense/en/latest/routing/static.html). Static routes direct traffic to a specific network through a gateway.",
		Attributes: map[string]schema.Attribute{
			"routes": schema.ListNestedAttribute{
				Description: "List of static routes.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"network": schema.StringAttribute{
							Description: RouteModel{}.descriptions()["network"].Description,
							Computed:    true,
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
				},
			},
		},
	}
}

func (d *RoutesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *RoutesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RoutesModel

	routes, err := d.client.GetRoutes(ctx)
	if addError(&resp.Diagnostics, "Unable to get static routes", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *routes)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
