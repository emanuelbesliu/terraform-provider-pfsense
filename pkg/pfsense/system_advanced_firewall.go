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
	DefaultAdvancedFirewallOptimization      = "normal"
	DefaultAdvancedFirewallMaxMSS            = 1400
	DefaultAdvancedFirewallBogonsInterval    = "monthly"
	DefaultAdvancedFirewallNATReflection     = "proxy"
	DefaultAdvancedFirewallReflectionTimeout = 2000
)

// AdvancedFirewall represents the System > Advanced > Firewall & NAT configuration.
type AdvancedFirewall struct {
	// Packet Processing
	ScrubNoDF           bool
	ScrubRNID           bool
	Optimization        string
	DisableScrub        bool
	AdaptiveStart       int // 0 = default (60% of max states)
	AdaptiveEnd         int // 0 = default (120% of max states)
	MaximumStates       int // 0 = system default
	MaximumTableEntries int // 0 = system default
	MaximumFragments    int // 0 = default (5000)

	// VPN Packet Processing
	VPNScrubNoDF          bool
	VPNFragmentReassemble bool
	MaxMSSEnable          bool
	MaxMSS                int // default 1400, range 576-65535

	// Advanced Options
	DisableFirewall        bool
	BypassStaticRoutes     bool // NOTE: stored in filter/ not system/
	DisableVPNRules        bool
	DisableReplyTo         bool
	DisableNegate          bool
	NoAPIPA                bool
	AliasesResolveInterval int // 0 = default (300s)
	CheckAliasesURLCert    bool

	// Bogon Networks
	BogonsInterval string // monthly, weekly, daily

	// NAT
	NATReflection             string // disable, proxy, purenat
	ReflectionTimeout         int    // 0 = default (2000), only for proxy mode
	EnableBINATReflection     bool
	EnableNATReflectionHelper bool
	TFTPInterface             string // comma-separated interface list

	// State Timeouts (all 0 = default)
	TCPFirstTimeout       int
	TCPOpeningTimeout     int
	TCPEstablishedTimeout int
	TCPClosingTimeout     int
	TCPFinWaitTimeout     int
	TCPClosedTimeout      int
	TCPTSDiffTimeout      int
	UDPFirstTimeout       int
	UDPSingleTimeout      int
	UDPMultipleTimeout    int
	ICMPFirstTimeout      int
	ICMPErrorTimeout      int
	OtherFirstTimeout     int
	OtherSingleTimeout    int
	OtherMultipleTimeout  int
}

