package pfsense

import (
	"context"
	"fmt"
	"strings"
)

type natOutboundRuleResponse struct {
	Interface      string `json:"interface"`
	Protocol       string `json:"protocol"`
	Source         any    `json:"source"`
	SourcePort     string `json:"sourceport"`
	Destination    any    `json:"destination"`
	DstPort        string `json:"dstport"`
	Target         string `json:"target"`
	TargetIP       string `json:"targetip"`
	TargetIPSubnet string `json:"targetip_subnet"` //nolint:tagliatelle
	NATPort        string `json:"natport"`
	PoolOpts       string `json:"poolopts"`
	SourceHashKey  string `json:"source_hash_key"` //nolint:tagliatelle
	StaticNATPort  string `json:"staticnatport"`
	NoSync         string `json:"nosync"`
	NoNAT          string `json:"nonat"`
	Disabled       string `json:"disabled"`
	Description    string `json:"descr"`
	ControlID      int    `json:"controlID"` //nolint:tagliatelle
}

type NATOutboundRule struct {
	Interface      string
	Protocol       string
	SourceAddress  string
	SourcePort     string
	SourceNot      bool
	DestAddress    string
	DestPort       string
	DestNot        bool
	Target         string
	TargetIP       string
	TargetIPSubnet string
	NATPort        string
	PoolOpts       string
	SourceHashKey  string
	StaticNATPort  bool
	NoSync         bool
	NoNAT          bool
	Disabled       bool
	Description    string
	controlID      int
}

func (NATOutboundRule) Protocols() []string {
	return []string{"", "tcp", "udp", "tcp/udp", "icmp", "esp", "ah", "gre", "ipv6", "igmp", "pim", "ospf"}
}

func (NATOutboundRule) PoolOptions() []string {
	return []string{"", "round-robin", "round-robin sticky-address", "random", "random sticky-address", "source-hash", "bitmask"}
}

func (NATOutboundRule) Targets() []string {
	return []string{"", "(self)", "other-subnet"}
}

func (NATOutboundRule) Modes() []string {
	return []string{"automatic", "hybrid", "advanced", "disabled"}
}

func (r *NATOutboundRule) SetInterface(iface string) error {
	if err := ValidateInterface(iface); err != nil {
		return fmt.Errorf("%w, invalid interface: %w", ErrClientValidation, err)
	}

	r.Interface = iface

	return nil
}

func (r *NATOutboundRule) SetProtocol(proto string) error {
	for _, v := range r.Protocols() {
		if proto == v {
			r.Protocol = proto

			return nil
		}
	}

	return fmt.Errorf("%w, protocol must be one of: %s (empty string for any)", ErrClientValidation, strings.Join(r.Protocols(), ", "))
}

func (r *NATOutboundRule) SetSourceAddress(addr string) error {
	r.SourceAddress = addr

	return nil
}

func (r *NATOutboundRule) SetSourcePort(port string) error {
	r.SourcePort = port

	return nil
}

func (r *NATOutboundRule) SetSourceNot(not bool) error {
	r.SourceNot = not

	return nil
}

func (r *NATOutboundRule) SetDestAddress(addr string) error {
	r.DestAddress = addr

	return nil
}

func (r *NATOutboundRule) SetDestPort(port string) error {
	r.DestPort = port

	return nil
}

func (r *NATOutboundRule) SetDestNot(not bool) error {
	r.DestNot = not

	return nil
}

func (r *NATOutboundRule) SetTarget(target string) error {
	r.Target = target

	return nil
}

func (r *NATOutboundRule) SetTargetIP(ip string) error {
	r.TargetIP = ip

	return nil
}

func (r *NATOutboundRule) SetTargetIPSubnet(subnet string) error {
	r.TargetIPSubnet = subnet

	return nil
}

func (r *NATOutboundRule) SetNATPort(port string) error {
	r.NATPort = port

	return nil
}

func (r *NATOutboundRule) SetPoolOpts(opts string) error {
	for _, v := range r.PoolOptions() {
		if opts == v {
			r.PoolOpts = opts

			return nil
		}
	}

	return fmt.Errorf("%w, pool options must be one of: %s (empty string for default)", ErrClientValidation, strings.Join(r.PoolOptions(), ", "))
}

func (r *NATOutboundRule) SetSourceHashKey(key string) error {
	r.SourceHashKey = key

	return nil
}

func (r *NATOutboundRule) SetStaticNATPort(static bool) error {
	r.StaticNATPort = static

	return nil
}

func (r *NATOutboundRule) SetNoSync(nosync bool) error {
	r.NoSync = nosync

	return nil
}

func (r *NATOutboundRule) SetNoNAT(nonat bool) error {
	r.NoNAT = nonat

	return nil
}

