package pfsense

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const (
	BridgeProtocolRSTP = "rstp"
	BridgeProtocolSTP  = "stp"
)

var (
	bridgeIfRegex          = regexp.MustCompile(`^bridge[0-9]+$`)
	ValidBridgeProtocols   = []string{BridgeProtocolRSTP, BridgeProtocolSTP}
	errBridgeMemberSubset  = "must be a subset of the bridge members"
	errBridgeMemberMissing = "does not reference an existing interface"
	errBridgeSpanMember    = "cannot also be a bridge member"
)

type bridgeResponse struct {
	BridgeIf     string          `json:"bridgeif"`
	Members      string          `json:"members"`
	Description  string          `json:"descr"`
	EnableSTP    json.RawMessage `json:"enablestp"`
	IP6LinkLocal json.RawMessage `json:"ip6linklocal"`
	Protocol     string          `json:"proto"`
	Priority     string          `json:"priority"`
	HelloTime    string          `json:"hellotime"`
	ForwardDelay string          `json:"fwdelay"`
	MaxAge       string          `json:"maxage"`
	HoldCount    string          `json:"holdcnt"`
	MaxAddresses string          `json:"maxaddr"`
	CacheExpire  string          `json:"timeout"`
	STP          string          `json:"stp"`
	Static       string          `json:"static"`
	Private      string          `json:"private"`
	Span         string          `json:"span"`
	Edge         string          `json:"edge"`
	AutoEdge     string          `json:"autoedge"`
	PTP          string          `json:"ptp"`
	AutoPTP      string          `json:"autoptp"`
	ControlID    int             `json:"controlID"` //nolint:tagliatelle
}

type Bridge struct {
	BridgeIf           string
	Members            []string
	Description        string
	EnableSTP          bool
	IP6LinkLocal       bool
	Protocol           string
	Priority           *int
	HelloTime          *int
	ForwardDelay       *int
	MaxAge             *int
	HoldCount          *int
	MaxAddresses       *int
	CacheExpire        *int
	STPInterfaces      []string
	StaticInterfaces   []string
	PrivateInterfaces  []string
	SpanInterfaces     []string
	EdgeInterfaces     []string
	AutoEdgeInterfaces []string
	PTPInterfaces      []string
	AutoPTPInterfaces  []string
	controlID          int
}

type Bridges []Bridge

func (b *Bridge) SetMembers(members []string) error {
	if len(members) == 0 {
		return fmt.Errorf("%w, bridge must have at least one member", ErrClientValidation)
	}

	b.Members = members

	return nil
}

func (b *Bridge) SetDescription(description string) error {
	b.Description = description

	return nil
}

func (b *Bridge) SetProtocol(protocol string) error {
	if protocol == "" {
		return nil
	}

	for _, valid := range ValidBridgeProtocols {
		if protocol == valid {
			b.Protocol = protocol

			return nil
		}
	}

	return fmt.Errorf("%w, bridge protocol must be one of %s", ErrClientValidation, strings.Join(ValidBridgeProtocols, ", "))
}

func (b *Bridge) memberSet() map[string]struct{} {
	set := make(map[string]struct{}, len(b.Members))
	for _, m := range b.Members {
		set[m] = struct{}{}
	}

	return set
}

func (b *Bridge) validateSubset(field string, ifaces []string) error {
	members := b.memberSet()
	for _, iface := range ifaces {
		if _, ok := members[iface]; !ok {
			return fmt.Errorf("%w, %s interface '%s' %s", ErrClientValidation, field, iface, errBridgeMemberSubset)
		}
	}

	return nil
}

func (b *Bridge) validateDisjoint(field string, ifaces []string) error {
	members := b.memberSet()
	for _, iface := range ifaces {
		if _, ok := members[iface]; ok {
			return fmt.Errorf("%w, %s interface '%s' %s", ErrClientValidation, field, iface, errBridgeSpanMember)
		}
	}

	return nil
}

func (bridges Bridges) GetByBridgeIf(bridgeIf string) (*Bridge, error) {
	for i := range bridges {
		if bridges[i].BridgeIf == bridgeIf {
			return &bridges[i], nil
		}
	}

	return nil, fmt.Errorf("bridge %w with interface '%s'", ErrNotFound, bridgeIf)
}

