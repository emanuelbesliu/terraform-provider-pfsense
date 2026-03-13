package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type UserModel struct {
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Scope          types.String `tfsdk:"scope"`
	UID            types.String `tfsdk:"uid"`
	Disabled       types.Bool   `tfsdk:"disabled"`
	Expires        types.String `tfsdk:"expires"`
	AuthorizedKeys types.String `tfsdk:"authorized_keys"`
	IPSecPSK       types.String `tfsdk:"ipsec_psk"`
	Privileges     types.List   `tfsdk:"privileges"`
	Groups         types.List   `tfsdk:"groups"`
	CustomSettings types.Bool   `tfsdk:"custom_settings"`
	WebGUICss      types.String `tfsdk:"webgui_css"`
	DashboardCols  types.String `tfsdk:"dashboard_columns"`
	KeepHistory    types.Bool   `tfsdk:"keep_history"`
}

func (UserModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"name": {
			Description: "Username. Must contain only alphanumeric characters, dots, hyphens, and underscores (max 32 characters).",
		},
		"description": {
			Description: "Full name or description of the user.",
		},
		"scope": {
			Description: "User scope ('system' for built-in users, 'user' for user-created).",
		},
		"uid": {
			Description: "Numeric user ID (auto-assigned).",
		},
		"disabled": {
			Description: "Whether the user account is disabled.",
		},
		"expires": {
			Description: "Account expiration date in MM/DD/YYYY format. Empty string or omitted means no expiration.",
		},
		"authorized_keys": {
			Description: "SSH authorized keys (raw public key text, base64-encoded for storage in pfSense config).",
		},
		"ipsec_psk": {
			Description: "IPsec pre-shared key. Must contain only ASCII characters.",
		},
		"privileges": {
			Description: "List of privilege names assigned to the user (e.g. 'page-all' for full admin access).",
		},
		"groups": {
			Description: "List of group names the user belongs to (excluding the implicit 'all' group).",
		},
		"custom_settings": {
			Description: "Enable per-user GUI settings (CSS theme, dashboard columns, command history).",
		},
		"webgui_css": {
			Description: "Web GUI CSS theme file (e.g. 'pfSense.css'). Only effective when custom_settings is enabled.",
		},
		"dashboard_columns": {
			Description: "Number of dashboard columns (default '2'). Only effective when custom_settings is enabled.",
		},
		"keep_history": {
			Description: "Keep command history for this user. Only effective when custom_settings is enabled.",
		},
		"password": {
			Description: "User password. Required when creating a new user. Write-only (cannot be read back from pfSense, only the bcrypt hash is stored).",
		},
	}
}

func (UserModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":              types.StringType,
		"description":       types.StringType,
		"scope":             types.StringType,
		"uid":               types.StringType,
		"disabled":          types.BoolType,
		"expires":           types.StringType,
		"authorized_keys":   types.StringType,
		"ipsec_psk":         types.StringType,
		"privileges":        types.ListType{ElemType: types.StringType},
		"groups":            types.ListType{ElemType: types.StringType},
		"custom_settings":   types.BoolType,
		"webgui_css":        types.StringType,
		"dashboard_columns": types.StringType,
		"keep_history":      types.BoolType,
	}
}

