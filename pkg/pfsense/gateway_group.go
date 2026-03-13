package pfsense

import (
	"context"
	"fmt"
	"regexp"
	"strings"
)

var gatewayGroupNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

const (
	MinGatewayGroupTier = 1
	MaxGatewayGroupTier = 5
)

type gatewayGroupResponse struct {
	Name               string   `json:"name"`
	Description        string   `json:"descr"`
	Trigger            string   `json:"trigger"`
	KeepFailoverStates string   `json:"keep_failover_states"`
	Items              []string `json:"item"`
	ControlID          int      `json:"controlID"` //nolint:tagliatelle
}

type GatewayGroupMember struct {
	Gateway   string
	Tier      int
	VirtualIP string
}

type GatewayGroup struct {
	Name               string
	Description        string
	Trigger            string
	KeepFailoverStates string
	Members            []GatewayGroupMember
	controlID          int
}

func (GatewayGroup) Triggers() []string {
	return []string{"down", "downloss", "downlatency", "downlosslatency"}
}

func (GatewayGroup) KeepFailoverStatesOptions() []string {
	return []string{"", "keep", "kill"}
}

func (gw *GatewayGroup) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("%w, name must not be empty", ErrClientValidation)
	}

	if len(name) > 31 {
		return fmt.Errorf("%w, name must be at most 31 characters", ErrClientValidation)
	}

	if !gatewayGroupNameRegex.MatchString(name) {
		return fmt.Errorf("%w, name must start with a letter or underscore and contain only alphanumeric characters and underscores", ErrClientValidation)
	}

	gw.Name = name

	return nil
}

func (gw *GatewayGroup) SetDescription(description string) error {
	if len(description) > 200 {
		return fmt.Errorf("%w, description must be at most 200 characters", ErrClientValidation)
	}

	gw.Description = description

	return nil
}

func (gw *GatewayGroup) SetTrigger(trigger string) error {
	for _, t := range gw.Triggers() {
		if trigger == t {
			gw.Trigger = trigger

			return nil
		}
	}

	return fmt.Errorf("%w, trigger must be one of: %s", ErrClientValidation, strings.Join(gw.Triggers(), ", "))
}

func (gw *GatewayGroup) SetKeepFailoverStates(kfs string) error {
	for _, o := range gw.KeepFailoverStatesOptions() {
		if kfs == o {
			gw.KeepFailoverStates = kfs

			return nil
		}
	}

	return fmt.Errorf("%w, keep_failover_states must be one of: %s", ErrClientValidation,
		strings.Join(gw.KeepFailoverStatesOptions(), ", "))
}

func (gw *GatewayGroup) SetMembers(members []GatewayGroupMember) error {
	if len(members) == 0 {
		return fmt.Errorf("%w, at least one gateway member is required", ErrClientValidation)
	}

	for _, m := range members {
		if m.Tier < MinGatewayGroupTier || m.Tier > MaxGatewayGroupTier {
			return fmt.Errorf("%w, tier must be between %d and %d", ErrClientValidation, MinGatewayGroupTier, MaxGatewayGroupTier)
		}
	}

	gw.Members = members

	return nil
}

type GatewayGroups []GatewayGroup

func (groups GatewayGroups) GetByName(name string) (*GatewayGroup, error) {
	for _, g := range groups {
		if g.Name == name {
			return &g, nil
		}
	}

	return nil, fmt.Errorf("gateway group %w with name '%s'", ErrNotFound, name)
}

func (groups GatewayGroups) GetControlIDByName(name string) (*int, error) {
	for _, g := range groups {
		if g.Name == name {
			return &g.controlID, nil
		}
	}

	return nil, fmt.Errorf("gateway group %w with name '%s'", ErrNotFound, name)
}

func parseGatewayGroupMember(item string) (GatewayGroupMember, error) {
	parts := strings.SplitN(item, "|", 3) //nolint:mnd

	if len(parts) < 2 { //nolint:mnd
		return GatewayGroupMember{}, fmt.Errorf("%w, invalid gateway group member format: '%s'", ErrUnableToParse, item)
	}

	tier := 0

	if parts[1] != "" {
		t := 0
		if _, err := fmt.Sscanf(parts[1], "%d", &t); err != nil {
			return GatewayGroupMember{}, fmt.Errorf("%w, unable to parse tier from '%s'", ErrUnableToParse, item)
		}

		tier = t
	}

	vip := ""
	if len(parts) == 3 { //nolint:mnd
		vip = parts[2]
	}

	return GatewayGroupMember{
		Gateway:   parts[0],
		Tier:      tier,
		VirtualIP: vip,
	}, nil
}

