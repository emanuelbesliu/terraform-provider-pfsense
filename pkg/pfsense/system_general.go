package pfsense

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	DefaultSystemHostname         = "pfSense"
	DefaultSystemDomain           = "home.arpa"
	DefaultSystemTimezone         = "Etc/UTC"
	DefaultSystemTimeservers      = "2.pfsense.pool.ntp.org"
	DefaultSystemLanguage         = "en_US"
	DefaultSystemDNSAllowOverride = true
	DefaultSystemDNSLocalhost     = ""
	DefaultSystemWebGUICSS        = "pfSense.css"
	DefaultSystemLoginCSS         = "1e3f75"
	DefaultSystemDashboardColumns = 2
	DefaultSystemHostnameInMenu   = "none"

	MaxDNSServers = 4
)

// SystemGeneral represents the System > General Setup configuration.
type SystemGeneral struct {
	Hostname     string
	Domain       string
	DNSServers   []DNSServerEntry
	DNSOverride  bool
	DNSLocalhost string

	// Localization
	Timezone    string
	Timeservers string
	Language    string

	// webConfigurator
	WebGUICSS                      string
	LoginCSS                       string
	LoginShowHost                  bool
	WebGUIFixedMenu                bool
	DashboardColumns               int
	WebGUILeftColumnHyper          bool
	DisableAliasPopupDetail        bool
	DashboardAvailableWidgetsPanel bool
	SystemLogsFilterPanel          bool
	SystemLogsManageLogPanel       bool
	StatusMonitoringSettingsPanel  bool
	RowOrderDragging               bool
	InterfacesSort                 bool
	RequireStateFilter             bool
	HostnameInMenu                 string
}

// DNSServerEntry represents a single DNS server with optional gateway and hostname.
type DNSServerEntry struct {
	Address  string
	Hostname string
	Gateway  string
}

func (sg *SystemGeneral) SetHostname(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("%w, hostname must not be empty", ErrClientValidation)
	}

	if err := ValidateDNSLabel(hostname); err != nil {
		return fmt.Errorf("%w, invalid hostname: %w", ErrClientValidation, err)
	}

	sg.Hostname = hostname

	return nil
}

func (sg *SystemGeneral) SetDomain(domain string) error {
	if domain == "" {
		return fmt.Errorf("%w, domain must not be empty", ErrClientValidation)
	}

	if err := ValidateDomain(domain); err != nil {
		return fmt.Errorf("%w, invalid domain: %w", ErrClientValidation, err)
	}

	sg.Domain = domain

	return nil
}

func (sg *SystemGeneral) SetTimezone(timezone string) error {
	if timezone == "" {
		return fmt.Errorf("%w, timezone must not be empty", ErrClientValidation)
	}

	sg.Timezone = timezone

	return nil
}

func (sg *SystemGeneral) SetTimeservers(timeservers string) error {
	sg.Timeservers = timeservers

	return nil
}

func (sg *SystemGeneral) SetLanguage(language string) error {
	if language == "" {
		return fmt.Errorf("%w, language must not be empty", ErrClientValidation)
	}

	sg.Language = language

	return nil
}

func (sg *SystemGeneral) SetDashboardColumns(columns int) error {
	if columns < 1 || columns > 6 {
		return fmt.Errorf("%w, dashboard columns must be between 1 and 6", ErrClientValidation)
	}

	sg.DashboardColumns = columns

	return nil
}

