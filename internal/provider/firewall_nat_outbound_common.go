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

type FirewallNATOutboundRuleModel struct {
	Interface      types.String `tfsdk:"interface"`
	Protocol       types.String `tfsdk:"protocol"`
	SourceAddress  types.String `tfsdk:"source_address"`
	SourcePort     types.String `tfsdk:"source_port"`
	SourceNot      types.Bool   `tfsdk:"source_not"`
	DestAddress    types.String `tfsdk:"destination_address"`
	DestPort       types.String `tfsdk:"destination_port"`
	DestNot        types.Bool   `tfsdk:"destination_not"`
	Target         types.String `tfsdk:"target"`
	TargetIP       types.String `tfsdk:"target_ip"`
	TargetIPSubnet types.String `tfsdk:"target_ip_subnet"`
	NATPort        types.String `tfsdk:"nat_port"`
	PoolOpts       types.String `tfsdk:"pool_options"`
	SourceHashKey  types.String `tfsdk:"source_hash_key"`
	StaticNATPort  types.Bool   `tfsdk:"static_nat_port"`
	NoSync         types.Bool   `tfsdk:"no_sync"`
	NoNAT          types.Bool   `tfsdk:"no_nat"`
	Disabled       types.Bool   `tfsdk:"disabled"`
	Description    types.String `tfsdk:"description"`
}

func (FirewallNATOutboundRuleModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"interface": {
			Description: "Network interface this outbound NAT rule applies to (e.g. 'wan', 'lan', 'opt1').",
		},
		"protocol": {
			Description:         fmt.Sprintf("Protocol. Options: %s. Empty string matches any protocol.", wrapElementsJoin(pfsense.NATOutboundRule{}.Protocols(), "'")),
			MarkdownDescription: fmt.Sprintf("Protocol. Options: %s. Empty string matches any protocol.", wrapElementsJoin(pfsense.NATOutboundRule{}.Protocols(), "`")),
		},
		"source_address": {
			Description: "Source network. Can be 'any', a CIDR network (e.g. '192.168.1.0/24'), an alias, or a special pfSense interface address. Defaults to 'any'.",
		},
		"source_port": {
			Description: "Source port or port range. Only applicable for TCP, UDP, or TCP/UDP protocols.",
		},
		"source_not": {
			Description: "Invert the source address match.",
		},
		"destination_address": {
			Description: "Destination network. Can be 'any', a CIDR network, an alias, or a special pfSense interface address. Defaults to 'any'.",
		},
		"destination_port": {
			Description: "Destination port or port range. Only applicable for TCP, UDP, or TCP/UDP protocols.",
		},
		"destination_not": {
			Description: "Invert the destination address match.",
		},
		"target": {
			Description: "Translation address type. Empty string means interface address, '(self)' for self, 'other-subnet' for other subnet, or a specific IP address.",
		},
		"target_ip": {
			Description: "Translation IP address when target is a specific IP or 'other-subnet'.",
		},
		"target_ip_subnet": {
			Description: "Translation subnet bits when target is 'other-subnet' (e.g. '24').",
		},
		"nat_port": {
			Description: "Translation port. If empty, the original port is preserved.",
		},
		"pool_options": {
			Description:         fmt.Sprintf("Pool options for translation address. Options: %s.", wrapElementsJoin(pfsense.NATOutboundRule{}.PoolOptions(), "'")),
			MarkdownDescription: fmt.Sprintf("Pool options for translation address. Options: %s.", wrapElementsJoin(pfsense.NATOutboundRule{}.PoolOptions(), "`")),
		},
		"source_hash_key": {
			Description: "Source hash key (hex string) used when pool_options is 'source-hash'.",
		},
		"static_nat_port": {
			Description: "Do not randomize the source port for outgoing connections (static port).",
		},
		"no_sync": {
			Description: "Do not synchronize this rule to HA peer.",
		},
		"no_nat": {
			Description: "Disable NAT for matching traffic (exclusion rule).",
		},
		"disabled": {
			Description: "Disable this outbound NAT rule.",
		},
		"description": {
			Description: "Description used as the unique identifier for this NAT outbound rule.",
		},
	}
}

func (FirewallNATOutboundRuleModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interface":           types.StringType,
		"protocol":            types.StringType,
		"source_address":      types.StringType,
		"source_port":         types.StringType,
		"source_not":          types.BoolType,
		"destination_address": types.StringType,
		"destination_port":    types.StringType,
		"destination_not":     types.BoolType,
		"target":              types.StringType,
		"target_ip":           types.StringType,
		"target_ip_subnet":    types.StringType,
		"nat_port":            types.StringType,
		"pool_options":        types.StringType,
		"source_hash_key":     types.StringType,
		"static_nat_port":     types.BoolType,
		"no_sync":             types.BoolType,
		"no_nat":              types.BoolType,
		"disabled":            types.BoolType,
		"description":         types.StringType,
	}
}

