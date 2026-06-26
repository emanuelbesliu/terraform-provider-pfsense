package pfsense

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type natNPtResponse struct {
	Interface   string          `json:"interface"`
	Source      any             `json:"source"`
	Destination any             `json:"destination"`
	Disabled    json.RawMessage `json:"disabled"`
	Description string          `json:"descr"`
	ControlID   int             `json:"controlID"` //nolint:tagliatelle
}

type NATNPt struct {
	Interface         string
	SourcePrefix      string
	SourceNot         bool
	DestinationPrefix string
	DestinationNot    bool
	Disabled          bool
	Description       string
	ControlID         int
}

type NATNPts []NATNPt

func (r *NATNPt) SetInterface(iface string) error {
	if err := ValidateInterface(iface); err != nil {
		return fmt.Errorf("%w, invalid interface: %w", ErrClientValidation, err)
	}

	r.Interface = iface

	return nil
}

func (r *NATNPt) SetSourcePrefix(prefix string) error {
	r.SourcePrefix = prefix

	return nil
}

func (r *NATNPt) SetSourceNot(not bool) error {
	r.SourceNot = not

	return nil
}

func (r *NATNPt) SetDestinationPrefix(prefix string) error {
	r.DestinationPrefix = prefix

	return nil
}

func (r *NATNPt) SetDestinationNot(not bool) error {
	r.DestinationNot = not

	return nil
}

func (r *NATNPt) SetDisabled(disabled bool) error {
	r.Disabled = disabled

	return nil
}

func (r *NATNPt) SetDescription(desc string) error {
	r.Description = desc

	return nil
}

func (rules NATNPts) GetByDescription(desc string) (*NATNPt, error) {
	for i := range rules {
		if rules[i].Description == desc {
			return &rules[i], nil
		}
	}

	return nil, fmt.Errorf("%w, NAT NPt rule with description %q not found", ErrNotFound, desc)
}

func (rules NATNPts) GetControlIDByDescription(desc string) (int, error) {
	for _, r := range rules {
		if r.Description == desc {
			return r.ControlID, nil
		}
	}

	return -1, fmt.Errorf("%w, NAT NPt rule with description %q not found", ErrNotFound, desc)
}

func parseNATNPtResponse(resp natNPtResponse) (NATNPt, error) {
	var r NATNPt

	if err := r.SetInterface(resp.Interface); err != nil {
		return r, err
	}

	srcAddr, _, srcNot := parseSourceOrDest(resp.Source)
	r.SourcePrefix = srcAddr
	r.SourceNot = srcNot

	dstAddr, _, dstNot := parseSourceOrDest(resp.Destination)
	r.DestinationPrefix = dstAddr
	r.DestinationNot = dstNot

	r.Disabled = rawIsPresent(resp.Disabled)
	r.Description = resp.Description
	r.ControlID = resp.ControlID

	return r, nil
}

// ====================================================================
// Source/Destination PHP helpers
// ====================================================================

// natNPtPHPAddress builds the PHP for a source or destination IPv6 prefix.
// A value containing "/" is stored as a CIDR address; otherwise it is treated
// as a special network (e.g. a track6 interface delegated prefix).
func natNPtPHPAddress(field string, prefix string, not bool) string {
	var b strings.Builder

	if not {
		fmt.Fprintf(&b, "$rule['%s']['not'] = '';", field)
	}

	if strings.Contains(prefix, "/") {
		fmt.Fprintf(&b, "$rule['%s']['address'] = '%s';", field, phpEscape(prefix))
	} else {
		fmt.Fprintf(&b, "$rule['%s']['network'] = '%s';", field, phpEscape(prefix))
	}

	return b.String()
}

func natNPtBuildRule(req NATNPt) string {
	var b strings.Builder

	b.WriteString("$rule = array();")
	fmt.Fprintf(&b, "$rule['interface'] = '%s';", phpEscape(req.Interface))
	fmt.Fprintf(&b, "$rule['descr'] = '%s';", phpEscape(req.Description))

	b.WriteString("$rule['source'] = array();")
	b.WriteString(natNPtPHPAddress("source", req.SourcePrefix, req.SourceNot))

	b.WriteString("$rule['destination'] = array();")
	b.WriteString(natNPtPHPAddress("destination", req.DestinationPrefix, req.DestinationNot))

	if req.Disabled {
		b.WriteString("$rule['disabled'] = '';")
	}

	return b.String()
}

// ====================================================================
// CRUD operations
// ====================================================================

