package pfsense

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var authServerNameRegex = regexp.MustCompile(`^[a-zA-Z0-9.\-_ ]+$`)

const (
	AuthServerMaxNameLength = 255
)

// authServerResponse represents the raw JSON response from pfSense PHP commands.
type authServerResponse struct {
	RefID                    string      `json:"refid"`
	Name                     string      `json:"name"`
	Type                     string      `json:"type"`
	Host                     string      `json:"host"`
	LDAPPort                 string      `json:"ldap_port"`
	LDAPURLType              string      `json:"ldap_urltype"`
	LDAPProtVer              string      `json:"ldap_protver"`
	LDAPScope                string      `json:"ldap_scope"`
	LDAPBaseDN               string      `json:"ldap_basedn"`
	LDAPAuthCN               string      `json:"ldap_authcn"`
	LDAPBindDN               string      `json:"ldap_binddn"`
	LDAPBindPW               string      `json:"ldap_bindpw"`
	LDAPCARef                string      `json:"ldap_caref"`
	LDAPTimeout              string      `json:"ldap_timeout"`
	LDAPExtendedEnabled      interface{} `json:"ldap_extended_enabled"`
	LDAPExtendedQuery        string      `json:"ldap_extended_query"`
	LDAPAttrUser             string      `json:"ldap_attr_user"`
	LDAPAttrGroup            string      `json:"ldap_attr_group"`
	LDAPAttrMember           string      `json:"ldap_attr_member"`
	LDAPAttrGroupObj         string      `json:"ldap_attr_groupobj"`
	LDAPPamGroupDN           string      `json:"ldap_pam_groupdn"`
	LDAPUTF8                 interface{} `json:"ldap_utf8"`
	LDAPNoStripAt            interface{} `json:"ldap_nostrip_at"`
	LDAPAllowUnauthenticated interface{} `json:"ldap_allow_unauthenticated"`
	LDAPRFC2307              interface{} `json:"ldap_rfc2307"`
	LDAPRFC2307UserDN        interface{} `json:"ldap_rfc2307_userdn"`
	RadiusProtocol           string      `json:"radius_protocol"`
	RadiusAuthPort           string      `json:"radius_auth_port"`
	RadiusAcctPort           string      `json:"radius_acct_port"`
	RadiusSecret             string      `json:"radius_secret"`
	RadiusTimeout            string      `json:"radius_timeout"`
	RadiusNASIPAttribute     string      `json:"radius_nasip_attribute"`
	RadiusSrvcs              string      `json:"radius_srvcs"`
	ControlID                int         `json:"controlID"` //nolint:tagliatelle
}

// AuthServer is the clean domain struct for a pfSense authentication server.
type AuthServer struct {
	RefID                    string
	Name                     string
	Type                     string // "ldap" or "radius"
	Host                     string
	LDAPPort                 string
	LDAPURLType              string
	LDAPProtVer              string
	LDAPScope                string
	LDAPBaseDN               string
	LDAPAuthCN               string // semicolon-separated containers
	LDAPBindDN               string
	LDAPBindPW               string // sensitive
	LDAPCARef                string
	LDAPTimeout              string
	LDAPExtendedEnabled      bool
	LDAPExtendedQuery        string
	LDAPAttrUser             string
	LDAPAttrGroup            string
	LDAPAttrMember           string
	LDAPAttrGroupObj         string
	LDAPPamGroupDN           string
	LDAPUTF8                 bool
	LDAPNoStripAt            bool
	LDAPAllowUnauthenticated bool
	LDAPRFC2307              bool
	LDAPRFC2307UserDN        bool
	RadiusProtocol           string
	RadiusAuthPort           string
	RadiusAcctPort           string
	RadiusSecret             string // sensitive
	RadiusTimeout            string
	RadiusNASIPAttribute     string
	RadiusSrvcs              string // "both", "auth", or "acct"
	controlID                int
}

func (a *AuthServer) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("%w, auth server name is required", ErrClientValidation)
	}

	if len(name) > AuthServerMaxNameLength {
		return fmt.Errorf("%w, auth server name must be at most %d characters", ErrClientValidation, AuthServerMaxNameLength)
	}

	if !authServerNameRegex.MatchString(name) {
		return fmt.Errorf("%w, auth server name must contain only alphanumeric characters, dots, hyphens, underscores, and spaces", ErrClientValidation)
	}

	a.Name = name

	return nil
}

