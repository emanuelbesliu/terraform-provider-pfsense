package pfsense

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	DefaultDNSResolverGeneralEnable                  = true
	DefaultDNSResolverGeneralPort                    = 53
	DefaultDNSResolverGeneralEnableSSL               = false
	DefaultDNSResolverGeneralTLSPort                 = 853
	DefaultDNSResolverGeneralActiveInterface         = "all"
	DefaultDNSResolverGeneralOutgoingInterface       = "all"
	DefaultDNSResolverGeneralSystemDomainLocalZone   = "transparent"
	DefaultDNSResolverGeneralDNSSEC                  = false
	DefaultDNSResolverGeneralForwarding              = false
	DefaultDNSResolverGeneralForwardTLSUpstream      = false
	DefaultDNSResolverGeneralRegisterDHCPLeases      = false
	DefaultDNSResolverGeneralRegisterDHCPStaticMaps  = false
	DefaultDNSResolverGeneralRegisterOpenVPNClients  = false
	DefaultDNSResolverGeneralStrictOutgoingInterface = false
	DefaultDNSResolverGeneralPython                  = false
)

// DNSResolverGeneral represents the Services > DNS Resolver general settings.
type DNSResolverGeneral struct {
	Enable                  bool
	Port                    int
	EnableSSL               bool
	TLSPort                 int
	SSLCertRef              string
	ActiveInterfaces        []string
	OutgoingInterfaces      []string
	SystemDomainLocalZone   string
	DNSSEC                  bool
	Forwarding              bool
	ForwardTLSUpstream      bool
	RegisterDHCPLeases      bool
	RegisterDHCPStaticMaps  bool
	RegisterOpenVPNClients  bool
	StrictOutgoingInterface bool
	Python                  bool
	PythonOrder             string
	PythonScript            string
	CustomOptions           string
}

func (d DNSResolverGeneral) SystemDomainLocalZoneOptions() []string {
	return []string{
		"transparent",
		"always_transparent",
		"typetransparent",
		"redirect",
		"inform",
		"inform_deny",
		"inform_redirect",
		"always_refuse",
		"always_nxdomain",
		"always_deny",
		"always_null",
		"nodefault",
		"deny",
		"refuse",
		"static",
	}
}

// dnsResolverGeneralResponse is the JSON shape returned by the PHP read command.
type dnsResolverGeneralResponse struct {
	Enable                  *string `json:"enable"`
	Port                    string  `json:"port"`
	EnableSSL               *string `json:"enablessl"`
	TLSPort                 string  `json:"tlsport"`
	SSLCertRef              string  `json:"sslcertref"`
	ActiveInterface         string  `json:"active_interface"`
	OutgoingInterface       string  `json:"outgoing_interface"`
	SystemDomainLocalZone   string  `json:"system_domain_local_zone_type"`
	DNSSEC                  *string `json:"dnssec"`
	Forwarding              *string `json:"forwarding"`
	ForwardTLSUpstream      *string `json:"forward_tls_upstream"`
	RegisterDHCPLeases      *string `json:"regdhcp"`
	RegisterDHCPStaticMaps  *string `json:"regdhcpstatic"`
	RegisterOpenVPNClients  *string `json:"regovpnclients"`
	StrictOutgoingInterface *string `json:"strictout"`
	Python                  *string `json:"python"`
	PythonOrder             string  `json:"python_order"`
	PythonScript            string  `json:"python_script"`
	CustomOptions           string  `json:"custom_options"`
}

