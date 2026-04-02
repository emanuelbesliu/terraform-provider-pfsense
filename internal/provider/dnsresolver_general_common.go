package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

// DNSResolverGeneralModel represents the Terraform model for DNS resolver general settings.
type DNSResolverGeneralModel struct {
	Enable                  types.Bool   `tfsdk:"enable"`
	Port                    types.Int64  `tfsdk:"port"`
	EnableSSL               types.Bool   `tfsdk:"enable_ssl"`
	TLSPort                 types.Int64  `tfsdk:"tls_port"`
	SSLCertRef              types.String `tfsdk:"ssl_cert_ref"`
	ActiveInterfaces        types.List   `tfsdk:"active_interfaces"`
	OutgoingInterfaces      types.List   `tfsdk:"outgoing_interfaces"`
	SystemDomainLocalZone   types.String `tfsdk:"system_domain_local_zone_type"`
	DNSSEC                  types.Bool   `tfsdk:"dnssec"`
	Forwarding              types.Bool   `tfsdk:"forwarding"`
	ForwardTLSUpstream      types.Bool   `tfsdk:"forward_tls_upstream"`
	RegisterDHCPLeases      types.Bool   `tfsdk:"register_dhcp_leases"`
	RegisterDHCPStaticMaps  types.Bool   `tfsdk:"register_dhcp_static_maps"`
	RegisterOpenVPNClients  types.Bool   `tfsdk:"register_openvpn_clients"`
	StrictOutgoingInterface types.Bool   `tfsdk:"strict_outgoing_interface"`
	Python                  types.Bool   `tfsdk:"python"`
	PythonOrder             types.String `tfsdk:"python_order"`
	PythonScript            types.String `tfsdk:"python_script"`
	CustomOptions           types.String `tfsdk:"custom_options"`
}

func (DNSResolverGeneralModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"enable": {
			Description: fmt.Sprintf("Enable the DNS resolver. Defaults to '%t'.", pfsense.DefaultDNSResolverGeneralEnable),
		},
		"port": {
			Description: fmt.Sprintf("Listen port number. Defaults to '%d'.", pfsense.DefaultDNSResolverGeneralPort),
		},
		"enable_ssl": {
			Description: fmt.Sprintf("Enable SSL/TLS service. Defaults to '%t'.", pfsense.DefaultDNSResolverGeneralEnableSSL),
		},
		"tls_port": {
			Description: fmt.Sprintf("SSL/TLS listen port number. Defaults to '%d'.", pfsense.DefaultDNSResolverGeneralTLSPort),
		},
		"ssl_cert_ref": {
			Description: "SSL/TLS certificate reference ID.",
		},
		"active_interfaces": {
			Description:         "Network interfaces to listen on. Use 'all' to listen on all interfaces.",
			MarkdownDescription: "Network interfaces to listen on. Use `all` to listen on all interfaces.",
		},
		"outgoing_interfaces": {
			Description:         "Network interfaces to use for outgoing queries. Use 'all' to use all interfaces.",
			MarkdownDescription: "Network interfaces to use for outgoing queries. Use `all` to use all interfaces.",
		},
		"system_domain_local_zone_type": {
			Description:         fmt.Sprintf("Local zone type for the system domain. Options: %s. Defaults to '%s'.", wrapElementsJoin(pfsense.DNSResolverGeneral{}.SystemDomainLocalZoneOptions(), "'"), pfsense.DefaultDNSResolverGeneralSystemDomainLocalZone),
			MarkdownDescription: fmt.Sprintf("Local zone type for the system domain. Options: %s. Defaults to `%s`.", wrapElementsJoin(pfsense.DNSResolverGeneral{}.SystemDomainLocalZoneOptions(), "`"), pfsense.DefaultDNSResolverGeneralSystemDomainLocalZone),
		},
		"dnssec": {
			Description: fmt.Sprintf("Enable DNSSEC support. Defaults to '%t'.", pfsense.DefaultDNSResolverGeneralDNSSEC),
		},
		"forwarding": {
			Description: fmt.Sprintf("Enable DNS query forwarding to upstream servers. Defaults to '%t'.", pfsense.DefaultDNSResolverGeneralForwarding),
		},
		"forward_tls_upstream": {
			Description: fmt.Sprintf("Use SSL/TLS for DNS query forwarding to upstream servers. Defaults to '%t'.", pfsense.DefaultDNSResolverGeneralForwardTLSUpstream),
		},
		"register_dhcp_leases": {
			Description: fmt.Sprintf("Register DHCP leases in the DNS resolver. Defaults to '%t'.", pfsense.DefaultDNSResolverGeneralRegisterDHCPLeases),
		},
		"register_dhcp_static_maps": {
			Description: fmt.Sprintf("Register DHCP static mappings in the DNS resolver. Defaults to '%t'.", pfsense.DefaultDNSResolverGeneralRegisterDHCPStaticMaps),
		},
		"register_openvpn_clients": {
			Description: fmt.Sprintf("Register connected OpenVPN clients in the DNS resolver. Defaults to '%t'.", pfsense.DefaultDNSResolverGeneralRegisterOpenVPNClients),
		},
		"strict_outgoing_interface": {
			Description: fmt.Sprintf("Strictly bind outgoing queries to the configured outgoing network interfaces. Defaults to '%t'.", pfsense.DefaultDNSResolverGeneralStrictOutgoingInterface),
		},
		"python": {
			Description: fmt.Sprintf("Enable Python module. Defaults to '%t'.", pfsense.DefaultDNSResolverGeneralPython),
		},
		"python_order": {
			Description: "Python module order.",
		},
		"python_script": {
			Description: "Python module script path.",
		},
		"custom_options": {
			Description: "Custom Unbound configuration options.",
		},
	}
}

