package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type InterfaceBridgeModel struct {
	BridgeIf           types.String `tfsdk:"bridge_if"`
	Members            types.List   `tfsdk:"members"`
	Description        types.String `tfsdk:"description"`
	EnableSTP          types.Bool   `tfsdk:"enable_stp"`
	IP6LinkLocal       types.Bool   `tfsdk:"ip6_link_local"`
	Protocol           types.String `tfsdk:"protocol"`
	Priority           types.Int64  `tfsdk:"priority"`
	HelloTime          types.Int64  `tfsdk:"hello_time"`
	ForwardDelay       types.Int64  `tfsdk:"forward_delay"`
	MaxAge             types.Int64  `tfsdk:"max_age"`
	HoldCount          types.Int64  `tfsdk:"hold_count"`
	MaxAddresses       types.Int64  `tfsdk:"max_addresses"`
	CacheExpire        types.Int64  `tfsdk:"cache_expire"`
	STPInterfaces      types.List   `tfsdk:"stp_interfaces"`
	StaticInterfaces   types.List   `tfsdk:"static_interfaces"`
	PrivateInterfaces  types.List   `tfsdk:"private_interfaces"`
	SpanInterfaces     types.List   `tfsdk:"span_interfaces"`
	EdgeInterfaces     types.List   `tfsdk:"edge_interfaces"`
	AutoEdgeInterfaces types.List   `tfsdk:"auto_edge_interfaces"`
	PTPInterfaces      types.List   `tfsdk:"ptp_interfaces"`
	AutoPTPInterfaces  types.List   `tfsdk:"auto_ptp_interfaces"`
}

func (InterfaceBridgeModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"bridge_if": {
			Description: "Interface name automatically assigned by pfSense to the bridge (e.g. 'bridge0'). Used as the unique identifier.",
		},
		"members": {
			Description: "List of logical interface names (e.g. 'lan', 'opt1') that are members of the bridge. At least one is required.",
		},
		"description": {
			Description: descriptionDescription,
		},
		"enable_stp": {
			Description: "Enable Spanning Tree Protocol (STP) on the bridge.",
		},
		"ip6_link_local": {
			Description: "Enable link-local IPv6 addressing on the bridge interface.",
		},
		"protocol": {
			Description: "Spanning Tree protocol version. Options: `rstp` (Rapid Spanning Tree), `stp` (legacy Spanning Tree).",
		},
		"priority": {
			Description: "Bridge priority for Spanning Tree, between 0 and 61440 in increments of 4096.",
		},
		"hello_time": {
			Description: "Time in seconds between broadcasting of Spanning Tree hello messages (RSTP only ignores this), between 1 and 2.",
		},
		"forward_delay": {
			Description: "Time in seconds that must pass before an interface begins forwarding packets when Spanning Tree is enabled, between 4 and 30.",
		},
		"max_age": {
			Description: "Maximum age in seconds of a Spanning Tree protocol message, between 6 and 40.",
		},
		"hold_count": {
			Description: "Transmit hold count for Spanning Tree, the number of packets transmitted before being rate limited, between 1 and 10.",
		},
		"max_addresses": {
			Description: "Maximum number of bridge address cache entries.",
		},
		"cache_expire": {
			Description: "Time in seconds before bridge address cache entries are expired, between 0 and 3600.",
		},
		"stp_interfaces": {
			Description: "Member interfaces on which Spanning Tree is enabled. Must be a subset of members.",
		},
		"static_interfaces": {
			Description: "Member interfaces with a static (sticky) address cache entry. Must be a subset of members.",
		},
		"private_interfaces": {
			Description: "Member interfaces marked private, which cannot communicate with any other private port. Must be a subset of members.",
		},
		"span_interfaces": {
			Description: "Member interfaces used as span ports, transmitting a copy of every frame received by the bridge. Span ports must NOT be bridge members.",
		},
		"edge_interfaces": {
			Description: "Member interfaces treated as edge ports, connecting directly to end stations and unable to create bridging loops. Must be a subset of members.",
		},
		"auto_edge_interfaces": {
			Description: "Member interfaces that automatically detect edge status. Must be a subset of members.",
		},
		"ptp_interfaces": {
			Description: "Member interfaces treated as point-to-point links. Must be a subset of members.",
		},
		"auto_ptp_interfaces": {
			Description: "Member interfaces that automatically detect point-to-point status. Must be a subset of members.",
		},
	}
}

func (InterfaceBridgeModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"bridge_if":            types.StringType,
		"members":              types.ListType{ElemType: types.StringType},
		"description":          types.StringType,
		"enable_stp":           types.BoolType,
		"ip6_link_local":       types.BoolType,
		"protocol":             types.StringType,
		"priority":             types.Int64Type,
		"hello_time":           types.Int64Type,
		"forward_delay":        types.Int64Type,
		"max_age":              types.Int64Type,
		"hold_count":           types.Int64Type,
		"max_addresses":        types.Int64Type,
		"cache_expire":         types.Int64Type,
		"stp_interfaces":       types.ListType{ElemType: types.StringType},
		"static_interfaces":    types.ListType{ElemType: types.StringType},
		"private_interfaces":   types.ListType{ElemType: types.StringType},
		"span_interfaces":      types.ListType{ElemType: types.StringType},
		"edge_interfaces":      types.ListType{ElemType: types.StringType},
		"auto_edge_interfaces": types.ListType{ElemType: types.StringType},
		"ptp_interfaces":       types.ListType{ElemType: types.StringType},
		"auto_ptp_interfaces":  types.ListType{ElemType: types.StringType},
	}
}

