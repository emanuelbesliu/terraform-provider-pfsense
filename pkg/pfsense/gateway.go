package pfsense

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var gatewayNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

const (
	DefaultGatewayWeight = 1
	MinGatewayWeight     = 1
	MaxGatewayWeight     = 30

	DefaultGatewayLatencyLow    = 200
	DefaultGatewayLatencyHigh   = 500
	DefaultGatewayLossLow       = 10
	DefaultGatewayLossHigh      = 20
	DefaultGatewayInterval      = 500
	DefaultGatewayLossInterval  = 2000
	DefaultGatewayTimePeriod    = 60000
	DefaultGatewayAlertInterval = 1000
	DefaultGatewayDataPayload   = 1
)

type gatewayResponse struct {
	Interface      string  `json:"interface"`
	IPProtocol     string  `json:"ipprotocol"`
	Name           string  `json:"name"`
	Gateway        string  `json:"gateway"`
	Description    string  `json:"descr"`
	Disabled       *string `json:"disabled"`
	DefaultGW      *string `json:"defaultgw"`
	Monitor        string  `json:"monitor"`
	MonitorDisable *string `json:"monitor_disable"`
	ActionDisable  *string `json:"action_disable"`
	ForceDown      *string `json:"force_down"`
	Weight         string  `json:"weight"`
	NonLocalGW     *string `json:"nonlocalgateway"`
	LatencyLow     string  `json:"latencylow"`
	LatencyHigh    string  `json:"latencyhigh"`
	LossLow        string  `json:"losslow"`
	LossHigh       string  `json:"losshigh"`
	Interval       string  `json:"interval"`
	LossInterval   string  `json:"loss_interval"`
	TimePeriod     string  `json:"time_period"`
	AlertInterval  string  `json:"alert_interval"`
	DataPayload    string  `json:"data_payload"`
	ControlID      int     `json:"controlID"` //nolint:tagliatelle
}

type Gateway struct {
	Interface      string
	IPProtocol     string
	Name           string
	GatewayIP      string
	Description    string
	Disabled       bool
	DefaultGW      bool
	Monitor        string
	MonitorDisable bool
	ActionDisable  bool
	ForceDown      bool
	Weight         int
	NonLocalGW     bool
	LatencyLow     int
	LatencyHigh    int
	LossLow        int
	LossHigh       int
	Interval       int
	LossInterval   int
	TimePeriod     int
	AlertInterval  int
	DataPayload    int
	controlID      int
}

func (Gateway) IPProtocols() []string {
	return []string{"inet", "inet6"}
}

func (gw *Gateway) SetInterface(iface string) error {
	if err := ValidateInterface(iface); err != nil {
		return fmt.Errorf("%w, invalid interface: %w", ErrClientValidation, err)
	}

	gw.Interface = iface

	return nil
}

func (gw *Gateway) SetIPProtocol(ipprotocol string) error {
	valid := false
	for _, p := range gw.IPProtocols() {
		if ipprotocol == p {
			valid = true

			break
		}
	}

	if !valid {
		return fmt.Errorf("%w, ip protocol must be one of: %s", ErrClientValidation, strings.Join(gw.IPProtocols(), ", "))
	}

	gw.IPProtocol = ipprotocol

	return nil
}

func (gw *Gateway) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("%w, name must not be empty", ErrClientValidation)
	}

	if len(name) > 31 {
		return fmt.Errorf("%w, name must be at most 31 characters", ErrClientValidation)
	}

	if !gatewayNameRegex.MatchString(name) {
		return fmt.Errorf("%w, name must start with a letter or underscore and contain only alphanumeric characters and underscores", ErrClientValidation)
	}

	gw.Name = name

	return nil
}

func (gw *Gateway) SetGatewayIP(gateway string) error {
	if gateway == "" {
		return fmt.Errorf("%w, gateway must not be empty", ErrClientValidation)
	}

	// Accept the literal "dynamic" for DHCP/PPPoE interfaces, or a valid IP address.
	if gateway != "dynamic" {
		if err := ValidateIPAddress(gateway, "Any"); err != nil {
			return fmt.Errorf("%w, gateway must be a valid IP address or 'dynamic': %w", ErrClientValidation, err)
		}
	}

	gw.GatewayIP = gateway

	return nil
}