func (DNSResolverGeneralModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enable":                        types.BoolType,
		"port":                          types.Int64Type,
		"enable_ssl":                    types.BoolType,
		"tls_port":                      types.Int64Type,
		"ssl_cert_ref":                  types.StringType,
		"active_interfaces":             types.ListType{ElemType: types.StringType},
		"outgoing_interfaces":           types.ListType{ElemType: types.StringType},
		"system_domain_local_zone_type": types.StringType,
		"dnssec":                        types.BoolType,
		"forwarding":                    types.BoolType,
		"forward_tls_upstream":          types.BoolType,
		"register_dhcp_leases":          types.BoolType,
		"register_dhcp_static_maps":     types.BoolType,
		"register_openvpn_clients":      types.BoolType,
		"strict_outgoing_interface":     types.BoolType,
		"python":                        types.BoolType,
		"python_order":                  types.StringType,
		"python_script":                 types.StringType,
		"custom_options":                types.StringType,
	}
}

func (m *DNSResolverGeneralModel) Set(ctx context.Context, dg pfsense.DNSResolverGeneral) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Enable = types.BoolValue(dg.Enable)
	m.Port = types.Int64Value(int64(dg.Port))
	m.EnableSSL = types.BoolValue(dg.EnableSSL)
	m.TLSPort = types.Int64Value(int64(dg.TLSPort))

	if dg.SSLCertRef != "" {
		m.SSLCertRef = types.StringValue(dg.SSLCertRef)
	} else {
		m.SSLCertRef = types.StringNull()
	}

	// Active interfaces
	if len(dg.ActiveInterfaces) > 0 {
		ifaceList, d := types.ListValueFrom(ctx, types.StringType, dg.ActiveInterfaces)
		diags.Append(d...)
		m.ActiveInterfaces = ifaceList
	} else {
		m.ActiveInterfaces = types.ListNull(types.StringType)
	}

	// Outgoing interfaces
	if len(dg.OutgoingInterfaces) > 0 {
		ifaceList, d := types.ListValueFrom(ctx, types.StringType, dg.OutgoingInterfaces)
		diags.Append(d...)
		m.OutgoingInterfaces = ifaceList
	} else {
		m.OutgoingInterfaces = types.ListNull(types.StringType)
	}

	m.SystemDomainLocalZone = types.StringValue(dg.SystemDomainLocalZone)
	m.DNSSEC = types.BoolValue(dg.DNSSEC)
	m.Forwarding = types.BoolValue(dg.Forwarding)
	m.ForwardTLSUpstream = types.BoolValue(dg.ForwardTLSUpstream)
	m.RegisterDHCPLeases = types.BoolValue(dg.RegisterDHCPLeases)
	m.RegisterDHCPStaticMaps = types.BoolValue(dg.RegisterDHCPStaticMaps)
	m.RegisterOpenVPNClients = types.BoolValue(dg.RegisterOpenVPNClients)
	m.StrictOutgoingInterface = types.BoolValue(dg.StrictOutgoingInterface)
	m.Python = types.BoolValue(dg.Python)

	if dg.PythonOrder != "" {
		m.PythonOrder = types.StringValue(dg.PythonOrder)
	} else {
		m.PythonOrder = types.StringNull()
	}

	if dg.PythonScript != "" {
		m.PythonScript = types.StringValue(dg.PythonScript)
	} else {
		m.PythonScript = types.StringNull()
	}

	if dg.CustomOptions != "" {
		m.CustomOptions = types.StringValue(dg.CustomOptions)
	} else {
		m.CustomOptions = types.StringNull()
	}

	return diags
}

