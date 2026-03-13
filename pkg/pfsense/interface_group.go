package pfsense

import (
	"context"
	"fmt"
	"strings"
)

type interfaceGroupResponse struct {
	Name        string `json:"ifname"`
	Members     string `json:"members"`
	Description string `json:"descr"`
	ControlID   int    `json:"controlID"` //nolint:tagliatelle
}

type InterfaceGroup struct {
	Name        string
	Members     []string
	Description string
	controlID   int
}

func (g *InterfaceGroup) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("%w, interface group name is required", ErrClientValidation)
	}

	g.Name = name

	return nil
}

func (g *InterfaceGroup) SetMembers(members []string) error {
	g.Members = members

	return nil
}

func (g *InterfaceGroup) SetDescription(description string) error {
	g.Description = description

	return nil
}

type InterfaceGroups []InterfaceGroup

func (groups InterfaceGroups) GetByName(name string) (*InterfaceGroup, error) {
	for _, g := range groups {
		if g.Name == name {
			return &g, nil
		}
	}

	return nil, fmt.Errorf("interface group %w with name '%s'", ErrNotFound, name)
}

func (groups InterfaceGroups) GetControlIDByName(name string) (*int, error) {
	for _, g := range groups {
		if g.Name == name {
			return &g.controlID, nil
		}
	}

	return nil, fmt.Errorf("interface group %w with name '%s'", ErrNotFound, name)
}

func parseInterfaceGroupResponse(resp interfaceGroupResponse) (InterfaceGroup, error) {
	var group InterfaceGroup

	if err := group.SetName(resp.Name); err != nil {
		return group, err
	}

	if resp.Members != "" {
		members := strings.Fields(resp.Members)
		if err := group.SetMembers(members); err != nil {
			return group, err
		}
	}

	if err := group.SetDescription(resp.Description); err != nil {
		return group, err
	}

	group.controlID = resp.ControlID

	return group, nil
}

func (pf *Client) getInterfaceGroups(ctx context.Context) (*InterfaceGroups, error) {
	command := "$output = array();" +
		"$groups = config_get_path('ifgroups/ifgroupentry', array());" +
		"if (!is_array($groups)) { $groups = array(); }" +
		"foreach ($groups as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var groupResp []interfaceGroupResponse
	if err := pf.executePHPCommand(ctx, command, &groupResp); err != nil {
		return nil, err
	}

	groups := make(InterfaceGroups, 0, len(groupResp))
	for _, resp := range groupResp {
		g, err := parseInterfaceGroupResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w interface group response, %w", ErrUnableToParse, err)
		}

		groups = append(groups, g)
	}

	return &groups, nil
}

func (pf *Client) GetInterfaceGroups(ctx context.Context) (*InterfaceGroups, error) {
	defer pf.read(&pf.mutexes.InterfaceGroup)()

	groups, err := pf.getInterfaceGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w interface groups, %w", ErrGetOperationFailed, err)
	}

	return groups, nil
}

func (pf *Client) GetInterfaceGroup(ctx context.Context, name string) (*InterfaceGroup, error) {
	defer pf.read(&pf.mutexes.InterfaceGroup)()

	groups, err := pf.getInterfaceGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w interface groups, %w", ErrGetOperationFailed, err)
	}

	g, err := groups.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w interface group, %w", ErrGetOperationFailed, err)
	}

	return g, nil
}

func (pf *Client) CreateInterfaceGroup(ctx context.Context, req InterfaceGroup) (*InterfaceGroup, error) {
	defer pf.write(&pf.mutexes.InterfaceGroup)()

	membersStr := strings.Join(req.Members, " ")

	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"$ifgroupentry = array();"+
			"$ifgroupentry['ifname'] = '%s';"+
			"$ifgroupentry['members'] = '%s';"+
			"$ifgroupentry['descr'] = '%s';"+
			"$groups = config_get_path('ifgroups/ifgroupentry', array());"+
			"if (!is_array($groups)) { config_set_path('ifgroups/ifgroupentry', array()); }"+
			"config_set_path('ifgroups/ifgroupentry/', $ifgroupentry);"+
			"interface_group_setup($ifgroupentry);"+
			"write_config('Terraform: created interface group %s');"+
			"print(json_encode(true));",
		phpEscape(req.Name),
		phpEscape(membersStr),
		phpEscape(req.Description),
		phpEscape(req.Name),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w interface group, %w", ErrCreateOperationFailed, err)
	}

	groups, err := pf.getInterfaceGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w interface groups after creating, %w", ErrGetOperationFailed, err)
	}

	g, err := groups.GetByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w interface group after creating, %w", ErrGetOperationFailed, err)
	}

	return g, nil
}