func (gw *Gateway) SetDescription(description string) error {
	if len(description) > 200 {
		return fmt.Errorf("%w, description must be at most 200 characters", ErrClientValidation)
	}

	gw.Description = description

	return nil
}

func (gw *Gateway) SetDisabled(disabled bool) error {
	gw.Disabled = disabled

	return nil
}

func (gw *Gateway) SetMonitor(monitor string) error {
	if monitor != "" {
		if err := ValidateIPAddress(monitor, "Any"); err != nil {
			return fmt.Errorf("%w, monitor must be a valid IP address: %w", ErrClientValidation, err)
		}
	}

	gw.Monitor = monitor

	return nil
}

func (gw *Gateway) SetMonitorDisable(monitorDisable bool) error {
	gw.MonitorDisable = monitorDisable

	return nil
}

func (gw *Gateway) SetActionDisable(actionDisable bool) error {
	gw.ActionDisable = actionDisable

	return nil
}

func (gw *Gateway) SetForceDown(forceDown bool) error {
	gw.ForceDown = forceDown

	return nil
}

func (gw *Gateway) SetWeight(weight int) error {
	if weight < MinGatewayWeight || weight > MaxGatewayWeight {
		return fmt.Errorf("%w, weight must be between %d and %d", ErrClientValidation, MinGatewayWeight, MaxGatewayWeight)
	}

	gw.Weight = weight

	return nil
}

func (gw *Gateway) SetNonLocalGW(nonLocalGW bool) error {
	gw.NonLocalGW = nonLocalGW

	return nil
}

func (gw *Gateway) SetDefaultGW(defaultGW bool) error {
	gw.DefaultGW = defaultGW

	return nil
}

func (gw *Gateway) SetLatencyLow(latencyLow int) error {
	if latencyLow < 0 {
		return fmt.Errorf("%w, latency low threshold must be non-negative", ErrClientValidation)
	}

	gw.LatencyLow = latencyLow

	return nil
}

func (gw *Gateway) SetLatencyHigh(latencyHigh int) error {
	if latencyHigh < 0 {
		return fmt.Errorf("%w, latency high threshold must be non-negative", ErrClientValidation)
	}

	gw.LatencyHigh = latencyHigh

	return nil
}

func (gw *Gateway) SetLossLow(lossLow int) error {
	if lossLow < 0 || lossLow > 100 {
		return fmt.Errorf("%w, loss low threshold must be between 0 and 100", ErrClientValidation)
	}

	gw.LossLow = lossLow

	return nil
}

func (gw *Gateway) SetLossHigh(lossHigh int) error {
	if lossHigh < 0 || lossHigh > 100 {
		return fmt.Errorf("%w, loss high threshold must be between 0 and 100", ErrClientValidation)
	}

	gw.LossHigh = lossHigh

	return nil
}

func (gw *Gateway) SetInterval(interval int) error {
	if interval < 1 {
		return fmt.Errorf("%w, probe interval must be at least 1", ErrClientValidation)
	}

	gw.Interval = interval

	return nil
}

func (gw *Gateway) SetLossInterval(lossInterval int) error {
	if lossInterval < 0 {
		return fmt.Errorf("%w, loss interval must be non-negative", ErrClientValidation)
	}

	gw.LossInterval = lossInterval

	return nil
}

func (gw *Gateway) SetTimePeriod(timePeriod int) error {
	if timePeriod < 0 {
		return fmt.Errorf("%w, time period must be non-negative", ErrClientValidation)
	}

	gw.TimePeriod = timePeriod

	return nil
}

func (gw *Gateway) SetAlertInterval(alertInterval int) error {
	if alertInterval < 0 {
		return fmt.Errorf("%w, alert interval must be non-negative", ErrClientValidation)
	}

	gw.AlertInterval = alertInterval

	return nil
}

