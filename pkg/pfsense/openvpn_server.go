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

// ============================================================================
// Flexible JSON helpers
// ============================================================================

// flexBool tolerates pfSense storing presence-based booleans as either a JSON
// boolean (true/false) or a string ("yes", "", "1", "0").
type flexBool bool

func (b *flexBool) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	switch strings.ToLower(s) {
	case "", "false", "0", "null", "no", "disabled":
		*b = false
	default:
		*b = true
	}

	return nil
}

// flexString tolerates pfSense returning either a string or a number for a
// field.
type flexString string

func (s *flexString) UnmarshalJSON(data []byte) error {
	*s = flexString(strings.Trim(string(data), `"`))

	return nil
}

// ============================================================================
// Response types (JSON from PHP config read)
// ============================================================================

type openVPNServerResponse struct {
	VPNID                flexString `json:"vpnid"`
	Disable              flexBool   `json:"disable"`
	Mode                 string     `json:"mode"`
	AuthMode             string     `json:"authmode"`
	DevMode              string     `json:"dev_mode"`
	Protocol             string     `json:"protocol"`
	Interface            string     `json:"interface"`
	IPAddr               string     `json:"ipaddr"`
	LocalPort            flexString `json:"local_port"`
	Description          string     `json:"description"`
	CustomOptions        string     `json:"custom_options"`
	TLS                  string     `json:"tls"`
	TLSType              string     `json:"tls_type"`
	TLSAuthKeyDir        string     `json:"tlsauth_keydir"`
	CARef                string     `json:"caref"`
	CRLRef               string     `json:"crlref"`
	OCSPCheck            flexBool   `json:"ocspcheck"`
	OCSPURL              string     `json:"ocspurl"`
	CertRef              string     `json:"certref"`
	DHLength             flexString `json:"dh_length"`
	ECDHCurve            string     `json:"ecdh_curve"`
	CertDepth            flexString `json:"cert_depth"`
	StrictUserCN         flexBool   `json:"strictusercn"`
	RemoteCertTLS        flexBool   `json:"remote_cert_tls"`
	SharedKey            string     `json:"shared_key"`
	DataCiphers          string     `json:"data_ciphers"`
	DataCiphersFallback  string     `json:"data_ciphers_fallback"`
	Digest               string     `json:"digest"`
	TunnelNetwork        string     `json:"tunnel_network"`
	TunnelNetworkV6      string     `json:"tunnel_networkv6"`
	LocalNetwork         string     `json:"local_network"`
	LocalNetworkV6       string     `json:"local_networkv6"`
	RemoteNetwork        string     `json:"remote_network"`
	RemoteNetworkV6      string     `json:"remote_networkv6"`
	GWRedir              flexBool   `json:"gwredir"`
	GWRedir6             flexBool   `json:"gwredir6"`
	Topology             string     `json:"topology"`
	MaxClients           flexString `json:"maxclients"`
	ConnLimit            flexString `json:"connlimit"`
	Client2Client        flexBool   `json:"client2client"`
	DuplicateCN          flexBool   `json:"duplicate_cn"`
	DynamicIP            flexBool   `json:"dynamic_ip"`
	Compression          string     `json:"compression"`
	CompressionPush      flexBool   `json:"compression_push"`
	AllowCompression     string     `json:"allow_compression"`
	PassTOS              flexBool   `json:"passtos"`
	DNSDomainEnable      flexBool   `json:"dns_domain_enable"`
	DNSDomain            string     `json:"dns_domain"`
	DNSServerEnable      flexBool   `json:"dns_server_enable"`
	DNSServer1           string     `json:"dns_server1"`
	DNSServer2           string     `json:"dns_server2"`
	DNSServer3           string     `json:"dns_server3"`
	DNSServer4           string     `json:"dns_server4"`
	NTPServerEnable      flexBool   `json:"ntp_server_enable"`
	NTPServer1           string     `json:"ntp_server1"`
	NTPServer2           string     `json:"ntp_server2"`
	PushRegisterDNS      flexBool   `json:"push_register_dns"`
	PushBlockOutsideDNS  flexBool   `json:"push_blockoutsidedns"`
	UsernameAsCommonName string     `json:"username_as_common_name"`
	CreateGW             string     `json:"create_gw"`
	VerbosityLevel       flexString `json:"verbosity_level"`
}