func (r *NATOutboundRule) SetDisabled(disabled bool) error {
	r.Disabled = disabled

	return nil
}

func (r *NATOutboundRule) SetDescription(desc string) error {
	r.Description = desc

	return nil
}

type NATOutboundRules []NATOutboundRule

func (rules NATOutboundRules) GetByDescription(desc string) (*NATOutboundRule, error) {
	for _, r := range rules {
		if r.Description == desc {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("NAT outbound rule %w with description '%s'", ErrNotFound, desc)
}

func (rules NATOutboundRules) GetControlIDByDescription(desc string) (*int, error) {
	for _, r := range rules {
		if r.Description == desc {
			return &r.controlID, nil
		}
	}

	return nil, fmt.Errorf("NAT outbound rule %w with description '%s'", ErrNotFound, desc)
}

func parseNATOutboundRuleResponse(resp natOutboundRuleResponse) (NATOutboundRule, error) {
	var r NATOutboundRule

	if err := r.SetInterface(resp.Interface); err != nil {
		return r, err
	}

	if err := r.SetProtocol(resp.Protocol); err != nil {
		return r, err
	}

	// Parse source
	srcAddr, srcPort, srcNot := parseSourceOrDest(resp.Source)
	r.SourceAddress = srcAddr

	if resp.SourcePort != "" {
		r.SourcePort = resp.SourcePort
	} else {
		r.SourcePort = srcPort
	}

	r.SourceNot = srcNot

	// Parse destination
	dstAddr, dstPort, dstNot := parseSourceOrDest(resp.Destination)
	r.DestAddress = dstAddr

	if resp.DstPort != "" {
		r.DestPort = resp.DstPort
	} else {
		r.DestPort = dstPort
	}

	r.DestNot = dstNot

	r.Target = resp.Target
	r.TargetIP = resp.TargetIP
	r.TargetIPSubnet = resp.TargetIPSubnet
	r.NATPort = resp.NATPort
	r.PoolOpts = resp.PoolOpts
	r.SourceHashKey = resp.SourceHashKey
	r.StaticNATPort = resp.StaticNATPort != ""
	r.NoSync = resp.NoSync != ""
	r.NoNAT = resp.NoNAT != ""
	r.Disabled = resp.Disabled != ""
	r.Description = resp.Description
	r.controlID = resp.ControlID

	return r, nil
}

// natOutboundPHPSource builds PHP code to set the source fields on $rule.
func natOutboundPHPSource(req NATOutboundRule) string {
	var b strings.Builder

	if req.SourceNot {
		b.WriteString("$rule['source']['not'] = '';")
	}

	if req.SourceAddress == "" || req.SourceAddress == "any" {
		b.WriteString("$rule['source']['any'] = '';")
	} else if strings.Contains(req.SourceAddress, "/") {
		fmt.Fprintf(&b, "$rule['source']['network'] = '%s';", phpEscape(req.SourceAddress))
	} else if isSpecialAddress(req.SourceAddress) {
		fmt.Fprintf(&b, "$rule['source']['network'] = '%s';", phpEscape(req.SourceAddress))
	} else {
		fmt.Fprintf(&b, "$rule['source']['network'] = '%s';", phpEscape(req.SourceAddress))
	}

	return b.String()
}

// natOutboundPHPDest builds PHP code to set the destination fields on $rule.
func natOutboundPHPDest(req NATOutboundRule) string {
	var b strings.Builder

	if req.DestNot {
		b.WriteString("$rule['destination']['not'] = '';")
	}

	if req.DestAddress == "" || req.DestAddress == "any" {
		b.WriteString("$rule['destination']['any'] = '';")
	} else if strings.Contains(req.DestAddress, "/") {
		parts := strings.SplitN(req.DestAddress, "/", 2)
		fmt.Fprintf(&b, "$rule['destination']['address'] = '%s/%s';", phpEscape(parts[0]), phpEscape(parts[1]))
	} else if isSpecialAddress(req.DestAddress) {
		fmt.Fprintf(&b, "$rule['destination']['network'] = '%s';", phpEscape(req.DestAddress))
	} else {
		fmt.Fprintf(&b, "$rule['destination']['address'] = '%s';", phpEscape(req.DestAddress))
	}

	return b.String()
}

func (pf *Client) getNATOutboundRules(ctx context.Context) (*NATOutboundRules, error) {
	command := "$output = array();" +
		"$rules = config_get_path('nat/outbound/rule', array());" +
		"foreach ($rules as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var ruleResp []natOutboundRuleResponse
	if err := pf.executePHPCommand(ctx, command, &ruleResp); err != nil {
		return nil, err
	}

	rules := make(NATOutboundRules, 0, len(ruleResp))
	for _, resp := range ruleResp {
		r, err := parseNATOutboundRuleResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w NAT outbound rule response, %w", ErrUnableToParse, err)
		}

		rules = append(rules, r)
	}

	return &rules, nil
}

func (pf *Client) GetNATOutboundRules(ctx context.Context) (*NATOutboundRules, error) {
	defer pf.read(&pf.mutexes.NATOutbound)()

	rules, err := pf.getNATOutboundRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT outbound rules, %w", ErrGetOperationFailed, err)
	}

	return rules, nil
}

