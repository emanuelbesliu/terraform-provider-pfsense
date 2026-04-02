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
	DefaultAdvancedAdminWebGUIProto           = "http"
	DefaultAdvancedAdminWebGUIPort            = 0 // 0 means use default (80/443)
	DefaultAdvancedAdminMaxProcs              = 2
	DefaultAdvancedAdminSSHPort               = 0 // 0 means default (22)
	DefaultAdvancedAdminSSHdKeyOnly           = "disabled"
	DefaultAdvancedAdminSerialSpeed           = 115200
	DefaultAdvancedAdminPrimaryConsole        = "video"
	DefaultAdvancedAdminSshguardThreshold     = 0 // 0 means default (30)
	DefaultAdvancedAdminSshguardBlocktime     = 0 // 0 means default (120)
	DefaultAdvancedAdminSshguardDetectionTime = 0 // 0 means default (1800)
)

// AdvancedAdmin represents the System > Advanced > Admin Access configuration.
type AdvancedAdmin struct {
	// webConfigurator
	WebGUIProto         string
	SSLCertRef          string
	WebGUIPort          int
	MaxProcs            int
	DisableHTTPRedirect bool
	DisableHSTS         bool
	OCSPStaple          bool
	LoginAutocomplete   bool
	QuietLogin          bool
	Roaming             bool
	NoAntiLockout       bool
	NoDNSRebindCheck    bool
	NoHTTPRefererCheck  bool
	AlternateHostnames  string
	PageNameFirst       bool

	// Secure Shell
	SSHEnabled          bool
	SSHdKeyOnly         string
	SSHdAgentForwarding bool
	SSHPort             int

	// Login Protection
	SshguardThreshold     int
	SshguardBlocktime     int
	SshguardDetectionTime int
	SshguardWhitelist     string

	// Serial Communications
	EnableSerial   bool
	SerialSpeed    int
	PrimaryConsole string

	// Console Options
	DisableConsoleMenu bool
}

func (a *AdvancedAdmin) SetWebGUIProto(proto string) error {
	valid := []string{"http", "https"}
	for _, v := range valid {
		if proto == v {
			a.WebGUIProto = proto

			return nil
		}
	}

	return fmt.Errorf("%w, webgui protocol must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedAdmin) WebGUIProtoOptions() []string {
	return []string{"http", "https"}
}

func (a *AdvancedAdmin) SetSSHdKeyOnly(mode string) error {
	valid := []string{"disabled", "enabled", "both"}
	for _, v := range valid {
		if mode == v {
			a.SSHdKeyOnly = mode

			return nil
		}
	}

	return fmt.Errorf("%w, sshd key only must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedAdmin) SSHdKeyOnlyOptions() []string {
	return []string{"disabled", "enabled", "both"}
}

func (a *AdvancedAdmin) SetSerialSpeed(speed int) error {
	valid := []int{115200, 57600, 38400, 19200, 14400, 9600}
	for _, v := range valid {
		if speed == v {
			a.SerialSpeed = speed

			return nil
		}
	}

	return fmt.Errorf("%w, serial speed must be one of: 115200, 57600, 38400, 19200, 14400, 9600", ErrClientValidation)
}

func (AdvancedAdmin) SerialSpeedOptions() []int {
	return []int{115200, 57600, 38400, 19200, 14400, 9600}
}

func (a *AdvancedAdmin) SetPrimaryConsole(console string) error {
	valid := []string{"serial", "video"}
	for _, v := range valid {
		if console == v {
			a.PrimaryConsole = console

			return nil
		}
	}

	return fmt.Errorf("%w, primary console must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedAdmin) PrimaryConsoleOptions() []string {
	return []string{"serial", "video"}
}

func (a *AdvancedAdmin) SetWebGUIPort(port int) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("%w, webgui port must be between 0 and 65535 (0 for default)", ErrClientValidation)
	}

	a.WebGUIPort = port

	return nil
}

func (a *AdvancedAdmin) SetMaxProcs(procs int) error {
	if procs < 1 || procs > 500 {
		return fmt.Errorf("%w, max processes must be between 1 and 500", ErrClientValidation)
	}

	a.MaxProcs = procs

	return nil
}

func (a *AdvancedAdmin) SetSSHPort(port int) error {
	if port < 0 || port > 65535 {
		return fmt.Errorf("%w, ssh port must be between 0 and 65535 (0 for default)", ErrClientValidation)
	}

	a.SSHPort = port

	return nil
}

// advancedAdminResponse is the JSON shape returned by the PHP read command.
type advancedAdminResponse struct {
	WebGUI             *webguiAdminResponse `json:"webgui"`
	SSH                json.RawMessage      `json:"ssh"`
	SerialSpeed        json.RawMessage      `json:"serialspeed"`
	PrimaryConsole     json.RawMessage      `json:"primaryconsole"`
	EnableSerial       json.RawMessage      `json:"enableserial"`
	DisableConsoleMenu json.RawMessage      `json:"disableconsolemenu"`
	SshguardThreshold  json.RawMessage      `json:"sshguard_threshold"`
	SshguardBlocktime  json.RawMessage      `json:"sshguard_blocktime"`
	SshguardDetection  json.RawMessage      `json:"sshguard_detection_time"`
	SshguardWhitelist  string               `json:"sshguard_whitelist"`
}

type webguiAdminResponse struct {
	Protocol            string          `json:"protocol"`
	SSLCertRef          string          `json:"ssl-certref"`
	Port                string          `json:"port"`
	MaxProcs            string          `json:"max_procs"`
	DisableHTTPRedirect json.RawMessage `json:"disablehttpredirect"`
	DisableHSTS         json.RawMessage `json:"disablehsts"`
	OCSPStaple          json.RawMessage `json:"ocsp-staple"`
	LoginAutocomplete   json.RawMessage `json:"loginautocomplete"`
	QuietLogin          json.RawMessage `json:"quietlogin"`
	Roaming             string          `json:"roaming"`
	NoAntiLockout       json.RawMessage `json:"noantilockout"`
	NoDNSRebindCheck    json.RawMessage `json:"nodnsrebindcheck"`
	NoHTTPRefererCheck  json.RawMessage `json:"nohttpreferercheck"`
	AlternateHostnames  string          `json:"althostnames"`
	PageNameFirst       json.RawMessage `json:"pagenamefirst"`
}

type sshAdminResponse struct {
	Enable              string `json:"enable"`
	SSHdKeyOnly         string `json:"sshdkeyonly"`
	SSHdAgentForwarding string `json:"sshdagentforwarding"`
	Port                string `json:"port"`
}

// rawIsPresent checks if a json.RawMessage represents a truthy value.
// In pfSense config, booleans are presence-based: the key existing in the XML
// means true, regardless of value (could be "", "true", "yes", etc.).
// When serialized via json_encode: null = key absent = false, anything else
// (including "" empty string) = key present = true.
func rawIsPresent(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}

	s := string(raw)
	// null means the key was absent in config
	if s == "null" {
		return false
	}

	// Any other value — including "" (empty string) — means the key exists
	return true
}

