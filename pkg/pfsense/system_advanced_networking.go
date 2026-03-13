package pfsense

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	DefaultAdvancedNetworkingDHCPBackend = "isc"
)

// AdvancedNetworking represents the System > Advanced > Networking configuration.
type AdvancedNetworking struct {
	// DHCP Options
	DHCPBackend      string // "isc" or "kea"
	IgnoreISCWarning bool
	RADVDDebug       bool
	DHCP6Debug       bool
	DHCP6NoRelease   bool
	GlobalV6DUID     string // Raw DUID string
	IPv6DUIDType     int    // 0=raw, 1=LLT, 2=EN, 3=LL, 4=UUID

	// IPv6 Options
	IPv6Allow              bool
	IPv6NATEnable          bool
	IPv6NATIPAddr          string // IPv4 address of tunnel peer
	PreferIPv4             bool
	IPv6DontCreateLocalDNS bool

	// Network Interfaces
	DisableChecksumOffloading     bool
	DisableSegmentationOffloading bool
	DisableLargeReceiveOffloading bool
	HNALTQEnable                  bool
	SharedNet                     bool
	IPChangeKillStates            bool
	UseIfPPPoE                    bool
}

func (a *AdvancedNetworking) SetDHCPBackend(backend string) error {
	valid := []string{"isc", "kea"}
	for _, v := range valid {
		if backend == v {
			a.DHCPBackend = backend

			return nil
		}
	}

	return fmt.Errorf("%w, DHCP backend must be one of: %s", ErrClientValidation, strings.Join(valid, ", "))
}

func (AdvancedNetworking) DHCPBackendOptions() []string {
	return []string{"isc", "kea"}
}

// advancedNetworkingResponse is the JSON shape returned by the PHP read command.
type advancedNetworkingResponse struct {
	// DHCP Options
	DHCPBackend      json.RawMessage `json:"dhcpbackend"`
	IgnoreISCWarning json.RawMessage `json:"ignoreiscwarning"`
	RADVDDebug       json.RawMessage `json:"radvddebug"`
	DHCP6Debug       json.RawMessage `json:"dhcp6debug"`
	DHCP6NoRelease   json.RawMessage `json:"dhcp6norelease"`
	GlobalV6DUID     json.RawMessage `json:"global-v6duid"`

	// IPv6 Options
	IPv6Allow              json.RawMessage `json:"ipv6allow"`
	IPv6NATEnable          json.RawMessage `json:"ipv6nat_enable"`
	IPv6NATIPAddr          json.RawMessage `json:"ipv6nat_ipaddr"`
	PreferIPv4             json.RawMessage `json:"prefer_ipv4"`
	IPv6DontCreateLocalDNS json.RawMessage `json:"ipv6dontcreatelocaldns"`

	// Network Interfaces
	DisableChecksumOffloading     json.RawMessage `json:"disablechecksumoffloading"`
	DisableSegmentationOffloading json.RawMessage `json:"disablesegmentationoffloading"`
	DisableLargeReceiveOffloading json.RawMessage `json:"disablelargereceiveoffloading"`
	HNALTQEnable                  json.RawMessage `json:"hnaltqenable"`
	SharedNet                     json.RawMessage `json:"sharednet"`
	IPChangeKillStates            json.RawMessage `json:"ip_change_kill_states"`
	UseIfPPPoE                    json.RawMessage `json:"use_if_pppoe"`
}