func parseGatewayGroupResponse(resp gatewayGroupResponse) (GatewayGroup, error) {
	var group GatewayGroup

	if err := group.SetName(resp.Name); err != nil {
		return group, err
	}

	if err := group.SetDescription(resp.Description); err != nil {
		return group, err
	}

	if resp.Trigger != "" {
		if err := group.SetTrigger(resp.Trigger); err != nil {
			return group, err
		}
	}

	if resp.KeepFailoverStates != "" {
		if err := group.SetKeepFailoverStates(resp.KeepFailoverStates); err != nil {
			return group, err
		}
	}

	members := make([]GatewayGroupMember, 0, len(resp.Items))
	for _, item := range resp.Items {
		m, err := parseGatewayGroupMember(item)
		if err != nil {
			return group, err
		}

		members = append(members, m)
	}

	group.Members = members
	group.controlID = resp.ControlID

	return group, nil
}

func (pf *Client) getGatewayGroups(ctx context.Context) (*GatewayGroups, error) {
	command := "$output = array();" +
		"$groups = config_get_path('gateways/gateway_group', array());" +
		"foreach ($groups as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"if (!is_array($v['item'])) { $v['item'] = !empty($v['item']) ? array($v['item']) : array(); }" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var groupResp []gatewayGroupResponse
	if err := pf.executePHPCommand(ctx, command, &groupResp); err != nil {
		return nil, err
	}

	groups := make(GatewayGroups, 0, len(groupResp))
	for _, resp := range groupResp {
		g, err := parseGatewayGroupResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w gateway group response, %w", ErrUnableToParse, err)
		}

		groups = append(groups, g)
	}

	return &groups, nil
}

func (pf *Client) GetGatewayGroups(ctx context.Context) (*GatewayGroups, error) {
	defer pf.read(&pf.mutexes.GatewayGroup)()

	groups, err := pf.getGatewayGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w gateway groups, %w", ErrGetOperationFailed, err)
	}

	return groups, nil
}

func (pf *Client) GetGatewayGroup(ctx context.Context, name string) (*GatewayGroup, error) {
	defer pf.read(&pf.mutexes.GatewayGroup)()

	groups, err := pf.getGatewayGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w gateway groups, %w", ErrGetOperationFailed, err)
	}

	g, err := groups.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w gateway group, %w", ErrGetOperationFailed, err)
	}

	return g, nil
}

// gatewayGroupItemStrings converts members to pfSense's pipe-delimited format.
func gatewayGroupItemStrings(members []GatewayGroupMember) []string {
	items := make([]string, 0, len(members))
	for _, m := range members {
		items = append(items, fmt.Sprintf("%s|%d|%s", m.Gateway, m.Tier, m.VirtualIP))
	}

	return items
}

func (pf *Client) CreateGatewayGroup(ctx context.Context, req GatewayGroup) (*GatewayGroup, error) {
	defer pf.write(&pf.mutexes.GatewayGroup)()

	items := gatewayGroupItemStrings(req.Members)

	// Build PHP array representation of items.
	phpItems := "array("
	for i, item := range items {
		if i > 0 {
			phpItems += ","
		}

		phpItems += fmt.Sprintf("'%s'", item)
	}

	phpItems += ")"

	kfs := ""
	if req.KeepFailoverStates != "" {
		kfs = fmt.Sprintf("$group['keep_failover_states'] = '%s';", req.KeepFailoverStates)
	}

	command := fmt.Sprintf(
		"$group = array();"+
			"$group['name'] = '%s';"+
			"$group['item'] = %s;"+
			"$group['trigger'] = '%s';"+
			"%s"+
			"$group['descr'] = '%s';"+
			"config_set_path('gateways/gateway_group/', $group);"+
			"mark_subsystem_dirty('staticroutes');"+
			"mark_subsystem_dirty('gwgroup.%s');"+
			"write_config('Terraform: created gateway group %s');"+
			"print(json_encode(true));",
		phpEscape(req.Name),
		phpItems,
		phpEscape(req.Trigger),
		kfs,
		phpEscape(req.Description),
		phpEscape(req.Name),
		phpEscape(req.Name),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w gateway group, %w", ErrCreateOperationFailed, err)
	}

	groups, err := pf.getGatewayGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w gateway groups after creating, %w", ErrGetOperationFailed, err)
	}

	g, err := groups.GetByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w gateway group after creating, %w", ErrGetOperationFailed, err)
	}

	return g, nil
}

