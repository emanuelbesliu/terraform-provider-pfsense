package pfsense

import (
	"context"
	"fmt"
	"strings"
)

// natOneToOneResponse represents the JSON response from pfSense for a 1:1 NAT rule
type natOneToOneResponse struct {
	External      any    `json:"external"`
	Interface     string `json:"interface"`
	IPProtocol    string `json:"ipprotocol"`
	Source        any    `json:"source"`
	Destination   any    `json:"destination"`
	Description   string `json:"descr"`
	Disabled      string `json:"disabled"`
	NoBinat       string `json:"nobinat"`
	NATReflection string `json:"natreflection"`
	ControlID     int    `json:"controlID"` //nolint:tagliatelle
}

// NATOneToOne represents a 1:1 NAT (BINAT) rule
type NATOneToOne struct {
	External      string
	Interface     string
	IPProtocol    string
	SourceAddress string
	SourceNot     bool
	DestAddress   string
	DestNot       bool
	Description   string
	Disabled      bool
	NoBinat       bool
	NATReflection string
	controlID     int
}

// IPProtocols returns valid IP protocol options for 1:1 NAT
func (NATOneToOne) IPProtocols() []string {
	return []string{"inet", "inet6", "inet46"}
}

// NATReflectionModes returns valid NAT reflection mode options
func (NATOneToOne) NATReflectionModes() []string {
	return []string{"", "enable", "disable", "purenat"}
}

// SetExternal validates and sets the external IP address
func (r *NATOneToOne) SetExternal(addr string) error {
	if addr == "" {
		return fmt.Errorf("%w, external IP address is required", ErrClientValidation)
	}

	r.External = addr

	return nil
}

// SetInterface validates and sets the interface
func (r *NATOneToOne) SetInterface(iface string) error {
	if err := ValidateInterface(iface); err != nil {
		return fmt.Errorf("%w, invalid interface: %w", ErrClientValidation, err)
	}

	r.Interface = iface

	return nil
}

// SetIPProtocol validates and sets the IP protocol
func (r *NATOneToOne) SetIPProtocol(ipp string) error {
	for _, v := range r.IPProtocols() {
		if ipp == v {
			r.IPProtocol = ipp

			return nil
		}
	}

	return fmt.Errorf("%w, ip protocol must be one of: %s", ErrClientValidation, strings.Join(r.IPProtocols(), ", "))
}

// SetSourceAddress sets the source address
func (r *NATOneToOne) SetSourceAddress(addr string) error {
	r.SourceAddress = addr

	return nil
}

// SetSourceNot sets the source negation flag
func (r *NATOneToOne) SetSourceNot(not bool) error {
	r.SourceNot = not

	return nil
}

// SetDestAddress sets the destination address
func (r *NATOneToOne) SetDestAddress(addr string) error {
	r.DestAddress = addr

	return nil
}

// SetDestNot sets the destination negation flag
func (r *NATOneToOne) SetDestNot(not bool) error {
	r.DestNot = not

	return nil
}

// SetDescription sets the rule description
func (r *NATOneToOne) SetDescription(desc string) error {
	r.Description = desc

	return nil
}

// SetDisabled sets the disabled flag
func (r *NATOneToOne) SetDisabled(disabled bool) error {
	r.Disabled = disabled

	return nil
}

// SetNoBinat sets the nobinat flag (disable 1:1 NAT)
func (r *NATOneToOne) SetNoBinat(nobinat bool) error {
	r.NoBinat = nobinat

	return nil
}

// SetNATReflection validates and sets the NAT reflection mode
func (r *NATOneToOne) SetNATReflection(mode string) error {
	for _, v := range r.NATReflectionModes() {
		if mode == v {
			r.NATReflection = mode

			return nil
		}
	}

	return fmt.Errorf("%w, NAT reflection must be one of: %s (empty string for system default)", ErrClientValidation, strings.Join(r.NATReflectionModes(), ", "))
}

// NATOneToOnes is a slice of NATOneToOne rules
type NATOneToOnes []NATOneToOne

// GetByDescription finds a rule by its description
func (rules NATOneToOnes) GetByDescription(desc string) (*NATOneToOne, error) {
	for _, r := range rules {
		if r.Description == desc {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("NAT 1:1 %w with description '%s'", ErrNotFound, desc)
}

// GetControlIDByDescription finds a rule's control ID by description
func (rules NATOneToOnes) GetControlIDByDescription(desc string) (*int, error) {
	for _, r := range rules {
		if r.Description == desc {
			return &r.controlID, nil
		}
	}

	return nil, fmt.Errorf("NAT 1:1 %w with description '%s'", ErrNotFound, desc)
}

// parseNATOneToOneResponse converts a JSON response to a NATOneToOne struct
func parseNATOneToOneResponse(resp natOneToOneResponse) (NATOneToOne, error) {
	var r NATOneToOne

	// Parse external address
	extAddr, _, _ := parseSourceOrDest(resp.External)
	r.External = extAddr

	if err := r.SetInterface(resp.Interface); err != nil {
		return r, err
	}

	if err := r.SetIPProtocol(resp.IPProtocol); err != nil {
		return r, err
	}

	// Parse source
	srcAddr, _, srcNot := parseSourceOrDest(resp.Source)
	r.SourceAddress = srcAddr
	r.SourceNot = srcNot

	// Parse destination
	dstAddr, _, dstNot := parseSourceOrDest(resp.Destination)
	r.DestAddress = dstAddr
	r.DestNot = dstNot

	if err := r.SetDescription(resp.Description); err != nil {
		return r, err
	}

	r.Disabled = resp.Disabled != ""
	r.NoBinat = resp.NoBinat != ""
	r.NATReflection = resp.NATReflection
	r.controlID = resp.ControlID

	return r, nil
}

// getNATOneToOnes retrieves all 1:1 NAT rules from pfSense
func (pf *Client) getNATOneToOnes(ctx context.Context) (*NATOneToOnes, error) {
	command := "$output = array();" +
		"$rules = config_get_path('nat/onetoone', array());" +
		"foreach ($rules as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var ruleResp []natOneToOneResponse
	if err := pf.executePHPCommand(ctx, command, &ruleResp); err != nil {
		return nil, err
	}

	rules := make(NATOneToOnes, 0, len(ruleResp))
	for _, resp := range ruleResp {
		r, err := parseNATOneToOneResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w NAT 1:1 response, %w", ErrUnableToParse, err)
		}

		rules = append(rules, r)
	}

	return &rules, nil
}

// GetNATOneToOnes retrieves all 1:1 NAT rules
func (pf *Client) GetNATOneToOnes(ctx context.Context) (*NATOneToOnes, error) {
	defer pf.read(&pf.mutexes.NATOneToOne)()

	rules, err := pf.getNATOneToOnes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT 1:1 rules, %w", ErrGetOperationFailed, err)
	}

	return rules, nil
}

