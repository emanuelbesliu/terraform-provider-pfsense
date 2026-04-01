package pfsense

import (
	"context"
	"fmt"
)

type cronJobResponse struct {
	Minute    string `json:"minute"`
	Hour      string `json:"hour"`
	MDay      string `json:"mday"`
	Month     string `json:"month"`
	WDay      string `json:"wday"`
	Who       string `json:"who"`
	Command   string `json:"command"`
	ControlID int    `json:"controlID"` //nolint:tagliatelle
}

type CronJob struct {
	Minute    string
	Hour      string
	MDay      string
	Month     string
	WDay      string
	Who       string
	Command   string
	controlID int
}

func (c *CronJob) SetMinute(minute string) error {
	if minute == "" {
		return fmt.Errorf("%w, cron job minute is required", ErrClientValidation)
	}

	c.Minute = minute

	return nil
}

func (c *CronJob) SetHour(hour string) error {
	if hour == "" {
		return fmt.Errorf("%w, cron job hour is required", ErrClientValidation)
	}

	c.Hour = hour

	return nil
}

func (c *CronJob) SetMDay(mday string) error {
	if mday == "" {
		return fmt.Errorf("%w, cron job day of month is required", ErrClientValidation)
	}

	c.MDay = mday

	return nil
}

func (c *CronJob) SetMonth(month string) error {
	if month == "" {
		return fmt.Errorf("%w, cron job month is required", ErrClientValidation)
	}

	c.Month = month

	return nil
}

func (c *CronJob) SetWDay(wday string) error {
	if wday == "" {
		return fmt.Errorf("%w, cron job day of week is required", ErrClientValidation)
	}

	c.WDay = wday

	return nil
}

func (c *CronJob) SetWho(who string) error {
	if who == "" {
		return fmt.Errorf("%w, cron job user is required", ErrClientValidation)
	}

	c.Who = who

	return nil
}

func (c *CronJob) SetCommand(command string) error {
	if command == "" {
		return fmt.Errorf("%w, cron job command is required", ErrClientValidation)
	}

	c.Command = command

	return nil
}

type CronJobs []CronJob

func (jobs CronJobs) GetByCommand(command string) (*CronJob, error) {
	for _, j := range jobs {
		if j.Command == command {
			return &j, nil
		}
	}

	return nil, fmt.Errorf("cron job %w with command '%s'", ErrNotFound, command)
}

func (jobs CronJobs) GetControlIDByCommand(command string) (*int, error) {
	for _, j := range jobs {
		if j.Command == command {
			return &j.controlID, nil
		}
	}

	return nil, fmt.Errorf("cron job %w with command '%s'", ErrNotFound, command)
}

func parseCronJobResponse(resp cronJobResponse) (CronJob, error) {
	var job CronJob

	// Command is always required — skip entries with no command (invalid/corrupt).
	if resp.Command == "" {
		return job, fmt.Errorf("%w, cron job command is required", ErrClientValidation)
	}

	// Populate fields directly — tolerate empty schedule fields from existing config.
	// Setters are used for create/update validation; parsing must be lenient.
	job.Minute = resp.Minute
	job.Hour = resp.Hour
	job.MDay = resp.MDay
	job.Month = resp.Month
	job.WDay = resp.WDay
	job.Who = resp.Who
	job.Command = resp.Command
	job.controlID = resp.ControlID

	return job, nil
}

func (pf *Client) getCronJobs(ctx context.Context) (*CronJobs, error) {
	command := "$output = array();" +
		"$items = config_get_path('cron/item', array());" +
		"foreach ($items as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var jobResp []cronJobResponse
	if err := pf.executePHPCommand(ctx, command, &jobResp); err != nil {
		return nil, err
	}

	jobs := make(CronJobs, 0, len(jobResp))
	for _, resp := range jobResp {
		j, err := parseCronJobResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w cron job response, %w", ErrUnableToParse, err)
		}

		jobs = append(jobs, j)
	}

	return &jobs, nil
}

func (pf *Client) GetCronJobs(ctx context.Context) (*CronJobs, error) {
	defer pf.read(&pf.mutexes.CronJob)()

	jobs, err := pf.getCronJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w cron jobs, %w", ErrGetOperationFailed, err)
	}

	return jobs, nil
}

func (pf *Client) GetCronJob(ctx context.Context, command string) (*CronJob, error) {
	defer pf.read(&pf.mutexes.CronJob)()

	jobs, err := pf.getCronJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w cron jobs, %w", ErrGetOperationFailed, err)
	}

	j, err := jobs.GetByCommand(command)
	if err != nil {
		return nil, fmt.Errorf("%w cron job, %w", ErrGetOperationFailed, err)
	}

	return j, nil
}