func (m DNSResolverGeneralModel) Value(ctx context.Context, dg *pfsense.DNSResolverGeneral) diag.Diagnostics {
	var diags diag.Diagnostics

	dg.Enable = m.Enable.ValueBool()
	dg.Port = int(m.Port.ValueInt64())
	dg.EnableSSL = m.EnableSSL.ValueBool()
	dg.TLSPort = int(m.TLSPort.ValueInt64())

	if !m.SSLCertRef.IsNull() {
		dg.SSLCertRef = m.SSLCertRef.ValueString()
	}

	// Active interfaces
	if !m.ActiveInterfaces.IsNull() && !m.ActiveInterfaces.IsUnknown() {
		var ifaces []string
		diags.Append(m.ActiveInterfaces.ElementsAs(ctx, &ifaces, false)...)
		dg.ActiveInterfaces = ifaces
	} else {
		dg.ActiveInterfaces = []string{"all"}
	}

	// Outgoing interfaces
	if !m.OutgoingInterfaces.IsNull() && !m.OutgoingInterfaces.IsUnknown() {
		var ifaces []string
		diags.Append(m.OutgoingInterfaces.ElementsAs(ctx, &ifaces, false)...)
		dg.OutgoingInterfaces = ifaces
	} else {
		dg.OutgoingInterfaces = []string{"all"}
	}

	dg.SystemDomainLocalZone = m.SystemDomainLocalZone.ValueString()
	dg.DNSSEC = m.DNSSEC.ValueBool()
	dg.Forwarding = m.Forwarding.ValueBool()
	dg.ForwardTLSUpstream = m.ForwardTLSUpstream.ValueBool()
	dg.RegisterDHCPLeases = m.RegisterDHCPLeases.ValueBool()
	dg.RegisterDHCPStaticMaps = m.RegisterDHCPStaticMaps.ValueBool()
	dg.RegisterOpenVPNClients = m.RegisterOpenVPNClients.ValueBool()
	dg.StrictOutgoingInterface = m.StrictOutgoingInterface.ValueBool()
	dg.Python = m.Python.ValueBool()

	if !m.PythonOrder.IsNull() {
		dg.PythonOrder = m.PythonOrder.ValueString()
	}

	if !m.PythonScript.IsNull() {
		dg.PythonScript = m.PythonScript.ValueString()
	}

	if !m.CustomOptions.IsNull() {
		dg.CustomOptions = m.CustomOptions.ValueString()
	}

	return diags
}
