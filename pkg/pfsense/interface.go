package pfsense

import (
	"context"
	"fmt"
	"strings"
)

// IPv4 configuration types for interface assignments.
var InterfaceIPv4Types = []string{"none", "staticv4", "dhcp", "pppoe", "pptp", "l2tp"}

// IPv6 configuration types for interface assignments.
var InterfaceIPv6Types = []string{"none", "staticv6", "dhcp6", "slaac", "6to4", "6rd", "track6"}

type interfaceResponse struct {
	If          string  `json:"if"`
	Descr       string  `json:"descr"`
	Enable      *string `json:"enable"`
	IPAddr      string  `json:"ipaddr"`
	Subnet      string  `json:"subnet"`
	Gateway     string  `json:"gateway"`
	IPAddrV6    string  `json:"ipaddrv6"`
	SubnetV6    string  `json:"subnetv6"`
	GatewayV6   string  `json:"gatewayv6"`
	SpoofMAC    string  `json:"spoofmac"`
	MTU         string  `json:"mtu"`
	MSS         string  `json:"mss"`
	Media       string  `json:"media"`
	BlockPriv   *string `json:"blockpriv"`
	BlockBogons *string `json:"blockbogons"`
}

// Interface represents a pfSense interface assignment with its configuration.
type Interface struct {
	LogicalName string
	If          string
	Description string
	Enabled     bool
	IPv4Type    string
	IPAddr      string
	Subnet      string
	Gateway     string
	IPv6Type    string
	IPAddrV6    string
	SubnetV6    string
	GatewayV6   string
	SpoofMAC    string
	MTU         int
	MSS         int
	Media       string
	BlockPriv   bool
	BlockBogons bool
}

func (iface *Interface) SetIf(ifName string) error {
	if ifName == "" {
		return fmt.Errorf("%w, interface port (if) is required", ErrClientValidation)
	}

	iface.If = ifName

	return nil
}

func (iface *Interface) SetDescription(description string) error {
	iface.Description = description

	return nil
}

func (iface *Interface) SetEnabled(enabled bool) error {
	iface.Enabled = enabled

	return nil
}

func (iface *Interface) SetIPv4Type(ipv4Type string) error {
	if ipv4Type == "" {
		ipv4Type = "none"
	}

	for _, t := range InterfaceIPv4Types {
		if ipv4Type == t {
			iface.IPv4Type = ipv4Type

			return nil
		}
	}

	return fmt.Errorf("%w, IPv4 type must be one of: %s", ErrClientValidation, strings.Join(InterfaceIPv4Types, ", "))
}

func (iface *Interface) SetIPv6Type(ipv6Type string) error {
	if ipv6Type == "" {
		ipv6Type = "none"
	}

	for _, t := range InterfaceIPv6Types {
		if ipv6Type == t {
			iface.IPv6Type = ipv6Type

			return nil
		}
	}

	return fmt.Errorf("%w, IPv6 type must be one of: %s", ErrClientValidation, strings.Join(InterfaceIPv6Types, ", "))
}

type Interfaces []Interface

func (ifaces Interfaces) GetByLogicalName(name string) (*Interface, error) {
	for _, i := range ifaces {
		if i.LogicalName == name {
			return &i, nil
		}
	}

	return nil, fmt.Errorf("interface %w with logical name '%s'", ErrNotFound, name)
}