func (pf *Client) getNATNPts(ctx context.Context) (*NATNPts, error) {
	command := "$output = array();" +
		"$rules = config_get_path('nat/npt', array());" +
		"foreach ($rules as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var ruleResp []natNPtResponse
	if err := pf.executePHPCommand(ctx, command, &ruleResp); err != nil {
		return nil, err
	}

	rules := make(NATNPts, 0, len(ruleResp))
	for _, resp := range ruleResp {
		r, err := parseNATNPtResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w NAT NPt rule response, %w", ErrUnableToParse, err)
		}

		rules = append(rules, r)
	}

	return &rules, nil
}

func (pf *Client) GetNATNPts(ctx context.Context) (*NATNPts, error) {
	defer pf.read(&pf.mutexes.NATNPt)()

	rules, err := pf.getNATNPts(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT NPt rules, %w", ErrGetOperationFailed, err)
	}

	return rules, nil
}

func (pf *Client) GetNATNPt(ctx context.Context, description string) (*NATNPt, error) {
	defer pf.read(&pf.mutexes.NATNPt)()

	rules, err := pf.getNATNPts(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT NPt rules, %w", ErrGetOperationFailed, err)
	}

	r, err := rules.GetByDescription(description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT NPt rule, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) CreateNATNPt(ctx context.Context, req NATNPt) (*NATNPt, error) {
	defer pf.write(&pf.mutexes.NATNPt)()

	command := "require_once('filter.inc');" +
		natNPtBuildRule(req) +
		"$rules = config_get_path('nat/npt', array());" +
		"$rules[] = $rule;" +
		"config_set_path('nat/npt', $rules);" +
		"write_config('Terraform: create NAT NPt rule');" +
		"mark_subsystem_dirty('natconf');" +
		"filter_configure();" +
		"clear_subsystem_dirty('natconf');" +
		"print(json_encode(array('status' => 'ok')));"

	var result map[string]string
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w NAT NPt rule, %w", ErrCreateOperationFailed, err)
	}

	rules, err := pf.getNATNPts(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w reading back NAT NPt rule, %w", ErrCreateOperationFailed, err)
	}

	r, err := rules.GetByDescription(req.Description)
	if err != nil {
		return nil, fmt.Errorf("%w reading back NAT NPt rule, %w", ErrCreateOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) UpdateNATNPt(ctx context.Context, description string, req NATNPt) (*NATNPt, error) {
	defer pf.write(&pf.mutexes.NATNPt)()

	rules, err := pf.getNATNPts(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w NAT NPt rules, %w", ErrUpdateOperationFailed, err)
	}

	controlID, err := rules.GetControlIDByDescription(description)
	if err != nil {
		return nil, fmt.Errorf("%w NAT NPt rule, %w", ErrUpdateOperationFailed, err)
	}

	command := "require_once('filter.inc');" +
		natNPtBuildRule(req) +
		fmt.Sprintf("config_set_path('nat/npt/%d', $rule);", controlID) +
		"write_config('Terraform: update NAT NPt rule');" +
		"mark_subsystem_dirty('natconf');" +
		"filter_configure();" +
		"clear_subsystem_dirty('natconf');" +
		"print(json_encode(array('status' => 'ok')));"

	var result map[string]string
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w NAT NPt rule, %w", ErrUpdateOperationFailed, err)
	}

	updatedRules, err := pf.getNATNPts(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w reading back NAT NPt rule, %w", ErrUpdateOperationFailed, err)
	}

	r, err := updatedRules.GetByDescription(req.Description)
	if err != nil {
		return nil, fmt.Errorf("%w reading back NAT NPt rule, %w", ErrUpdateOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) DeleteNATNPt(ctx context.Context, description string) error {
	defer pf.write(&pf.mutexes.NATNPt)()

	rules, err := pf.getNATNPts(ctx)
	if err != nil {
		return fmt.Errorf("%w NAT NPt rules, %w", ErrDeleteOperationFailed, err)
	}

	controlID, err := rules.GetControlIDByDescription(description)
	if err != nil {
		return fmt.Errorf("%w NAT NPt rule, %w", ErrDeleteOperationFailed, err)
	}

	command := "require_once('filter.inc');" +
		fmt.Sprintf("config_del_path('nat/npt/%d');", controlID) +
		"write_config('Terraform: delete NAT NPt rule');" +
		"mark_subsystem_dirty('natconf');" +
		"filter_configure();" +
		"clear_subsystem_dirty('natconf');" +
		"print(json_encode(array('status' => 'ok')));"

	var result map[string]string
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w NAT NPt rule, %w", ErrDeleteOperationFailed, err)
	}

	return nil
}
