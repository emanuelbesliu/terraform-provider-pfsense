package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type WakeOnLanModel struct {
	Interface   types.String `tfsdk:"interface"`
	MAC         types.String `tfsdk:"mac"`
	Description types.String `tfsdk:"description"`
}

func (WakeOnLanModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"interface": {
			Description: "Network interface for the Wake-on-LAN entry (e.g. 'lan', 'opt1').",
		},
		"mac": {
			Description: "MAC address of the device to wake.",
		},
		"description": {
			Description: "Description of the Wake-on-LAN entry.",
		},
	}
}

func (WakeOnLanModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"interface":   types.StringType,
		"mac":         types.StringType,
		"description": types.StringType,
	}
}

func (m *WakeOnLanModel) Set(_ context.Context, entry pfsense.WakeOnLanEntry) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Interface = types.StringValue(entry.Interface)
	m.MAC = types.StringValue(entry.MAC)
	m.Description = types.StringValue(entry.Description)

	return diags
}

func (m WakeOnLanModel) Value(_ context.Context, entry *pfsense.WakeOnLanEntry) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("interface"),
		"Interface cannot be parsed",
		entry.SetInterface(m.Interface.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("mac"),
		"MAC address cannot be parsed",
		entry.SetMAC(m.MAC.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("description"),
		"Description cannot be parsed",
		entry.SetDescription(m.Description.ValueString()),
	)

	return diags
}
