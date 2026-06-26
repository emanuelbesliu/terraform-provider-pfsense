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
	_ datasource.DataSource              = (*FirewallSchedulesDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallSchedulesDataSource)(nil)
)

type FirewallSchedulesModel struct {
	Schedules types.List `tfsdk:"schedules"`
}

func NewFirewallSchedulesDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallSchedulesDataSource{}
}

type FirewallSchedulesDataSource struct {
	client *pfsense.Client
}

func (m *FirewallSchedulesModel) Set(ctx context.Context, schedules pfsense.Schedules) diag.Diagnostics {
	var diags diag.Diagnostics

	models := []FirewallScheduleModel{}
	for _, s := range schedules {
		var model FirewallScheduleModel
		diags.Append(model.Set(ctx, s)...)
		models = append(models, model)
	}

	listValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: FirewallScheduleModel{}.AttrTypes()}, models)
	diags.Append(newDiags...)
	m.Schedules = listValue

	return diags
}

func (d *FirewallSchedulesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_schedules", req.ProviderTypeName)
}

func (d *FirewallSchedulesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves all firewall schedules.",
		MarkdownDescription: "Retrieves all [firewall schedules](https://docs.netgate.com/pfsense/en/latest/firewall/time-based-rules.html).",
		Attributes: map[string]schema.Attribute{
			"schedules": schema.ListNestedAttribute{
				Description: "List of firewall schedules.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: FirewallScheduleModel{}.descriptions()["name"].Description,
							Computed:    true,
						},
						"description": schema.StringAttribute{
							Description: FirewallScheduleModel{}.descriptions()["description"].Description,
							Computed:    true,
						},
						"time_range": schema.ListNestedAttribute{
							Description: FirewallScheduleModel{}.descriptions()["time_range"].Description,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: firewallScheduleTimeRangeDataSourceAttributes(),
							},
						},
						"label": schema.StringAttribute{
							Description: FirewallScheduleModel{}.descriptions()["label"].Description,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *FirewallSchedulesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallSchedulesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallSchedulesModel

	schedules, err := d.client.GetSchedules(ctx)
	if addError(&resp.Diagnostics, "Unable to get firewall schedules", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *schedules)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
