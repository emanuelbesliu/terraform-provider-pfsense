package pfsense

import (
	"context"
	"fmt"
	"strings"
)

type natPortForwardResponse struct {
	Interface        string `json:"interface"`
	IPProtocol       string `json:"ipprotocol"`
	Protocol         string `json:"protocol"`
	Source           any    `json:"source"`
	Destination      any    `json:"destination"`
	Target           string `json:"target"`
	LocalPort        string `json:"local-port"` //nolint:tagliatelle
	Description      string `json:"descr"`
	Disabled         string `json:"disabled"`
	NoRDR            string `json:"nordr"`
	NATReflection    string `json:"natreflection"`
	AssociatedRuleID string `json:"associated-rule-id"` //nolint:tagliatelle
	ControlID        int    `json:"controlID"`          //nolint:tagliatelle
}

type NATPortForward struct {
	Interface        string
	IPProtocol       string
	Protocol         string
	SourceAddress    string
	SourcePort       string
	SourceNot        bool
	DestAddress      string
	DestPort         string
	DestNot          bool
	Target           string
	LocalPort        string
	Description      string
	Disabled         bool
	NoRDR            bool
	NATReflection    string
	AssociatedRuleID string
	controlID        int
}

func (NATPortForward) IPProtocols() []string {
	return []string{"inet", "inet6", "inet46"}
}

func (NATPortForward) Protocols() []string {
	return []string{"tcp", "udp", "tcp/udp", "icmp", "esp", "ah", "gre", "ipv6", "igmp", "pim", "ospf"}
}

func (NATPortForward) NATReflectionModes() []string {
	return []string{"", "enable", "disable", "purenat"}
}

func (NATPortForward) AssociatedRuleIDOptions() []string {
	return []string{"pass", ""}
}

func (r *NATPortForward) SetInterface(iface string) error {
	if err := ValidateInterface(iface); err != nil {
		return fmt.Errorf("%w, invalid interface: %w", ErrClientValidation, err)
	}

	r.Interface = iface

	return nil
}

func (r *NATPortForward) SetIPProtocol(ipp string) error {
	for _, v := range r.IPProtocols() {
		if ipp == v {
			r.IPProtocol = ipp

			return nil
		}
	}

	return fmt.Errorf("%w, ip protocol must be one of: %s", ErrClientValidation, strings.Join(r.IPProtocols(), ", "))
}

func (r *NATPortForward) SetProtocol(proto string) error {
	for _, v := range r.Protocols() {
		if proto == v {
			r.Protocol = proto

			return nil
		}
	}

	return fmt.Errorf("%w, protocol must be one of: %s", ErrClientValidation, strings.Join(r.Protocols(), ", "))
}

func (r *NATPortForward) SetSourceAddress(addr string) error {
	r.SourceAddress = addr

	return nil
}

func (r *NATPortForward) SetSourcePort(port string) error {
	r.SourcePort = port

	return nil
}

func (r *NATPortForward) SetSourceNot(not bool) error {
	r.SourceNot = not

	return nil
}

func (r *NATPortForward) SetDestAddress(addr string) error {
	r.DestAddress = addr

	return nil
}

func (r *NATPortForward) SetDestPort(port string) error {
	r.DestPort = port

	return nil
}

func (r *NATPortForward) SetDestNot(not bool) error {
	r.DestNot = not

	return nil
}

func (r *NATPortForward) SetTarget(target string) error {
	if target == "" {
		return fmt.Errorf("%w, target (redirect target IP) is required", ErrClientValidation)
	}

	r.Target = target

	return nil
}

func (r *NATPortForward) SetLocalPort(port string) error {
	r.LocalPort = port

	return nil
}

func (r *NATPortForward) SetDescription(desc string) error {
	r.Description = desc

	return nil
}

func (r *NATPortForward) SetDisabled(disabled bool) error {
	r.Disabled = disabled

	return nil
}

