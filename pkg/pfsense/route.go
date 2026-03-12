package pfsense

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type routeResponse struct {
	Network     string  `json:"network"`
	Gateway     string  `json:"gateway"`
	Description string  `json:"descr"`
	Disabled    *string `json:"disabled"`
	ControlID   int     `json:"controlID"` //nolint:tagliatelle
}

type Route struct {
	Network     string
	Gateway     string
	Description string
	Disabled    bool
	controlID   int
}

func (r *Route) SetNetwork(network string) error {
	r.Network = network

	return nil
}

func (r *Route) SetGateway(gateway string) error {
	r.Gateway = gateway

	return nil
}

func (r *Route) SetDescription(description string) error {
	r.Description = description

	return nil
}

func (r *Route) SetDisabled(disabled bool) error {
	r.Disabled = disabled

	return nil
}

type Routes []Route

func (routes Routes) GetByNetworkAndGateway(network string, gateway string) (*Route, error) {
	for _, r := range routes {
		if r.Network == network && r.Gateway == gateway {
			return &r, nil
		}
	}

	return nil, fmt.Errorf("route %w with network '%s' and gateway '%s'", ErrNotFound, network, gateway)
}

func (routes Routes) GetControlIDByNetworkAndGateway(network string, gateway string) (*int, error) {
	for _, r := range routes {
		if r.Network == network && r.Gateway == gateway {
			return &r.controlID, nil
		}
	}

	return nil, fmt.Errorf("route %w with network '%s' and gateway '%s'", ErrNotFound, network, gateway)
}

func parseRouteResponse(resp routeResponse, index int) (Route, error) {
	var r Route

	if err := r.SetNetwork(resp.Network); err != nil {
		return r, err
	}

	if err := r.SetGateway(resp.Gateway); err != nil {
		return r, err
	}

	if err := r.SetDescription(resp.Description); err != nil {
		return r, err
	}

	r.Disabled = resp.Disabled != nil
	r.controlID = index

	return r, nil
}

func (pf *Client) getRoutes(ctx context.Context) (*Routes, error) {
	command := "$output = array();" +
		"if (isset($config['staticroutes']['route']) && is_array($config['staticroutes']['route'])) {" +
		"foreach ($config['staticroutes']['route'] as $k => $v) {" +
		"$v['controlID'] = $k; array_push($output, $v);" +
		"}};" +
		"print_r(json_encode($output));"

	var routeResp []routeResponse
	if err := pf.executePHPCommand(ctx, command, &routeResp); err != nil {
		return nil, err
	}

	routes := make(Routes, 0, len(routeResp))
	for index, resp := range routeResp {
		r, err := parseRouteResponse(resp, index)
		if err != nil {
			return nil, fmt.Errorf("%w route response, %w", ErrUnableToParse, err)
		}

		routes = append(routes, r)
	}

	return &routes, nil
}

func (pf *Client) GetRoutes(ctx context.Context) (*Routes, error) {
	defer pf.read(&pf.mutexes.Route)()

	routes, err := pf.getRoutes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w routes, %w", ErrGetOperationFailed, err)
	}

	return routes, nil
}

