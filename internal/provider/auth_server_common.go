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

type AuthServerModel struct {
	Name                     types.String `tfsdk:"name"`
	Type                     types.String `tfsdk:"type"`
	Host                     types.String `tfsdk:"host"`
	RefID                    types.String `tfsdk:"refid"`
	LDAPPort                 types.String `tfsdk:"ldap_port"`
	LDAPURLType              types.String `tfsdk:"ldap_urltype"`
	LDAPProtVer              types.String `tfsdk:"ldap_protver"`
	LDAPScope                types.String `tfsdk:"ldap_scope"`
	LDAPBaseDN               types.String `tfsdk:"ldap_basedn"`
	LDAPAuthCN               types.String `tfsdk:"ldap_authcn"`
	LDAPBindDN               types.String `tfsdk:"ldap_binddn"`
	LDAPBindPW               types.String `tfsdk:"ldap_bindpw"`
	LDAPCARef                types.String `tfsdk:"ldap_caref"`
	LDAPTimeout              types.String `tfsdk:"ldap_timeout"`
	LDAPExtendedEnabled      types.Bool   `tfsdk:"ldap_extended_enabled"`
	LDAPExtendedQuery        types.String `tfsdk:"ldap_extended_query"`
	LDAPAttrUser             types.String `tfsdk:"ldap_attr_user"`
	LDAPAttrGroup            types.String `tfsdk:"ldap_attr_group"`
	LDAPAttrMember           types.String `tfsdk:"ldap_attr_member"`
	LDAPAttrGroupObj         types.String `tfsdk:"ldap_attr_groupobj"`
	LDAPPamGroupDN           types.String `tfsdk:"ldap_pam_groupdn"`
	LDAPUTF8                 types.Bool   `tfsdk:"ldap_utf8"`
	LDAPNoStripAt            types.Bool   `tfsdk:"ldap_nostrip_at"`
	LDAPAllowUnauthenticated types.Bool   `tfsdk:"ldap_allow_unauthenticated"`
	LDAPRFC2307              types.Bool   `tfsdk:"ldap_rfc2307"`
	LDAPRFC2307UserDN        types.Bool   `tfsdk:"ldap_rfc2307_userdn"`
	RadiusProtocol           types.String `tfsdk:"radius_protocol"`
	RadiusAuthPort           types.String `tfsdk:"radius_auth_port"`
	RadiusAcctPort           types.String `tfsdk:"radius_acct_port"`
	RadiusSecret             types.String `tfsdk:"radius_secret"`
	RadiusTimeout            types.String `tfsdk:"radius_timeout"`
	RadiusNASIPAttribute     types.String `tfsdk:"radius_nasip_attribute"`
	RadiusSrvcs              types.String `tfsdk:"radius_srvcs"`
}

