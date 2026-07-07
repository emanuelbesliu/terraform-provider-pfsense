package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type OpenVPNServersModel struct {
	All types.List `tfsdk:"all"`
}

type OpenVPNServerModel struct {
	VPNID                types.String `tfsdk:"vpn_id"`
	Disable              types.Bool   `tfsdk:"disable"`
	Mode                 types.String `tfsdk:"mode"`
	AuthMode             types.List   `tfsdk:"auth_mode"`
	DevMode              types.String `tfsdk:"dev_mode"`
	Protocol             types.String `tfsdk:"protocol"`
	Interface            types.String `tfsdk:"interface"`
	IPAddr               types.String `tfsdk:"ip_address"`
	LocalPort            types.String `tfsdk:"local_port"`
	Description          types.String `tfsdk:"description"`
	CustomOptions        types.String `tfsdk:"custom_options"`
	TLS                  types.String `tfsdk:"tls"`
	TLSType              types.String `tfsdk:"tls_type"`
	TLSAuthKeyDir        types.String `tfsdk:"tls_auth_keydir"`
	CARef                types.String `tfsdk:"ca_ref"`
	CRLRef               types.String `tfsdk:"crl_ref"`
	OCSPCheck            types.Bool   `tfsdk:"ocsp_check"`
	OCSPURL              types.String `tfsdk:"ocsp_url"`
	CertRef              types.String `tfsdk:"cert_ref"`
	DHLength             types.String `tfsdk:"dh_length"`
	ECDHCurve            types.String `tfsdk:"ecdh_curve"`
	CertDepth            types.String `tfsdk:"cert_depth"`
	StrictUserCN         types.Bool   `tfsdk:"strict_user_cn"`
	RemoteCertTLS        types.Bool   `tfsdk:"remote_cert_tls"`
	SharedKey            types.String `tfsdk:"shared_key"`
	DataCiphers          types.List   `tfsdk:"data_ciphers"`
	DataCiphersFallback  types.String `tfsdk:"data_ciphers_fallback"`
	Digest               types.String `tfsdk:"digest"`
	TunnelNetwork        types.String `tfsdk:"tunnel_network"`
	TunnelNetworkV6      types.String `tfsdk:"tunnel_network_v6"`
	LocalNetwork         types.String `tfsdk:"local_network"`
	LocalNetworkV6       types.String `tfsdk:"local_network_v6"`
	RemoteNetwork        types.String `tfsdk:"remote_network"`
	RemoteNetworkV6      types.String `tfsdk:"remote_network_v6"`
	GWRedir              types.Bool   `tfsdk:"gw_redir"`
	GWRedir6             types.Bool   `tfsdk:"gw_redir_v6"`
	Topology             types.String `tfsdk:"topology"`
	MaxClients           types.String `tfsdk:"max_clients"`
	ConnLimit            types.String `tfsdk:"connection_limit"`
	Client2Client        types.Bool   `tfsdk:"client_to_client"`
	DuplicateCN          types.Bool   `tfsdk:"duplicate_cn"`
	DynamicIP            types.Bool   `tfsdk:"dynamic_ip"`
	Compression          types.String `tfsdk:"compression"`
	CompressionPush      types.Bool   `tfsdk:"compression_push"`
	AllowCompression     types.String `tfsdk:"allow_compression"`
	PassTOS              types.Bool   `tfsdk:"pass_tos"`
	DNSDomainEnable      types.Bool   `tfsdk:"dns_domain_enable"`
	DNSDomain            types.String `tfsdk:"dns_domain"`
	DNSServerEnable      types.Bool   `tfsdk:"dns_server_enable"`
	DNSServer1           types.String `tfsdk:"dns_server1"`
	DNSServer2           types.String `tfsdk:"dns_server2"`
	DNSServer3           types.String `tfsdk:"dns_server3"`
	DNSServer4           types.String `tfsdk:"dns_server4"`
	NTPServerEnable      types.Bool   `tfsdk:"ntp_server_enable"`
	NTPServer1           types.String `tfsdk:"ntp_server1"`
	NTPServer2           types.String `tfsdk:"ntp_server2"`
	PushRegisterDNS      types.Bool   `tfsdk:"push_register_dns"`
	PushBlockOutsideDNS  types.Bool   `tfsdk:"push_block_outside_dns"`
	UsernameAsCommonName types.Bool   `tfsdk:"username_as_common_name"`
	CreateGW             types.String `tfsdk:"create_gw"`
	VerbosityLevel       types.String `tfsdk:"verbosity_level"`
}

