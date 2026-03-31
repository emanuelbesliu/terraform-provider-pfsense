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

type FirewallRuleModel struct {
	Tracker       types.String `tfsdk:"tracker"`
	Type          types.String `tfsdk:"type"`
	Interface     types.String `tfsdk:"interface"`
	IPProtocol    types.String `tfsdk:"ipprotocol"`
	Protocol      types.String `tfsdk:"protocol"`
	SourceAddress types.String `tfsdk:"source_address"`
	SourcePort    types.String `tfsdk:"source_port"`
	SourceNot     types.Bool   `tfsdk:"source_not"`
	DestAddress   types.String `tfsdk:"destination_address"`
	DestPort      types.String `tfsdk:"destination_port"`
	DestNot       types.Bool   `tfsdk:"destination_not"`
	Description   types.String `tfsdk:"description"`
	Disabled      types.Bool   `tfsdk:"disabled"`
	Log           types.Bool   `tfsdk:"log"`
}

func (FirewallRuleModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"tracker": {
			Description: "Unique tracking ID for the firewall rule. Assigned automatically by pfSense.",
		},
		"type": {
			Description:         fmt.Sprintf("Rule action type. Options: %s.", wrapElementsJoin(pfsense.FirewallRule{}.Types(), "'")),
			MarkdownDescription: fmt.Sprintf("Rule action type. Options: %s.", wrapElementsJoin(pfsense.FirewallRule{}.Types(), "`")),
		},
		"interface": {
			Description: "Network interface this rule applies to.",
		},
		"ipprotocol": {
			Description:         fmt.Sprintf("IP address family. Options: %s.", wrapElementsJoin(pfsense.FirewallRule{}.IPProtocols(), "'")),
			MarkdownDescription: fmt.Sprintf("IP address family. Options: %s.", wrapElementsJoin(pfsense.FirewallRule{}.IPProtocols(), "`")),
		},
		"protocol": {
			Description:         fmt.Sprintf("Protocol. Options: %s. Defaults to 'any'.", wrapElementsJoin(pfsense.FirewallRule{}.Protocols(), "'")),
			MarkdownDescription: fmt.Sprintf("Protocol. Options: %s. Defaults to `any`.", wrapElementsJoin(pfsense.FirewallRule{}.Protocols(), "`")),
		},
		"source_address": {
			Description: "Source address. Can be 'any', a single IP address, a CIDR network (e.g. '10.0.0.0/24'), an alias name, or a special pfSense interface address (e.g. 'lanip', 'wanip'). Defaults to 'any'.",
		},
		"source_port": {
			Description: "Source port or port range (e.g. '80' or '8000-9000'). Only applicable for TCP, UDP, or TCP/UDP protocols.",
		},
		"source_not": {
			Description: "Invert the source address match.",
		},
		"destination_address": {
			Description: "Destination address. Can be 'any', a single IP address, a CIDR network (e.g. '10.0.0.0/24'), an alias name, or a special pfSense interface address (e.g. 'lanip', 'wanip'). Defaults to 'any'.",
		},
		"destination_port": {
			Description: "Destination port or port range (e.g. '443' or '8000-9000'). Only applicable for TCP, UDP, or TCP/UDP protocols.",
		},
		"destination_not": {
			Description: "Invert the destination address match.",
		},
		"description": {
			Description: descriptionDescription,
		},
		"disabled": {
			Description: "Disable this firewall rule.",
		},
		"log": {
			Description: "Log packets matched by this rule.",
		},
	}
}

func (FirewallRuleModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tracker":             types.StringType,
		"type":                types.StringType,
		"interface":           types.StringType,
		"ipprotocol":          types.StringType,
		"protocol":            types.StringType,
		"source_address":      types.StringType,
		"source_port":         types.StringType,
		"source_not":          types.BoolType,
		"destination_address": types.StringType,
		"destination_port":    types.StringType,
		"destination_not":     types.BoolType,
		"description":         types.StringType,
		"disabled":            types.BoolType,
		"log":                 types.BoolType,
	}
}

func (m *FirewallRuleModel) Set(_ context.Context, r pfsense.FirewallRule) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Tracker = types.StringValue(r.Tracker)
	m.Type = types.StringValue(r.Type)
	m.Interface = types.StringValue(r.Interface)
	m.IPProtocol = types.StringValue(r.IPProtocol)
	m.Protocol = types.StringValue(r.Protocol)
	m.SourceAddress = types.StringValue(r.SourceAddress)

	if r.SourcePort != "" {
		m.SourcePort = types.StringValue(r.SourcePort)
	}

	m.SourceNot = types.BoolValue(r.SourceNot)
	m.DestAddress = types.StringValue(r.DestAddress)

	if r.DestPort != "" {
		m.DestPort = types.StringValue(r.DestPort)
	}

	m.DestNot = types.BoolValue(r.DestNot)

	if r.Description != "" {
		m.Description = types.StringValue(r.Description)
	}

	m.Disabled = types.BoolValue(r.Disabled)
	m.Log = types.BoolValue(r.Log)

	return diags
}

func (m FirewallRuleModel) Value(_ context.Context, r *pfsense.FirewallRule) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("type"),
		"Type cannot be parsed",
		r.SetType(m.Type.ValueString()),
	)

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

	addPathError(
		&diags,
		path.Root("source_address"),
		"Source address cannot be parsed",
		r.SetSourceAddress(m.SourceAddress.ValueString()),
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

	if !m.Description.IsNull() {
		addPathError(
			&diags,
			path.Root("description"),
			"Description cannot be parsed",
			r.SetDescription(m.Description.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("disabled"),
		"Disabled cannot be parsed",
		r.SetDisabled(m.Disabled.ValueBool()),
	)

	addPathError(
		&diags,
		path.Root("log"),
		"Log cannot be parsed",
		r.SetLog(m.Log.ValueBool()),
	)

	if !m.Tracker.IsNull() && !m.Tracker.IsUnknown() {
		addPathError(
			&diags,
			path.Root("tracker"),
			"Tracker cannot be parsed",
			r.SetTracker(m.Tracker.ValueString()),
		)
	}

	return diags
}