func (pf *Client) GetRoute(ctx context.Context, network string, gateway string) (*Route, error) {
	defer pf.read(&pf.mutexes.Route)()

	routes, err := pf.getRoutes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w routes, %w", ErrGetOperationFailed, err)
	}

	r, err := routes.GetByNetworkAndGateway(network, gateway)
	if err != nil {
		return nil, fmt.Errorf("%w route, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

func routeFormValues(routeReq Route) url.Values {
	// pfSense expects network address and subnet prefix as separate form fields.
	// Config stores them combined as CIDR (e.g. "10.10.0.0/24"), but
	// system_routes_edit.php requires "network" and "network_subnet" separately.
	network := routeReq.Network
	subnet := ""

	if parts := strings.SplitN(routeReq.Network, "/", 2); len(parts) == 2 {
		network = parts[0]
		subnet = parts[1]
	}

	values := url.Values{
		"network":        {network},
		"network_subnet": {subnet},
		"gateway":        {routeReq.Gateway},
		"descr":          {routeReq.Description},
		"save":           {"Save"},
	}

	if routeReq.Disabled {
		values.Set("disabled", "yes")
	}

	return values
}

func (pf *Client) createOrUpdateRoute(ctx context.Context, routeReq Route, controlID *int) error {
	relativeURL := url.URL{Path: "system_routes_edit.php"}
	values := routeFormValues(routeReq)

	if controlID != nil {
		q := relativeURL.Query()
		q.Set("id", strconv.Itoa(*controlID))
		relativeURL.RawQuery = q.Encode()
	}

	doc, err := pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return err
	}

	return scrapeHTMLValidationErrors(doc)
}

func (pf *Client) CreateRoute(ctx context.Context, routeReq Route) (*Route, error) {
	defer pf.write(&pf.mutexes.Route)()

	if err := pf.createOrUpdateRoute(ctx, routeReq, nil); err != nil {
		return nil, fmt.Errorf("%w route, %w", ErrCreateOperationFailed, err)
	}

	routes, err := pf.getRoutes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w routes after creating, %w", ErrGetOperationFailed, err)
	}

	r, err := routes.GetByNetworkAndGateway(routeReq.Network, routeReq.Gateway)
	if err != nil {
		return nil, fmt.Errorf("%w route after creating, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) UpdateRoute(ctx context.Context, routeReq Route, oldNetwork string, oldGateway string) (*Route, error) {
	defer pf.write(&pf.mutexes.Route)()

	routes, err := pf.getRoutes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w routes, %w", ErrGetOperationFailed, err)
	}

	controlID, err := routes.GetControlIDByNetworkAndGateway(oldNetwork, oldGateway)
	if err != nil {
		return nil, fmt.Errorf("%w route, %w", ErrGetOperationFailed, err)
	}

	if err := pf.createOrUpdateRoute(ctx, routeReq, controlID); err != nil {
		return nil, fmt.Errorf("%w route, %w", ErrUpdateOperationFailed, err)
	}

	routes, err = pf.getRoutes(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w routes after updating, %w", ErrGetOperationFailed, err)
	}

	r, err := routes.GetByNetworkAndGateway(routeReq.Network, routeReq.Gateway)
	if err != nil {
		return nil, fmt.Errorf("%w route after updating, %w", ErrGetOperationFailed, err)
	}

	return r, nil
}

func (pf *Client) DeleteRoute(ctx context.Context, network string, gateway string) error {
	defer pf.write(&pf.mutexes.Route)()

	routes, err := pf.getRoutes(ctx)
	if err != nil {
		return fmt.Errorf("%w routes, %w", ErrGetOperationFailed, err)
	}

	controlID, err := routes.GetControlIDByNetworkAndGateway(network, gateway)
	if err != nil {
		return fmt.Errorf("%w route, %w", ErrGetOperationFailed, err)
	}

	relativeURL := url.URL{Path: "system_routes.php"}
	values := url.Values{
		"act": {"del"},
		"id":  {strconv.Itoa(*controlID)},
	}

	_, err = pf.callHTML(ctx, http.MethodPost, relativeURL, &values)
	if err != nil {
		return fmt.Errorf("%w route, %w", ErrDeleteOperationFailed, err)
	}

	routes, err = pf.getRoutes(ctx)
	if err != nil {
		return fmt.Errorf("%w routes after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := routes.GetByNetworkAndGateway(network, gateway); err == nil {
		return fmt.Errorf("%w route, still exists", ErrDeleteOperationFailed)
	}

	return nil
}

func (pf *Client) ApplyRouteChanges(ctx context.Context) error {
	pf.mutexes.RouteApply.Lock()
	defer pf.mutexes.RouteApply.Unlock()

	command := "require_once(\"filter.inc\");" +
		"$retval = 0;" +
		"if (file_exists(\"{$g['tmp_path']}/.system_routes.apply\")) {" +
		"$toapplylist = unserialize(file_get_contents(\"{$g['tmp_path']}/.system_routes.apply\"));" +
		"foreach ($toapplylist as $toapply) mwexec(\"{$toapply}\");" +
		"@unlink(\"{$g['tmp_path']}/.system_routes.apply\");" +
		"}" +
		"$retval |= system_routing_configure();" +
		"$retval |= filter_configure();" +
		"setup_gateways_monitor();" +
		"if ($retval == 0) clear_subsystem_dirty('staticroutes');" +
		"print(json_encode($retval));"

	var result int
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply route changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}
