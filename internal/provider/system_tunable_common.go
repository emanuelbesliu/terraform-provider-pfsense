package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type SystemTunableModel struct {
	Tunable      types.String `tfsdk:"tunable"`
	TunableValue types.String `tfsdk:"value"`
	Description  types.String `tfsdk:"description"`
}

func (SystemTunableModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"tunable": {
			Description: "Tunable name (sysctl MIB name, e.g. 'net.inet.tcp.syncookies').",
		},
		"value": {
			Description: "Value to set for the tunable. Use 'default' to reset to the system default.",
		},
		"description": {
			Description: descriptionDescription,
		},
	}
}

func (SystemTunableModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"tunable":     types.StringType,
		"value":       types.StringType,
		"description": types.StringType,
	}
}

func (m *SystemTunableModel) Set(_ context.Context, tunable pfsense.SystemTunable) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Tunable = types.StringValue(tunable.Tunable)
	m.TunableValue = types.StringValue(tunable.Value)

	if tunable.Description != "" {
		m.Description = types.StringValue(tunable.Description)
	} else {
		m.Description = types.StringNull()
	}

	return diags
}

func (m SystemTunableModel) Value(_ context.Context, tunable *pfsense.SystemTunable) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("tunable"),
		"Tunable name cannot be parsed",
		tunable.SetTunable(m.Tunable.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("value"),
		"Tunable value cannot be parsed",
		tunable.SetValue(m.TunableValue.ValueString()),
	)

	if !m.Description.IsNull() {
		addPathError(
			&diags,
			path.Root("description"),
			"Description cannot be parsed",
			tunable.SetDescription(m.Description.ValueString()),
		)
	}

	return diags
}
