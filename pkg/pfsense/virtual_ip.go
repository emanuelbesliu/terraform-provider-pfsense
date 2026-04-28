package pfsense

import (
	"context"
	"fmt"
)

const (
	VirtualIPModeIPAlias  = "ipalias"
	VirtualIPModeCarp     = "carp"
	VirtualIPModeProxyARP = "proxyarp"
	VirtualIPModeOther    = "other"
	MinVHID               = 1
	MaxVHID               = 255
	MinAdvSkew            = 0
	MaxAdvSkew            = 254
	MinAdvBase            = 1
	MaxAdvBase            = 254
)

var ValidVirtualIPModes = []string{
	VirtualIPModeIPAlias,
	VirtualIPModeCarp,
	VirtualIPModeProxyARP,
	VirtualIPModeOther,
}

type virtualIPResponse struct {
	Mode        string `json:"mode"`
	Interface   string `json:"interface"`
	VHID        string `json:"vhid"`
	AdvSkew     string `json:"advskew"`
	AdvBase     string `json:"advbase"`
	Password    string `json:"password"`
	Subnet      string `json:"subnet"`
	SubnetBits  string `json:"subnet_bits"`
	Description string `json:"descr"`
	Type        string `json:"type"`
	UniqueID    string `json:"uniqid"`
	ControlID   int    `json:"controlID"` //nolint:tagliatelle
}

type VirtualIP struct {
	Mode        string
	Interface   string
	VHID        *int
	AdvSkew     *int
	AdvBase     *int
	Password    string
	Subnet      string
	SubnetBits  int
	Description string
	Type        string
	UniqueID    string
	controlID   int
}

func (v *VirtualIP) SetMode(mode string) error {
	for _, valid := range ValidVirtualIPModes {
		if mode == valid {
			v.Mode = mode

			return nil
		}
	}

	return fmt.Errorf("%w, mode must be one of: ipalias, carp, proxyarp, other", ErrClientValidation)
}

func (v *VirtualIP) SetInterface(iface string) error {
	if iface == "" {
		return fmt.Errorf("%w, interface is required", ErrClientValidation)
	}

	v.Interface = iface

	return nil
}

func (v *VirtualIP) SetVHID(vhid *int) error {
	if vhid != nil && (*vhid < MinVHID || *vhid > MaxVHID) {
		return fmt.Errorf("%w, VHID must be between %d and %d", ErrClientValidation, MinVHID, MaxVHID)
	}

	v.VHID = vhid

	return nil
}

func (v *VirtualIP) SetAdvSkew(advskew *int) error {
	if advskew != nil && (*advskew < MinAdvSkew || *advskew > MaxAdvSkew) {
		return fmt.Errorf("%w, advertisement skew must be between %d and %d", ErrClientValidation, MinAdvSkew, MaxAdvSkew)
	}

	v.AdvSkew = advskew

	return nil
}

func (v *VirtualIP) SetAdvBase(advbase *int) error {
	if advbase != nil && (*advbase < MinAdvBase || *advbase > MaxAdvBase) {
		return fmt.Errorf("%w, advertisement base must be between %d and %d", ErrClientValidation, MinAdvBase, MaxAdvBase)
	}

	v.AdvBase = advbase

	return nil
}

func (v *VirtualIP) SetPassword(password string) error {
	v.Password = password

	return nil
}

func (v *VirtualIP) SetSubnet(subnet string) error {
	if subnet == "" {
		return fmt.Errorf("%w, subnet (IP address) is required", ErrClientValidation)
	}

	v.Subnet = subnet

	return nil
}

func (v *VirtualIP) SetSubnetBits(bits int) error {
	if bits < 1 || bits > 128 {
		return fmt.Errorf("%w, subnet bits must be between 1 and 128", ErrClientValidation)
	}

	v.SubnetBits = bits

	return nil
}

func (v *VirtualIP) SetDescription(description string) error {
	v.Description = description

	return nil
}

type VirtualIPs []VirtualIP

func (vips VirtualIPs) GetByUniqueID(uniqid string) (*VirtualIP, error) {
	for _, v := range vips {
		if v.UniqueID == uniqid {
			return &v, nil
		}
	}

	return nil, fmt.Errorf("Virtual IP %w with unique ID '%s'", ErrNotFound, uniqid)
}

func (vips VirtualIPs) GetControlIDByUniqueID(uniqid string) (*int, error) {
	for _, v := range vips {
		if v.UniqueID == uniqid {
			return &v.controlID, nil
		}
	}

	return nil, fmt.Errorf("Virtual IP %w with unique ID '%s'", ErrNotFound, uniqid)
}