func parseAdvancedAdminResponse(resp advancedAdminResponse) (AdvancedAdmin, error) {
	var a AdvancedAdmin

	// webConfigurator
	if resp.WebGUI != nil {
		wg := resp.WebGUI

		proto := wg.Protocol
		if proto == "" {
			proto = DefaultAdvancedAdminWebGUIProto
		}

		if err := a.SetWebGUIProto(proto); err != nil {
			return a, err
		}

		a.SSLCertRef = wg.SSLCertRef

		if wg.Port != "" {
			p, err := strconv.Atoi(wg.Port)
			if err != nil {
				return a, fmt.Errorf("%w, unable to parse webgui port '%s'", ErrUnableToParse, wg.Port)
			}

			a.WebGUIPort = p
		}

		procs := DefaultAdvancedAdminMaxProcs
		if wg.MaxProcs != "" {
			p, err := strconv.Atoi(wg.MaxProcs)
			if err != nil {
				return a, fmt.Errorf("%w, unable to parse max procs '%s'", ErrUnableToParse, wg.MaxProcs)
			}

			procs = p
		}

		a.MaxProcs = procs

		a.DisableHTTPRedirect = rawIsPresent(wg.DisableHTTPRedirect)
		a.DisableHSTS = rawIsPresent(wg.DisableHSTS)
		a.OCSPStaple = rawIsPresent(wg.OCSPStaple)
		a.LoginAutocomplete = rawIsPresent(wg.LoginAutocomplete)
		a.QuietLogin = rawIsPresent(wg.QuietLogin)

		// Roaming: "enabled" means true, anything else means false
		a.Roaming = (wg.Roaming == "enabled")

		a.NoAntiLockout = rawIsPresent(wg.NoAntiLockout)
		a.NoDNSRebindCheck = rawIsPresent(wg.NoDNSRebindCheck)
		a.NoHTTPRefererCheck = rawIsPresent(wg.NoHTTPRefererCheck)
		a.AlternateHostnames = wg.AlternateHostnames
		a.PageNameFirst = rawIsPresent(wg.PageNameFirst)
	} else {
		a.WebGUIProto = DefaultAdvancedAdminWebGUIProto
		a.MaxProcs = DefaultAdvancedAdminMaxProcs
		a.Roaming = true // pfSense defaults to enabled
	}

	// SSH — handle both object and string/null formats.
	// pfSense may return the ssh config as an object {"enable":"enabled",...}
	// or as a simple string "enabled" or "" depending on version/config state.
	if len(resp.SSH) > 0 && string(resp.SSH) != "null" {
		// Try to unmarshal as an object first
		var sshObj sshAdminResponse
		if err := json.Unmarshal(resp.SSH, &sshObj); err == nil {
			a.SSHEnabled = (sshObj.Enable == "enabled")

			if sshObj.SSHdKeyOnly != "" {
				a.SSHdKeyOnly = sshObj.SSHdKeyOnly
			} else {
				a.SSHdKeyOnly = DefaultAdvancedAdminSSHdKeyOnly
			}

			a.SSHdAgentForwarding = (sshObj.SSHdAgentForwarding == "enabled")

			if sshObj.Port != "" {
				p, err := strconv.Atoi(sshObj.Port)
				if err != nil {
					return a, fmt.Errorf("%w, unable to parse ssh port '%s'", ErrUnableToParse, sshObj.Port)
				}

				a.SSHPort = p
			}
		} else {
			// Fell through — SSH field is a string (e.g. "enabled" or "")
			sshStr := strings.Trim(string(resp.SSH), `"`)
			a.SSHEnabled = (sshStr == "enabled")
			a.SSHdKeyOnly = DefaultAdvancedAdminSSHdKeyOnly
		}
	} else {
		a.SSHdKeyOnly = DefaultAdvancedAdminSSHdKeyOnly
	}

	// Serial
	a.EnableSerial = rawIsPresent(resp.EnableSerial)

	speedStr := strings.Trim(string(resp.SerialSpeed), `"`)
	if speedStr != "" && speedStr != "null" {
		s, err := strconv.Atoi(speedStr)
		if err != nil {
			return a, fmt.Errorf("%w, unable to parse serial speed '%s'", ErrUnableToParse, speedStr)
		}

		a.SerialSpeed = s
	} else {
		a.SerialSpeed = DefaultAdvancedAdminSerialSpeed
	}

	consoleStr := strings.Trim(string(resp.PrimaryConsole), `"`)
	if consoleStr != "" && consoleStr != "null" {
		a.PrimaryConsole = consoleStr
	} else {
		a.PrimaryConsole = DefaultAdvancedAdminPrimaryConsole
	}

	// Console
	a.DisableConsoleMenu = rawIsPresent(resp.DisableConsoleMenu)

	// Login Protection
	parseIntField := func(raw json.RawMessage) int {
		s := strings.Trim(string(raw), `"`)
		if s == "" || s == "null" {
			return 0
		}

		v, _ := strconv.Atoi(s)

		return v
	}

	a.SshguardThreshold = parseIntField(resp.SshguardThreshold)
	a.SshguardBlocktime = parseIntField(resp.SshguardBlocktime)
	a.SshguardDetectionTime = parseIntField(resp.SshguardDetection)
	a.SshguardWhitelist = resp.SshguardWhitelist

	return a, nil
}

