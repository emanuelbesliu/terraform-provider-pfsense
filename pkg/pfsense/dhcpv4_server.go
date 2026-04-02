package pfsense

import (
	"context"
	"fmt"
	"net/http"
	"net/netip"
	"net/url"
	"strings"
	"time"
)

// ======================================================================
// Constants
// ======================================================================

const (
	DHCPv4ServerMaxDNSServers  = 4
	DHCPv4ServerMaxNTPServers  = 4
	DHCPv4ServerMaxWINSServers = 2

	dhcpv4ServerDomainSearchListSep = ";"
)

// ======================================================================
// Response struct (JSON shape from PHP)
// ======================================================================

type dhcpv4ServerResponse struct {
	Enable            *string  `json:"enable"`            // presence-based: key exists = enabled
	RangeFrom         string   `json:"range_from"`        // range.from
	RangeTo           string   `json:"range_to"`          // range.to
	DNSServers        []string `json:"dnsserver"`         // array of IP strings
	Gateway           string   `json:"gateway"`           // IP string
	DomainName        string   `json:"domain"`            // domain string
	DomainSearchList  string   `json:"domainsearchlist"`  // semicolon-separated
	DefaultLeaseTime  string   `json:"defaultleasetime"`  // seconds as string
	MaximumLeaseTime  string   `json:"maxleasetime"`      // seconds as string
	WINSServers       []string `json:"winsserver"`        // array of IP strings
	NTPServers        []string `json:"ntpserver"`         // array of IP strings
	TFTPServer        string   `json:"tftp"`              // hostname/IP
	LDAPServer        string   `json:"ldap"`              // URI
	MACAllow          string   `json:"mac_allow"`         // comma-separated MACs
	MACDeny           string   `json:"mac_deny"`          // comma-separated MACs
	DenyUnknown       *string  `json:"denyunknown"`       // presence-based
	IgnoreClientUIDs  *string  `json:"ignoreclientuids"`  // presence-based
	StaticARP         *string  `json:"staticarp"`         // presence-based
	Netmask           string   `json:"netmask"`           // subnet mask override
	FailoverPeerIP    string   `json:"failover_peerip"`   // failover peer IP
	Netboot           *string  `json:"netboot"`           // presence-based
	NextServer        string   `json:"nextserver"`        // PXE next-server
	Filename          string   `json:"filename"`          // PXE default boot file
	Filename32        string   `json:"filename32"`        // PXE 32-bit boot file
	Filename32Arm     string   `json:"filename32arm"`     // PXE 32-bit ARM boot file
	Filename64        string   `json:"filename64"`        // PXE 64-bit boot file
	Filename64Arm     string   `json:"filename64arm"`     // PXE 64-bit ARM boot file
	RootPath          string   `json:"rootpath"`          // PXE root path
	UEFIHTTPBoot      string   `json:"uefihttpboot"`      // UEFI HTTP boot URL
	DDNSDomain        string   `json:"ddnsdomain"`        // DDNS domain
	DDNSDomainPrimary string   `json:"ddnsdomainprimary"` // DDNS primary server
	DDNSDomainKeyName string   `json:"ddnsdomainkeyname"` // DDNS key name
	DDNSDomainKey     string   `json:"ddnsdomainkey"`     // DDNS key
	DDNSClientUpdates string   `json:"ddnsclientupdates"` // DDNS client updates policy
}

// ======================================================================
// Domain struct
// ======================================================================

// DHCPv4Server represents the DHCP server configuration for a single interface.
type DHCPv4Server struct {
	Interface        string
	Enable           bool
	RangeFrom        netip.Addr
	RangeTo          netip.Addr
	DNSServers       []netip.Addr
	Gateway          netip.Addr
	DomainName       string
	DomainSearchList []string
	DefaultLeaseTime time.Duration
	MaximumLeaseTime time.Duration
	WINSServers      []netip.Addr
	NTPServers       []netip.Addr
	TFTPServer       string
	LDAPServer       string
	MACAllow         string
	MACDeny          string
	DenyUnknown      bool
	IgnoreClientUIDs bool
	StaticARP        bool
}

// ======================================================================
// Stringify helpers
// ======================================================================

func (s DHCPv4Server) StringifyRangeFrom() string {
	return safeAddrString(s.RangeFrom)
}

func (s DHCPv4Server) StringifyRangeTo() string {
	return safeAddrString(s.RangeTo)
}

func (s DHCPv4Server) StringifyGateway() string {
	return safeAddrString(s.Gateway)
}

func (s DHCPv4Server) StringifyDNSServers() []string {
	result := make([]string, 0, len(s.DNSServers))
	for _, addr := range s.DNSServers {
		result = append(result, safeAddrString(addr))
	}

	return result
}

