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

type FirewallNATOneToOneRuleModel struct {
	Interface          types.String `tfsdk:"interface"`
	External           types.String `tfsdk:"external"`
	IPProtocol         types.String `tfsdk:"ipprotocol"`
	SourceAddress      types.String `tfsdk:"source_address"`
	SourceNot          types.Bool   `tfsdk:"source_not"`
	DestinationAddress types.String `tfsdk:"destination_address"`
	DestinationNot     types.Bool   `tfsdk:"destination_not"`
	Disabled           types.Bool   `tfsdk:"disabled"`
	NoBinat            types.Bool   `tfsdk:"no_binat"`
	Description        types.String `tfsdk:"description"`
	NATReflection      types.String `tfsdk:"nat_reflection"`
}

func (FirewallNATOneToOneRuleModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"interface": {
			Description: "Network interface this 1:1 NAT rule applies to (e.g. 'wan', 'lan', 'opt1').",
		},
		"external": {
			Description: "External subnet IP for the 1:1 mapping. Enter the external (usually on a WAN) subnet for the 1:1 mapping (e.g. '10.0.0.1' or '10.0.0.0/24').",
		},
		"ipprotocol": {
			Description:         fmt.Sprintf("IP protocol. Options: %s.", wrapElementsJoin(pfsense.NATOneToOne{}.IPProtocols(), "'")),
			MarkdownDescription: fmt.Sprintf("IP protocol. Options: %s.", wrapElementsJoin(pfsense.NATOneToOne{}.IPProtocols(), "`")),
		},
		"source_address": {
			Description: "Source network. Can be 'any', a CIDR network (e.g. '192.168.1.0/24'), an alias, or a special pfSense interface address. Defaults to 'any'.",
		},
		"source_not": {
			Description: "Invert the source address match.",
		},
		"destination_address": {
			Description: "Destination network. Can be 'any', a CIDR network, an alias, or a special pfSense interface address. Defaults to 'any'.",
		},
		"destination_not": {
			Description: "Invert the destination address match.",
		},
		"disabled": {
			Description: "Disable this 1:1 NAT rule.",
		},
		"no_binat": {
			Description: "Do not create a binat entry for this rule (one-way 1:1 NAT).",
		},
		"description": {
			Description: "Description used as the unique identifier for this NAT 1:1 rule.",
		},
		"nat_reflection": {
			Description:         fmt.Sprintf("NAT reflection mode. Options: %s. Empty string uses system default.", wrapElementsJoin(pfsense.NATOneToOne{}.NATReflectionModes(), "'")),
			MarkdownDescription: fmt.Sprintf("NAT reflection mode. Options: %s. Empty string uses system default.", wrapElementsJoin(pfsense.NATOneToOne{}.NATReflectionModes(), "`")),
		},
	}
}

func (FirewallNATOneToOneRuleModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interface":           types.StringType,
		"external":            types.StringType,
		"ipprotocol":          types.StringType,
		"source_address":      types.StringType,
		"source_not":          types.BoolType,
		"destination_address": types.StringType,
		"destination_not":     types.BoolType,
		"disabled":            types.BoolType,
		"no_binat":            types.BoolType,
		"description":         types.StringType,
		"nat_reflection":      types.StringType,
	}
}

func (m *FirewallNATOneToOneRuleModel) Set(_ context.Context, r pfsense.NATOneToOne) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Interface = types.StringValue(r.Interface)
	m.External = types.StringValue(r.External)
	m.IPProtocol = types.StringValue(r.IPProtocol)
	m.SourceAddress = types.StringValue(r.SourceAddress)
	m.SourceNot = types.BoolValue(r.SourceNot)
	m.DestinationAddress = types.StringValue(r.DestinationAddress)
	m.DestinationNot = types.BoolValue(r.DestinationNot)
	m.Disabled = types.BoolValue(r.Disabled)
	m.NoBinat = types.BoolValue(r.NoBinat)
	m.Description = types.StringValue(r.Description)

	if r.NATReflection != "" {
		m.NATReflection = types.StringValue(r.NATReflection)
	} else {
		m.NATReflection = types.StringNull()
	}

	return diags
}

func (m FirewallNATOneToOneRuleModel) Value(_ context.Context, r *pfsense.NATOneToOne) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(&diags, path.Root("interface"), "Interface cannot be parsed", r.SetInterface(m.Interface.ValueString()))
	addPathError(&diags, path.Root("external"), "External cannot be parsed", r.SetExternal(m.External.ValueString()))
	addPathError(&diags, path.Root("ipprotocol"), "IP protocol cannot be parsed", r.SetIPProtocol(m.IPProtocol.ValueString()))

	srcAddr := m.SourceAddress.ValueString()
	if m.SourceAddress.IsNull() || srcAddr == "" {
		srcAddr = "any"
	}

	addPathError(&diags, path.Root("source_address"), "Source address cannot be parsed", r.SetSourceAddress(srcAddr))
	addPathError(&diags, path.Root("source_not"), "Source not cannot be parsed", r.SetSourceNot(m.SourceNot.ValueBool()))

	dstAddr := m.DestinationAddress.ValueString()
	if m.DestinationAddress.IsNull() || dstAddr == "" {
		dstAddr = "any"
	}

	addPathError(&diags, path.Root("destination_address"), "Destination address cannot be parsed", r.SetDestinationAddress(dstAddr))
	addPathError(&diags, path.Root("destination_not"), "Destination not cannot be parsed", r.SetDestinationNot(m.DestinationNot.ValueBool()))
	addPathError(&diags, path.Root("disabled"), "Disabled cannot be parsed", r.SetDisabled(m.Disabled.ValueBool()))
	addPathError(&diags, path.Root("no_binat"), "No binat cannot be parsed", r.SetNoBinat(m.NoBinat.ValueBool()))
	addPathError(&diags, path.Root("description"), "Description cannot be parsed", r.SetDescription(m.Description.ValueString()))

	if !m.NATReflection.IsNull() {
		addPathError(&diags, path.Root("nat_reflection"), "NAT reflection cannot be parsed", r.SetNATReflection(m.NATReflection.ValueString()))
	}

	return diags
}