func (pf *Client) getAdvancedAdmin(ctx context.Context) (*AdvancedAdmin, error) {
	command := "$sys = config_get_path('system', array());" +
		"$out = array(" +
		"'webgui' => isset($sys['webgui']) ? $sys['webgui'] : null," +
		"'ssh' => isset($sys['ssh']) ? $sys['ssh'] : null," +
		"'serialspeed' => isset($sys['serialspeed']) ? $sys['serialspeed'] : null," +
		"'primaryconsole' => isset($sys['primaryconsole']) ? $sys['primaryconsole'] : null," +
		"'enableserial' => isset($sys['enableserial']) ? $sys['enableserial'] : null," +
		"'disableconsolemenu' => isset($sys['disableconsolemenu']) ? $sys['disableconsolemenu'] : null," +
		"'sshguard_threshold' => isset($sys['sshguard_threshold']) ? $sys['sshguard_threshold'] : ''," +
		"'sshguard_blocktime' => isset($sys['sshguard_blocktime']) ? $sys['sshguard_blocktime'] : ''," +
		"'sshguard_detection_time' => isset($sys['sshguard_detection_time']) ? $sys['sshguard_detection_time'] : ''," +
		"'sshguard_whitelist' => isset($sys['sshguard_whitelist']) ? $sys['sshguard_whitelist'] : ''" +
		");" +
		"print(json_encode($out));"

	var resp advancedAdminResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	a, err := parseAdvancedAdminResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("%w advanced admin response, %w", ErrUnableToParse, err)
	}

	return &a, nil
}

func (pf *Client) GetAdvancedAdmin(ctx context.Context) (*AdvancedAdmin, error) {
	defer pf.read(&pf.mutexes.AdvancedAdmin)()

	a, err := pf.getAdvancedAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w advanced admin, %w", ErrGetOperationFailed, err)
	}

	return a, nil
}