func (pf *Client) GetNATOutboundRule(ctx context.Context, description string) (*NATOutboundRule, error) {
	defer pf.read(&pf.mutexes.NATOutbound)()

	rules, err := pf.getNATOutboundRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT outbound rules, %w", ErrGetOperationFailed, err)
	}

	r, err := rules.GetByDescription(description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT outbound rule, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) GetNATOutboundMode(ctx context.Context) (string, error) {
	defer pf.read(&pf.mutexes.NATOutbound)()

	command := "$mode = config_get_path('nat/outbound/mode', 'automatic');" +
		"print(json_encode($mode));"

	var mode string
	if err := pf.executePHPCommand(ctx, command, &mode); err != nil {
		return "", fmt.Errorf("%w NAT outbound mode, %w", ErrGetOperationFailed, err)
	}

	if mode == "" {
		mode = "automatic"
	}

	return mode, nil
}

func (pf *Client) SetNATOutboundMode(ctx context.Context, mode string) error {
	defer pf.write(&pf.mutexes.NATOutbound)()

	valid := false
	for _, v := range (NATOutboundRule{}).Modes() {
		if mode == v {
			valid = true

			break
		}
	}

	if !valid {
		return fmt.Errorf("%w, NAT outbound mode must be one of: %s", ErrClientValidation, strings.Join((NATOutboundRule{}).Modes(), ", "))
	}

	command := "require_once('filter.inc');" +
		fmt.Sprintf("config_set_path('nat/outbound/mode', '%s');", phpEscape(mode)) +
		fmt.Sprintf("write_config('Terraform: set NAT outbound mode to %s');", phpEscape(mode)) +
		"filter_configure();" +
		"print(json_encode(true));"

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w NAT outbound mode, %w", ErrUpdateOperationFailed, err)
	}

	return nil
}