// GetNATOneToOne retrieves a specific 1:1 NAT rule by description
func (pf *Client) GetNATOneToOne(ctx context.Context, description string) (*NATOneToOne, error) {
	defer pf.read(&pf.mutexes.NATOneToOne)()

	rules, err := pf.getNATOneToOnes(ctx)
	if err != nil {
		return nil, err
	}

	rule, err := rules.GetByDescription(description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT 1:1 rule, %w", ErrGetOperationFailed, err)
	}

	return rule, nil
}

// CreateNATOneToOne creates a new 1:1 NAT rule
func (pf *Client) CreateNATOneToOne(ctx context.Context, rule *NATOneToOne) error {
	defer pf.write(&pf.mutexes.NATOneToOne)()

	// Build PHP command to add the rule
	command := fmt.Sprintf(
		"$rule = array('external' => '%s', 'interface' => '%s', 'ipprotocol' => '%s', 'source' => array('address' => '%s'), 'destination' => array('address' => '%s'), 'descr' => '%s');"+
			"if (%v) { $rule['disabled'] = ''; }"+
			"if (%v) { $rule['nobinat'] = ''; }"+
			"if ('%s' !== '') { $rule['natreflection'] = '%s'; }"+
			"$rules = config_get_path('nat/onetoone', array()); $rules[] = $rule; config_set_path('nat/onetoone', $rules); write_config(); mark_subsystem_dirty('natconf'); filter_configure(); print('success');",
		phpEscape(rule.External),
		phpEscape(rule.Interface),
		phpEscape(rule.IPProtocol),
		phpEscape(rule.SourceAddress),
		phpEscape(rule.DestAddress),
		phpEscape(rule.Description),
		rule.Disabled,
		rule.NoBinat,
		rule.NATReflection,
		phpEscape(rule.NATReflection),
	)

	var result string
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w NAT 1:1 rule, %w", ErrCreateOperationFailed, err)
	}

	return nil
}

// UpdateNATOneToOne updates an existing 1:1 NAT rule
func (pf *Client) UpdateNATOneToOne(ctx context.Context, oldDescription string, rule *NATOneToOne) error {
	defer pf.write(&pf.mutexes.NATOneToOne)()

	// Build PHP command to update the rule
	command := fmt.Sprintf(
		"$rules = config_get_path('nat/onetoone', array()); "+
			"foreach ($rules as $k => $v) { "+
			"if ($v['descr'] === '%s') { "+
			"$rules[$k] = array('external' => '%s', 'interface' => '%s', 'ipprotocol' => '%s', 'source' => array('address' => '%s'), 'destination' => array('address' => '%s'), 'descr' => '%s'); "+
			"if (%v) { $rules[$k]['disabled'] = ''; } "+
			"if (%v) { $rules[$k]['nobinat'] = ''; } "+
			"if ('%s' !== '') { $rules[$k]['natreflection'] = '%s'; } "+
			"break; "+
			"} "+
			"} "+
			"config_set_path('nat/onetoone', $rules); write_config(); mark_subsystem_dirty('natconf'); filter_configure(); print('success');",
		phpEscape(oldDescription),
		phpEscape(rule.External),
		phpEscape(rule.Interface),
		phpEscape(rule.IPProtocol),
		phpEscape(rule.SourceAddress),
		phpEscape(rule.DestAddress),
		phpEscape(rule.Description),
		rule.Disabled,
		rule.NoBinat,
		rule.NATReflection,
		phpEscape(rule.NATReflection),
	)

	var result string
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w NAT 1:1 rule, %w", ErrUpdateOperationFailed, err)
	}

	return nil
}

// DeleteNATOneToOne deletes a 1:1 NAT rule by description
func (pf *Client) DeleteNATOneToOne(ctx context.Context, description string) error {
	defer pf.write(&pf.mutexes.NATOneToOne)()

	// Build PHP command to delete the rule
	command := fmt.Sprintf(
		"$rules = config_get_path('nat/onetoone', array()); "+
			"foreach ($rules as $k => $v) { "+
			"if ($v['descr'] === '%s') { "+
			"unset($rules[$k]); "+
			"break; "+
			"} "+
			"} "+
			"config_set_path('nat/onetoone', array_values($rules)); write_config(); mark_subsystem_dirty('natconf'); filter_configure(); print('success');",
		phpEscape(description),
	)

	var result string
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w NAT 1:1 rule, %w", ErrDeleteOperationFailed, err)
	}

	return nil
}