func parseInterfaceResponse(logicalName string, resp interfaceResponse) (Interface, error) {
	var iface Interface

	iface.LogicalName = logicalName

	if err := iface.SetIf(resp.If); err != nil {
		return iface, err
	}

	if err := iface.SetDescription(resp.Descr); err != nil {
		return iface, err
	}

	iface.Enabled = resp.Enable != nil

	// Determine IPv4 type from ipaddr value.
	switch resp.IPAddr {
	case "dhcp":
		iface.IPv4Type = "dhcp"
	case "pppoe":
		iface.IPv4Type = "pppoe"
	case "pptp":
		iface.IPv4Type = "pptp"
	case "l2tp":
		iface.IPv4Type = "l2tp"
	case "", "none":
		iface.IPv4Type = "none"
	default:
		iface.IPv4Type = "staticv4"
		iface.IPAddr = resp.IPAddr
	}

	if resp.Subnet != "" {
		subnet := 0
		if _, err := fmt.Sscanf(resp.Subnet, "%d", &subnet); err == nil {
			iface.Subnet = resp.Subnet
		}
	}

	iface.Gateway = resp.Gateway

	// Determine IPv6 type from ipaddrv6 value.
	switch resp.IPAddrV6 {
	case "dhcp6":
		iface.IPv6Type = "dhcp6"
	case "slaac":
		iface.IPv6Type = "slaac"
	case "6to4":
		iface.IPv6Type = "6to4"
	case "6rd":
		iface.IPv6Type = "6rd"
	case "track6":
		iface.IPv6Type = "track6"
	case "", "none":
		iface.IPv6Type = "none"
	default:
		iface.IPv6Type = "staticv6"
		iface.IPAddrV6 = resp.IPAddrV6
	}

	if resp.SubnetV6 != "" {
		iface.SubnetV6 = resp.SubnetV6
	}

	iface.GatewayV6 = resp.GatewayV6
	iface.SpoofMAC = resp.SpoofMAC

	if resp.MTU != "" {
		mtu := 0
		if _, err := fmt.Sscanf(resp.MTU, "%d", &mtu); err == nil {
			iface.MTU = mtu
		}
	}

	if resp.MSS != "" {
		mss := 0
		if _, err := fmt.Sscanf(resp.MSS, "%d", &mss); err == nil {
			iface.MSS = mss
		}
	}

	iface.Media = resp.Media
	iface.BlockPriv = resp.BlockPriv != nil
	iface.BlockBogons = resp.BlockBogons != nil

	return iface, nil
}

func (pf *Client) getInterfaces(ctx context.Context) (*Interfaces, error) {
	command := "$output = array();" +
		"foreach (config_get_path('interfaces', array()) as $ifname => $ifcfg) {" +
		"$ifcfg['_logical_name'] = $ifname;" +
		"array_push($output, $ifcfg);" +
		"};" +
		"print(json_encode($output));"

	var rawResp []map[string]interface{}
	if err := pf.executePHPCommand(ctx, command, &rawResp); err != nil {
		return nil, err
	}

	ifaces := make(Interfaces, 0, len(rawResp))
	for _, raw := range rawResp {
		logicalName, _ := raw["_logical_name"].(string)

		// Re-encode individual interface for proper parsing.
		resp := interfaceResponse{}
		resp.If, _ = raw["if"].(string)
		resp.Descr, _ = raw["descr"].(string)

		if _, ok := raw["enable"]; ok {
			empty := ""
			resp.Enable = &empty
		}

		resp.IPAddr, _ = raw["ipaddr"].(string)
		resp.Subnet, _ = raw["subnet"].(string)
		resp.Gateway, _ = raw["gateway"].(string)
		resp.IPAddrV6, _ = raw["ipaddrv6"].(string)
		resp.SubnetV6, _ = raw["subnetv6"].(string)
		resp.GatewayV6, _ = raw["gatewayv6"].(string)
		resp.SpoofMAC, _ = raw["spoofmac"].(string)
		resp.MTU, _ = raw["mtu"].(string)
		resp.MSS, _ = raw["mss"].(string)
		resp.Media, _ = raw["media"].(string)

		if _, ok := raw["blockpriv"]; ok {
			empty := ""
			resp.BlockPriv = &empty
		}

		if _, ok := raw["blockbogons"]; ok {
			empty := ""
			resp.BlockBogons = &empty
		}

		iface, err := parseInterfaceResponse(logicalName, resp)
		if err != nil {
			return nil, fmt.Errorf("%w interface '%s' response, %w", ErrUnableToParse, logicalName, err)
		}

		ifaces = append(ifaces, iface)
	}

	return &ifaces, nil
}