func (gw *Gateway) SetDataPayload(dataPayload int) error {
	if dataPayload < 0 {
		return fmt.Errorf("%w, data payload must be non-negative", ErrClientValidation)
	}

	gw.DataPayload = dataPayload

	return nil
}

type Gateways []Gateway

func (gateways Gateways) GetByName(name string) (*Gateway, error) {
	for _, gw := range gateways {
		if gw.Name == name {
			return &gw, nil
		}
	}

	return nil, fmt.Errorf("gateway %w with name '%s'", ErrNotFound, name)
}

func (gateways Gateways) GetControlIDByName(name string) (*int, error) {
	for _, gw := range gateways {
		if gw.Name == name {
			return &gw.controlID, nil
		}
	}

	return nil, fmt.Errorf("gateway %w with name '%s'", ErrNotFound, name)
}

func parseGatewayResponse(resp gatewayResponse, index int) (Gateway, error) {
	var gw Gateway

	if err := gw.SetName(resp.Name); err != nil {
		return gw, err
	}

	if err := gw.SetInterface(resp.Interface); err != nil {
		return gw, err
	}

	if err := gw.SetIPProtocol(resp.IPProtocol); err != nil {
		return gw, err
	}

	if err := gw.SetGatewayIP(resp.Gateway); err != nil {
		return gw, err
	}

	if err := gw.SetDescription(resp.Description); err != nil {
		return gw, err
	}

	gw.Disabled = resp.Disabled != nil
	// DefaultGW is determined by the system-level defaultgw4/defaultgw6 config,
	// not the per-gateway legacy flag. Set by the caller (getGateways).
	gw.MonitorDisable = resp.MonitorDisable != nil
	gw.ActionDisable = resp.ActionDisable != nil
	gw.ForceDown = resp.ForceDown != nil
	gw.NonLocalGW = resp.NonLocalGW != nil

	if err := gw.SetMonitor(resp.Monitor); err != nil {
		return gw, err
	}

	weight := DefaultGatewayWeight
	if resp.Weight != "" {
		w, err := strconv.Atoi(resp.Weight)
		if err != nil {
			return gw, fmt.Errorf("%w, unable to parse gateway weight '%s'", ErrUnableToParse, resp.Weight)
		}

		weight = w
	}

	if err := gw.SetWeight(weight); err != nil {
		return gw, err
	}

	if err := parseIntField(resp.LatencyLow, DefaultGatewayLatencyLow, gw.SetLatencyLow); err != nil {
		return gw, err
	}

	if err := parseIntField(resp.LatencyHigh, DefaultGatewayLatencyHigh, gw.SetLatencyHigh); err != nil {
		return gw, err
	}

	if err := parseIntField(resp.LossLow, DefaultGatewayLossLow, gw.SetLossLow); err != nil {
		return gw, err
	}

	if err := parseIntField(resp.LossHigh, DefaultGatewayLossHigh, gw.SetLossHigh); err != nil {
		return gw, err
	}

	if err := parseIntField(resp.Interval, DefaultGatewayInterval, gw.SetInterval); err != nil {
		return gw, err
	}

	if err := parseIntField(resp.LossInterval, DefaultGatewayLossInterval, gw.SetLossInterval); err != nil {
		return gw, err
	}

	if err := parseIntField(resp.TimePeriod, DefaultGatewayTimePeriod, gw.SetTimePeriod); err != nil {
		return gw, err
	}

	if err := parseIntField(resp.AlertInterval, DefaultGatewayAlertInterval, gw.SetAlertInterval); err != nil {
		return gw, err
	}

	if err := parseIntField(resp.DataPayload, DefaultGatewayDataPayload, gw.SetDataPayload); err != nil {
		return gw, err
	}

	gw.controlID = index

	return gw, nil
}

func parseIntField(raw string, defaultVal int, setter func(int) error) error {
	val := defaultVal
	if raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil {
			return fmt.Errorf("%w, unable to parse '%s'", ErrUnableToParse, raw)
		}

		val = v
	}

	return setter(val)
}