func (s DHCPv4Server) StringifyWINSServers() []string {
	result := make([]string, 0, len(s.WINSServers))
	for _, addr := range s.WINSServers {
		result = append(result, safeAddrString(addr))
	}

	return result
}

func (s DHCPv4Server) StringifyNTPServers() []string {
	result := make([]string, 0, len(s.NTPServers))
	for _, addr := range s.NTPServers {
		result = append(result, safeAddrString(addr))
	}

	return result
}

func (s DHCPv4Server) formatDomainSearchList() string {
	return strings.Join(s.DomainSearchList, dhcpv4ServerDomainSearchListSep)
}

func (s DHCPv4Server) formatDefaultLeaseTime() string {
	if s.DefaultLeaseTime == 0 {
		return ""
	}

	return fmt.Sprintf("%.0f", s.DefaultLeaseTime.Seconds())
}

func (s DHCPv4Server) formatMaximumLeaseTime() string {
	if s.MaximumLeaseTime == 0 {
		return ""
	}

	return fmt.Sprintf("%.0f", s.MaximumLeaseTime.Seconds())
}

// ======================================================================
// Setters
// ======================================================================

func (s *DHCPv4Server) SetInterface(iface string) error {
	s.Interface = iface

	return nil
}

func (s *DHCPv4Server) SetEnable(enable bool) error {
	s.Enable = enable

	return nil
}

func (s *DHCPv4Server) SetRangeFrom(addr string) error {
	if addr == "" {
		return nil
	}

	parsed, err := netip.ParseAddr(addr)
	if err != nil {
		return err
	}

	s.RangeFrom = parsed

	return nil
}

func (s *DHCPv4Server) SetRangeTo(addr string) error {
	if addr == "" {
		return nil
	}

	parsed, err := netip.ParseAddr(addr)
	if err != nil {
		return err
	}

	s.RangeTo = parsed

	return nil
}

func (s *DHCPv4Server) SetDNSServers(servers []string) error {
	for _, srv := range servers {
		addr, err := netip.ParseAddr(srv)
		if err != nil {
			return err
		}

		s.DNSServers = append(s.DNSServers, addr)
	}

	return nil
}

func (s *DHCPv4Server) SetGateway(gateway string) error {
	if gateway == "" {
		return nil
	}

	addr, err := netip.ParseAddr(gateway)
	if err != nil {
		return err
	}

	s.Gateway = addr

	return nil
}

func (s *DHCPv4Server) SetDomainName(domain string) error {
	s.DomainName = domain

	return nil
}

func (s *DHCPv4Server) SetDomainSearchList(list []string) error {
	s.DomainSearchList = list

	return nil
}

func (s *DHCPv4Server) SetDefaultLeaseTime(leaseTime string) error {
	duration, err := time.ParseDuration(leaseTime)
	if err != nil {
		return err
	}

	s.DefaultLeaseTime = duration

	return nil
}

func (s *DHCPv4Server) SetMaximumLeaseTime(leaseTime string) error {
	duration, err := time.ParseDuration(leaseTime)
	if err != nil {
		return err
	}

	s.MaximumLeaseTime = duration

	return nil
}

func (s *DHCPv4Server) SetWINSServers(servers []string) error {
	for _, srv := range servers {
		addr, err := netip.ParseAddr(srv)
		if err != nil {
			return err
		}

		s.WINSServers = append(s.WINSServers, addr)
	}

	return nil
}

func (s *DHCPv4Server) SetNTPServers(servers []string) error {
	for _, srv := range servers {
		addr, err := netip.ParseAddr(srv)
		if err != nil {
			return err
		}

		s.NTPServers = append(s.NTPServers, addr)
	}

	return nil
}

func (s *DHCPv4Server) SetTFTPServer(tftp string) error {
	s.TFTPServer = tftp

	return nil
}

func (s *DHCPv4Server) SetLDAPServer(ldap string) error {
	s.LDAPServer = ldap

	return nil
}

func (s *DHCPv4Server) SetMACAllow(mac string) error {
	s.MACAllow = mac

	return nil
}

func (s *DHCPv4Server) SetMACDeny(mac string) error {
	s.MACDeny = mac

	return nil
}

func (s *DHCPv4Server) SetDenyUnknown(deny bool) error {
	s.DenyUnknown = deny

	return nil
}

func (s *DHCPv4Server) SetIgnoreClientUIDs(ignore bool) error {
	s.IgnoreClientUIDs = ignore

	return nil
}

func (s *DHCPv4Server) SetStaticARP(staticARP bool) error {
	s.StaticARP = staticARP

	return nil
}

// ======================================================================
// Response parsing
// ======================================================================