func (OpenVPNServerModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"vpn_id": {
			Description: "Unique numeric identifier (vpnid) pfSense assigns to this OpenVPN server instance.",
		},
		"disable": {
			Description: "Whether this OpenVPN server instance is disabled.",
		},
		"mode": {
			Description: "Server mode (p2p_tls, p2p_shared_key, server_tls, server_user, server_tls_user).",
		},
		"auth_mode": {
			Description: "Authentication backends (authentication server names) used for user authentication modes.",
		},
		"dev_mode": {
			Description: "Device mode (tun for routed, tap for bridged).",
		},
		"protocol": {
			Description: "Protocol and IP version (e.g. UDP4, TCP4, UDP6, TCP6).",
		},
		"interface": {
			Description: "Interface (or 'any', 'localhost') on which the OpenVPN server listens.",
		},
		"ip_address": {
			Description: "Specific virtual IP address on the interface to bind to (optional).",
		},
		"local_port": {
			Description: "Local port on which the OpenVPN server listens.",
		},
		"description": {
			Description: descriptionDescription,
		},
		"custom_options": {
			Description: "Additional raw OpenVPN configuration options.",
		},
		"tls": {
			Description: "TLS key (tls-auth/tls-crypt) used to secure the control channel.",
		},
		"tls_type": {
			Description: "TLS key usage mode (auth for tls-auth, crypt for tls-crypt).",
		},
		"tls_auth_keydir": {
			Description: "TLS key direction (0, 1, or omitted for bidirectional).",
		},
		"ca_ref": {
			Description: "Reference (refid) of the Certificate Authority used to verify clients.",
		},
		"crl_ref": {
			Description: "Reference (refid) of the Certificate Revocation List.",
		},
		"ocsp_check": {
			Description: "Whether to check client certificates against an OCSP responder.",
		},
		"ocsp_url": {
			Description: "OCSP responder URL used for client certificate checks.",
		},
		"cert_ref": {
			Description: "Reference (refid) of the server certificate.",
		},
		"dh_length": {
			Description: "Diffie-Hellman parameter length in bits, or 'none' for ECDH only.",
		},
		"ecdh_curve": {
			Description: "ECDH curve name to use for key exchange.",
		},
		"cert_depth": {
			Description: "Maximum certificate verification depth for client certificates.",
		},
		"strict_user_cn": {
			Description: "Whether to enforce that the certificate CN matches the authenticated username.",
		},
		"remote_cert_tls": {
			Description: "Whether to require clients to present a certificate with server extended key usage.",
		},
		"shared_key": {
			Description: "Pre-shared key used for shared-key (p2p_shared_key) mode.",
		},
		"data_ciphers": {
			Description: "List of allowed data channel encryption algorithms (negotiable ciphers).",
		},
		"data_ciphers_fallback": {
			Description: "Fallback data channel cipher for clients that do not support cipher negotiation.",
		},
		"digest": {
			Description: "Authentication digest algorithm for the data channel.",
		},
		"tunnel_network": {
			Description: "IPv4 tunnel network (CIDR) used to assign virtual addresses to clients.",
		},
		"tunnel_network_v6": {
			Description: "IPv6 tunnel network (CIDR) used to assign virtual addresses to clients.",
		},
		"local_network": {
			Description: "IPv4 local networks (CIDR, comma-separated) pushed to clients as accessible routes.",
		},
		"local_network_v6": {
			Description: "IPv6 local networks (CIDR, comma-separated) pushed to clients as accessible routes.",
		},
		"remote_network": {
			Description: "IPv4 remote networks (CIDR, comma-separated) reachable behind clients.",
		},
		"remote_network_v6": {
			Description: "IPv6 remote networks (CIDR, comma-separated) reachable behind clients.",
		},
		"gw_redir": {
			Description: "Whether to force all client IPv4 traffic through the tunnel (redirect gateway).",
		},
		"gw_redir_v6": {
			Description: "Whether to force all client IPv6 traffic through the tunnel (redirect gateway).",
		},
		"topology": {
			Description: "Tunnel network topology (subnet or net30).",
		},
		"max_clients": {
			Description: "Maximum number of concurrently connected clients.",
		},
		"connection_limit": {
			Description: "Maximum number of connections from a single client at once.",
		},
		"client_to_client": {
			Description: "Whether connected clients are allowed to communicate with each other.",
		},
		"duplicate_cn": {
			Description: "Whether multiple concurrent connections using the same common name are allowed.",
		},
		"dynamic_ip": {
			Description: "Whether to allow connected clients to retain their connection when their IP changes.",
		},
		"compression": {
			Description: "Compression setting for the tunnel.",
		},
		"compression_push": {
			Description: "Whether to push the compression setting to clients.",
		},
		"allow_compression": {
			Description: "Compression policy (asym, yes, no).",
		},
		"pass_tos": {
			Description: "Whether to set the TOS IP header value of tunnel packets to match the encapsulated packet.",
		},
		"dns_domain_enable": {
			Description: "Whether to push a DNS default domain to clients.",
		},
		"dns_domain": {
			Description: "DNS default domain pushed to clients.",
		},
		"dns_server_enable": {
			Description: "Whether to push DNS servers to clients.",
		},
		"dns_server1": {
			Description: "First DNS server pushed to clients.",
		},
		"dns_server2": {
			Description: "Second DNS server pushed to clients.",
		},
		"dns_server3": {
			Description: "Third DNS server pushed to clients.",
		},
		"dns_server4": {
			Description: "Fourth DNS server pushed to clients.",
		},
		"ntp_server_enable": {
			Description: "Whether to push NTP servers to clients.",
		},
		"ntp_server1": {
			Description: "First NTP server pushed to clients.",
		},
		"ntp_server2": {
			Description: "Second NTP server pushed to clients.",
		},
		"push_register_dns": {
			Description: "Whether to push the register-dns option to Windows clients.",
		},
		"push_block_outside_dns": {
			Description: "Whether to push the block-outside-dns option to Windows clients.",
		},
		"username_as_common_name": {
			Description: "Whether to use the authenticated username instead of the certificate common name.",
		},
		"create_gw": {
			Description: "Which gateway types to create for this instance (both, v4only, v6only).",
		},
		"verbosity_level": {
			Description: "OpenVPN log verbosity level.",
		},
	}
}