func (pf *Client) GetInterfaces(ctx context.Context) (*Interfaces, error) {
	defer pf.read(&pf.mutexes.Interface)()

	ifaces, err := pf.getInterfaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w interfaces, %w", ErrGetOperationFailed, err)
	}

	return ifaces, nil
}

func (pf *Client) GetInterface(ctx context.Context, logicalName string) (*Interface, error) {
	defer pf.read(&pf.mutexes.Interface)()

	ifaces, err := pf.getInterfaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w interfaces, %w", ErrGetOperationFailed, err)
	}

	iface, err := ifaces.GetByLogicalName(logicalName)
	if err != nil {
		return nil, fmt.Errorf("%w interface, %w", ErrGetOperationFailed, err)
	}

	return iface, nil
}

func (pf *Client) CreateInterface(ctx context.Context, req Interface) (*Interface, error) {
	defer pf.write(&pf.mutexes.Interface)()

	// Add a new interface assignment (assigns next available optN).
	// Then configure it with the requested settings.
	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"require_once('filter.inc');"+
			// Find next available optN.
			"$ifaces = config_get_path('interfaces', array());"+
			"$new_name = '';"+
			"for ($i = 1; $i <= 4096; $i++) {"+
			"if (!isset($ifaces['opt' . $i])) { $new_name = 'opt' . $i; break; }"+
			"}"+
			"if ($new_name === '') {"+
			"print(json_encode('No available interface slots'));"+
			"} else {"+
			// Create the interface assignment.
			"$ifcfg = array();"+
			"$ifcfg['if'] = '%s';"+
			"$ifcfg['descr'] = '%s';"+
			"%s"+ // enable
			"%s"+ // ipaddr
			"%s"+ // subnet
			"%s"+ // gateway
			"%s"+ // ipaddrv6
			"%s"+ // subnetv6
			"%s"+ // gatewayv6
			"%s"+ // spoofmac
			"%s"+ // mtu
			"%s"+ // mss
			"%s"+ // blockpriv
			"%s"+ // blockbogons
			"config_set_path('interfaces/' . $new_name, $ifcfg);"+
			"write_config('Terraform: created interface assignment ' . $new_name);"+
			"print(json_encode($new_name));"+
			"}",
		phpEscape(req.If),
		phpEscape(req.Description),
		phpBoolField("enable", req.Enabled),
		phpIPAddrField(req.IPv4Type, req.IPAddr),
		phpOptionalField("subnet", req.Subnet),
		phpOptionalField("gateway", req.Gateway),
		phpIPAddrV6Field(req.IPv6Type, req.IPAddrV6),
		phpOptionalField("subnetv6", req.SubnetV6),
		phpOptionalField("gatewayv6", req.GatewayV6),
		phpOptionalField("spoofmac", req.SpoofMAC),
		phpIntField("mtu", req.MTU),
		phpIntField("mss", req.MSS),
		phpBoolField("blockpriv", req.BlockPriv),
		phpBoolField("blockbogons", req.BlockBogons),
	)

	var result interface{}
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w interface, %w", ErrCreateOperationFailed, err)
	}

	logicalName, ok := result.(string)
	if !ok || logicalName == "" {
		return nil, fmt.Errorf("%w interface, unexpected result", ErrCreateOperationFailed)
	}

	// Check for error message.
	if logicalName == "No available interface slots" {
		return nil, fmt.Errorf("%w interface, %s", ErrCreateOperationFailed, logicalName)
	}

	ifaces, err := pf.getInterfaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w interfaces after creating, %w", ErrGetOperationFailed, err)
	}

	iface, err := ifaces.GetByLogicalName(logicalName)
	if err != nil {
		return nil, fmt.Errorf("%w interface after creating, %w", ErrGetOperationFailed, err)
	}

	return iface, nil
}

