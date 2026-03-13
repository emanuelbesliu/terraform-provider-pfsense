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

// SystemAdvancedAdminModel represents the Terraform model for system advanced admin access settings.
type SystemAdvancedAdminModel struct {
	// webConfigurator
	WebGUIProto         types.String `tfsdk:"webgui_protocol"`
	SSLCertRef          types.String `tfsdk:"ssl_certificate"`
	WebGUIPort          types.Int64  `tfsdk:"webgui_port"`
	MaxProcs            types.Int64  `tfsdk:"max_processes"`
	DisableHTTPRedirect types.Bool   `tfsdk:"disable_http_redirect"`
	DisableHSTS         types.Bool   `tfsdk:"disable_hsts"`
	OCSPStaple          types.Bool   `tfsdk:"ocsp_staple"`
	LoginAutocomplete   types.Bool   `tfsdk:"login_autocomplete"`
	QuietLogin          types.Bool   `tfsdk:"quiet_login"`
	Roaming             types.Bool   `tfsdk:"roaming"`
	NoAntiLockout       types.Bool   `tfsdk:"disable_anti_lockout"`
	NoDNSRebindCheck    types.Bool   `tfsdk:"disable_dns_rebind_check"`
	NoHTTPRefererCheck  types.Bool   `tfsdk:"disable_http_referer_check"`
	AlternateHostnames  types.String `tfsdk:"alternate_hostnames"`
	PageNameFirst       types.Bool   `tfsdk:"page_name_first"`

	// Secure Shell
	SSHEnabled          types.Bool   `tfsdk:"ssh_enabled"`
	SSHdKeyOnly         types.String `tfsdk:"sshd_key_only"`
	SSHdAgentForwarding types.Bool   `tfsdk:"sshd_agent_forwarding"`
	SSHPort             types.Int64  `tfsdk:"ssh_port"`

	// Login Protection
	SshguardThreshold     types.Int64  `tfsdk:"login_protection_threshold"`
	SshguardBlocktime     types.Int64  `tfsdk:"login_protection_blocktime"`
	SshguardDetectionTime types.Int64  `tfsdk:"login_protection_detection_time"`
	SshguardWhitelist     types.String `tfsdk:"login_protection_pass_list"`

	// Serial Communications
	EnableSerial   types.Bool   `tfsdk:"serial_terminal"`
	SerialSpeed    types.Int64  `tfsdk:"serial_speed"`
	PrimaryConsole types.String `tfsdk:"primary_console"`

	// Console Options
	DisableConsoleMenu types.Bool `tfsdk:"disable_console_menu"`
}