func (bridges Bridges) GetControlIDByBridgeIf(bridgeIf string) (int, error) {
	for _, b := range bridges {
		if b.BridgeIf == bridgeIf {
			return b.controlID, nil
		}
	}

	return -1, fmt.Errorf("bridge %w with interface '%s'", ErrNotFound, bridgeIf)
}

func parseBridgeMemberList(value string) []string {
	if value == "" {
		return nil
	}

	return strings.Split(value, ",")
}

func parseBridgeInt(field, value string) (*int, error) {
	if value == "" {
		return nil, nil //nolint:nilnil
	}

	parsed := 0
	if _, err := fmt.Sscanf(value, "%d", &parsed); err != nil {
		return nil, fmt.Errorf("%w %s from '%s'", ErrUnableToParse, field, value)
	}

	return &parsed, nil
}

func parseBridgeResponse(resp bridgeResponse) (Bridge, error) {
	var bridge Bridge

	bridge.BridgeIf = resp.BridgeIf
	bridge.Members = parseBridgeMemberList(resp.Members)
	bridge.Description = resp.Description
	bridge.EnableSTP = rawIsPresent(resp.EnableSTP)
	bridge.IP6LinkLocal = rawIsPresent(resp.IP6LinkLocal)
	bridge.Protocol = resp.Protocol

	for _, field := range []struct {
		name  string
		value string
		dest  **int
	}{
		{"priority", resp.Priority, &bridge.Priority},
		{"hello time", resp.HelloTime, &bridge.HelloTime},
		{"forward delay", resp.ForwardDelay, &bridge.ForwardDelay},
		{"max age", resp.MaxAge, &bridge.MaxAge},
		{"hold count", resp.HoldCount, &bridge.HoldCount},
		{"max addresses", resp.MaxAddresses, &bridge.MaxAddresses},
		{"cache expire", resp.CacheExpire, &bridge.CacheExpire},
	} {
		parsed, err := parseBridgeInt(field.name, field.value)
		if err != nil {
			return bridge, err
		}

		*field.dest = parsed
	}

	bridge.STPInterfaces = parseBridgeMemberList(resp.STP)
	bridge.StaticInterfaces = parseBridgeMemberList(resp.Static)
	bridge.PrivateInterfaces = parseBridgeMemberList(resp.Private)
	bridge.SpanInterfaces = parseBridgeMemberList(resp.Span)
	bridge.EdgeInterfaces = parseBridgeMemberList(resp.Edge)
	bridge.AutoEdgeInterfaces = parseBridgeMemberList(resp.AutoEdge)
	bridge.PTPInterfaces = parseBridgeMemberList(resp.PTP)
	bridge.AutoPTPInterfaces = parseBridgeMemberList(resp.AutoPTP)
	bridge.controlID = resp.ControlID

	return bridge, nil
}

func bridgeBuildInt(b *strings.Builder, key string, value *int) {
	if value != nil {
		fmt.Fprintf(b, "$bridge['%s'] = '%d';", key, *value)
	}
}

func bridgeBuildMemberList(b *strings.Builder, key string, ifaces []string) {
	if len(ifaces) > 0 {
		fmt.Fprintf(b, "$bridge['%s'] = '%s';", key, phpEscape(strings.Join(ifaces, ",")))
	}
}

func bridgeBuild(req Bridge) string {
	var b strings.Builder

	b.WriteString("$bridge = array();")
	fmt.Fprintf(&b, "$bridge['members'] = '%s';", phpEscape(strings.Join(req.Members, ",")))
	fmt.Fprintf(&b, "$bridge['descr'] = '%s';", phpEscape(req.Description))

	if req.EnableSTP {
		b.WriteString("$bridge['enablestp'] = true;")
	}

	if req.IP6LinkLocal {
		b.WriteString("$bridge['ip6linklocal'] = true;")
	}

	if req.Protocol != "" {
		fmt.Fprintf(&b, "$bridge['proto'] = '%s';", phpEscape(req.Protocol))
	}

	bridgeBuildInt(&b, "priority", req.Priority)
	bridgeBuildInt(&b, "hellotime", req.HelloTime)
	bridgeBuildInt(&b, "fwdelay", req.ForwardDelay)
	bridgeBuildInt(&b, "maxage", req.MaxAge)
	bridgeBuildInt(&b, "holdcnt", req.HoldCount)
	bridgeBuildInt(&b, "maxaddr", req.MaxAddresses)
	bridgeBuildInt(&b, "timeout", req.CacheExpire)

	bridgeBuildMemberList(&b, "stp", req.STPInterfaces)
	bridgeBuildMemberList(&b, "static", req.StaticInterfaces)
	bridgeBuildMemberList(&b, "private", req.PrivateInterfaces)
	bridgeBuildMemberList(&b, "span", req.SpanInterfaces)
	bridgeBuildMemberList(&b, "edge", req.EdgeInterfaces)
	bridgeBuildMemberList(&b, "autoedge", req.AutoEdgeInterfaces)
	bridgeBuildMemberList(&b, "ptp", req.PTPInterfaces)
	bridgeBuildMemberList(&b, "autoptp", req.AutoPTPInterfaces)

	return b.String()
}