func (a *AdvancedFirewall) SetOptimization(opt string) error {
	valid := []string{"normal", "high-latency", "aggressive", "conservative"}
	for _, v := range valid {
		if opt == v {
			a.Optimization = opt

			return nil
		}
	}

	return fmt.Errorf("%w, optimization must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedFirewall) OptimizationOptions() []string {
	return []string{"normal", "high-latency", "aggressive", "conservative"}
}

func (a *AdvancedFirewall) SetBogonsInterval(interval string) error {
	valid := []string{"monthly", "weekly", "daily"}
	for _, v := range valid {
		if interval == v {
			a.BogonsInterval = interval

			return nil
		}
	}

	return fmt.Errorf("%w, bogons interval must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedFirewall) BogonsIntervalOptions() []string {
	return []string{"monthly", "weekly", "daily"}
}

func (a *AdvancedFirewall) SetNATReflection(mode string) error {
	valid := []string{"disable", "proxy", "purenat"}
	for _, v := range valid {
		if mode == v {
			a.NATReflection = mode

			return nil
		}
	}

	return fmt.Errorf("%w, NAT reflection mode must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedFirewall) NATReflectionOptions() []string {
	return []string{"disable", "proxy", "purenat"}
}

func (a *AdvancedFirewall) SetMaxMSS(mss int) error {
	if mss < 576 || mss > 65535 {
		return fmt.Errorf("%w, max MSS must be between 576 and 65535", ErrClientValidation)
	}

	a.MaxMSS = mss

	return nil
}

// advancedFirewallResponse is the JSON shape returned by the PHP read command.
type advancedFirewallResponse struct {
	// Packet Processing
	ScrubNoDF           json.RawMessage `json:"scrubnodf"`
	ScrubRNID           json.RawMessage `json:"scrubrnid"`
	Optimization        json.RawMessage `json:"optimization"`
	DisableScrub        json.RawMessage `json:"disablescrub"`
	AdaptiveStart       json.RawMessage `json:"adaptivestart"`
	AdaptiveEnd         json.RawMessage `json:"adaptiveend"`
	MaximumStates       json.RawMessage `json:"maximumstates"`
	MaximumTableEntries json.RawMessage `json:"maximumtableentries"`
	MaximumFragments    json.RawMessage `json:"maximumfrags"`

	// VPN Packet Processing
	VPNScrubNoDF          json.RawMessage `json:"vpn_scrubnodf"`
	VPNFragmentReassemble json.RawMessage `json:"vpn_fragment_reassemble"`
	MaxMSSEnable          json.RawMessage `json:"maxmss_enable"`
	MaxMSS                json.RawMessage `json:"maxmss"`

	// Advanced Options
	DisableFirewall        json.RawMessage `json:"disablefilter"`
	BypassStaticRoutes     json.RawMessage `json:"bypassstaticroutes"`
	DisableVPNRules        json.RawMessage `json:"disablevpnrules"`
	DisableReplyTo         json.RawMessage `json:"disablereplyto"`
	DisableNegate          json.RawMessage `json:"disablenegate"`
	NoAPIPA                json.RawMessage `json:"no_apipa_block"`
	AliasesResolveInterval json.RawMessage `json:"aliasesresolveinterval"`
	CheckAliasesURLCert    json.RawMessage `json:"checkaliasesurlcert"`

	// Bogon Networks
	BogonsInterval json.RawMessage `json:"bogonsinterval"`

	// NAT - composite fields
	DisableNATReflection      json.RawMessage `json:"disablenatreflection"`
	EnableNATReflectionPure   json.RawMessage `json:"enablenatreflectionpurenat"`
	ReflectionTimeout         json.RawMessage `json:"reflectiontimeout"`
	EnableBINATReflection     json.RawMessage `json:"enablebinatreflection"`
	EnableNATReflectionHelper json.RawMessage `json:"enablenatreflectionhelper"`
	TFTPInterface             json.RawMessage `json:"tftpinterface"`

	// State Timeouts
	TCPFirstTimeout       json.RawMessage `json:"tcpfirsttimeout"`
	TCPOpeningTimeout     json.RawMessage `json:"tcpopeningtimeout"`
	TCPEstablishedTimeout json.RawMessage `json:"tcpestablishedtimeout"`
	TCPClosingTimeout     json.RawMessage `json:"tcpclosingtimeout"`
	TCPFinWaitTimeout     json.RawMessage `json:"tcpfinwaittimeout"`
	TCPClosedTimeout      json.RawMessage `json:"tcpclosedtimeout"`
	TCPTSDiffTimeout      json.RawMessage `json:"tcptsdifftimeout"`
	UDPFirstTimeout       json.RawMessage `json:"udpfirsttimeout"`
	UDPSingleTimeout      json.RawMessage `json:"udpsingletimeout"`
	UDPMultipleTimeout    json.RawMessage `json:"udpmultipletimeout"`
	ICMPFirstTimeout      json.RawMessage `json:"icmpfirsttimeout"`
	ICMPErrorTimeout      json.RawMessage `json:"icmperrortimeout"`
	OtherFirstTimeout     json.RawMessage `json:"otherfirsttimeout"`
	OtherSingleTimeout    json.RawMessage `json:"othersingletimeout"`
	OtherMultipleTimeout  json.RawMessage `json:"othermultipletimeout"`
}

// rawToInt parses a json.RawMessage as an integer.
// Returns 0 for null/empty/absent values.
func rawToInt(raw json.RawMessage) int {
	s := strings.Trim(string(raw), `"`)
	if s == "" || s == "null" {
		return 0
	}

	v, _ := strconv.Atoi(s)

	return v
}

// rawToString parses a json.RawMessage as a string.
// Returns empty string for null/absent values.
func rawToString(raw json.RawMessage) string {
	s := strings.Trim(string(raw), `"`)
	if s == "null" {
		return ""
	}

	return s
}

func parseAdvancedFirewallResponse(resp advancedFirewallResponse) (AdvancedFirewall, error) {
	var a AdvancedFirewall

	// Packet Processing
	a.ScrubNoDF = rawIsPresent(resp.ScrubNoDF)
	a.ScrubRNID = rawIsPresent(resp.ScrubRNID)

	opt := rawToString(resp.Optimization)
	if opt == "" {
		opt = DefaultAdvancedFirewallOptimization
	}

	if err := a.SetOptimization(opt); err != nil {
		return a, err
	}

	a.DisableScrub = rawIsPresent(resp.DisableScrub)
	a.AdaptiveStart = rawToInt(resp.AdaptiveStart)
	a.AdaptiveEnd = rawToInt(resp.AdaptiveEnd)
	a.MaximumStates = rawToInt(resp.MaximumStates)
	a.MaximumTableEntries = rawToInt(resp.MaximumTableEntries)
	a.MaximumFragments = rawToInt(resp.MaximumFragments)

	// VPN Packet Processing
	a.VPNScrubNoDF = rawIsPresent(resp.VPNScrubNoDF)
	a.VPNFragmentReassemble = rawIsPresent(resp.VPNFragmentReassemble)
	a.MaxMSSEnable = rawIsPresent(resp.MaxMSSEnable)

	maxmss := rawToInt(resp.MaxMSS)
	if maxmss == 0 {
		maxmss = DefaultAdvancedFirewallMaxMSS
	}

	a.MaxMSS = maxmss

	// Advanced Options
	a.DisableFirewall = rawIsPresent(resp.DisableFirewall)
	a.BypassStaticRoutes = rawIsPresent(resp.BypassStaticRoutes)
	a.DisableVPNRules = rawIsPresent(resp.DisableVPNRules)
	a.DisableReplyTo = rawIsPresent(resp.DisableReplyTo)
	a.DisableNegate = rawIsPresent(resp.DisableNegate)
	a.NoAPIPA = rawIsPresent(resp.NoAPIPA)
	a.AliasesResolveInterval = rawToInt(resp.AliasesResolveInterval)
	a.CheckAliasesURLCert = rawIsPresent(resp.CheckAliasesURLCert)

	// Bogon Networks
	bogons := rawToString(resp.BogonsInterval)
	if bogons == "" {
		bogons = DefaultAdvancedFirewallBogonsInterval
	}

	if err := a.SetBogonsInterval(bogons); err != nil {
		return a, err
	}

	// NAT - composite natreflection logic
	disableNAT := rawIsPresent(resp.DisableNATReflection)
	enablePure := rawIsPresent(resp.EnableNATReflectionPure)

	if disableNAT {
		a.NATReflection = "disable"
	} else if enablePure {
		a.NATReflection = "purenat"
	} else {
		a.NATReflection = "proxy"
	}

	a.ReflectionTimeout = rawToInt(resp.ReflectionTimeout)
	a.EnableBINATReflection = rawIsPresent(resp.EnableBINATReflection)
	a.EnableNATReflectionHelper = rawIsPresent(resp.EnableNATReflectionHelper)
	a.TFTPInterface = rawToString(resp.TFTPInterface)

	// State Timeouts
	a.TCPFirstTimeout = rawToInt(resp.TCPFirstTimeout)
	a.TCPOpeningTimeout = rawToInt(resp.TCPOpeningTimeout)
	a.TCPEstablishedTimeout = rawToInt(resp.TCPEstablishedTimeout)
	a.TCPClosingTimeout = rawToInt(resp.TCPClosingTimeout)
	a.TCPFinWaitTimeout = rawToInt(resp.TCPFinWaitTimeout)
	a.TCPClosedTimeout = rawToInt(resp.TCPClosedTimeout)
	a.TCPTSDiffTimeout = rawToInt(resp.TCPTSDiffTimeout)
	a.UDPFirstTimeout = rawToInt(resp.UDPFirstTimeout)
	a.UDPSingleTimeout = rawToInt(resp.UDPSingleTimeout)
	a.UDPMultipleTimeout = rawToInt(resp.UDPMultipleTimeout)
	a.ICMPFirstTimeout = rawToInt(resp.ICMPFirstTimeout)
	a.ICMPErrorTimeout = rawToInt(resp.ICMPErrorTimeout)
	a.OtherFirstTimeout = rawToInt(resp.OtherFirstTimeout)
	a.OtherSingleTimeout = rawToInt(resp.OtherSingleTimeout)
	a.OtherMultipleTimeout = rawToInt(resp.OtherMultipleTimeout)

	return a, nil
}

func (pf *Client) getAdvancedFirewall(ctx context.Context) (*AdvancedFirewall, error) {
	command := "$sys = config_get_path('system', array());" +
		"$filter = config_get_path('filter', array());" +
		"$bogons = isset($sys['bogons']) && isset($sys['bogons']['interval']) ? $sys['bogons']['interval'] : null;" +
		"$out = array(" +
		"'scrubnodf' => isset($sys['scrubnodf']) ? $sys['scrubnodf'] : null," +
		"'scrubrnid' => isset($sys['scrubrnid']) ? $sys['scrubrnid'] : null," +
		"'optimization' => isset($sys['optimization']) ? $sys['optimization'] : null," +
		"'disablescrub' => isset($sys['disablescrub']) ? $sys['disablescrub'] : null," +
		"'adaptivestart' => isset($sys['adaptivestart']) ? $sys['adaptivestart'] : null," +
		"'adaptiveend' => isset($sys['adaptiveend']) ? $sys['adaptiveend'] : null," +
		"'maximumstates' => isset($sys['maximumstates']) ? $sys['maximumstates'] : null," +
		"'maximumtableentries' => isset($sys['maximumtableentries']) ? $sys['maximumtableentries'] : null," +
		"'maximumfrags' => isset($sys['maximumfrags']) ? $sys['maximumfrags'] : null," +
		"'vpn_scrubnodf' => isset($sys['vpn_scrubnodf']) ? $sys['vpn_scrubnodf'] : null," +
		"'vpn_fragment_reassemble' => isset($sys['vpn_fragment_reassemble']) ? $sys['vpn_fragment_reassemble'] : null," +
		"'maxmss_enable' => isset($sys['maxmss_enable']) ? $sys['maxmss_enable'] : null," +
		"'maxmss' => isset($sys['maxmss']) ? $sys['maxmss'] : null," +
		"'disablefilter' => isset($sys['disablefilter']) ? $sys['disablefilter'] : null," +
		"'bypassstaticroutes' => isset($filter['bypassstaticroutes']) ? $filter['bypassstaticroutes'] : null," +
		"'disablevpnrules' => isset($sys['disablevpnrules']) ? $sys['disablevpnrules'] : null," +
		"'disablereplyto' => isset($sys['disablereplyto']) ? $sys['disablereplyto'] : null," +
		"'disablenegate' => isset($sys['disablenegate']) ? $sys['disablenegate'] : null," +
		"'no_apipa_block' => isset($sys['no_apipa_block']) ? $sys['no_apipa_block'] : null," +
		"'aliasesresolveinterval' => isset($sys['aliasesresolveinterval']) ? $sys['aliasesresolveinterval'] : null," +
		"'checkaliasesurlcert' => isset($sys['checkaliasesurlcert']) ? $sys['checkaliasesurlcert'] : null," +
		"'bogonsinterval' => $bogons," +
		"'disablenatreflection' => isset($sys['disablenatreflection']) ? $sys['disablenatreflection'] : null," +
		"'enablenatreflectionpurenat' => isset($sys['enablenatreflectionpurenat']) ? $sys['enablenatreflectionpurenat'] : null," +
		"'reflectiontimeout' => isset($sys['reflectiontimeout']) ? $sys['reflectiontimeout'] : null," +
		"'enablebinatreflection' => isset($sys['enablebinatreflection']) ? $sys['enablebinatreflection'] : null," +
		"'enablenatreflectionhelper' => isset($sys['enablenatreflectionhelper']) ? $sys['enablenatreflectionhelper'] : null," +
		"'tftpinterface' => isset($sys['tftpinterface']) ? $sys['tftpinterface'] : null," +
		"'tcpfirsttimeout' => isset($sys['tcpfirsttimeout']) ? $sys['tcpfirsttimeout'] : null," +
		"'tcpopeningtimeout' => isset($sys['tcpopeningtimeout']) ? $sys['tcpopeningtimeout'] : null," +
		"'tcpestablishedtimeout' => isset($sys['tcpestablishedtimeout']) ? $sys['tcpestablishedtimeout'] : null," +
		"'tcpclosingtimeout' => isset($sys['tcpclosingtimeout']) ? $sys['tcpclosingtimeout'] : null," +
		"'tcpfinwaittimeout' => isset($sys['tcpfinwaittimeout']) ? $sys['tcpfinwaittimeout'] : null," +
		"'tcpclosedtimeout' => isset($sys['tcpclosedtimeout']) ? $sys['tcpclosedtimeout'] : null," +
		"'tcptsdifftimeout' => isset($sys['tcptsdifftimeout']) ? $sys['tcptsdifftimeout'] : null," +
		"'udpfirsttimeout' => isset($sys['udpfirsttimeout']) ? $sys['udpfirsttimeout'] : null," +
		"'udpsingletimeout' => isset($sys['udpsingletimeout']) ? $sys['udpsingletimeout'] : null," +
		"'udpmultipletimeout' => isset($sys['udpmultipletimeout']) ? $sys['udpmultipletimeout'] : null," +
		"'icmpfirsttimeout' => isset($sys['icmpfirsttimeout']) ? $sys['icmpfirsttimeout'] : null," +
		"'icmperrortimeout' => isset($sys['icmperrortimeout']) ? $sys['icmperrortimeout'] : null," +
		"'otherfirsttimeout' => isset($sys['otherfirsttimeout']) ? $sys['otherfirsttimeout'] : null," +
		"'othersingletimeout' => isset($sys['othersingletimeout']) ? $sys['othersingletimeout'] : null," +
		"'othermultipletimeout' => isset($sys['othermultipletimeout']) ? $sys['othermultipletimeout'] : null" +
		");" +
		"print(json_encode($out));"

	var resp advancedFirewallResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	a, err := parseAdvancedFirewallResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("%w advanced firewall response, %w", ErrUnableToParse, err)
	}

	return &a, nil
}

func (pf *Client) GetAdvancedFirewall(ctx context.Context) (*AdvancedFirewall, error) {
	defer pf.read(&pf.mutexes.AdvancedFirewall)()

	a, err := pf.getAdvancedFirewall(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w advanced firewall, %w", ErrGetOperationFailed, err)
	}

	return a, nil
}

func advancedFirewallFormValues(a AdvancedFirewall) url.Values {
	values := url.Values{
		"save": {"Save"},
	}

	// Packet Processing
	if a.ScrubNoDF {
		values.Set("scrubnodf", "yes")
	}

	if a.ScrubRNID {
		values.Set("scrubrnid", "yes")
	}

	values.Set("optimization", a.Optimization)

	if a.DisableScrub {
		values.Set("disablescrub", "yes")
	}

	if a.AdaptiveStart > 0 {
		values.Set("adaptivestart", strconv.Itoa(a.AdaptiveStart))
	}

	if a.AdaptiveEnd > 0 {
		values.Set("adaptiveend", strconv.Itoa(a.AdaptiveEnd))
	}

	if a.MaximumStates > 0 {
		values.Set("maximumstates", strconv.Itoa(a.MaximumStates))
	}

	if a.MaximumTableEntries > 0 {
		values.Set("maximumtableentries", strconv.Itoa(a.MaximumTableEntries))
	}

	if a.MaximumFragments > 0 {
		values.Set("maximumfrags", strconv.Itoa(a.MaximumFragments))
	}

	// VPN Packet Processing
	if a.VPNScrubNoDF {
		values.Set("vpn_scrubnodf", "yes")
	}

	if a.VPNFragmentReassemble {
		values.Set("vpn_fragment_reassemble", "yes")
	}

	if a.MaxMSSEnable {
		values.Set("maxmss_enable", "yes")
	}

	if a.MaxMSS > 0 {
		values.Set("maxmss", strconv.Itoa(a.MaxMSS))
	}

	// Advanced Options
	if a.DisableFirewall {
		values.Set("disablefilter", "yes")
	}

	if a.BypassStaticRoutes {
		values.Set("bypassstaticroutes", "yes")
	}

	if a.DisableVPNRules {
		values.Set("disablevpnrules", "yes")
	}

	if a.DisableReplyTo {
		values.Set("disablereplyto", "yes")
	}

	if a.DisableNegate {
		values.Set("disablenegate", "yes")
	}

	if a.NoAPIPA {
		values.Set("no_apipa_block", "yes")
	}

	if a.AliasesResolveInterval > 0 {
		values.Set("aliasesresolveinterval", strconv.Itoa(a.AliasesResolveInterval))
	}

	if a.CheckAliasesURLCert {
		values.Set("checkaliasesurlcert", "yes")
	}

	// Bogon Networks
	values.Set("bogonsinterval", a.BogonsInterval)

	// NAT - composite natreflection
	switch a.NATReflection {
	case "disable":
		values.Set("natreflection", "disable")
	case "proxy":
		values.Set("natreflection", "proxy")
	case "purenat":
		values.Set("natreflection", "purenat")
	}

	if a.ReflectionTimeout > 0 {
		values.Set("reflectiontimeout", strconv.Itoa(a.ReflectionTimeout))
	}

	if a.EnableBINATReflection {
		values.Set("enablebinatreflection", "yes")
	}

	if a.EnableNATReflectionHelper {
		values.Set("enablenatreflectionhelper", "yes")
	}

	if a.TFTPInterface != "" {
		// tftpinterface is a multi-select — POST as individual values
		interfaces := strings.Split(a.TFTPInterface, ",")
		for _, iface := range interfaces {
			iface = strings.TrimSpace(iface)
			if iface != "" {
				values.Add("tftpinterface[]", iface)
			}
		}
	}

	// State Timeouts — only send if non-zero
	setIfPositive := func(key string, val int) {
		if val > 0 {
			values.Set(key, strconv.Itoa(val))
		}
	}

	setIfPositive("tcpfirsttimeout", a.TCPFirstTimeout)
	setIfPositive("tcpopeningtimeout", a.TCPOpeningTimeout)
	setIfPositive("tcpestablishedtimeout", a.TCPEstablishedTimeout)
	setIfPositive("tcpclosingtimeout", a.TCPClosingTimeout)
	setIfPositive("tcpfinwaittimeout", a.TCPFinWaitTimeout)
	setIfPositive("tcpclosedtimeout", a.TCPClosedTimeout)
	setIfPositive("tcptsdifftimeout", a.TCPTSDiffTimeout)
	setIfPositive("udpfirsttimeout", a.UDPFirstTimeout)
	setIfPositive("udpsingletimeout", a.UDPSingleTimeout)
	setIfPositive("udpmultipletimeout", a.UDPMultipleTimeout)
	setIfPositive("icmpfirsttimeout", a.ICMPFirstTimeout)
	setIfPositive("icmperrortimeout", a.ICMPErrorTimeout)
	setIfPositive("otherfirsttimeout", a.OtherFirstTimeout)
	setIfPositive("othersingletimeout", a.OtherSingleTimeout)
	setIfPositive("othermultipletimeout", a.OtherMultipleTimeout)

	return values
}

func (pf *Client) UpdateAdvancedFirewall(ctx context.Context, a AdvancedFirewall) (*AdvancedFirewall, error) {
	defer pf.write(&pf.mutexes.AdvancedFirewall)()

	relativeURL := url.URL{Path: "system_advanced_firewall.php"}
	values := advancedFirewallFormValues(a)

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return nil, fmt.Errorf("%w advanced firewall, %w", ErrUpdateOperationFailed, err)
	}

	if err := scrapeHTMLValidationErrors(doc); err != nil {
		return nil, fmt.Errorf("%w advanced firewall, %w", ErrUpdateOperationFailed, err)
	}

	result, err := pf.getAdvancedFirewall(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w advanced firewall after updating, %w", ErrGetOperationFailed, err)
	}

	return result, nil
}

func (pf *Client) ApplyAdvancedFirewallChanges(ctx context.Context) error {
	pf.mutexes.AdvancedFirewallApply.Lock()
	defer pf.mutexes.AdvancedFirewallApply.Unlock()

	// Mirrors saveSystemAdvancedFirewall() in system_advanced_firewall.inc:
	// filter_configure(), setup_loader_settings(), system_setup_sysctl(), filterdns restart, cron update
	command := "$retval = 0;" +
		"$retval |= filter_configure();" +
		"if (is_module_loaded('pf')) {" +
		"setup_loader_settings();" +
		"}" +
		"system_setup_sysctl();" +
		"print(json_encode($retval));"

	var result int
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply advanced firewall changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}