func (r *NATPortForward) SetNoRDR(nordr bool) error {
	r.NoRDR = nordr

	return nil
}

func (r *NATPortForward) SetNATReflection(mode string) error {
	for _, v := range r.NATReflectionModes() {
		if mode == v {
			r.NATReflection = mode

			return nil
		}
	}

	return fmt.Errorf("%w, NAT reflection must be one of: %s (empty string for system default)", ErrClientValidation, strings.Join(r.NATReflectionModes(), ", "))
}

func (r *NATPortForward) SetAssociatedRuleID(id string) error {
	for _, v := range r.AssociatedRuleIDOptions() {
		if id == v {
			r.AssociatedRuleID = id

			return nil
		}
	}

	return fmt.Errorf("%w, associated rule ID must be one of: 'pass' (auto-create), '' (none)", ErrClientValidation)
}

type NATPortForwards []NATPortForward

func (rules NATPortForwards) GetByDescription(desc string) (*NATPortForward, error) {
	for _, r := range rules {
		if r.Description == desc {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("NAT port forward %w with description '%s'", ErrNotFound, desc)
}

func (rules NATPortForwards) GetControlIDByDescription(desc string) (*int, error) {
	for _, r := range rules {
		if r.Description == desc {
			return &r.controlID, nil
		}
	}

	return nil, fmt.Errorf("NAT port forward %w with description '%s'", ErrNotFound, desc)
}

func parseNATPortForwardResponse(resp natPortForwardResponse) (NATPortForward, error) {
	var r NATPortForward

	if err := r.SetInterface(resp.Interface); err != nil {
		return r, err
	}

	if err := r.SetIPProtocol(resp.IPProtocol); err != nil {
		return r, err
	}

	if err := r.SetProtocol(resp.Protocol); err != nil {
		return r, err
	}

	// Parse source
	srcAddr, srcPort, srcNot := parseSourceOrDest(resp.Source)
	r.SourceAddress = srcAddr
	r.SourcePort = srcPort
	r.SourceNot = srcNot

	// Parse destination
	dstAddr, dstPort, dstNot := parseSourceOrDest(resp.Destination)
	r.DestAddress = dstAddr
	r.DestPort = dstPort
	r.DestNot = dstNot

	if err := r.SetTarget(resp.Target); err != nil {
		return r, err
	}

	r.LocalPort = resp.LocalPort

	if err := r.SetDescription(resp.Description); err != nil {
		return r, err
	}

	r.Disabled = resp.Disabled != ""
	r.NoRDR = resp.NoRDR != ""
	r.NATReflection = resp.NATReflection
	r.AssociatedRuleID = resp.AssociatedRuleID
	r.controlID = resp.ControlID

	return r, nil
}

func (pf *Client) getNATPortForwards(ctx context.Context) (*NATPortForwards, error) {
	command := "$output = array();" +
		"$rules = config_get_path('nat/rule', array());" +
		"foreach ($rules as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var ruleResp []natPortForwardResponse
	if err := pf.executePHPCommand(ctx, command, &ruleResp); err != nil {
		return nil, err
	}

	rules := make(NATPortForwards, 0, len(ruleResp))
	for _, resp := range ruleResp {
		r, err := parseNATPortForwardResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w NAT port forward response, %w", ErrUnableToParse, err)
		}

		rules = append(rules, r)
	}

	return &rules, nil
}

func (pf *Client) GetNATPortForwards(ctx context.Context) (*NATPortForwards, error) {
	defer pf.read(&pf.mutexes.NATPortForward)()

	rules, err := pf.getNATPortForwards(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT port forwards, %w", ErrGetOperationFailed, err)
	}

	return rules, nil
}

func (pf *Client) GetNATPortForward(ctx context.Context, description string) (*NATPortForward, error) {
	defer pf.read(&pf.mutexes.NATPortForward)()

	rules, err := pf.getNATPortForwards(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT port forwards, %w", ErrGetOperationFailed, err)
	}

	r, err := rules.GetByDescription(description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT port forward, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

// natPortForwardPHPSource builds PHP code to set the source fields on $rule.
func natPortForwardPHPSource(req NATPortForward) string {
	var b strings.Builder

	if req.SourceNot {
		b.WriteString("$rule['source']['not'] = '';")
	}

	if req.SourceAddress == "any" {
		b.WriteString("$rule['source']['any'] = '';")
	} else if strings.Contains(req.SourceAddress, "/") {
		parts := strings.SplitN(req.SourceAddress, "/", 2)
		fmt.Fprintf(&b, "$rule['source']['address'] = '%s/%s';", phpEscape(parts[0]), phpEscape(parts[1]))
	} else if isSpecialAddress(req.SourceAddress) {
		fmt.Fprintf(&b, "$rule['source']['network'] = '%s';", phpEscape(req.SourceAddress))
	} else {
		fmt.Fprintf(&b, "$rule['source']['address'] = '%s';", phpEscape(req.SourceAddress))
	}

	if req.SourcePort != "" {
		fmt.Fprintf(&b, "$rule['source']['port'] = '%s';", phpEscape(req.SourcePort))
	}

	return b.String()
}

// natPortForwardPHPDest builds PHP code to set the destination fields on $rule.
func natPortForwardPHPDest(req NATPortForward) string {
	var b strings.Builder

	if req.DestNot {
		b.WriteString("$rule['destination']['not'] = '';")
	}

	if req.DestAddress == "any" {
		b.WriteString("$rule['destination']['any'] = '';")
	} else if strings.Contains(req.DestAddress, "/") {
		parts := strings.SplitN(req.DestAddress, "/", 2)
		fmt.Fprintf(&b, "$rule['destination']['address'] = '%s/%s';", phpEscape(parts[0]), phpEscape(parts[1]))
	} else if isSpecialAddress(req.DestAddress) {
		fmt.Fprintf(&b, "$rule['destination']['network'] = '%s';", phpEscape(req.DestAddress))
	} else {
		fmt.Fprintf(&b, "$rule['destination']['address'] = '%s';", phpEscape(req.DestAddress))
	}

	if req.DestPort != "" {
		fmt.Fprintf(&b, "$rule['destination']['port'] = '%s';", phpEscape(req.DestPort))
	}

	return b.String()
}

func (pf *Client) CreateNATPortForward(ctx context.Context, req NATPortForward) (*NATPortForward, error) {
	defer pf.write(&pf.mutexes.NATPortForward)()

	command := "require_once('filter.inc');" +
		"$rule = array();" +
		fmt.Sprintf("$rule['interface'] = '%s';", phpEscape(req.Interface)) +
		fmt.Sprintf("$rule['ipprotocol'] = '%s';", phpEscape(req.IPProtocol)) +
		fmt.Sprintf("$rule['protocol'] = '%s';", phpEscape(req.Protocol)) +
		"$rule['source'] = array();" +
		natPortForwardPHPSource(req) +
		"$rule['destination'] = array();" +
		natPortForwardPHPDest(req) +
		fmt.Sprintf("$rule['target'] = '%s';", phpEscape(req.Target)) +
		fmt.Sprintf("$rule['local-port'] = '%s';", phpEscape(req.LocalPort)) +
		fmt.Sprintf("$rule['descr'] = '%s';", phpEscape(req.Description))

	if req.Disabled {
		command += "$rule['disabled'] = '';"
	}

	if req.NoRDR {
		command += "$rule['nordr'] = '';"
	}

	if req.NATReflection != "" {
		command += fmt.Sprintf("$rule['natreflection'] = '%s';", phpEscape(req.NATReflection))
	}

	if req.AssociatedRuleID != "" {
		command += fmt.Sprintf("$rule['associated-rule-id'] = '%s';", phpEscape(req.AssociatedRuleID))
	}

	command += "config_set_path('nat/rule/', $rule);" +
		fmt.Sprintf("write_config('Terraform: created NAT port forward %s');", phpEscape(req.Description)) +
		"filter_configure();" +
		"print(json_encode(true));"

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w NAT port forward, %w", ErrCreateOperationFailed, err)
	}

	rules, err := pf.getNATPortForwards(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT port forwards after creating, %w", ErrGetOperationFailed, err)
	}

	r, err := rules.GetByDescription(req.Description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT port forward after creating, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) UpdateNATPortForward(ctx context.Context, req NATPortForward) (*NATPortForward, error) {
	defer pf.write(&pf.mutexes.NATPortForward)()

	rules, err := pf.getNATPortForwards(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT port forwards, %w", ErrGetOperationFailed, err)
	}

	controlID, err := rules.GetControlIDByDescription(req.Description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT port forward, %w", ErrGetOperationFailed, err)
	}

	command := "require_once('filter.inc');" +
		"$rule = array();" +
		fmt.Sprintf("$rule['interface'] = '%s';", phpEscape(req.Interface)) +
		fmt.Sprintf("$rule['ipprotocol'] = '%s';", phpEscape(req.IPProtocol)) +
		fmt.Sprintf("$rule['protocol'] = '%s';", phpEscape(req.Protocol)) +
		"$rule['source'] = array();" +
		natPortForwardPHPSource(req) +
		"$rule['destination'] = array();" +
		natPortForwardPHPDest(req) +
		fmt.Sprintf("$rule['target'] = '%s';", phpEscape(req.Target)) +
		fmt.Sprintf("$rule['local-port'] = '%s';", phpEscape(req.LocalPort)) +
		fmt.Sprintf("$rule['descr'] = '%s';", phpEscape(req.Description))

	if req.Disabled {
		command += "$rule['disabled'] = '';"
	}

	if req.NoRDR {
		command += "$rule['nordr'] = '';"
	}

	if req.NATReflection != "" {
		command += fmt.Sprintf("$rule['natreflection'] = '%s';", phpEscape(req.NATReflection))
	}

	if req.AssociatedRuleID != "" {
		command += fmt.Sprintf("$rule['associated-rule-id'] = '%s';", phpEscape(req.AssociatedRuleID))
	}

	command += fmt.Sprintf("config_set_path('nat/rule/%d', $rule);", *controlID) +
		fmt.Sprintf("write_config('Terraform: updated NAT port forward %s');", phpEscape(req.Description)) +
		"filter_configure();" +
		"print(json_encode(true));"

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w NAT port forward, %w", ErrUpdateOperationFailed, err)
	}

	rules, err = pf.getNATPortForwards(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT port forwards after updating, %w", ErrGetOperationFailed, err)
	}

	r, err := rules.GetByDescription(req.Description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT port forward after updating, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) DeleteNATPortForward(ctx context.Context, description string) error {
	defer pf.write(&pf.mutexes.NATPortForward)()

	rules, err := pf.getNATPortForwards(ctx)
	if err != nil {
		return fmt.Errorf("%w NAT port forwards, %w", ErrGetOperationFailed, err)
	}

	controlID, err := rules.GetControlIDByDescription(description)
	if err != nil {
		return fmt.Errorf("%w NAT port forward, %w", ErrGetOperationFailed, err)
	}

	command := "require_once('filter.inc');" +
		fmt.Sprintf("config_del_path('nat/rule/%d');", *controlID) +
		fmt.Sprintf("write_config('Terraform: deleted NAT port forward %s');", phpEscape(description)) +
		"filter_configure();" +
		"print(json_encode(true));"

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w NAT port forward, %w", ErrDeleteOperationFailed, err)
	}

	rules, err = pf.getNATPortForwards(ctx)
	if err != nil {
		return fmt.Errorf("%w NAT port forwards after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := rules.GetByDescription(description); err == nil {
		return fmt.Errorf("%w NAT port forward, still exists", ErrDeleteOperationFailed)
	}

	return nil
}
