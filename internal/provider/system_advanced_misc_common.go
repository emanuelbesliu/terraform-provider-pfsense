package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

// SystemAdvancedMiscModel represents the Terraform model for system advanced miscellaneous settings.
type SystemAdvancedMiscModel struct {
	// Proxy Support
	ProxyURL  types.String `tfsdk:"proxy_url"`
	ProxyPort types.Int64  `tfsdk:"proxy_port"`
	ProxyUser types.String `tfsdk:"proxy_user"`
	ProxyPass types.String `tfsdk:"proxy_pass"`

	// Load Balancing
	LBUseSticky types.Bool  `tfsdk:"lb_use_sticky"`
	SrcTrack    types.Int64 `tfsdk:"src_track"`

	// Power Savings - Intel Speed Shift (hardware-dependent)
	HWPState             types.String `tfsdk:"hwpstate"`
	HWPStateControlLevel types.String `tfsdk:"hwpstate_control_level"`
	HWPStateEPP          types.Int64  `tfsdk:"hwpstate_epp"`

	// Power Savings - PowerD
	PowerdEnable      types.Bool   `tfsdk:"powerd_enable"`
	PowerdACMode      types.String `tfsdk:"powerd_ac_mode"`
	PowerdBatteryMode types.String `tfsdk:"powerd_battery_mode"`
	PowerdNormalMode  types.String `tfsdk:"powerd_normal_mode"`

	// Cryptographic & Thermal Hardware
	CryptoHardware  types.String `tfsdk:"crypto_hardware"`
	ThermalHardware types.String `tfsdk:"thermal_hardware"`

	// Security Mitigations
	PTIDisabled types.Bool   `tfsdk:"pti_disabled"`
	MDSDisable  types.String `tfsdk:"mds_disable"`

	// Schedules
	ScheduleStates types.Bool `tfsdk:"schedule_states"`

	// Gateway Monitoring
	GWDownKillStates           types.String `tfsdk:"gw_down_kill_states"`
	SkipRulesGWDown            types.Bool   `tfsdk:"skip_rules_gw_down"`
	DPingerDontAddStaticRoutes types.Bool   `tfsdk:"dpinger_dont_add_static_routes"`

	// RAM Disk Settings
	UseMFSTmpVar        types.Bool  `tfsdk:"use_mfs_tmpvar"`
	UseMFSTmpSize       types.Int64 `tfsdk:"use_mfs_tmp_size"`
	UseMFSVarSize       types.Int64 `tfsdk:"use_mfs_var_size"`
	RRDBackup           types.Int64 `tfsdk:"rrd_backup_interval"`
	DHCPBackup          types.Int64 `tfsdk:"dhcp_backup_interval"`
	LogsBackup          types.Int64 `tfsdk:"logs_backup_interval"`
	CaptivePortalBackup types.Int64 `tfsdk:"captive_portal_backup_interval"`

	// Hardware Settings
	HardDiskStandby types.String `tfsdk:"hard_disk_standby"`

	// PHP Settings
	PHPMemoryLimit types.Int64 `tfsdk:"php_memory_limit"`

	// Installation Feedback
	DoNotSendUniqueID types.Bool `tfsdk:"do_not_send_unique_id"`
}