func (pf *Client) UpdateInterface(ctx context.Context, req Interface) (*Interface, error) {
	defer pf.write(&pf.mutexes.Interface)()

	// Verify the interface exists.
	ifaces, err := pf.getInterfaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w interfaces, %w", ErrGetOperationFailed, err)
	}

	if _, err := ifaces.GetByLogicalName(req.LogicalName); err != nil {
		return nil, fmt.Errorf("%w interface, %w", ErrGetOperationFailed, err)
	}

	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"require_once('filter.inc');"+
			"$ifcfg = array();"+
			"$ifcfg['if'] = '%s';"+
			"$ifcfg['descr'] = '%s';"+
			"%s"+ // enable
			"%s"+ // ipaddr
			"%s"+ // subnet
			"%s"+ // gateway
			"%s"+ // ipaddrv6
			"%s"+ // subnetv6
			"%s"+ // gatewayv6
			"%s"+ // spoofmac
			"%s"+ // mtu
			"%s"+ // mss
			"%s"+ // blockpriv
			"%s"+ // blockbogons
			"config_set_path('interfaces/%s', $ifcfg);"+
			"write_config('Terraform: updated interface %s');"+
			"print(json_encode(true));",
		phpEscape(req.If),
		phpEscape(req.Description),
		phpBoolField("enable", req.Enabled),
		phpIPAddrField(req.IPv4Type, req.IPAddr),
		phpOptionalField("subnet", req.Subnet),
		phpOptionalField("gateway", req.Gateway),
		phpIPAddrV6Field(req.IPv6Type, req.IPAddrV6),
		phpOptionalField("subnetv6", req.SubnetV6),
		phpOptionalField("gatewayv6", req.GatewayV6),
		phpOptionalField("spoofmac", req.SpoofMAC),
		phpIntField("mtu", req.MTU),
		phpIntField("mss", req.MSS),
		phpBoolField("blockpriv", req.BlockPriv),
		phpBoolField("blockbogons", req.BlockBogons),
		phpEscape(req.LogicalName),
		phpEscape(req.LogicalName),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w interface, %w", ErrUpdateOperationFailed, err)
	}

	ifaces, err = pf.getInterfaces(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w interfaces after updating, %w", ErrGetOperationFailed, err)
	}

	iface, err := ifaces.GetByLogicalName(req.LogicalName)
	if err != nil {
		return nil, fmt.Errorf("%w interface after updating, %w", ErrGetOperationFailed, err)
	}

	return iface, nil
}

