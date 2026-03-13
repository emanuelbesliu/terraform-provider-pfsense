package pfsense

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var groupNameRegex = regexp.MustCompile(`^[a-zA-Z0-9.\-_]+$`)

const (
	GroupMaxNameLength = 16
)

type groupResponse struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Scope       string      `json:"scope"`
	GID         string      `json:"gid"`
	MemberUIDs  interface{} `json:"member"`
	MemberNames interface{} `json:"member_names"`
	Privileges  interface{} `json:"priv"`
	ControlID   int         `json:"controlID"` //nolint:tagliatelle
}

type Group struct {
	Name        string
	Description string
	Scope       string
	GID         string
	Members     []string // Usernames (resolved from UIDs for reads, usernames for writes).
	Privileges  []string
	controlID   int
}

func (g *Group) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("%w, group name is required", ErrClientValidation)
	}

	if len(name) > GroupMaxNameLength {
		return fmt.Errorf("%w, group name must be at most %d characters", ErrClientValidation, GroupMaxNameLength)
	}

	if !groupNameRegex.MatchString(name) {
		return fmt.Errorf("%w, group name must contain only alphanumeric characters, dots, hyphens, and underscores", ErrClientValidation)
	}

	g.Name = name

	return nil
}

func (g *Group) SetDescription(description string) error {
	g.Description = description

	return nil
}

func (g *Group) SetMembers(members []string) error {
	g.Members = members

	return nil
}

func (g *Group) SetPrivileges(privileges []string) error {
	g.Privileges = privileges

	return nil
}

type Groups []Group

func (groups Groups) GetByName(name string) (*Group, error) {
	for _, g := range groups {
		if g.Name == name {
			return &g, nil
		}
	}

	return nil, fmt.Errorf("group %w with name '%s'", ErrNotFound, name)
}

func (groups Groups) GetControlIDByName(name string) (*int, error) {
	for _, g := range groups {
		if g.Name == name {
			return &g.controlID, nil
		}
	}

	return nil, fmt.Errorf("group %w with name '%s'", ErrNotFound, name)
}

func parseGroupResponse(resp groupResponse) (Group, error) {
	var group Group

	if err := group.SetName(resp.Name); err != nil {
		return group, err
	}

	if err := group.SetDescription(resp.Description); err != nil {
		return group, err
	}

	group.Scope = resp.Scope
	group.GID = resp.GID

	// Use resolved member names (usernames) instead of raw UIDs.
	group.Members = parseStringSliceField(resp.MemberNames)
	group.Privileges = parseStringSliceField(resp.Privileges)
	group.controlID = resp.ControlID

	return group, nil
}

func (pf *Client) getGroups(ctx context.Context) (*Groups, error) {
	// Read groups and resolve member UIDs to usernames.
	command := "$output = array();" +
		"$groups = config_get_path('system/group', array());" +
		"$users = config_get_path('system/user', array());" +
		"foreach ($groups as $k => $v) {" +
		"$v['controlID'] = $k;" +
		// Normalize: pfSense stores 'descr' but we want 'description' in JSON.
		"$v['description'] = isset($v['descr']) ? $v['descr'] : '';" +
		// Ensure member is always an array.
		"if (!isset($v['member'])) { $v['member'] = array(); }" +
		"elseif (!is_array($v['member'])) { $v['member'] = array($v['member']); }" +
		// Resolve UIDs to usernames.
		"$member_names = array();" +
		"foreach ($v['member'] as $uid) {" +
		"foreach ($users as $u) {" +
		"if ($u['uid'] == $uid) { $member_names[] = $u['name']; break; }" +
		"};" +
		"};" +
		"$v['member_names'] = $member_names;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var groupResp []groupResponse
	if err := pf.executePHPCommand(ctx, command, &groupResp); err != nil {
		return nil, err
	}

	groups := make(Groups, 0, len(groupResp))
	for _, resp := range groupResp {
		g, err := parseGroupResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w group response, %w", ErrUnableToParse, err)
		}

		groups = append(groups, g)
	}

	return &groups, nil
}

func (pf *Client) GetGroups(ctx context.Context) (*Groups, error) {
	defer pf.read(&pf.mutexes.Group)()

	groups, err := pf.getGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w groups, %w", ErrGetOperationFailed, err)
	}

	return groups, nil
}