type gatewaysWithDefaults struct {
	DefaultGW4   string            `json:"defaultgw4"`
	DefaultGW6   string            `json:"defaultgw6"`
	GatewayItems []gatewayResponse `json:"gateway_items"`
}

func (pf *Client) getGateways(ctx context.Context) (*Gateways, error) {
	command := "$items = array();" +
		"if (isset($config['gateways']['gateway_item']) && is_array($config['gateways']['gateway_item'])) {" +
		"foreach ($config['gateways']['gateway_item'] as $k => $v) {" +
		"$v['controlID'] = $k; array_push($items, $v);" +
		"}};" +
		"$out = array(" +
		"'defaultgw4' => isset($config['gateways']['defaultgw4']) ? $config['gateways']['defaultgw4'] : ''," +
		"'defaultgw6' => isset($config['gateways']['defaultgw6']) ? $config['gateways']['defaultgw6'] : ''," +
		"'gateway_items' => $items" +
		");" +
		"print_r(json_encode($out));"

	var resp gatewaysWithDefaults
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	gateways := make(Gateways, 0, len(resp.GatewayItems))
	for _, gwResp := range resp.GatewayItems {
		gw, err := parseGatewayResponse(gwResp, gwResp.ControlID)
		if err != nil {
			return nil, fmt.Errorf("%w gateway response, %w", ErrUnableToParse, err)
		}

		// Determine default gateway status from system-level config, not the
		// per-gateway legacy flag (which is stripped by return_gateways_array).
		switch gw.IPProtocol {
		case "inet":
			gw.DefaultGW = resp.DefaultGW4 == gw.Name
		case "inet6":
			gw.DefaultGW = resp.DefaultGW6 == gw.Name
		}

		gateways = append(gateways, gw)
	}

	return &gateways, nil
}

// getGatewayEditID returns the index into get_gateways(GW_CACHE_INDEXED)
// for the given gateway name. The edit page (system_gateways_edit.php) and
// validate_gateway() use this index, which differs from the config.xml
// array index when gateways are reordered by get_gateways().
func (pf *Client) getGatewayEditID(ctx context.Context, name string) (*int, error) {
	type gwIndexEntry struct {
		Name string `json:"name"`
	}

	command := "require_once('gwlb.inc');" +
		"refresh_gateways();" +
		"$a = get_gateways(GW_CACHE_INDEXED);" +
		"$out = array();" +
		"foreach ($a as $k => $v) { $out[] = array('name' => $v['name']); }" +
		"print(json_encode($out));"

	var entries []gwIndexEntry
	if err := pf.executePHPCommand(ctx, command, &entries); err != nil {
		return nil, err
	}

	for i, entry := range entries {
		if entry.Name == name {
			return &i, nil
		}
	}

	return nil, fmt.Errorf("gateway %w with name '%s'", ErrNotFound, name)
}

func (pf *Client) GetGateways(ctx context.Context) (*Gateways, error) {
	defer pf.read(&pf.mutexes.Gateway)()

	gateways, err := pf.getGateways(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w gateways, %w", ErrGetOperationFailed, err)
	}

	return gateways, nil
}

func (pf *Client) GetGateway(ctx context.Context, name string) (*Gateway, error) {
	defer pf.read(&pf.mutexes.Gateway)()

	gateways, err := pf.getGateways(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w gateways, %w", ErrGetOperationFailed, err)
	}

	gw, err := gateways.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w gateway, %w", ErrGetOperationFailed, err)
	}

	return gw, nil
}

