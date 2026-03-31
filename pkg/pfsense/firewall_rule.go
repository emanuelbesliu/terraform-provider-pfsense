package pfsense

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	firewallRuleSourcePortSep      = "-"
	firewallRuleDestinationPortSep = "-"
)

type firewallRuleResponse struct {
	Type             string  `json:"type"`
	Interface        string  `json:"interface"`
	IPProtocol       string  `json:"ipprotocol"`
	Protocol         string  `json:"protocol"`
	Source           any     `json:"source"`
	Destination      any     `json:"destination"`
	Description      string  `json:"descr"`
	Disabled         *string `json:"disabled"`
	Log              *string `json:"log"`
	Tracker          string  `json:"tracker"`
	ControlID        int     `json:"controlID"`          //nolint:tagliatelle
	AssociatedRuleID string  `json:"associated-rule-id"` //nolint:tagliatelle
}

type FirewallRule struct {
	Type          string
	Interface     string
	IPProtocol    string
	Protocol      string
	SourceAddress string
	SourcePort    string
	SourceNot     bool
	DestAddress   string
	DestPort      string
	DestNot       bool
	Description   string
	Disabled      bool
	Log           bool
	Tracker       string
	controlID     int
}

func (FirewallRule) Types() []string {
	return []string{"pass", "block", "reject"}
}

func (FirewallRule) IPProtocols() []string {
	return []string{"inet", "inet6", "inet46"}
}

func (FirewallRule) Protocols() []string {
	return []string{"any", "tcp", "udp", "tcp/udp", "icmp", "esp", "ah", "gre", "ipv6", "igmp", "pim", "ospf", "carp", "pfsync"}
}

func (r *FirewallRule) SetType(t string) error {
	valid := false
	for _, v := range r.Types() {
		if t == v {
			valid = true

			break
		}
	}

	if !valid {
		return fmt.Errorf("%w, type must be one of: %s", ErrClientValidation, strings.Join(r.Types(), ", "))
	}

	r.Type = t

	return nil
}

func (r *FirewallRule) SetInterface(iface string) error {
	if err := ValidateInterface(iface); err != nil {
		return fmt.Errorf("%w, invalid interface: %w", ErrClientValidation, err)
	}

	r.Interface = iface

	return nil
}

func (r *FirewallRule) SetIPProtocol(ipprotocol string) error {
	valid := false
	for _, p := range r.IPProtocols() {
		if ipprotocol == p {
			valid = true

			break
		}
	}

	if !valid {
		return fmt.Errorf("%w, ip protocol must be one of: %s", ErrClientValidation, strings.Join(r.IPProtocols(), ", "))
	}

	r.IPProtocol = ipprotocol

	return nil
}

func (r *FirewallRule) SetProtocol(protocol string) error {
	valid := false
	for _, p := range r.Protocols() {
		if protocol == p {
			valid = true

			break
		}
	}

	if !valid {
		return fmt.Errorf("%w, protocol must be one of: %s", ErrClientValidation, strings.Join(r.Protocols(), ", "))
	}

	r.Protocol = protocol

	return nil
}

func (r *FirewallRule) SetSourceAddress(addr string) error {
	r.SourceAddress = addr

	return nil
}

func (r *FirewallRule) SetSourcePort(port string) error {
	r.SourcePort = port

	return nil
}

func (r *FirewallRule) SetSourceNot(not bool) error {
	r.SourceNot = not

	return nil
}

func (r *FirewallRule) SetDestAddress(addr string) error {
	r.DestAddress = addr

	return nil
}

func (r *FirewallRule) SetDestPort(port string) error {
	r.DestPort = port

	return nil
}

func (r *FirewallRule) SetDestNot(not bool) error {
	r.DestNot = not

	return nil
}

func (r *FirewallRule) SetDescription(description string) error {
	r.Description = description

	return nil
}

func (r *FirewallRule) SetDisabled(disabled bool) error {
	r.Disabled = disabled

	return nil
}

func (r *FirewallRule) SetLog(log bool) error {
	r.Log = log

	return nil
}

