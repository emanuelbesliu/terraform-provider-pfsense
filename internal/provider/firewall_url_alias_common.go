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

type FirewallURLAliasModel struct {
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	Type            types.String `tfsdk:"type"`
	Entries         types.List   `tfsdk:"entries"`
	UpdateFrequency types.Int64  `tfsdk:"update_frequency"`
}

type FirewallURLAliasEntryModel struct {
	URL         types.String `tfsdk:"url"`
	Description types.String `tfsdk:"description"`
}

func (FirewallURLAliasModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"name": {
			Description: "Name of URL alias.",
		},
		"description": {
			Description: descriptionDescription,
		},
		"type": {
			Description:         fmt.Sprintf("Type of URL alias. Options: %s.", wrapElementsJoin(pfsense.FirewallURLAlias{}.Types(), "'")),
			MarkdownDescription: fmt.Sprintf("Type of URL alias. Options: %s.", wrapElementsJoin(pfsense.FirewallURLAlias{}.Types(), "`")),
		},
		"entries": {
			Description: "URL(s) to fetch and import.",
		},
		"update_frequency": {
			Description: "How often the URL table is updated, in days. Only applicable to 'urltable' and 'urltable_ports' types. Defaults to 7.",
		},
	}
}

func (FirewallURLAliasEntryModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"url": {
			Description: "URL to fetch. Must be a valid URL.",
		},
		"description": {
			Description: descriptionDescription,
		},
	}
}

func (FirewallURLAliasModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":             types.StringType,
		"description":      types.StringType,
		"type":             types.StringType,
		"entries":          types.ListType{ElemType: types.ObjectType{AttrTypes: FirewallURLAliasEntryModel{}.AttrTypes()}},
		"update_frequency": types.Int64Type,
	}
}

func (FirewallURLAliasEntryModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"url":         types.StringType,
		"description": types.StringType,
	}
}

func (m *FirewallURLAliasModel) Set(ctx context.Context, urlAlias pfsense.FirewallURLAlias) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Name = types.StringValue(urlAlias.Name)

	if urlAlias.Description != "" {
		m.Description = types.StringValue(urlAlias.Description)
	}

	m.Type = types.StringValue(urlAlias.Type)

	if urlAlias.Type == "urltable" || urlAlias.Type == "urltable_ports" {
		m.UpdateFrequency = types.Int64Value(int64(urlAlias.UpdateFrequency))
	} else {
		m.UpdateFrequency = types.Int64Null()
	}

	urlAliasEntryModels := []FirewallURLAliasEntryModel{}
	for _, entry := range urlAlias.Entries {
		var entryModel FirewallURLAliasEntryModel
		diags.Append(entryModel.Set(ctx, entry)...)
		urlAliasEntryModels = append(urlAliasEntryModels, entryModel)
	}

	urlAliasEntriesValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: FirewallURLAliasEntryModel{}.AttrTypes()}, urlAliasEntryModels)
	diags.Append(newDiags...)
	m.Entries = urlAliasEntriesValue

	return diags
}

func (m *FirewallURLAliasEntryModel) Set(_ context.Context, entry pfsense.FirewallURLAliasEntry) diag.Diagnostics {
	var diags diag.Diagnostics

	m.URL = types.StringValue(entry.URL)

	if entry.Description != "" {
		m.Description = types.StringValue(entry.Description)
	}

	return diags
}

func (m FirewallURLAliasModel) Value(ctx context.Context, urlAlias *pfsense.FirewallURLAlias) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("name"),
		"Name cannot be parsed",
		urlAlias.SetName(m.Name.ValueString()),
	)

	if !m.Description.IsNull() {
		addPathError(
			&diags,
			path.Root("description"),
			"Description cannot be parsed",
			urlAlias.SetDescription(m.Description.ValueString()),
		)
	}

	addPathError(
		&diags,
		path.Root("type"),
		"Type cannot be parsed",
		urlAlias.SetType(m.Type.ValueString()),
	)

	if !m.UpdateFrequency.IsNull() && !m.UpdateFrequency.IsUnknown() {
		addPathError(
			&diags,
			path.Root("update_frequency"),
			"Update frequency cannot be parsed",
			urlAlias.SetUpdateFrequency(int(m.UpdateFrequency.ValueInt64())),
		)
	} else if m.Type.ValueString() == "urltable" || m.Type.ValueString() == "urltable_ports" {
		// Default to 7 days if not specified for urltable types.
		urlAlias.UpdateFrequency = 7
	}

	var urlAliasEntryModels []FirewallURLAliasEntryModel
	if !m.Entries.IsNull() {
		diags.Append(m.Entries.ElementsAs(ctx, &urlAliasEntryModels, false)...)
	}

	urlAlias.Entries = make([]pfsense.FirewallURLAliasEntry, 0, len(urlAliasEntryModels))
	for index, entryModel := range urlAliasEntryModels {
		var entry pfsense.FirewallURLAliasEntry

		diags.Append(entryModel.Value(ctx, &entry, path.Root("entries").AtListIndex(index))...)
		urlAlias.Entries = append(urlAlias.Entries, entry)
	}

	return diags
}

func (m FirewallURLAliasEntryModel) Value(_ context.Context, entry *pfsense.FirewallURLAliasEntry, attrPath path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		attrPath.AtName("url"),
		"Entry URL cannot be parsed",
		entry.SetURL(m.URL.ValueString()),
	)

	if !m.Description.IsNull() {
		addPathError(
			&diags,
			attrPath.AtName("description"),
			"Entry description cannot be parsed",
			entry.SetDescription(m.Description.ValueString()),
		)
	}

	return diags
}
