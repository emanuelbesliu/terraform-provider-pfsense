package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*SystemGeneralDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*SystemGeneralDataSource)(nil)
)

type SystemGeneralDataSourceModel struct {
	SystemGeneralModel
}

func NewSystemGeneralDataSource() datasource.DataSource { //nolint:ireturn
	return &SystemGeneralDataSource{}
}

type SystemGeneralDataSource struct {
	client *pfsense.Client
}

func (d *SystemGeneralDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_general", req.ProviderTypeName)
}

func (d *SystemGeneralDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Retrieves the system general setup configuration including hostname, domain, DNS servers, localization, and webConfigurator settings.",
		MarkdownDescription: "Retrieves the [system general setup](https://docs.netgate.com/pfsense/en/latest/config/general.html) configuration including hostname, domain, DNS servers, localization, and webConfigurator settings.",
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["hostname"].Description,
				Computed:    true,
			},
			"domain": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["domain"].Description,
				Computed:    true,
			},
			"dns_servers": schema.ListNestedAttribute{
				Description: SystemGeneralModel{}.descriptions()["dns_servers"].Description,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							Description: SystemGeneralDNSServerModel{}.descriptions()["address"].Description,
							Computed:    true,
						},
						"hostname": schema.StringAttribute{
							Description: SystemGeneralDNSServerModel{}.descriptions()["hostname"].Description,
							Computed:    true,
						},
						"gateway": schema.StringAttribute{
							Description: SystemGeneralDNSServerModel{}.descriptions()["gateway"].Description,
							Computed:    true,
						},
					},
				},
			},
			"dns_override": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["dns_override"].Description,
				Computed:    true,
			},
			"dns_localhost": schema.StringAttribute{
				Description:         SystemGeneralModel{}.descriptions()["dns_localhost"].Description,
				MarkdownDescription: SystemGeneralModel{}.descriptions()["dns_localhost"].MarkdownDescription,
				Computed:            true,
			},
			"timezone": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["timezone"].Description,
				Computed:    true,
			},
			"timeservers": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["timeservers"].Description,
				Computed:    true,
			},
			"language": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["language"].Description,
				Computed:    true,
			},
			"webgui_theme": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["webgui_theme"].Description,
				Computed:    true,
			},
			"login_color": schema.StringAttribute{
				Description: SystemGeneralModel{}.descriptions()["login_color"].Description,
				Computed:    true,
			},
			"login_show_host": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["login_show_host"].Description,
				Computed:    true,
			},
			"webgui_fixed_menu": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["webgui_fixed_menu"].Description,
				Computed:    true,
			},
			"dashboard_columns": schema.Int64Attribute{
				Description: SystemGeneralModel{}.descriptions()["dashboard_columns"].Description,
				Computed:    true,
			},
			"webgui_left_column_hyper": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["webgui_left_column_hyper"].Description,
				Computed:    true,
			},
			"disable_alias_popup_detail": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["disable_alias_popup_detail"].Description,
				Computed:    true,
			},
			"dashboard_available_widgets_panel": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["dashboard_available_widgets_panel"].Description,
				Computed:    true,
			},
			"system_logs_filter_panel": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["system_logs_filter_panel"].Description,
				Computed:    true,
			},
			"system_logs_manage_log_panel": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["system_logs_manage_log_panel"].Description,
				Computed:    true,
			},
			"status_monitoring_settings_panel": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["status_monitoring_settings_panel"].Description,
				Computed:    true,
			},
			"row_order_dragging": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["row_order_dragging"].Description,
				Computed:    true,
			},
			"interfaces_sort": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["interfaces_sort"].Description,
				Computed:    true,
			},
			"require_state_filter": schema.BoolAttribute{
				Description: SystemGeneralModel{}.descriptions()["require_state_filter"].Description,
				Computed:    true,
			},
			"hostname_in_menu": schema.StringAttribute{
				Description:         SystemGeneralModel{}.descriptions()["hostname_in_menu"].Description,
				MarkdownDescription: SystemGeneralModel{}.descriptions()["hostname_in_menu"].MarkdownDescription,
				Computed:            true,
			},
		},
	}
}

func (d *SystemGeneralDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *SystemGeneralDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemGeneralDataSourceModel

	sg, err := d.client.GetSystemGeneral(ctx)
	if addError(&resp.Diagnostics, "Unable to get system general settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *sg)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