func parseDNSResolverGeneralResponse(resp dnsResolverGeneralResponse) DNSResolverGeneral {
	var dg DNSResolverGeneral

	// Presence-based booleans: pointer != nil means the key was present in JSON (enabled)
	dg.Enable = resp.Enable != nil
	dg.EnableSSL = resp.EnableSSL != nil
	dg.DNSSEC = resp.DNSSEC != nil
	dg.Forwarding = resp.Forwarding != nil
	dg.ForwardTLSUpstream = resp.ForwardTLSUpstream != nil
	dg.RegisterDHCPLeases = resp.RegisterDHCPLeases != nil
	dg.RegisterDHCPStaticMaps = resp.RegisterDHCPStaticMaps != nil
	dg.RegisterOpenVPNClients = resp.RegisterOpenVPNClients != nil
	dg.StrictOutgoingInterface = resp.StrictOutgoingInterface != nil
	dg.Python = resp.Python != nil

	// Port
	dg.Port = DefaultDNSResolverGeneralPort
	if resp.Port != "" {
		p := 0
		if _, err := fmt.Sscanf(resp.Port, "%d", &p); err == nil && p > 0 {
			dg.Port = p
		}
	}

	// TLS Port
	dg.TLSPort = DefaultDNSResolverGeneralTLSPort
	if resp.TLSPort != "" {
		p := 0
		if _, err := fmt.Sscanf(resp.TLSPort, "%d", &p); err == nil && p > 0 {
			dg.TLSPort = p
		}
	}

	// SSL cert ref
	dg.SSLCertRef = resp.SSLCertRef

	// Interfaces — stored as comma-separated in XML config
	if resp.ActiveInterface != "" {
		dg.ActiveInterfaces = strings.Split(resp.ActiveInterface, ",")
	} else {
		dg.ActiveInterfaces = []string{"all"}
	}

	if resp.OutgoingInterface != "" {
		dg.OutgoingInterfaces = strings.Split(resp.OutgoingInterface, ",")
	} else {
		dg.OutgoingInterfaces = []string{"all"}
	}

	// System domain local zone type
	dg.SystemDomainLocalZone = resp.SystemDomainLocalZone
	if dg.SystemDomainLocalZone == "" {
		dg.SystemDomainLocalZone = DefaultDNSResolverGeneralSystemDomainLocalZone
	}

	// Python fields
	dg.PythonOrder = resp.PythonOrder
	dg.PythonScript = resp.PythonScript

	// Custom options — stored as base64 in XML, the PHP command decodes it for us
	dg.CustomOptions = resp.CustomOptions

	return dg
}

func (pf *Client) getDNSResolverGeneral(ctx context.Context) (*DNSResolverGeneral, error) {
	command := "$ub = config_get_path('unbound', array());" +
		"$out = array(" +
		"'enable' => array_key_exists('enable', $ub) ? $ub['enable'] : null," +
		"'port' => isset($ub['port']) ? $ub['port'] : ''," +
		"'enablessl' => array_key_exists('enablessl', $ub) ? $ub['enablessl'] : null," +
		"'tlsport' => isset($ub['tlsport']) ? $ub['tlsport'] : ''," +
		"'sslcertref' => isset($ub['sslcertref']) ? $ub['sslcertref'] : ''," +
		"'active_interface' => isset($ub['active_interface']) ? $ub['active_interface'] : ''," +
		"'outgoing_interface' => isset($ub['outgoing_interface']) ? $ub['outgoing_interface'] : ''," +
		"'system_domain_local_zone_type' => isset($ub['system_domain_local_zone_type']) ? $ub['system_domain_local_zone_type'] : ''," +
		"'dnssec' => array_key_exists('dnssec', $ub) ? $ub['dnssec'] : null," +
		"'forwarding' => array_key_exists('forwarding', $ub) ? $ub['forwarding'] : null," +
		"'forward_tls_upstream' => array_key_exists('forward_tls_upstream', $ub) ? $ub['forward_tls_upstream'] : null," +
		"'regdhcp' => array_key_exists('regdhcp', $ub) ? $ub['regdhcp'] : null," +
		"'regdhcpstatic' => array_key_exists('regdhcpstatic', $ub) ? $ub['regdhcpstatic'] : null," +
		"'regovpnclients' => array_key_exists('regovpnclients', $ub) ? $ub['regovpnclients'] : null," +
		"'strictout' => array_key_exists('strictout', $ub) ? $ub['strictout'] : null," +
		"'python' => array_key_exists('python', $ub) ? $ub['python'] : null," +
		"'python_order' => isset($ub['python_order']) ? $ub['python_order'] : ''," +
		"'python_script' => isset($ub['python_script']) ? $ub['python_script'] : ''," +
		"'custom_options' => isset($ub['custom_options']) ? base64_decode($ub['custom_options']) : ''" +
		");" +
		"print(json_encode($out));"

	var resp dnsResolverGeneralResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	dg := parseDNSResolverGeneralResponse(resp)

	return &dg, nil
}

