package pfsense

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

var userNameRegex = regexp.MustCompile(`^[a-zA-Z0-9.\-_]+$`)

const (
	UserMaxNameLength = 32
)

type userResponse struct {
	Name           string      `json:"name"`
	Description    string      `json:"descr"`
	Scope          string      `json:"scope"`
	UID            string      `json:"uid"`
	Disabled       interface{} `json:"disabled"`
	Expires        string      `json:"expires"`
	AuthorizedKeys string      `json:"authorizedkeys"`
	IPSecPSK       string      `json:"ipsecpsk"`
	Privileges     interface{} `json:"priv"`
	Groups         interface{} `json:"groups"`
	CustomSettings interface{} `json:"customsettings"`
	WebGUICss      string      `json:"webguicss"`
	DashboardCols  string      `json:"dashboardcolumns"`
	KeepHistory    interface{} `json:"keephistory"`
	ControlID      int         `json:"controlID"` //nolint:tagliatelle
}

type User struct {
	Name           string
	Description    string
	Scope          string
	UID            string
	Disabled       bool
	Expires        string
	AuthorizedKeys string
	IPSecPSK       string
	Privileges     []string
	Groups         []string
	CustomSettings bool
	WebGUICss      string
	DashboardCols  string
	KeepHistory    bool
	controlID      int
}

func (u *User) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("%w, username is required", ErrClientValidation)
	}

	if len(name) > UserMaxNameLength {
		return fmt.Errorf("%w, username must be at most %d characters", ErrClientValidation, UserMaxNameLength)
	}

	if !userNameRegex.MatchString(name) {
		return fmt.Errorf("%w, username must contain only alphanumeric characters, dots, hyphens, and underscores", ErrClientValidation)
	}

	u.Name = name

	return nil
}

func (u *User) SetDescription(description string) error {
	u.Description = description

	return nil
}

func (u *User) SetDisabled(disabled bool) error {
	u.Disabled = disabled

	return nil
}

func (u *User) SetExpires(expires string) error {
	u.Expires = expires

	return nil
}

func (u *User) SetAuthorizedKeys(authorizedKeys string) error {
	u.AuthorizedKeys = authorizedKeys

	return nil
}

func (u *User) SetIPSecPSK(ipsecPSK string) error {
	if ipsecPSK != "" {
		for _, r := range ipsecPSK {
			if r > 127 { //nolint:mnd
				return fmt.Errorf("%w, IPsec pre-shared key must contain only ASCII characters", ErrClientValidation)
			}
		}
	}

	u.IPSecPSK = ipsecPSK

	return nil
}

func (u *User) SetPrivileges(privileges []string) error {
	u.Privileges = privileges

	return nil
}

func (u *User) SetGroups(groups []string) error {
	u.Groups = groups

	return nil
}

func (u *User) SetCustomSettings(customSettings bool) error {
	u.CustomSettings = customSettings

	return nil
}

func (u *User) SetWebGUICss(webguicss string) error {
	u.WebGUICss = webguicss

	return nil
}

func (u *User) SetDashboardCols(dashboardCols string) error {
	u.DashboardCols = dashboardCols

	return nil
}

func (u *User) SetKeepHistory(keepHistory bool) error {
	u.KeepHistory = keepHistory

	return nil
}

type Users []User