func (pf *Client) UpdateInterfaceGroup(ctx context.Context, req InterfaceGroup) (*InterfaceGroup, error) {
	defer pf.write(&pf.mutexes.InterfaceGroup)()

	groups, err := pf.getInterfaceGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w interface groups, %w", ErrGetOperationFailed, err)
	}

	controlID, err := groups.GetControlIDByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w interface group, %w", ErrGetOperationFailed, err)
	}

	membersStr := strings.Join(req.Members, " ")

	// Remove old group from OS interfaces, then re-setup with new config.
	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"$old = config_get_path('ifgroups/ifgroupentry/%d', array());"+
			"if (!empty($old['members'])) {"+
			"foreach (explode(' ', $old['members']) as $intf) {"+
			"$realif = get_real_interface($intf);"+
			"if (!empty($realif)) { mwexec(\"/sbin/ifconfig \" . escapeshellarg($realif) . \" -group \" . escapeshellarg($old['ifname'])); }"+
			"}}"+
			"$ifgroupentry = array();"+
			"$ifgroupentry['ifname'] = '%s';"+
			"$ifgroupentry['members'] = '%s';"+
			"$ifgroupentry['descr'] = '%s';"+
			"config_set_path('ifgroups/ifgroupentry/%d', $ifgroupentry);"+
			"interface_group_setup($ifgroupentry);"+
			"write_config('Terraform: updated interface group %s');"+
			"print(json_encode(true));",
		*controlID,
		phpEscape(req.Name),
		phpEscape(membersStr),
		phpEscape(req.Description),
		*controlID,
		phpEscape(req.Name),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w interface group, %w", ErrUpdateOperationFailed, err)
	}

	groups, err = pf.getInterfaceGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w interface groups after updating, %w", ErrGetOperationFailed, err)
	}

	g, err := groups.GetByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w interface group after updating, %w", ErrGetOperationFailed, err)
	}

	return g, nil
}

func (pf *Client) DeleteInterfaceGroup(ctx context.Context, name string) error {
	defer pf.write(&pf.mutexes.InterfaceGroup)()

	groups, err := pf.getInterfaceGroups(ctx)
	if err != nil {
		return fmt.Errorf("%w interface groups, %w", ErrGetOperationFailed, err)
	}

	controlID, err := groups.GetControlIDByName(name)
	if err != nil {
		return fmt.Errorf("%w interface group, %w", ErrGetOperationFailed, err)
	}

	// Remove group from all member interfaces via ifconfig, then delete from config.
	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"$group = config_get_path('ifgroups/ifgroupentry/%d', array());"+
			"if (!empty($group['members'])) {"+
			"foreach (explode(' ', $group['members']) as $intf) {"+
			"$realif = get_real_interface($intf);"+
			"if (!empty($realif)) { mwexec(\"/sbin/ifconfig \" . escapeshellarg($realif) . \" -group \" . escapeshellarg($group['ifname'])); }"+
			"}}"+
			"config_del_path('ifgroups/ifgroupentry/%d');"+
			"write_config('Terraform: deleted interface group %s');"+
			"print(json_encode(true));",
		*controlID,
		*controlID,
		phpEscape(name),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w interface group, %w", ErrDeleteOperationFailed, err)
	}

	// Verify deletion.
	groups, err = pf.getInterfaceGroups(ctx)
	if err != nil {
		return fmt.Errorf("%w interface groups after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := groups.GetByName(name); err == nil {
		return fmt.Errorf("%w interface group, still exists", ErrDeleteOperationFailed)
	}

	return nil
}