func gatewayFormValues(gwReq Gateway) url.Values {
	values := url.Values{
		"interface":  {gwReq.Interface},
		"ipprotocol": {gwReq.IPProtocol},
		"name":       {gwReq.Name},
		"gateway":    {gwReq.GatewayIP},
		"descr":      {gwReq.Description},
		"weight":     {strconv.Itoa(gwReq.Weight)},
		"save":       {"Save"},
	}

	if gwReq.Disabled {
		values.Set("disabled", "yes")
	}

	if gwReq.DefaultGW {
		values.Set("defaultgw", "yes")
	}

	if gwReq.Monitor != "" {
		values.Set("monitor", gwReq.Monitor)
	}

	if gwReq.MonitorDisable {
		values.Set("monitor_disable", "yes")
	}

	if gwReq.ActionDisable {
		values.Set("action_disable", "yes")
	}

	if gwReq.ForceDown {
		values.Set("force_down", "yes")
	}

	if gwReq.NonLocalGW {
		values.Set("nonlocalgateway", "yes")
	}

	if gwReq.LatencyLow != DefaultGatewayLatencyLow {
		values.Set("latencylow", strconv.Itoa(gwReq.LatencyLow))
	}

	if gwReq.LatencyHigh != DefaultGatewayLatencyHigh {
		values.Set("latencyhigh", strconv.Itoa(gwReq.LatencyHigh))
	}

	if gwReq.LossLow != DefaultGatewayLossLow {
		values.Set("losslow", strconv.Itoa(gwReq.LossLow))
	}

	if gwReq.LossHigh != DefaultGatewayLossHigh {
		values.Set("losshigh", strconv.Itoa(gwReq.LossHigh))
	}

	if gwReq.Interval != DefaultGatewayInterval {
		values.Set("interval", strconv.Itoa(gwReq.Interval))
	}

	if gwReq.LossInterval != DefaultGatewayLossInterval {
		values.Set("loss_interval", strconv.Itoa(gwReq.LossInterval))
	}

	if gwReq.TimePeriod != DefaultGatewayTimePeriod {
		values.Set("time_period", strconv.Itoa(gwReq.TimePeriod))
	}

	if gwReq.AlertInterval != DefaultGatewayAlertInterval {
		values.Set("alert_interval", strconv.Itoa(gwReq.AlertInterval))
	}

	if gwReq.DataPayload != DefaultGatewayDataPayload {
		values.Set("data_payload", strconv.Itoa(gwReq.DataPayload))
	}

	return values
}

func (pf *Client) createOrUpdateGateway(ctx context.Context, gwReq Gateway, controlID *int) error {
	relativeURL := url.URL{Path: "system_gateways_edit.php"}
	values := gatewayFormValues(gwReq)

	if controlID != nil {
		q := relativeURL.Query()
		q.Set("id", strconv.Itoa(*controlID))
		relativeURL.RawQuery = q.Encode()
		values.Set("id", strconv.Itoa(*controlID))
	}

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return err
	}

	return scrapeHTMLValidationErrors(doc)
}

func (pf *Client) CreateGateway(ctx context.Context, gwReq Gateway) (*Gateway, error) {
	defer pf.write(&pf.mutexes.Gateway)()

	if err := pf.createOrUpdateGateway(ctx, gwReq, nil); err != nil {
		return nil, fmt.Errorf("%w gateway, %w", ErrCreateOperationFailed, err)
	}

	if err := pf.ensureDefaultGateway(ctx, gwReq.Name, gwReq.IPProtocol, gwReq.DefaultGW); err != nil {
		return nil, fmt.Errorf("%w gateway, %w", ErrCreateOperationFailed, err)
	}

	gateways, err := pf.getGateways(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w gateways after creating, %w", ErrGetOperationFailed, err)
	}

	gw, err := gateways.GetByName(gwReq.Name)
	if err != nil {
		return nil, fmt.Errorf("%w gateway after creating, %w", ErrGetOperationFailed, err)
	}

	return gw, nil
}