func (pf *Client) validateBridge(ctx context.Context, req Bridge, errWrap error) error {
	ifaces, err := pf.getInterfaces(ctx)
	if err != nil {
		return fmt.Errorf("%w interfaces for member validation, %w", ErrGetOperationFailed, err)
	}

	for _, member := range req.Members {
		if _, err := ifaces.GetByLogicalName(member); err != nil {
			return fmt.Errorf("%w bridge, member '%s' %s", errWrap, member, errBridgeMemberMissing)
		}
	}

	for _, field := range []struct {
		name   string
		ifaces []string
	}{
		{"stp", req.STPInterfaces},
		{"static", req.StaticInterfaces},
		{"private", req.PrivateInterfaces},
		{"edge", req.EdgeInterfaces},
		{"autoedge", req.AutoEdgeInterfaces},
		{"ptp", req.PTPInterfaces},
		{"autoptp", req.AutoPTPInterfaces},
	} {
		if err := req.validateSubset(field.name, field.ifaces); err != nil {
			return fmt.Errorf("%w bridge, %w", errWrap, err)
		}
	}

	// Span interfaces mirror all bridge traffic and must NOT be bridge members.
	if err := req.validateDisjoint("span", req.SpanInterfaces); err != nil {
		return fmt.Errorf("%w bridge, %w", errWrap, err)
	}

	return nil
}

func (pf *Client) getBridges(ctx context.Context) (*Bridges, error) {
	command := "$output = array();" +
		"$bridges = config_get_path('bridges/bridge', array());" +
		"foreach ($bridges as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var bridgeResp []bridgeResponse
	if err := pf.executePHPCommand(ctx, command, &bridgeResp); err != nil {
		return nil, err
	}

	bridges := make(Bridges, 0, len(bridgeResp))
	for _, resp := range bridgeResp {
		b, err := parseBridgeResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w bridge response, %w", ErrUnableToParse, err)
		}

		bridges = append(bridges, b)
	}

	return &bridges, nil
}

func (pf *Client) GetBridges(ctx context.Context) (*Bridges, error) {
	defer pf.read(&pf.mutexes.InterfaceBridge)()

	bridges, err := pf.getBridges(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w bridges, %w", ErrGetOperationFailed, err)
	}

	return bridges, nil
}

func (pf *Client) GetBridge(ctx context.Context, bridgeIf string) (*Bridge, error) {
	defer pf.read(&pf.mutexes.InterfaceBridge)()

	bridges, err := pf.getBridges(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w bridges, %w", ErrGetOperationFailed, err)
	}

	b, err := bridges.GetByBridgeIf(bridgeIf)
	if err != nil {
		return nil, fmt.Errorf("%w bridge, %w", ErrGetOperationFailed, err)
	}

	return b, nil
}