func (pf *Client) GetGroup(ctx context.Context, name string) (*Group, error) {
	defer pf.read(&pf.mutexes.Group)()

	groups, err := pf.getGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w groups, %w", ErrGetOperationFailed, err)
	}

	g, err := groups.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w group, %w", ErrGetOperationFailed, err)
	}

	return g, nil
}

func (pf *Client) CreateGroup(ctx context.Context, req Group) (*Group, error) {
	defer pf.write(&pf.mutexes.Group)()

	// Check for duplicate group name.
	existingGroups, err := pf.getGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w groups for duplicate check, %w", ErrGetOperationFailed, err)
	}

	if _, err := existingGroups.GetByName(req.Name); err == nil {
		return nil, fmt.Errorf("%w group, a group with name '%s' already exists", ErrCreateOperationFailed, req.Name)
	}

	var cmd strings.Builder

	cmd.WriteString("require_once('auth.inc');")

	// Build group array.
	cmd.WriteString("$group = array();")
	cmd.WriteString(fmt.Sprintf("$group['name'] = '%s';", phpEscape(req.Name)))
	cmd.WriteString(fmt.Sprintf("$group['descr'] = '%s';", phpEscape(req.Description)))
	cmd.WriteString("$group['scope'] = 'local';")

	// Assign GID from nextgid counter.
	cmd.WriteString("$nextgid = config_get_path('system/nextgid');")
	cmd.WriteString("$group['gid'] = $nextgid++;")
	cmd.WriteString("config_set_path('system/nextgid', $nextgid);")

	// Members: resolve usernames to UIDs via PHP.
	if len(req.Members) > 0 {
		cmd.WriteString("$users = config_get_path('system/user', array());")
		cmd.WriteString("$member_uids = array();")

		for _, username := range req.Members {
			cmd.WriteString(fmt.Sprintf(
				"foreach ($users as $u) { if ($u['name'] == '%s') { $member_uids[] = $u['uid']; break; } };",
				phpEscape(username),
			))
		}

		cmd.WriteString("$group['member'] = $member_uids;")
	}

	// Privileges.
	if len(req.Privileges) > 0 {
		cmd.WriteString("$group['priv'] = array(")
		for i, priv := range req.Privileges {
			if i > 0 {
				cmd.WriteString(",")
			}
			cmd.WriteString(fmt.Sprintf("'%s'", phpEscape(priv)))
		}
		cmd.WriteString(");")
	}

	// Append group to config.
	cmd.WriteString("config_set_path('system/group/', $group);")

	// Write config.
	cmd.WriteString(fmt.Sprintf("write_config('Terraform: created group %s');", phpEscape(req.Name)))

	cmd.WriteString("print(json_encode(true));")

	var result bool
	if err := pf.executePHPCommand(ctx, cmd.String(), &result); err != nil {
		return nil, fmt.Errorf("%w group, %w", ErrCreateOperationFailed, err)
	}

	groups, err := pf.getGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w groups after creating, %w", ErrGetOperationFailed, err)
	}

	g, err := groups.GetByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w group after creating, %w", ErrGetOperationFailed, err)
	}

	return g, nil
}

func (pf *Client) UpdateGroup(ctx context.Context, req Group) (*Group, error) {
	defer pf.write(&pf.mutexes.Group)()

	groups, err := pf.getGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w groups, %w", ErrGetOperationFailed, err)
	}

	controlID, err := groups.GetControlIDByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w group, %w", ErrGetOperationFailed, err)
	}

	var cmd strings.Builder

	cmd.WriteString("require_once('auth.inc');")

	// Load existing group to preserve read-only fields.
	cmd.WriteString(fmt.Sprintf("$group = config_get_path('system/group/%d');", *controlID))

	// Update mutable fields.
	cmd.WriteString(fmt.Sprintf("$group['descr'] = '%s';", phpEscape(req.Description)))

	// Members: resolve usernames to UIDs via PHP.
	if len(req.Members) > 0 {
		cmd.WriteString("$users = config_get_path('system/user', array());")
		cmd.WriteString("$member_uids = array();")

		for _, username := range req.Members {
			cmd.WriteString(fmt.Sprintf(
				"foreach ($users as $u) { if ($u['name'] == '%s') { $member_uids[] = $u['uid']; break; } };",
				phpEscape(username),
			))
		}

		cmd.WriteString("$group['member'] = $member_uids;")
	} else {
		cmd.WriteString("unset($group['member']);")
	}

	// Privileges.
	if len(req.Privileges) > 0 {
		cmd.WriteString("$group['priv'] = array(")
		for i, priv := range req.Privileges {
			if i > 0 {
				cmd.WriteString(",")
			}
			cmd.WriteString(fmt.Sprintf("'%s'", phpEscape(priv)))
		}
		cmd.WriteString(");")
	} else {
		cmd.WriteString("unset($group['priv']);")
	}

	// Write group back.
	cmd.WriteString(fmt.Sprintf("config_set_path('system/group/%d', $group);", *controlID))

	// Write config.
	cmd.WriteString(fmt.Sprintf("write_config('Terraform: updated group %s');", phpEscape(req.Name)))

	cmd.WriteString("print(json_encode(true));")

	var result bool
	if err := pf.executePHPCommand(ctx, cmd.String(), &result); err != nil {
		return nil, fmt.Errorf("%w group, %w", ErrUpdateOperationFailed, err)
	}

	groups, err = pf.getGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w groups after updating, %w", ErrGetOperationFailed, err)
	}

	g, err := groups.GetByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w group after updating, %w", ErrGetOperationFailed, err)
	}

	return g, nil
}

