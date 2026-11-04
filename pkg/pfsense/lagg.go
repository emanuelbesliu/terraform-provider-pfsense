package pfsense

import (
	"context"
	"fmt"
	"strings"
)

// Lagg models a pfSense LAGG. pfSense stores LAGGs as an unkeyed list under
// laggs/lagg; each entry carries its own stable name in the 'laggif' field
// (lagg0, lagg1, ...), exactly like the pfSense GUI (interfaces_lagg_edit.php)
// which derives the name from interface_lagg_configure(). We mirror that
// contract: field 'proto', 'members' CSV, 'laggif' name, 'descr'.
type Lagg struct {
	Name        string // laggif, e.g. "lagg0" (assigned by pfSense on create)
	Members     []string
	Protocol    string // stored in pfSense as 'proto'
	Description string
}

var laggProtocols = []string{"none", "failover", "lacp", "loadbalance", "protect", "prefixlacp"}

// phpLiteral escapes a string for direct inclusion in single-quoted PHP source.
func phpLiteral(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, `\`, `\\`), "'", `\'`)
}

func (l *Lagg) SetMembers(members []string) error {
	if len(members) == 0 {
		return fmt.Errorf("%w, at least one member interface is required", ErrClientValidation)
	}
	for _, m := range members {
		if m == "" {
			return fmt.Errorf("%w, member interface names may not be empty", ErrClientValidation)
		}
	}
	l.Members = members
	return nil
}

func (l *Lagg) SetProtocol(protocol string) error {
	for _, p := range laggProtocols {
		if p == protocol {
			l.Protocol = protocol
			return nil
		}
	}
	return fmt.Errorf("%w, unknown lagg protocol '%s' (options: %s)", ErrClientValidation, protocol, strings.Join(laggProtocols, ", "))
}

func (l *Lagg) SetDescription(description string) error {
	l.Description = description
	return nil
}

func (l *Lagg) membersCSV() string {
	return strings.Join(l.Members, ",")
}

type laggResponse struct {
	Name    string `json:"laggif"`
	Members string `json:"members"` // pfSense stores a CSV string
	Proto   string `json:"proto"`
	Descr   string `json:"descr"`
}

func parseLaggResponse(resp laggResponse, fallbackName string) Lagg {
	members := []string{}
	if resp.Members != "" {
		members = strings.Split(resp.Members, ",")
	}
	name := resp.Name
	if name == "" {
		// legacy entry without laggif (pre-existing from manual config edits):
		// pfSense resolves these by position on apply; surface a stable fallback.
		name = fallbackName
	}
	return Lagg{Name: name, Members: members, Protocol: resp.Proto, Description: resp.Descr}
}

func (pf *Client) getLaggs(ctx context.Context) ([]Lagg, error) {
	command := "$laggs = config_get_path('laggs/lagg', array());" +
		"$output = array(); $i = 0;" +
		"foreach ($laggs as $v) {" +
		"$v['__pos'] = 'lagg' . $i; $i++;" +
		"array_push($output, $v);" +
		"}" +
		"print(json_encode($output));"

	var resp []struct {
		laggResponse
		Pos string `json:"__pos"`
	}
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}
	laggs := make([]Lagg, 0, len(resp))
	for _, r := range resp {
		laggs = append(laggs, parseLaggResponse(r.laggResponse, r.Pos))
	}
	return laggs, nil
}

func findLagg(laggs []Lagg, name string) (*Lagg, error) {
	for _, l := range laggs {
		if l.Name == name {
			out := l
			return &out, nil
		}
	}
	return nil, fmt.Errorf("lagg %w with name '%s'", ErrNotFound, name)
}

func (pf *Client) GetLaggs(ctx context.Context) ([]Lagg, error) {
	defer pf.read(&pf.mutexes.Lagg)()
	laggs, err := pf.getLaggs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w laggs, %w", ErrGetOperationFailed, err)
	}
	return laggs, nil
}

func (pf *Client) GetLagg(ctx context.Context, name string) (*Lagg, error) {
	defer pf.read(&pf.mutexes.Lagg)()
	laggs, err := pf.getLaggs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w laggs, %w", ErrGetOperationFailed, err)
	}
	return findLagg(laggs, name)
}