func (AuthServerModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"name": {
			Description: "Descriptive name for the authentication server.",
		},
		"type": {
			Description: "Server type: 'ldap' or 'radius'.",
		},
		"host": {
			Description: "Hostname or IP address of the authentication server.",
		},
		"refid": {
			Description: "Unique reference ID (auto-generated).",
		},
		"ldap_port": {
			Description: "LDAP server port number (e.g. '389' for TCP/STARTTLS, '636' for SSL).",
		},
		"ldap_urltype": {
			Description: "LDAP transport type: 'Standard TCP', 'STARTTLS Encrypted', or 'SSL/TLS Encrypted'.",
		},
		"ldap_protver": {
			Description: "LDAP protocol version: '2' or '3'.",
		},
		"ldap_scope": {
			Description: "LDAP search scope: 'one' (One Level) or 'subtree' (Entire Subtree).",
		},
		"ldap_basedn": {
			Description: "LDAP base DN for searches (e.g. 'DC=example,DC=com').",
		},
		"ldap_authcn": {
			Description: "LDAP authentication containers, semicolon-separated (e.g. 'CN=Users;OU=Staff').",
		},
		"ldap_binddn": {
			Description: "LDAP bind DN for authenticated searches (leave empty for anonymous bind).",
		},
		"ldap_bindpw": {
			Description: "LDAP bind password. Sensitive/write-only for creates/updates.",
		},
		"ldap_caref": {
			Description: "CA reference ID for LDAP SSL/STARTTLS connections.",
		},
		"ldap_timeout": {
			Description: "LDAP connection timeout in seconds (default: '25').",
		},
		"ldap_extended_enabled": {
			Description: "Enable LDAP extended query.",
		},
		"ldap_extended_query": {
			Description: "LDAP extended query string (e.g. 'memberOf=CN=VPNUsers,OU=Groups,DC=example,DC=com').",
		},
		"ldap_attr_user": {
			Description: "LDAP user naming attribute (e.g. 'cn', 'samAccountName').",
		},
		"ldap_attr_group": {
			Description: "LDAP group naming attribute (e.g. 'cn').",
		},
		"ldap_attr_member": {
			Description: "LDAP group member attribute (e.g. 'member', 'memberOf', 'uniqueMember').",
		},
		"ldap_attr_groupobj": {
			Description: "LDAP group object class (default: 'posixGroup').",
		},
		"ldap_pam_groupdn": {
			Description: "LDAP group DN for shell authentication access.",
		},
		"ldap_utf8": {
			Description: "UTF8 encode LDAP parameters.",
		},
		"ldap_nostrip_at": {
			Description: "Do not strip '@' from usernames.",
		},
		"ldap_allow_unauthenticated": {
			Description: "Allow unauthenticated bind.",
		},
		"ldap_rfc2307": {
			Description: "Use RFC 2307 style group lookups.",
		},
		"ldap_rfc2307_userdn": {
			Description: "Use DN for username search in RFC 2307 mode.",
		},
		"radius_protocol": {
			Description: "RADIUS authentication protocol: 'PAP', 'CHAP_MD5', 'MSCHAPv1', or 'MSCHAPv2'.",
		},
		"radius_auth_port": {
			Description: "RADIUS authentication port (default: '1812').",
		},
		"radius_acct_port": {
			Description: "RADIUS accounting port (default: '1813').",
		},
		"radius_secret": {
			Description: "RADIUS shared secret. Sensitive/write-only for creates/updates.",
		},
		"radius_timeout": {
			Description: "RADIUS timeout in seconds (default: '5').",
		},
		"radius_nasip_attribute": {
			Description: "RADIUS NAS-IP-Address attribute (interface name or IP).",
		},
		"radius_srvcs": {
			Description: "RADIUS services: 'both', 'auth', or 'acct'.",
		},
	}
}

func (AuthServerModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":                       types.StringType,
		"type":                       types.StringType,
		"host":                       types.StringType,
		"refid":                      types.StringType,
		"ldap_port":                  types.StringType,
		"ldap_urltype":               types.StringType,
		"ldap_protver":               types.StringType,
		"ldap_scope":                 types.StringType,
		"ldap_basedn":                types.StringType,
		"ldap_authcn":                types.StringType,
		"ldap_binddn":                types.StringType,
		"ldap_bindpw":                types.StringType,
		"ldap_caref":                 types.StringType,
		"ldap_timeout":               types.StringType,
		"ldap_extended_enabled":      types.BoolType,
		"ldap_extended_query":        types.StringType,
		"ldap_attr_user":             types.StringType,
		"ldap_attr_group":            types.StringType,
		"ldap_attr_member":           types.StringType,
		"ldap_attr_groupobj":         types.StringType,
		"ldap_pam_groupdn":           types.StringType,
		"ldap_utf8":                  types.BoolType,
		"ldap_nostrip_at":            types.BoolType,
		"ldap_allow_unauthenticated": types.BoolType,
		"ldap_rfc2307":               types.BoolType,
		"ldap_rfc2307_userdn":        types.BoolType,
		"radius_protocol":            types.StringType,
		"radius_auth_port":           types.StringType,
		"radius_acct_port":           types.StringType,
		"radius_secret":              types.StringType,
		"radius_timeout":             types.StringType,
		"radius_nasip_attribute":     types.StringType,
		"radius_srvcs":               types.StringType,
	}
}