func parseAdvancedNetworkingResponse(resp advancedNetworkingResponse) (AdvancedNetworking, error) {
	var a AdvancedNetworking

	// DHCP Options
	backend := rawToString(resp.DHCPBackend)
	if backend == "" {
		backend = DefaultAdvancedNetworkingDHCPBackend
	}

	if err := a.SetDHCPBackend(backend); err != nil {
		return a, err
	}

	a.IgnoreISCWarning = rawIsPresent(resp.IgnoreISCWarning)
	a.RADVDDebug = rawIsPresent(resp.RADVDDebug)
	a.DHCP6Debug = rawIsPresent(resp.DHCP6Debug)
	a.DHCP6NoRelease = rawIsPresent(resp.DHCP6NoRelease)
	a.GlobalV6DUID = rawToString(resp.GlobalV6DUID)

	// IPv6 Options
	a.IPv6Allow = rawIsPresent(resp.IPv6Allow)
	a.IPv6NATEnable = rawIsPresent(resp.IPv6NATEnable)
	a.IPv6NATIPAddr = rawToString(resp.IPv6NATIPAddr)
	a.PreferIPv4 = rawIsPresent(resp.PreferIPv4)
	a.IPv6DontCreateLocalDNS = rawIsPresent(resp.IPv6DontCreateLocalDNS)

	// Network Interfaces
	a.DisableChecksumOffloading = rawIsPresent(resp.DisableChecksumOffloading)
	a.DisableSegmentationOffloading = rawIsPresent(resp.DisableSegmentationOffloading)
	a.DisableLargeReceiveOffloading = rawIsPresent(resp.DisableLargeReceiveOffloading)
	a.HNALTQEnable = rawIsPresent(resp.HNALTQEnable)
	a.SharedNet = rawIsPresent(resp.SharedNet)
	a.IPChangeKillStates = rawIsPresent(resp.IPChangeKillStates)
	a.UseIfPPPoE = rawIsPresent(resp.UseIfPPPoE)

	return a, nil
}

func (pf *Client) getAdvancedNetworking(ctx context.Context) (*AdvancedNetworking, error) {
	command := "$sys = config_get_path('system', array());" +
		"$diag = config_get_path('diag', array());" +
		"$dhcpbackend = config_get_path('dhcpbackend', 'isc');" +
		"$ipv6nat_enable = isset($diag['ipv6nat']) && isset($diag['ipv6nat']['enable']) ? true : null;" +
		"$ipv6nat_ipaddr = isset($diag['ipv6nat']) && isset($diag['ipv6nat']['ipaddr']) ? $diag['ipv6nat']['ipaddr'] : null;" +
		"$out = array(" +
		"'dhcpbackend' => $dhcpbackend," +
		"'ignoreiscwarning' => isset($sys['ignoreiscwarning']) ? $sys['ignoreiscwarning'] : null," +
		"'radvddebug' => isset($sys['radvddebug']) ? $sys['radvddebug'] : null," +
		"'dhcp6debug' => isset($sys['dhcp6debug']) ? $sys['dhcp6debug'] : null," +
		"'dhcp6norelease' => isset($sys['dhcp6norelease']) ? $sys['dhcp6norelease'] : null," +
		"'global-v6duid' => isset($sys['global-v6duid']) ? $sys['global-v6duid'] : null," +
		"'ipv6allow' => isset($sys['ipv6allow']) ? $sys['ipv6allow'] : null," +
		"'ipv6nat_enable' => $ipv6nat_enable," +
		"'ipv6nat_ipaddr' => $ipv6nat_ipaddr," +
		"'prefer_ipv4' => isset($sys['prefer_ipv4']) ? $sys['prefer_ipv4'] : null," +
		"'ipv6dontcreatelocaldns' => isset($sys['ipv6dontcreatelocaldns']) ? $sys['ipv6dontcreatelocaldns'] : null," +
		"'disablechecksumoffloading' => isset($sys['disablechecksumoffloading']) ? $sys['disablechecksumoffloading'] : null," +
		"'disablesegmentationoffloading' => isset($sys['disablesegmentationoffloading']) ? $sys['disablesegmentationoffloading'] : null," +
		"'disablelargereceiveoffloading' => isset($sys['disablelargereceiveoffloading']) ? $sys['disablelargereceiveoffloading'] : null," +
		"'hnaltqenable' => isset($sys['hn_altq_enable']) ? $sys['hn_altq_enable'] : null," +
		"'sharednet' => isset($sys['sharednet']) ? $sys['sharednet'] : null," +
		"'ip_change_kill_states' => isset($sys['ip_change_kill_states']) ? $sys['ip_change_kill_states'] : null," +
		"'use_if_pppoe' => isset($sys['use_if_pppoe']) ? $sys['use_if_pppoe'] : null" +
		");" +
		"print(json_encode($out));"

	var resp advancedNetworkingResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	a, err := parseAdvancedNetworkingResponse(resp)
	if err != nil {
		return nil, fmt.Errorf("%w advanced networking response, %w", ErrUnableToParse, err)
	}

	return &a, nil
}