func (pf *Client) GetDNSResolverGeneral(ctx context.Context) (*DNSResolverGeneral, error) {
	defer pf.read(&pf.mutexes.DNSResolverGeneral)()

	dg, err := pf.getDNSResolverGeneral(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w dns resolver general, %w", ErrGetOperationFailed, err)
	}

	return dg, nil
}

func dnsResolverGeneralFormValues(dg DNSResolverGeneral) url.Values {
	values := url.Values{
		"save": {"Save"},
	}

	// Enable
	if dg.Enable {
		values.Set("enable", "yes")
	}

	// Port
	if dg.Port != DefaultDNSResolverGeneralPort {
		values.Set("port", fmt.Sprintf("%d", dg.Port))
	} else {
		values.Set("port", "")
	}

	// SSL/TLS
	if dg.EnableSSL {
		values.Set("enablessl", "yes")
	}

	if dg.TLSPort != DefaultDNSResolverGeneralTLSPort {
		values.Set("tlsport", fmt.Sprintf("%d", dg.TLSPort))
	} else {
		values.Set("tlsport", "")
	}

	if dg.SSLCertRef != "" {
		values.Set("sslcertref", dg.SSLCertRef)
	}

	// Interfaces — multi-select fields use [] suffix
	for _, iface := range dg.ActiveInterfaces {
		values.Add("active_interface[]", iface)
	}

	for _, iface := range dg.OutgoingInterfaces {
		values.Add("outgoing_interface[]", iface)
	}

	// System domain local zone type
	values.Set("system_domain_local_zone_type", dg.SystemDomainLocalZone)

	// DNSSEC
	if dg.DNSSEC {
		values.Set("dnssec", "yes")
	}

	// Forwarding
	if dg.Forwarding {
		values.Set("forwarding", "yes")
	}

	if dg.ForwardTLSUpstream {
		values.Set("forward_tls_upstream", "yes")
	}

	// DHCP registration
	if dg.RegisterDHCPLeases {
		values.Set("regdhcp", "yes")
	}

	if dg.RegisterDHCPStaticMaps {
		values.Set("regdhcpstatic", "yes")
	}

	// OpenVPN clients
	if dg.RegisterOpenVPNClients {
		values.Set("regovpnclients", "yes")
	}

	// Strict outgoing interface
	if dg.StrictOutgoingInterface {
		values.Set("strictout", "yes")
	}

	// Python module
	if dg.Python {
		values.Set("python", "yes")
	}

	if dg.PythonOrder != "" {
		values.Set("python_order", dg.PythonOrder)
	}

	if dg.PythonScript != "" {
		values.Set("python_script", dg.PythonScript)
	}

	// Custom options
	if dg.CustomOptions != "" {
		values.Set("custom_options", dg.CustomOptions)
	}

	return values
}

func (pf *Client) UpdateDNSResolverGeneral(ctx context.Context, dg DNSResolverGeneral) (*DNSResolverGeneral, error) {
	defer pf.write(&pf.mutexes.DNSResolverGeneral)()

	relativeURL := url.URL{Path: "services_unbound.php"}
	values := dnsResolverGeneralFormValues(dg)

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return nil, fmt.Errorf("%w dns resolver general, %w", ErrUpdateOperationFailed, err)
	}

	if err := scrapeHTMLValidationErrors(doc); err != nil {
		return nil, fmt.Errorf("%w dns resolver general, %w", ErrUpdateOperationFailed, err)
	}

	result, err := pf.getDNSResolverGeneral(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w dns resolver general after updating, %w", ErrGetOperationFailed, err)
	}

	return result, nil
}
