package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

// DHCPv4ServerModel represents the Terraform model for DHCP server settings on a single interface.
type DHCPv4ServerModel struct {
	Interface        types.String         `tfsdk:"interface"`
	Enable           types.Bool           `tfsdk:"enable"`
	RangeFrom        types.String         `tfsdk:"range_from"`
	RangeTo          types.String         `tfsdk:"range_to"`
	DNSServers       types.List           `tfsdk:"dns_servers"`
	Gateway          types.String         `tfsdk:"gateway"`
	DomainName       types.String         `tfsdk:"domain_name"`
	DomainSearchList types.List           `tfsdk:"domain_search_list"`
	DefaultLeaseTime timetypes.GoDuration `tfsdk:"default_lease_time"`
	MaximumLeaseTime timetypes.GoDuration `tfsdk:"maximum_lease_time"`
	WINSServers      types.List           `tfsdk:"wins_servers"`
	NTPServers       types.List           `tfsdk:"ntp_servers"`
	TFTPServer       types.String         `tfsdk:"tftp_server"`
	LDAPServer       types.String         `tfsdk:"ldap_server"`
	MACAllow         types.String         `tfsdk:"mac_allow"`
	MACDeny          types.String         `tfsdk:"mac_deny"`
	DenyUnknown      types.Bool           `tfsdk:"deny_unknown"`
	IgnoreClientUIDs types.Bool           `tfsdk:"ignore_client_uids"`
	StaticARP        types.Bool           `tfsdk:"static_arp"`
}

func (DHCPv4ServerModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"interface": {
			Description: "Network interface. Each interface has its own separate DHCP server configuration.",
		},
		"enable": {
			Description: "Enable DHCP server on this interface.",
		},
		"range_from": {
			Description: "Start of the DHCP address pool range.",
		},
		"range_to": {
			Description: "End of the DHCP address pool range.",
		},
		"dns_servers": {
			Description: fmt.Sprintf("DNS servers provided to DHCP clients (up to %d).", pfsense.DHCPv4ServerMaxDNSServers),
		},
		"gateway": {
			Description: "Gateway IPv4 address provided to DHCP clients. Leave empty to use the interface IP address.",
		},
		"domain_name": {
			Description: "Domain name provided to DHCP clients.",
		},
		"domain_search_list": {
			Description: "DNS search domains provided to DHCP clients.",
		},
		"default_lease_time": {
			Description: "Default lease time for clients that do not ask for a specific expiration time.",
		},
		"maximum_lease_time": {
			Description: "Maximum lease time for clients that ask for a specific expiration time.",
		},
		"wins_servers": {
			Description: fmt.Sprintf("WINS (Windows Internet Name Service) servers provided to DHCP clients (up to %d).", pfsense.DHCPv4ServerMaxWINSServers),
		},
		"ntp_servers": {
			Description: fmt.Sprintf("NTP (Network Time Protocol) servers provided to DHCP clients (up to %d).", pfsense.DHCPv4ServerMaxNTPServers),
		},
		"tftp_server": {
			Description: "TFTP server hostname or IP address provided to DHCP clients.",
		},
		"ldap_server": {
			Description: "LDAP server URI provided to DHCP clients.",
		},
		"mac_allow": {
			Description: "Comma-separated list of partial MAC addresses to allow. Patterns matched from the beginning (e.g. '00:1A:2B').",
		},
		"mac_deny": {
			Description: "Comma-separated list of partial MAC addresses to deny. Patterns matched from the beginning (e.g. '00:1A:2B').",
		},
		"deny_unknown": {
			Description: "Deny leases to unknown MAC addresses. Only known clients (those with static mappings) will receive leases.",
		},
		"ignore_client_uids": {
			Description: "Ignore client UIDs (option 61) in DHCP requests. Use only MAC addresses for identification.",
		},
		"static_arp": {
			Description:         "Enable Static ARP entries. All clients must have static mappings. Only the machines listed below will be able to communicate with the firewall on this interface.",
			MarkdownDescription: "Enable Static ARP entries. **All** clients must have static mappings. Only the machines listed will be able to communicate with the firewall on this interface.",
		},
	}
}

