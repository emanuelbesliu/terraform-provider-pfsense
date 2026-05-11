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

// FirewallNAT1to1Model represents the Terraform model for a 1:1 NAT rule
type FirewallNAT1to1Model struct {
	External      types.String `tfsdk:"external"`
	Interface     types.String `tfsdk:"interface"`
	IPProtocol    types.String `tfsdk:"ipprotocol"`
	SourceAddress types.String `tfsdk:"source_address"`
	SourceNot     types.Bool   `tfsdk:"source_not"`
	DestAddress   types.String `tfsdk:"destination_address"`
	DestNot       types.Bool   `tfsdk:"destination_not"`
	Description   types.String `tfsdk:"description"`
	Disabled      types.Bool   `tfsdk:"disabled"`
	NoBinat       types.Bool   `tfsdk:"no_binat"`
	NATReflection types.String `tfsdk:"nat_reflection"`
}

// descriptions returns attribute descriptions for the 1:1 NAT model
func (FirewallNAT1to1Model) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"external": {
			Description: "External IP address for the 1:1 NAT rule. This is the IP address that will be translated.",
		},
		"interface": {
			Description: "Network interface this NAT rule applies to (e.g. 'wan', 'lan', 'opt1').",
		},
		"ipprotocol": {
			Description:         fmt.Sprintf("IP address family. Options: %s.", wrapElementsJoin(pfsense.NATOneToOne{}.IPProtocols(), "'")),
			MarkdownDescription: fmt.Sprintf("IP address family. Options: %s.", wrapElementsJoin(pfsense.NATOneToOne{}.IPProtocols(), "`")),
		},
		"source_address": {
			Description: "Source address for the rule. Can be 'any', a single IP, a CIDR network, an alias, or a special pfSense interface address.",
		},
		"source_not": {
			Description: "Invert the source address match.",
		},
		"destination_address": {
			Description: "Destination address for the rule. Can be 'any', a single IP, a CIDR network, an alias, or a special pfSense interface address.",
		},
		"destination_not": {
			Description: "Invert the destination address match.",
		},
		"description": {
			Description: "Description used as the unique identifier for this 1:1 NAT rule.",
		},
		"disabled": {
			Description: "Disable this 1:1 NAT rule.",
		},
		"no_binat": {
			Description: "Disable 1:1 NAT (BINAT). When enabled, the rule acts as a negation — matching traffic is NOT translated.",
		},
		"nat_reflection": {
			Description:         fmt.Sprintf("NAT reflection mode. Options: %s. Empty string means system default.", wrapElementsJoin(pfsense.NATOneToOne{}.NATReflectionModes(), "'")),
			MarkdownDescription: fmt.Sprintf("NAT reflection mode. Options: %s. Empty string means system default.", wrapElementsJoin(pfsense.NATOneToOne{}.NATReflectionModes(), "`")),
		},
	}
}

// AttrTypes returns the attribute types for the 1:1 NAT model
func (FirewallNAT1to1Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"external":            types.StringType,
		"interface":           types.StringType,
		"ipprotocol":          types.StringType,
		"source_address":      types.StringType,
		"source_not":          types.BoolType,
		"destination_address": types.StringType,
		"destination_not":     types.BoolType,
		"description":         types.StringType,
		"disabled":            types.BoolType,
		"no_binat":            types.BoolType,
		"nat_reflection":      types.StringType,
	}
}