func (a *AuthServer) SetType(t string) error {
	if t != "ldap" && t != "radius" {
		return fmt.Errorf("%w, auth server type must be 'ldap' or 'radius'", ErrClientValidation)
	}

	a.Type = t

	return nil
}

func (a *AuthServer) SetHost(host string) error {
	a.Host = host

	return nil
}

func (a *AuthServer) SetLDAPPort(port string) error {
	a.LDAPPort = port

	return nil
}

func (a *AuthServer) SetLDAPURLType(urlType string) error {
	a.LDAPURLType = urlType

	return nil
}

func (a *AuthServer) SetLDAPProtVer(protVer string) error {
	a.LDAPProtVer = protVer

	return nil
}

func (a *AuthServer) SetLDAPScope(scope string) error {
	a.LDAPScope = scope

	return nil
}

func (a *AuthServer) SetLDAPBaseDN(baseDN string) error {
	a.LDAPBaseDN = baseDN

	return nil
}

func (a *AuthServer) SetLDAPAuthCN(authCN string) error {
	a.LDAPAuthCN = authCN

	return nil
}

func (a *AuthServer) SetLDAPBindDN(bindDN string) error {
	a.LDAPBindDN = bindDN

	return nil
}

func (a *AuthServer) SetLDAPBindPW(bindPW string) error {
	a.LDAPBindPW = bindPW

	return nil
}

func (a *AuthServer) SetLDAPCARef(caRef string) error {
	a.LDAPCARef = caRef

	return nil
}

func (a *AuthServer) SetLDAPTimeout(timeout string) error {
	a.LDAPTimeout = timeout

	return nil
}

func (a *AuthServer) SetLDAPExtendedEnabled(enabled bool) error {
	a.LDAPExtendedEnabled = enabled

	return nil
}

func (a *AuthServer) SetLDAPExtendedQuery(query string) error {
	a.LDAPExtendedQuery = query

	return nil
}

func (a *AuthServer) SetLDAPAttrUser(attr string) error {
	a.LDAPAttrUser = attr

	return nil
}

func (a *AuthServer) SetLDAPAttrGroup(attr string) error {
	a.LDAPAttrGroup = attr

	return nil
}

func (a *AuthServer) SetLDAPAttrMember(attr string) error {
	a.LDAPAttrMember = attr

	return nil
}

func (a *AuthServer) SetLDAPAttrGroupObj(attr string) error {
	a.LDAPAttrGroupObj = attr

	return nil
}

func (a *AuthServer) SetLDAPPamGroupDN(dn string) error {
	a.LDAPPamGroupDN = dn

	return nil
}

func (a *AuthServer) SetLDAPUTF8(utf8 bool) error {
	a.LDAPUTF8 = utf8

	return nil
}

func (a *AuthServer) SetLDAPNoStripAt(noStripAt bool) error {
	a.LDAPNoStripAt = noStripAt

	return nil
}

func (a *AuthServer) SetLDAPAllowUnauthenticated(allow bool) error {
	a.LDAPAllowUnauthenticated = allow

	return nil
}

func (a *AuthServer) SetLDAPRFC2307(rfc2307 bool) error {
	a.LDAPRFC2307 = rfc2307

	return nil
}

func (a *AuthServer) SetLDAPRFC2307UserDN(userDN bool) error {
	a.LDAPRFC2307UserDN = userDN

	return nil
}

func (a *AuthServer) SetRadiusProtocol(protocol string) error {
	a.RadiusProtocol = protocol

	return nil
}

func (a *AuthServer) SetRadiusAuthPort(port string) error {
	a.RadiusAuthPort = port

	return nil
}

func (a *AuthServer) SetRadiusAcctPort(port string) error {
	a.RadiusAcctPort = port

	return nil
}

func (a *AuthServer) SetRadiusSecret(secret string) error {
	a.RadiusSecret = secret

	return nil
}

func (a *AuthServer) SetRadiusTimeout(timeout string) error {
	a.RadiusTimeout = timeout

	return nil
}