func (OpenVPNServerModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"vpn_id":                  types.StringType,
		"disable":                 types.BoolType,
		"mode":                    types.StringType,
		"auth_mode":               types.ListType{ElemType: types.StringType},
		"dev_mode":                types.StringType,
		"protocol":                types.StringType,
		"interface":               types.StringType,
		"ip_address":              types.StringType,
		"local_port":              types.StringType,
		"description":             types.StringType,
		"custom_options":          types.StringType,
		"tls":                     types.StringType,
		"tls_type":                types.StringType,
		"tls_auth_keydir":         types.StringType,
		"ca_ref":                  types.StringType,
		"crl_ref":                 types.StringType,
		"ocsp_check":              types.BoolType,
		"ocsp_url":                types.StringType,
		"cert_ref":                types.StringType,
		"dh_length":               types.StringType,
		"ecdh_curve":              types.StringType,
		"cert_depth":              types.StringType,
		"strict_user_cn":          types.BoolType,
		"remote_cert_tls":         types.BoolType,
		"shared_key":              types.StringType,
		"data_ciphers":            types.ListType{ElemType: types.StringType},
		"data_ciphers_fallback":   types.StringType,
		"digest":                  types.StringType,
		"tunnel_network":          types.StringType,
		"tunnel_network_v6":       types.StringType,
		"local_network":           types.StringType,
		"local_network_v6":        types.StringType,
		"remote_network":          types.StringType,
		"remote_network_v6":       types.StringType,
		"gw_redir":                types.BoolType,
		"gw_redir_v6":             types.BoolType,
		"topology":                types.StringType,
		"max_clients":             types.StringType,
		"connection_limit":        types.StringType,
		"client_to_client":        types.BoolType,
		"duplicate_cn":            types.BoolType,
		"dynamic_ip":              types.BoolType,
		"compression":             types.StringType,
		"compression_push":        types.BoolType,
		"allow_compression":       types.StringType,
		"pass_tos":                types.BoolType,
		"dns_domain_enable":       types.BoolType,
		"dns_domain":              types.StringType,
		"dns_server_enable":       types.BoolType,
		"dns_server1":             types.StringType,
		"dns_server2":             types.StringType,
		"dns_server3":             types.StringType,
		"dns_server4":             types.StringType,
		"ntp_server_enable":       types.BoolType,
		"ntp_server1":             types.StringType,
		"ntp_server2":             types.StringType,
		"push_register_dns":       types.BoolType,
		"push_block_outside_dns":  types.BoolType,
		"username_as_common_name": types.BoolType,
		"create_gw":               types.StringType,
		"verbosity_level":         types.StringType,
	}
}