// toClient converts the Terraform model to a pfsense.NATOneToOne struct
func (m *FirewallNAT1to1Model) toClient(ctx context.Context, diags *diag.Diagnostics) *pfsense.NATOneToOne {
	rule := &pfsense.NATOneToOne{}

	if !m.External.IsNull() && !m.External.IsUnknown() {
		if err := rule.SetExternal(m.External.ValueString()); err != nil {
			diags.AddAttributeError(
				path.Root("external"),
				"Invalid External Address",
				fmt.Sprintf("Unable to set external address: %s", err.Error()),
			)
		}
	}

	if !m.Interface.IsNull() && !m.Interface.IsUnknown() {
		if err := rule.SetInterface(m.Interface.ValueString()); err != nil {
			diags.AddAttributeError(
				path.Root("interface"),
				"Invalid Interface",
				fmt.Sprintf("Unable to set interface: %s", err.Error()),
			)
		}
	}

	if !m.IPProtocol.IsNull() && !m.IPProtocol.IsUnknown() {
		if err := rule.SetIPProtocol(m.IPProtocol.ValueString()); err != nil {
			diags.AddAttributeError(
				path.Root("ipprotocol"),
				"Invalid IP Protocol",
				fmt.Sprintf("Unable to set IP protocol: %s", err.Error()),
			)
		}
	}

	if !m.SourceAddress.IsNull() && !m.SourceAddress.IsUnknown() {
		if err := rule.SetSourceAddress(m.SourceAddress.ValueString()); err != nil {
			diags.AddAttributeError(
				path.Root("source_address"),
				"Invalid Source Address",
				fmt.Sprintf("Unable to set source address: %s", err.Error()),
			)
		}
	}

	if !m.SourceNot.IsNull() && !m.SourceNot.IsUnknown() {
		if err := rule.SetSourceNot(m.SourceNot.ValueBool()); err != nil {
			diags.AddAttributeError(
				path.Root("source_not"),
				"Invalid Source Not",
				fmt.Sprintf("Unable to set source not: %s", err.Error()),
			)
		}
	}

	if !m.DestAddress.IsNull() && !m.DestAddress.IsUnknown() {
		if err := rule.SetDestAddress(m.DestAddress.ValueString()); err != nil {
			diags.AddAttributeError(
				path.Root("destination_address"),
				"Invalid Destination Address",
				fmt.Sprintf("Unable to set destination address: %s", err.Error()),
			)
		}
	}

	if !m.DestNot.IsNull() && !m.DestNot.IsUnknown() {
		if err := rule.SetDestNot(m.DestNot.ValueBool()); err != nil {
			diags.AddAttributeError(
				path.Root("destination_not"),
				"Invalid Destination Not",
				fmt.Sprintf("Unable to set destination not: %s", err.Error()),
			)
		}
	}

	if !m.Description.IsNull() && !m.Description.IsUnknown() {
		if err := rule.SetDescription(m.Description.ValueString()); err != nil {
			diags.AddAttributeError(
				path.Root("description"),
				"Invalid Description",
				fmt.Sprintf("Unable to set description: %s", err.Error()),
			)
		}
	}

	if !m.Disabled.IsNull() && !m.Disabled.IsUnknown() {
		if err := rule.SetDisabled(m.Disabled.ValueBool()); err != nil {
			diags.AddAttributeError(
				path.Root("disabled"),
				"Invalid Disabled",
				fmt.Sprintf("Unable to set disabled: %s", err.Error()),
			)
		}
	}

	if !m.NoBinat.IsNull() && !m.NoBinat.IsUnknown() {
		if err := rule.SetNoBinat(m.NoBinat.ValueBool()); err != nil {
			diags.AddAttributeError(
				path.Root("no_binat"),
				"Invalid No BINAT",
				fmt.Sprintf("Unable to set no_binat: %s", err.Error()),
			)
		}
	}

	if !m.NATReflection.IsNull() && !m.NATReflection.IsUnknown() {
		if err := rule.SetNATReflection(m.NATReflection.ValueString()); err != nil {
			diags.AddAttributeError(
				path.Root("nat_reflection"),
				"Invalid NAT Reflection",
				fmt.Sprintf("Unable to set NAT reflection: %s", err.Error()),
			)
		}
	}

	return rule
}

// fromClient converts a pfsense.NATOneToOne struct to the Terraform model
func (m *FirewallNAT1to1Model) fromClient(ctx context.Context, rule *pfsense.NATOneToOne, diags *diag.Diagnostics) {
	m.External = types.StringValue(rule.External)
	m.Interface = types.StringValue(rule.Interface)
	m.IPProtocol = types.StringValue(rule.IPProtocol)
	m.SourceAddress = types.StringValue(rule.SourceAddress)
	m.SourceNot = types.BoolValue(rule.SourceNot)
	m.DestAddress = types.StringValue(rule.DestAddress)
	m.DestNot = types.BoolValue(rule.DestNot)
	m.Description = types.StringValue(rule.Description)
	m.Disabled = types.BoolValue(rule.Disabled)
	m.NoBinat = types.BoolValue(rule.NoBinat)
	m.NATReflection = types.StringValue(rule.NATReflection)
}
