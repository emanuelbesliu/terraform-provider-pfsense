package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type InterfaceModel struct {
	LogicalName types.String `tfsdk:"logical_name"`
	If          types.String `tfsdk:"port"`
	Description types.String `tfsdk:"description"`
	Enabled     types.Bool   `tfsdk:"enabled"`
	IPv4Type    types.String `tfsdk:"ipv4_type"`
	IPAddr      types.String `tfsdk:"ipv4_address"`
	Subnet      types.String `tfsdk:"ipv4_subnet"`
	Gateway     types.String `tfsdk:"ipv4_gateway"`
	IPv6Type    types.String `tfsdk:"ipv6_type"`
	IPAddrV6    types.String `tfsdk:"ipv6_address"`
	SubnetV6    types.String `tfsdk:"ipv6_subnet"`
	GatewayV6   types.String `tfsdk:"ipv6_gateway"`
	SpoofMAC    types.String `tfsdk:"spoof_mac"`
	MTU         types.Int64  `tfsdk:"mtu"`
	MSS         types.Int64  `tfsdk:"mss"`
	BlockPriv   types.Bool   `tfsdk:"block_private"`
	BlockBogons types.Bool   `tfsdk:"block_bogons"`
}

func (InterfaceModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"logical_name": {
			Description: "Logical name of the interface assignment (e.g. 'wan', 'lan', 'opt1'). Automatically assigned on creation.",
		},
		"port": {
			Description: "Physical or VLAN interface port name (e.g. 'vmx0', 'vmx2.100').",
		},
		"description": {
			Description: "Friendly description for the interface.",
		},
		"enabled": {
			Description: "Whether the interface is enabled.",
		},
		"ipv4_type": {
			Description:         fmt.Sprintf("IPv4 configuration type. Options: %s.", wrapElementsJoin(pfsense.InterfaceIPv4Types, "'")),
			MarkdownDescription: fmt.Sprintf("IPv4 configuration type. Options: %s.", wrapElementsJoin(pfsense.InterfaceIPv4Types, "`")),
		},
		"ipv4_address": {
			Description: "IPv4 address when ipv4_type is 'staticv4'.",
		},
		"ipv4_subnet": {
			Description: "IPv4 subnet prefix length (e.g. '24') when ipv4_type is 'staticv4'.",
		},
		"ipv4_gateway": {
			Description: "IPv4 gateway name for this interface.",
		},
		"ipv6_type": {
			Description:         fmt.Sprintf("IPv6 configuration type. Options: %s.", wrapElementsJoin(pfsense.InterfaceIPv6Types, "'")),
			MarkdownDescription: fmt.Sprintf("IPv6 configuration type. Options: %s.", wrapElementsJoin(pfsense.InterfaceIPv6Types, "`")),
		},
		"ipv6_address": {
			Description: "IPv6 address when ipv6_type is 'staticv6'.",
		},
		"ipv6_subnet": {
			Description: "IPv6 subnet prefix length (e.g. '64') when ipv6_type is 'staticv6'.",
		},
		"ipv6_gateway": {
			Description: "IPv6 gateway name for this interface.",
		},
		"spoof_mac": {
			Description: "Spoofed MAC address for this interface.",
		},
		"mtu": {
			Description: "Maximum Transmission Unit for this interface.",
		},
		"mss": {
			Description: "Maximum Segment Size for TCP connections on this interface.",
		},
		"block_private": {
			Description: "Block private networks (RFC 1918) from entering via this interface.",
		},
		"block_bogons": {
			Description: "Block bogon networks from entering via this interface.",
		},
	}
}

func (InterfaceModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"logical_name":  types.StringType,
		"port":          types.StringType,
		"description":   types.StringType,
		"enabled":       types.BoolType,
		"ipv4_type":     types.StringType,
		"ipv4_address":  types.StringType,
		"ipv4_subnet":   types.StringType,
		"ipv4_gateway":  types.StringType,
		"ipv6_type":     types.StringType,
		"ipv6_address":  types.StringType,
		"ipv6_subnet":   types.StringType,
		"ipv6_gateway":  types.StringType,
		"spoof_mac":     types.StringType,
		"mtu":           types.Int64Type,
		"mss":           types.Int64Type,
		"block_private": types.BoolType,
		"block_bogons":  types.BoolType,
	}
}