// CreateLagg appends the LAGG the same way interfaces_lagg_edit.php does:
// next free laggN name computed over existing laggif values, entry stored with
// laggif/proto/members/descr. When apply is true, interface_lagg_configure()
// brings the aggregation up immediately (like the GUI save button does).
func (pf *Client) CreateLagg(ctx context.Context, lagg Lagg, apply bool) (*Lagg, error) {
	defer pf.write(&pf.mutexes.Lagg)()

	command := fmt.Sprintf(
		"require_once('interfaces.inc'); global $config;"+
			"$items = config_get_path('laggs/lagg', array());"+
			"$taken = array(); foreach ($items as $v) { if (isset($v['laggif'])) { $taken[$v['laggif']] = true; } }"+
			"$i = 0; while (isset($taken['lagg'.$i])) { $i++; }"+
			"$n = 'lagg'.$i;"+
			"$entry = array('laggif' => $n, 'members' => '%s', 'proto' => '%s', 'descr' => '%s');"+
			"$items[] = $entry;"+
			"config_set_path('laggs/lagg', $items);"+
			"write_config('Terraform: created lagg '.$n);"+
			"if (%t) { interface_lagg_configure($entry); }"+
			"echo json_encode($entry);"+
			"exit(0);",
		phpLiteral(lagg.membersCSV()), phpLiteral(lagg.Protocol), phpLiteral(lagg.Description), apply)

	var resp laggResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, fmt.Errorf("%w lagg, %w", ErrCreateOperationFailed, err)
	}
	// server-side truth: re-read the list (internal getter: we hold the write lock)
	laggsAfter, err := pf.getLaggs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w lagg after create, %w", ErrCreateOperationFailed, err)
	}
	return findLagg(laggsAfter, resp.Name)
}

func (pf *Client) UpdateLagg(ctx context.Context, lagg Lagg, apply bool) (*Lagg, error) {
	defer pf.write(&pf.mutexes.Lagg)()
	if lagg.Name == "" {
		return nil, fmt.Errorf("%w, lagg name is required for update", ErrClientValidation)
	}

	command := fmt.Sprintf(
		"require_once('interfaces.inc'); global $config;"+
			"$items = config_get_path('laggs/lagg', array());"+
			"$found = -1; foreach ($items as $k => $v) { if (($v['laggif'] ?? '') === '%s') { $found = $k; break; } }"+
			"if ($found === -1) { throw new Exception('lagg not found'); }"+
			"$entry = array('laggif' => '%s', 'members' => '%s', 'proto' => '%s', 'descr' => '%s');"+
			"$items[$found] = $entry;"+
			"config_set_path('laggs/lagg', $items);"+
			"write_config('Terraform: updated lagg %s');"+
			"if (%t) { interface_lagg_configure($entry); }"+
			"echo json_encode(0);"+
			"exit(0);",
		phpLiteral(lagg.Name), phpLiteral(lagg.Name), phpLiteral(lagg.membersCSV()),
		phpLiteral(lagg.Protocol), phpLiteral(lagg.Description), phpLiteral(lagg.Name), apply)

	var result int
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w lagg, %w", ErrUpdateOperationFailed, err)
	}
	laggsAfter, err := pf.getLaggs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w lagg after update, %w", ErrUpdateOperationFailed, err)
	}
	return findLagg(laggsAfter, lagg.Name)
}

// DeleteLagg removes the config entry (matching GUI delete which also
// destroys the kernel interface when it exists; we replicate both, guarded).
func (pf *Client) DeleteLagg(ctx context.Context, name string) error {
	defer pf.write(&pf.mutexes.Lagg)()

	command := fmt.Sprintf(
		"require_once('interfaces.inc'); require_once('util.inc'); global $config;"+
			"$items = config_get_path('laggs/lagg', array());"+
			"$found = -1; foreach ($items as $k => $v) { if (($v['laggif'] ?? '') === '%s') { $found = $k; break; } }"+
			"if ($found !== -1) {"+
			" if (function_exists('pfSense_interface_exists') && pfSense_interface_exists($items[$found]['laggif'])) { pfSense_interface_destroy($items[$found]['laggif']); }"+
			" unset($items[$found]);"+
			" $items = array_values($items);"+
			" config_set_path('laggs/lagg', $items);"+
			" write_config('Terraform: deleted lagg %s');"+
			"}"+
			"echo json_encode(0);"+
			"exit(0);",
		phpLiteral(name), phpLiteral(name))

	var result int
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w lagg, %w", ErrDeleteOperationFailed, err)
	}
	laggsAfter, _ := pf.getLaggs(ctx) // still holding the write mutex: use the internal getter
	if laggsAfter != nil {
		if _, stillThere := findLagg(laggsAfter, name); stillThere == nil {
			return fmt.Errorf("failed to delete lagg %s, still exists after delete", name)
		}
	}
	return nil
}
