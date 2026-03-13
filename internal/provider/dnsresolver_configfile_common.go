package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type DNSResolverConfigFilesModel struct {
	All types.List `tfsdk:"all"`
}

type DNSResolverConfigFileModel struct {
	Name    types.String `tfsdk:"name"`
	Content types.String `tfsdk:"content"`
}

func (DNSResolverConfigFileModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"name": {
			Description: "Name of config file.",
		},
		"content": {
			Description:         "Contents of file. Must specify Unbound clause(s). Comments start with '#' and last to the end of line.",
			MarkdownDescription: "Contents of file. Must specify Unbound clause(s). Comments start with `#` and last to the end of line.",
		},
	}
}

func (DNSResolverConfigFileModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":    types.StringType,
		"content": types.StringType,
	}
}

func (m *DNSResolverConfigFilesModel) Set(ctx context.Context, configFiles pfsense.ConfigFiles) diag.Diagnostics {
	var diags diag.Diagnostics

	configFileModels := []DNSResolverConfigFileModel{}
	for _, configFile := range configFiles {
		var configFileModel DNSResolverConfigFileModel
		diags.Append(configFileModel.Set(ctx, configFile)...)
		configFileModels = append(configFileModels, configFileModel)
	}

	configFilesValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: DNSResolverConfigFileModel{}.AttrTypes()}, configFileModels)
	diags.Append(newDiags...)
	m.All = configFilesValue

	return diags
}

func (r *DNSResolverConfigFileModel) Set(_ context.Context, configFile pfsense.ConfigFile) diag.Diagnostics {
	var diags diag.Diagnostics

	r.Name = types.StringValue(configFile.Name)
	r.Content = types.StringValue(configFile.Content)

	return diags
}

func (r DNSResolverConfigFileModel) Value(_ context.Context, configFile *pfsense.ConfigFile) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("name"),
		"Name cannot be parsed",
		configFile.SetName(r.Name.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("content"),
		"Content cannot be parsed",
		configFile.SetContent(r.Content.ValueString()),
	)

	return diags
}