func parseDHCPv4ServerResponse(iface string, resp dhcpv4ServerResponse) (DHCPv4Server, error) {
	var srv DHCPv4Server

	unableToParseErr := fmt.Errorf("%w DHCPv4 server response", ErrUnableToParse)

	if err := srv.SetInterface(iface); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	// enable is presence-based: the key exists in config = enabled
	if err := srv.SetEnable(resp.Enable != nil); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetRangeFrom(resp.RangeFrom); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetRangeTo(resp.RangeTo); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetDNSServers(removeEmptyStrings(resp.DNSServers)); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetGateway(resp.Gateway); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetDomainName(resp.DomainName); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetDomainSearchList(safeSplit(resp.DomainSearchList, dhcpv4ServerDomainSearchListSep)); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetDefaultLeaseTime(durationSeconds(resp.DefaultLeaseTime)); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetMaximumLeaseTime(durationSeconds(resp.MaximumLeaseTime)); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetWINSServers(removeEmptyStrings(resp.WINSServers)); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetNTPServers(removeEmptyStrings(resp.NTPServers)); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetTFTPServer(resp.TFTPServer); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetLDAPServer(resp.LDAPServer); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetMACAllow(resp.MACAllow); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetMACDeny(resp.MACDeny); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetDenyUnknown(resp.DenyUnknown != nil); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetIgnoreClientUIDs(resp.IgnoreClientUIDs != nil); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	if err := srv.SetStaticARP(resp.StaticARP != nil); err != nil {
		return srv, fmt.Errorf("%w, %w", unableToParseErr, err)
	}

	return srv, nil
}

// ======================================================================
// Client methods
// ======================================================================

func (pf *Client) getDHCPv4Server(ctx context.Context, iface string) (*DHCPv4Server, error) {
	command := fmt.Sprintf(
		"$dhcpd = config_get_path('dhcpd/%s', array());"+
			"$out = array("+
			"'enable' => isset($dhcpd['enable']) ? $dhcpd['enable'] : null,"+
			"'range_from' => isset($dhcpd['range']['from']) ? $dhcpd['range']['from'] : '',"+
			"'range_to' => isset($dhcpd['range']['to']) ? $dhcpd['range']['to'] : '',"+
			"'dnsserver' => isset($dhcpd['dnsserver']) ? $dhcpd['dnsserver'] : array(),"+
			"'gateway' => isset($dhcpd['gateway']) ? $dhcpd['gateway'] : '',"+
			"'domain' => isset($dhcpd['domain']) ? $dhcpd['domain'] : '',"+
			"'domainsearchlist' => isset($dhcpd['domainsearchlist']) ? $dhcpd['domainsearchlist'] : '',"+
			"'defaultleasetime' => isset($dhcpd['defaultleasetime']) ? $dhcpd['defaultleasetime'] : '',"+
			"'maxleasetime' => isset($dhcpd['maxleasetime']) ? $dhcpd['maxleasetime'] : '',"+
			"'winsserver' => isset($dhcpd['winsserver']) ? $dhcpd['winsserver'] : array(),"+
			"'ntpserver' => isset($dhcpd['ntpserver']) ? $dhcpd['ntpserver'] : array(),"+
			"'tftp' => isset($dhcpd['tftp']) ? $dhcpd['tftp'] : '',"+
			"'ldap' => isset($dhcpd['ldap']) ? $dhcpd['ldap'] : '',"+
			"'mac_allow' => isset($dhcpd['mac_allow']) ? $dhcpd['mac_allow'] : '',"+
			"'mac_deny' => isset($dhcpd['mac_deny']) ? $dhcpd['mac_deny'] : '',"+
			"'denyunknown' => isset($dhcpd['denyunknown']) ? $dhcpd['denyunknown'] : null,"+
			"'ignoreclientuids' => isset($dhcpd['ignoreclientuids']) ? $dhcpd['ignoreclientuids'] : null,"+
			"'staticarp' => isset($dhcpd['staticarp']) ? $dhcpd['staticarp'] : null,"+
			"'netmask' => isset($dhcpd['netmask']) ? $dhcpd['netmask'] : '',"+
			"'failover_peerip' => isset($dhcpd['failover_peerip']) ? $dhcpd['failover_peerip'] : '',"+
			"'netboot' => isset($dhcpd['netboot']) ? $dhcpd['netboot'] : null,"+
			"'nextserver' => isset($dhcpd['nextserver']) ? $dhcpd['nextserver'] : '',"+
			"'filename' => isset($dhcpd['filename']) ? $dhcpd['filename'] : '',"+
			"'filename32' => isset($dhcpd['filename32']) ? $dhcpd['filename32'] : '',"+
			"'filename32arm' => isset($dhcpd['filename32arm']) ? $dhcpd['filename32arm'] : '',"+
			"'filename64' => isset($dhcpd['filename64']) ? $dhcpd['filename64'] : '',"+
			"'filename64arm' => isset($dhcpd['filename64arm']) ? $dhcpd['filename64arm'] : '',"+
			"'rootpath' => isset($dhcpd['rootpath']) ? $dhcpd['rootpath'] : '',"+
			"'uefihttpboot' => isset($dhcpd['uefihttpboot']) ? $dhcpd['uefihttpboot'] : '',"+
			"'ddnsdomain' => isset($dhcpd['ddnsdomain']) ? $dhcpd['ddnsdomain'] : '',"+
			"'ddnsdomainprimary' => isset($dhcpd['ddnsdomainprimary']) ? $dhcpd['ddnsdomainprimary'] : '',"+
			"'ddnsdomainkeyname' => isset($dhcpd['ddnsdomainkeyname']) ? $dhcpd['ddnsdomainkeyname'] : '',"+
			"'ddnsdomainkey' => isset($dhcpd['ddnsdomainkey']) ? $dhcpd['ddnsdomainkey'] : '',"+
			"'ddnsclientupdates' => isset($dhcpd['ddnsclientupdates']) ? $dhcpd['ddnsclientupdates'] : ''"+
			");"+
			"print(json_encode($out));",
		phpEscape(iface),
	)

	var resp dhcpv4ServerResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	srv, err := parseDHCPv4ServerResponse(iface, resp)
	if err != nil {
		return nil, err
	}

	return &srv, nil
}

