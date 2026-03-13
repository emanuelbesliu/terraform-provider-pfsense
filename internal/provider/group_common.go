package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type GroupModel struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Scope       types.String `tfsdk:"scope"`
	GID         types.String `tfsdk:"gid"`
	Members     types.List   `tfsdk:"members"`
	Privileges  types.List   `tfsdk:"privileges"`
}

func (GroupModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"name": {
			Description: "Group name. Must contain only alphanumeric characters, dots, hyphens, and underscores (max 16 characters).",
		},
		"description": {
			Description: "Description of the group.",
		},
		"scope": {
			Description: "Group scope ('system' for built-in groups, 'local' for user-created).",
		},
		"gid": {
			Description: "Numeric group ID (auto-assigned).",
		},
		"members": {
			Description: "List of usernames that are members of this group.",
		},
		"privileges": {
			Description: "List of privilege names assigned to the group (e.g. 'page-all' for full admin access).",
		},
	}
}

func (GroupModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"description": types.StringType,
		"scope":       types.StringType,
		"gid":         types.StringType,
		"members":     types.ListType{ElemType: types.StringType},
		"privileges":  types.ListType{ElemType: types.StringType},
	}
}

func (m *GroupModel) Set(_ context.Context, group pfsense.Group) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Name = types.StringValue(group.Name)
	m.Scope = types.StringValue(group.Scope)
	m.GID = types.StringValue(group.GID)

	if group.Description != "" {
		m.Description = types.StringValue(group.Description)
	} else {
		m.Description = types.StringNull()
	}

	// Members (usernames).
	if len(group.Members) > 0 {
		memberValues := make([]attr.Value, 0, len(group.Members))
		for _, member := range group.Members {
			memberValues = append(memberValues, types.StringValue(member))
		}

		memberList, newDiags := types.ListValue(types.StringType, memberValues)
		diags.Append(newDiags...)
		m.Members = memberList
	} else {
		m.Members = types.ListNull(types.StringType)
	}

	// Privileges.
	if len(group.Privileges) > 0 {
		privValues := make([]attr.Value, 0, len(group.Privileges))
		for _, p := range group.Privileges {
			privValues = append(privValues, types.StringValue(p))
		}

		privList, newDiags := types.ListValue(types.StringType, privValues)
		diags.Append(newDiags...)
		m.Privileges = privList
	} else {
		m.Privileges = types.ListNull(types.StringType)
	}

	return diags
}

func (m GroupModel) Value(ctx context.Context, group *pfsense.Group) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("name"),
		"Group name cannot be parsed",
		group.SetName(m.Name.ValueString()),
	)

	if !m.Description.IsNull() {
		addPathError(
			&diags,
			path.Root("description"),
			"Description cannot be parsed",
			group.SetDescription(m.Description.ValueString()),
		)
	}

	// Members.
	var memberStrings []string
	if !m.Members.IsNull() {
		diags.Append(m.Members.ElementsAs(ctx, &memberStrings, false)...)
	}

	addPathError(
		&diags,
		path.Root("members"),
		"Members cannot be parsed",
		group.SetMembers(memberStrings),
	)

	// Privileges.
	var privStrings []string
	if !m.Privileges.IsNull() {
		diags.Append(m.Privileges.ElementsAs(ctx, &privStrings, false)...)
	}

	addPathError(
		&diags,
		path.Root("privileges"),
		"Privileges cannot be parsed",
		group.SetPrivileges(privStrings),
	)

	return diags
}