func (pf *Client) GetAdvancedNetworking(ctx context.Context) (*AdvancedNetworking, error) {
	defer pf.read(&pf.mutexes.AdvancedNetworking)()

	a, err := pf.getAdvancedNetworking(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w advanced networking, %w", ErrGetOperationFailed, err)
	}

	return a, nil
}

func advancedNetworkingFormValues(a AdvancedNetworking) url.Values {
	values := url.Values{
		"save": {"Save"},
	}

	// DHCP Options
	values.Set("dhcpbackend", a.DHCPBackend)

	if a.IgnoreISCWarning {
		values.Set("ignoreiscwarning", "yes")
	}

	if a.RADVDDebug {
		values.Set("radvddebug", "yes")
	}

	if a.DHCP6Debug {
		values.Set("dhcp6debug", "yes")
	}

	if a.DHCP6NoRelease {
		values.Set("dhcp6norelease", "yes")
	}

	// DUID handling — only send raw DUID if set, with type 0 (raw)
	if a.GlobalV6DUID != "" {
		values.Set("ipv6duidtype", "0")
		values.Set("global-v6duid", a.GlobalV6DUID)
	} else {
		values.Set("ipv6duidtype", "0")
	}

	// IPv6 Options
	if a.IPv6Allow {
		values.Set("ipv6allow", "yes")
	}

	if a.IPv6NATEnable {
		values.Set("ipv6nat_enable", "yes")
	}

	if a.IPv6NATIPAddr != "" {
		values.Set("ipv6nat_ipaddr", a.IPv6NATIPAddr)
	}

	if a.PreferIPv4 {
		values.Set("prefer_ipv4", "yes")
	}

	if a.IPv6DontCreateLocalDNS {
		values.Set("ipv6dontcreatelocaldns", "yes")
	}

	// Network Interfaces
	if a.DisableChecksumOffloading {
		values.Set("disablechecksumoffloading", "yes")
	}

	if a.DisableSegmentationOffloading {
		values.Set("disablesegmentationoffloading", "yes")
	}

	if a.DisableLargeReceiveOffloading {
		values.Set("disablelargereceiveoffloading", "yes")
	}

	if a.HNALTQEnable {
		values.Set("hnaltqenable", "yes")
	}

	if a.SharedNet {
		values.Set("sharednet", "yes")
	}

	if a.IPChangeKillStates {
		values.Set("ip_change_kill_states", "yes")
	}

	if a.UseIfPPPoE {
		values.Set("use_if_pppoe", "yes")
	}

	return values
}

func (pf *Client) UpdateAdvancedNetworking(ctx context.Context, a AdvancedNetworking) (*AdvancedNetworking, error) {
	defer pf.write(&pf.mutexes.AdvancedNetworking)()

	relativeURL := url.URL{Path: "system_advanced_network.php"}
	values := advancedNetworkingFormValues(a)

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return nil, fmt.Errorf("%w advanced networking, %w", ErrUpdateOperationFailed, err)
	}

	if err := scrapeHTMLValidationErrors(doc); err != nil {
		return nil, fmt.Errorf("%w advanced networking, %w", ErrUpdateOperationFailed, err)
	}

	result, err := pf.getAdvancedNetworking(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w advanced networking after updating, %w", ErrGetOperationFailed, err)
	}

	return result, nil
}

func (pf *Client) ApplyAdvancedNetworkingChanges(ctx context.Context) error {
	pf.mutexes.AdvancedNetworkingApply.Lock()
	defer pf.mutexes.AdvancedNetworkingApply.Unlock()

	// Mirrors saveAdvancedNetworking() in system_advanced_network.inc:
	// filter_configure(), prefer_ipv4_or_ipv6(), setup_loader_settings()
	command := "$retval = 0;" +
		"$retval |= filter_configure();" +
		"prefer_ipv4_or_ipv6();" +
		"setup_loader_settings();" +
		"print(json_encode($retval));"

	var result int
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply advanced networking changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}