func (SystemAdvancedAdminModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"webgui_protocol": {
			Description:         fmt.Sprintf("WebGUI protocol. Options: %s. Defaults to '%s'.", wrapElementsJoin(pfsense.AdvancedAdmin{}.WebGUIProtoOptions(), "'"), pfsense.DefaultAdvancedAdminWebGUIProto),
			MarkdownDescription: fmt.Sprintf("WebGUI protocol. Options: %s. Defaults to `%s`.", wrapElementsJoin(pfsense.AdvancedAdmin{}.WebGUIProtoOptions(), "`"), pfsense.DefaultAdvancedAdminWebGUIProto),
		},
		"ssl_certificate": {
			Description: "SSL certificate reference ID for HTTPS. Only applicable when webgui_protocol is 'https'.",
		},
		"webgui_port": {
			Description: "WebGUI TCP port. Set to 0 for the default port (80 for HTTP, 443 for HTTPS).",
		},
		"max_processes": {
			Description: fmt.Sprintf("Maximum number of webConfigurator processes. Defaults to '%d'.", pfsense.DefaultAdvancedAdminMaxProcs),
		},
		"disable_http_redirect": {
			Description: "Disable the HTTP to HTTPS redirect on the webConfigurator when HTTPS is enabled.",
		},
		"disable_hsts": {
			Description: "Disable HTTP Strict Transport Security (HSTS) on the webConfigurator.",
		},
		"ocsp_staple": {
			Description: "Enable OCSP stapling on the webConfigurator.",
		},
		"login_autocomplete": {
			Description: "Enable login autocomplete on the webConfigurator login page.",
		},
		"quiet_login": {
			Description: "Disable logging of webConfigurator successful logins.",
		},
		"roaming": {
			Description: "Allow webConfigurator session roaming between IP addresses. Defaults to 'true'.",
		},
		"disable_anti_lockout": {
			Description: "Disable the webConfigurator anti-lockout rule.",
		},
		"disable_dns_rebind_check": {
			Description: "Disable DNS rebinding attack check on the webConfigurator.",
		},
		"disable_http_referer_check": {
			Description: "Disable HTTP_REFERER enforcement check on the webConfigurator.",
		},
		"alternate_hostnames": {
			Description: "Space-separated list of alternate hostnames for the webConfigurator DNS rebind and HTTP_REFERER checks.",
		},
		"page_name_first": {
			Description: "Display page name first in the browser tab title.",
		},
		"ssh_enabled": {
			Description: "Enable Secure Shell (SSH) server.",
		},
		"sshd_key_only": {
			Description:         fmt.Sprintf("SSH authentication method. Options: %s. Defaults to '%s'.", wrapElementsJoin(pfsense.AdvancedAdmin{}.SSHdKeyOnlyOptions(), "'"), pfsense.DefaultAdvancedAdminSSHdKeyOnly),
			MarkdownDescription: fmt.Sprintf("SSH authentication method. Options: %s. Defaults to `%s`.", wrapElementsJoin(pfsense.AdvancedAdmin{}.SSHdKeyOnlyOptions(), "`"), pfsense.DefaultAdvancedAdminSSHdKeyOnly),
		},
		"sshd_agent_forwarding": {
			Description: "Enable SSH agent forwarding.",
		},
		"ssh_port": {
			Description: "SSH TCP port. Set to 0 for the default port (22).",
		},
		"login_protection_threshold": {
			Description: "Number of failed login attempts before blocking. Set to 0 for the pfSense default (30).",
		},
		"login_protection_blocktime": {
			Description: "Duration in seconds to block after exceeding threshold. Set to 0 for the pfSense default (120).",
		},
		"login_protection_detection_time": {
			Description: "Detection time window in seconds. Set to 0 for the pfSense default (1800).",
		},
		"login_protection_pass_list": {
			Description: "Space-separated list of IP addresses or CIDR networks to whitelist from login protection.",
		},
		"serial_terminal": {
			Description: "Enable serial terminal.",
		},
		"serial_speed": {
			Description:         fmt.Sprintf("Serial port speed in bps. Options: %s. Defaults to '%d'.", wrapElementsJoin(wrapElements(intSliceToStringSlice(pfsense.AdvancedAdmin{}.SerialSpeedOptions()), ""), ""), pfsense.DefaultAdvancedAdminSerialSpeed),
			MarkdownDescription: fmt.Sprintf("Serial port speed in bps. Options: %s. Defaults to `%d`.", wrapElementsJoin(intSliceToStringSlice(pfsense.AdvancedAdmin{}.SerialSpeedOptions()), "`"), pfsense.DefaultAdvancedAdminSerialSpeed),
		},
		"primary_console": {
			Description:         fmt.Sprintf("Primary console. Options: %s. Defaults to '%s'.", wrapElementsJoin(pfsense.AdvancedAdmin{}.PrimaryConsoleOptions(), "'"), pfsense.DefaultAdvancedAdminPrimaryConsole),
			MarkdownDescription: fmt.Sprintf("Primary console. Options: %s. Defaults to `%s`.", wrapElementsJoin(pfsense.AdvancedAdmin{}.PrimaryConsoleOptions(), "`"), pfsense.DefaultAdvancedAdminPrimaryConsole),
		},
		"disable_console_menu": {
			Description: "Disable console menu (password protect the console menu).",
		},
	}
}