func (sg *SystemGeneral) SetHostnameInMenu(value string) error {
	valid := []string{"none", "hostname", "fqdn"}
	for _, v := range valid {
		if value == v {
			sg.HostnameInMenu = value

			return nil
		}
	}

	return fmt.Errorf("%w, hostname_in_menu must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (SystemGeneral) HostnameInMenuOptions() []string {
	return []string{"none", "hostname", "fqdn"}
}

func (SystemGeneral) DNSLocalhostOptions() []string {
	return []string{"", "local", "remote", "none"}
}

// systemGeneralResponse is the JSON shape returned by the PHP read command.
type systemGeneralResponse struct {
	Hostname     string          `json:"hostname"`
	Domain       string          `json:"domain"`
	DNSServers   json.RawMessage `json:"dnsserver"`
	DNSOverride  string          `json:"dnsallowoverride"`
	DNSLocalhost string          `json:"dnslocalhost"`
	Timezone     string          `json:"timezone"`
	Timeservers  string          `json:"timeservers"`
	Language     string          `json:"language"`
	WebGUI       *webGUIResponse `json:"webgui"`
}

type webGUIResponse struct {
	WebGUICSS                      string `json:"webguicss"`
	LoginCSS                       string `json:"logincss"`
	LoginShowHost                  string `json:"loginshowhost"`
	WebGUIFixedMenu                string `json:"webguifixedmenu"`
	DashboardColumns               string `json:"dashboardcolumns"`
	WebGUILeftColumnHyper          string `json:"webguileftcolumnhyper"`
	DisableAliasPopupDetail        string `json:"disablealiaspopupdetail"`
	DashboardAvailableWidgetsPanel string `json:"dashboardavailablewidgetspanel"`
	SystemLogsFilterPanel          string `json:"systemlogsfilterpanel"`
	SystemLogsManageLogPanel       string `json:"systemlogsmanagelogpanel"`
	StatusMonitoringSettingsPanel  string `json:"statusmonitoringsettingspanel"`
	RowOrderDragging               string `json:"roworderdragging"`
	InterfacesSort                 string `json:"interfacessort"`
	RequireStateFilter             string `json:"requirestatefilter"`
	HostnameInMenu                 string `json:"hostname_in_menu"`
}

func parseSystemGeneralResponse(resp systemGeneralResponse) (SystemGeneral, error) {
	var sg SystemGeneral

	// Hostname
	hostname := resp.Hostname
	if hostname == "" {
		hostname = DefaultSystemHostname
	}

	if err := sg.SetHostname(hostname); err != nil {
		return sg, err
	}

	// Domain
	domain := resp.Domain
	if domain == "" {
		domain = DefaultSystemDomain
	}

	if err := sg.SetDomain(domain); err != nil {
		return sg, err
	}

	// DNS servers - can be a JSON array of strings or a single string
	sg.DNSServers = parseDNSServers(resp.DNSServers)

	// DNS override - present means enabled
	sg.DNSOverride = resp.DNSOverride != ""

	// DNS localhost
	sg.DNSLocalhost = resp.DNSLocalhost

	// Timezone
	timezone := resp.Timezone
	if timezone == "" {
		timezone = DefaultSystemTimezone
	}

	if err := sg.SetTimezone(timezone); err != nil {
		return sg, err
	}

	// Timeservers
	timeservers := resp.Timeservers
	if timeservers == "" {
		timeservers = DefaultSystemTimeservers
	}

	sg.Timeservers = timeservers

	// Language
	language := resp.Language
	if language == "" {
		language = DefaultSystemLanguage
	}

	if err := sg.SetLanguage(language); err != nil {
		return sg, err
	}

	// WebGUI settings
	if resp.WebGUI != nil {
		wg := resp.WebGUI

		sg.WebGUICSS = wg.WebGUICSS
		if sg.WebGUICSS == "" {
			sg.WebGUICSS = DefaultSystemWebGUICSS
		}

		sg.LoginCSS = wg.LoginCSS
		if sg.LoginCSS == "" {
			sg.LoginCSS = DefaultSystemLoginCSS
		}

		sg.LoginShowHost = wg.LoginShowHost != ""
		sg.WebGUIFixedMenu = wg.WebGUIFixedMenu != ""

		columns := DefaultSystemDashboardColumns
		if wg.DashboardColumns != "" {
			c, err := strconv.Atoi(wg.DashboardColumns)
			if err != nil {
				return sg, fmt.Errorf("%w, unable to parse dashboard columns '%s'", ErrUnableToParse, wg.DashboardColumns)
			}

			columns = c
		}

		if err := sg.SetDashboardColumns(columns); err != nil {
			return sg, err
		}

		sg.WebGUILeftColumnHyper = wg.WebGUILeftColumnHyper != ""
		sg.DisableAliasPopupDetail = wg.DisableAliasPopupDetail != ""
		sg.DashboardAvailableWidgetsPanel = wg.DashboardAvailableWidgetsPanel != ""
		sg.SystemLogsFilterPanel = wg.SystemLogsFilterPanel != ""
		sg.SystemLogsManageLogPanel = wg.SystemLogsManageLogPanel != ""
		sg.StatusMonitoringSettingsPanel = wg.StatusMonitoringSettingsPanel != ""
		sg.RowOrderDragging = wg.RowOrderDragging != ""
		sg.InterfacesSort = wg.InterfacesSort != ""
		sg.RequireStateFilter = wg.RequireStateFilter != ""

		hostnameInMenu := wg.HostnameInMenu
		if hostnameInMenu == "" {
			hostnameInMenu = DefaultSystemHostnameInMenu
		}

		if err := sg.SetHostnameInMenu(hostnameInMenu); err != nil {
			return sg, err
		}
	} else {
		// Defaults when webgui section is missing
		sg.WebGUICSS = DefaultSystemWebGUICSS
		sg.LoginCSS = DefaultSystemLoginCSS
		sg.DashboardColumns = DefaultSystemDashboardColumns
		sg.HostnameInMenu = DefaultSystemHostnameInMenu
	}

	return sg, nil
}

func parseDNSServers(raw json.RawMessage) []DNSServerEntry {
	if len(raw) == 0 {
		return nil
	}

	// Try as array of strings first
	var arr []string
	if err := json.Unmarshal(raw, &arr); err == nil {
		entries := make([]DNSServerEntry, 0, len(arr))
		for _, addr := range arr {
			if addr != "" {
				entries = append(entries, DNSServerEntry{Address: addr})
			}
		}

		return entries
	}

	// Try as single string
	var single string
	if err := json.Unmarshal(raw, &single); err == nil && single != "" {
		return []DNSServerEntry{{Address: single}}
	}

	return nil
}

func (pf *Client) getSystemGeneral(ctx context.Context) (*SystemGeneral, error) {
	command := "$sys = config_get_path('system', array());" +
		"$out = array(" +
		"'hostname' => isset($sys['hostname']) ? $sys['hostname'] : ''," +
		"'domain' => isset($sys['domain']) ? $sys['domain'] : ''," +
		"'dnsserver' => isset($sys['dnsserver']) ? $sys['dnsserver'] : array()," +
		"'dnsallowoverride' => isset($sys['dnsallowoverride']) ? $sys['dnsallowoverride'] : ''," +
		"'dnslocalhost' => isset($sys['dnslocalhost']) ? $sys['dnslocalhost'] : ''," +
		"'timezone' => isset($sys['timezone']) ? $sys['timezone'] : ''," +
		"'timeservers' => isset($sys['timeservers']) ? $sys['timeservers'] : ''," +
		"'language' => isset($sys['language']) ? $sys['language'] : ''," +
		"'webgui' => isset($sys['webgui']) ? $sys['webgui'] : null" +
		");" +
		"print(json_encode($out));"

	var resp systemGeneralResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	sg, err := parseSystemGeneralResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("%w system general response, %w", ErrUnableToParse, err)
	}

	// Now read per-server gateway and hostname fields via a second PHP command.
	// These are stored separately in config as dns1gw, dns2gw, etc.
	if len(sg.DNSServers) > 0 {
		gwCommand := "$sys = config_get_path('system', array());" +
			"$out = array();" +
			"for ($i = 1; $i <= 4; $i++) {" +
			"$out[] = array(" +
			"'gateway' => isset($sys['dns' . $i . 'gw']) ? $sys['dns' . $i . 'gw'] : ''," +
			"'hostname' => isset($sys['dns' . $i . 'host']) ? $sys['dns' . $i . 'host'] : ''" +
			");" +
			"}" +
			"print(json_encode($out));"

		type dnsGWEntry struct {
			Gateway  string `json:"gateway"`
			Hostname string `json:"hostname"`
		}

		var gwEntries []dnsGWEntry
		if err := pf.executePHPCommand(ctx, gwCommand, &gwEntries); err != nil {
			return nil, err
		}

		for i := range sg.DNSServers {
			if i < len(gwEntries) {
				if gwEntries[i].Gateway != "" && gwEntries[i].Gateway != "none" {
					sg.DNSServers[i].Gateway = gwEntries[i].Gateway
				}

				sg.DNSServers[i].Hostname = gwEntries[i].Hostname
			}
		}
	}

	return &sg, nil
}

