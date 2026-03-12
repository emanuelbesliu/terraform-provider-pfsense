package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type RouteModel struct {
	Network     types.String `tfsdk:"network"`
	Gateway     types.String `tfsdk:"gateway"`
	Description types.String `tfsdk:"description"`
	Disabled    types.Bool   `tfsdk:"disabled"`
}

func (RouteModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"network": {
			Description: "Destination network for this static route, in CIDR format.",
		},
		"gateway": {
			Description: "Gateway to use for this static route. Must be the name of a configured gateway.",
		},
		"description": {
			Description: descriptionDescription,
		},
		"disabled": {
			Description: "Disable this static route.",
		},
	}
}

func (RouteModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"network":     types.StringType,
		"gateway":     types.StringType,
		"description": types.StringType,
		"disabled":    types.BoolType,
	}
}

func (m *RouteModel) Set(_ context.Context, route pfsense.Route) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Network = types.StringValue(route.Network)
	m.Gateway = types.StringValue(route.Gateway)

	if route.Description != "" {
		m.Description = types.StringValue(route.Description)
	}

	m.Disabled = types.BoolValue(route.Disabled)

	return diags
}

func (m RouteModel) Value(_ context.Context, route *pfsense.Route) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("network"),
		"Network cannot be parsed",
		route.SetNetwork(m.Network.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("gateway"),
		"Gateway cannot be parsed",
		route.SetGateway(m.Gateway.ValueString()),
	)

	if !m.Description.IsNull() {
		addPathError(
			&diags,
			path.Root("description"),
			"Description cannot be parsed",
			route.SetDescription(m.Description.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("disabled"),
		"Disabled cannot be parsed",
		route.SetDisabled(m.Disabled.ValueBool()),
	)

	return diags
}