func (pf *Client) GetDHCPv4Server(ctx context.Context, iface string) (*DHCPv4Server, error) {
	defer pf.read(&pf.mutexes.DHCPv4Server)()

	srv, err := pf.getDHCPv4Server(ctx, iface)
	if err != nil {
		return nil, fmt.Errorf("%w '%s' DHCPv4 server, %w", ErrGetOperationFailed, iface, err)
	}

	return srv, nil
}

func dhcpv4ServerFormValues(srv DHCPv4Server) url.Values {
	values := url.Values{
		"if":               {srv.Interface},
		"range_from":       {srv.StringifyRangeFrom()},
		"range_to":         {srv.StringifyRangeTo()},
		"gateway":          {srv.StringifyGateway()},
		"domain":           {srv.DomainName},
		"domainsearchlist": {srv.formatDomainSearchList()},
		"deftime":          {srv.formatDefaultLeaseTime()},
		"maxtime":          {srv.formatMaximumLeaseTime()},
		"tftp":             {srv.TFTPServer},
		"ldap":             {srv.LDAPServer},
		"mac_allow":        {srv.MACAllow},
		"mac_deny":         {srv.MACDeny},
		"save":             {"Save"},
	}

	if srv.Enable {
		values.Set("enable", "yes")
	}

	if srv.DenyUnknown {
		values.Set("denyunknown", "yes")
	}

	if srv.IgnoreClientUIDs {
		values.Set("ignoreclientuids", "yes")
	}

	if srv.StaticARP {
		values.Set("staticarp", "yes")
	}

	for index, dns := range srv.DNSServers {
		values.Set(fmt.Sprintf("dns%d", index+1), safeAddrString(dns))
	}

	for index, wins := range srv.WINSServers {
		values.Set(fmt.Sprintf("wins%d", index+1), safeAddrString(wins))
	}

	for index, ntp := range srv.NTPServers {
		values.Set(fmt.Sprintf("ntp%d", index+1), safeAddrString(ntp))
	}

	return values
}

func (pf *Client) UpdateDHCPv4Server(ctx context.Context, srv DHCPv4Server) (*DHCPv4Server, error) {
	defer pf.write(&pf.mutexes.DHCPv4Server)()

	relativeURL := url.URL{Path: "services_dhcp.php"}
	query := relativeURL.Query()
	query.Set("if", srv.Interface)
	relativeURL.RawQuery = query.Encode()

	values := dhcpv4ServerFormValues(srv)

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return nil, fmt.Errorf("%w '%s' DHCPv4 server, %w", ErrUpdateOperationFailed, srv.Interface, err)
	}

	if err := scrapeHTMLValidationErrors(doc); err != nil {
		return nil, fmt.Errorf("%w '%s' DHCPv4 server, %w", ErrUpdateOperationFailed, srv.Interface, err)
	}

	result, err := pf.getDHCPv4Server(ctx, srv.Interface)
	if err != nil {
		return nil, fmt.Errorf("%w '%s' DHCPv4 server after updating, %w", ErrGetOperationFailed, srv.Interface, err)
	}

	return result, nil
}
