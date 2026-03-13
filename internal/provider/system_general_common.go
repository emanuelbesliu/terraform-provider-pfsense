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

// SystemGeneralDNSServerModel represents a single DNS server entry.
type SystemGeneralDNSServerModel struct {
	Address  types.String `tfsdk:"address"`
	Hostname types.String `tfsdk:"hostname"`
	Gateway  types.String `tfsdk:"gateway"`
}

func (SystemGeneralDNSServerModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"address": {
			Description: "IP address of the DNS server.",
		},
		"hostname": {
			Description: "Hostname for DNS over TLS (DoT) requests.",
		},
		"gateway": {
			Description: "Gateway to use for reaching this DNS server. Leave empty for default routing.",
		},
	}
}

func (SystemGeneralDNSServerModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"address":  types.StringType,
		"hostname": types.StringType,
		"gateway":  types.StringType,
	}
}

// SystemGeneralModel represents the Terraform model for system general settings.
type SystemGeneralModel struct {
	Hostname     types.String `tfsdk:"hostname"`
	Domain       types.String `tfsdk:"domain"`
	DNSServers   types.List   `tfsdk:"dns_servers"`
	DNSOverride  types.Bool   `tfsdk:"dns_override"`
	DNSLocalhost types.String `tfsdk:"dns_localhost"`

	// Localization
	Timezone    types.String `tfsdk:"timezone"`
	Timeservers types.String `tfsdk:"timeservers"`
	Language    types.String `tfsdk:"language"`

	// webConfigurator
	WebGUICSS                      types.String `tfsdk:"webgui_theme"`
	LoginCSS                       types.String `tfsdk:"login_color"`
	LoginShowHost                  types.Bool   `tfsdk:"login_show_host"`
	WebGUIFixedMenu                types.Bool   `tfsdk:"webgui_fixed_menu"`
	DashboardColumns               types.Int64  `tfsdk:"dashboard_columns"`
	WebGUILeftColumnHyper          types.Bool   `tfsdk:"webgui_left_column_hyper"`
	DisableAliasPopupDetail        types.Bool   `tfsdk:"disable_alias_popup_detail"`
	DashboardAvailableWidgetsPanel types.Bool   `tfsdk:"dashboard_available_widgets_panel"`
	SystemLogsFilterPanel          types.Bool   `tfsdk:"system_logs_filter_panel"`
	SystemLogsManageLogPanel       types.Bool   `tfsdk:"system_logs_manage_log_panel"`
	StatusMonitoringSettingsPanel  types.Bool   `tfsdk:"status_monitoring_settings_panel"`
	RowOrderDragging               types.Bool   `tfsdk:"row_order_dragging"`
	InterfacesSort                 types.Bool   `tfsdk:"interfaces_sort"`
	RequireStateFilter             types.Bool   `tfsdk:"require_state_filter"`
	HostnameInMenu                 types.String `tfsdk:"hostname_in_menu"`
}

