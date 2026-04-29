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

type FirewallNATPortForwardModel struct {
	Interface        types.String `tfsdk:"interface"`
	IPProtocol       types.String `tfsdk:"ipprotocol"`
	Protocol         types.String `tfsdk:"protocol"`
	SourceAddress    types.String `tfsdk:"source_address"`
	SourcePort       types.String `tfsdk:"source_port"`
	SourceNot        types.Bool   `tfsdk:"source_not"`
	DestAddress      types.String `tfsdk:"destination_address"`
	DestPort         types.String `tfsdk:"destination_port"`
	DestNot          types.Bool   `tfsdk:"destination_not"`
	Target           types.String `tfsdk:"target"`
	LocalPort        types.String `tfsdk:"local_port"`
	Description      types.String `tfsdk:"description"`
	Disabled         types.Bool   `tfsdk:"disabled"`
	NoRDR            types.Bool   `tfsdk:"no_rdr"`
	NATReflection    types.String `tfsdk:"nat_reflection"`
	AssociatedRuleID types.String `tfsdk:"associated_rule_id"`
}

func (FirewallNATPortForwardModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"interface": {
			Description: "Network interface this NAT rule applies to (e.g. 'wan', 'lan', 'opt1').",
		},
		"ipprotocol": {
			Description:         fmt.Sprintf("IP address family. Options: %s.", wrapElementsJoin(pfsense.NATPortForward{}.IPProtocols(), "'")),
			MarkdownDescription: fmt.Sprintf("IP address family. Options: %s.", wrapElementsJoin(pfsense.NATPortForward{}.IPProtocols(), "`")),
		},
		"protocol": {
			Description:         fmt.Sprintf("Protocol. Options: %s.", wrapElementsJoin(pfsense.NATPortForward{}.Protocols(), "'")),
			MarkdownDescription: fmt.Sprintf("Protocol. Options: %s.", wrapElementsJoin(pfsense.NATPortForward{}.Protocols(), "`")),
		},
		"source_address": {
			Description: "Source address. Can be 'any', a single IP, a CIDR network, an alias, or a special pfSense interface address. Defaults to 'any'.",
		},
		"source_port": {
			Description: "Source port or port range (e.g. '80' or '8000-9000'). Only applicable for TCP, UDP, or TCP/UDP protocols.",
		},
		"source_not": {
			Description: "Invert the source address match.",
		},
		"destination_address": {
			Description: "Destination address. Can be 'any', a single IP, a CIDR network, an alias, or a special pfSense interface address (e.g. 'wanip').",
		},
		"destination_port": {
			Description: "Destination port or port range. The external port that triggers this NAT rule.",
		},
		"destination_not": {
			Description: "Invert the destination address match.",
		},
		"target": {
			Description: "Redirect target IP address. The internal IP address to forward traffic to.",
		},
		"local_port": {
			Description: "Redirect target port. The internal port to forward traffic to. If empty, the destination port is used.",
		},
		"description": {
			Description: "Description used as the unique identifier for this NAT port forward rule.",
		},
		"disabled": {
			Description: "Disable this NAT port forward rule.",
		},
		"no_rdr": {
			Description: "Disable redirection (NOT mode). When enabled, the rule acts as a negation — matching traffic is NOT redirected.",
		},
		"nat_reflection": {
			Description:         fmt.Sprintf("NAT reflection mode. Options: %s. Empty string means system default.", wrapElementsJoin(pfsense.NATPortForward{}.NATReflectionModes(), "'")),
			MarkdownDescription: fmt.Sprintf("NAT reflection mode. Options: %s. Empty string means system default.", wrapElementsJoin(pfsense.NATPortForward{}.NATReflectionModes(), "`")),
		},
		"associated_rule_id": {
			Description: "Associated filter rule behavior. 'pass' creates an associated firewall rule automatically. Empty string means no associated rule (manual rule creation required).",
		},
	}
}