func (pf *Client) CreateNATOutboundRule(ctx context.Context, req NATOutboundRule) (*NATOutboundRule, error) {
	defer pf.write(&pf.mutexes.NATOutbound)()

	command := "require_once('filter.inc');" +
		"$rule = array();" +
		fmt.Sprintf("$rule['interface'] = '%s';", phpEscape(req.Interface)) +
		fmt.Sprintf("$rule['protocol'] = '%s';", phpEscape(req.Protocol)) +
		"$rule['source'] = array();" +
		natOutboundPHPSource(req) +
		"$rule['destination'] = array();" +
		natOutboundPHPDest(req) +
		fmt.Sprintf("$rule['descr'] = '%s';", phpEscape(req.Description))

	if req.SourcePort != "" {
		command += fmt.Sprintf("$rule['sourceport'] = '%s';", phpEscape(req.SourcePort))
	}

	if req.DestPort != "" {
		command += fmt.Sprintf("$rule['dstport'] = '%s';", phpEscape(req.DestPort))
	}

	if req.Target != "" {
		command += fmt.Sprintf("$rule['target'] = '%s';", phpEscape(req.Target))
	}

	if req.TargetIP != "" {
		command += fmt.Sprintf("$rule['targetip'] = '%s';", phpEscape(req.TargetIP))
	}

	if req.TargetIPSubnet != "" {
		command += fmt.Sprintf("$rule['targetip_subnet'] = '%s';", phpEscape(req.TargetIPSubnet))
	}

	if req.NATPort != "" {
		command += fmt.Sprintf("$rule['natport'] = '%s';", phpEscape(req.NATPort))
	}

	if req.PoolOpts != "" {
		command += fmt.Sprintf("$rule['poolopts'] = '%s';", phpEscape(req.PoolOpts))
	}

	if req.SourceHashKey != "" {
		command += fmt.Sprintf("$rule['source_hash_key'] = '%s';", phpEscape(req.SourceHashKey))
	}

	if req.StaticNATPort {
		command += "$rule['staticnatport'] = '';"
	}

	if req.NoSync {
		command += "$rule['nosync'] = '';"
	}

	if req.NoNAT {
		command += "$rule['nonat'] = '';"
	}

	if req.Disabled {
		command += "$rule['disabled'] = '';"
	}

	command += "config_set_path('nat/outbound/rule/', $rule);" +
		fmt.Sprintf("write_config('Terraform: created NAT outbound rule %s');", phpEscape(req.Description)) +
		"filter_configure();" +
		"print(json_encode(true));"

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w NAT outbound rule, %w", ErrCreateOperationFailed, err)
	}

	rules, err := pf.getNATOutboundRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT outbound rules after creating, %w", ErrGetOperationFailed, err)
	}

	r, err := rules.GetByDescription(req.Description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT outbound rule after creating, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) UpdateNATOutboundRule(ctx context.Context, req NATOutboundRule) (*NATOutboundRule, error) {
	defer pf.write(&pf.mutexes.NATOutbound)()

	rules, err := pf.getNATOutboundRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT outbound rules, %w", ErrGetOperationFailed, err)
	}

	controlID, err := rules.GetControlIDByDescription(req.Description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT outbound rule, %w", ErrGetOperationFailed, err)
	}

	command := "require_once('filter.inc');" +
		"$rule = array();" +
		fmt.Sprintf("$rule['interface'] = '%s';", phpEscape(req.Interface)) +
		fmt.Sprintf("$rule['protocol'] = '%s';", phpEscape(req.Protocol)) +
		"$rule['source'] = array();" +
		natOutboundPHPSource(req) +
		"$rule['destination'] = array();" +
		natOutboundPHPDest(req) +
		fmt.Sprintf("$rule['descr'] = '%s';", phpEscape(req.Description))

	if req.SourcePort != "" {
		command += fmt.Sprintf("$rule['sourceport'] = '%s';", phpEscape(req.SourcePort))
	}

	if req.DestPort != "" {
		command += fmt.Sprintf("$rule['dstport'] = '%s';", phpEscape(req.DestPort))
	}

	if req.Target != "" {
		command += fmt.Sprintf("$rule['target'] = '%s';", phpEscape(req.Target))
	}

	if req.TargetIP != "" {
		command += fmt.Sprintf("$rule['targetip'] = '%s';", phpEscape(req.TargetIP))
	}

	if req.TargetIPSubnet != "" {
		command += fmt.Sprintf("$rule['targetip_subnet'] = '%s';", phpEscape(req.TargetIPSubnet))
	}

	if req.NATPort != "" {
		command += fmt.Sprintf("$rule['natport'] = '%s';", phpEscape(req.NATPort))
	}

	if req.PoolOpts != "" {
		command += fmt.Sprintf("$rule['poolopts'] = '%s';", phpEscape(req.PoolOpts))
	}

	if req.SourceHashKey != "" {
		command += fmt.Sprintf("$rule['source_hash_key'] = '%s';", phpEscape(req.SourceHashKey))
	}

	if req.StaticNATPort {
		command += "$rule['staticnatport'] = '';"
	}

	if req.NoSync {
		command += "$rule['nosync'] = '';"
	}

	if req.NoNAT {
		command += "$rule['nonat'] = '';"
	}

	if req.Disabled {
		command += "$rule['disabled'] = '';"
	}

	command += fmt.Sprintf("config_set_path('nat/outbound/rule/%d', $rule);", *controlID) +
		fmt.Sprintf("write_config('Terraform: updated NAT outbound rule %s');", phpEscape(req.Description)) +
		"filter_configure();" +
		"print(json_encode(true));"

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w NAT outbound rule, %w", ErrUpdateOperationFailed, err)
	}

	rules, err = pf.getNATOutboundRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT outbound rules after updating, %w", ErrGetOperationFailed, err)
	}

	r, err := rules.GetByDescription(req.Description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT outbound rule after updating, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) DeleteNATOutboundRule(ctx context.Context, description string) error {
	defer pf.write(&pf.mutexes.NATOutbound)()

	rules, err := pf.getNATOutboundRules(ctx)
	if err != nil {
		return fmt.Errorf("%w NAT outbound rules, %w", ErrGetOperationFailed, err)
	}

	controlID, err := rules.GetControlIDByDescription(description)
	if err != nil {
		return fmt.Errorf("%w NAT outbound rule, %w", ErrGetOperationFailed, err)
	}

	command := "require_once('filter.inc');" +
		fmt.Sprintf("config_del_path('nat/outbound/rule/%d');", *controlID) +
		fmt.Sprintf("write_config('Terraform: deleted NAT outbound rule %s');", phpEscape(description)) +
		"filter_configure();" +
		"print(json_encode(true));"

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w NAT outbound rule, %w", ErrDeleteOperationFailed, err)
	}

	rules, err = pf.getNATOutboundRules(ctx)
	if err != nil {
		return fmt.Errorf("%w NAT outbound rules after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := rules.GetByDescription(description); err == nil {
		return fmt.Errorf("%w NAT outbound rule, still exists", ErrDeleteOperationFailed)
	}

	return nil
}