func (users Users) GetByName(name string) (*User, error) {
	for _, u := range users {
		if u.Name == name {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("user %w with name '%s'", ErrNotFound, name)
}

func (users Users) GetControlIDByName(name string) (*int, error) {
	for _, u := range users {
		if u.Name == name {
			return &u.controlID, nil
		}
	}

	return nil, fmt.Errorf("user %w with name '%s'", ErrNotFound, name)
}

func parseBoolField(v interface{}) bool {
	if v == nil {
		return false
	}

	switch val := v.(type) {
	case bool:
		return val
	case string:
		return val != ""
	default:
		return false
	}
}

func parseStringSliceField(v interface{}) []string {
	if v == nil {
		return nil
	}

	switch val := v.(type) {
	case []interface{}:
		result := make([]string, 0, len(val))
		for _, item := range val {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}

		return result
	case string:
		if val != "" {
			return []string{val}
		}

		return nil
	default:
		return nil
	}
}

func parseUserResponse(resp userResponse) (User, error) {
	var user User

	if err := user.SetName(resp.Name); err != nil {
		return user, err
	}

	if err := user.SetDescription(resp.Description); err != nil {
		return user, err
	}

	user.Scope = resp.Scope
	user.UID = resp.UID
	user.Disabled = parseBoolField(resp.Disabled)
	user.Expires = resp.Expires

	// Authorized keys are base64-encoded in config.
	if resp.AuthorizedKeys != "" {
		decoded, err := base64.StdEncoding.DecodeString(resp.AuthorizedKeys)
		if err != nil {
			// If decode fails, store as-is.
			user.AuthorizedKeys = resp.AuthorizedKeys
		} else {
			user.AuthorizedKeys = string(decoded)
		}
	}

	user.IPSecPSK = resp.IPSecPSK
	user.Privileges = parseStringSliceField(resp.Privileges)
	user.Groups = parseStringSliceField(resp.Groups)
	user.CustomSettings = parseBoolField(resp.CustomSettings)
	user.WebGUICss = resp.WebGUICss
	user.DashboardCols = resp.DashboardCols
	user.KeepHistory = parseBoolField(resp.KeepHistory)
	user.controlID = resp.ControlID

	return user, nil
}

func (pf *Client) getUsers(ctx context.Context) (*Users, error) {
	// Read users with their group memberships resolved.
	command := "$output = array();" +
		"$users = config_get_path('system/user', array());" +
		"$groups = config_get_path('system/group', array());" +
		"foreach ($users as $k => $v) {" +
		"$v['controlID'] = $k;" +
		// Normalize boolean fields: key presence means true in pfSense config.
		"$v['disabled'] = array_key_exists('disabled', $v);" +
		"$v['customsettings'] = array_key_exists('customsettings', $v);" +
		"$v['keephistory'] = array_key_exists('keephistory', $v);" +
		// Resolve group memberships for this user.
		"$user_groups = array();" +
		"foreach ($groups as $g) {" +
		"if (isset($g['member']) && is_array($g['member'])) {" +
		"if (in_array($v['uid'], $g['member'])) {" +
		"$user_groups[] = $g['name'];" +
		"}" +
		"} elseif (isset($g['member']) && $g['member'] == $v['uid']) {" +
		"$user_groups[] = $g['name'];" +
		"}" +
		// The 'all' group implicitly contains all users.
		"if ($g['name'] == 'all' && !in_array('all', $user_groups)) {" +
		"$user_groups[] = 'all';" +
		"}" +
		"};" +
		"$v['groups'] = $user_groups;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var userResp []userResponse
	if err := pf.executePHPCommand(ctx, command, &userResp); err != nil {
		return nil, err
	}

	users := make(Users, 0, len(userResp))
	for _, resp := range userResp {
		u, err := parseUserResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w user response, %w", ErrUnableToParse, err)
		}

		users = append(users, u)
	}

	return &users, nil
}

func (pf *Client) GetUsers(ctx context.Context) (*Users, error) {
	defer pf.read(&pf.mutexes.User)()

	users, err := pf.getUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w users, %w", ErrGetOperationFailed, err)
	}

	return users, nil
}

func (pf *Client) GetUser(ctx context.Context, name string) (*User, error) {
	defer pf.read(&pf.mutexes.User)()

	users, err := pf.getUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w users, %w", ErrGetOperationFailed, err)
	}

	u, err := users.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w user, %w", ErrGetOperationFailed, err)
	}

	return u, nil
}