func (SystemGeneralModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"hostname": {
			Description: fmt.Sprintf("System hostname. Defaults to '%s'.", pfsense.DefaultSystemHostname),
		},
		"domain": {
			Description: fmt.Sprintf("System domain. Defaults to '%s'.", pfsense.DefaultSystemDomain),
		},
		"dns_servers": {
			Description: fmt.Sprintf("DNS server configuration (up to %d entries).", pfsense.MaxDNSServers),
		},
		"dns_override": {
			Description: fmt.Sprintf("Allow DNS server list to be overridden by DHCP/PPP on WAN. Defaults to '%t'.", pfsense.DefaultSystemDNSAllowOverride),
		},
		"dns_localhost": {
			Description:         fmt.Sprintf("DNS resolution behavior. Options: %s. Empty string means use default behavior.", wrapElementsJoin(pfsense.SystemGeneral{}.DNSLocalhostOptions(), "'")),
			MarkdownDescription: fmt.Sprintf("DNS resolution behavior. Options: %s. Empty string means use default behavior.", wrapElementsJoin(pfsense.SystemGeneral{}.DNSLocalhostOptions(), "`")),
		},
		"timezone": {
			Description: fmt.Sprintf("System timezone. Defaults to '%s'.", pfsense.DefaultSystemTimezone),
		},
		"timeservers": {
			Description: fmt.Sprintf("Space-separated list of NTP servers. Defaults to '%s'.", pfsense.DefaultSystemTimeservers),
		},
		"language": {
			Description: fmt.Sprintf("WebGUI language. Defaults to '%s'.", pfsense.DefaultSystemLanguage),
		},
		"webgui_theme": {
			Description: fmt.Sprintf("WebGUI theme CSS file. Defaults to '%s'.", pfsense.DefaultSystemWebGUICSS),
		},
		"login_color": {
			Description: fmt.Sprintf("Login page color (hex without #). Defaults to '%s'.", pfsense.DefaultSystemLoginCSS),
		},
		"login_show_host": {
			Description: "Show hostname on login page.",
		},
		"webgui_fixed_menu": {
			Description: "Fix the navigation bar at the top of the page.",
		},
		"dashboard_columns": {
			Description: fmt.Sprintf("Number of dashboard columns (1-6). Defaults to '%d'.", pfsense.DefaultSystemDashboardColumns),
		},
		"webgui_left_column_hyper": {
			Description: "Enable left column navigation hyperlinks.",
		},
		"disable_alias_popup_detail": {
			Description: "Disable alias popup detail on mouse hover.",
		},
		"dashboard_available_widgets_panel": {
			Description: "Show the available widgets panel on the dashboard.",
		},
		"system_logs_filter_panel": {
			Description: "Show the log filter panel by default in System Logs.",
		},
		"system_logs_manage_log_panel": {
			Description: "Show the manage log panel by default in System Logs.",
		},
		"status_monitoring_settings_panel": {
			Description: "Show the monitoring settings panel by default.",
		},
		"row_order_dragging": {
			Description: "Disable dragging of table rows to reorder.",
		},
		"interfaces_sort": {
			Description: "Sort interfaces alphabetically in the navigation.",
		},
		"require_state_filter": {
			Description: "Require a filter to be entered before showing results in firewall states.",
		},
		"hostname_in_menu": {
			Description:         fmt.Sprintf("Display hostname in the menu bar. Options: %s. Defaults to '%s'.", wrapElementsJoin(pfsense.SystemGeneral{}.HostnameInMenuOptions(), "'"), pfsense.DefaultSystemHostnameInMenu),
			MarkdownDescription: fmt.Sprintf("Display hostname in the menu bar. Options: %s. Defaults to `%s`.", wrapElementsJoin(pfsense.SystemGeneral{}.HostnameInMenuOptions(), "`"), pfsense.DefaultSystemHostnameInMenu),
		},
	}
}

func (SystemGeneralModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"hostname":                          types.StringType,
		"domain":                            types.StringType,
		"dns_servers":                       types.ListType{ElemType: types.ObjectType{AttrTypes: SystemGeneralDNSServerModel{}.AttrTypes()}},
		"dns_override":                      types.BoolType,
		"dns_localhost":                     types.StringType,
		"timezone":                          types.StringType,
		"timeservers":                       types.StringType,
		"language":                          types.StringType,
		"webgui_theme":                      types.StringType,
		"login_color":                       types.StringType,
		"login_show_host":                   types.BoolType,
		"webgui_fixed_menu":                 types.BoolType,
		"dashboard_columns":                 types.Int64Type,
		"webgui_left_column_hyper":          types.BoolType,
		"disable_alias_popup_detail":        types.BoolType,
		"dashboard_available_widgets_panel": types.BoolType,
		"system_logs_filter_panel":          types.BoolType,
		"system_logs_manage_log_panel":      types.BoolType,
		"status_monitoring_settings_panel":  types.BoolType,
		"row_order_dragging":                types.BoolType,
		"interfaces_sort":                   types.BoolType,
		"require_state_filter":              types.BoolType,
		"hostname_in_menu":                  types.StringType,
	}
}

func (m *SystemGeneralModel) Set(ctx context.Context, sg pfsense.SystemGeneral) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Hostname = types.StringValue(sg.Hostname)
	m.Domain = types.StringValue(sg.Domain)
	m.DNSOverride = types.BoolValue(sg.DNSOverride)

	if sg.DNSLocalhost != "" {
		m.DNSLocalhost = types.StringValue(sg.DNSLocalhost)
	} else {
		m.DNSLocalhost = types.StringNull()
	}

	// DNS servers
	if len(sg.DNSServers) > 0 {
		dnsModels := make([]SystemGeneralDNSServerModel, 0, len(sg.DNSServers))
		for _, entry := range sg.DNSServers {
			model := SystemGeneralDNSServerModel{
				Address: types.StringValue(entry.Address),
			}

			if entry.Hostname != "" {
				model.Hostname = types.StringValue(entry.Hostname)
			} else {
				model.Hostname = types.StringNull()
			}

			if entry.Gateway != "" {
				model.Gateway = types.StringValue(entry.Gateway)
			} else {
				model.Gateway = types.StringNull()
			}

			dnsModels = append(dnsModels, model)
		}

		dnsList, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: SystemGeneralDNSServerModel{}.AttrTypes()}, dnsModels)
		diags.Append(d...)
		m.DNSServers = dnsList
	} else {
		m.DNSServers = types.ListNull(types.ObjectType{AttrTypes: SystemGeneralDNSServerModel{}.AttrTypes()})
	}

	// Localization
	m.Timezone = types.StringValue(sg.Timezone)
	m.Timeservers = types.StringValue(sg.Timeservers)
	m.Language = types.StringValue(sg.Language)

	// webConfigurator
	m.WebGUICSS = types.StringValue(sg.WebGUICSS)
	m.LoginCSS = types.StringValue(sg.LoginCSS)
	m.LoginShowHost = types.BoolValue(sg.LoginShowHost)
	m.WebGUIFixedMenu = types.BoolValue(sg.WebGUIFixedMenu)
	m.DashboardColumns = types.Int64Value(int64(sg.DashboardColumns))
	m.WebGUILeftColumnHyper = types.BoolValue(sg.WebGUILeftColumnHyper)
	m.DisableAliasPopupDetail = types.BoolValue(sg.DisableAliasPopupDetail)
	m.DashboardAvailableWidgetsPanel = types.BoolValue(sg.DashboardAvailableWidgetsPanel)
	m.SystemLogsFilterPanel = types.BoolValue(sg.SystemLogsFilterPanel)
	m.SystemLogsManageLogPanel = types.BoolValue(sg.SystemLogsManageLogPanel)
	m.StatusMonitoringSettingsPanel = types.BoolValue(sg.StatusMonitoringSettingsPanel)
	m.RowOrderDragging = types.BoolValue(sg.RowOrderDragging)
	m.InterfacesSort = types.BoolValue(sg.InterfacesSort)
	m.RequireStateFilter = types.BoolValue(sg.RequireStateFilter)
	m.HostnameInMenu = types.StringValue(sg.HostnameInMenu)

	return diags
}

