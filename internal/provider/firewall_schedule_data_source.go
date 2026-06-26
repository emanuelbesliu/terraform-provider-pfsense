package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*FirewallScheduleDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*FirewallScheduleDataSource)(nil)
)

type FirewallScheduleDataSourceModel struct {
	FirewallScheduleModel
}

func NewFirewallScheduleDataSource() datasource.DataSource { //nolint:ireturn
	return &FirewallScheduleDataSource{}
}

type FirewallScheduleDataSource struct {
	client *pfsense.Client
}

func firewallScheduleTimeRangeDataSourceAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"position": schema.StringAttribute{
			Description: FirewallScheduleTimeRangeModel{}.descriptions()["position"].Description,
			Computed:    true,
		},
		"month": schema.StringAttribute{
			Description: FirewallScheduleTimeRangeModel{}.descriptions()["month"].Description,
			Computed:    true,
		},
		"day": schema.StringAttribute{
			Description: FirewallScheduleTimeRangeModel{}.descriptions()["day"].Description,
			Computed:    true,
		},
		"start_time": schema.StringAttribute{
			Description: FirewallScheduleTimeRangeModel{}.descriptions()["start_time"].Description,
			Computed:    true,
		},
		"stop_time": schema.StringAttribute{
			Description: FirewallScheduleTimeRangeModel{}.descriptions()["stop_time"].Description,
			Computed:    true,
		},
		"range_description": schema.StringAttribute{
			Description: FirewallScheduleTimeRangeModel{}.descriptions()["range_description"].Description,
			Computed:    true,
		},
	}
}

func (d *FirewallScheduleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_firewall_schedule", req.ProviderTypeName)
}

func (d *FirewallScheduleDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves a single firewall schedule by its name.",
		MarkdownDescription: "Retrieves a single [firewall schedule](https://docs.netgate.com/pfsense/en/latest/firewall/time-based-rules.html) by its name.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: FirewallScheduleModel{}.descriptions()["name"].Description,
				Required:    true,
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
	}
}

func (d *FirewallScheduleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *FirewallScheduleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data FirewallScheduleDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	schedule, err := d.client.GetSchedule(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Unable to get firewall schedule", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *schedule)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