func (pf *Client) CreateCronJob(ctx context.Context, req CronJob) (*CronJob, error) {
	defer pf.write(&pf.mutexes.CronJob)()

	// Check for duplicate command.
	existingJobs, err := pf.getCronJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w cron jobs for duplicate check, %w", ErrGetOperationFailed, err)
	}

	if _, err := existingJobs.GetByCommand(req.Command); err == nil {
		return nil, fmt.Errorf("%w cron job, a cron job with command '%s' already exists", ErrCreateOperationFailed, req.Command)
	}

	command := fmt.Sprintf(
		"$item = array();"+
			"$item['minute'] = '%s';"+
			"$item['hour'] = '%s';"+
			"$item['mday'] = '%s';"+
			"$item['month'] = '%s';"+
			"$item['wday'] = '%s';"+
			"$item['who'] = '%s';"+
			"$item['command'] = '%s';"+
			"config_set_path('cron/item/', $item);"+
			"write_config('Terraform: created cron job');"+
			"print(json_encode(true));",
		phpEscape(req.Minute),
		phpEscape(req.Hour),
		phpEscape(req.MDay),
		phpEscape(req.Month),
		phpEscape(req.WDay),
		phpEscape(req.Who),
		phpEscape(req.Command),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w cron job, %w", ErrCreateOperationFailed, err)
	}

	// Apply cron changes.
	if err := pf.applyCronChanges(ctx); err != nil {
		return nil, fmt.Errorf("%w cron job, %w", ErrCreateOperationFailed, err)
	}

	jobs, err := pf.getCronJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w cron jobs after creating, %w", ErrGetOperationFailed, err)
	}

	j, err := jobs.GetByCommand(req.Command)
	if err != nil {
		return nil, fmt.Errorf("%w cron job after creating, %w", ErrGetOperationFailed, err)
	}

	return j, nil
}

func (pf *Client) UpdateCronJob(ctx context.Context, req CronJob) (*CronJob, error) {
	defer pf.write(&pf.mutexes.CronJob)()

	jobs, err := pf.getCronJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w cron jobs, %w", ErrGetOperationFailed, err)
	}

	controlID, err := jobs.GetControlIDByCommand(req.Command)
	if err != nil {
		return nil, fmt.Errorf("%w cron job, %w", ErrGetOperationFailed, err)
	}

	command := fmt.Sprintf(
		"$item = array();"+
			"$item['minute'] = '%s';"+
			"$item['hour'] = '%s';"+
			"$item['mday'] = '%s';"+
			"$item['month'] = '%s';"+
			"$item['wday'] = '%s';"+
			"$item['who'] = '%s';"+
			"$item['command'] = '%s';"+
			"config_set_path('cron/item/%d', $item);"+
			"write_config('Terraform: updated cron job');"+
			"print(json_encode(true));",
		phpEscape(req.Minute),
		phpEscape(req.Hour),
		phpEscape(req.MDay),
		phpEscape(req.Month),
		phpEscape(req.WDay),
		phpEscape(req.Who),
		phpEscape(req.Command),
		*controlID,
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w cron job, %w", ErrUpdateOperationFailed, err)
	}

	// Apply cron changes.
	if err := pf.applyCronChanges(ctx); err != nil {
		return nil, fmt.Errorf("%w cron job, %w", ErrUpdateOperationFailed, err)
	}

	jobs, err = pf.getCronJobs(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w cron jobs after updating, %w", ErrGetOperationFailed, err)
	}

	j, err := jobs.GetByCommand(req.Command)
	if err != nil {
		return nil, fmt.Errorf("%w cron job after updating, %w", ErrGetOperationFailed, err)
	}

	return j, nil
}

func (pf *Client) DeleteCronJob(ctx context.Context, cmdStr string) error {
	defer pf.write(&pf.mutexes.CronJob)()

	jobs, err := pf.getCronJobs(ctx)
	if err != nil {
		return fmt.Errorf("%w cron jobs, %w", ErrGetOperationFailed, err)
	}

	controlID, err := jobs.GetControlIDByCommand(cmdStr)
	if err != nil {
		return fmt.Errorf("%w cron job, %w", ErrGetOperationFailed, err)
	}

	command := fmt.Sprintf(
		"config_del_path('cron/item/%d');"+
			"$items = config_get_path('cron/item', array());"+
			"config_set_path('cron/item', array_values($items));"+
			"write_config('Terraform: deleted cron job');"+
			"print(json_encode(true));",
		*controlID,
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w cron job, %w", ErrDeleteOperationFailed, err)
	}

	// Apply cron changes.
	if err := pf.applyCronChanges(ctx); err != nil {
		return fmt.Errorf("%w cron job, %w", ErrDeleteOperationFailed, err)
	}

	// Verify deletion.
	jobs, err = pf.getCronJobs(ctx)
	if err != nil {
		return fmt.Errorf("%w cron jobs after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := jobs.GetByCommand(cmdStr); err == nil {
		return fmt.Errorf("%w cron job, still exists", ErrDeleteOperationFailed)
	}

	return nil
}

func (pf *Client) applyCronChanges(ctx context.Context) error {
	command := "require_once('cron.inc');" +
		"configure_cron();" +
		"print(json_encode(true));"

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w, failed to apply cron changes, %w", ErrApplyOperationFailed, err)
	}

	return nil
}