func (FirewallNATPortForwardModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interface":           types.StringType,
		"ipprotocol":          types.StringType,
		"protocol":            types.StringType,
		"source_address":      types.StringType,
		"source_port":         types.StringType,
		"source_not":          types.BoolType,
		"destination_address": types.StringType,
		"destination_port":    types.StringType,
		"destination_not":     types.BoolType,
		"target":              types.StringType,
		"local_port":          types.StringType,
		"description":         types.StringType,
		"disabled":            types.BoolType,
		"no_rdr":              types.BoolType,
		"nat_reflection":      types.StringType,
		"associated_rule_id":  types.StringType,
	}
}

func (m *FirewallNATPortForwardModel) Set(_ context.Context, r pfsense.NATPortForward) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Interface = types.StringValue(r.Interface)
	m.IPProtocol = types.StringValue(r.IPProtocol)
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

	if r.LocalPort != "" {
		m.LocalPort = types.StringValue(r.LocalPort)
	} else {
		m.LocalPort = types.StringNull()
	}

	m.Description = types.StringValue(r.Description)
	m.Disabled = types.BoolValue(r.Disabled)
	m.NoRDR = types.BoolValue(r.NoRDR)

	if r.NATReflection != "" {
		m.NATReflection = types.StringValue(r.NATReflection)
	} else {
		m.NATReflection = types.StringNull()
	}

	m.AssociatedRuleID = types.StringValue(r.AssociatedRuleID)

	return diags
}

func (m FirewallNATPortForwardModel) Value(_ context.Context, r *pfsense.NATPortForward) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("interface"),
		"Interface cannot be parsed",
		r.SetInterface(m.Interface.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("ipprotocol"),
		"IP protocol cannot be parsed",
		r.SetIPProtocol(m.IPProtocol.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("protocol"),
		"Protocol cannot be parsed",
		r.SetProtocol(m.Protocol.ValueString()),
	)

	srcAddr := m.SourceAddress.ValueString()
	if m.SourceAddress.IsNull() || srcAddr == "" {
		srcAddr = "any"
	}

	addPathError(
		&diags,
		path.Root("source_address"),
		"Source address cannot be parsed",
		r.SetSourceAddress(srcAddr),
	)

	if !m.SourcePort.IsNull() {
		addPathError(
			&diags,
			path.Root("source_port"),
			"Source port cannot be parsed",
			r.SetSourcePort(m.SourcePort.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("source_not"),
		"Source not cannot be parsed",
		r.SetSourceNot(m.SourceNot.ValueBool()),
	)

	addPathError(
		&diags,
		path.Root("destination_address"),
		"Destination address cannot be parsed",
		r.SetDestAddress(m.DestAddress.ValueString()),
	)

	if !m.DestPort.IsNull() {
		addPathError(
			&diags,
			path.Root("destination_port"),
			"Destination port cannot be parsed",
			r.SetDestPort(m.DestPort.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("destination_not"),
		"Destination not cannot be parsed",
		r.SetDestNot(m.DestNot.ValueBool()),
	)

	addPathError(
		&diags,
		path.Root("target"),
		"Target cannot be parsed",
		r.SetTarget(m.Target.ValueString()),
	)

	if !m.LocalPort.IsNull() {
		addPathError(
			&diags,
			path.Root("local_port"),
			"Local port cannot be parsed",
			r.SetLocalPort(m.LocalPort.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("description"),
		"Description cannot be parsed",
		r.SetDescription(m.Description.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("disabled"),
		"Disabled cannot be parsed",
		r.SetDisabled(m.Disabled.ValueBool()),
	)

	addPathError(
		&diags,
		path.Root("no_rdr"),
		"No RDR cannot be parsed",
		r.SetNoRDR(m.NoRDR.ValueBool()),
	)

	if !m.NATReflection.IsNull() {
		addPathError(
			&diags,
			path.Root("nat_reflection"),
			"NAT reflection cannot be parsed",
			r.SetNATReflection(m.NATReflection.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("associated_rule_id"),
		"Associated rule ID cannot be parsed",
		r.SetAssociatedRuleID(m.AssociatedRuleID.ValueString()),
	)

	return diags
}