func (m *OpenVPNServersModel) Set(ctx context.Context, servers pfsense.OpenVPNServers) diag.Diagnostics {
	var diags diag.Diagnostics

	serverModels := []OpenVPNServerModel{}
	for _, server := range servers {
		var serverModel OpenVPNServerModel
		diags.Append(serverModel.Set(ctx, server)...)
		serverModels = append(serverModels, serverModel)
	}

	serversValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: OpenVPNServerModel{}.AttrTypes()}, serverModels)
	diags.Append(newDiags...)
	m.All = serversValue

	return diags
}

func (m *OpenVPNServerModel) Set(ctx context.Context, server pfsense.OpenVPNServer) diag.Diagnostics {
	var diags diag.Diagnostics

	m.VPNID = types.StringValue(server.VPNID)
	m.Disable = types.BoolValue(server.Disable)
	m.Mode = types.StringValue(server.Mode)
	m.DevMode = types.StringValue(server.DevMode)
	m.Protocol = types.StringValue(server.Protocol)
	m.Interface = types.StringValue(server.Interface)
	m.OCSPCheck = types.BoolValue(server.OCSPCheck)
	m.StrictUserCN = types.BoolValue(server.StrictUserCN)
	m.RemoteCertTLS = types.BoolValue(server.RemoteCertTLS)
	m.GWRedir = types.BoolValue(server.GWRedir)
	m.GWRedir6 = types.BoolValue(server.GWRedir6)
	m.Client2Client = types.BoolValue(server.Client2Client)
	m.DuplicateCN = types.BoolValue(server.DuplicateCN)
	m.DynamicIP = types.BoolValue(server.DynamicIP)
	m.CompressionPush = types.BoolValue(server.CompressionPush)
	m.PassTOS = types.BoolValue(server.PassTOS)
	m.DNSDomainEnable = types.BoolValue(server.DNSDomainEnable)
	m.DNSServerEnable = types.BoolValue(server.DNSServerEnable)
	m.NTPServerEnable = types.BoolValue(server.NTPServerEnable)
	m.PushRegisterDNS = types.BoolValue(server.PushRegisterDNS)
	m.PushBlockOutsideDNS = types.BoolValue(server.PushBlockOutsideDNS)
	m.UsernameAsCommonName = types.BoolValue(server.UsernameAsCommonName == "enabled")

	setOptionalString(&m.IPAddr, server.IPAddr)
	setOptionalString(&m.LocalPort, server.LocalPort)
	setOptionalString(&m.Description, server.Description)
	setOptionalString(&m.CustomOptions, server.CustomOptions)
	setOptionalString(&m.TLS, server.TLS)
	setOptionalString(&m.TLSType, server.TLSType)
	setOptionalString(&m.TLSAuthKeyDir, server.TLSAuthKeyDir)
	setOptionalString(&m.CARef, server.CARef)
	setOptionalString(&m.CRLRef, server.CRLRef)
	setOptionalString(&m.OCSPURL, server.OCSPURL)
	setOptionalString(&m.CertRef, server.CertRef)
	setOptionalString(&m.DHLength, server.DHLength)
	setOptionalString(&m.ECDHCurve, server.ECDHCurve)
	setOptionalString(&m.CertDepth, server.CertDepth)
	setOptionalString(&m.SharedKey, server.SharedKey)
	setOptionalString(&m.DataCiphersFallback, server.DataCiphersFallback)
	setOptionalString(&m.Digest, server.Digest)
	setOptionalString(&m.TunnelNetwork, server.TunnelNetwork)
	setOptionalString(&m.TunnelNetworkV6, server.TunnelNetworkV6)
	setOptionalString(&m.LocalNetwork, server.LocalNetwork)
	setOptionalString(&m.LocalNetworkV6, server.LocalNetworkV6)
	setOptionalString(&m.RemoteNetwork, server.RemoteNetwork)
	setOptionalString(&m.RemoteNetworkV6, server.RemoteNetworkV6)
	setOptionalString(&m.Topology, server.Topology)
	setOptionalString(&m.MaxClients, server.MaxClients)
	setOptionalString(&m.ConnLimit, server.ConnLimit)
	setOptionalString(&m.Compression, server.Compression)
	setOptionalString(&m.AllowCompression, server.AllowCompression)
	setOptionalString(&m.DNSDomain, server.DNSDomain)
	setOptionalString(&m.DNSServer1, server.DNSServer1)
	setOptionalString(&m.DNSServer2, server.DNSServer2)
	setOptionalString(&m.DNSServer3, server.DNSServer3)
	setOptionalString(&m.DNSServer4, server.DNSServer4)
	setOptionalString(&m.NTPServer1, server.NTPServer1)
	setOptionalString(&m.NTPServer2, server.NTPServer2)
	setOptionalString(&m.CreateGW, server.CreateGW)
	setOptionalString(&m.VerbosityLevel, server.VerbosityLevel)

	authModeValue, newDiags := types.ListValueFrom(ctx, types.StringType, server.AuthMode)
	diags.Append(newDiags...)
	m.AuthMode = authModeValue

	dataCiphersValue, newDiags := types.ListValueFrom(ctx, types.StringType, server.DataCiphers)
	diags.Append(newDiags...)
	m.DataCiphers = dataCiphersValue

	return diags
}