// ============================================================================
// Domain type
// ============================================================================

type OpenVPNServer struct {
	VPNID                string
	Disable              bool
	Mode                 string
	AuthMode             []string
	DevMode              string
	Protocol             string
	Interface            string
	IPAddr               string
	LocalPort            string
	Description          string
	CustomOptions        string
	TLS                  string
	TLSType              string
	TLSAuthKeyDir        string
	CARef                string
	CRLRef               string
	OCSPCheck            bool
	OCSPURL              string
	CertRef              string
	DHLength             string
	ECDHCurve            string
	CertDepth            string
	StrictUserCN         bool
	RemoteCertTLS        bool
	SharedKey            string
	DataCiphers          []string
	DataCiphersFallback  string
	Digest               string
	TunnelNetwork        string
	TunnelNetworkV6      string
	LocalNetwork         string
	LocalNetworkV6       string
	RemoteNetwork        string
	RemoteNetworkV6      string
	GWRedir              bool
	GWRedir6             bool
	Topology             string
	MaxClients           string
	ConnLimit            string
	Client2Client        bool
	DuplicateCN          bool
	DynamicIP            bool
	Compression          string
	CompressionPush      bool
	AllowCompression     string
	PassTOS              bool
	DNSDomainEnable      bool
	DNSDomain            string
	DNSServerEnable      bool
	DNSServer1           string
	DNSServer2           string
	DNSServer3           string
	DNSServer4           string
	NTPServerEnable      bool
	NTPServer1           string
	NTPServer2           string
	PushRegisterDNS      bool
	PushBlockOutsideDNS  bool
	UsernameAsCommonName string
	CreateGW             string
	VerbosityLevel       string
	// controlID is the array index used by pfSense to edit/delete the server.
	controlID int
}

type OpenVPNServers []OpenVPNServer

func (servers OpenVPNServers) GetByVPNID(vpnID string) (*OpenVPNServer, error) {
	for _, server := range servers {
		if server.VPNID == vpnID {
			return &server, nil
		}
	}

	return nil, fmt.Errorf("openvpn server %w with vpnid '%s'", ErrNotFound, vpnID)
}

func (servers OpenVPNServers) GetControlIDByVPNID(vpnID string) (*int, error) {
	for _, server := range servers {
		if server.VPNID == vpnID {
			controlID := server.controlID

			return &controlID, nil
		}
	}

	return nil, fmt.Errorf("openvpn server %w with vpnid '%s'", ErrNotFound, vpnID)
}

// ============================================================================
// Parsing
// ============================================================================

func splitCommaList(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}

	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}

	return out
}

