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

type GatewayGroupMemberModel struct {
	Gateway   types.String `tfsdk:"gateway"`
	Tier      types.Int64  `tfsdk:"tier"`
	VirtualIP types.String `tfsdk:"virtual_ip"`
}

func (GatewayGroupMemberModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"gateway": {
			Description: "Name of the gateway to include in this group.",
		},
		"tier": {
			Description: fmt.Sprintf("Priority tier for this gateway (%d-%d). Gateways in the same tier are load balanced; lower tiers are preferred for failover.", pfsense.MinGatewayGroupTier, pfsense.MaxGatewayGroupTier),
		},
		"virtual_ip": {
			Description: "Virtual IP address to use for this gateway member. Leave empty to use the gateway's interface address.",
		},
	}
}

func (GatewayGroupMemberModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"gateway":    types.StringType,
		"tier":       types.Int64Type,
		"virtual_ip": types.StringType,
	}
}

func (m *GatewayGroupMemberModel) Set(_ context.Context, member pfsense.GatewayGroupMember) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Gateway = types.StringValue(member.Gateway)
	m.Tier = types.Int64Value(int64(member.Tier))

	if member.VirtualIP != "" {
		m.VirtualIP = types.StringValue(member.VirtualIP)
	} else {
		m.VirtualIP = types.StringNull()
	}

	return diags
}

func (m GatewayGroupMemberModel) Value(_ context.Context, member *pfsense.GatewayGroupMember, attrPath path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	member.Gateway = m.Gateway.ValueString()
	member.Tier = int(m.Tier.ValueInt64())

	if !m.VirtualIP.IsNull() {
		member.VirtualIP = m.VirtualIP.ValueString()
	}

	return diags
}

type GatewayGroupModel struct {
	Name               types.String `tfsdk:"name"`
	Description        types.String `tfsdk:"description"`
	Trigger            types.String `tfsdk:"trigger"`
	KeepFailoverStates types.String `tfsdk:"keep_failover_states"`
	Members            types.List   `tfsdk:"members"`
}

func (GatewayGroupModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"name": {
			Description: "Name of the gateway group. Must be unique and cannot conflict with gateway names.",
		},
		"description": {
			Description: descriptionDescription,
		},
		"trigger": {
			Description:         fmt.Sprintf("Trigger level that determines when a member gateway is excluded from the group. Options: %s.", wrapElementsJoin(pfsense.GatewayGroup{}.Triggers(), "'")),
			MarkdownDescription: fmt.Sprintf("Trigger level that determines when a member gateway is excluded from the group. Options: %s.", wrapElementsJoin(pfsense.GatewayGroup{}.Triggers(), "`")),
		},
		"keep_failover_states": {
			Description: "Control whether firewall states are kept or killed when a gateway recovers from failure. Empty string uses global behavior.",
		},
		"members": {
			Description: "List of gateway members with their priority tiers.",
		},
	}
}

func (GatewayGroupModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                 types.StringType,
		"description":          types.StringType,
		"trigger":              types.StringType,
		"keep_failover_states": types.StringType,
		"members":              types.ListType{ElemType: types.ObjectType{AttrTypes: GatewayGroupMemberModel{}.AttrTypes()}},
	}
}

func (m *GatewayGroupModel) Set(ctx context.Context, group pfsense.GatewayGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Name = types.StringValue(group.Name)

	if group.Description != "" {
		m.Description = types.StringValue(group.Description)
	} else {
		m.Description = types.StringNull()
	}

	m.Trigger = types.StringValue(group.Trigger)

	if group.KeepFailoverStates != "" {
		m.KeepFailoverStates = types.StringValue(group.KeepFailoverStates)
	} else {
		m.KeepFailoverStates = types.StringNull()
	}

	memberModels := []GatewayGroupMemberModel{}
	for _, member := range group.Members {
		var memberModel GatewayGroupMemberModel
		diags.Append(memberModel.Set(ctx, member)...)
		memberModels = append(memberModels, memberModel)
	}

	membersValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: GatewayGroupMemberModel{}.AttrTypes()}, memberModels)
	diags.Append(newDiags...)
	m.Members = membersValue

	return diags
}

func (m GatewayGroupModel) Value(ctx context.Context, group *pfsense.GatewayGroup) diag.Diagnostics {
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

	addPathError(
		&diags,
		path.Root("trigger"),
		"Trigger cannot be parsed",
		group.SetTrigger(m.Trigger.ValueString()),
	)

	if !m.KeepFailoverStates.IsNull() {
		addPathError(
			&diags,
			path.Root("keep_failover_states"),
			"Keep failover states cannot be parsed",
			group.SetKeepFailoverStates(m.KeepFailoverStates.ValueString()),
		)
	}

	var memberModels []GatewayGroupMemberModel
	if !m.Members.IsNull() {
		diags.Append(m.Members.ElementsAs(ctx, &memberModels, false)...)
	}

	members := make([]pfsense.GatewayGroupMember, 0, len(memberModels))
	for index, memberModel := range memberModels {
		var member pfsense.GatewayGroupMember
		diags.Append(memberModel.Value(ctx, &member, path.Root("members").AtListIndex(index))...)
		members = append(members, member)
	}

	addPathError(
		&diags,
		path.Root("members"),
		"Members cannot be parsed",
		group.SetMembers(members),
	)

	return diags
}
