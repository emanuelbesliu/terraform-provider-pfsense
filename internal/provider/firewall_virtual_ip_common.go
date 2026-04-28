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

type FirewallVirtualIPModel struct {
	Mode        types.String `tfsdk:"mode"`
	Interface   types.String `tfsdk:"interface"`
	VHID        types.Int64  `tfsdk:"vhid"`
	AdvSkew     types.Int64  `tfsdk:"advskew"`
	AdvBase     types.Int64  `tfsdk:"advbase"`
	Password    types.String `tfsdk:"password"`
	Subnet      types.String `tfsdk:"subnet"`
	SubnetBits  types.Int64  `tfsdk:"subnet_bits"`
	Description types.String `tfsdk:"description"`
	UniqueID    types.String `tfsdk:"unique_id"`
}

func (FirewallVirtualIPModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"mode": {
			Description: fmt.Sprintf("Virtual IP mode. Must be one of: %s, %s, %s, %s.", pfsense.VirtualIPModeIPAlias, pfsense.VirtualIPModeCarp, pfsense.VirtualIPModeProxyARP, pfsense.VirtualIPModeOther),
		},
		"interface": {
			Description: "Network interface for the virtual IP (e.g. 'wan', 'lan', 'opt1').",
		},
		"vhid": {
			Description: fmt.Sprintf("VHID group number (CARP only), between %d and %d.", pfsense.MinVHID, pfsense.MaxVHID),
		},
		"advskew": {
			Description: fmt.Sprintf("Advertisement skew (CARP only), between %d and %d.", pfsense.MinAdvSkew, pfsense.MaxAdvSkew),
		},
		"advbase": {
			Description: fmt.Sprintf("Advertisement frequency in seconds (CARP only), between %d and %d.", pfsense.MinAdvBase, pfsense.MaxAdvBase),
		},
		"password": {
			Description: "VHID group password (CARP only).",
		},
		"subnet": {
			Description: "IP address for the virtual IP (e.g. '10.0.1.100').",
		},
		"subnet_bits": {
			Description: "CIDR prefix length (e.g. 32 for a single host, 24 for a /24 subnet).",
		},
		"description": {
			Description: descriptionDescription,
		},
		"unique_id": {
			Description: "Unique identifier assigned by pfSense.",
		},
	}
}

func (FirewallVirtualIPModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"mode":        types.StringType,
		"interface":   types.StringType,
		"vhid":        types.Int64Type,
		"advskew":     types.Int64Type,
		"advbase":     types.Int64Type,
		"password":    types.StringType,
		"subnet":      types.StringType,
		"subnet_bits": types.Int64Type,
		"description": types.StringType,
		"unique_id":   types.StringType,
	}
}

func (m *FirewallVirtualIPModel) Set(_ context.Context, vip pfsense.VirtualIP) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Mode = types.StringValue(vip.Mode)
	m.Interface = types.StringValue(vip.Interface)

	if vip.VHID != nil {
		m.VHID = types.Int64Value(int64(*vip.VHID))
	} else {
		m.VHID = types.Int64Null()
	}

	if vip.AdvSkew != nil {
		m.AdvSkew = types.Int64Value(int64(*vip.AdvSkew))
	} else {
		m.AdvSkew = types.Int64Null()
	}

	if vip.AdvBase != nil {
		m.AdvBase = types.Int64Value(int64(*vip.AdvBase))
	} else {
		m.AdvBase = types.Int64Null()
	}

	if vip.Password != "" {
		m.Password = types.StringValue(vip.Password)
	} else {
		m.Password = types.StringNull()
	}

	m.Subnet = types.StringValue(vip.Subnet)
	m.SubnetBits = types.Int64Value(int64(vip.SubnetBits))

	if vip.Description != "" {
		m.Description = types.StringValue(vip.Description)
	} else {
		m.Description = types.StringNull()
	}

	m.UniqueID = types.StringValue(vip.UniqueID)

	return diags
}

func (m FirewallVirtualIPModel) Value(_ context.Context, vip *pfsense.VirtualIP) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("mode"),
		"Mode cannot be parsed",
		vip.SetMode(m.Mode.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("interface"),
		"Interface cannot be parsed",
		vip.SetInterface(m.Interface.ValueString()),
	)

	if !m.VHID.IsNull() {
		vhid := int(m.VHID.ValueInt64())
		addPathError(
			&diags,
			path.Root("vhid"),
			"VHID cannot be parsed",
			vip.SetVHID(&vhid),
		)
	}

	if !m.AdvSkew.IsNull() {
		advskew := int(m.AdvSkew.ValueInt64())
		addPathError(
			&diags,
			path.Root("advskew"),
			"Advertisement skew cannot be parsed",
			vip.SetAdvSkew(&advskew),
		)
	}

	if !m.AdvBase.IsNull() {
		advbase := int(m.AdvBase.ValueInt64())
		addPathError(
			&diags,
			path.Root("advbase"),
			"Advertisement base cannot be parsed",
			vip.SetAdvBase(&advbase),
		)
	}

	if !m.Password.IsNull() {
		addPathError(
			&diags,
			path.Root("password"),
			"Password cannot be parsed",
			vip.SetPassword(m.Password.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("subnet"),
		"Subnet cannot be parsed",
		vip.SetSubnet(m.Subnet.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("subnet_bits"),
		"Subnet bits cannot be parsed",
		vip.SetSubnetBits(int(m.SubnetBits.ValueInt64())),
	)

	if !m.Description.IsNull() {
		addPathError(
			&diags,
			path.Root("description"),
			"Description cannot be parsed",
			vip.SetDescription(m.Description.ValueString()),
		)
	}

	return diags
}