func parseOpenVPNServerResponse(resp openVPNServerResponse, controlID int) OpenVPNServer {
	return OpenVPNServer{
		VPNID:                string(resp.VPNID),
		Disable:              bool(resp.Disable),
		Mode:                 resp.Mode,
		AuthMode:             splitCommaList(resp.AuthMode),
		DevMode:              resp.DevMode,
		Protocol:             resp.Protocol,
		Interface:            resp.Interface,
		IPAddr:               resp.IPAddr,
		LocalPort:            string(resp.LocalPort),
		Description:          resp.Description,
		CustomOptions:        resp.CustomOptions,
		TLS:                  resp.TLS,
		TLSType:              resp.TLSType,
		TLSAuthKeyDir:        resp.TLSAuthKeyDir,
		CARef:                resp.CARef,
		CRLRef:               resp.CRLRef,
		OCSPCheck:            bool(resp.OCSPCheck),
		OCSPURL:              resp.OCSPURL,
		CertRef:              resp.CertRef,
		DHLength:             string(resp.DHLength),
		ECDHCurve:            resp.ECDHCurve,
		CertDepth:            string(resp.CertDepth),
		StrictUserCN:         bool(resp.StrictUserCN),
		RemoteCertTLS:        bool(resp.RemoteCertTLS),
		SharedKey:            resp.SharedKey,
		DataCiphers:          splitCommaList(resp.DataCiphers),
		DataCiphersFallback:  resp.DataCiphersFallback,
		Digest:               resp.Digest,
		TunnelNetwork:        resp.TunnelNetwork,
		TunnelNetworkV6:      resp.TunnelNetworkV6,
		LocalNetwork:         resp.LocalNetwork,
		LocalNetworkV6:       resp.LocalNetworkV6,
		RemoteNetwork:        resp.RemoteNetwork,
		RemoteNetworkV6:      resp.RemoteNetworkV6,
		GWRedir:              bool(resp.GWRedir),
		GWRedir6:             bool(resp.GWRedir6),
		Topology:             resp.Topology,
		MaxClients:           string(resp.MaxClients),
		ConnLimit:            string(resp.ConnLimit),
		Client2Client:        bool(resp.Client2Client),
		DuplicateCN:          bool(resp.DuplicateCN),
		DynamicIP:            bool(resp.DynamicIP),
		Compression:          resp.Compression,
		CompressionPush:      bool(resp.CompressionPush),
		AllowCompression:     resp.AllowCompression,
		PassTOS:              bool(resp.PassTOS),
		DNSDomainEnable:      bool(resp.DNSDomainEnable),
		DNSDomain:            resp.DNSDomain,
		DNSServerEnable:      bool(resp.DNSServerEnable),
		DNSServer1:           resp.DNSServer1,
		DNSServer2:           resp.DNSServer2,
		DNSServer3:           resp.DNSServer3,
		DNSServer4:           resp.DNSServer4,
		NTPServerEnable:      bool(resp.NTPServerEnable),
		NTPServer1:           resp.NTPServer1,
		NTPServer2:           resp.NTPServer2,
		PushRegisterDNS:      bool(resp.PushRegisterDNS),
		PushBlockOutsideDNS:  bool(resp.PushBlockOutsideDNS),
		UsernameAsCommonName: resp.UsernameAsCommonName,
		CreateGW:             resp.CreateGW,
		VerbosityLevel:       string(resp.VerbosityLevel),
		controlID:            controlID,
	}
}

// ============================================================================
// Form values for POST
// ============================================================================