func (SystemAdvancedMiscModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"proxy_url": {
			Description: "Proxy URL or hostname for the system to use when making outbound HTTP/HTTPS requests (e.g., for package manager, update checks).",
		},
		"proxy_port": {
			Description: "Proxy port number. Only applicable when proxy_url is set.",
		},
		"proxy_user": {
			Description: "Proxy authentication username. Only applicable when proxy_url is set.",
		},
		"proxy_pass": {
			Description: "Proxy authentication password. Only applicable when proxy_url and proxy_user are set.",
		},
		"lb_use_sticky": {
			Description: "Use sticky connections for load balancer. Successive connections from the same source will be redirected to the same web server.",
		},
		"src_track": {
			Description: "Source tracking timeout in seconds for sticky connections. Set to 0 for the default. Only applicable when lb_use_sticky is true.",
		},
		"hwpstate": {
			Description: "Intel Speed Shift Technology (HWP) state. This field is hardware-dependent and may not be available on all systems. Values are 'enabled' or 'disabled', empty if unsupported.",
		},
		"hwpstate_control_level": {
			Description: "Intel Speed Shift control level. '0' for per-core control, '1' for per-package control. Only applicable when hwpstate is 'enabled'.",
		},
		"hwpstate_epp": {
			Description: "Intel Speed Shift Energy/Performance Preference (0-100). 0 is maximum performance, 100 is maximum energy savings. -1 indicates unsupported/not configured.",
		},
		"powerd_enable": {
			Description: "Enable the powerd utility to manage CPU frequency dynamically based on system load.",
		},
		"powerd_ac_mode": {
			Description:         fmt.Sprintf("PowerD mode when running on AC power. Options: %s. Defaults to '%s'.", wrapElementsJoin(pfsense.AdvancedMisc{}.PowerdModeOptions(), "'"), pfsense.DefaultAdvancedMiscPowerdMode),
			MarkdownDescription: fmt.Sprintf("PowerD mode when running on AC power. Options: %s. Defaults to `%s`.", wrapElementsJoin(pfsense.AdvancedMisc{}.PowerdModeOptions(), "`"), pfsense.DefaultAdvancedMiscPowerdMode),
		},
		"powerd_battery_mode": {
			Description:         fmt.Sprintf("PowerD mode when running on battery power. Options: %s. Defaults to '%s'.", wrapElementsJoin(pfsense.AdvancedMisc{}.PowerdModeOptions(), "'"), pfsense.DefaultAdvancedMiscPowerdMode),
			MarkdownDescription: fmt.Sprintf("PowerD mode when running on battery power. Options: %s. Defaults to `%s`.", wrapElementsJoin(pfsense.AdvancedMisc{}.PowerdModeOptions(), "`"), pfsense.DefaultAdvancedMiscPowerdMode),
		},
		"powerd_normal_mode": {
			Description:         fmt.Sprintf("PowerD mode when power source is unknown. Options: %s. Defaults to '%s'.", wrapElementsJoin(pfsense.AdvancedMisc{}.PowerdModeOptions(), "'"), pfsense.DefaultAdvancedMiscPowerdMode),
			MarkdownDescription: fmt.Sprintf("PowerD mode when power source is unknown. Options: %s. Defaults to `%s`.", wrapElementsJoin(pfsense.AdvancedMisc{}.PowerdModeOptions(), "`"), pfsense.DefaultAdvancedMiscPowerdMode),
		},
		"crypto_hardware": {
			Description:         fmt.Sprintf("Cryptographic hardware acceleration module. Options: %s. Empty string means no acceleration.", wrapElementsJoin(pfsense.AdvancedMisc{}.CryptoHardwareOptions(), "'")),
			MarkdownDescription: fmt.Sprintf("Cryptographic hardware acceleration module. Options: %s. Empty string means no acceleration.", wrapElementsJoin(pfsense.AdvancedMisc{}.CryptoHardwareOptions(), "`")),
		},
		"thermal_hardware": {
			Description:         fmt.Sprintf("Thermal sensor hardware module. Options: %s. Empty string means no thermal sensor.", wrapElementsJoin(pfsense.AdvancedMisc{}.ThermalHardwareOptions(), "'")),
			MarkdownDescription: fmt.Sprintf("Thermal sensor hardware module. Options: %s. Empty string means no thermal sensor.", wrapElementsJoin(pfsense.AdvancedMisc{}.ThermalHardwareOptions(), "`")),
		},
		"pti_disabled": {
			Description: "Disable kernel Page Table Isolation (PTI). PTI is a mitigation for Meltdown (CVE-2017-5754). Disabling may improve performance but reduces security.",
		},
		"mds_disable": {
			Description:         fmt.Sprintf("Microarchitectural Data Sampling (MDS) mitigation mode. Options: %s. Empty string uses system default.", wrapElementsJoin(pfsense.AdvancedMisc{}.MDSDisableOptions(), "'")),
			MarkdownDescription: fmt.Sprintf("Microarchitectural Data Sampling (MDS) mitigation mode. Options: %s. Empty string uses system default.", wrapElementsJoin(pfsense.AdvancedMisc{}.MDSDisableOptions(), "`")),
		},
		"schedule_states": {
			Description: "Do not kill connections when a schedule expires. By default, when a schedule expires, connections permitted by that schedule are killed.",
		},
		"gw_down_kill_states": {
			Description:         fmt.Sprintf("Flush states when a gateway goes down. Options: %s. Empty string means no action.", wrapElementsJoin(pfsense.AdvancedMisc{}.GWDownKillStatesOptions(), "'")),
			MarkdownDescription: fmt.Sprintf("Flush states when a gateway goes down. Options: %s. Empty string means no action.", wrapElementsJoin(pfsense.AdvancedMisc{}.GWDownKillStatesOptions(), "`")),
		},
		"skip_rules_gw_down": {
			Description: "Do not create rules when gateway is down. By default, when a rule has a specific gateway, the rule is created and traffic is sent to the default gateway. This option overrides that behavior.",
		},
		"dpinger_dont_add_static_routes": {
			Description: "Do not add static routes for gateway monitor addresses. By default, static routes are added to ensure gateway monitoring traffic uses the correct interface.",
		},
		"use_mfs_tmpvar": {
			Description: "Use RAM disks for /tmp and /var. Useful for embedded systems or to reduce disk writes.",
		},
		"use_mfs_tmp_size": {
			Description: "Size of the /tmp RAM disk in MiB. Set to 0 for the default (40 MiB). Only applicable when use_mfs_tmpvar is true.",
		},
		"use_mfs_var_size": {
			Description: "Size of the /var RAM disk in MiB. Set to 0 for the default (60 MiB). Only applicable when use_mfs_tmpvar is true.",
		},
		"rrd_backup_interval": {
			Description: "RRD data backup interval in hours (0-24). Set to 0 to disable periodic RRD backup. Only applicable when use_mfs_tmpvar is true.",
		},
		"dhcp_backup_interval": {
			Description: "DHCP leases backup interval in hours (0-24). Set to 0 to disable periodic DHCP backup. Only applicable when use_mfs_tmpvar is true.",
		},
		"logs_backup_interval": {
			Description: "Log files backup interval in hours (0-24). Set to 0 to disable periodic log backup. Only applicable when use_mfs_tmpvar is true.",
		},
		"captive_portal_backup_interval": {
			Description: "Captive portal data backup interval in hours (0-24). Set to 0 to disable periodic captive portal backup. Only applicable when use_mfs_tmpvar is true.",
		},
		"hard_disk_standby": {
			Description:         fmt.Sprintf("Hard disk standby time. Options: %s. Empty string means 'Always On' (no standby).", wrapElementsJoin(pfsense.AdvancedMisc{}.HardDiskStandbyOptions(), "'")),
			MarkdownDescription: fmt.Sprintf("Hard disk standby time. Options: %s. Empty string means 'Always On' (no standby).", wrapElementsJoin(pfsense.AdvancedMisc{}.HardDiskStandbyOptions(), "`")),
		},
		"php_memory_limit": {
			Description: "PHP memory limit in MiB. Set to 0 for the system default. Higher values may be needed for large configurations.",
		},
		"do_not_send_unique_id": {
			Description: "Do not send the pfSense unique identifier with crash reports and update checks.",
		},
	}
}

