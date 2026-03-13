package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

// SystemAdvancedNetworkingModel represents the Terraform model for system advanced networking settings.
type SystemAdvancedNetworkingModel struct {
	// DHCP Options
	DHCPBackend      types.String `tfsdk:"dhcp_backend"`
	IgnoreISCWarning types.Bool   `tfsdk:"ignore_isc_warning"`
	RADVDDebug       types.Bool   `tfsdk:"radvd_debug"`
	DHCP6Debug       types.Bool   `tfsdk:"dhcp6_debug"`
	DHCP6NoRelease   types.Bool   `tfsdk:"dhcp6_no_release"`
	GlobalV6DUID     types.String `tfsdk:"global_v6_duid"`

	// IPv6 Options
	IPv6Allow              types.Bool   `tfsdk:"ipv6_allow"`
	IPv6NATEnable          types.Bool   `tfsdk:"ipv6_nat_enable"`
	IPv6NATIPAddr          types.String `tfsdk:"ipv6_nat_ip_address"`
	PreferIPv4             types.Bool   `tfsdk:"prefer_ipv4"`
	IPv6DontCreateLocalDNS types.Bool   `tfsdk:"ipv6_dont_create_local_dns"`

	// Network Interfaces
	DisableChecksumOffloading     types.Bool `tfsdk:"disable_checksum_offloading"`
	DisableSegmentationOffloading types.Bool `tfsdk:"disable_segmentation_offloading"`
	DisableLargeReceiveOffloading types.Bool `tfsdk:"disable_large_receive_offloading"`
	HNALTQEnable                  types.Bool `tfsdk:"hn_altq_enable"`
	SharedNet                     types.Bool `tfsdk:"suppress_arp_messages"`
	IPChangeKillStates            types.Bool `tfsdk:"ip_change_kill_states"`
	UseIfPPPoE                    types.Bool `tfsdk:"use_if_pppoe"`
}

func (SystemAdvancedNetworkingModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"dhcp_backend": {
			Description:         fmt.Sprintf("DHCP server backend. Options: %s. Defaults to '%s'. ISC DHCP has reached end-of-life and will be removed from a future version of pfSense.", wrapElementsJoin(pfsense.AdvancedNetworking{}.DHCPBackendOptions(), "'"), pfsense.DefaultAdvancedNetworkingDHCPBackend),
			MarkdownDescription: fmt.Sprintf("DHCP server backend. Options: %s. Defaults to `%s`. ISC DHCP has reached end-of-life and will be removed from a future version of pfSense.", wrapElementsJoin(pfsense.AdvancedNetworking{}.DHCPBackendOptions(), "`"), pfsense.DefaultAdvancedNetworkingDHCPBackend),
		},
		"ignore_isc_warning": {
			Description: "Ignore the ISC DHCP deprecation warning.",
		},
		"radvd_debug": {
			Description: "Log all radvd log levels for debugging.",
		},
		"dhcp6_debug": {
			Description: "Start the DHCP6 client in debug mode.",
		},
		"dhcp6_no_release": {
			Description: "Do not allow PD/Address release. Prevents dhcp6c from sending a release to the ISP on exit, which some ISPs use to release the allocated address or prefix.",
		},
		"global_v6_duid": {
			Description: "DHCPv6 Unique Identifier (DUID) in raw colon-separated hex format. Used by the firewall when requesting an IPv6 address. If not set, a dynamic DUID-LLT is automatically created.",
		},
		"ipv6_allow": {
			Description: "Allow IPv6 traffic. When unchecked, all IPv6 traffic is blocked by the firewall. This does not disable IPv6 features on the firewall itself.",
		},
		"ipv6_nat_enable": {
			Description: "Enable IPv6 over IPv4 tunneling (RFC 2893). Creates a mechanism for IPv4 NAT encapsulation of IPv6 packets.",
		},
		"ipv6_nat_ip_address": {
			Description: "IPv4 address of the tunnel peer for IPv6 over IPv4 tunneling. Only applicable when ipv6_nat_enable is true.",
		},
		"prefer_ipv4": {
			Description: "Prefer to use IPv4 even if IPv6 is available. By default, if IPv6 is configured and a hostname resolves both IPv6 and IPv4 addresses, IPv6 will be used.",
		},
		"ipv6_dont_create_local_dns": {
			Description: "Do not generate local IPv6 DNS entries for LAN interfaces. Useful when a LAN interface's IPv6 configuration is set to Track and the tracked interface loses connectivity.",
		},
		"disable_checksum_offloading": {
			Description: "Disable hardware checksum offloading. Checksum offloading is broken in some hardware, particularly some Realtek cards. Requires a system reboot to take effect.",
		},
		"disable_segmentation_offloading": {
			Description: "Disable hardware TCP segmentation offloading (TSO, TSO4, TSO6). This offloading is broken in some hardware drivers. Requires a system reboot to take effect.",
		},
		"disable_large_receive_offloading": {
			Description: "Disable hardware large receive offloading (LRO). This offloading is broken in some hardware drivers. Requires a system reboot to take effect.",
		},
		"hn_altq_enable": {
			Description: "Enable ALTQ support for vtnet/hn NICs. Disables the multiqueue API and may reduce the system's capability to handle traffic. Requires a system reboot to take effect.",
		},
		"suppress_arp_messages": {
			Description: "Suppress ARP log messages when multiple interfaces reside on the same broadcast domain.",
		},
		"ip_change_kill_states": {
			Description: "Reset all states when a WAN IP address changes, instead of only states associated with the previous IP address.",
		},
		"use_if_pppoe": {
			Description: "Use the if_pppoe kernel module for PPPoE client connections instead of the deprecated mpd5. Changing this option interrupts connectivity for affected interfaces and requires a system reboot.",
		},
	}
}

func (SystemAdvancedNetworkingModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"dhcp_backend":                     types.StringType,
		"ignore_isc_warning":               types.BoolType,
		"radvd_debug":                      types.BoolType,
		"dhcp6_debug":                      types.BoolType,
		"dhcp6_no_release":                 types.BoolType,
		"global_v6_duid":                   types.StringType,
		"ipv6_allow":                       types.BoolType,
		"ipv6_nat_enable":                  types.BoolType,
		"ipv6_nat_ip_address":              types.StringType,
		"prefer_ipv4":                      types.BoolType,
		"ipv6_dont_create_local_dns":       types.BoolType,
		"disable_checksum_offloading":      types.BoolType,
		"disable_segmentation_offloading":  types.BoolType,
		"disable_large_receive_offloading": types.BoolType,
		"hn_altq_enable":                   types.BoolType,
		"suppress_arp_messages":            types.BoolType,
		"ip_change_kill_states":            types.BoolType,
		"use_if_pppoe":                     types.BoolType,
	}
}

func (m *SystemAdvancedNetworkingModel) Set(_ context.Context, a pfsense.AdvancedNetworking) diag.Diagnostics {
	var diags diag.Diagnostics

	// DHCP Options
	m.DHCPBackend = types.StringValue(a.DHCPBackend)
	m.IgnoreISCWarning = types.BoolValue(a.IgnoreISCWarning)
	m.RADVDDebug = types.BoolValue(a.RADVDDebug)
	m.DHCP6Debug = types.BoolValue(a.DHCP6Debug)
	m.DHCP6NoRelease = types.BoolValue(a.DHCP6NoRelease)

	if a.GlobalV6DUID != "" {
		m.GlobalV6DUID = types.StringValue(a.GlobalV6DUID)
	} else {
		m.GlobalV6DUID = types.StringNull()
	}

	// IPv6 Options
	m.IPv6Allow = types.BoolValue(a.IPv6Allow)
	m.IPv6NATEnable = types.BoolValue(a.IPv6NATEnable)

	if a.IPv6NATIPAddr != "" {
		m.IPv6NATIPAddr = types.StringValue(a.IPv6NATIPAddr)
	} else {
		m.IPv6NATIPAddr = types.StringNull()
	}

	m.PreferIPv4 = types.BoolValue(a.PreferIPv4)
	m.IPv6DontCreateLocalDNS = types.BoolValue(a.IPv6DontCreateLocalDNS)

	// Network Interfaces
	m.DisableChecksumOffloading = types.BoolValue(a.DisableChecksumOffloading)
	m.DisableSegmentationOffloading = types.BoolValue(a.DisableSegmentationOffloading)
	m.DisableLargeReceiveOffloading = types.BoolValue(a.DisableLargeReceiveOffloading)
	m.HNALTQEnable = types.BoolValue(a.HNALTQEnable)
	m.SharedNet = types.BoolValue(a.SharedNet)
	m.IPChangeKillStates = types.BoolValue(a.IPChangeKillStates)
	m.UseIfPPPoE = types.BoolValue(a.UseIfPPPoE)

	return diags
}

func (m SystemAdvancedNetworkingModel) Value(_ context.Context, a *pfsense.AdvancedNetworking) diag.Diagnostics {
	var diags diag.Diagnostics

	// DHCP Options
	addPathError(
		&diags,
		path.Root("dhcp_backend"),
		"DHCP backend cannot be parsed",
		a.SetDHCPBackend(m.DHCPBackend.ValueString()),
	)

	a.IgnoreISCWarning = m.IgnoreISCWarning.ValueBool()
	a.RADVDDebug = m.RADVDDebug.ValueBool()
	a.DHCP6Debug = m.DHCP6Debug.ValueBool()
	a.DHCP6NoRelease = m.DHCP6NoRelease.ValueBool()

	if !m.GlobalV6DUID.IsNull() {
		a.GlobalV6DUID = m.GlobalV6DUID.ValueString()
	}

	// IPv6 Options
	a.IPv6Allow = m.IPv6Allow.ValueBool()
	a.IPv6NATEnable = m.IPv6NATEnable.ValueBool()

	if !m.IPv6NATIPAddr.IsNull() {
		a.IPv6NATIPAddr = m.IPv6NATIPAddr.ValueString()
	}

	a.PreferIPv4 = m.PreferIPv4.ValueBool()
	a.IPv6DontCreateLocalDNS = m.IPv6DontCreateLocalDNS.ValueBool()

	// Network Interfaces
	a.DisableChecksumOffloading = m.DisableChecksumOffloading.ValueBool()
	a.DisableSegmentationOffloading = m.DisableSegmentationOffloading.ValueBool()
	a.DisableLargeReceiveOffloading = m.DisableLargeReceiveOffloading.ValueBool()
	a.HNALTQEnable = m.HNALTQEnable.ValueBool()
	a.SharedNet = m.SharedNet.ValueBool()
	a.IPChangeKillStates = m.IPChangeKillStates.ValueBool()
	a.UseIfPPPoE = m.UseIfPPPoE.ValueBool()

	return diags
}