func (m OpenVPNServerModel) Value(ctx context.Context, server *pfsense.OpenVPNServer) diag.Diagnostics {
	var diags diag.Diagnostics

	server.Disable = m.Disable.ValueBool()
	server.Mode = m.Mode.ValueString()
	server.DevMode = m.DevMode.ValueString()
	server.Protocol = m.Protocol.ValueString()
	server.Interface = m.Interface.ValueString()
	server.OCSPCheck = m.OCSPCheck.ValueBool()
	server.StrictUserCN = m.StrictUserCN.ValueBool()
	server.RemoteCertTLS = m.RemoteCertTLS.ValueBool()
	server.GWRedir = m.GWRedir.ValueBool()
	server.GWRedir6 = m.GWRedir6.ValueBool()
	server.Client2Client = m.Client2Client.ValueBool()
	server.DuplicateCN = m.DuplicateCN.ValueBool()
	server.DynamicIP = m.DynamicIP.ValueBool()
	server.CompressionPush = m.CompressionPush.ValueBool()
	server.PassTOS = m.PassTOS.ValueBool()
	server.DNSDomainEnable = m.DNSDomainEnable.ValueBool()
	server.DNSServerEnable = m.DNSServerEnable.ValueBool()
	server.NTPServerEnable = m.NTPServerEnable.ValueBool()
	server.PushRegisterDNS = m.PushRegisterDNS.ValueBool()
	server.PushBlockOutsideDNS = m.PushBlockOutsideDNS.ValueBool()

	if m.UsernameAsCommonName.ValueBool() {
		server.UsernameAsCommonName = "enabled"
	} else {
		server.UsernameAsCommonName = "disabled"
	}

	if !m.VPNID.IsNull() && !m.VPNID.IsUnknown() {
		server.VPNID = m.VPNID.ValueString()
	}

	getOptionalString(m.IPAddr, &server.IPAddr)
	getOptionalString(m.LocalPort, &server.LocalPort)
	getOptionalString(m.Description, &server.Description)
	getOptionalString(m.CustomOptions, &server.CustomOptions)
	getOptionalString(m.TLS, &server.TLS)
	getOptionalString(m.TLSType, &server.TLSType)
	getOptionalString(m.TLSAuthKeyDir, &server.TLSAuthKeyDir)
	getOptionalString(m.CARef, &server.CARef)
	getOptionalString(m.CRLRef, &server.CRLRef)
	getOptionalString(m.OCSPURL, &server.OCSPURL)
	getOptionalString(m.CertRef, &server.CertRef)
	getOptionalString(m.DHLength, &server.DHLength)
	getOptionalString(m.ECDHCurve, &server.ECDHCurve)
	getOptionalString(m.CertDepth, &server.CertDepth)
	getOptionalString(m.SharedKey, &server.SharedKey)
	getOptionalString(m.DataCiphersFallback, &server.DataCiphersFallback)
	getOptionalString(m.Digest, &server.Digest)
	getOptionalString(m.TunnelNetwork, &server.TunnelNetwork)
	getOptionalString(m.TunnelNetworkV6, &server.TunnelNetworkV6)
	getOptionalString(m.LocalNetwork, &server.LocalNetwork)
	getOptionalString(m.LocalNetworkV6, &server.LocalNetworkV6)
	getOptionalString(m.RemoteNetwork, &server.RemoteNetwork)
	getOptionalString(m.RemoteNetworkV6, &server.RemoteNetworkV6)
	getOptionalString(m.Topology, &server.Topology)
	getOptionalString(m.MaxClients, &server.MaxClients)
	getOptionalString(m.ConnLimit, &server.ConnLimit)
	getOptionalString(m.Compression, &server.Compression)
	getOptionalString(m.AllowCompression, &server.AllowCompression)
	getOptionalString(m.DNSDomain, &server.DNSDomain)
	getOptionalString(m.DNSServer1, &server.DNSServer1)
	getOptionalString(m.DNSServer2, &server.DNSServer2)
	getOptionalString(m.DNSServer3, &server.DNSServer3)
	getOptionalString(m.DNSServer4, &server.DNSServer4)
	getOptionalString(m.NTPServer1, &server.NTPServer1)
	getOptionalString(m.NTPServer2, &server.NTPServer2)
	getOptionalString(m.CreateGW, &server.CreateGW)
	getOptionalString(m.VerbosityLevel, &server.VerbosityLevel)

	if !m.AuthMode.IsNull() && !m.AuthMode.IsUnknown() {
		diags.Append(m.AuthMode.ElementsAs(ctx, &server.AuthMode, false)...)
	}

	if !m.DataCiphers.IsNull() && !m.DataCiphers.IsUnknown() {
		diags.Append(m.DataCiphers.ElementsAs(ctx, &server.DataCiphers, false)...)
	}

	return diags
}

func setOptionalString(field *types.String, value string) {
	if value != "" {
		*field = types.StringValue(value)
	} else {
		*field = types.StringNull()
	}
}

func getOptionalString(field types.String, value *string) {
	if !field.IsNull() && !field.IsUnknown() {
		*value = field.ValueString()
	}
}