func (SystemAdvancedAdminModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"webgui_protocol":                 types.StringType,
		"ssl_certificate":                 types.StringType,
		"webgui_port":                     types.Int64Type,
		"max_processes":                   types.Int64Type,
		"disable_http_redirect":           types.BoolType,
		"disable_hsts":                    types.BoolType,
		"ocsp_staple":                     types.BoolType,
		"login_autocomplete":              types.BoolType,
		"quiet_login":                     types.BoolType,
		"roaming":                         types.BoolType,
		"disable_anti_lockout":            types.BoolType,
		"disable_dns_rebind_check":        types.BoolType,
		"disable_http_referer_check":      types.BoolType,
		"alternate_hostnames":             types.StringType,
		"page_name_first":                 types.BoolType,
		"ssh_enabled":                     types.BoolType,
		"sshd_key_only":                   types.StringType,
		"sshd_agent_forwarding":           types.BoolType,
		"ssh_port":                        types.Int64Type,
		"login_protection_threshold":      types.Int64Type,
		"login_protection_blocktime":      types.Int64Type,
		"login_protection_detection_time": types.Int64Type,
		"login_protection_pass_list":      types.StringType,
		"serial_terminal":                 types.BoolType,
		"serial_speed":                    types.Int64Type,
		"primary_console":                 types.StringType,
		"disable_console_menu":            types.BoolType,
	}
}

func (m *SystemAdvancedAdminModel) Set(_ context.Context, a pfsense.AdvancedAdmin) diag.Diagnostics {
	var diags diag.Diagnostics

	// webConfigurator
	m.WebGUIProto = types.StringValue(a.WebGUIProto)

	if a.SSLCertRef != "" {
		m.SSLCertRef = types.StringValue(a.SSLCertRef)
	} else {
		m.SSLCertRef = types.StringNull()
	}

	m.WebGUIPort = types.Int64Value(int64(a.WebGUIPort))
	m.MaxProcs = types.Int64Value(int64(a.MaxProcs))
	m.DisableHTTPRedirect = types.BoolValue(a.DisableHTTPRedirect)
	m.DisableHSTS = types.BoolValue(a.DisableHSTS)
	m.OCSPStaple = types.BoolValue(a.OCSPStaple)
	m.LoginAutocomplete = types.BoolValue(a.LoginAutocomplete)
	m.QuietLogin = types.BoolValue(a.QuietLogin)
	m.Roaming = types.BoolValue(a.Roaming)
	m.NoAntiLockout = types.BoolValue(a.NoAntiLockout)
	m.NoDNSRebindCheck = types.BoolValue(a.NoDNSRebindCheck)
	m.NoHTTPRefererCheck = types.BoolValue(a.NoHTTPRefererCheck)

	if a.AlternateHostnames != "" {
		m.AlternateHostnames = types.StringValue(a.AlternateHostnames)
	} else {
		m.AlternateHostnames = types.StringNull()
	}

	m.PageNameFirst = types.BoolValue(a.PageNameFirst)

	// SSH
	m.SSHEnabled = types.BoolValue(a.SSHEnabled)
	m.SSHdKeyOnly = types.StringValue(a.SSHdKeyOnly)
	m.SSHdAgentForwarding = types.BoolValue(a.SSHdAgentForwarding)
	m.SSHPort = types.Int64Value(int64(a.SSHPort))

	// Login Protection
	m.SshguardThreshold = types.Int64Value(int64(a.SshguardThreshold))
	m.SshguardBlocktime = types.Int64Value(int64(a.SshguardBlocktime))
	m.SshguardDetectionTime = types.Int64Value(int64(a.SshguardDetectionTime))

	if a.SshguardWhitelist != "" {
		m.SshguardWhitelist = types.StringValue(a.SshguardWhitelist)
	} else {
		m.SshguardWhitelist = types.StringNull()
	}

	// Serial
	m.EnableSerial = types.BoolValue(a.EnableSerial)
	m.SerialSpeed = types.Int64Value(int64(a.SerialSpeed))
	m.PrimaryConsole = types.StringValue(a.PrimaryConsole)

	// Console
	m.DisableConsoleMenu = types.BoolValue(a.DisableConsoleMenu)

	return diags
}