func (m *AuthServerModel) Set(_ context.Context, server pfsense.AuthServer) diag.Diagnostics {
	var diags diag.Diagnostics

	m.Name = types.StringValue(server.Name)
	m.Type = types.StringValue(server.Type)
	m.RefID = types.StringValue(server.RefID)

	if server.Host != "" {
		m.Host = types.StringValue(server.Host)
	} else {
		m.Host = types.StringNull()
	}

	// LDAP fields.
	m.setOptionalString(&m.LDAPPort, server.LDAPPort)
	m.setOptionalString(&m.LDAPURLType, server.LDAPURLType)
	m.setOptionalString(&m.LDAPProtVer, server.LDAPProtVer)
	m.setOptionalString(&m.LDAPScope, server.LDAPScope)
	m.setOptionalString(&m.LDAPBaseDN, server.LDAPBaseDN)
	m.setOptionalString(&m.LDAPAuthCN, server.LDAPAuthCN)
	m.setOptionalString(&m.LDAPBindDN, server.LDAPBindDN)
	m.setOptionalString(&m.LDAPBindPW, server.LDAPBindPW)
	m.setOptionalString(&m.LDAPCARef, server.LDAPCARef)
	m.setOptionalString(&m.LDAPTimeout, server.LDAPTimeout)
	m.LDAPExtendedEnabled = types.BoolValue(server.LDAPExtendedEnabled)
	m.setOptionalString(&m.LDAPExtendedQuery, server.LDAPExtendedQuery)
	m.setOptionalString(&m.LDAPAttrUser, server.LDAPAttrUser)
	m.setOptionalString(&m.LDAPAttrGroup, server.LDAPAttrGroup)
	m.setOptionalString(&m.LDAPAttrMember, server.LDAPAttrMember)
	m.setOptionalString(&m.LDAPAttrGroupObj, server.LDAPAttrGroupObj)
	m.setOptionalString(&m.LDAPPamGroupDN, server.LDAPPamGroupDN)
	m.LDAPUTF8 = types.BoolValue(server.LDAPUTF8)
	m.LDAPNoStripAt = types.BoolValue(server.LDAPNoStripAt)
	m.LDAPAllowUnauthenticated = types.BoolValue(server.LDAPAllowUnauthenticated)
	m.LDAPRFC2307 = types.BoolValue(server.LDAPRFC2307)
	m.LDAPRFC2307UserDN = types.BoolValue(server.LDAPRFC2307UserDN)

	// RADIUS fields.
	m.setOptionalString(&m.RadiusProtocol, server.RadiusProtocol)
	m.setOptionalString(&m.RadiusAuthPort, server.RadiusAuthPort)
	m.setOptionalString(&m.RadiusAcctPort, server.RadiusAcctPort)
	m.setOptionalString(&m.RadiusSecret, server.RadiusSecret)
	m.setOptionalString(&m.RadiusTimeout, server.RadiusTimeout)
	m.setOptionalString(&m.RadiusNASIPAttribute, server.RadiusNASIPAttribute)
	m.setOptionalString(&m.RadiusSrvcs, server.RadiusSrvcs)

	return diags
}

func (m *AuthServerModel) setOptionalString(field *types.String, value string) {
	if value != "" {
		*field = types.StringValue(value)
	} else {
		*field = types.StringNull()
	}
}