func (pf *Client) UpdateGatewayGroup(ctx context.Context, req GatewayGroup) (*GatewayGroup, error) {
	defer pf.write(&pf.mutexes.GatewayGroup)()

	groups, err := pf.getGatewayGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w gateway groups, %w", ErrGetOperationFailed, err)
	}

	controlID, err := groups.GetControlIDByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w gateway group, %w", ErrGetOperationFailed, err)
	}

	items := gatewayGroupItemStrings(req.Members)

	phpItems := "array("
	for i, item := range items {
		if i > 0 {
			phpItems += ","
		}

		phpItems += fmt.Sprintf("'%s'", item)
	}

	phpItems += ")"

	kfs := "$group['keep_failover_states'] = '';"
	if req.KeepFailoverStates != "" {
		kfs = fmt.Sprintf("$group['keep_failover_states'] = '%s';", req.KeepFailoverStates)
	}

	command := fmt.Sprintf(
		"$group = array();"+
			"$group['name'] = '%s';"+
			"$group['item'] = %s;"+
			"$group['trigger'] = '%s';"+
			"%s"+
			"$group['descr'] = '%s';"+
			"config_set_path('gateways/gateway_group/%d', $group);"+
			"mark_subsystem_dirty('staticroutes');"+
			"mark_subsystem_dirty('gwgroup.%s');"+
			"write_config('Terraform: updated gateway group %s');"+
			"print(json_encode(true));",
		phpEscape(req.Name),
		phpItems,
		phpEscape(req.Trigger),
		kfs,
		phpEscape(req.Description),
		*controlID,
		phpEscape(req.Name),
		phpEscape(req.Name),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w gateway group, %w", ErrUpdateOperationFailed, err)
	}

	groups, err = pf.getGatewayGroups(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w gateway groups after updating, %w", ErrGetOperationFailed, err)
	}

	g, err := groups.GetByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w gateway group after updating, %w", ErrGetOperationFailed, err)
	}

	return g, nil
}

func (pf *Client) DeleteGatewayGroup(ctx context.Context, name string) error {
	defer pf.write(&pf.mutexes.GatewayGroup)()

	groups, err := pf.getGatewayGroups(ctx)
	if err != nil {
		return fmt.Errorf("%w gateway groups, %w", ErrGetOperationFailed, err)
	}

	controlID, err := groups.GetControlIDByName(name)
	if err != nil {
		return fmt.Errorf("%w gateway group, %w", ErrGetOperationFailed, err)
	}

	// Delete via PHP (mirrors system_gateway_groups.php delete logic).
	// Also strips the group from any firewall rules that reference it.
	command := fmt.Sprintf(
		"$groups = config_get_path('gateways/gateway_group', array());"+
			"$name = '%s';"+
			"if ((config_get_path('gateways/defaultgw4', '') === $name) || "+
			"(config_get_path('gateways/defaultgw6', '') === $name)) {"+
			"print(json_encode('Cannot delete a gateway group that is the default gateway'));"+
			"} else {"+
			"foreach (config_get_path('filter/rule', array()) as $idx => $rule) {"+
			"if (isset($rule['gateway']) && $rule['gateway'] === $name) {"+
			"config_del_path(\"filter/rule/{$idx}/gateway\");"+
			"}}"+
			"config_del_path('gateways/gateway_group/%d');"+
			"mark_subsystem_dirty('staticroutes');"+
			"write_config('Terraform: deleted gateway group %s');"+
			"print(json_encode(true));"+
			"}",
		phpEscape(name),
		*controlID,
		phpEscape(name),
	)

	var result interface{}
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w gateway group, %w", ErrDeleteOperationFailed, err)
	}

	// Check if the result is an error string.
	if errMsg, ok := result.(string); ok {
		return fmt.Errorf("%w gateway group, %s", ErrDeleteOperationFailed, errMsg)
	}

	// Verify deletion.
	groups, err = pf.getGatewayGroups(ctx)
	if err != nil {
		return fmt.Errorf("%w gateway groups after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := groups.GetByName(name); err == nil {
		return fmt.Errorf("%w gateway group, still exists", ErrDeleteOperationFailed)
	}

	return nil
}

func (pf *Client) ApplyGatewayGroupChanges(ctx context.Context) error {
	pf.mutexes.GatewayGroupApply.Lock()
	defer pf.mutexes.GatewayGroupApply.Unlock()

	// Same apply logic as gateway changes — gateway groups share the routing subsystem.
	command := "require_once(\"filter.inc\");" +
		"require_once(\"openvpn.inc\");" +
		"$retval = 0;" +
		"$retval |= system_routing_configure();" +
		"$retval |= system_resolvconf_generate();" +
		"$retval |= filter_configure();" +
		"setup_gateways_monitor();" +
		"send_event(\"service reload dyndnsall\");" +
		"send_event(\"service reload ipsecdns\");" +
		"if ($retval == 0) clear_subsystem_dirty('staticroutes');" +
		"print(json_encode($retval));"

	var result int
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply gateway group changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}

// phpEscape escapes backslashes and single quotes in strings for safe PHP embedding.
// Backslashes must be escaped first, then single quotes, to prevent sequences
// like \\' from breaking out of PHP string context.
func phpEscape(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "\\", "\\\\"), "'", "\\'")
}
