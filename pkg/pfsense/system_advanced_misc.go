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
	DefaultAdvancedMiscPowerdMode = "hadp"
)

// AdvancedMisc represents the System > Advanced > Miscellaneous configuration.
type AdvancedMisc struct {
	// Proxy Support
	ProxyURL  string // hostname or IP
	ProxyPort int    // 0 = not set
	ProxyUser string
	ProxyPass string // sensitive

	// Load Balancing
	LBUseSticky bool
	SrcTrack    int // source tracking timeout in seconds, 0 = default

	// Power Savings - Intel Speed Shift (hardware-dependent)
	HWPState             string // "enabled" or "disabled", empty if unsupported
	HWPStateControlLevel string // "0" (core) or "1" (package), empty if unsupported
	HWPStateEPP          int    // 0-100, energy/performance preference, -1 if unsupported

	// Power Savings - PowerD
	PowerdEnable      bool
	PowerdACMode      string // hadp/adp/min/max
	PowerdBatteryMode string
	PowerdNormalMode  string

	// Cryptographic & Thermal Hardware
	CryptoHardware  string // aesni/cryptodev/aesni_cryptodev/empty
	ThermalHardware string // coretemp/amdtemp/empty

	// Security Mitigations
	PTIDisabled bool   // disable kernel PTI
	MDSDisable  string // "" (default), "0"-"3"

	// Schedules
	ScheduleStates bool // don't kill connections when schedule expires

	// Gateway Monitoring
	GWDownKillStates           string // none/down/all/empty
	SkipRulesGWDown            bool
	DPingerDontAddStaticRoutes bool

	// RAM Disk Settings
	UseMFSTmpVar        bool
	UseMFSTmpSize       int // MiB, 0 = default (40)
	UseMFSVarSize       int // MiB, 0 = default (60)
	RRDBackup           int // 0 = disabled, 1-24 hours
	DHCPBackup          int
	LogsBackup          int
	CaptivePortalBackup int

	// Hardware Settings
	HardDiskStandby string // "" (always on) or standby time code

	// PHP Settings
	PHPMemoryLimit int // 0 = default, 128-max MiB

	// Installation Feedback
	DoNotSendUniqueID bool
}

func (a *AdvancedMisc) SetPowerdACMode(mode string) error {
	return a.setPowerdMode(&a.PowerdACMode, mode, "AC")
}

func (a *AdvancedMisc) SetPowerdBatteryMode(mode string) error {
	return a.setPowerdMode(&a.PowerdBatteryMode, mode, "battery")
}

func (a *AdvancedMisc) SetPowerdNormalMode(mode string) error {
	return a.setPowerdMode(&a.PowerdNormalMode, mode, "normal/unknown")
}

func (a *AdvancedMisc) setPowerdMode(field *string, mode string, label string) error {
	valid := AdvancedMisc{}.PowerdModeOptions()
	for _, v := range valid {
		if mode == v {
			*field = mode

			return nil
		}
	}

	return fmt.Errorf("%w, %s power mode must be one of: %s", ErrClientValidation, label, strings.Join(valid, ", "))
}

func (AdvancedMisc) PowerdModeOptions() []string {
	return []string{"hadp", "adp", "min", "max"}
}