func (m AuthServerModel) Value(ctx context.Context, server *pfsense.AuthServer) diag.Diagnostics {
	var diags diag.Diagnostics

	addPathError(
		&diags,
		path.Root("name"),
		"Auth server name cannot be parsed",
		server.SetName(m.Name.ValueString()),
	)

	addPathError(
		&diags,
		path.Root("type"),
		"Auth server type cannot be parsed",
		server.SetType(m.Type.ValueString()),
	)

	if !m.Host.IsNull() {
		addPathError(
			&diags,
			path.Root("host"),
			"Host cannot be parsed",
			server.SetHost(m.Host.ValueString()),
		)
	}

	// LDAP fields.
	m.setOptionalStringValue(ctx, &diags, "ldap_port", m.LDAPPort, server.SetLDAPPort)
	m.setOptionalStringValue(ctx, &diags, "ldap_urltype", m.LDAPURLType, server.SetLDAPURLType)
	m.setOptionalStringValue(ctx, &diags, "ldap_protver", m.LDAPProtVer, server.SetLDAPProtVer)
	m.setOptionalStringValue(ctx, &diags, "ldap_scope", m.LDAPScope, server.SetLDAPScope)
	m.setOptionalStringValue(ctx, &diags, "ldap_basedn", m.LDAPBaseDN, server.SetLDAPBaseDN)
	m.setOptionalStringValue(ctx, &diags, "ldap_authcn", m.LDAPAuthCN, server.SetLDAPAuthCN)
	m.setOptionalStringValue(ctx, &diags, "ldap_binddn", m.LDAPBindDN, server.SetLDAPBindDN)
	m.setOptionalStringValue(ctx, &diags, "ldap_bindpw", m.LDAPBindPW, server.SetLDAPBindPW)
	m.setOptionalStringValue(ctx, &diags, "ldap_caref", m.LDAPCARef, server.SetLDAPCARef)
	m.setOptionalStringValue(ctx, &diags, "ldap_timeout", m.LDAPTimeout, server.SetLDAPTimeout)

	addPathError(
		&diags,
		path.Root("ldap_extended_enabled"),
		"LDAP extended enabled cannot be parsed",
		server.SetLDAPExtendedEnabled(m.LDAPExtendedEnabled.ValueBool()),
	)

	m.setOptionalStringValue(ctx, &diags, "ldap_extended_query", m.LDAPExtendedQuery, server.SetLDAPExtendedQuery)
	m.setOptionalStringValue(ctx, &diags, "ldap_attr_user", m.LDAPAttrUser, server.SetLDAPAttrUser)
	m.setOptionalStringValue(ctx, &diags, "ldap_attr_group", m.LDAPAttrGroup, server.SetLDAPAttrGroup)
	m.setOptionalStringValue(ctx, &diags, "ldap_attr_member", m.LDAPAttrMember, server.SetLDAPAttrMember)
	m.setOptionalStringValue(ctx, &diags, "ldap_attr_groupobj", m.LDAPAttrGroupObj, server.SetLDAPAttrGroupObj)
	m.setOptionalStringValue(ctx, &diags, "ldap_pam_groupdn", m.LDAPPamGroupDN, server.SetLDAPPamGroupDN)

	addPathError(&diags, path.Root("ldap_utf8"), "LDAP UTF8 cannot be parsed", server.SetLDAPUTF8(m.LDAPUTF8.ValueBool()))
	addPathError(&diags, path.Root("ldap_nostrip_at"), "LDAP no strip at cannot be parsed", server.SetLDAPNoStripAt(m.LDAPNoStripAt.ValueBool()))
	addPathError(&diags, path.Root("ldap_allow_unauthenticated"), "LDAP allow unauthenticated cannot be parsed", server.SetLDAPAllowUnauthenticated(m.LDAPAllowUnauthenticated.ValueBool()))
	addPathError(&diags, path.Root("ldap_rfc2307"), "LDAP RFC2307 cannot be parsed", server.SetLDAPRFC2307(m.LDAPRFC2307.ValueBool()))
	addPathError(&diags, path.Root("ldap_rfc2307_userdn"), "LDAP RFC2307 UserDN cannot be parsed", server.SetLDAPRFC2307UserDN(m.LDAPRFC2307UserDN.ValueBool()))

	// RADIUS fields.
	m.setOptionalStringValue(ctx, &diags, "radius_protocol", m.RadiusProtocol, server.SetRadiusProtocol)
	m.setOptionalStringValue(ctx, &diags, "radius_auth_port", m.RadiusAuthPort, server.SetRadiusAuthPort)
	m.setOptionalStringValue(ctx, &diags, "radius_acct_port", m.RadiusAcctPort, server.SetRadiusAcctPort)
	m.setOptionalStringValue(ctx, &diags, "radius_secret", m.RadiusSecret, server.SetRadiusSecret)
	m.setOptionalStringValue(ctx, &diags, "radius_timeout", m.RadiusTimeout, server.SetRadiusTimeout)
	m.setOptionalStringValue(ctx, &diags, "radius_nasip_attribute", m.RadiusNASIPAttribute, server.SetRadiusNASIPAttribute)
	m.setOptionalStringValue(ctx, &diags, "radius_srvcs", m.RadiusSrvcs, server.SetRadiusSrvcs)

	return diags
}

func (m AuthServerModel) setOptionalStringValue(_ context.Context, diags *diag.Diagnostics, attrName string, field types.String, setter func(string) error) {
	if !field.IsNull() {
		addPathError(
			diags,
			path.Root(attrName),
			fmt.Sprintf("%s cannot be parsed", attrName),
			setter(field.ValueString()),
		)
	}
}