func (pf *Client) UpdateGateway(ctx context.Context, gwReq Gateway) (*Gateway, error) {
	defer pf.write(&pf.mutexes.Gateway)()

	editID, err := pf.getGatewayEditID(ctx, gwReq.Name)
	if err != nil {
		return nil, fmt.Errorf("%w gateway, %w", ErrGetOperationFailed, err)
	}

	if err := pf.createOrUpdateGateway(ctx, gwReq, editID); err != nil {
		return nil, fmt.Errorf("%w gateway, %w", ErrUpdateOperationFailed, err)
	}

	if err := pf.ensureDefaultGateway(ctx, gwReq.Name, gwReq.IPProtocol, gwReq.DefaultGW); err != nil {
		return nil, fmt.Errorf("%w gateway, %w", ErrUpdateOperationFailed, err)
	}

	gateways, err := pf.getGateways(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w gateways after updating, %w", ErrGetOperationFailed, err)
	}

	gw, err := gateways.GetByName(gwReq.Name)
	if err != nil {
		return nil, fmt.Errorf("%w gateway after updating, %w", ErrGetOperationFailed, err)
	}

	return gw, nil
}

func (pf *Client) DeleteGateway(ctx context.Context, name string) error {
	defer pf.write(&pf.mutexes.Gateway)()

	// Use getGatewayEditID to get the index into get_gateways(GW_CACHE_INDEXED),
	// which is what system_gateways.php uses for delete operations (not config.xml index).
	editID, err := pf.getGatewayEditID(ctx, name)
	if err != nil {
		return fmt.Errorf("%w gateway, %w", ErrGetOperationFailed, err)
	}

	relativeURL := url.URL{Path: "system_gateways.php"}
	values := url.Values{
		"act": {"del"},
		"id":  {strconv.Itoa(*editID)},
	}

	_, err = pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return fmt.Errorf("%w gateway, %w", ErrDeleteOperationFailed, err)
	}

	gateways, err := pf.getGateways(ctx)
	if err != nil {
		return fmt.Errorf("%w gateways after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := gateways.GetByName(name); err == nil {
		return fmt.Errorf("%w gateway, still exists", ErrDeleteOperationFailed)
	}

	return nil
}

// ensureDefaultGateway ensures that the system-level default gateway config
// (defaultgw4 / defaultgw6) matches the desired state for the given gateway.
// When setDefault is true and the form POST to save_gateway() already handled
// it, this is a no-op (save_gateway sets defaultgw4/6 when defaultgw=yes).
// When setDefault is false and the gateway IS currently the system default,
// this clears the system default to "" (Automatic).
func (pf *Client) ensureDefaultGateway(ctx context.Context, name string, ipProtocol string, setDefault bool) error {
	if setDefault {
		// save_gateway() already sets defaultgw4/defaultgw6 when defaultgw=yes
		// is in the form POST. Nothing else to do.
		return nil
	}

	// Need to check if this gateway IS currently the system default and clear it.
	configKey := "defaultgw4"
	if ipProtocol == "inet6" {
		configKey = "defaultgw6"
	}

	// Read the current system default, clear it only if it matches this gateway.
	command := fmt.Sprintf(
		"$key = '%s';"+
			"$current = config_get_path('gateways/' . $key, '');"+
			"if ($current === '%s') {"+
			"config_set_path('gateways/' . $key, '');"+
			"write_config('Terraform: cleared default gateway');"+
			"print(json_encode(true));"+
			"} else {"+
			"print(json_encode(false));"+
			"}",
		configKey, name,
	)

	var changed bool
	if err := pf.executePHPCommand(ctx, command, &changed); err != nil {
		return fmt.Errorf("failed to clear default gateway for %s, %w", name, err)
	}

	return nil
}

func (pf *Client) ApplyGatewayChanges(ctx context.Context) error {
	pf.mutexes.GatewayApply.Lock()
	defer pf.mutexes.GatewayApply.Unlock()

	command := "require_once(\"filter.inc\");" +
		"$retval = 0;" +
		"$retval |= system_routing_configure();" +
		"$retval |= system_resolvconf_generate();" +
		"$retval |= filter_configure();" +
		"setup_gateways_monitor();" +
		"send_event(\"service reload dyndnsall\");" +
		"if ($retval == 0) clear_subsystem_dirty('staticroutes');" +
		"print(json_encode($retval));"

	var result int
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply gateway changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}
