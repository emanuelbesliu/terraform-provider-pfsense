package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type FirewallNATNPtRuleModel struct {
	Interface         types.String `tfsdk:"interface"`
	SourcePrefix      types.String `tfsdk:"source_prefix"`
	SourceNot         types.Bool   `tfsdk:"source_not"`
	DestinationPrefix types.String `tfsdk:"destination_prefix"`
	DestinationNot    types.Bool   `tfsdk:"destination_not"`
	Disabled          types.Bool   `tfsdk:"disabled"`
	Description       types.String `tfsdk:"description"`
}

func (FirewallNATNPtRuleModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"interface": {
			Description: "Network interface this NPt rule applies to (e.g. 'wan', 'lan', 'opt1'). Typically the WAN.",
		},
		"source_prefix": {
			Description: "Internal (LAN) ULA IPv6 prefix in CIDR notation (e.g. 'fd00:1::/64'). The prefix size is applied to the external prefix.",
		},
		"source_not": {
			Description: "Invert the source prefix match.",
		},
		"destination_prefix": {
			Description: "External global unicast routable IPv6 prefix in CIDR notation (e.g. '2001:db8::/64'). The prefix size must equal the source prefix size. May also be an interface name to use a delegated (track6) prefix.",
		},
		"destination_not": {
			Description: "Invert the destination prefix match.",
		},
		"disabled": {
			Description: "Disable this NPt rule.",
		},
		"description": {
			Description: "Description used as the unique identifier for this NAT NPt rule.",
		},
	}
}

func (FirewallNATNPtRuleModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interface":          types.StringType,
		"source_prefix":      types.StringType,
		"source_not":         types.BoolType,
		"destination_prefix": types.StringType,
		"destination_not":    types.BoolType,
		"disabled":           types.BoolType,
		"description":        types.StringType,
	}
}

func (m *FirewallNATNPtRuleModel) Set(_ context.Context, r pfsense.NATNPt) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Interface = types.StringValue(r.Interface)
	m.SourcePrefix = types.StringValue(r.SourcePrefix)
	m.SourceNot = types.BoolValue(r.SourceNot)
	m.DestinationPrefix = types.StringValue(r.DestinationPrefix)
	m.DestinationNot = types.BoolValue(r.DestinationNot)
	m.Disabled = types.BoolValue(r.Disabled)
	m.Description = types.StringValue(r.Description)

	return diags
}

func (m FirewallNATNPtRuleModel) Value(_ context.Context, r *pfsense.NATNPt) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(&diags, path.Root("interface"), "Interface cannot be parsed", r.SetInterface(m.Interface.ValueString()))
	addPathError(&diags, path.Root("source_prefix"), "Source prefix cannot be parsed", r.SetSourcePrefix(m.SourcePrefix.ValueString()))
	addPathError(&diags, path.Root("source_not"), "Source not cannot be parsed", r.SetSourceNot(m.SourceNot.ValueBool()))
	addPathError(&diags, path.Root("destination_prefix"), "Destination prefix cannot be parsed", r.SetDestinationPrefix(m.DestinationPrefix.ValueString()))
	addPathError(&diags, path.Root("destination_not"), "Destination not cannot be parsed", r.SetDestinationNot(m.DestinationNot.ValueBool()))
	addPathError(&diags, path.Root("disabled"), "Disabled cannot be parsed", r.SetDisabled(m.Disabled.ValueBool()))
	addPathError(&diags, path.Root("description"), "Description cannot be parsed", r.SetDescription(m.Description.ValueString()))

	return diags
}