func (SystemAdvancedMiscModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"proxy_url":                      types.StringType,
		"proxy_port":                     types.Int64Type,
		"proxy_user":                     types.StringType,
		"proxy_pass":                     types.StringType,
		"lb_use_sticky":                  types.BoolType,
		"src_track":                      types.Int64Type,
		"hwpstate":                       types.StringType,
		"hwpstate_control_level":         types.StringType,
		"hwpstate_epp":                   types.Int64Type,
		"powerd_enable":                  types.BoolType,
		"powerd_ac_mode":                 types.StringType,
		"powerd_battery_mode":            types.StringType,
		"powerd_normal_mode":             types.StringType,
		"crypto_hardware":                types.StringType,
		"thermal_hardware":               types.StringType,
		"pti_disabled":                   types.BoolType,
		"mds_disable":                    types.StringType,
		"schedule_states":                types.BoolType,
		"gw_down_kill_states":            types.StringType,
		"skip_rules_gw_down":             types.BoolType,
		"dpinger_dont_add_static_routes": types.BoolType,
		"use_mfs_tmpvar":                 types.BoolType,
		"use_mfs_tmp_size":               types.Int64Type,
		"use_mfs_var_size":               types.Int64Type,
		"rrd_backup_interval":            types.Int64Type,
		"dhcp_backup_interval":           types.Int64Type,
		"logs_backup_interval":           types.Int64Type,
		"captive_portal_backup_interval": types.Int64Type,
		"hard_disk_standby":              types.StringType,
		"php_memory_limit":               types.Int64Type,
		"do_not_send_unique_id":          types.BoolType,
	}
}