func (m *InterfaceModel) Set(_ context.Context, iface pfsense.Interface) diag.Diagnostics {
	var diags diag.Diagnostics

	m.LogicalName = types.StringValue(iface.LogicalName)
	m.If = types.StringValue(iface.If)

	if iface.Description != "" {
		m.Description = types.StringValue(iface.Description)
	} else {
		m.Description = types.StringNull()
	}

	m.Enabled = types.BoolValue(iface.Enabled)

	m.IPv4Type = types.StringValue(iface.IPv4Type)

	if iface.IPAddr != "" {
		m.IPAddr = types.StringValue(iface.IPAddr)
	} else {
		m.IPAddr = types.StringNull()
	}

	if iface.Subnet != "" {
		m.Subnet = types.StringValue(iface.Subnet)
	} else {
		m.Subnet = types.StringNull()
	}

	if iface.Gateway != "" {
		m.Gateway = types.StringValue(iface.Gateway)
	} else {
		m.Gateway = types.StringNull()
	}

	m.IPv6Type = types.StringValue(iface.IPv6Type)

	if iface.IPAddrV6 != "" {
		m.IPAddrV6 = types.StringValue(iface.IPAddrV6)
	} else {
		m.IPAddrV6 = types.StringNull()
	}

	if iface.SubnetV6 != "" {
		m.SubnetV6 = types.StringValue(iface.SubnetV6)
	} else {
		m.SubnetV6 = types.StringNull()
	}

	if iface.GatewayV6 != "" {
		m.GatewayV6 = types.StringValue(iface.GatewayV6)
	} else {
		m.GatewayV6 = types.StringNull()
	}

	if iface.SpoofMAC != "" {
		m.SpoofMAC = types.StringValue(iface.SpoofMAC)
	} else {
		m.SpoofMAC = types.StringNull()
	}

	if iface.MTU > 0 {
		m.MTU = types.Int64Value(int64(iface.MTU))
	} else {
		m.MTU = types.Int64Null()
	}

	if iface.MSS > 0 {
		m.MSS = types.Int64Value(int64(iface.MSS))
	} else {
		m.MSS = types.Int64Null()
	}

	m.BlockPriv = types.BoolValue(iface.BlockPriv)
	m.BlockBogons = types.BoolValue(iface.BlockBogons)

	return diags
}

func (m InterfaceModel) Value(_ context.Context, iface *pfsense.Interface) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("port"),
		"Port (if) cannot be parsed",
		iface.SetIf(m.If.ValueString()),
	)

	if !m.Description.IsNull() {
		addPathError(
			&diags,
			path.Root("description"),
			"Description cannot be parsed",
			iface.SetDescription(m.Description.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("enabled"),
		"Enabled cannot be parsed",
		iface.SetEnabled(m.Enabled.ValueBool()),
	)

	ipv4Type := "none"
	if !m.IPv4Type.IsNull() {
		ipv4Type = m.IPv4Type.ValueString()
	}

	addPathError(
		&diags,
		path.Root("ipv4_type"),
		"IPv4 type cannot be parsed",
		iface.SetIPv4Type(ipv4Type),
	)

	if !m.IPAddr.IsNull() {
		iface.IPAddr = m.IPAddr.ValueString()
	}

	if !m.Subnet.IsNull() {
		iface.Subnet = m.Subnet.ValueString()
	}

	if !m.Gateway.IsNull() {
		iface.Gateway = m.Gateway.ValueString()
	}

	ipv6Type := "none"
	if !m.IPv6Type.IsNull() {
		ipv6Type = m.IPv6Type.ValueString()
	}

	addPathError(
		&diags,
		path.Root("ipv6_type"),
		"IPv6 type cannot be parsed",
		iface.SetIPv6Type(ipv6Type),
	)

	if !m.IPAddrV6.IsNull() {
		iface.IPAddrV6 = m.IPAddrV6.ValueString()
	}

	if !m.SubnetV6.IsNull() {
		iface.SubnetV6 = m.SubnetV6.ValueString()
	}

	if !m.GatewayV6.IsNull() {
		iface.GatewayV6 = m.GatewayV6.ValueString()
	}

	if !m.SpoofMAC.IsNull() {
		iface.SpoofMAC = m.SpoofMAC.ValueString()
	}

	if !m.MTU.IsNull() {
		iface.MTU = int(m.MTU.ValueInt64())
	}

	if !m.MSS.IsNull() {
		iface.MSS = int(m.MSS.ValueInt64())
	}

	iface.BlockPriv = m.BlockPriv.ValueBool()
	iface.BlockBogons = m.BlockBogons.ValueBool()

	// Set logical name if available (used during updates).
	if !m.LogicalName.IsNull() && !m.LogicalName.IsUnknown() {
		iface.LogicalName = m.LogicalName.ValueString()
	}

	return diags
}

// interfaceIPv4Types returns the list of supported IPv4 types as a single-line string for validator messages.
func interfaceIPv4Types() string {
	return strings.Join(pfsense.InterfaceIPv4Types, ", ")
}

// interfaceIPv6Types returns the list of supported IPv6 types as a single-line string for validator messages.
func interfaceIPv6Types() string {
	return strings.Join(pfsense.InterfaceIPv6Types, ", ")
}