func (pf *Client) CreateBridge(ctx context.Context, req Bridge) (*Bridge, error) {
	defer pf.write(&pf.mutexes.InterfaceBridge)()

	if err := pf.validateBridge(ctx, req, ErrCreateOperationFailed); err != nil {
		return nil, err
	}

	command := "require_once('interfaces.inc');" +
		bridgeBuild(req) +
		"$bridge['bridgeif'] = '';" +
		"interface_bridge_configure($bridge);" +
		"if (empty($bridge['bridgeif']) || !preg_match('/^bridge[0-9]+$/', $bridge['bridgeif'])) {" +
		"print(json_encode(array('error' => 'unable to create bridge interface'))); return; }" +
		"$bridges = config_get_path('bridges/bridge', array());" +
		"$bridges[] = $bridge;" +
		"config_set_path('bridges/bridge', $bridges);" +
		"write_config('Terraform: create bridge');" +
		"$confif = convert_real_interface_to_friendly_interface_name($bridge['bridgeif']);" +
		"if ($confif <> '') { interface_configure($confif); }" +
		"print(json_encode($bridge['bridgeif']));"

	var bridgeIf string
	if err := pf.executePHPCommand(ctx, command, &bridgeIf); err != nil {
		return nil, fmt.Errorf("%w bridge, %w", ErrCreateOperationFailed, err)
	}

	if !bridgeIfRegex.MatchString(bridgeIf) {
		return nil, fmt.Errorf("%w bridge, unexpected interface name '%s'", ErrCreateOperationFailed, bridgeIf)
	}

	bridges, err := pf.getBridges(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w bridges after creating, %w", ErrGetOperationFailed, err)
	}

	b, err := bridges.GetByBridgeIf(bridgeIf)
	if err != nil {
		return nil, fmt.Errorf("%w bridge after creating, %w", ErrGetOperationFailed, err)
	}

	return b, nil
}

func (pf *Client) UpdateBridge(ctx context.Context, bridgeIf string, req Bridge) (*Bridge, error) {
	defer pf.write(&pf.mutexes.InterfaceBridge)()

	if err := pf.validateBridge(ctx, req, ErrUpdateOperationFailed); err != nil {
		return nil, err
	}

	bridges, err := pf.getBridges(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w bridges, %w", ErrUpdateOperationFailed, err)
	}

	controlID, err := bridges.GetControlIDByBridgeIf(bridgeIf)
	if err != nil {
		return nil, fmt.Errorf("%w bridge, %w", ErrUpdateOperationFailed, err)
	}

	command := "require_once('interfaces.inc');" +
		bridgeBuild(req) +
		fmt.Sprintf("$bridge['bridgeif'] = '%s';", phpEscape(bridgeIf)) +
		"interface_bridge_configure($bridge);" +
		fmt.Sprintf("config_set_path('bridges/bridge/%d', $bridge);", controlID) +
		"write_config('Terraform: update bridge');" +
		"$confif = convert_real_interface_to_friendly_interface_name($bridge['bridgeif']);" +
		"if ($confif <> '') { interface_configure($confif); }" +
		"print(json_encode($bridge['bridgeif']));"

	var updatedIf string
	if err := pf.executePHPCommand(ctx, command, &updatedIf); err != nil {
		return nil, fmt.Errorf("%w bridge, %w", ErrUpdateOperationFailed, err)
	}

	bridges, err = pf.getBridges(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w bridges after updating, %w", ErrGetOperationFailed, err)
	}

	b, err := bridges.GetByBridgeIf(bridgeIf)
	if err != nil {
		return nil, fmt.Errorf("%w bridge after updating, %w", ErrGetOperationFailed, err)
	}

	return b, nil
}

func (pf *Client) DeleteBridge(ctx context.Context, bridgeIf string) error {
	defer pf.write(&pf.mutexes.InterfaceBridge)()

	bridges, err := pf.getBridges(ctx)
	if err != nil {
		return fmt.Errorf("%w bridges, %w", ErrDeleteOperationFailed, err)
	}

	controlID, err := bridges.GetControlIDByBridgeIf(bridgeIf)
	if err != nil {
		return fmt.Errorf("%w bridge, %w", ErrDeleteOperationFailed, err)
	}

	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"if (does_interface_exist('%s')) { pfSense_interface_destroy('%s'); }"+
			"config_del_path('bridges/bridge/%d');"+
			"write_config('Terraform: delete bridge');"+
			"print(json_encode(true));",
		phpEscape(bridgeIf),
		phpEscape(bridgeIf),
		controlID,
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w bridge, %w", ErrDeleteOperationFailed, err)
	}

	bridges, err = pf.getBridges(ctx)
	if err != nil {
		return fmt.Errorf("%w bridges after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := bridges.GetByBridgeIf(bridgeIf); err == nil {
		return fmt.Errorf("%w bridge, still exists", ErrDeleteOperationFailed)
	}

	return nil
}