func (m *SystemAdvancedMiscModel) Set(_ context.Context, a pfsense.AdvancedMisc) diag.Diagnostics {
	var diags diag.Diagnostics

	// Proxy Support
	if a.ProxyURL != "" {
		m.ProxyURL = types.StringValue(a.ProxyURL)
	} else {
		m.ProxyURL = types.StringNull()
	}

	m.ProxyPort = types.Int64Value(int64(a.ProxyPort))

	if a.ProxyUser != "" {
		m.ProxyUser = types.StringValue(a.ProxyUser)
	} else {
		m.ProxyUser = types.StringNull()
	}

	if a.ProxyPass != "" {
		m.ProxyPass = types.StringValue(a.ProxyPass)
	} else {
		m.ProxyPass = types.StringNull()
	}

	// Load Balancing
	m.LBUseSticky = types.BoolValue(a.LBUseSticky)
	m.SrcTrack = types.Int64Value(int64(a.SrcTrack))

	// Intel Speed Shift (hardware-dependent)
	if a.HWPState != "" {
		m.HWPState = types.StringValue(a.HWPState)
	} else {
		m.HWPState = types.StringNull()
	}

	if a.HWPStateControlLevel != "" {
		m.HWPStateControlLevel = types.StringValue(a.HWPStateControlLevel)
	} else {
		m.HWPStateControlLevel = types.StringNull()
	}

	m.HWPStateEPP = types.Int64Value(int64(a.HWPStateEPP))

	// PowerD
	m.PowerdEnable = types.BoolValue(a.PowerdEnable)
	m.PowerdACMode = types.StringValue(a.PowerdACMode)
	m.PowerdBatteryMode = types.StringValue(a.PowerdBatteryMode)
	m.PowerdNormalMode = types.StringValue(a.PowerdNormalMode)

	// Cryptographic & Thermal Hardware
	m.CryptoHardware = types.StringValue(a.CryptoHardware)
	m.ThermalHardware = types.StringValue(a.ThermalHardware)

	// Security Mitigations
	m.PTIDisabled = types.BoolValue(a.PTIDisabled)
	m.MDSDisable = types.StringValue(a.MDSDisable)

	// Schedules
	m.ScheduleStates = types.BoolValue(a.ScheduleStates)

	// Gateway Monitoring
	m.GWDownKillStates = types.StringValue(a.GWDownKillStates)
	m.SkipRulesGWDown = types.BoolValue(a.SkipRulesGWDown)
	m.DPingerDontAddStaticRoutes = types.BoolValue(a.DPingerDontAddStaticRoutes)

	// RAM Disk Settings
	m.UseMFSTmpVar = types.BoolValue(a.UseMFSTmpVar)
	m.UseMFSTmpSize = types.Int64Value(int64(a.UseMFSTmpSize))
	m.UseMFSVarSize = types.Int64Value(int64(a.UseMFSVarSize))
	m.RRDBackup = types.Int64Value(int64(a.RRDBackup))
	m.DHCPBackup = types.Int64Value(int64(a.DHCPBackup))
	m.LogsBackup = types.Int64Value(int64(a.LogsBackup))
	m.CaptivePortalBackup = types.Int64Value(int64(a.CaptivePortalBackup))

	// Hardware Settings
	m.HardDiskStandby = types.StringValue(a.HardDiskStandby)

	// PHP Settings
	m.PHPMemoryLimit = types.Int64Value(int64(a.PHPMemoryLimit))

	// Installation Feedback
	m.DoNotSendUniqueID = types.BoolValue(a.DoNotSendUniqueID)

	return diags
}