func (m SystemGeneralModel) Value(ctx context.Context, sg *pfsense.SystemGeneral) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("hostname"),
		"Hostname cannot be parsed",
		sg.SetHostname(m.Hostname.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("domain"),
		"Domain cannot be parsed",
		sg.SetDomain(m.Domain.ValueString()),
	)

	sg.DNSOverride = m.DNSOverride.ValueBool()

	if !m.DNSLocalhost.IsNull() {
		sg.DNSLocalhost = m.DNSLocalhost.ValueString()
	}

	// DNS servers
	if !m.DNSServers.IsNull() && !m.DNSServers.IsUnknown() {
		var dnsModels []SystemGeneralDNSServerModel
		diags.Append(m.DNSServers.ElementsAs(ctx, &dnsModels, false)...)

		if !diags.HasError() {
			entries := make([]pfsense.DNSServerEntry, 0, len(dnsModels))
			for _, dm := range dnsModels {
				entry := pfsense.DNSServerEntry{
					Address: dm.Address.ValueString(),
				}

				if !dm.Hostname.IsNull() {
					entry.Hostname = dm.Hostname.ValueString()
				}

				if !dm.Gateway.IsNull() {
					entry.Gateway = dm.Gateway.ValueString()
				}

				entries = append(entries, entry)
			}

			sg.DNSServers = entries
		}
	}

	addPathError(
		&diags,
		path.Root("timezone"),
		"Timezone cannot be parsed",
		sg.SetTimezone(m.Timezone.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("timeservers"),
		"Timeservers cannot be parsed",
		sg.SetTimeservers(m.Timeservers.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("language"),
		"Language cannot be parsed",
		sg.SetLanguage(m.Language.ValueString()),
	)

	// webConfigurator
	sg.WebGUICSS = m.WebGUICSS.ValueString()
	sg.LoginCSS = m.LoginCSS.ValueString()
	sg.LoginShowHost = m.LoginShowHost.ValueBool()
	sg.WebGUIFixedMenu = m.WebGUIFixedMenu.ValueBool()

	addPathError(
		&diags,
		path.Root("dashboard_columns"),
		"Dashboard columns cannot be parsed",
		sg.SetDashboardColumns(int(m.DashboardColumns.ValueInt64())),
	)

	sg.WebGUILeftColumnHyper = m.WebGUILeftColumnHyper.ValueBool()
	sg.DisableAliasPopupDetail = m.DisableAliasPopupDetail.ValueBool()
	sg.DashboardAvailableWidgetsPanel = m.DashboardAvailableWidgetsPanel.ValueBool()
	sg.SystemLogsFilterPanel = m.SystemLogsFilterPanel.ValueBool()
	sg.SystemLogsManageLogPanel = m.SystemLogsManageLogPanel.ValueBool()
	sg.StatusMonitoringSettingsPanel = m.StatusMonitoringSettingsPanel.ValueBool()
	sg.RowOrderDragging = m.RowOrderDragging.ValueBool()
	sg.InterfacesSort = m.InterfacesSort.ValueBool()
	sg.RequireStateFilter = m.RequireStateFilter.ValueBool()

	addPathError(
		&diags,
		path.Root("hostname_in_menu"),
		"Hostname in menu cannot be parsed",
		sg.SetHostnameInMenu(m.HostnameInMenu.ValueString()),
	)

	return diags
}
