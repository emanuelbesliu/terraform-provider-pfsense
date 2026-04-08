package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type PackageModel struct {
	Name             types.String `tfsdk:"name"`
	InstalledVersion types.String `tfsdk:"installed_version"`
	Description      types.String `tfsdk:"description"`
}

func (PackageModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"name": {
			Description: "Package name (e.g. 'pfSense-pkg-saml2-auth').",
		},
		"installed_version": {
			Description: "Currently installed version of the package.",
		},
		"description": {
			Description: "Package description.",
		},
	}
}

func (PackageModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":              types.StringType,
		"installed_version": types.StringType,
		"description":       types.StringType,
	}
}

func (m *PackageModel) Set(_ context.Context, pkg pfsense.Package) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Name = types.StringValue(pkg.Name)

	if pkg.InstalledVersion != "" {
		m.InstalledVersion = types.StringValue(pkg.InstalledVersion)
	} else {
		m.InstalledVersion = types.StringNull()
	}

	if pkg.Description != "" {
		m.Description = types.StringValue(pkg.Description)
	} else {
		m.Description = types.StringNull()
	}

	return diags
}

func (m PackageModel) Value(_ context.Context, pkg *pfsense.Package) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("name"),
		"Package name cannot be parsed",
		pkg.SetName(m.Name.ValueString()),
	)

	return diags
}