func (m *FirewallNATOutboundRuleModel) Set(_ context.Context, r pfsense.NATOutboundRule) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Interface = types.StringValue(r.Interface)
	m.Protocol = types.StringValue(r.Protocol)
	m.SourceAddress = types.StringValue(r.SourceAddress)

	if r.SourcePort != "" {
		m.SourcePort = types.StringValue(r.SourcePort)
	} else {
		m.SourcePort = types.StringNull()
	}

	m.SourceNot = types.BoolValue(r.SourceNot)
	m.DestAddress = types.StringValue(r.DestAddress)

	if r.DestPort != "" {
		m.DestPort = types.StringValue(r.DestPort)
	} else {
		m.DestPort = types.StringNull()
	}

	m.DestNot = types.BoolValue(r.DestNot)

	m.Target = types.StringValue(r.Target)

	if r.TargetIP != "" {
		m.TargetIP = types.StringValue(r.TargetIP)
	} else {
		m.TargetIP = types.StringNull()
	}

	if r.TargetIPSubnet != "" {
		m.TargetIPSubnet = types.StringValue(r.TargetIPSubnet)
	} else {
		m.TargetIPSubnet = types.StringNull()
	}

	if r.NATPort != "" {
		m.NATPort = types.StringValue(r.NATPort)
	} else {
		m.NATPort = types.StringNull()
	}

	if r.PoolOpts != "" {
		m.PoolOpts = types.StringValue(r.PoolOpts)
	} else {
		m.PoolOpts = types.StringNull()
	}

	if r.SourceHashKey != "" {
		m.SourceHashKey = types.StringValue(r.SourceHashKey)
	} else {
		m.SourceHashKey = types.StringNull()
	}

	m.StaticNATPort = types.BoolValue(r.StaticNATPort)
	m.NoSync = types.BoolValue(r.NoSync)
	m.NoNAT = types.BoolValue(r.NoNAT)
	m.Disabled = types.BoolValue(r.Disabled)
	m.Description = types.StringValue(r.Description)

	return diags
}

func (m FirewallNATOutboundRuleModel) Value(_ context.Context, r *pfsense.NATOutboundRule) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(&diags, path.Root("interface"), "Interface cannot be parsed", r.SetInterface(m.Interface.ValueString()))
	addPathError(&diags, path.Root("protocol"), "Protocol cannot be parsed", r.SetProtocol(m.Protocol.ValueString()))

	srcAddr := m.SourceAddress.ValueString()
	if m.SourceAddress.IsNull() || srcAddr == "" {
		srcAddr = "any"
	}

	addPathError(&diags, path.Root("source_address"), "Source address cannot be parsed", r.SetSourceAddress(srcAddr))

	if !m.SourcePort.IsNull() {
		addPathError(&diags, path.Root("source_port"), "Source port cannot be parsed", r.SetSourcePort(m.SourcePort.ValueString()))
	}

	addPathError(&diags, path.Root("source_not"), "Source not cannot be parsed", r.SetSourceNot(m.SourceNot.ValueBool()))

	dstAddr := m.DestAddress.ValueString()
	if m.DestAddress.IsNull() || dstAddr == "" {
		dstAddr = "any"
	}

	addPathError(&diags, path.Root("destination_address"), "Destination address cannot be parsed", r.SetDestAddress(dstAddr))

	if !m.DestPort.IsNull() {
		addPathError(&diags, path.Root("destination_port"), "Destination port cannot be parsed", r.SetDestPort(m.DestPort.ValueString()))
	}

	addPathError(&diags, path.Root("destination_not"), "Destination not cannot be parsed", r.SetDestNot(m.DestNot.ValueBool()))
	addPathError(&diags, path.Root("target"), "Target cannot be parsed", r.SetTarget(m.Target.ValueString()))

	if !m.TargetIP.IsNull() {
		addPathError(&diags, path.Root("target_ip"), "Target IP cannot be parsed", r.SetTargetIP(m.TargetIP.ValueString()))
	}

	if !m.TargetIPSubnet.IsNull() {
		addPathError(&diags, path.Root("target_ip_subnet"), "Target IP subnet cannot be parsed", r.SetTargetIPSubnet(m.TargetIPSubnet.ValueString()))
	}

	if !m.NATPort.IsNull() {
		addPathError(&diags, path.Root("nat_port"), "NAT port cannot be parsed", r.SetNATPort(m.NATPort.ValueString()))
	}

	if !m.PoolOpts.IsNull() {
		addPathError(&diags, path.Root("pool_options"), "Pool options cannot be parsed", r.SetPoolOpts(m.PoolOpts.ValueString()))
	}

	if !m.SourceHashKey.IsNull() {
		addPathError(&diags, path.Root("source_hash_key"), "Source hash key cannot be parsed", r.SetSourceHashKey(m.SourceHashKey.ValueString()))
	}

	addPathError(&diags, path.Root("static_nat_port"), "Static NAT port cannot be parsed", r.SetStaticNATPort(m.StaticNATPort.ValueBool()))
	addPathError(&diags, path.Root("no_sync"), "No sync cannot be parsed", r.SetNoSync(m.NoSync.ValueBool()))
	addPathError(&diags, path.Root("no_nat"), "No NAT cannot be parsed", r.SetNoNAT(m.NoNAT.ValueBool()))
	addPathError(&diags, path.Root("disabled"), "Disabled cannot be parsed", r.SetDisabled(m.Disabled.ValueBool()))
	addPathError(&diags, path.Root("description"), "Description cannot be parsed", r.SetDescription(m.Description.ValueString()))

	return diags
}