func (pf *Client) DeleteInterface(ctx context.Context, logicalName string) error {
	defer pf.write(&pf.mutexes.Interface)()

	// Prevent deleting wan or lan.
	if logicalName == "wan" || logicalName == "lan" {
		return fmt.Errorf("%w interface, cannot delete '%s' — it is a system interface", ErrDeleteOperationFailed, logicalName)
	}

	ifaces, err := pf.getInterfaces(ctx)
	if err != nil {
		return fmt.Errorf("%w interfaces, %w", ErrGetOperationFailed, err)
	}

	if _, err := ifaces.GetByLogicalName(logicalName); err != nil {
		return fmt.Errorf("%w interface, %w", ErrGetOperationFailed, err)
	}

	// Check for dependencies (groups, bridges, GRE, GIF) and delete the interface.
	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"require_once('filter.inc');"+
			"$ifname = '%s';"+
			"$realif = get_real_interface($ifname);"+
			// Check if in use by interface group.
			"$groups = config_get_path('ifgroups/ifgroupentry', array());"+
			"if (!is_array($groups)) { $groups = array(); }"+
			"foreach ($groups as $g) {"+
			"if (!empty($g['members']) && in_array($ifname, explode(' ', $g['members']))) {"+
			"print(json_encode('Cannot delete interface: member of interface group ' . $g['ifname']));"+
			"exit;"+
			"}}"+
			// Delete interface config.
			"config_del_path('interfaces/' . $ifname);"+
			// Remove from DHCP.
			"config_del_path('dhcpd/' . $ifname);"+
			"config_del_path('dhcpdv6/' . $ifname);"+
			// Remove firewall rules referencing this interface.
			"$rules = config_get_path('filter/rule', array());"+
			"$newrules = array();"+
			"foreach ($rules as $rule) {"+
			"if (!isset($rule['interface']) || $rule['interface'] !== $ifname) { $newrules[] = $rule; }"+
			"}"+
			"config_set_path('filter/rule', $newrules);"+
			// Remove NAT rules referencing this interface.
			"$natrules = config_get_path('nat/rule', array());"+
			"$newnat = array();"+
			"foreach ($natrules as $rule) {"+
			"if (!isset($rule['interface']) || $rule['interface'] !== $ifname) { $newnat[] = $rule; }"+
			"}"+
			"config_set_path('nat/rule', $newnat);"+
			"write_config('Terraform: deleted interface assignment ' . $ifname);"+
			"print(json_encode(true));",
		phpEscape(logicalName),
	)

	var result interface{}
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w interface, %w", ErrDeleteOperationFailed, err)
	}

	if errMsg, ok := result.(string); ok {
		return fmt.Errorf("%w interface, %s", ErrDeleteOperationFailed, errMsg)
	}

	return nil
}

func (pf *Client) ApplyInterfaceChanges(ctx context.Context, logicalName string) error {
	pf.mutexes.InterfaceApply.Lock()
	defer pf.mutexes.InterfaceApply.Unlock()

	command := fmt.Sprintf(
		"require_once('interfaces.inc');"+
			"require_once('filter.inc');"+
			"require_once('rrd.inc');"+
			"require_once('shaper.inc');"+
			"interface_configure('%s');"+
			"filter_configure();"+
			"print(json_encode(0));",
		phpEscape(logicalName),
	)

	var result int
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply interface changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}

// Helper functions for building PHP config arrays.

func phpBoolField(field string, value bool) string {
	if value {
		return fmt.Sprintf("$ifcfg['%s'] = '';", field)
	}

	return ""
}

func phpOptionalField(field string, value string) string {
	if value != "" {
		return fmt.Sprintf("$ifcfg['%s'] = '%s';", field, phpEscape(value))
	}

	return ""
}

func phpIntField(field string, value int) string {
	if value > 0 {
		return fmt.Sprintf("$ifcfg['%s'] = '%d';", field, value)
	}

	return ""
}

func phpIPAddrField(ipv4Type string, ipAddr string) string {
	switch ipv4Type {
	case "staticv4":
		return fmt.Sprintf("$ifcfg['ipaddr'] = '%s';", phpEscape(ipAddr))
	case "dhcp":
		return "$ifcfg['ipaddr'] = 'dhcp';"
	case "pppoe":
		return "$ifcfg['ipaddr'] = 'pppoe';"
	case "pptp":
		return "$ifcfg['ipaddr'] = 'pptp';"
	case "l2tp":
		return "$ifcfg['ipaddr'] = 'l2tp';"
	default:
		return ""
	}
}

func phpIPAddrV6Field(ipv6Type string, ipAddrV6 string) string {
	switch ipv6Type {
	case "staticv6":
		return fmt.Sprintf("$ifcfg['ipaddrv6'] = '%s';", phpEscape(ipAddrV6))
	case "dhcp6":
		return "$ifcfg['ipaddrv6'] = 'dhcp6';"
	case "slaac":
		return "$ifcfg['ipaddrv6'] = 'slaac';"
	case "6to4":
		return "$ifcfg['ipaddrv6'] = '6to4';"
	case "6rd":
		return "$ifcfg['ipaddrv6'] = '6rd';"
	case "track6":
		return "$ifcfg['ipaddrv6'] = 'track6';"
	default:
		return ""
	}
}
