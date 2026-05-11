package pfsense

import (
	"context"
	"fmt"
	"strings"
)

type natOneToOneResponse struct {
	Interface     string `json:"interface"`
	External      string `json:"external"`
	Source        any    `json:"source"`
	Destination   any    `json:"destination"`
	IPProtocol    string `json:"ipprotocol"`
	Disabled      string `json:"disabled"`
	NoBinat       string `json:"nobinat"`
	Description   string `json:"descr"`
	NATReflection string `json:"natreflection"`
	ControlID     int    `json:"controlID"` //nolint:tagliatelle
}

type NATOneToOne struct {
	Interface          string
	External           string
	SourceAddress      string
	SourceNot          bool
	DestinationAddress string
	DestinationNot     bool
	IPProtocol         string
	Disabled           bool
	NoBinat            bool
	Description        string
	NATReflection      string
	ControlID          int
}

type NATOneToOnes []NATOneToOne

func (NATOneToOne) IPProtocols() []string {
	return []string{"inet", "inet6", "inet46"}
}

func (NATOneToOne) NATReflectionModes() []string {
	return []string{"", "enable", "disable", "purenat"}
}

func (r *NATOneToOne) SetInterface(iface string) error {
	if err := ValidateInterface(iface); err != nil {
		return fmt.Errorf("%w, invalid interface: %w", ErrClientValidation, err)
	}

	r.Interface = iface

	return nil
}

func (r *NATOneToOne) SetExternal(ext string) error {
	r.External = ext

	return nil
}

func (r *NATOneToOne) SetIPProtocol(ipp string) error {
	for _, v := range r.IPProtocols() {
		if ipp == v {
			r.IPProtocol = ipp

			return nil
		}
	}

	return fmt.Errorf("%w, IP protocol must be one of: %s", ErrClientValidation, strings.Join(r.IPProtocols(), ", "))
}

func (r *NATOneToOne) SetSourceAddress(addr string) error {
	r.SourceAddress = addr

	return nil
}

func (r *NATOneToOne) SetSourceNot(not bool) error {
	r.SourceNot = not

	return nil
}

func (r *NATOneToOne) SetDestinationAddress(addr string) error {
	r.DestinationAddress = addr

	return nil
}

func (r *NATOneToOne) SetDestinationNot(not bool) error {
	r.DestinationNot = not

	return nil
}

func (r *NATOneToOne) SetDisabled(disabled bool) error {
	r.Disabled = disabled

	return nil
}

func (r *NATOneToOne) SetNoBinat(nobinat bool) error {
	r.NoBinat = nobinat

	return nil
}

func (r *NATOneToOne) SetDescription(desc string) error {
	r.Description = desc

	return nil
}

func (r *NATOneToOne) SetNATReflection(mode string) error {
	for _, v := range r.NATReflectionModes() {
		if mode == v {
			r.NATReflection = mode

			return nil
		}
	}

	return fmt.Errorf("%w, NAT reflection must be one of: %s", ErrClientValidation, strings.Join(r.NATReflectionModes(), ", "))
}

func (rules NATOneToOnes) GetByDescription(desc string) (*NATOneToOne, error) {
	for i := range rules {
		if rules[i].Description == desc {
			return &rules[i], nil
		}
	}

	return nil, fmt.Errorf("%w, NAT 1:1 rule with description %q not found", ErrNotFound, desc)
}

func (rules NATOneToOnes) GetControlIDByDescription(desc string) (int, error) {
	for _, r := range rules {
		if r.Description == desc {
			return r.ControlID, nil
		}
	}

	return -1, fmt.Errorf("%w, NAT 1:1 rule with description %q not found", ErrNotFound, desc)
}

func parseNATOneToOneResponse(resp natOneToOneResponse) (NATOneToOne, error) {
	var r NATOneToOne

	if err := r.SetInterface(resp.Interface); err != nil {
		return r, err
	}

	r.External = resp.External

	if resp.IPProtocol != "" {
		if err := r.SetIPProtocol(resp.IPProtocol); err != nil {
			return r, err
		}
	} else {
		r.IPProtocol = "inet"
	}

	srcAddr, _, srcNot := parseSourceOrDest(resp.Source)
	r.SourceAddress = srcAddr
	r.SourceNot = srcNot

	dstAddr, _, dstNot := parseSourceOrDest(resp.Destination)
	r.DestinationAddress = dstAddr
	r.DestinationNot = dstNot

	r.Disabled = resp.Disabled != ""
	r.NoBinat = resp.NoBinat != ""
	r.Description = resp.Description
	r.NATReflection = resp.NATReflection
	r.ControlID = resp.ControlID

	return r, nil
}

// ====================================================================
// Source/Destination PHP helpers
// ====================================================================

func natOneToOnePHPSource(req NATOneToOne) string {
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
		fmt.Fprintf(&b, "$rule['source']['address'] = '%s';", phpEscape(req.SourceAddress))
	}

	return b.String()
}

func natOneToOnePHPDest(req NATOneToOne) string {
	var b strings.Builder

	if req.DestinationNot {
		b.WriteString("$rule['destination']['not'] = '';")
	}

	if req.DestinationAddress == "" || req.DestinationAddress == "any" {
		b.WriteString("$rule['destination']['any'] = '';")
	} else if strings.Contains(req.DestinationAddress, "/") {
		parts := strings.SplitN(req.DestinationAddress, "/", 2)
		fmt.Fprintf(&b, "$rule['destination']['address'] = '%s/%s';", phpEscape(parts[0]), phpEscape(parts[1]))
	} else if isSpecialAddress(req.DestinationAddress) {
		fmt.Fprintf(&b, "$rule['destination']['network'] = '%s';", phpEscape(req.DestinationAddress))
	} else {
		fmt.Fprintf(&b, "$rule['destination']['address'] = '%s';", phpEscape(req.DestinationAddress))
	}

	return b.String()
}