func openVPNServerFormValues(server OpenVPNServer) url.Values {
	values := url.Values{
		"save":     {"Save"},
		"mode":     {server.Mode},
		"dev_mode": {server.DevMode},
		"protocol": {server.Protocol},
	}

	setIf := func(key, value string) {
		if value != "" {
			values.Set(key, value)
		}
	}

	setBool := func(key string, value bool) {
		if value {
			values.Set(key, "yes")
		}
	}

	// Interface can carry a virtual IP suffix ("iface|ipaddr").
	iface := server.Interface
	if server.IPAddr != "" {
		iface = fmt.Sprintf("%s|%s", server.Interface, server.IPAddr)
	}
	values.Set("interface", iface)

	setIf("local_port", server.LocalPort)
	setIf("description", server.Description)
	setIf("custom_options", server.CustomOptions)
	setBool("disable", server.Disable)

	for _, mode := range server.AuthMode {
		values.Add("authmode[]", mode)
	}

	// TLS / crypto.
	setIf("tls", server.TLS)
	setIf("tls_type", server.TLSType)
	setIf("tlsauth_keydir", server.TLSAuthKeyDir)
	if server.TLS != "" {
		values.Set("tlsauth_enable", "yes")
	}
	setIf("caref", server.CARef)
	setIf("crlref", server.CRLRef)
	setBool("ocspcheck", server.OCSPCheck)
	setIf("ocspurl", server.OCSPURL)
	setIf("certref", server.CertRef)
	setIf("dh_length", server.DHLength)
	setIf("ecdh_curve", server.ECDHCurve)
	setIf("cert_depth", server.CertDepth)
	setBool("strictusercn", server.StrictUserCN)
	setBool("remote_cert_tls", server.RemoteCertTLS)
	setIf("shared_key", server.SharedKey)

	for _, cipher := range server.DataCiphers {
		values.Add("data_ciphers[]", cipher)
	}
	setIf("data_ciphers_fallback", server.DataCiphersFallback)
	setIf("digest", server.Digest)

	// Tunnel / routing.
	setIf("tunnel_network", server.TunnelNetwork)
	setIf("tunnel_networkv6", server.TunnelNetworkV6)
	setIf("local_network", server.LocalNetwork)
	setIf("local_networkv6", server.LocalNetworkV6)
	setIf("remote_network", server.RemoteNetwork)
	setIf("remote_networkv6", server.RemoteNetworkV6)
	setBool("gwredir", server.GWRedir)
	setBool("gwredir6", server.GWRedir6)
	setIf("topology", server.Topology)
	setIf("maxclients", server.MaxClients)
	setIf("connlimit", server.ConnLimit)
	setBool("client2client", server.Client2Client)
	setBool("duplicate_cn", server.DuplicateCN)
	setBool("dynamic_ip", server.DynamicIP)

	// Compression.
	setIf("compression", server.Compression)
	setBool("compression_push", server.CompressionPush)
	setIf("allow_compression", server.AllowCompression)
	setBool("passtos", server.PassTOS)

	// Client settings (push).
	setBool("dns_domain_enable", server.DNSDomainEnable)
	setIf("dns_domain", server.DNSDomain)
	setBool("dns_server_enable", server.DNSServerEnable)
	setIf("dns_server1", server.DNSServer1)
	setIf("dns_server2", server.DNSServer2)
	setIf("dns_server3", server.DNSServer3)
	setIf("dns_server4", server.DNSServer4)
	setBool("ntp_server_enable", server.NTPServerEnable)
	setIf("ntp_server1", server.NTPServer1)
	setIf("ntp_server2", server.NTPServer2)
	setBool("push_register_dns", server.PushRegisterDNS)
	setBool("push_blockoutsidedns", server.PushBlockOutsideDNS)
	setBool("username_as_common_name", server.UsernameAsCommonName == "enabled")
	setIf("create_gw", server.CreateGW)
	setIf("verbosity_level", server.VerbosityLevel)

	return values
}

// ============================================================================
// Client methods
// ============================================================================

func (pf *Client) getOpenVPNServers(ctx context.Context) (*OpenVPNServers, error) {
	unableToParseResErr := fmt.Errorf("%w openvpn server response", ErrUnableToParse)

	command := `
$servers = config_get_path('openvpn/openvpn-server', array());
if (!is_array($servers)) { $servers = array(); }
print(json_encode(array_values($servers)));
`
	var resp []openVPNServerResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	servers := make(OpenVPNServers, 0, len(resp))
	for index, r := range resp {
		servers = append(servers, parseOpenVPNServerResponse(r, index))
	}

	// Sanity check that the parsed data round-trips.
	if _, err := json.Marshal(servers); err != nil {
		return nil, fmt.Errorf("%w, %w", unableToParseResErr, err)
	}

	return &servers, nil
}

func (pf *Client) GetOpenVPNServers(ctx context.Context) (*OpenVPNServers, error) {
	defer pf.read(&pf.mutexes.OpenVPNServer)()

	servers, err := pf.getOpenVPNServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w openvpn servers, %w", ErrGetOperationFailed, err)
	}

	return servers, nil
}

func (pf *Client) GetOpenVPNServer(ctx context.Context, vpnID string) (*OpenVPNServer, error) {
	defer pf.read(&pf.mutexes.OpenVPNServer)()

	servers, err := pf.getOpenVPNServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w openvpn servers, %w", ErrGetOperationFailed, err)
	}

	server, err := servers.GetByVPNID(vpnID)
	if err != nil {
		return nil, fmt.Errorf("%w openvpn server, %w", ErrGetOperationFailed, err)
	}

	return server, nil
}