func advancedAdminFormValues(a AdvancedAdmin) url.Values {
	values := url.Values{
		"save": {"Save"},
	}

	// webConfigurator
	values.Set("webguiproto", a.WebGUIProto)

	if a.SSLCertRef != "" {
		values.Set("ssl-certref", a.SSLCertRef)
	}

	if a.WebGUIPort > 0 {
		values.Set("webguiport", strconv.Itoa(a.WebGUIPort))
	}

	values.Set("max_procs", strconv.Itoa(a.MaxProcs))

	if a.DisableHTTPRedirect {
		values.Set("webgui-redirect", "yes")
	}

	if a.DisableHSTS {
		values.Set("webgui-hsts", "yes")
	}

	if a.OCSPStaple {
		values.Set("ocsp-staple", "yes")
	}

	if a.LoginAutocomplete {
		values.Set("loginautocomplete", "yes")
	}

	if a.QuietLogin {
		values.Set("webgui-login-messages", "yes")
	}

	if a.Roaming {
		values.Set("roaming", "yes")
	}

	if a.NoAntiLockout {
		values.Set("noantilockout", "yes")
	}

	if a.NoDNSRebindCheck {
		values.Set("nodnsrebindcheck", "yes")
	}

	if a.NoHTTPRefererCheck {
		values.Set("nohttpreferercheck", "yes")
	}

	if a.AlternateHostnames != "" {
		values.Set("althostnames", a.AlternateHostnames)
	}

	if a.PageNameFirst {
		values.Set("pagenamefirst", "yes")
	}

	// SSH
	if a.SSHEnabled {
		values.Set("enablesshd", "yes")
	}

	values.Set("sshdkeyonly", a.SSHdKeyOnly)

	if a.SSHdAgentForwarding {
		values.Set("sshdagentforwarding", "yes")
	}

	if a.SSHPort > 0 {
		values.Set("sshport", strconv.Itoa(a.SSHPort))
	}

	// Login Protection
	if a.SshguardThreshold > 0 {
		values.Set("sshguard_threshold", strconv.Itoa(a.SshguardThreshold))
	}

	if a.SshguardBlocktime > 0 {
		values.Set("sshguard_blocktime", strconv.Itoa(a.SshguardBlocktime))
	}

	if a.SshguardDetectionTime > 0 {
		values.Set("sshguard_detection_time", strconv.Itoa(a.SshguardDetectionTime))
	}

	// Pass list — POST uses address0/address_subnet0, address1/address_subnet1, etc.
	if a.SshguardWhitelist != "" {
		entries := strings.Fields(a.SshguardWhitelist)
		for i, entry := range entries {
			parts := strings.SplitN(entry, "/", 2)
			values.Set("address"+strconv.Itoa(i), parts[0])

			if len(parts) > 1 {
				values.Set("address_subnet"+strconv.Itoa(i), parts[1])
			} else {
				values.Set("address_subnet"+strconv.Itoa(i), "32")
			}
		}
	}

	// Serial
	if a.EnableSerial {
		values.Set("enableserial", "yes")
	}

	values.Set("serialspeed", strconv.Itoa(a.SerialSpeed))
	values.Set("primaryconsole", a.PrimaryConsole)

	// Console
	if a.DisableConsoleMenu {
		values.Set("disableconsolemenu", "yes")
	}

	return values
}

func (pf *Client) UpdateAdvancedAdmin(ctx context.Context, a AdvancedAdmin) (*AdvancedAdmin, error) {
	defer pf.write(&pf.mutexes.AdvancedAdmin)()

	relativeURL := url.URL{Path: "system_advanced_admin.php"}
	values := advancedAdminFormValues(a)

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return nil, fmt.Errorf("%w advanced admin, %w", ErrUpdateOperationFailed, err)
	}

	if err := scrapeHTMLValidationErrors(doc); err != nil {
		return nil, fmt.Errorf("%w advanced admin, %w", ErrUpdateOperationFailed, err)
	}

	result, err := pf.getAdvancedAdmin(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w advanced admin after updating, %w", ErrGetOperationFailed, err)
	}

	return result, nil
}

func (pf *Client) ApplyAdvancedAdminChanges(ctx context.Context) error {
	pf.mutexes.AdvancedAdminApply.Lock()
	defer pf.mutexes.AdvancedAdminApply.Unlock()

	command := "$retval = 0;" +
		"$retval |= filter_configure();" +
		"$retval |= system_syslogd_start(true);" +
		"console_configure();" +
		"if (config_path_enabled('dnsmasq', 'enable')) {" +
		"services_dnsmasq_configure();" +
		"} elseif (config_path_enabled('unbound', 'enable')) {" +
		"services_unbound_configure();" +
		"}" +
		"print(json_encode($retval));"

	var result int
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply advanced admin changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}