func parseVirtualIPResponse(resp virtualIPResponse) (VirtualIP, error) {
	var vip VirtualIP

	if err := vip.SetMode(resp.Mode); err != nil {
		return vip, err
	}

	if err := vip.SetInterface(resp.Interface); err != nil {
		return vip, err
	}

	if resp.VHID != "" {
		vhid := 0
		if _, err := fmt.Sscanf(resp.VHID, "%d", &vhid); err != nil {
			return vip, fmt.Errorf("%w, unable to parse VHID from '%s'", ErrUnableToParse, resp.VHID)
		}

		if err := vip.SetVHID(&vhid); err != nil {
			return vip, err
		}
	}

	if resp.AdvSkew != "" {
		advskew := 0
		if _, err := fmt.Sscanf(resp.AdvSkew, "%d", &advskew); err != nil {
			return vip, fmt.Errorf("%w, unable to parse advertisement skew from '%s'", ErrUnableToParse, resp.AdvSkew)
		}

		if err := vip.SetAdvSkew(&advskew); err != nil {
			return vip, err
		}
	}

	if resp.AdvBase != "" {
		advbase := 0
		if _, err := fmt.Sscanf(resp.AdvBase, "%d", &advbase); err != nil {
			return vip, fmt.Errorf("%w, unable to parse advertisement base from '%s'", ErrUnableToParse, resp.AdvBase)
		}

		if err := vip.SetAdvBase(&advbase); err != nil {
			return vip, err
		}
	}

	if err := vip.SetPassword(resp.Password); err != nil {
		return vip, err
	}

	if err := vip.SetSubnet(resp.Subnet); err != nil {
		return vip, err
	}

	bits := 0
	if _, err := fmt.Sscanf(resp.SubnetBits, "%d", &bits); err != nil {
		return vip, fmt.Errorf("%w, unable to parse subnet bits from '%s'", ErrUnableToParse, resp.SubnetBits)
	}

	if err := vip.SetSubnetBits(bits); err != nil {
		return vip, err
	}

	if err := vip.SetDescription(resp.Description); err != nil {
		return vip, err
	}

	vip.Type = resp.Type
	vip.UniqueID = resp.UniqueID
	vip.controlID = resp.ControlID

	return vip, nil
}

func (pf *Client) getVirtualIPs(ctx context.Context) (*VirtualIPs, error) {
	command := "$output = array();" +
		"$vips = config_get_path('virtualip/vip', array());" +
		"foreach ($vips as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var vipResp []virtualIPResponse
	if err := pf.executePHPCommand(ctx, command, &vipResp); err != nil {
		return nil, err
	}

	vips := make(VirtualIPs, 0, len(vipResp))
	for _, resp := range vipResp {
		v, err := parseVirtualIPResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w Virtual IP response, %w", ErrUnableToParse, err)
		}

		vips = append(vips, v)
	}

	return &vips, nil
}

func (pf *Client) GetVirtualIPs(ctx context.Context) (*VirtualIPs, error) {
	defer pf.read(&pf.mutexes.VirtualIP)()

	vips, err := pf.getVirtualIPs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w Virtual IPs, %w", ErrGetOperationFailed, err)
	}

	return vips, nil
}

func (pf *Client) GetVirtualIP(ctx context.Context, uniqid string) (*VirtualIP, error) {
	defer pf.read(&pf.mutexes.VirtualIP)()

	vips, err := pf.getVirtualIPs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w Virtual IPs, %w", ErrGetOperationFailed, err)
	}

	v, err := vips.GetByUniqueID(uniqid)
	if err != nil {
		return nil, fmt.Errorf("%w Virtual IP, %w", ErrGetOperationFailed, err)
	}

	return v, nil
}

func (pf *Client) CreateVirtualIP(ctx context.Context, req VirtualIP) (*VirtualIP, error) {
	defer pf.write(&pf.mutexes.VirtualIP)()

	vhidStr := ""
	if req.VHID != nil {
		vhidStr = fmt.Sprintf("%d", *req.VHID)
	}

	advskewStr := ""
	if req.AdvSkew != nil {
		advskewStr = fmt.Sprintf("%d", *req.AdvSkew)
	}

	advbaseStr := ""
	if req.AdvBase != nil {
		advbaseStr = fmt.Sprintf("%d", *req.AdvBase)
	}

	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"require_once('filter.inc');"+
			"$vip = array();"+
			"$vip['mode'] = '%s';"+
			"$vip['interface'] = '%s';"+
			"$vip['vhid'] = '%s';"+
			"$vip['advskew'] = '%s';"+
			"$vip['advbase'] = '%s';"+
			"$vip['password'] = '%s';"+
			"$vip['subnet'] = '%s';"+
			"$vip['subnet_bits'] = '%d';"+
			"$vip['descr'] = '%s';"+
			"$vip['type'] = 'single';"+
			"$vip['uniqid'] = uniqid('', true);"+
			"config_set_path('virtualip/vip/', $vip);"+
			"if ($vip['mode'] == 'ipalias') { interface_ipalias_configure($vip); }"+
			"if ($vip['mode'] == 'carp') { interface_carp_configure($vip); }"+
			"if ($vip['mode'] == 'proxyarp') { interface_proxyarp_configure($vip['interface']); }"+
			"write_config('Terraform: created Virtual IP %s on %s');"+
			"filter_configure();"+
			"print(json_encode($vip['uniqid']));",
		phpEscape(req.Mode),
		phpEscape(req.Interface),
		phpEscape(vhidStr),
		phpEscape(advskewStr),
		phpEscape(advbaseStr),
		phpEscape(req.Password),
		phpEscape(req.Subnet),
		req.SubnetBits,
		phpEscape(req.Description),
		phpEscape(req.Subnet),
		phpEscape(req.Interface),
	)

	var uniqid string
	if err := pf.executePHPCommand(ctx, command, &uniqid); err != nil {
		return nil, fmt.Errorf("%w Virtual IP, %w", ErrCreateOperationFailed, err)
	}

	vips, err := pf.getVirtualIPs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w Virtual IPs after creating, %w", ErrGetOperationFailed, err)
	}

	v, err := vips.GetByUniqueID(uniqid)
	if err != nil {
		return nil, fmt.Errorf("%w Virtual IP after creating, %w", ErrGetOperationFailed, err)
	}

	return v, nil
}

