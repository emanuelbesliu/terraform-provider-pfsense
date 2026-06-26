package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type DynamicDNSModel struct {
	ID                types.Int64  `tfsdk:"id"`
	Type              types.String `tfsdk:"type"`
	Interface         types.String `tfsdk:"interface"`
	Host              types.String `tfsdk:"host"`
	DomainName        types.String `tfsdk:"domain_name"`
	Username          types.String `tfsdk:"username"`
	Password          types.String `tfsdk:"password"`
	MX                types.String `tfsdk:"mx"`
	Wildcard          types.Bool   `tfsdk:"wildcard"`
	Proxied           types.Bool   `tfsdk:"proxied"`
	VerboseLog        types.Bool   `tfsdk:"verbose_log"`
	CurlIPResolveV4   types.Bool   `tfsdk:"curl_ipresolve_v4"`
	CurlSSLVerifyPeer types.Bool   `tfsdk:"curl_ssl_verifypeer"`
	ZoneID            types.String `tfsdk:"zone_id"`
	TTL               types.String `tfsdk:"ttl"`
	MaxCacheAge       types.String `tfsdk:"max_cache_age"`
	UpdateURL         types.String `tfsdk:"update_url"`
	ResultMatch       types.String `tfsdk:"result_match"`
	RequestIf         types.String `tfsdk:"request_interface"`
	CurlProxy         types.String `tfsdk:"curl_proxy"`
	Description       types.String `tfsdk:"description"`
	Disabled          types.Bool   `tfsdk:"disabled"`
}

func (DynamicDNSModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"id": {
			Description: "Numeric index of the dynamic DNS entry in pfSense configuration.",
		},
		"type": {
			Description: "Dynamic DNS service type (e.g. 'cloudflare-v6', 'dyndns', 'noip', 'custom').",
		},
		"interface": {
			Description: "Interface to monitor for IP changes (e.g. 'wan', 'opt1').",
		},
		"host": {
			Description: "Hostname to update.",
		},
		"domain_name": {
			Description: "Domain name for the DNS record.",
		},
		"username": {
			Description: "Username or API token for the DNS provider.",
		},
		"password": {
			Description: "Password or API key for the DNS provider.",
		},
		"mx": {
			Description: "MX record hostname.",
		},
		"wildcard": {
			Description: "Enable wildcard DNS entry.",
		},
		"proxied": {
			Description: "Enable proxy mode (Cloudflare-specific).",
		},
		"verbose_log": {
			Description: "Enable verbose logging.",
		},
		"curl_ipresolve_v4": {
			Description: "Force IPv4 resolution for cURL.",
		},
		"curl_ssl_verifypeer": {
			Description: "Verify SSL peer certificate for cURL.",
		},
		"zone_id": {
			Description: "Zone ID (Cloudflare-specific).",
		},
		"ttl": {
			Description: "TTL for the DNS record.",
		},
		"max_cache_age": {
			Description: "Maximum cache age in seconds before forcing an update.",
		},
		"update_url": {
			Description: "Custom update URL (for custom DynDNS type).",
		},
		"result_match": {
			Description: "Expected result match string (for custom DynDNS type).",
		},
		"request_interface": {
			Description: "Interface to use for sending update requests.",
		},
		"curl_proxy": {
			Description: "HTTP proxy for cURL requests.",
		},
		"description": {
			Description: "Description of the dynamic DNS entry.",
		},
		"disabled": {
			Description: "Whether the dynamic DNS entry is disabled.",
		},
	}
}

func (DynamicDNSModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                  types.Int64Type,
		"type":                types.StringType,
		"interface":           types.StringType,
		"host":                types.StringType,
		"domain_name":         types.StringType,
		"username":            types.StringType,
		"password":            types.StringType,
		"mx":                  types.StringType,
		"wildcard":            types.BoolType,
		"proxied":             types.BoolType,
		"verbose_log":         types.BoolType,
		"curl_ipresolve_v4":   types.BoolType,
		"curl_ssl_verifypeer": types.BoolType,
		"zone_id":             types.StringType,
		"ttl":                 types.StringType,
		"max_cache_age":       types.StringType,
		"update_url":          types.StringType,
		"result_match":        types.StringType,
		"request_interface":   types.StringType,
		"curl_proxy":          types.StringType,
		"description":         types.StringType,
		"disabled":            types.BoolType,
	}
}

func (m *DynamicDNSModel) Set(_ context.Context, entry pfsense.DynamicDNS) diag.Diagnostics {
	var diags diag.Diagnostics

	m.ID = types.Int64Value(int64(entry.ControlID()))
	m.Type = types.StringValue(entry.Type)
	m.Interface = types.StringValue(entry.Interface)
	m.Host = types.StringValue(entry.Host)
	m.DomainName = types.StringValue(entry.DomainName)
	m.Username = types.StringValue(entry.Username)
	m.Password = types.StringValue(entry.Password)
	m.MX = types.StringValue(entry.MX)
	m.Wildcard = types.BoolValue(entry.Wildcard)
	m.Proxied = types.BoolValue(entry.Proxied)
	m.VerboseLog = types.BoolValue(entry.VerboseLog)
	m.CurlIPResolveV4 = types.BoolValue(entry.CurlIPResolveV4)
	m.CurlSSLVerifyPeer = types.BoolValue(entry.CurlSSLVerifyPeer)
	m.ZoneID = types.StringValue(entry.ZoneID)
	m.TTL = types.StringValue(entry.TTL)
	m.MaxCacheAge = types.StringValue(entry.MaxCacheAge)
	m.UpdateURL = types.StringValue(entry.UpdateURL)
	m.ResultMatch = types.StringValue(entry.ResultMatch)
	m.RequestIf = types.StringValue(entry.RequestIf)
	m.CurlProxy = types.StringValue(entry.CurlProxy)
	m.Description = types.StringValue(entry.Description)
	m.Disabled = types.BoolValue(entry.Disabled)

	return diags
}

func (m DynamicDNSModel) Value(_ context.Context, entry *pfsense.DynamicDNS) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("type"),
		"Type cannot be parsed",
		entry.SetType(m.Type.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("interface"),
		"Interface cannot be parsed",
		entry.SetInterface(m.Interface.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("host"),
		"Host cannot be parsed",
		entry.SetHost(m.Host.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("domain_name"),
		"Domain name cannot be parsed",
		entry.SetDomainName(m.DomainName.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("username"),
		"Username cannot be parsed",
		entry.SetUsername(m.Username.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("password"),
		"Password cannot be parsed",
		entry.SetPassword(m.Password.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("description"),
		"Description cannot be parsed",
		entry.SetDescription(m.Description.ValueString()),
	)

	entry.MX = m.MX.ValueString()
	entry.Wildcard = m.Wildcard.ValueBool()
	entry.Proxied = m.Proxied.ValueBool()
	entry.VerboseLog = m.VerboseLog.ValueBool()
	entry.CurlIPResolveV4 = m.CurlIPResolveV4.ValueBool()
	entry.CurlSSLVerifyPeer = m.CurlSSLVerifyPeer.ValueBool()
	entry.ZoneID = m.ZoneID.ValueString()
	entry.TTL = m.TTL.ValueString()
	entry.MaxCacheAge = m.MaxCacheAge.ValueString()
	entry.UpdateURL = m.UpdateURL.ValueString()
	entry.ResultMatch = m.ResultMatch.ValueString()
	entry.RequestIf = m.RequestIf.ValueString()
	entry.CurlProxy = m.CurlProxy.ValueString()
	entry.Disabled = m.Disabled.ValueBool()

	return diags
}