func (pf *Client) createOrUpdateOpenVPNServer(ctx context.Context, serverReq OpenVPNServer, controlID *int) error {
	relativeURL := url.URL{Path: "vpn_openvpn_server.php"}

	q := relativeURL.Query()
	if controlID != nil {
		q.Set("act", "edit")
		q.Set("id", strconv.Itoa(*controlID))
	} else {
		q.Set("act", "new")
	}
	relativeURL.RawQuery = q.Encode()

	values := openVPNServerFormValues(serverReq)

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return err
	}

	return scrapeHTMLValidationErrors(doc)
}

func (pf *Client) CreateOpenVPNServer(ctx context.Context, serverReq OpenVPNServer) (*OpenVPNServer, error) {
	defer pf.write(&pf.mutexes.OpenVPNServer)()

	before, err := pf.getOpenVPNServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w openvpn servers before creating, %w", ErrGetOperationFailed, err)
	}

	existing := make(map[string]struct{}, len(*before))
	for _, server := range *before {
		existing[server.VPNID] = struct{}{}
	}

	if err := pf.createOrUpdateOpenVPNServer(ctx, serverReq, nil); err != nil {
		return nil, fmt.Errorf("%w openvpn server, %w", ErrCreateOperationFailed, err)
	}

	after, err := pf.getOpenVPNServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w openvpn servers after creating, %w", ErrGetOperationFailed, err)
	}

	// The newly created server is the one whose vpnid was not present before.
	var found *OpenVPNServer
	for i := len(*after) - 1; i >= 0; i-- {
		server := (*after)[i]
		if _, ok := existing[server.VPNID]; !ok {
			found = &server

			break
		}
	}

	if found == nil {
		return nil, fmt.Errorf("%w openvpn server after creating, could not find newly created entry", ErrGetOperationFailed)
	}

	return found, nil
}

func (pf *Client) UpdateOpenVPNServer(ctx context.Context, serverReq OpenVPNServer) (*OpenVPNServer, error) {
	defer pf.write(&pf.mutexes.OpenVPNServer)()

	servers, err := pf.getOpenVPNServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w openvpn servers, %w", ErrGetOperationFailed, err)
	}

	controlID, err := servers.GetControlIDByVPNID(serverReq.VPNID)
	if err != nil {
		return nil, fmt.Errorf("%w openvpn server, %w", ErrGetOperationFailed, err)
	}

	if err := pf.createOrUpdateOpenVPNServer(ctx, serverReq, controlID); err != nil {
		return nil, fmt.Errorf("%w openvpn server, %w", ErrUpdateOperationFailed, err)
	}

	servers, err = pf.getOpenVPNServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w openvpn servers after updating, %w", ErrGetOperationFailed, err)
	}

	server, err := servers.GetByVPNID(serverReq.VPNID)
	if err != nil {
		return nil, fmt.Errorf("%w openvpn server after updating, %w", ErrGetOperationFailed, err)
	}

	return server, nil
}

func (pf *Client) deleteOpenVPNServer(ctx context.Context, controlID int) error {
	relativeURL := url.URL{Path: "vpn_openvpn_server.php"}
	values := url.Values{
		"act": {"del"},
		"id":  {strconv.Itoa(controlID)},
	}

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return err
	}

	return scrapeHTMLValidationErrors(doc)
}

func (pf *Client) DeleteOpenVPNServer(ctx context.Context, vpnID string) error {
	defer pf.write(&pf.mutexes.OpenVPNServer)()

	servers, err := pf.getOpenVPNServers(ctx)
	if err != nil {
		return fmt.Errorf("%w openvpn servers, %w", ErrGetOperationFailed, err)
	}

	controlID, err := servers.GetControlIDByVPNID(vpnID)
	if err != nil {
		return fmt.Errorf("%w openvpn server, %w", ErrGetOperationFailed, err)
	}

	if err := pf.deleteOpenVPNServer(ctx, *controlID); err != nil {
		return fmt.Errorf("%w openvpn server, %w", ErrDeleteOperationFailed, err)
	}

	servers, err = pf.getOpenVPNServers(ctx)
	if err != nil {
		return fmt.Errorf("%w openvpn servers after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := servers.GetByVPNID(vpnID); err == nil {
		return fmt.Errorf("%w openvpn server, still exists", ErrDeleteOperationFailed)
	}

	return nil
}