func (m *UserModel) Set(_ context.Context, user pfsense.User) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Name = types.StringValue(user.Name)
	m.Scope = types.StringValue(user.Scope)
	m.UID = types.StringValue(user.UID)
	m.Disabled = types.BoolValue(user.Disabled)

	if user.Description != "" {
		m.Description = types.StringValue(user.Description)
	} else {
		m.Description = types.StringNull()
	}

	if user.Expires != "" {
		m.Expires = types.StringValue(user.Expires)
	} else {
		m.Expires = types.StringNull()
	}

	if user.AuthorizedKeys != "" {
		m.AuthorizedKeys = types.StringValue(user.AuthorizedKeys)
	} else {
		m.AuthorizedKeys = types.StringNull()
	}

	if user.IPSecPSK != "" {
		m.IPSecPSK = types.StringValue(user.IPSecPSK)
	} else {
		m.IPSecPSK = types.StringNull()
	}

	m.CustomSettings = types.BoolValue(user.CustomSettings)
	m.KeepHistory = types.BoolValue(user.KeepHistory)

	if user.WebGUICss != "" {
		m.WebGUICss = types.StringValue(user.WebGUICss)
	} else {
		m.WebGUICss = types.StringNull()
	}

	if user.DashboardCols != "" {
		m.DashboardCols = types.StringValue(user.DashboardCols)
	} else {
		m.DashboardCols = types.StringNull()
	}

	// Privileges - filter out empty strings.
	if len(user.Privileges) > 0 {
		privValues := make([]attr.Value, 0, len(user.Privileges))
		for _, p := range user.Privileges {
			privValues = append(privValues, types.StringValue(p))
		}

		privList, newDiags := types.ListValue(types.StringType, privValues)
		diags.Append(newDiags...)
		m.Privileges = privList
	} else {
		m.Privileges = types.ListNull(types.StringType)
	}

	// Groups - exclude the implicit 'all' group.
	filteredGroups := make([]string, 0, len(user.Groups))
	for _, g := range user.Groups {
		if g != "all" {
			filteredGroups = append(filteredGroups, g)
		}
	}

	if len(filteredGroups) > 0 {
		groupValues := make([]attr.Value, 0, len(filteredGroups))
		for _, g := range filteredGroups {
			groupValues = append(groupValues, types.StringValue(g))
		}

		groupList, newDiags := types.ListValue(types.StringType, groupValues)
		diags.Append(newDiags...)
		m.Groups = groupList
	} else {
		m.Groups = types.ListNull(types.StringType)
	}

	return diags
}

func (m UserModel) Value(ctx context.Context, user *pfsense.User) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("name"),
		"Username cannot be parsed",
		user.SetName(m.Name.ValueString()),
	)

	if !m.Description.IsNull() {
		addPathError(
			&diags,
			path.Root("description"),
			"Description cannot be parsed",
			user.SetDescription(m.Description.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("disabled"),
		"Disabled cannot be parsed",
		user.SetDisabled(m.Disabled.ValueBool()),
	)

	if !m.Expires.IsNull() {
		addPathError(
			&diags,
			path.Root("expires"),
			"Expires cannot be parsed",
			user.SetExpires(m.Expires.ValueString()),
		)
	}

	if !m.AuthorizedKeys.IsNull() {
		addPathError(
			&diags,
			path.Root("authorized_keys"),
			"Authorized keys cannot be parsed",
			user.SetAuthorizedKeys(m.AuthorizedKeys.ValueString()),
		)
	}

	if !m.IPSecPSK.IsNull() {
		addPathError(
			&diags,
			path.Root("ipsec_psk"),
			"IPsec PSK cannot be parsed",
			user.SetIPSecPSK(m.IPSecPSK.ValueString()),
		)
	}

	// Privileges.
	var privStrings []string
	if !m.Privileges.IsNull() {
		diags.Append(m.Privileges.ElementsAs(ctx, &privStrings, false)...)
	}

	addPathError(
		&diags,
		path.Root("privileges"),
		"Privileges cannot be parsed",
		user.SetPrivileges(privStrings),
	)

	// Groups.
	var groupStrings []string
	if !m.Groups.IsNull() {
		diags.Append(m.Groups.ElementsAs(ctx, &groupStrings, false)...)
	}

	addPathError(
		&diags,
		path.Root("groups"),
		"Groups cannot be parsed",
		user.SetGroups(groupStrings),
	)

	addPathError(
		&diags,
		path.Root("custom_settings"),
		"Custom settings cannot be parsed",
		user.SetCustomSettings(m.CustomSettings.ValueBool()),
	)

	if !m.WebGUICss.IsNull() {
		addPathError(
			&diags,
			path.Root("webgui_css"),
			"WebGUI CSS cannot be parsed",
			user.SetWebGUICss(m.WebGUICss.ValueString()),
		)
	}

	if !m.DashboardCols.IsNull() {
		addPathError(
			&diags,
			path.Root("dashboard_columns"),
			"Dashboard columns cannot be parsed",
			user.SetDashboardCols(m.DashboardCols.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("keep_history"),
		"Keep history cannot be parsed",
		user.SetKeepHistory(m.KeepHistory.ValueBool()),
	)

	return diags
}
