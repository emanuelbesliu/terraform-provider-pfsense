package pfsense

import (
	"context"
	"fmt"
)

type wakeOnLanResponse struct {
	Interface   string `json:"interface"`
	MAC         string `json:"mac"`
	Description string `json:"descr"`
	ControlID   int    `json:"controlID"` //nolint:tagliatelle
}

type WakeOnLanEntry struct {
	Interface   string
	MAC         string
	Description string
	controlID   int
}

func (w *WakeOnLanEntry) SetInterface(iface string) error {
	if err := ValidateInterface(iface); err != nil {
		return err
	}

	w.Interface = iface

	return nil
}

func (w *WakeOnLanEntry) SetMAC(mac string) error {
	if err := ValidateMACAddress(mac); err != nil {
		return err
	}

	w.MAC = mac

	return nil
}

func (w *WakeOnLanEntry) SetDescription(desc string) error {
	w.Description = desc

	return nil
}

type WakeOnLanEntries []WakeOnLanEntry

func (entries WakeOnLanEntries) GetByMAC(mac string) (*WakeOnLanEntry, error) {
	for _, e := range entries {
		if e.MAC == mac {
			return &e, nil
		}
	}

	return nil, fmt.Errorf("wake on lan entry %w with mac '%s'", ErrNotFound, mac)
}

func (entries WakeOnLanEntries) GetControlIDByMAC(mac string) (*int, error) {
	for _, e := range entries {
		if e.MAC == mac {
			return &e.controlID, nil
		}
	}

	return nil, fmt.Errorf("wake on lan entry %w with mac '%s'", ErrNotFound, mac)
}

func parseWakeOnLanResponse(resp wakeOnLanResponse) (WakeOnLanEntry, error) {
	var entry WakeOnLanEntry

	if resp.MAC == "" {
		return entry, fmt.Errorf("%w, wake on lan entry MAC is required", ErrClientValidation)
	}

	entry.Interface = resp.Interface
	entry.MAC = resp.MAC
	entry.Description = resp.Description
	entry.controlID = resp.ControlID

	return entry, nil
}

func (pf *Client) getWakeOnLanEntries(ctx context.Context) (*WakeOnLanEntries, error) {
	command := "$output = array();" +
		"$items = config_get_path('wol/wolentry', array());" +
		"if (!empty($items) && isset($items['interface'])) { $items = array($items); }" +
		"foreach ($items as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var resp []wakeOnLanResponse
	if err := pf.executePHPCommand(ctx, command, &resp); err != nil {
		return nil, err
	}

	entries := make(WakeOnLanEntries, 0, len(resp))
	for _, r := range resp {
		e, err := parseWakeOnLanResponse(r)
		if err != nil {
			return nil, fmt.Errorf("%w wake on lan response, %w", ErrUnableToParse, err)
		}

		entries = append(entries, e)
	}

	return &entries, nil
}

func (pf *Client) GetWakeOnLanEntries(ctx context.Context) (*WakeOnLanEntries, error) {
	defer pf.read(&pf.mutexes.WakeOnLan)()

	entries, err := pf.getWakeOnLanEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w wake on lan entries, %w", ErrGetOperationFailed, err)
	}

	return entries, nil
}

func (pf *Client) GetWakeOnLanEntry(ctx context.Context, mac string) (*WakeOnLanEntry, error) {
	defer pf.read(&pf.mutexes.WakeOnLan)()

	entries, err := pf.getWakeOnLanEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w wake on lan entries, %w", ErrGetOperationFailed, err)
	}

	entry, err := entries.GetByMAC(mac)
	if err != nil {
		return nil, fmt.Errorf("%w wake on lan entry, %w", ErrGetOperationFailed, err)
	}

	return entry, nil
}