func (pf *Client) CreateUser(ctx context.Context, req User, password string) (*User, error) {
	defer pf.write(&pf.mutexes.User)()

	// Check for duplicate username.
	existingUsers, err := pf.getUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w users for duplicate check, %w", ErrGetOperationFailed, err)
	}

	if _, err := existingUsers.GetByName(req.Name); err == nil {
		return nil, fmt.Errorf("%w user, a user with name '%s' already exists", ErrCreateOperationFailed, req.Name)
	}

	if password == "" {
		return nil, fmt.Errorf("%w user, password is required for new users", ErrCreateOperationFailed)
	}

	// Build the PHP command to create the user.
	var cmd strings.Builder

	cmd.WriteString("require_once('auth.inc');")

	// Build user array.
	cmd.WriteString("$user = array();")
	cmd.WriteString(fmt.Sprintf("$user['name'] = '%s';", phpEscape(req.Name)))
	cmd.WriteString(fmt.Sprintf("$user['descr'] = '%s';", phpEscape(req.Description)))
	cmd.WriteString("$user['scope'] = 'user';")

	// Assign UID from nextuid counter.
	cmd.WriteString("$nextuid = config_get_path('system/nextuid');")
	cmd.WriteString("$user['uid'] = $nextuid++;")
	cmd.WriteString("config_set_path('system/nextuid', $nextuid);")

	// Disabled.
	if req.Disabled {
		cmd.WriteString("$user['disabled'] = '';")
	}

	// Expires.
	if req.Expires != "" {
		cmd.WriteString(fmt.Sprintf("$user['expires'] = '%s';", phpEscape(req.Expires)))
	} else {
		cmd.WriteString("$user['expires'] = '';")
	}

	// Dashboard columns default.
	if req.DashboardCols != "" {
		cmd.WriteString(fmt.Sprintf("$user['dashboardcolumns'] = '%s';", phpEscape(req.DashboardCols)))
	} else {
		cmd.WriteString("$user['dashboardcolumns'] = '2';")
	}

	// Authorized keys (base64-encode).
	if req.AuthorizedKeys != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(req.AuthorizedKeys))
		cmd.WriteString(fmt.Sprintf("$user['authorizedkeys'] = '%s';", phpEscape(encoded)))
	} else {
		cmd.WriteString("$user['authorizedkeys'] = '';")
	}

	// IPsec PSK.
	if req.IPSecPSK != "" {
		cmd.WriteString(fmt.Sprintf("$user['ipsecpsk'] = '%s';", phpEscape(req.IPSecPSK)))
	} else {
		cmd.WriteString("$user['ipsecpsk'] = '';")
	}

	// WebGUI CSS.
	if req.WebGUICss != "" {
		cmd.WriteString(fmt.Sprintf("$user['webguicss'] = '%s';", phpEscape(req.WebGUICss)))
	} else {
		cmd.WriteString("$user['webguicss'] = 'pfSense.css';")
	}

	// Custom settings.
	if req.CustomSettings {
		cmd.WriteString("$user['customsettings'] = '';")
	}

	// Keep history.
	if req.KeepHistory {
		cmd.WriteString("$user['keephistory'] = '';")
	}

	// Privileges.
	if len(req.Privileges) > 0 {
		cmd.WriteString("$user['priv'] = array(")
		for i, priv := range req.Privileges {
			if i > 0 {
				cmd.WriteString(",")
			}
			cmd.WriteString(fmt.Sprintf("'%s'", phpEscape(priv)))
		}
		cmd.WriteString(");")
	}

	// Append user to config.
	cmd.WriteString("config_set_path('system/user/', $user);")

	// Set password (bcrypt hash).
	cmd.WriteString("$user_item_config = array();")
	cmd.WriteString("$users = config_get_path('system/user', array());")
	cmd.WriteString("$idx = count($users) - 1;")
	cmd.WriteString("$user_item_config['idx'] = $idx;")
	cmd.WriteString("$user_item_config['item'] = &$users[$idx];")
	cmd.WriteString(fmt.Sprintf("local_user_set_password($user_item_config, '%s');", phpEscape(password)))
	cmd.WriteString("config_set_path('system/user', $users);")

	// Sort users alphabetically by name.
	cmd.WriteString("$all_users = config_get_path('system/user', array());")
	cmd.WriteString("usort($all_users, function($a, $b) { return strcmp($a['name'], $b['name']); });")
	cmd.WriteString("config_set_path('system/user', $all_users);")

	// Add user to 'all' group.
	cmd.WriteString("$groups = config_get_path('system/group', array());")
	cmd.WriteString("foreach ($groups as $gidx => &$group) {")
	cmd.WriteString("if ($group['name'] == 'all') {")
	cmd.WriteString("if (!isset($group['member'])) { $group['member'] = array(); }")
	cmd.WriteString("if (!is_array($group['member'])) { $group['member'] = array($group['member']); }")
	cmd.WriteString("$group['member'][] = $user['uid'];")
	cmd.WriteString("break;")
	cmd.WriteString("}")
	cmd.WriteString("};")
	cmd.WriteString("config_set_path('system/group', $groups);")

	// Handle additional group memberships.
	if len(req.Groups) > 0 {
		cmd.WriteString("$desired_groups = array(")
		for i, g := range req.Groups {
			if i > 0 {
				cmd.WriteString(",")
			}
			cmd.WriteString(fmt.Sprintf("'%s'", phpEscape(g)))
		}
		cmd.WriteString(");")
		// Re-read user to get assigned UID.
		cmd.WriteString("$all_users = config_get_path('system/user', array());")
		cmd.WriteString("$created_user = null;")
		cmd.WriteString("foreach ($all_users as $u) {")
		cmd.WriteString(fmt.Sprintf("if ($u['name'] == '%s') { $created_user = $u; break; }", phpEscape(req.Name)))
		cmd.WriteString("};")
		cmd.WriteString("if ($created_user) { local_user_set_groups($created_user, $desired_groups); };")
	}

	// Write config.
	cmd.WriteString(fmt.Sprintf("write_config('Terraform: created user %s');", phpEscape(req.Name)))

	cmd.WriteString("print(json_encode(true));")

	var result bool
	if err := pf.executePHPCommand(ctx, cmd.String(), &result); err != nil {
		return nil, fmt.Errorf("%w user, %w", ErrCreateOperationFailed, err)
	}

	users, err := pf.getUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w users after creating, %w", ErrGetOperationFailed, err)
	}

	u, err := users.GetByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w user after creating, %w", ErrGetOperationFailed, err)
	}

	return u, nil
}