func (a *AuthServer) SetRadiusNASIPAttribute(attr string) error {
	a.RadiusNASIPAttribute = attr

	return nil
}

func (a *AuthServer) SetRadiusSrvcs(srvcs string) error {
	a.RadiusSrvcs = srvcs

	return nil
}

// AuthServers is a named type for a slice of AuthServer.
type AuthServers []AuthServer

func (servers AuthServers) GetByName(name string) (*AuthServer, error) {
	for _, s := range servers {
		if s.Name == name {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("auth server %w with name '%s'", ErrNotFound, name)
}

func (servers AuthServers) GetControlIDByName(name string) (*int, error) {
	for _, s := range servers {
		if s.Name == name {
			return &s.controlID, nil
		}
	}

	return nil, fmt.Errorf("auth server %w with name '%s'", ErrNotFound, name)
}

func parseAuthServerResponse(resp authServerResponse) (AuthServer, error) {
	var server AuthServer

	if err := server.SetName(resp.Name); err != nil {
		return server, err
	}

	if err := server.SetType(resp.Type); err != nil {
		return server, err
	}

	server.RefID = resp.RefID
	server.Host = resp.Host
	server.controlID = resp.ControlID

	// LDAP fields.
	server.LDAPPort = resp.LDAPPort
	server.LDAPURLType = resp.LDAPURLType
	server.LDAPProtVer = resp.LDAPProtVer
	server.LDAPScope = resp.LDAPScope
	server.LDAPBaseDN = resp.LDAPBaseDN
	server.LDAPAuthCN = resp.LDAPAuthCN
	server.LDAPBindDN = resp.LDAPBindDN
	server.LDAPBindPW = resp.LDAPBindPW
	server.LDAPCARef = resp.LDAPCARef
	server.LDAPTimeout = resp.LDAPTimeout
	server.LDAPExtendedEnabled = parseBoolField(resp.LDAPExtendedEnabled)
	server.LDAPExtendedQuery = resp.LDAPExtendedQuery
	server.LDAPAttrUser = resp.LDAPAttrUser
	server.LDAPAttrGroup = resp.LDAPAttrGroup
	server.LDAPAttrMember = resp.LDAPAttrMember
	server.LDAPAttrGroupObj = resp.LDAPAttrGroupObj
	server.LDAPPamGroupDN = resp.LDAPPamGroupDN
	server.LDAPUTF8 = parseBoolField(resp.LDAPUTF8)
	server.LDAPNoStripAt = parseBoolField(resp.LDAPNoStripAt)
	server.LDAPAllowUnauthenticated = parseBoolField(resp.LDAPAllowUnauthenticated)
	server.LDAPRFC2307 = parseBoolField(resp.LDAPRFC2307)
	server.LDAPRFC2307UserDN = parseBoolField(resp.LDAPRFC2307UserDN)

	// RADIUS fields.
	server.RadiusProtocol = resp.RadiusProtocol
	server.RadiusAuthPort = resp.RadiusAuthPort
	server.RadiusAcctPort = resp.RadiusAcctPort
	server.RadiusSecret = resp.RadiusSecret
	server.RadiusTimeout = resp.RadiusTimeout
	server.RadiusNASIPAttribute = resp.RadiusNASIPAttribute
	server.RadiusSrvcs = resp.RadiusSrvcs

	return server, nil
}

func (pf *Client) getAuthServers(ctx context.Context) (*AuthServers, error) {
	command := "$output = array();" +
		"$servers = config_get_path('system/authserver', array());" +
		"foreach ($servers as $k => $v) {" +
		"$v['controlID'] = $k;" +
		// Normalize boolean fields: presence-based -> proper booleans for JSON.
		"$v['ldap_extended_enabled'] = array_key_exists('ldap_extended_enabled', $v);" +
		"$v['ldap_utf8'] = array_key_exists('ldap_utf8', $v);" +
		"$v['ldap_nostrip_at'] = array_key_exists('ldap_nostrip_at', $v);" +
		"$v['ldap_allow_unauthenticated'] = array_key_exists('ldap_allow_unauthenticated', $v);" +
		"$v['ldap_rfc2307'] = array_key_exists('ldap_rfc2307', $v);" +
		"$v['ldap_rfc2307_userdn'] = array_key_exists('ldap_rfc2307_userdn', $v);" +
		// Ensure optional string fields have defaults for consistent JSON.
		"if (!isset($v['ldap_port'])) { $v['ldap_port'] = ''; }" +
		"if (!isset($v['ldap_urltype'])) { $v['ldap_urltype'] = ''; }" +
		"if (!isset($v['ldap_protver'])) { $v['ldap_protver'] = ''; }" +
		"if (!isset($v['ldap_scope'])) { $v['ldap_scope'] = ''; }" +
		"if (!isset($v['ldap_basedn'])) { $v['ldap_basedn'] = ''; }" +
		"if (!isset($v['ldap_authcn'])) { $v['ldap_authcn'] = ''; }" +
		"if (!isset($v['ldap_binddn'])) { $v['ldap_binddn'] = ''; }" +
		"if (!isset($v['ldap_bindpw'])) { $v['ldap_bindpw'] = ''; }" +
		"if (!isset($v['ldap_caref'])) { $v['ldap_caref'] = ''; }" +
		"if (!isset($v['ldap_timeout'])) { $v['ldap_timeout'] = ''; }" +
		"if (!isset($v['ldap_extended_query'])) { $v['ldap_extended_query'] = ''; }" +
		"if (!isset($v['ldap_attr_user'])) { $v['ldap_attr_user'] = ''; }" +
		"if (!isset($v['ldap_attr_group'])) { $v['ldap_attr_group'] = ''; }" +
		"if (!isset($v['ldap_attr_member'])) { $v['ldap_attr_member'] = ''; }" +
		"if (!isset($v['ldap_attr_groupobj'])) { $v['ldap_attr_groupobj'] = ''; }" +
		"if (!isset($v['ldap_pam_groupdn'])) { $v['ldap_pam_groupdn'] = ''; }" +
		"if (!isset($v['radius_protocol'])) { $v['radius_protocol'] = ''; }" +
		"if (!isset($v['radius_auth_port'])) { $v['radius_auth_port'] = ''; }" +
		"if (!isset($v['radius_acct_port'])) { $v['radius_acct_port'] = ''; }" +
		"if (!isset($v['radius_secret'])) { $v['radius_secret'] = ''; }" +
		"if (!isset($v['radius_timeout'])) { $v['radius_timeout'] = ''; }" +
		"if (!isset($v['radius_nasip_attribute'])) { $v['radius_nasip_attribute'] = ''; }" +
		"if (!isset($v['radius_srvcs'])) { $v['radius_srvcs'] = ''; }" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var serverResp []authServerResponse
	if err := pf.executePHPCommand(ctx, command, &serverResp); err != nil {
		return nil, err
	}

	servers := make(AuthServers, 0, len(serverResp))
	for _, resp := range serverResp {
		s, err := parseAuthServerResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w auth server response, %w", ErrUnableToParse, err)
		}

		servers = append(servers, s)
	}

	return &servers, nil
}

func (pf *Client) GetAuthServers(ctx context.Context) (*AuthServers, error) {
	defer pf.read(&pf.mutexes.AuthServer)()

	servers, err := pf.getAuthServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w auth servers, %w", ErrGetOperationFailed, err)
	}

	return servers, nil
}