func (m SystemAdvancedAdminModel) Value(_ context.Context, a *pfsense.AdvancedAdmin) diag.Diagnostics {
	var diags diag.Diagnostics

	// webConfigurator
	addPathError(
		&diags,
		path.Root("webgui_protocol"),
		"WebGUI protocol cannot be parsed",
		a.SetWebGUIProto(m.WebGUIProto.ValueString()),
	)

	if !m.SSLCertRef.IsNull() {
		a.SSLCertRef = m.SSLCertRef.ValueString()
	}

	addPathError(
		&diags,
		path.Root("webgui_port"),
		"WebGUI port cannot be parsed",
		a.SetWebGUIPort(int(m.WebGUIPort.ValueInt64())),
	)

	addPathError(
		&diags,
		path.Root("max_processes"),
		"Max processes cannot be parsed",
		a.SetMaxProcs(int(m.MaxProcs.ValueInt64())),
	)

	a.DisableHTTPRedirect = m.DisableHTTPRedirect.ValueBool()
	a.DisableHSTS = m.DisableHSTS.ValueBool()
	a.OCSPStaple = m.OCSPStaple.ValueBool()
	a.LoginAutocomplete = m.LoginAutocomplete.ValueBool()
	a.QuietLogin = m.QuietLogin.ValueBool()
	a.Roaming = m.Roaming.ValueBool()
	a.NoAntiLockout = m.NoAntiLockout.ValueBool()
	a.NoDNSRebindCheck = m.NoDNSRebindCheck.ValueBool()
	a.NoHTTPRefererCheck = m.NoHTTPRefererCheck.ValueBool()

	if !m.AlternateHostnames.IsNull() {
		a.AlternateHostnames = m.AlternateHostnames.ValueString()
	}

	a.PageNameFirst = m.PageNameFirst.ValueBool()

	// SSH
	a.SSHEnabled = m.SSHEnabled.ValueBool()

	addPathError(
		&diags,
		path.Root("sshd_key_only"),
		"SSH key only mode cannot be parsed",
		a.SetSSHdKeyOnly(m.SSHdKeyOnly.ValueString()),
	)

	a.SSHdAgentForwarding = m.SSHdAgentForwarding.ValueBool()

	addPathError(
		&diags,
		path.Root("ssh_port"),
		"SSH port cannot be parsed",
		a.SetSSHPort(int(m.SSHPort.ValueInt64())),
	)

	// Login Protection
	a.SshguardThreshold = int(m.SshguardThreshold.ValueInt64())
	a.SshguardBlocktime = int(m.SshguardBlocktime.ValueInt64())
	a.SshguardDetectionTime = int(m.SshguardDetectionTime.ValueInt64())

	if !m.SshguardWhitelist.IsNull() {
		a.SshguardWhitelist = m.SshguardWhitelist.ValueString()
	}

	// Serial
	a.EnableSerial = m.EnableSerial.ValueBool()

	addPathError(
		&diags,
		path.Root("serial_speed"),
		"Serial speed cannot be parsed",
		a.SetSerialSpeed(int(m.SerialSpeed.ValueInt64())),
	)

	addPathError(
		&diags,
		path.Root("primary_console"),
		"Primary console cannot be parsed",
		a.SetPrimaryConsole(m.PrimaryConsole.ValueString()),
	)

	// Console
	a.DisableConsoleMenu = m.DisableConsoleMenu.ValueBool()

	return diags
}

// intSliceToStringSlice converts a slice of ints to a slice of strings.
func intSliceToStringSlice(ints []int) []string {
	strs := make([]string, 0, len(ints))
	for _, i := range ints {
		strs = append(strs, fmt.Sprintf("%d", i))
	}

	return strs
}