func natOneToOneBuildRule(req NATOneToOne) string {
	var b strings.Builder

	b.WriteString("$rule = array();")
	fmt.Fprintf(&b, "$rule['interface'] = '%s';", phpEscape(req.Interface))
	fmt.Fprintf(&b, "$rule['external'] = '%s';", phpEscape(req.External))
	fmt.Fprintf(&b, "$rule['ipprotocol'] = '%s';", phpEscape(req.IPProtocol))
	fmt.Fprintf(&b, "$rule['descr'] = '%s';", phpEscape(req.Description))

	b.WriteString("$rule['source'] = array();")
	b.WriteString(natOneToOnePHPSource(req))

	b.WriteString("$rule['destination'] = array();")
	b.WriteString(natOneToOnePHPDest(req))

	if req.Disabled {
		b.WriteString("$rule['disabled'] = '';")
	}

	if req.NoBinat {
		b.WriteString("$rule['nobinat'] = '';")
	}

	if req.NATReflection != "" {
		fmt.Fprintf(&b, "$rule['natreflection'] = '%s';", phpEscape(req.NATReflection))
	}

	return b.String()
}

// ====================================================================
// CRUD operations
// ====================================================================

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
			return nil, fmt.Errorf("%w NAT 1:1 rule response, %w", ErrUnableToParse, err)
		}

		rules = append(rules, r)
	}

	return &rules, nil
}

func (pf *Client) GetNATOneToOnes(ctx context.Context) (*NATOneToOnes, error) {
	defer pf.read(&pf.mutexes.NATOneToOne)()

	rules, err := pf.getNATOneToOnes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT 1:1 rules, %w", ErrGetOperationFailed, err)
	}

	return rules, nil
}

func (pf *Client) GetNATOneToOne(ctx context.Context, description string) (*NATOneToOne, error) {
	defer pf.read(&pf.mutexes.NATOneToOne)()

	rules, err := pf.getNATOneToOnes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT 1:1 rules, %w", ErrGetOperationFailed, err)
	}

	r, err := rules.GetByDescription(description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT 1:1 rule, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) CreateNATOneToOne(ctx context.Context, req NATOneToOne) (*NATOneToOne, error) {
	defer pf.write(&pf.mutexes.NATOneToOne)()

	command := "require_once('filter.inc');" +
		natOneToOneBuildRule(req) +
		"$rules = config_get_path('nat/onetoone', array());" +
		"$rules[] = $rule;" +
		"config_set_path('nat/onetoone', $rules);" +
		"write_config('Terraform: create NAT 1:1 rule');" +
		"mark_subsystem_dirty('natconf');" +
		"filter_configure();" +
		"clear_subsystem_dirty('natconf');" +
		"print(json_encode(array('status' => 'ok')));"

	var result map[string]string
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w NAT 1:1 rule, %w", ErrCreateOperationFailed, err)
	}

	rules, err := pf.getNATOneToOnes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w reading back NAT 1:1 rule, %w", ErrCreateOperationFailed, err)
	}

	r, err := rules.GetByDescription(req.Description)
	if err != nil {
		return nil, fmt.Errorf("%w reading back NAT 1:1 rule, %w", ErrCreateOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) UpdateNATOneToOne(ctx context.Context, description string, req NATOneToOne) (*NATOneToOne, error) {
	defer pf.write(&pf.mutexes.NATOneToOne)()

	rules, err := pf.getNATOneToOnes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT 1:1 rules, %w", ErrUpdateOperationFailed, err)
	}

	controlID, err := rules.GetControlIDByDescription(description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT 1:1 rule, %w", ErrUpdateOperationFailed, err)
	}

	command := "require_once('filter.inc');" +
		natOneToOneBuildRule(req) +
		fmt.Sprintf("config_set_path('nat/onetoone/%d', $rule);", controlID) +
		"write_config('Terraform: update NAT 1:1 rule');" +
		"mark_subsystem_dirty('natconf');" +
		"filter_configure();" +
		"clear_subsystem_dirty('natconf');" +
		"print(json_encode(array('status' => 'ok')));"

	var result map[string]string
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w NAT 1:1 rule, %w", ErrUpdateOperationFailed, err)
	}

	updatedRules, err := pf.getNATOneToOnes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w reading back NAT 1:1 rule, %w", ErrUpdateOperationFailed, err)
	}

	r, err := updatedRules.GetByDescription(req.Description)
	if err != nil {
		return nil, fmt.Errorf("%w reading back NAT 1:1 rule, %w", ErrUpdateOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) DeleteNATOneToOne(ctx context.Context, description string) error {
	defer pf.write(&pf.mutexes.NATOneToOne)()

	rules, err := pf.getNATOneToOnes(ctx)
	if err != nil {
		return fmt.Errorf("%w NAT 1:1 rules, %w", ErrDeleteOperationFailed, err)
	}

	controlID, err := rules.GetControlIDByDescription(description)
	if err != nil {
		return fmt.Errorf("%w NAT 1:1 rule, %w", ErrDeleteOperationFailed, err)
	}

	command := "require_once('filter.inc');" +
		fmt.Sprintf("config_del_path('nat/onetoone/%d');", controlID) +
		"write_config('Terraform: delete NAT 1:1 rule');" +
		"mark_subsystem_dirty('natconf');" +
		"filter_configure();" +
		"clear_subsystem_dirty('natconf');" +
		"print(json_encode(array('status' => 'ok')));"

	var result map[string]string
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w NAT 1:1 rule, %w", ErrDeleteOperationFailed, err)
	}

	return nil
}