func (pf *Client) UpdateUser(ctx context.Context, req User, password string) (*User, error) {
	defer pf.write(&pf.mutexes.User)()

	users, err := pf.getUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w users, %w", ErrGetOperationFailed, err)
	}

	controlID, err := users.GetControlIDByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w user, %w", ErrGetOperationFailed, err)
	}

	var cmd strings.Builder

	cmd.WriteString("require_once('auth.inc');")

	// Load existing user to preserve read-only fields.
	cmd.WriteString(fmt.Sprintf("$user = config_get_path('system/user/%d');", *controlID))

	// Update mutable fields.
	cmd.WriteString(fmt.Sprintf("$user['descr'] = '%s';", phpEscape(req.Description)))

	// Disabled.
	if req.Disabled {
		cmd.WriteString("$user['disabled'] = '';")
	} else {
		cmd.WriteString("unset($user['disabled']);")
	}

	// Expires.
	if req.Expires != "" {
		cmd.WriteString(fmt.Sprintf("$user['expires'] = '%s';", phpEscape(req.Expires)))
	} else {
		cmd.WriteString("$user['expires'] = '';")
	}

	// Dashboard columns.
	if req.DashboardCols != "" {
		cmd.WriteString(fmt.Sprintf("$user['dashboardcolumns'] = '%s';", phpEscape(req.DashboardCols)))
	} else {
		cmd.WriteString("$user['dashboardcolumns'] = '2';")
	}

	// Authorized keys (base64-encode).
	if req.AuthorizedKeys != "" {
		encoded := base64.StdEncoding.EncodeToString([]byte(req.AuthorizedKeys))
		cmd.WriteString(fmt.Sprintf("$user['authorizedkeys'] = '%s';", phpEscape(encoded)))
	} else {
		cmd.WriteString("$user['authorizedkeys'] = '';")
	}

	// IPsec PSK.
	if req.IPSecPSK != "" {
		cmd.WriteString(fmt.Sprintf("$user['ipsecpsk'] = '%s';", phpEscape(req.IPSecPSK)))
	} else {
		cmd.WriteString("$user['ipsecpsk'] = '';")
	}

	// WebGUI CSS.
	if req.WebGUICss != "" {
		cmd.WriteString(fmt.Sprintf("$user['webguicss'] = '%s';", phpEscape(req.WebGUICss)))
	} else {
		cmd.WriteString("$user['webguicss'] = 'pfSense.css';")
	}

	// Custom settings.
	if req.CustomSettings {
		cmd.WriteString("$user['customsettings'] = '';")
	} else {
		cmd.WriteString("unset($user['customsettings']);")
	}

	// Keep history.
	if req.KeepHistory {
		cmd.WriteString("$user['keephistory'] = '';")
	} else {
		cmd.WriteString("unset($user['keephistory']);")
	}

	// Privileges.
	if len(req.Privileges) > 0 {
		cmd.WriteString("$user['priv'] = array(")
		for i, priv := range req.Privileges {
			if i > 0 {
				cmd.WriteString(",")
			}
			cmd.WriteString(fmt.Sprintf("'%s'", phpEscape(priv)))
		}
		cmd.WriteString(");")
	} else {
		cmd.WriteString("unset($user['priv']);")
	}

	// Write user back.
	cmd.WriteString(fmt.Sprintf("config_set_path('system/user/%d', $user);", *controlID))

	// Update password if provided.
	if password != "" {
		cmd.WriteString("$users = config_get_path('system/user', array());")
		// Find the current index (may have shifted).
		cmd.WriteString(fmt.Sprintf("$idx = null; foreach ($users as $k => $u) { if ($u['name'] == '%s') { $idx = $k; break; } };", phpEscape(req.Name)))
		cmd.WriteString("if ($idx !== null) {")
		cmd.WriteString("$user_item_config = array();")
		cmd.WriteString("$user_item_config['idx'] = $idx;")
		cmd.WriteString("$user_item_config['item'] = &$users[$idx];")
		cmd.WriteString(fmt.Sprintf("local_user_set_password($user_item_config, '%s');", phpEscape(password)))
		cmd.WriteString("config_set_path('system/user', $users);")
		cmd.WriteString("};")
	}

	// Sort users alphabetically by name.
	cmd.WriteString("$all_users = config_get_path('system/user', array());")
	cmd.WriteString("usort($all_users, function($a, $b) { return strcmp($a['name'], $b['name']); });")
	cmd.WriteString("config_set_path('system/user', $all_users);")

	// Update group memberships.
	cmd.WriteString("$all_users = config_get_path('system/user', array());")
	cmd.WriteString("$current_user = null;")
	cmd.WriteString("foreach ($all_users as $u) {")
	cmd.WriteString(fmt.Sprintf("if ($u['name'] == '%s') { $current_user = $u; break; }", phpEscape(req.Name)))
	cmd.WriteString("};")
	cmd.WriteString("if ($current_user) {")

	if len(req.Groups) > 0 {
		cmd.WriteString("$desired_groups = array(")
		for i, g := range req.Groups {
			if i > 0 {
				cmd.WriteString(",")
			}
			cmd.WriteString(fmt.Sprintf("'%s'", phpEscape(g)))
		}
		cmd.WriteString(");")
		cmd.WriteString("local_user_set_groups($current_user, $desired_groups);")
	} else {
		cmd.WriteString("local_user_set_groups($current_user, array());")
	}

	cmd.WriteString("};")

	// Write config.
	cmd.WriteString(fmt.Sprintf("write_config('Terraform: updated user %s');", phpEscape(req.Name)))

	cmd.WriteString("print(json_encode(true));")

	var result bool
	if err := pf.executePHPCommand(ctx, cmd.String(), &result); err != nil {
		return nil, fmt.Errorf("%w user, %w", ErrUpdateOperationFailed, err)
	}

	users, err = pf.getUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w users after updating, %w", ErrGetOperationFailed, err)
	}

	u, err := users.GetByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w user after updating, %w", ErrGetOperationFailed, err)
	}

	return u, nil
}