func (r *FirewallRule) SetTracker(tracker string) error {
	r.Tracker = tracker

	return nil
}

type FirewallRules []FirewallRule

func (rules FirewallRules) GetByTracker(tracker string) (*FirewallRule, error) {
	for _, r := range rules {
		if r.Tracker == tracker {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("firewall rule %w with tracker '%s'", ErrNotFound, tracker)
}

func (rules FirewallRules) GetControlIDByTracker(tracker string) (*int, error) {
	for _, r := range rules {
		if r.Tracker == tracker {
			return &r.controlID, nil
		}
	}

	return nil, fmt.Errorf("firewall rule %w with tracker '%s'", ErrNotFound, tracker)
}

// parseSourceOrDest extracts the address, port, and "not" flag from a pfSense
// source or destination field. pfSense stores these as either a map or an
// empty string (when no source/dest config).
func parseSourceOrDest(raw any) (address string, port string, not bool) {
	m, ok := raw.(map[string]any)
	if !ok {
		return "any", "", false
	}

	// Check for "not" flag
	if _, exists := m["not"]; exists {
		not = true
	}

	// Check for address - pfSense stores as "address" or "any" or "network"
	if addr, exists := m["address"]; exists {
		address = fmt.Sprintf("%v", addr)
	} else if _, exists := m["any"]; exists {
		address = "any"
	} else if net, exists := m["network"]; exists {
		address = fmt.Sprintf("%v", net)
	} else {
		address = "any"
	}

	// Port
	if p, exists := m["port"]; exists {
		port = fmt.Sprintf("%v", p)
	}

	return address, port, not
}

func parseFirewallRuleResponse(resp firewallRuleResponse) (FirewallRule, error) {
	var r FirewallRule

	if err := r.SetType(resp.Type); err != nil {
		return r, err
	}

	if err := r.SetInterface(resp.Interface); err != nil {
		return r, err
	}

	if err := r.SetIPProtocol(resp.IPProtocol); err != nil {
		return r, err
	}

	protocol := resp.Protocol
	if protocol == "" {
		protocol = "any"
	}

	if err := r.SetProtocol(protocol); err != nil {
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

	if err := r.SetDescription(resp.Description); err != nil {
		return r, err
	}

	r.Disabled = resp.Disabled != nil
	r.Log = resp.Log != nil
	r.Tracker = resp.Tracker
	r.controlID = resp.ControlID

	return r, nil
}

func (pf *Client) getFirewallRules(ctx context.Context) (*FirewallRules, error) {
	command := "$output = array();" +
		"if (isset($config['filter']['rule']) && is_array($config['filter']['rule'])) {" +
		"foreach ($config['filter']['rule'] as $k => $v) {" +
		"$v['controlID'] = $k; array_push($output, $v);" +
		"}};" +
		"print_r(json_encode($output));"

	var ruleResp []firewallRuleResponse
	if err := pf.executePHPCommand(ctx, command, &ruleResp); err != nil {
		return nil, err
	}

	rules := make(FirewallRules, 0, len(ruleResp))
	for _, resp := range ruleResp {
		// Skip NAT-associated rules — these are auto-generated by pfSense
		// for NAT port forward entries and have no "type" field.
		if resp.AssociatedRuleID != "" {
			continue
		}

		r, err := parseFirewallRuleResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w firewall rule response, %w", ErrUnableToParse, err)
		}

		rules = append(rules, r)
	}

	return &rules, nil
}

func (pf *Client) GetFirewallRules(ctx context.Context) (*FirewallRules, error) {
	defer pf.read(&pf.mutexes.FirewallRule)()

	rules, err := pf.getFirewallRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w firewall rules, %w", ErrGetOperationFailed, err)
	}

	return rules, nil
}

func (pf *Client) GetFirewallRule(ctx context.Context, tracker string) (*FirewallRule, error) {
	defer pf.read(&pf.mutexes.FirewallRule)()

	rules, err := pf.getFirewallRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w firewall rules, %w", ErrGetOperationFailed, err)
	}

	r, err := rules.GetByTracker(tracker)
	if err != nil {
		return nil, fmt.Errorf("%w firewall rule, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

func firewallRuleFormValues(ruleReq FirewallRule) url.Values {
	values := url.Values{
		"type":       {ruleReq.Type},
		"interface":  {ruleReq.Interface},
		"ipprotocol": {ruleReq.IPProtocol},
		"descr":      {ruleReq.Description},
		"save":       {"Save"},
	}

	// Protocol: pfSense uses "proto" for the form field
	if ruleReq.Protocol == "any" {
		values.Set("proto", "any")
	} else {
		values.Set("proto", ruleReq.Protocol)
	}

	// Source
	if ruleReq.SourceNot {
		values.Set("srcnot", "yes")
	}

	if ruleReq.SourceAddress == "any" {
		values.Set("srctype", "any")
	} else if strings.Contains(ruleReq.SourceAddress, "/") {
		values.Set("srctype", "network")
		parts := strings.SplitN(ruleReq.SourceAddress, "/", 2)
		values.Set("src", parts[0])
		values.Set("srcmask", parts[1])
	} else if isSpecialAddress(ruleReq.SourceAddress) {
		values.Set("srctype", ruleReq.SourceAddress)
	} else {
		values.Set("srctype", "single")
		values.Set("src", ruleReq.SourceAddress)
		values.Set("srcmask", "32")
	}

	if ruleReq.SourcePort != "" {
		values.Set("srcbeginport", ruleReq.SourcePort)
		values.Set("srcendport", ruleReq.SourcePort)

		// Handle port ranges
		if strings.Contains(ruleReq.SourcePort, "-") {
			parts := strings.SplitN(ruleReq.SourcePort, "-", 2)
			values.Set("srcbeginport", parts[0])
			values.Set("srcendport", parts[1])
		}
	}

	// Destination
	if ruleReq.DestNot {
		values.Set("dstnot", "yes")
	}

	if ruleReq.DestAddress == "any" {
		values.Set("dsttype", "any")
	} else if strings.Contains(ruleReq.DestAddress, "/") {
		values.Set("dsttype", "network")
		parts := strings.SplitN(ruleReq.DestAddress, "/", 2)
		values.Set("dst", parts[0])
		values.Set("dstmask", parts[1])
	} else if isSpecialAddress(ruleReq.DestAddress) {
		values.Set("dsttype", ruleReq.DestAddress)
	} else {
		values.Set("dsttype", "single")
		values.Set("dst", ruleReq.DestAddress)
		values.Set("dstmask", "32")
	}

	if ruleReq.DestPort != "" {
		values.Set("dstbeginport", ruleReq.DestPort)
		values.Set("dstendport", ruleReq.DestPort)

		// Handle port ranges
		if strings.Contains(ruleReq.DestPort, "-") {
			parts := strings.SplitN(ruleReq.DestPort, "-", 2)
			values.Set("dstbeginport", parts[0])
			values.Set("dstendport", parts[1])
		}
	}

	if ruleReq.Disabled {
		values.Set("disabled", "yes")
	}

	if ruleReq.Log {
		values.Set("log", "yes")
	}

	return values
}

// isSpecialAddress checks if the address is a pfSense special address
// (interface address, network, etc.) rather than a literal IP/CIDR.
func isSpecialAddress(addr string) bool {
	specialSuffixes := []string{"ip", "net"}
	for _, suffix := range specialSuffixes {
		if strings.HasSuffix(addr, suffix) && len(addr) > len(suffix) {
			return true
		}
	}

	// Check for interface names used as source/dest (e.g., "lan", "wan", "opt1")
	specialPrefixes := []string{
		"(self)", "pppoe", "l2tp",
	}
	for _, prefix := range specialPrefixes {
		if strings.HasPrefix(addr, prefix) {
			return true
		}
	}

	return false
}

func (pf *Client) createOrUpdateFirewallRule(ctx context.Context, ruleReq FirewallRule, controlID *int) error {
	relativeURL := url.URL{Path: "firewall_rules_edit.php"}
	values := firewallRuleFormValues(ruleReq)

	if controlID != nil {
		q := relativeURL.Query()
		q.Set("id", strconv.Itoa(*controlID))
		relativeURL.RawQuery = q.Encode()
		values.Set("id", strconv.Itoa(*controlID))
	}

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return err
	}

	return scrapeHTMLValidationErrors(doc)
}

func (pf *Client) CreateFirewallRule(ctx context.Context, ruleReq FirewallRule) (*FirewallRule, error) {
	defer pf.write(&pf.mutexes.FirewallRule)()

	if err := pf.createOrUpdateFirewallRule(ctx, ruleReq, nil); err != nil {
		return nil, fmt.Errorf("%w firewall rule, %w", ErrCreateOperationFailed, err)
	}

	rules, err := pf.getFirewallRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w firewall rules after creating, %w", ErrGetOperationFailed, err)
	}

	// Find the newly created rule. Since we don't know the tracker yet,
	// find the last rule matching our interface and type (the new rule is
	// appended to the end of the list).
	var found *FirewallRule
	for i := len(*rules) - 1; i >= 0; i-- {
		r := (*rules)[i]
		if r.Interface == ruleReq.Interface &&
			r.Type == ruleReq.Type &&
			r.Description == ruleReq.Description {
			found = &r

			break
		}
	}

	if found == nil {
		return nil, fmt.Errorf("%w firewall rule after creating, could not find newly created rule", ErrGetOperationFailed)
	}

	return found, nil
}

func (pf *Client) UpdateFirewallRule(ctx context.Context, ruleReq FirewallRule) (*FirewallRule, error) {
	defer pf.write(&pf.mutexes.FirewallRule)()

	rules, err := pf.getFirewallRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w firewall rules, %w", ErrGetOperationFailed, err)
	}

	controlID, err := rules.GetControlIDByTracker(ruleReq.Tracker)
	if err != nil {
		return nil, fmt.Errorf("%w firewall rule, %w", ErrGetOperationFailed, err)
	}

	if err := pf.createOrUpdateFirewallRule(ctx, ruleReq, controlID); err != nil {
		return nil, fmt.Errorf("%w firewall rule, %w", ErrUpdateOperationFailed, err)
	}

	rules, err = pf.getFirewallRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w firewall rules after updating, %w", ErrGetOperationFailed, err)
	}

	r, err := rules.GetByTracker(ruleReq.Tracker)
	if err != nil {
		return nil, fmt.Errorf("%w firewall rule after updating, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) DeleteFirewallRule(ctx context.Context, tracker string) error {
	defer pf.write(&pf.mutexes.FirewallRule)()

	rules, err := pf.getFirewallRules(ctx)
	if err != nil {
		return fmt.Errorf("%w firewall rules, %w", ErrGetOperationFailed, err)
	}

	controlID, err := rules.GetControlIDByTracker(tracker)
	if err != nil {
		return fmt.Errorf("%w firewall rule, %w", ErrGetOperationFailed, err)
	}

	relativeURL := url.URL{Path: "firewall_rules.php"}
	values := url.Values{
		"act": {"del"},
		"id":  {strconv.Itoa(*controlID)},
	}

	_, err = pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return fmt.Errorf("%w firewall rule, %w", ErrDeleteOperationFailed, err)
	}

	rules, err = pf.getFirewallRules(ctx)
	if err != nil {
		return fmt.Errorf("%w firewall rules after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := rules.GetByTracker(tracker); err == nil {
		return fmt.Errorf("%w firewall rule, still exists", ErrDeleteOperationFailed)
	}

	return nil
}

func (pf *Client) ApplyFirewallRuleChanges(ctx context.Context) error {
	pf.mutexes.FirewallFilterReload.Lock()
	defer pf.mutexes.FirewallFilterReload.Unlock()

	command := "require_once(\"filter.inc\");" +
		"$retval = 0;" +
		"$retval |= filter_configure();" +
		"clear_subsystem_dirty('filter');" +
		"print(json_encode($retval));"

	var result int
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply firewall rule changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}