func (pf *Client) CreateWakeOnLanEntry(ctx context.Context, req WakeOnLanEntry) (*WakeOnLanEntry, error) {
	defer pf.write(&pf.mutexes.WakeOnLan)()

	// Check for duplicate MAC.
	existingEntries, err := pf.getWakeOnLanEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w wake on lan entries for duplicate check, %w", ErrGetOperationFailed, err)
	}

	if _, err := existingEntries.GetByMAC(req.MAC); err == nil {
		return nil, fmt.Errorf("%w wake on lan entry, an entry with MAC '%s' already exists", ErrCreateOperationFailed, req.MAC)
	}

	command := fmt.Sprintf(
		"$item = array('interface' => '%s', 'mac' => '%s', 'descr' => '%s');"+
			"$existing = config_get_path('wol/wolentry', array());"+
			"if (!is_array($existing)) { $existing = array(); }"+
			"if (isset($existing['interface'])) { $existing = array($existing); }"+
			"$existing[] = $item;"+
			"config_set_path('wol/wolentry', $existing);"+
			"write_config('Terraform: created wake on lan entry');"+
			"print(json_encode(true));",
		phpEscape(req.Interface),
		phpEscape(req.MAC),
		phpEscape(req.Description),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w wake on lan entry, %w", ErrCreateOperationFailed, err)
	}

	entries, err := pf.getWakeOnLanEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w wake on lan entries after creating, %w", ErrGetOperationFailed, err)
	}

	entry, err := entries.GetByMAC(req.MAC)
	if err != nil {
		return nil, fmt.Errorf("%w wake on lan entry after creating, %w", ErrGetOperationFailed, err)
	}

	return entry, nil
}

func (pf *Client) UpdateWakeOnLanEntry(ctx context.Context, req WakeOnLanEntry) (*WakeOnLanEntry, error) {
	defer pf.write(&pf.mutexes.WakeOnLan)()

	entries, err := pf.getWakeOnLanEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w wake on lan entries, %w", ErrGetOperationFailed, err)
	}

	controlID, err := entries.GetControlIDByMAC(req.MAC)
	if err != nil {
		return nil, fmt.Errorf("%w wake on lan entry, %w", ErrGetOperationFailed, err)
	}

	command := fmt.Sprintf(
		"$item = array();"+
			"$item['interface'] = '%s';"+
			"$item['mac'] = '%s';"+
			"$item['descr'] = '%s';"+
			"config_set_path('wol/wolentry/%d', $item);"+
			"write_config('Terraform: updated wake on lan entry');"+
			"print(json_encode(true));",
		phpEscape(req.Interface),
		phpEscape(req.MAC),
		phpEscape(req.Description),
		*controlID,
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w wake on lan entry, %w", ErrUpdateOperationFailed, err)
	}

	entries, err = pf.getWakeOnLanEntries(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w wake on lan entries after updating, %w", ErrGetOperationFailed, err)
	}

	entry, err := entries.GetByMAC(req.MAC)
	if err != nil {
		return nil, fmt.Errorf("%w wake on lan entry after updating, %w", ErrGetOperationFailed, err)
	}

	return entry, nil
}

func (pf *Client) DeleteWakeOnLanEntry(ctx context.Context, mac string) error {
	defer pf.write(&pf.mutexes.WakeOnLan)()

	entries, err := pf.getWakeOnLanEntries(ctx)
	if err != nil {
		return fmt.Errorf("%w wake on lan entries, %w", ErrGetOperationFailed, err)
	}

	controlID, err := entries.GetControlIDByMAC(mac)
	if err != nil {
		return fmt.Errorf("%w wake on lan entry, %w", ErrGetOperationFailed, err)
	}

	command := fmt.Sprintf(
		"config_del_path('wol/wolentry/%d');"+
			"$items = config_get_path('wol/wolentry', array());"+
			"config_set_path('wol/wolentry', array_values($items));"+
			"write_config('Terraform: deleted wake on lan entry');"+
			"print(json_encode(true));",
		*controlID,
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w wake on lan entry, %w", ErrDeleteOperationFailed, err)
	}

	// Verify deletion.
	entries, err = pf.getWakeOnLanEntries(ctx)
	if err != nil {
		return fmt.Errorf("%w wake on lan entries after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := entries.GetByMAC(mac); err == nil {
		return fmt.Errorf("%w wake on lan entry, still exists", ErrDeleteOperationFailed)
	}

	return nil
}
