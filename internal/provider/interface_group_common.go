package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type InterfaceGroupModel struct {
	Name        types.String `tfsdk:"name"`
	Members     types.List   `tfsdk:"members"`
	Description types.String `tfsdk:"description"`
}

func (InterfaceGroupModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"name": {
			Description: "Name of the interface group. Must be unique, maximum 15 characters.",
		},
		"members": {
			Description: "List of logical interface names (e.g. 'lan', 'opt1') that are members of this group.",
		},
		"description": {
			Description: descriptionDescription,
		},
	}
}

func (InterfaceGroupModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":        types.StringType,
		"members":     types.ListType{ElemType: types.StringType},
		"description": types.StringType,
	}
}

func (m *InterfaceGroupModel) Set(_ context.Context, group pfsense.InterfaceGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Name = types.StringValue(group.Name)

	if group.Description != "" {
		m.Description = types.StringValue(group.Description)
	} else {
		m.Description = types.StringNull()
	}

	memberValues := make([]attr.Value, 0, len(group.Members))
	for _, member := range group.Members {
		memberValues = append(memberValues, types.StringValue(member))
	}

	membersValue, newDiags := types.ListValue(types.StringType, memberValues)
	diags.Append(newDiags...)
	m.Members = membersValue

	return diags
}

func (m InterfaceGroupModel) Value(ctx context.Context, group *pfsense.InterfaceGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("name"),
		"Name cannot be parsed",
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

	return diags
}