func (pf *Client) GetAuthServer(ctx context.Context, name string) (*AuthServer, error) {
	defer pf.read(&pf.mutexes.AuthServer)()

	servers, err := pf.getAuthServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w auth servers, %w", ErrGetOperationFailed, err)
	}

	s, err := servers.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w auth server, %w", ErrGetOperationFailed, err)
	}

	return s, nil
}

// buildAuthServerPHP builds the PHP statements to set all fields on a $server variable.
func buildAuthServerPHP(req AuthServer) string {
	var cmd strings.Builder

	cmd.WriteString(fmt.Sprintf("$server['name'] = '%s';", phpEscape(req.Name)))
	cmd.WriteString(fmt.Sprintf("$server['type'] = '%s';", phpEscape(req.Type)))
	cmd.WriteString(fmt.Sprintf("$server['host'] = '%s';", phpEscape(req.Host)))

	if req.Type == "ldap" {
		cmd.WriteString(fmt.Sprintf("$server['ldap_port'] = '%s';", phpEscape(req.LDAPPort)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_urltype'] = '%s';", phpEscape(req.LDAPURLType)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_protver'] = '%s';", phpEscape(req.LDAPProtVer)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_scope'] = '%s';", phpEscape(req.LDAPScope)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_basedn'] = '%s';", phpEscape(req.LDAPBaseDN)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_authcn'] = '%s';", phpEscape(req.LDAPAuthCN)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_binddn'] = '%s';", phpEscape(req.LDAPBindDN)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_bindpw'] = '%s';", phpEscape(req.LDAPBindPW)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_caref'] = '%s';", phpEscape(req.LDAPCARef)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_timeout'] = '%s';", phpEscape(req.LDAPTimeout)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_extended_query'] = '%s';", phpEscape(req.LDAPExtendedQuery)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_attr_user'] = '%s';", phpEscape(req.LDAPAttrUser)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_attr_group'] = '%s';", phpEscape(req.LDAPAttrGroup)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_attr_member'] = '%s';", phpEscape(req.LDAPAttrMember)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_attr_groupobj'] = '%s';", phpEscape(req.LDAPAttrGroupObj)))
		cmd.WriteString(fmt.Sprintf("$server['ldap_pam_groupdn'] = '%s';", phpEscape(req.LDAPPamGroupDN)))

		// Boolean fields: presence-based in pfSense config.
		if req.LDAPExtendedEnabled {
			cmd.WriteString("$server['ldap_extended_enabled'] = '';")
		} else {
			cmd.WriteString("unset($server['ldap_extended_enabled']);")
		}

		if req.LDAPUTF8 {
			cmd.WriteString("$server['ldap_utf8'] = '';")
		} else {
			cmd.WriteString("unset($server['ldap_utf8']);")
		}

		if req.LDAPNoStripAt {
			cmd.WriteString("$server['ldap_nostrip_at'] = '';")
		} else {
			cmd.WriteString("unset($server['ldap_nostrip_at']);")
		}

		if req.LDAPAllowUnauthenticated {
			cmd.WriteString("$server['ldap_allow_unauthenticated'] = '';")
		} else {
			cmd.WriteString("unset($server['ldap_allow_unauthenticated']);")
		}

		if req.LDAPRFC2307 {
			cmd.WriteString("$server['ldap_rfc2307'] = '';")
		} else {
			cmd.WriteString("unset($server['ldap_rfc2307']);")
		}

		if req.LDAPRFC2307UserDN {
			cmd.WriteString("$server['ldap_rfc2307_userdn'] = '';")
		} else {
			cmd.WriteString("unset($server['ldap_rfc2307_userdn']);")
		}
	}

	if req.Type == "radius" {
		cmd.WriteString(fmt.Sprintf("$server['radius_protocol'] = '%s';", phpEscape(req.RadiusProtocol)))
		cmd.WriteString(fmt.Sprintf("$server['radius_auth_port'] = '%s';", phpEscape(req.RadiusAuthPort)))
		cmd.WriteString(fmt.Sprintf("$server['radius_acct_port'] = '%s';", phpEscape(req.RadiusAcctPort)))
		cmd.WriteString(fmt.Sprintf("$server['radius_secret'] = '%s';", phpEscape(req.RadiusSecret)))
		cmd.WriteString(fmt.Sprintf("$server['radius_timeout'] = '%s';", phpEscape(req.RadiusTimeout)))
		cmd.WriteString(fmt.Sprintf("$server['radius_nasip_attribute'] = '%s';", phpEscape(req.RadiusNASIPAttribute)))
		cmd.WriteString(fmt.Sprintf("$server['radius_srvcs'] = '%s';", phpEscape(req.RadiusSrvcs)))
	}

	return cmd.String()
}

