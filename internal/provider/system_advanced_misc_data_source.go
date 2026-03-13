package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ datasource.DataSource              = (*SystemAdvancedMiscDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*SystemAdvancedMiscDataSource)(nil)
)

type SystemAdvancedMiscDataSourceModel struct {
	SystemAdvancedMiscModel
}

func NewSystemAdvancedMiscDataSource() datasource.DataSource { //nolint:ireturn
	return &SystemAdvancedMiscDataSource{}
}

type SystemAdvancedMiscDataSource struct {
	client *pfsense.Client
}

func (d *SystemAdvancedMiscDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_advanced_misc", req.ProviderTypeName)
}

func (d *SystemAdvancedMiscDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	descriptions := SystemAdvancedMiscModel{}.descriptions()

	resp.Schema = schema.Schema{
		Description:         "Retrieves the system advanced miscellaneous configuration including proxy settings, load balancing, power savings, cryptographic hardware, security mitigations, gateway monitoring, and RAM disk settings.",
		MarkdownDescription: "Retrieves the [system advanced miscellaneous](https://docs.netgate.com/pfsense/en/latest/config/advanced/miscellaneous.html) configuration including proxy settings, load balancing, power savings, cryptographic hardware, security mitigations, gateway monitoring, and RAM disk settings.",
		Attributes: map[string]schema.Attribute{
			// Proxy Support
			"proxy_url": schema.StringAttribute{
				Description: descriptions["proxy_url"].Description,
				Computed:    true,
			},
			"proxy_port": schema.Int64Attribute{
				Description: descriptions["proxy_port"].Description,
				Computed:    true,
			},
			"proxy_user": schema.StringAttribute{
				Description: descriptions["proxy_user"].Description,
				Computed:    true,
			},
			"proxy_pass": schema.StringAttribute{
				Description: descriptions["proxy_pass"].Description,
				Computed:    true,
				Sensitive:   true,
			},

			// Load Balancing
			"lb_use_sticky": schema.BoolAttribute{
				Description: descriptions["lb_use_sticky"].Description,
				Computed:    true,
			},
			"src_track": schema.Int64Attribute{
				Description: descriptions["src_track"].Description,
				Computed:    true,
			},

			// Intel Speed Shift (hardware-dependent)
			"hwpstate": schema.StringAttribute{
				Description: descriptions["hwpstate"].Description,
				Computed:    true,
			},
			"hwpstate_control_level": schema.StringAttribute{
				Description: descriptions["hwpstate_control_level"].Description,
				Computed:    true,
			},
			"hwpstate_epp": schema.Int64Attribute{
				Description: descriptions["hwpstate_epp"].Description,
				Computed:    true,
			},

			// PowerD
			"powerd_enable": schema.BoolAttribute{
				Description: descriptions["powerd_enable"].Description,
				Computed:    true,
			},
			"powerd_ac_mode": schema.StringAttribute{
				Description:         descriptions["powerd_ac_mode"].Description,
				MarkdownDescription: descriptions["powerd_ac_mode"].MarkdownDescription,
				Computed:            true,
			},
			"powerd_battery_mode": schema.StringAttribute{
				Description:         descriptions["powerd_battery_mode"].Description,
				MarkdownDescription: descriptions["powerd_battery_mode"].MarkdownDescription,
				Computed:            true,
			},
			"powerd_normal_mode": schema.StringAttribute{
				Description:         descriptions["powerd_normal_mode"].Description,
				MarkdownDescription: descriptions["powerd_normal_mode"].MarkdownDescription,
				Computed:            true,
			},

			// Cryptographic & Thermal Hardware
			"crypto_hardware": schema.StringAttribute{
				Description:         descriptions["crypto_hardware"].Description,
				MarkdownDescription: descriptions["crypto_hardware"].MarkdownDescription,
				Computed:            true,
			},
			"thermal_hardware": schema.StringAttribute{
				Description:         descriptions["thermal_hardware"].Description,
				MarkdownDescription: descriptions["thermal_hardware"].MarkdownDescription,
				Computed:            true,
			},

			// Security Mitigations
			"pti_disabled": schema.BoolAttribute{
				Description: descriptions["pti_disabled"].Description,
				Computed:    true,
			},
			"mds_disable": schema.StringAttribute{
				Description:         descriptions["mds_disable"].Description,
				MarkdownDescription: descriptions["mds_disable"].MarkdownDescription,
				Computed:            true,
			},

			// Schedules
			"schedule_states": schema.BoolAttribute{
				Description: descriptions["schedule_states"].Description,
				Computed:    true,
			},

			// Gateway Monitoring
			"gw_down_kill_states": schema.StringAttribute{
				Description:         descriptions["gw_down_kill_states"].Description,
				MarkdownDescription: descriptions["gw_down_kill_states"].MarkdownDescription,
				Computed:            true,
			},
			"skip_rules_gw_down": schema.BoolAttribute{
				Description: descriptions["skip_rules_gw_down"].Description,
				Computed:    true,
			},
			"dpinger_dont_add_static_routes": schema.BoolAttribute{
				Description: descriptions["dpinger_dont_add_static_routes"].Description,
				Computed:    true,
			},

			// RAM Disk Settings
			"use_mfs_tmpvar": schema.BoolAttribute{
				Description: descriptions["use_mfs_tmpvar"].Description,
				Computed:    true,
			},
			"use_mfs_tmp_size": schema.Int64Attribute{
				Description: descriptions["use_mfs_tmp_size"].Description,
				Computed:    true,
			},
			"use_mfs_var_size": schema.Int64Attribute{
				Description: descriptions["use_mfs_var_size"].Description,
				Computed:    true,
			},
			"rrd_backup_interval": schema.Int64Attribute{
				Description: descriptions["rrd_backup_interval"].Description,
				Computed:    true,
			},
			"dhcp_backup_interval": schema.Int64Attribute{
				Description: descriptions["dhcp_backup_interval"].Description,
				Computed:    true,
			},
			"logs_backup_interval": schema.Int64Attribute{
				Description: descriptions["logs_backup_interval"].Description,
				Computed:    true,
			},
			"captive_portal_backup_interval": schema.Int64Attribute{
				Description: descriptions["captive_portal_backup_interval"].Description,
				Computed:    true,
			},

			// Hardware Settings
			"hard_disk_standby": schema.StringAttribute{
				Description:         descriptions["hard_disk_standby"].Description,
				MarkdownDescription: descriptions["hard_disk_standby"].MarkdownDescription,
				Computed:            true,
			},

			// PHP Settings
			"php_memory_limit": schema.Int64Attribute{
				Description: descriptions["php_memory_limit"].Description,
				Computed:    true,
			},

			// Installation Feedback
			"do_not_send_unique_id": schema.BoolAttribute{
				Description: descriptions["do_not_send_unique_id"].Description,
				Computed:    true,
			},
		},
	}
}

func (d *SystemAdvancedMiscDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	client, ok := configureDataSourceClient(req, resp)
	if !ok {
		return
	}

	d.client = client
}

func (d *SystemAdvancedMiscDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SystemAdvancedMiscDataSourceModel

	a, err := d.client.GetAdvancedMisc(ctx)
	if addError(&resp.Diagnostics, "Unable to get system advanced misc settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