func (m SystemAdvancedMiscModel) Value(_ context.Context, a *pfsense.AdvancedMisc) diag.Diagnostics {
	var diags diag.Diagnostics

	// Proxy Support
	if !m.ProxyURL.IsNull() {
		a.ProxyURL = m.ProxyURL.ValueString()
	}

	a.ProxyPort = int(m.ProxyPort.ValueInt64())

	if !m.ProxyUser.IsNull() {
		a.ProxyUser = m.ProxyUser.ValueString()
	}

	if !m.ProxyPass.IsNull() {
		a.ProxyPass = m.ProxyPass.ValueString()
	}

	// Load Balancing
	a.LBUseSticky = m.LBUseSticky.ValueBool()
	a.SrcTrack = int(m.SrcTrack.ValueInt64())

	// Intel Speed Shift (hardware-dependent)
	if !m.HWPState.IsNull() {
		a.HWPState = m.HWPState.ValueString()
	}

	if !m.HWPStateControlLevel.IsNull() {
		a.HWPStateControlLevel = m.HWPStateControlLevel.ValueString()
	}

	a.HWPStateEPP = int(m.HWPStateEPP.ValueInt64())

	// PowerD
	a.PowerdEnable = m.PowerdEnable.ValueBool()

	addPathError(
		&diags,
		path.Root("powerd_ac_mode"),
		"PowerD AC mode cannot be parsed",
		a.SetPowerdACMode(m.PowerdACMode.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("powerd_battery_mode"),
		"PowerD battery mode cannot be parsed",
		a.SetPowerdBatteryMode(m.PowerdBatteryMode.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("powerd_normal_mode"),
		"PowerD normal mode cannot be parsed",
		a.SetPowerdNormalMode(m.PowerdNormalMode.ValueString()),
	)

	// Cryptographic & Thermal Hardware
	addPathError(
		&diags,
		path.Root("crypto_hardware"),
		"Crypto hardware cannot be parsed",
		a.SetCryptoHardware(m.CryptoHardware.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("thermal_hardware"),
		"Thermal hardware cannot be parsed",
		a.SetThermalHardware(m.ThermalHardware.ValueString()),
	)

	// Security Mitigations
	a.PTIDisabled = m.PTIDisabled.ValueBool()

	addPathError(
		&diags,
		path.Root("mds_disable"),
		"MDS disable mode cannot be parsed",
		a.SetMDSDisable(m.MDSDisable.ValueString()),
	)

	// Schedules
	a.ScheduleStates = m.ScheduleStates.ValueBool()

	// Gateway Monitoring
	addPathError(
		&diags,
		path.Root("gw_down_kill_states"),
		"Gateway down kill states cannot be parsed",
		a.SetGWDownKillStates(m.GWDownKillStates.ValueString()),
	)

	a.SkipRulesGWDown = m.SkipRulesGWDown.ValueBool()
	a.DPingerDontAddStaticRoutes = m.DPingerDontAddStaticRoutes.ValueBool()

	// RAM Disk Settings
	a.UseMFSTmpVar = m.UseMFSTmpVar.ValueBool()
	a.UseMFSTmpSize = int(m.UseMFSTmpSize.ValueInt64())
	a.UseMFSVarSize = int(m.UseMFSVarSize.ValueInt64())
	a.RRDBackup = int(m.RRDBackup.ValueInt64())
	a.DHCPBackup = int(m.DHCPBackup.ValueInt64())
	a.LogsBackup = int(m.LogsBackup.ValueInt64())
	a.CaptivePortalBackup = int(m.CaptivePortalBackup.ValueInt64())

	// Hardware Settings
	addPathError(
		&diags,
		path.Root("hard_disk_standby"),
		"Hard disk standby cannot be parsed",
		a.SetHardDiskStandby(m.HardDiskStandby.ValueString()),
	)

	// PHP Settings
	a.PHPMemoryLimit = int(m.PHPMemoryLimit.ValueInt64())

	// Installation Feedback
	a.DoNotSendUniqueID = m.DoNotSendUniqueID.ValueBool()

	return diags
}