func (pf *Client) GetSystemGeneral(ctx context.Context) (*SystemGeneral, error) {
	defer pf.read(&pf.mutexes.SystemGeneral)()

	sg, err := pf.getSystemGeneral(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w system general, %w", ErrGetOperationFailed, err)
	}

	return sg, nil
}

func systemGeneralFormValues(sg SystemGeneral) url.Values {
	values := url.Values{
		"hostname":    {sg.Hostname},
		"domain":      {sg.Domain},
		"timezone":    {sg.Timezone},
		"timeservers": {sg.Timeservers},
		"language":    {sg.Language},
		"save":        {"Save"},
	}

	// DNS servers — form fields are 0-based (dns0, dns1, dns2, dns3)
	// Gateway fields are dnsgw0, dnsgw1, ... and hostname fields are dnshost0, dnshost1, ...
	// pfSense iterates with while(isset($_POST["dns{$i}"])) starting at 0,
	// so we must send contiguous 0-based fields and stop after the last server.
	for i := 0; i < MaxDNSServers; i++ {
		idx := strconv.Itoa(i)
		if i < len(sg.DNSServers) {
			values.Set("dns"+idx, sg.DNSServers[i].Address)

			gw := sg.DNSServers[i].Gateway
			if gw == "" {
				gw = "none"
			}

			values.Set("dnsgw"+idx, gw)
			values.Set("dnshost"+idx, sg.DNSServers[i].Hostname)
		} else {
			values.Set("dns"+idx, "")
			values.Set("dnsgw"+idx, "none")
			values.Set("dnshost"+idx, "")
		}
	}

	// DNS override
	if sg.DNSOverride {
		values.Set("dnsallowoverride", "yes")
	}

	// DNS localhost
	if sg.DNSLocalhost != "" {
		values.Set("dnslocalhost", sg.DNSLocalhost)
	}

	// webConfigurator settings
	values.Set("webguicss", sg.WebGUICSS)
	values.Set("logincss", sg.LoginCSS)
	values.Set("dashboardcolumns", strconv.Itoa(sg.DashboardColumns))

	if sg.LoginShowHost {
		values.Set("loginshowhost", "checked")
	}

	if sg.WebGUIFixedMenu {
		values.Set("webguifixedmenu", "fixed")
	}

	if sg.WebGUILeftColumnHyper {
		values.Set("webguileftcolumnhyper", "active")
	}

	if sg.DisableAliasPopupDetail {
		values.Set("disablealiaspopupdetail", "yes")
	}

	if sg.DashboardAvailableWidgetsPanel {
		values.Set("dashboardavailablewidgetspanel", "yes")
	}

	if sg.SystemLogsFilterPanel {
		values.Set("systemlogsfilterpanel", "yes")
	}

	if sg.SystemLogsManageLogPanel {
		values.Set("systemlogsmanagelogpanel", "yes")
	}

	if sg.StatusMonitoringSettingsPanel {
		values.Set("statusmonitoringsettingspanel", "yes")
	}

	if sg.RowOrderDragging {
		values.Set("roworderdragging", "yes")
	}

	if sg.InterfacesSort {
		values.Set("interfacessort", "yes")
	}

	if sg.RequireStateFilter {
		values.Set("requirestatefilter", "yes")
	}

	if sg.HostnameInMenu != "" && sg.HostnameInMenu != "none" {
		values.Set("hostname_in_menu", sg.HostnameInMenu)
	}

	return values
}