func (a *AdvancedMisc) SetCryptoHardware(module string) error {
	valid := AdvancedMisc{}.CryptoHardwareOptions()
	for _, v := range valid {
		if module == v {
			a.CryptoHardware = module

			return nil
		}
	}

	return fmt.Errorf("%w, crypto hardware must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedMisc) CryptoHardwareOptions() []string {
	return []string{"", "aesni", "cryptodev", "aesni_cryptodev"}
}

func (a *AdvancedMisc) SetThermalHardware(module string) error {
	valid := AdvancedMisc{}.ThermalHardwareOptions()
	for _, v := range valid {
		if module == v {
			a.ThermalHardware = module

			return nil
		}
	}

	return fmt.Errorf("%w, thermal hardware must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedMisc) ThermalHardwareOptions() []string {
	return []string{"", "coretemp", "amdtemp"}
}

func (a *AdvancedMisc) SetGWDownKillStates(mode string) error {
	valid := AdvancedMisc{}.GWDownKillStatesOptions()
	for _, v := range valid {
		if mode == v {
			a.GWDownKillStates = mode

			return nil
		}
	}

	return fmt.Errorf("%w, gateway down kill states must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedMisc) GWDownKillStatesOptions() []string {
	return []string{"", "none", "down", "all"}
}

func (a *AdvancedMisc) SetMDSDisable(mode string) error {
	valid := AdvancedMisc{}.MDSDisableOptions()
	for _, v := range valid {
		if mode == v {
			a.MDSDisable = mode

			return nil
		}
	}

	return fmt.Errorf("%w, MDS disable mode must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedMisc) MDSDisableOptions() []string {
	return []string{"", "0", "1", "2", "3"}
}

func (a *AdvancedMisc) SetHardDiskStandby(code string) error {
	valid := AdvancedMisc{}.HardDiskStandbyOptions()
	for _, v := range valid {
		if code == v {
			a.HardDiskStandby = code

			return nil
		}
	}

	return fmt.Errorf("%w, hard disk standby must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedMisc) HardDiskStandbyOptions() []string {
	// "" = always on, then values correspond to time periods:
	// 6=0.5min, 12=1min, 24=2min, 36=3min, 48=4min, 60=5min,
	// 90=7.5min, 120=10min, 180=15min, 240=20min, 241=30min, 242=60min
	return []string{"", "6", "12", "24", "36", "48", "60", "90", "120", "180", "240", "241", "242"}
}

// advancedMiscResponse is the JSON shape returned by the PHP read command.
type advancedMiscResponse struct {
	// Proxy Support
	ProxyURL  json.RawMessage `json:"proxyurl"`
	ProxyPort json.RawMessage `json:"proxyport"`
	ProxyUser json.RawMessage `json:"proxyuser"`
	ProxyPass json.RawMessage `json:"proxypass"`

	// Load Balancing
	LBUseSticky json.RawMessage `json:"lb_use_sticky"`
	SrcTrack    json.RawMessage `json:"srctrack"`

	// Intel Speed Shift
	HWPState             json.RawMessage `json:"hwpstate"`
	HWPStateControlLevel json.RawMessage `json:"hwpstate_control_level"`
	HWPStateEPP          json.RawMessage `json:"hwpstate_epp"`

	// PowerD
	PowerdEnable      json.RawMessage `json:"powerd_enable"`
	PowerdACMode      json.RawMessage `json:"powerd_ac_mode"`
	PowerdBatteryMode json.RawMessage `json:"powerd_battery_mode"`
	PowerdNormalMode  json.RawMessage `json:"powerd_normal_mode"`

	// Cryptographic & Thermal Hardware
	CryptoHardware  json.RawMessage `json:"crypto_hardware"`
	ThermalHardware json.RawMessage `json:"thermal_hardware"`

	// Security Mitigations
	PTIDisabled json.RawMessage `json:"pti_disabled"`
	MDSDisable  json.RawMessage `json:"mds_disable"`

	// Schedules
	ScheduleStates json.RawMessage `json:"schedule_states"`

	// Gateway Monitoring
	GWDownKillStates           json.RawMessage `json:"gw_down_kill_states"`
	SkipRulesGWDown            json.RawMessage `json:"skip_rules_gw_down"`
	DPingerDontAddStaticRoutes json.RawMessage `json:"dpinger_dont_add_static_routes"`

	// RAM Disk Settings
	UseMFSTmpVar        json.RawMessage `json:"use_mfs_tmpvar"`
	UseMFSTmpSize       json.RawMessage `json:"use_mfs_tmp_size"`
	UseMFSVarSize       json.RawMessage `json:"use_mfs_var_size"`
	RRDBackup           json.RawMessage `json:"rrdbackup"`
	DHCPBackup          json.RawMessage `json:"dhcpbackup"`
	LogsBackup          json.RawMessage `json:"logsbackup"`
	CaptivePortalBackup json.RawMessage `json:"captiveportalbackup"`

	// Hardware Settings
	HardDiskStandby json.RawMessage `json:"harddiskstandby"`

	// PHP Settings
	PHPMemoryLimit json.RawMessage `json:"php_memory_limit"`

	// Installation Feedback
	DoNotSendUniqueID json.RawMessage `json:"do_not_send_uniqueid"`
}

func parseAdvancedMiscResponse(resp advancedMiscResponse) (AdvancedMisc, error) {
	var a AdvancedMisc

	// Proxy Support
	a.ProxyURL = rawToString(resp.ProxyURL)
	a.ProxyPort = rawToInt(resp.ProxyPort)
	a.ProxyUser = rawToString(resp.ProxyUser)
	a.ProxyPass = rawToString(resp.ProxyPass)

	// Load Balancing
	a.LBUseSticky = rawIsPresent(resp.LBUseSticky)
	a.SrcTrack = rawToInt(resp.SrcTrack)

	// Intel Speed Shift — config stores "enabled"/"disabled" string, not presence-based
	hwpstate := rawToString(resp.HWPState)
	a.HWPState = hwpstate
	a.HWPStateControlLevel = rawToString(resp.HWPStateControlLevel)
	a.HWPStateEPP = rawToInt(resp.HWPStateEPP)
	if a.HWPStateEPP == 0 && !rawIsPresent(resp.HWPStateEPP) {
		a.HWPStateEPP = -1 // -1 indicates unsupported/not configured
	}

	// PowerD
	a.PowerdEnable = rawIsPresent(resp.PowerdEnable)

	acMode := rawToString(resp.PowerdACMode)
	if acMode == "" {
		acMode = DefaultAdvancedMiscPowerdMode
	}

	if err := a.SetPowerdACMode(acMode); err != nil {
		return a, err
	}

	batteryMode := rawToString(resp.PowerdBatteryMode)
	if batteryMode == "" {
		batteryMode = DefaultAdvancedMiscPowerdMode
	}

	if err := a.SetPowerdBatteryMode(batteryMode); err != nil {
		return a, err
	}

	normalMode := rawToString(resp.PowerdNormalMode)
	if normalMode == "" {
		normalMode = DefaultAdvancedMiscPowerdMode
	}

	if err := a.SetPowerdNormalMode(normalMode); err != nil {
		return a, err
	}

	// Cryptographic & Thermal Hardware
	if err := a.SetCryptoHardware(rawToString(resp.CryptoHardware)); err != nil {
		return a, err
	}

	if err := a.SetThermalHardware(rawToString(resp.ThermalHardware)); err != nil {
		return a, err
	}

	// Security Mitigations
	a.PTIDisabled = rawIsPresent(resp.PTIDisabled)

	if err := a.SetMDSDisable(rawToString(resp.MDSDisable)); err != nil {
		return a, err
	}

	// Schedules
	a.ScheduleStates = rawIsPresent(resp.ScheduleStates)

	// Gateway Monitoring
	if err := a.SetGWDownKillStates(rawToString(resp.GWDownKillStates)); err != nil {
		return a, err
	}

	a.SkipRulesGWDown = rawIsPresent(resp.SkipRulesGWDown)
	a.DPingerDontAddStaticRoutes = rawIsPresent(resp.DPingerDontAddStaticRoutes)

	// RAM Disk Settings
	a.UseMFSTmpVar = rawIsPresent(resp.UseMFSTmpVar)
	a.UseMFSTmpSize = rawToInt(resp.UseMFSTmpSize)
	a.UseMFSVarSize = rawToInt(resp.UseMFSVarSize)
	a.RRDBackup = rawToInt(resp.RRDBackup)
	a.DHCPBackup = rawToInt(resp.DHCPBackup)
	a.LogsBackup = rawToInt(resp.LogsBackup)
	a.CaptivePortalBackup = rawToInt(resp.CaptivePortalBackup)

	// Hardware Settings
	if err := a.SetHardDiskStandby(rawToString(resp.HardDiskStandby)); err != nil {
		return a, err
	}

	// PHP Settings
	a.PHPMemoryLimit = rawToInt(resp.PHPMemoryLimit)

	// Installation Feedback
	a.DoNotSendUniqueID = rawIsPresent(resp.DoNotSendUniqueID)

	return a, nil
}

func (pf *Client) getAdvancedMisc(ctx context.Context) (*AdvancedMisc, error) {
	command := "$sys = config_get_path('system', array());" +
		"$out = array(" +
		"'proxyurl' => isset($sys['proxyurl']) ? $sys['proxyurl'] : null," +
		"'proxyport' => isset($sys['proxyport']) ? $sys['proxyport'] : null," +
		"'proxyuser' => isset($sys['proxyuser']) ? $sys['proxyuser'] : null," +
		"'proxypass' => isset($sys['proxypass']) ? $sys['proxypass'] : null," +
		"'lb_use_sticky' => isset($sys['lb_use_sticky']) ? $sys['lb_use_sticky'] : null," +
		"'srctrack' => isset($sys['srctrack']) ? $sys['srctrack'] : null," +
		"'hwpstate' => isset($sys['hwpstate']) ? $sys['hwpstate'] : null," +
		"'hwpstate_control_level' => isset($sys['hwpstate_control_level']) ? $sys['hwpstate_control_level'] : null," +
		"'hwpstate_epp' => isset($sys['hwpstate_epp']) ? $sys['hwpstate_epp'] : null," +
		"'powerd_enable' => isset($sys['powerd_enable']) ? $sys['powerd_enable'] : null," +
		"'powerd_ac_mode' => isset($sys['powerd_ac_mode']) ? $sys['powerd_ac_mode'] : null," +
		"'powerd_battery_mode' => isset($sys['powerd_battery_mode']) ? $sys['powerd_battery_mode'] : null," +
		"'powerd_normal_mode' => isset($sys['powerd_normal_mode']) ? $sys['powerd_normal_mode'] : null," +
		"'crypto_hardware' => isset($sys['crypto_hardware']) ? $sys['crypto_hardware'] : null," +
		"'thermal_hardware' => isset($sys['thermal_hardware']) ? $sys['thermal_hardware'] : null," +
		"'pti_disabled' => isset($sys['pti_disabled']) ? $sys['pti_disabled'] : null," +
		"'mds_disable' => isset($sys['mds_disable']) ? $sys['mds_disable'] : null," +
		"'schedule_states' => isset($sys['schedule_states']) ? $sys['schedule_states'] : null," +
		"'gw_down_kill_states' => isset($sys['gw_down_kill_states']) ? $sys['gw_down_kill_states'] : null," +
		"'skip_rules_gw_down' => isset($sys['skip_rules_gw_down']) ? $sys['skip_rules_gw_down'] : null," +
		"'dpinger_dont_add_static_routes' => isset($sys['dpinger_dont_add_static_routes']) ? $sys['dpinger_dont_add_static_routes'] : null," +
		"'use_mfs_tmpvar' => isset($sys['use_mfs_tmpvar']) ? $sys['use_mfs_tmpvar'] : null," +
		"'use_mfs_tmp_size' => isset($sys['use_mfs_tmp_size']) ? $sys['use_mfs_tmp_size'] : null," +
		"'use_mfs_var_size' => isset($sys['use_mfs_var_size']) ? $sys['use_mfs_var_size'] : null," +
		"'rrdbackup' => isset($sys['rrdbackup']) ? $sys['rrdbackup'] : null," +
		"'dhcpbackup' => isset($sys['dhcpbackup']) ? $sys['dhcpbackup'] : null," +
		"'logsbackup' => isset($sys['logsbackup']) ? $sys['logsbackup'] : null," +
		"'captiveportalbackup' => isset($sys['captiveportalbackup']) ? $sys['captiveportalbackup'] : null," +
		"'harddiskstandby' => isset($sys['harddiskstandby']) ? $sys['harddiskstandby'] : null," +
		"'php_memory_limit' => isset($sys['php_memory_limit']) ? $sys['php_memory_limit'] : null," +
		"'do_not_send_uniqueid' => isset($sys['do_not_send_uniqueid']) ? $sys['do_not_send_uniqueid'] : null" +
		");" +
		"print(json_encode($out));"

	var resp advancedMiscResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	a, err := parseAdvancedMiscResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("%w advanced misc response, %w", ErrUnableToParse, err)
	}

	return &a, nil
}

func (pf *Client) GetAdvancedMisc(ctx context.Context) (*AdvancedMisc, error) {
	defer pf.read(&pf.mutexes.AdvancedMisc)()

	a, err := pf.getAdvancedMisc(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w advanced misc, %w", ErrGetOperationFailed, err)
	}

	return a, nil
}

func advancedMiscFormValues(a AdvancedMisc) url.Values {
	values := url.Values{
		"save": {"Save"},
	}

	// Proxy Support
	if a.ProxyURL != "" {
		values.Set("proxyurl", a.ProxyURL)
	}

	if a.ProxyPort > 0 {
		values.Set("proxyport", strconv.Itoa(a.ProxyPort))
	}

	if a.ProxyUser != "" {
		values.Set("proxyuser", a.ProxyUser)
	}

	if a.ProxyPass != "" {
		values.Set("proxypass", a.ProxyPass)
		values.Set("proxypass_confirm", a.ProxyPass)
	}

	// Load Balancing
	if a.LBUseSticky {
		values.Set("lb_use_sticky", "yes")
	}

	if a.SrcTrack > 0 {
		values.Set("srctrack", strconv.Itoa(a.SrcTrack))
	} else {
		values.Set("srctrack", "0")
	}

	// Intel Speed Shift — only send if configured
	if a.HWPState != "" {
		if a.HWPState == "enabled" {
			values.Set("hwpstate", "yes")
		}

		if a.HWPStateControlLevel != "" {
			values.Set("hwpstate_control_level", a.HWPStateControlLevel)
		}

		if a.HWPStateEPP >= 0 {
			values.Set("hwpstate_epp", strconv.Itoa(a.HWPStateEPP))
		}
	}

	// PowerD
	if a.PowerdEnable {
		values.Set("powerd_enable", "yes")
	}

	values.Set("powerd_ac_mode", a.PowerdACMode)
	values.Set("powerd_battery_mode", a.PowerdBatteryMode)
	values.Set("powerd_normal_mode", a.PowerdNormalMode)

	// Cryptographic & Thermal Hardware
	values.Set("crypto_hardware", a.CryptoHardware)
	values.Set("thermal_hardware", a.ThermalHardware)

	// Security Mitigations
	if a.PTIDisabled {
		values.Set("pti_disabled", "yes")
	}

	values.Set("mds_disable", a.MDSDisable)

	// Schedules
	if a.ScheduleStates {
		values.Set("schedule_states", "yes")
	}

	// Gateway Monitoring
	values.Set("gw_down_kill_states", a.GWDownKillStates)

	if a.SkipRulesGWDown {
		values.Set("skip_rules_gw_down", "yes")
	}

	if a.DPingerDontAddStaticRoutes {
		values.Set("dpinger_dont_add_static_routes", "yes")
	}

	// RAM Disk Settings
	if a.UseMFSTmpVar {
		values.Set("use_mfs_tmpvar", "yes")
	}

	if a.UseMFSTmpSize > 0 {
		values.Set("use_mfs_tmp_size", strconv.Itoa(a.UseMFSTmpSize))
	}

	if a.UseMFSVarSize > 0 {
		values.Set("use_mfs_var_size", strconv.Itoa(a.UseMFSVarSize))
	}

	values.Set("rrdbackup", strconv.Itoa(a.RRDBackup))
	values.Set("dhcpbackup", strconv.Itoa(a.DHCPBackup))
	values.Set("logsbackup", strconv.Itoa(a.LogsBackup))
	values.Set("captiveportalbackup", strconv.Itoa(a.CaptivePortalBackup))

	// Hardware Settings
	values.Set("harddiskstandby", a.HardDiskStandby)

	// PHP Settings
	if a.PHPMemoryLimit > 0 {
		values.Set("php_memory_limit", strconv.Itoa(a.PHPMemoryLimit))
	}

	// Installation Feedback
	if a.DoNotSendUniqueID {
		values.Set("do_not_send_uniqueid", "yes")
	}

	return values
}

func (pf *Client) UpdateAdvancedMisc(ctx context.Context, a AdvancedMisc) (*AdvancedMisc, error) {
	defer pf.write(&pf.mutexes.AdvancedMisc)()

	relativeURL := url.URL{Path: "system_advanced_misc.php"}
	values := advancedMiscFormValues(a)

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return nil, fmt.Errorf("%w advanced misc, %w", ErrUpdateOperationFailed, err)
	}

	if err := scrapeHTMLValidationErrors(doc); err != nil {
		return nil, fmt.Errorf("%w advanced misc, %w", ErrUpdateOperationFailed, err)
	}

	result, err := pf.getAdvancedMisc(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w advanced misc after updating, %w", ErrGetOperationFailed, err)
	}

	return result, nil
}

func (pf *Client) ApplyAdvancedMiscChanges(ctx context.Context) error {
	pf.mutexes.AdvancedMiscApply.Lock()
	defer pf.mutexes.AdvancedMiscApply.Unlock()

	// Mirrors saveSystemAdvancedMisc(): filter_configure, activate_powerd,
	// load_crypto, load_thermal_hardware, system_resolvconf_generate
	command := "$retval = 0;" +
		"system_resolvconf_generate(true);" +
		"$retval |= filter_configure();" +
		"activate_powerd();" +
		"load_crypto();" +
		"load_thermal_hardware();" +
		"print(json_encode($retval));"

	var result int
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply advanced misc changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}