func (pf *Client) CreateAuthServer(ctx context.Context, req AuthServer) (*AuthServer, error) {
	defer pf.write(&pf.mutexes.AuthServer)()

	// Check for duplicate name.
	existingServers, err := pf.getAuthServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w auth servers for duplicate check, %w", ErrGetOperationFailed, err)
	}

	if _, err := existingServers.GetByName(req.Name); err == nil {
		return nil, fmt.Errorf("%w auth server, an auth server with name '%s' already exists", ErrCreateOperationFailed, req.Name)
	}

	var cmd strings.Builder

	// Build server array.
	cmd.WriteString("$server = array();")

	// Generate a unique refid.
	cmd.WriteString("$server['refid'] = uniqid();")

	cmd.WriteString(buildAuthServerPHP(req))

	// Append to config.
	cmd.WriteString("config_set_path('system/authserver/', $server);")

	// Write config.
	cmd.WriteString(fmt.Sprintf("write_config('Terraform: created auth server %s');", phpEscape(req.Name)))

	cmd.WriteString("print(json_encode(true));")

	var result bool
	if err := pf.executePHPCommand(ctx, cmd.String(), &result); err != nil {
		return nil, fmt.Errorf("%w auth server, %w", ErrCreateOperationFailed, err)
	}

	servers, err := pf.getAuthServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w auth servers after creating, %w", ErrGetOperationFailed, err)
	}

	s, err := servers.GetByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w auth server after creating, %w", ErrGetOperationFailed, err)
	}

	return s, nil
}

