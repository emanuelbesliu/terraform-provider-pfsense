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

type VLANModel struct {
	ParentInterface types.String `tfsdk:"parent_interface"`
	Tag             types.Int64  `tfsdk:"tag"`
	PCP             types.Int64  `tfsdk:"pcp"`
	Description     types.String `tfsdk:"description"`
	VLANInterface   types.String `tfsdk:"vlan_interface"`
}

func (VLANModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"parent_interface": {
			Description: "Parent physical interface for the VLAN (e.g. 'vmx0', 'igb1').",
		},
		"tag": {
			Description: fmt.Sprintf("VLAN tag (ID), must be between %d and %d.", pfsense.MinVLANTag, pfsense.MaxVLANTag),
		},
		"pcp": {
			Description: fmt.Sprintf("VLAN Priority Code Point (802.1p), between %d and %d.", pfsense.MinVLANPCP, pfsense.MaxVLANPCP),
		},
		"description": {
			Description: descriptionDescription,
		},
		"vlan_interface": {
			Description: "Computed VLAN interface name (e.g. 'vmx0.100'). This is the OS-level interface created for the VLAN.",
		},
	}
}

func (VLANModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"parent_interface": types.StringType,
		"tag":              types.Int64Type,
		"pcp":              types.Int64Type,
		"description":      types.StringType,
		"vlan_interface":   types.StringType,
	}
}

func (m *VLANModel) Set(_ context.Context, vlan pfsense.VLAN) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ParentInterface = types.StringValue(vlan.ParentInterface)
	m.Tag = types.Int64Value(int64(vlan.Tag))

	if vlan.PCP != nil {
		m.PCP = types.Int64Value(int64(*vlan.PCP))
	} else {
		m.PCP = types.Int64Null()
	}

	if vlan.Description != "" {
		m.Description = types.StringValue(vlan.Description)
	} else {
		m.Description = types.StringNull()
	}

	m.VLANInterface = types.StringValue(vlan.VLANInterface)

	return diags
}

func (m VLANModel) Value(_ context.Context, vlan *pfsense.VLAN) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("parent_interface"),
		"Parent interface cannot be parsed",
		vlan.SetParentInterface(m.ParentInterface.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("tag"),
		"Tag cannot be parsed",
		vlan.SetTag(int(m.Tag.ValueInt64())),
	)

	if !m.PCP.IsNull() {
		pcp := int(m.PCP.ValueInt64())
		addPathError(
			&diags,
			path.Root("pcp"),
			"PCP cannot be parsed",
			vlan.SetPCP(&pcp),
		)
	}

	if !m.Description.IsNull() {
		addPathError(
			&diags,
			path.Root("description"),
			"Description cannot be parsed",
			vlan.SetDescription(m.Description.ValueString()),
		)
	}

	return diags
}