func (pf *Client) DeleteUser(ctx context.Context, name string) error {
	defer pf.write(&pf.mutexes.User)()

	users, err := pf.getUsers(ctx)
	if err != nil {
		return fmt.Errorf("%w users, %w", ErrGetOperationFailed, err)
	}

	controlID, err := users.GetControlIDByName(name)
	if err != nil {
		return fmt.Errorf("%w user, %w", ErrGetOperationFailed, err)
	}

	// Verify not a system user.
	user, _ := users.GetByName(name)
	if user != nil && user.Scope == "system" {
		return fmt.Errorf("%w user, cannot delete system user '%s'", ErrDeleteOperationFailed, name)
	}

	var cmd strings.Builder

	cmd.WriteString("require_once('auth.inc');")

	// Delete OS-level user.
	cmd.WriteString(fmt.Sprintf("local_user_del(config_get_path('system/user/%d'));", *controlID))

	// Remove from config.
	cmd.WriteString(fmt.Sprintf("config_del_path('system/user/%d');", *controlID))

	// Reindex array.
	cmd.WriteString("config_set_path('system/user', array_values(config_get_path('system/user', array())));")

	// Write config.
	cmd.WriteString(fmt.Sprintf("write_config('Terraform: deleted user %s');", phpEscape(name)))

	cmd.WriteString("print(json_encode(true));")

	var result bool
	if err := pf.executePHPCommand(ctx, cmd.String(), &result); err != nil {
		return fmt.Errorf("%w user, %w", ErrDeleteOperationFailed, err)
	}

	// Verify deletion.
	users, err = pf.getUsers(ctx)
	if err != nil {
		return fmt.Errorf("%w users after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := users.GetByName(name); err == nil {
		return fmt.Errorf("%w user, still exists", ErrDeleteOperationFailed)
	}

	return nil
}

func (pf *Client) ApplyUserChanges(ctx context.Context, name string) error {
	pf.mutexes.UserApply.Lock()
	defer pf.mutexes.UserApply.Unlock()

	var cmd strings.Builder

	cmd.WriteString("require_once('auth.inc');")

	// Find user by name and apply.
	cmd.WriteString("$users = config_get_path('system/user', array());")
	cmd.WriteString("$user = null;")
	cmd.WriteString("foreach ($users as $u) {")
	cmd.WriteString(fmt.Sprintf("if ($u['name'] == '%s') { $user = $u; break; }", phpEscape(name)))
	cmd.WriteString("};")
	cmd.WriteString("if ($user) {")
	cmd.WriteString("local_user_set($user);")
	cmd.WriteString("};")

	cmd.WriteString("print(json_encode(true));")

	var result bool
	if err := pf.executePHPCommand(ctx, cmd.String(), &result); err != nil {
		return fmt.Errorf("%w, failed to apply user changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}
