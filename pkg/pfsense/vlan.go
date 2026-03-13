package pfsense

import (
	"context"
	"fmt"
	"regexp"
)

// parentInterfaceRegex matches physical interface names like vmx0, igb1, em0, etc.
var parentInterfaceRegex = regexp.MustCompile(`^[a-zA-Z]+\d+$`)

const (
	MinVLANTag = 1
	MaxVLANTag = 4094
	MinVLANPCP = 0
	MaxVLANPCP = 7
)

type vlanResponse struct {
	ParentInterface string `json:"if"`
	Tag             string `json:"tag"`
	PCP             string `json:"pcp"`
	Description     string `json:"descr"`
	VLANInterface   string `json:"vlanif"`
	ControlID       int    `json:"controlID"` //nolint:tagliatelle
}

type VLAN struct {
	ParentInterface string
	Tag             int
	PCP             *int
	Description     string
	VLANInterface   string
	controlID       int
}

func (v *VLAN) SetParentInterface(iface string) error {
	if iface == "" {
		return fmt.Errorf("%w, parent interface is required", ErrClientValidation)
	}

	if !parentInterfaceRegex.MatchString(iface) {
		return fmt.Errorf("%w, parent interface must be a physical interface name (e.g. 'vmx0', 'igb1')", ErrClientValidation)
	}

	v.ParentInterface = iface

	return nil
}

func (v *VLAN) SetTag(tag int) error {
	if tag < MinVLANTag || tag > MaxVLANTag {
		return fmt.Errorf("%w, VLAN tag must be between %d and %d", ErrClientValidation, MinVLANTag, MaxVLANTag)
	}

	v.Tag = tag

	return nil
}

func (v *VLAN) SetPCP(pcp *int) error {
	if pcp != nil && (*pcp < MinVLANPCP || *pcp > MaxVLANPCP) {
		return fmt.Errorf("%w, VLAN PCP must be between %d and %d", ErrClientValidation, MinVLANPCP, MaxVLANPCP)
	}

	v.PCP = pcp

	return nil
}

func (v *VLAN) SetDescription(description string) error {
	v.Description = description

	return nil
}

type VLANs []VLAN

func (vlans VLANs) GetByInterface(vlanif string) (*VLAN, error) {
	for _, v := range vlans {
		if v.VLANInterface == vlanif {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("VLAN %w with interface '%s'", ErrNotFound, vlanif)
}

func (vlans VLANs) GetControlIDByInterface(vlanif string) (*int, error) {
	for _, v := range vlans {
		if v.VLANInterface == vlanif {
			return &v.controlID, nil
		}
	}

	return nil, fmt.Errorf("VLAN %w with interface '%s'", ErrNotFound, vlanif)
}

func parseVLANResponse(resp vlanResponse) (VLAN, error) {
	var vlan VLAN

	if err := vlan.SetParentInterface(resp.ParentInterface); err != nil {
		return vlan, err
	}

	tag := 0
	if _, err := fmt.Sscanf(resp.Tag, "%d", &tag); err != nil {
		return vlan, fmt.Errorf("%w, unable to parse VLAN tag from '%s'", ErrUnableToParse, resp.Tag)
	}

	if err := vlan.SetTag(tag); err != nil {
		return vlan, err
	}

	if resp.PCP != "" {
		pcp := 0
		if _, err := fmt.Sscanf(resp.PCP, "%d", &pcp); err != nil {
			return vlan, fmt.Errorf("%w, unable to parse VLAN PCP from '%s'", ErrUnableToParse, resp.PCP)
		}

		if err := vlan.SetPCP(&pcp); err != nil {
			return vlan, err
		}
	}

	if err := vlan.SetDescription(resp.Description); err != nil {
		return vlan, err
	}

	vlan.VLANInterface = resp.VLANInterface
	vlan.controlID = resp.ControlID

	return vlan, nil
}

func (pf *Client) getVLANs(ctx context.Context) (*VLANs, error) {
	command := "$output = array();" +
		"$vlans = config_get_path('vlans/vlan', array());" +
		"foreach ($vlans as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var vlanResp []vlanResponse
	if err := pf.executePHPCommand(ctx, command, &vlanResp); err != nil {
		return nil, err
	}

	vlans := make(VLANs, 0, len(vlanResp))
	for _, resp := range vlanResp {
		v, err := parseVLANResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w VLAN response, %w", ErrUnableToParse, err)
		}

		vlans = append(vlans, v)
	}

	return &vlans, nil
}

func (pf *Client) GetVLANs(ctx context.Context) (*VLANs, error) {
	defer pf.read(&pf.mutexes.VLAN)()

	vlans, err := pf.getVLANs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w VLANs, %w", ErrGetOperationFailed, err)
	}

	return vlans, nil
}

func (pf *Client) GetVLAN(ctx context.Context, vlanif string) (*VLAN, error) {
	defer pf.read(&pf.mutexes.VLAN)()

	vlans, err := pf.getVLANs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w VLANs, %w", ErrGetOperationFailed, err)
	}

	v, err := vlans.GetByInterface(vlanif)
	if err != nil {
		return nil, fmt.Errorf("%w VLAN, %w", ErrGetOperationFailed, err)
	}

	return v, nil
}