func (DHCPv4ServerModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interface":          types.StringType,
		"enable":             types.BoolType,
		"range_from":         types.StringType,
		"range_to":           types.StringType,
		"dns_servers":        types.ListType{ElemType: types.StringType},
		"gateway":            types.StringType,
		"domain_name":        types.StringType,
		"domain_search_list": types.ListType{ElemType: types.StringType},
		"default_lease_time": timetypes.GoDurationType{},
		"maximum_lease_time": timetypes.GoDurationType{},
		"wins_servers":       types.ListType{ElemType: types.StringType},
		"ntp_servers":        types.ListType{ElemType: types.StringType},
		"tftp_server":        types.StringType,
		"ldap_server":        types.StringType,
		"mac_allow":          types.StringType,
		"mac_deny":           types.StringType,
		"deny_unknown":       types.BoolType,
		"ignore_client_uids": types.BoolType,
		"static_arp":         types.BoolType,
	}
}

func (m *DHCPv4ServerModel) Set(ctx context.Context, srv pfsense.DHCPv4Server) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Interface = types.StringValue(srv.Interface)
	m.Enable = types.BoolValue(srv.Enable)

	if srv.StringifyRangeFrom() != "" {
		m.RangeFrom = types.StringValue(srv.StringifyRangeFrom())
	} else {
		m.RangeFrom = types.StringNull()
	}

	if srv.StringifyRangeTo() != "" {
		m.RangeTo = types.StringValue(srv.StringifyRangeTo())
	} else {
		m.RangeTo = types.StringNull()
	}

	// DNS servers
	dnsServers := srv.StringifyDNSServers()
	if len(dnsServers) > 0 {
		dnsValue, d := types.ListValueFrom(ctx, types.StringType, dnsServers)
		diags.Append(d...)
		m.DNSServers = dnsValue
	} else {
		m.DNSServers = types.ListNull(types.StringType)
	}

	if srv.StringifyGateway() != "" {
		m.Gateway = types.StringValue(srv.StringifyGateway())
	} else {
		m.Gateway = types.StringNull()
	}

	if srv.DomainName != "" {
		m.DomainName = types.StringValue(srv.DomainName)
	} else {
		m.DomainName = types.StringNull()
	}

	// Domain search list
	if len(srv.DomainSearchList) > 0 {
		dslValue, d := types.ListValueFrom(ctx, types.StringType, srv.DomainSearchList)
		diags.Append(d...)
		m.DomainSearchList = dslValue
	} else {
		m.DomainSearchList = types.ListNull(types.StringType)
	}

	// Lease times
	if srv.DefaultLeaseTime != 0 {
		m.DefaultLeaseTime = timetypes.NewGoDurationValue(srv.DefaultLeaseTime)
	} else {
		m.DefaultLeaseTime = timetypes.NewGoDurationNull()
	}

	if srv.MaximumLeaseTime != 0 {
		m.MaximumLeaseTime = timetypes.NewGoDurationValue(srv.MaximumLeaseTime)
	} else {
		m.MaximumLeaseTime = timetypes.NewGoDurationNull()
	}

	// WINS servers
	winsServers := srv.StringifyWINSServers()
	if len(winsServers) > 0 {
		winsValue, d := types.ListValueFrom(ctx, types.StringType, winsServers)
		diags.Append(d...)
		m.WINSServers = winsValue
	} else {
		m.WINSServers = types.ListNull(types.StringType)
	}

	// NTP servers
	ntpServers := srv.StringifyNTPServers()
	if len(ntpServers) > 0 {
		ntpValue, d := types.ListValueFrom(ctx, types.StringType, ntpServers)
		diags.Append(d...)
		m.NTPServers = ntpValue
	} else {
		m.NTPServers = types.ListNull(types.StringType)
	}

	if srv.TFTPServer != "" {
		m.TFTPServer = types.StringValue(srv.TFTPServer)
	} else {
		m.TFTPServer = types.StringNull()
	}

	if srv.LDAPServer != "" {
		m.LDAPServer = types.StringValue(srv.LDAPServer)
	} else {
		m.LDAPServer = types.StringNull()
	}

	if srv.MACAllow != "" {
		m.MACAllow = types.StringValue(srv.MACAllow)
	} else {
		m.MACAllow = types.StringNull()
	}

	if srv.MACDeny != "" {
		m.MACDeny = types.StringValue(srv.MACDeny)
	} else {
		m.MACDeny = types.StringNull()
	}

	m.DenyUnknown = types.BoolValue(srv.DenyUnknown)
	m.IgnoreClientUIDs = types.BoolValue(srv.IgnoreClientUIDs)
	m.StaticARP = types.BoolValue(srv.StaticARP)

	return diags
}