func (pf *Client) DeleteGroup(ctx context.Context, name string) error {
	defer pf.write(&pf.mutexes.Group)()

	groups, err := pf.getGroups(ctx)
	if err != nil {
		return fmt.Errorf("%w groups, %w", ErrGetOperationFailed, err)
	}

	controlID, err := groups.GetControlIDByName(name)
	if err != nil {
		return fmt.Errorf("%w group, %w", ErrGetOperationFailed, err)
	}

	// Verify not a system group.
	group, _ := groups.GetByName(name)
	if group != nil && group.Scope == "system" {
		return fmt.Errorf("%w group, cannot delete system group '%s'", ErrDeleteOperationFailed, name)
	}

	var cmd strings.Builder

	cmd.WriteString("require_once('auth.inc');")

	// Delete OS-level group.
	cmd.WriteString(fmt.Sprintf("local_group_del(config_get_path('system/group/%d'));", *controlID))

	// Remove from config.
	cmd.WriteString(fmt.Sprintf("config_del_path('system/group/%d');", *controlID))

	// Reindex array.
	cmd.WriteString("config_set_path('system/group', array_values(config_get_path('system/group', array())));")

	// Write config.
	cmd.WriteString(fmt.Sprintf("write_config('Terraform: deleted group %s');", phpEscape(name)))

	cmd.WriteString("print(json_encode(true));")

	var result bool
	if err := pf.executePHPCommand(ctx, cmd.String(), &result); err != nil {
		return fmt.Errorf("%w group, %w", ErrDeleteOperationFailed, err)
	}

	// Verify deletion.
	groups, err = pf.getGroups(ctx)
	if err != nil {
		return fmt.Errorf("%w groups after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := groups.GetByName(name); err == nil {
		return fmt.Errorf("%w group, still exists", ErrDeleteOperationFailed)
	}

	return nil
}

func (pf *Client) ApplyGroupChanges(ctx context.Context, name string) error {
	pf.mutexes.GroupApply.Lock()
	defer pf.mutexes.GroupApply.Unlock()

	var cmd strings.Builder

	cmd.WriteString("require_once('auth.inc');")

	// Find group by name and apply.
	cmd.WriteString("$groups = config_get_path('system/group', array());")
	cmd.WriteString("$group = null;")
	cmd.WriteString("foreach ($groups as $g) {")
	cmd.WriteString(fmt.Sprintf("if ($g['name'] == '%s') { $group = $g; break; }", phpEscape(name)))
	cmd.WriteString("};")
	cmd.WriteString("if ($group) {")
	// Apply group to OS (pw groupadd/groupmod).
	cmd.WriteString("local_group_set($group);")
	// Refresh member users so their group membership is synced at OS level.
	cmd.WriteString("if (isset($group['member']) && is_array($group['member'])) {")
	cmd.WriteString("$users = config_get_path('system/user', array());")
	cmd.WriteString("foreach ($group['member'] as $uid) {")
	cmd.WriteString("foreach ($users as $u) {")
	cmd.WriteString("if ($u['uid'] == $uid) { local_user_set($u); break; }")
	cmd.WriteString("};")
	cmd.WriteString("};")
	cmd.WriteString("};")
	cmd.WriteString("};")

	cmd.WriteString("print(json_encode(true));")

	var result bool
	if err := pf.executePHPCommand(ctx, cmd.String(), &result); err != nil {
		return fmt.Errorf("%w, failed to apply group changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}