func (pf *Client) UpdateAuthServer(ctx context.Context, req AuthServer) (*AuthServer, error) {
	defer pf.write(&pf.mutexes.AuthServer)()

	servers, err := pf.getAuthServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w auth servers, %w", ErrGetOperationFailed, err)
	}

	controlID, err := servers.GetControlIDByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w auth server, %w", ErrGetOperationFailed, err)
	}

	var cmd strings.Builder

	// Load existing server to preserve read-only fields (refid).
	cmd.WriteString(fmt.Sprintf("$server = config_get_path('system/authserver/%d');", *controlID))

	cmd.WriteString(buildAuthServerPHP(req))

	// Write server back.
	cmd.WriteString(fmt.Sprintf("config_set_path('system/authserver/%d', $server);", *controlID))

	// Write config.
	cmd.WriteString(fmt.Sprintf("write_config('Terraform: updated auth server %s');", phpEscape(req.Name)))

	cmd.WriteString("print(json_encode(true));")

	var result bool
	if err := pf.executePHPCommand(ctx, cmd.String(), &result); err != nil {
		return nil, fmt.Errorf("%w auth server, %w", ErrUpdateOperationFailed, err)
	}

	servers, err = pf.getAuthServers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w auth servers after updating, %w", ErrGetOperationFailed, err)
	}

	s, err := servers.GetByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w auth server after updating, %w", ErrGetOperationFailed, err)
	}

	return s, nil
}

func (pf *Client) DeleteAuthServer(ctx context.Context, name string) error {
	defer pf.write(&pf.mutexes.AuthServer)()

	servers, err := pf.getAuthServers(ctx)
	if err != nil {
		return fmt.Errorf("%w auth servers, %w", ErrGetOperationFailed, err)
	}

	controlID, err := servers.GetControlIDByName(name)
	if err != nil {
		return fmt.Errorf("%w auth server, %w", ErrGetOperationFailed, err)
	}

	var cmd strings.Builder

	// Remove from config.
	cmd.WriteString(fmt.Sprintf("config_del_path('system/authserver/%d');", *controlID))

	// Reindex array.
	cmd.WriteString("config_set_path('system/authserver', array_values(config_get_path('system/authserver', array())));")

	// Write config.
	cmd.WriteString(fmt.Sprintf("write_config('Terraform: deleted auth server %s');", phpEscape(name)))

	cmd.WriteString("print(json_encode(true));")

	var result bool
	if err := pf.executePHPCommand(ctx, cmd.String(), &result); err != nil {
		return fmt.Errorf("%w auth server, %w", ErrDeleteOperationFailed, err)
	}

	// Verify deletion.
	servers, err = pf.getAuthServers(ctx)
	if err != nil {
		return fmt.Errorf("%w auth servers after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := servers.GetByName(name); err == nil {
		return fmt.Errorf("%w auth server, still exists", ErrDeleteOperationFailed)
	}

	return nil
}