func (m DHCPv4ServerModel) Value(ctx context.Context, srv *pfsense.DHCPv4Server) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("interface"),
		"Interface cannot be parsed",
		srv.SetInterface(m.Interface.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("enable"),
		"Enable cannot be parsed",
		srv.SetEnable(m.Enable.ValueBool()),
	)

	if !m.RangeFrom.IsNull() {
		addPathError(
			&diags,
			path.Root("range_from"),
			"Range from cannot be parsed",
			srv.SetRangeFrom(m.RangeFrom.ValueString()),
		)
	}

	if !m.RangeTo.IsNull() {
		addPathError(
			&diags,
			path.Root("range_to"),
			"Range to cannot be parsed",
			srv.SetRangeTo(m.RangeTo.ValueString()),
		)
	}

	if !m.DNSServers.IsNull() {
		var dnsServers []string
		diags.Append(m.DNSServers.ElementsAs(ctx, &dnsServers, false)...)
		addPathError(
			&diags,
			path.Root("dns_servers"),
			"DNS servers cannot be parsed",
			srv.SetDNSServers(dnsServers),
		)
	}

	if !m.Gateway.IsNull() {
		addPathError(
			&diags,
			path.Root("gateway"),
			"Gateway cannot be parsed",
			srv.SetGateway(m.Gateway.ValueString()),
		)
	}

	if !m.DomainName.IsNull() {
		addPathError(
			&diags,
			path.Root("domain_name"),
			"Domain name cannot be parsed",
			srv.SetDomainName(m.DomainName.ValueString()),
		)
	}

	if !m.DomainSearchList.IsNull() {
		var domainSearchList []string
		diags.Append(m.DomainSearchList.ElementsAs(ctx, &domainSearchList, false)...)
		addPathError(
			&diags,
			path.Root("domain_search_list"),
			"Domain search list cannot be parsed",
			srv.SetDomainSearchList(domainSearchList),
		)
	}

	if !m.DefaultLeaseTime.IsNull() {
		addPathError(
			&diags,
			path.Root("default_lease_time"),
			"Default lease time cannot be parsed",
			srv.SetDefaultLeaseTime(m.DefaultLeaseTime.ValueString()),
		)
	}

	if !m.MaximumLeaseTime.IsNull() {
		addPathError(
			&diags,
			path.Root("maximum_lease_time"),
			"Maximum lease time cannot be parsed",
			srv.SetMaximumLeaseTime(m.MaximumLeaseTime.ValueString()),
		)
	}

	if !m.WINSServers.IsNull() {
		var winsServers []string
		diags.Append(m.WINSServers.ElementsAs(ctx, &winsServers, false)...)
		addPathError(
			&diags,
			path.Root("wins_servers"),
			"WINS servers cannot be parsed",
			srv.SetWINSServers(winsServers),
		)
	}

	if !m.NTPServers.IsNull() {
		var ntpServers []string
		diags.Append(m.NTPServers.ElementsAs(ctx, &ntpServers, false)...)
		addPathError(
			&diags,
			path.Root("ntp_servers"),
			"NTP servers cannot be parsed",
			srv.SetNTPServers(ntpServers),
		)
	}

	if !m.TFTPServer.IsNull() {
		addPathError(
			&diags,
			path.Root("tftp_server"),
			"TFTP server cannot be parsed",
			srv.SetTFTPServer(m.TFTPServer.ValueString()),
		)
	}

	if !m.LDAPServer.IsNull() {
		addPathError(
			&diags,
			path.Root("ldap_server"),
			"LDAP server cannot be parsed",
			srv.SetLDAPServer(m.LDAPServer.ValueString()),
		)
	}

	if !m.MACAllow.IsNull() {
		addPathError(
			&diags,
			path.Root("mac_allow"),
			"MAC allow cannot be parsed",
			srv.SetMACAllow(m.MACAllow.ValueString()),
		)
	}

	if !m.MACDeny.IsNull() {
		addPathError(
			&diags,
			path.Root("mac_deny"),
			"MAC deny cannot be parsed",
			srv.SetMACDeny(m.MACDeny.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("deny_unknown"),
		"Deny unknown cannot be parsed",
		srv.SetDenyUnknown(m.DenyUnknown.ValueBool()),
	)

	addPathError(
		&diags,
		path.Root("ignore_client_uids"),
		"Ignore client UIDs cannot be parsed",
		srv.SetIgnoreClientUIDs(m.IgnoreClientUIDs.ValueBool()),
	)

	addPathError(
		&diags,
		path.Root("static_arp"),
		"Static ARP cannot be parsed",
		srv.SetStaticARP(m.StaticARP.ValueBool()),
	)

	return diags
}