func (pf *Client) UpdateSystemGeneral(ctx context.Context, sg SystemGeneral) (*SystemGeneral, error) {
	defer pf.write(&pf.mutexes.SystemGeneral)()

	relativeURL := url.URL{Path: "system.php"}
	values := systemGeneralFormValues(sg)

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return nil, fmt.Errorf("%w system general, %w", ErrUpdateOperationFailed, err)
	}

	if err := scrapeHTMLValidationErrors(doc); err != nil {
		return nil, fmt.Errorf("%w system general, %w", ErrUpdateOperationFailed, err)
	}

	result, err := pf.getSystemGeneral(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w system general after updating, %w", ErrGetOperationFailed, err)
	}

	return result, nil
}

func (pf *Client) ApplySystemGeneralChanges(ctx context.Context) error {
	pf.mutexes.SystemGeneralApply.Lock()
	defer pf.mutexes.SystemGeneralApply.Unlock()

	command := "require_once('rrd.inc');" +
		"$retval = 0;" +
		"$retval |= system_hostname_configure();" +
		"$retval |= system_hosts_generate();" +
		"$retval |= system_resolvconf_generate();" +
		"if (is_array($config['dnsmasq']) && isset($config['dnsmasq']['enable'])) {" +
		"$retval |= services_dnsmasq_configure();" +
		"} elseif (is_array($config['unbound']) && isset($config['unbound']['enable'])) {" +
		"$retval |= services_unbound_configure();" +
		"}" +
		"$retval |= system_timezone_configure();" +
		"$retval |= system_ntp_configure();" +
		"$retval |= filter_configure();" +
		"print(json_encode($retval));"

	var result int
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply system general changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}