func (pf *Client) CreateVLAN(ctx context.Context, req VLAN) (*VLAN, error) {
	defer pf.write(&pf.mutexes.VLAN)()

	vlanif := fmt.Sprintf("%s.%d", req.ParentInterface, req.Tag)

	// Check for duplicate VLAN (same parent interface + tag).
	existingVLANs, err := pf.getVLANs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w VLANs for duplicate check, %w", ErrGetOperationFailed, err)
	}

	if _, err := existingVLANs.GetByInterface(vlanif); err == nil {
		return nil, fmt.Errorf("%w VLAN, a VLAN with parent '%s' and tag %d already exists (%s)", ErrCreateOperationFailed, req.ParentInterface, req.Tag, vlanif)
	}

	pcpStr := ""
	if req.PCP != nil {
		pcpStr = fmt.Sprintf("%d", *req.PCP)
	}

	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"$vlan = array();"+
			"$vlan['if'] = '%s';"+
			"$vlan['tag'] = '%d';"+
			"$vlan['pcp'] = '%s';"+
			"$vlan['descr'] = '%s';"+
			"$vlan['vlanif'] = '%s';"+
			"config_set_path('vlans/vlan/', $vlan);"+
			"interface_vlan_configure($vlan);"+
			"write_config('Terraform: created VLAN %s');"+
			"print(json_encode(true));",
		phpEscape(req.ParentInterface),
		req.Tag,
		phpEscape(pcpStr),
		phpEscape(req.Description),
		phpEscape(vlanif),
		phpEscape(vlanif),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w VLAN, %w", ErrCreateOperationFailed, err)
	}

	vlans, err := pf.getVLANs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w VLANs after creating, %w", ErrGetOperationFailed, err)
	}

	v, err := vlans.GetByInterface(vlanif)
	if err != nil {
		return nil, fmt.Errorf("%w VLAN after creating, %w", ErrGetOperationFailed, err)
	}

	return v, nil
}

func (pf *Client) UpdateVLAN(ctx context.Context, req VLAN) (*VLAN, error) {
	defer pf.write(&pf.mutexes.VLAN)()

	vlanif := fmt.Sprintf("%s.%d", req.ParentInterface, req.Tag)

	vlans, err := pf.getVLANs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w VLANs, %w", ErrGetOperationFailed, err)
	}

	controlID, err := vlans.GetControlIDByInterface(vlanif)
	if err != nil {
		return nil, fmt.Errorf("%w VLAN, %w", ErrGetOperationFailed, err)
	}

	pcpStr := ""
	if req.PCP != nil {
		pcpStr = fmt.Sprintf("%d", *req.PCP)
	}

	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"$vlan = array();"+
			"$vlan['if'] = '%s';"+
			"$vlan['tag'] = '%d';"+
			"$vlan['pcp'] = '%s';"+
			"$vlan['descr'] = '%s';"+
			"$vlan['vlanif'] = '%s';"+
			"config_set_path('vlans/vlan/%d', $vlan);"+
			"interface_vlan_configure($vlan);"+
			"write_config('Terraform: updated VLAN %s');"+
			"print(json_encode(true));",
		phpEscape(req.ParentInterface),
		req.Tag,
		phpEscape(pcpStr),
		phpEscape(req.Description),
		phpEscape(vlanif),
		*controlID,
		phpEscape(vlanif),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w VLAN, %w", ErrUpdateOperationFailed, err)
	}

	vlans, err = pf.getVLANs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w VLANs after updating, %w", ErrGetOperationFailed, err)
	}

	v, err := vlans.GetByInterface(vlanif)
	if err != nil {
		return nil, fmt.Errorf("%w VLAN after updating, %w", ErrGetOperationFailed, err)
	}

	return v, nil
}

func (pf *Client) DeleteVLAN(ctx context.Context, vlanif string) error {
	defer pf.write(&pf.mutexes.VLAN)()

	vlans, err := pf.getVLANs(ctx)
	if err != nil {
		return fmt.Errorf("%w VLANs, %w", ErrGetOperationFailed, err)
	}

	controlID, err := vlans.GetControlIDByInterface(vlanif)
	if err != nil {
		return fmt.Errorf("%w VLAN, %w", ErrGetOperationFailed, err)
	}

	// Check if VLAN is assigned to an interface before deleting.
	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"$vlanif = '%s';"+
			"$assigned = false;"+
			"foreach (config_get_path('interfaces', array()) as $ifname => $iface) {"+
			"if (isset($iface['if']) && $iface['if'] === $vlanif) { $assigned = true; break; }"+
			"}"+
			"if ($assigned) {"+
			"print(json_encode('Cannot delete VLAN that is assigned to an interface'));"+
			"} else {"+
			"pfSense_interface_destroy($vlanif);"+
			"config_del_path('vlans/vlan/%d');"+
			"write_config('Terraform: deleted VLAN %s');"+
			"print(json_encode(true));"+
			"}",
		phpEscape(vlanif),
		*controlID,
		phpEscape(vlanif),
	)

	var result interface{}
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w VLAN, %w", ErrDeleteOperationFailed, err)
	}

	if errMsg, ok := result.(string); ok {
		return fmt.Errorf("%w VLAN, %s", ErrDeleteOperationFailed, errMsg)
	}

	// Verify deletion.
	vlans, err = pf.getVLANs(ctx)
	if err != nil {
		return fmt.Errorf("%w VLANs after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := vlans.GetByInterface(vlanif); err == nil {
		return fmt.Errorf("%w VLAN, still exists", ErrDeleteOperationFailed)
	}

	return nil
}