func bridgeStringListValue(values []string, diags *diag.Diagnostics) types.List {
	if len(values) == 0 {
		return types.ListNull(types.StringType)
	}

	elements := make([]attr.Value, 0, len(values))
	for _, value := range values {
		elements = append(elements, types.StringValue(value))
	}

	listValue, newDiags := types.ListValue(types.StringType, elements)
	diags.Append(newDiags...)

	return listValue
}

func bridgeInt64Value(value *int) types.Int64 {
	if value == nil {
		return types.Int64Null()
	}

	return types.Int64Value(int64(*value))
}

func (m *InterfaceBridgeModel) Set(_ context.Context, bridge pfsense.Bridge) diag.Diagnostics {
	var diags diag.Diagnostics

	m.BridgeIf = types.StringValue(bridge.BridgeIf)
	m.Members = bridgeStringListValue(bridge.Members, &diags)

	if bridge.Description != "" {
		m.Description = types.StringValue(bridge.Description)
	} else {
		m.Description = types.StringNull()
	}

	m.EnableSTP = types.BoolValue(bridge.EnableSTP)
	m.IP6LinkLocal = types.BoolValue(bridge.IP6LinkLocal)

	if bridge.Protocol != "" {
		m.Protocol = types.StringValue(bridge.Protocol)
	} else {
		m.Protocol = types.StringNull()
	}

	m.Priority = bridgeInt64Value(bridge.Priority)
	m.HelloTime = bridgeInt64Value(bridge.HelloTime)
	m.ForwardDelay = bridgeInt64Value(bridge.ForwardDelay)
	m.MaxAge = bridgeInt64Value(bridge.MaxAge)
	m.HoldCount = bridgeInt64Value(bridge.HoldCount)
	m.MaxAddresses = bridgeInt64Value(bridge.MaxAddresses)
	m.CacheExpire = bridgeInt64Value(bridge.CacheExpire)

	m.STPInterfaces = bridgeStringListValue(bridge.STPInterfaces, &diags)
	m.StaticInterfaces = bridgeStringListValue(bridge.StaticInterfaces, &diags)
	m.PrivateInterfaces = bridgeStringListValue(bridge.PrivateInterfaces, &diags)
	m.SpanInterfaces = bridgeStringListValue(bridge.SpanInterfaces, &diags)
	m.EdgeInterfaces = bridgeStringListValue(bridge.EdgeInterfaces, &diags)
	m.AutoEdgeInterfaces = bridgeStringListValue(bridge.AutoEdgeInterfaces, &diags)
	m.PTPInterfaces = bridgeStringListValue(bridge.PTPInterfaces, &diags)
	m.AutoPTPInterfaces = bridgeStringListValue(bridge.AutoPTPInterfaces, &diags)

	return diags
}

func bridgeListToStrings(ctx context.Context, list types.List, diags *diag.Diagnostics) []string {
	if list.IsNull() || list.IsUnknown() {
		return nil
	}

	var values []string
	diags.Append(list.ElementsAs(ctx, &values, false)...)

	return values
}

func bridgeInt64Pointer(value types.Int64) *int {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	parsed := int(value.ValueInt64())

	return &parsed
}

func (m InterfaceBridgeModel) Value(ctx context.Context, bridge *pfsense.Bridge) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(&diags, path.Root("members"), "Members cannot be parsed", bridge.SetMembers(bridgeListToStrings(ctx, m.Members, &diags)))

	if !m.Description.IsNull() {
		addPathError(&diags, path.Root("description"), "Description cannot be parsed", bridge.SetDescription(m.Description.ValueString()))
	}

	bridge.EnableSTP = m.EnableSTP.ValueBool()
	bridge.IP6LinkLocal = m.IP6LinkLocal.ValueBool()

	if !m.Protocol.IsNull() {
		addPathError(&diags, path.Root("protocol"), "Protocol cannot be parsed", bridge.SetProtocol(m.Protocol.ValueString()))
	}

	bridge.Priority = bridgeInt64Pointer(m.Priority)
	bridge.HelloTime = bridgeInt64Pointer(m.HelloTime)
	bridge.ForwardDelay = bridgeInt64Pointer(m.ForwardDelay)
	bridge.MaxAge = bridgeInt64Pointer(m.MaxAge)
	bridge.HoldCount = bridgeInt64Pointer(m.HoldCount)
	bridge.MaxAddresses = bridgeInt64Pointer(m.MaxAddresses)
	bridge.CacheExpire = bridgeInt64Pointer(m.CacheExpire)

	bridge.STPInterfaces = bridgeListToStrings(ctx, m.STPInterfaces, &diags)
	bridge.StaticInterfaces = bridgeListToStrings(ctx, m.StaticInterfaces, &diags)
	bridge.PrivateInterfaces = bridgeListToStrings(ctx, m.PrivateInterfaces, &diags)
	bridge.SpanInterfaces = bridgeListToStrings(ctx, m.SpanInterfaces, &diags)
	bridge.EdgeInterfaces = bridgeListToStrings(ctx, m.EdgeInterfaces, &diags)
	bridge.AutoEdgeInterfaces = bridgeListToStrings(ctx, m.AutoEdgeInterfaces, &diags)
	bridge.PTPInterfaces = bridgeListToStrings(ctx, m.PTPInterfaces, &diags)
	bridge.AutoPTPInterfaces = bridgeListToStrings(ctx, m.AutoPTPInterfaces, &diags)

	return diags
}
