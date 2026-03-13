package pfsense

import (
	"context"
	"fmt"
)

type tunableResponse struct {
	Tunable     string `json:"tunable"`
	Value       string `json:"value"`
	Description string `json:"descr"`
	ControlID   int    `json:"controlID"` //nolint:tagliatelle
}

type SystemTunable struct {
	Tunable     string
	Value       string
	Description string
	controlID   int
}

func (t *SystemTunable) SetTunable(tunable string) error {
	if tunable == "" {
		return fmt.Errorf("%w, tunable name is required", ErrClientValidation)
	}

	t.Tunable = tunable

	return nil
}

func (t *SystemTunable) SetValue(value string) error {
	if value == "" {
		return fmt.Errorf("%w, tunable value is required", ErrClientValidation)
	}

	t.Value = value

	return nil
}

func (t *SystemTunable) SetDescription(description string) error {
	t.Description = description

	return nil
}

type SystemTunables []SystemTunable

func (tunables SystemTunables) GetByName(name string) (*SystemTunable, error) {
	for _, t := range tunables {
		if t.Tunable == name {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("system tunable %w with name '%s'", ErrNotFound, name)
}

func (tunables SystemTunables) GetControlIDByName(name string) (*int, error) {
	for _, t := range tunables {
		if t.Tunable == name {
			return &t.controlID, nil
		}
	}

	return nil, fmt.Errorf("system tunable %w with name '%s'", ErrNotFound, name)
}

func parseTunableResponse(resp tunableResponse) (SystemTunable, error) {
	var tunable SystemTunable

	if err := tunable.SetTunable(resp.Tunable); err != nil {
		return tunable, err
	}

	if err := tunable.SetValue(resp.Value); err != nil {
		return tunable, err
	}

	if err := tunable.SetDescription(resp.Description); err != nil {
		return tunable, err
	}

	tunable.controlID = resp.ControlID

	return tunable, nil
}

func (pf *Client) getTunables(ctx context.Context) (*SystemTunables, error) {
	command := "$output = array();" +
		"$items = config_get_path('sysctl/item', array());" +
		"foreach ($items as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var tunableResp []tunableResponse
	if err := pf.executePHPCommand(ctx, command, &tunableResp); err != nil {
		return nil, err
	}

	tunables := make(SystemTunables, 0, len(tunableResp))
	for _, resp := range tunableResp {
		t, err := parseTunableResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w system tunable response, %w", ErrUnableToParse, err)
		}

		tunables = append(tunables, t)
	}

	return &tunables, nil
}

func (pf *Client) GetSystemTunables(ctx context.Context) (*SystemTunables, error) {
	defer pf.read(&pf.mutexes.SystemTunable)()

	tunables, err := pf.getTunables(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w system tunables, %w", ErrGetOperationFailed, err)
	}

	return tunables, nil
}

func (pf *Client) GetSystemTunable(ctx context.Context, name string) (*SystemTunable, error) {
	defer pf.read(&pf.mutexes.SystemTunable)()

	tunables, err := pf.getTunables(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w system tunables, %w", ErrGetOperationFailed, err)
	}

	t, err := tunables.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w system tunable, %w", ErrGetOperationFailed, err)
	}

	return t, nil
}

func (pf *Client) CreateSystemTunable(ctx context.Context, req SystemTunable) (*SystemTunable, error) {
	defer pf.write(&pf.mutexes.SystemTunable)()

	// Check for duplicate tunable name.
	existingTunables, err := pf.getTunables(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w system tunables for duplicate check, %w", ErrGetOperationFailed, err)
	}

	if _, err := existingTunables.GetByName(req.Tunable); err == nil {
		return nil, fmt.Errorf("%w system tunable, a tunable with name '%s' already exists", ErrCreateOperationFailed, req.Tunable)
	}

	command := fmt.Sprintf(
		"$tunable = array();"+
			"$tunable['tunable'] = '%s';"+
			"$tunable['value'] = '%s';"+
			"$tunable['descr'] = '%s';"+
			"config_set_path('sysctl/item/', $tunable);"+
			"write_config('Terraform: created system tunable %s');"+
			"print(json_encode(true));",
		phpEscape(req.Tunable),
		phpEscape(req.Value),
		phpEscape(req.Description),
		phpEscape(req.Tunable),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w system tunable, %w", ErrCreateOperationFailed, err)
	}

	tunables, err := pf.getTunables(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w system tunables after creating, %w", ErrGetOperationFailed, err)
	}

	t, err := tunables.GetByName(req.Tunable)
	if err != nil {
		return nil, fmt.Errorf("%w system tunable after creating, %w", ErrGetOperationFailed, err)
	}

	return t, nil
}

func (pf *Client) UpdateSystemTunable(ctx context.Context, req SystemTunable) (*SystemTunable, error) {
	defer pf.write(&pf.mutexes.SystemTunable)()

	tunables, err := pf.getTunables(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w system tunables, %w", ErrGetOperationFailed, err)
	}

	controlID, err := tunables.GetControlIDByName(req.Tunable)
	if err != nil {
		return nil, fmt.Errorf("%w system tunable, %w", ErrGetOperationFailed, err)
	}

	command := fmt.Sprintf(
		"$tunable = array();"+
			"$tunable['tunable'] = '%s';"+
			"$tunable['value'] = '%s';"+
			"$tunable['descr'] = '%s';"+
			"config_set_path('sysctl/item/%d', $tunable);"+
			"write_config('Terraform: updated system tunable %s');"+
			"print(json_encode(true));",
		phpEscape(req.Tunable),
		phpEscape(req.Value),
		phpEscape(req.Description),
		*controlID,
		phpEscape(req.Tunable),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w system tunable, %w", ErrUpdateOperationFailed, err)
	}

	tunables, err = pf.getTunables(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w system tunables after updating, %w", ErrGetOperationFailed, err)
	}

	t, err := tunables.GetByName(req.Tunable)
	if err != nil {
		return nil, fmt.Errorf("%w system tunable after updating, %w", ErrGetOperationFailed, err)
	}

	return t, nil
}

func (pf *Client) DeleteSystemTunable(ctx context.Context, name string) error {
	defer pf.write(&pf.mutexes.SystemTunable)()

	tunables, err := pf.getTunables(ctx)
	if err != nil {
		return fmt.Errorf("%w system tunables, %w", ErrGetOperationFailed, err)
	}

	controlID, err := tunables.GetControlIDByName(name)
	if err != nil {
		return fmt.Errorf("%w system tunable, %w", ErrGetOperationFailed, err)
	}

	command := fmt.Sprintf(
		"config_del_path('sysctl/item/%d');"+
			"write_config('Terraform: deleted system tunable %s');"+
			"print(json_encode(true));",
		*controlID,
		phpEscape(name),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w system tunable, %w", ErrDeleteOperationFailed, err)
	}

	// Verify deletion.
	tunables, err = pf.getTunables(ctx)
	if err != nil {
		return fmt.Errorf("%w system tunables after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := tunables.GetByName(name); err == nil {
		return fmt.Errorf("%w system tunable, still exists", ErrDeleteOperationFailed)
	}

	return nil
}

func (pf *Client) ApplySystemTunableChanges(ctx context.Context) error {
	pf.mutexes.SystemTunableApply.Lock()
	defer pf.mutexes.SystemTunableApply.Unlock()

	command := "require_once(\"system.inc\");" +
		"system_setup_sysctl();" +
		"clear_subsystem_dirty('sysctl');" +
		"print(json_encode(true));"

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply system tunable changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}