func (pf *Client) UpdateVirtualIP(ctx context.Context, req VirtualIP) (*VirtualIP, error) {
	defer pf.write(&pf.mutexes.VirtualIP)()

	vips, err := pf.getVirtualIPs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w Virtual IPs, %w", ErrGetOperationFailed, err)
	}

	controlID, err := vips.GetControlIDByUniqueID(req.UniqueID)
	if err != nil {
		return nil, fmt.Errorf("%w Virtual IP, %w", ErrGetOperationFailed, err)
	}

	vhidStr := ""
	if req.VHID != nil {
		vhidStr = fmt.Sprintf("%d", *req.VHID)
	}

	advskewStr := ""
	if req.AdvSkew != nil {
		advskewStr = fmt.Sprintf("%d", *req.AdvSkew)
	}

	advbaseStr := ""
	if req.AdvBase != nil {
		advbaseStr = fmt.Sprintf("%d", *req.AdvBase)
	}

	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"require_once('filter.inc');"+
			"$oldvip = config_get_path('virtualip/vip/%d', array());"+
			"interface_vip_bring_down($oldvip);"+
			"$vip = array();"+
			"$vip['mode'] = '%s';"+
			"$vip['interface'] = '%s';"+
			"$vip['vhid'] = '%s';"+
			"$vip['advskew'] = '%s';"+
			"$vip['advbase'] = '%s';"+
			"$vip['password'] = '%s';"+
			"$vip['subnet'] = '%s';"+
			"$vip['subnet_bits'] = '%d';"+
			"$vip['descr'] = '%s';"+
			"$vip['type'] = 'single';"+
			"$vip['uniqid'] = '%s';"+
			"config_set_path('virtualip/vip/%d', $vip);"+
			"if ($vip['mode'] == 'ipalias') { interface_ipalias_configure($vip); }"+
			"if ($vip['mode'] == 'carp') { interface_carp_configure($vip); }"+
			"if ($vip['mode'] == 'proxyarp') { interface_proxyarp_configure($vip['interface']); }"+
			"write_config('Terraform: updated Virtual IP %s on %s');"+
			"filter_configure();"+
			"print(json_encode(true));",
		*controlID,
		phpEscape(req.Mode),
		phpEscape(req.Interface),
		phpEscape(vhidStr),
		phpEscape(advskewStr),
		phpEscape(advbaseStr),
		phpEscape(req.Password),
		phpEscape(req.Subnet),
		req.SubnetBits,
		phpEscape(req.Description),
		phpEscape(req.UniqueID),
		*controlID,
		phpEscape(req.Subnet),
		phpEscape(req.Interface),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w Virtual IP, %w", ErrUpdateOperationFailed, err)
	}

	vips, err = pf.getVirtualIPs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w Virtual IPs after updating, %w", ErrGetOperationFailed, err)
	}

	v, err := vips.GetByUniqueID(req.UniqueID)
	if err != nil {
		return nil, fmt.Errorf("%w Virtual IP after updating, %w", ErrGetOperationFailed, err)
	}

	return v, nil
}

func (pf *Client) DeleteVirtualIP(ctx context.Context, uniqid string) error {
	defer pf.write(&pf.mutexes.VirtualIP)()

	vips, err := pf.getVirtualIPs(ctx)
	if err != nil {
		return fmt.Errorf("%w Virtual IPs, %w", ErrGetOperationFailed, err)
	}

	controlID, err := vips.GetControlIDByUniqueID(uniqid)
	if err != nil {
		return fmt.Errorf("%w Virtual IP, %w", ErrGetOperationFailed, err)
	}

	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"require_once('filter.inc');"+
			"$vip = config_get_path('virtualip/vip/%d', array());"+
			"interface_vip_bring_down($vip);"+
			"config_del_path('virtualip/vip/%d');"+
			"write_config('Terraform: deleted Virtual IP %s');"+
			"filter_configure();"+
			"print(json_encode(true));",
		*controlID,
		*controlID,
		phpEscape(uniqid),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w Virtual IP, %w", ErrDeleteOperationFailed, err)
	}

	vips, err = pf.getVirtualIPs(ctx)
	if err != nil {
		return fmt.Errorf("%w Virtual IPs after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := vips.GetByUniqueID(uniqid); err == nil {
		return fmt.Errorf("%w Virtual IP, still exists", ErrDeleteOperationFailed)
	}

	return nil
}
