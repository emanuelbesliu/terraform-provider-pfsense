package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*SystemAdvancedMiscResource)(nil)
	_ resource.ResourceWithConfigure   = (*SystemAdvancedMiscResource)(nil)
	_ resource.ResourceWithImportState = (*SystemAdvancedMiscResource)(nil)
)

type SystemAdvancedMiscResourceModel struct {
	SystemAdvancedMiscModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewSystemAdvancedMiscResource() resource.Resource { //nolint:ireturn
	return &SystemAdvancedMiscResource{}
}

type SystemAdvancedMiscResource struct {
	client *pfsense.Client
}

func (r *SystemAdvancedMiscResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_advanced_misc", req.ProviderTypeName)
}

func (r *SystemAdvancedMiscResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	descriptions := SystemAdvancedMiscModel{}.descriptions()

	resp.Schema = schema.Schema{
		Description:         "System advanced miscellaneous configuration including proxy settings, load balancing, power savings, cryptographic hardware, security mitigations, gateway monitoring, RAM disk settings, and more. This is a singleton resource — only one instance can exist per pfSense installation.",
		MarkdownDescription: "[System advanced miscellaneous](https://docs.netgate.com/pfsense/en/latest/config/advanced/miscellaneous.html) configuration including proxy settings, load balancing, power savings, cryptographic hardware, security mitigations, gateway monitoring, RAM disk settings, and more. This is a **singleton** resource — only one instance can exist per pfSense installation.",
		Attributes: map[string]schema.Attribute{
			// Proxy Support
			"proxy_url": schema.StringAttribute{
				Description: descriptions["proxy_url"].Description,
				Optional:    true,
			},
			"proxy_port": schema.Int64Attribute{
				Description: descriptions["proxy_port"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},
			"proxy_user": schema.StringAttribute{
				Description: descriptions["proxy_user"].Description,
				Optional:    true,
			},
			"proxy_pass": schema.StringAttribute{
				Description: descriptions["proxy_pass"].Description,
				Optional:    true,
				Sensitive:   true,
			},

			// Load Balancing
			"lb_use_sticky": schema.BoolAttribute{
				Description: descriptions["lb_use_sticky"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"src_track": schema.Int64Attribute{
				Description: descriptions["src_track"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},

			// Intel Speed Shift (hardware-dependent)
			"hwpstate": schema.StringAttribute{
				Description: descriptions["hwpstate"].Description,
				Optional:    true,
			},
			"hwpstate_control_level": schema.StringAttribute{
				Description: descriptions["hwpstate_control_level"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("0", "1"),
				},
			},
			"hwpstate_epp": schema.Int64Attribute{
				Description: descriptions["hwpstate_epp"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(-1),
				Validators: []validator.Int64{
					int64validator.Between(-1, 100),
				},
			},

			// PowerD
			"powerd_enable": schema.BoolAttribute{
				Description: descriptions["powerd_enable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"powerd_ac_mode": schema.StringAttribute{
				Description:         descriptions["powerd_ac_mode"].Description,
				MarkdownDescription: descriptions["powerd_ac_mode"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedMiscPowerdMode),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedMisc{}.PowerdModeOptions()...),
				},
			},
			"powerd_battery_mode": schema.StringAttribute{
				Description:         descriptions["powerd_battery_mode"].Description,
				MarkdownDescription: descriptions["powerd_battery_mode"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedMiscPowerdMode),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedMisc{}.PowerdModeOptions()...),
				},
			},
			"powerd_normal_mode": schema.StringAttribute{
				Description:         descriptions["powerd_normal_mode"].Description,
				MarkdownDescription: descriptions["powerd_normal_mode"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedMiscPowerdMode),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedMisc{}.PowerdModeOptions()...),
				},
			},

			// Cryptographic & Thermal Hardware
			"crypto_hardware": schema.StringAttribute{
				Description:         descriptions["crypto_hardware"].Description,
				MarkdownDescription: descriptions["crypto_hardware"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedMisc{}.CryptoHardwareOptions()...),
				},
			},
			"thermal_hardware": schema.StringAttribute{
				Description:         descriptions["thermal_hardware"].Description,
				MarkdownDescription: descriptions["thermal_hardware"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedMisc{}.ThermalHardwareOptions()...),
				},
			},

			// Security Mitigations
			"pti_disabled": schema.BoolAttribute{
				Description: descriptions["pti_disabled"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"mds_disable": schema.StringAttribute{
				Description:         descriptions["mds_disable"].Description,
				MarkdownDescription: descriptions["mds_disable"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedMisc{}.MDSDisableOptions()...),
				},
			},

			// Schedules
			"schedule_states": schema.BoolAttribute{
				Description: descriptions["schedule_states"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},

			// Gateway Monitoring
			"gw_down_kill_states": schema.StringAttribute{
				Description:         descriptions["gw_down_kill_states"].Description,
				MarkdownDescription: descriptions["gw_down_kill_states"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedMisc{}.GWDownKillStatesOptions()...),
				},
			},
			"skip_rules_gw_down": schema.BoolAttribute{
				Description: descriptions["skip_rules_gw_down"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"dpinger_dont_add_static_routes": schema.BoolAttribute{
				Description: descriptions["dpinger_dont_add_static_routes"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},

			// RAM Disk Settings
			"use_mfs_tmpvar": schema.BoolAttribute{
				Description: descriptions["use_mfs_tmpvar"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"use_mfs_tmp_size": schema.Int64Attribute{
				Description: descriptions["use_mfs_tmp_size"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},
			"use_mfs_var_size": schema.Int64Attribute{
				Description: descriptions["use_mfs_var_size"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},
			"rrd_backup_interval": schema.Int64Attribute{
				Description: descriptions["rrd_backup_interval"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 24),
				},
			},
			"dhcp_backup_interval": schema.Int64Attribute{
				Description: descriptions["dhcp_backup_interval"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 24),
				},
			},
			"logs_backup_interval": schema.Int64Attribute{
				Description: descriptions["logs_backup_interval"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 24),
				},
			},
			"captive_portal_backup_interval": schema.Int64Attribute{
				Description: descriptions["captive_portal_backup_interval"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
				Validators: []validator.Int64{
					int64validator.Between(0, 24),
				},
			},

			// Hardware Settings
			"hard_disk_standby": schema.StringAttribute{
				Description:         descriptions["hard_disk_standby"].Description,
				MarkdownDescription: descriptions["hard_disk_standby"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(""),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedMisc{}.HardDiskStandbyOptions()...),
				},
			},

			// PHP Settings
			"php_memory_limit": schema.Int64Attribute{
				Description: descriptions["php_memory_limit"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(0),
			},

			// Installation Feedback
			"do_not_send_unique_id": schema.BoolAttribute{
				Description: descriptions["do_not_send_unique_id"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},

			// Apply
			"apply": schema.BoolAttribute{
				Description:         applyDescription,
				MarkdownDescription: applyMarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(defaultApply),
			},
		},
	}
}

func (r *SystemAdvancedMiscResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *SystemAdvancedMiscResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SystemAdvancedMiscResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var aReq pfsense.AdvancedMisc
	resp.Diagnostics.Append(data.Value(ctx, &aReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.UpdateAdvancedMisc(ctx, aReq)
	if addError(&resp.Diagnostics, "Error creating system advanced misc settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedMiscChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced misc changes", err)
	}
}

func (r *SystemAdvancedMiscResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SystemAdvancedMiscResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.GetAdvancedMisc(ctx)
	if addError(&resp.Diagnostics, "Error reading system advanced misc settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemAdvancedMiscResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SystemAdvancedMiscResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var aReq pfsense.AdvancedMisc
	resp.Diagnostics.Append(data.Value(ctx, &aReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.UpdateAdvancedMisc(ctx, aReq)
	if addError(&resp.Diagnostics, "Error updating system advanced misc settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedMiscChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced misc changes", err)
	}
}

func (r *SystemAdvancedMiscResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SystemAdvancedMiscResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Reset to pfSense defaults
	defaultMisc := pfsense.AdvancedMisc{
		PowerdACMode:      pfsense.DefaultAdvancedMiscPowerdMode,
		PowerdBatteryMode: pfsense.DefaultAdvancedMiscPowerdMode,
		PowerdNormalMode:  pfsense.DefaultAdvancedMiscPowerdMode,
		HWPStateEPP:       -1,
	}

	_, err := r.client.UpdateAdvancedMisc(ctx, defaultMisc)
	if addError(&resp.Diagnostics, "Error resetting system advanced misc settings to defaults", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedMiscChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced misc changes", err)
	}
}

func (r *SystemAdvancedMiscResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	a, err := r.client.GetAdvancedMisc(ctx)
	if addError(&resp.Diagnostics, "Error importing system advanced misc settings", err) {
		return
	}

	var data SystemAdvancedMiscResourceModel
	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Apply = types.BoolValue(defaultApply)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
