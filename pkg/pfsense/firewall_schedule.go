package pfsense

import (
	"context"
	"fmt"
	"strings"
)

type scheduleTimeRangeResponse struct {
	Position   string `json:"position"`
	Month      string `json:"month"`
	Day        string `json:"day"`
	Hour       string `json:"hour"`
	RangeDescr string `json:"rangedescr"`
}

type scheduleResponse struct {
	Name        string                      `json:"name"`
	Description string                      `json:"descr"`
	TimeRange   []scheduleTimeRangeResponse `json:"timerange"`
	SchedLabel  string                      `json:"schedlabel"`
	ControlID   int                         `json:"controlID"` //nolint:tagliatelle
}

// ScheduleTimeRange represents a single time range within a schedule. A time
// range is either a weekly recurring range (Position set) or a specific
// calendar range (Month and Day set), never both.
type ScheduleTimeRange struct {
	// Position is a comma separated list of weekday numbers where Monday is 1
	// and Sunday is 7 (e.g. "1,2,3,4,5" for weekdays). Used for weekly
	// recurring ranges.
	Position string
	// Month is a comma separated list of month numbers (1-12), parallel to Day.
	Month string
	// Day is a comma separated list of day-of-month numbers (1-31), parallel to
	// Month.
	Day string
	// StartTime is the start time in "H:MM" format (e.g. "0:00").
	StartTime string
	// StopTime is the stop time in "H:MM" format (e.g. "23:59").
	StopTime string
	// RangeDescription is an optional administrative description.
	RangeDescription string
}

type Schedule struct {
	Name        string
	Description string
	TimeRanges  []ScheduleTimeRange
	Label       string
	controlID   int
}

type Schedules []Schedule

func (s *Schedule) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("%w, name is required", ErrClientValidation)
	}

	switch strings.ToLower(name) {
	case "lan", "wan":
		return fmt.Errorf("%w, schedule may not be named LAN or WAN", ErrClientValidation)
	}

	s.Name = name

	return nil
}

func (s *Schedule) SetDescription(desc string) error {
	s.Description = desc

	return nil
}

func (s *Schedule) SetLabel(label string) error {
	s.Label = label

	return nil
}

func (s *Schedule) SetTimeRanges(ranges []ScheduleTimeRange) error {
	if len(ranges) == 0 {
		return fmt.Errorf("%w, schedule must have at least one time range", ErrClientValidation)
	}

	for i, tr := range ranges {
		hasPosition := tr.Position != ""
		hasMonthDay := tr.Month != "" || tr.Day != ""

		if hasPosition && hasMonthDay {
			return fmt.Errorf("%w, time range %d must set either position or month/day, not both", ErrClientValidation, i)
		}

		if !hasPosition && !hasMonthDay {
			return fmt.Errorf("%w, time range %d must set either position or month/day", ErrClientValidation, i)
		}

		if hasMonthDay && (tr.Month == "" || tr.Day == "") {
			return fmt.Errorf("%w, time range %d must set both month and day", ErrClientValidation, i)
		}

		if tr.StartTime == "" || tr.StopTime == "" {
			return fmt.Errorf("%w, time range %d must set start_time and stop_time", ErrClientValidation, i)
		}
	}

	s.TimeRanges = ranges

	return nil
}

func (schedules Schedules) GetByName(name string) (*Schedule, error) {
	for i := range schedules {
		if schedules[i].Name == name {
			return &schedules[i], nil
		}
	}

	return nil, fmt.Errorf("%w, schedule with name %q not found", ErrNotFound, name)
}

func (schedules Schedules) GetControlIDByName(name string) (int, error) {
	for _, s := range schedules {
		if s.Name == name {
			return s.controlID, nil
		}
	}

	return -1, fmt.Errorf("%w, schedule with name %q not found", ErrNotFound, name)
}

func parseScheduleResponse(resp scheduleResponse) (Schedule, error) {
	var s Schedule

	if err := s.SetName(resp.Name); err != nil {
		return s, err
	}

	s.Description = resp.Description
	s.Label = resp.SchedLabel
	s.controlID = resp.ControlID

	ranges := make([]ScheduleTimeRange, 0, len(resp.TimeRange))
	for _, tr := range resp.TimeRange {
		start, stop := parseScheduleHour(tr.Hour)
		ranges = append(ranges, ScheduleTimeRange{
			Position:         tr.Position,
			Month:            tr.Month,
			Day:              tr.Day,
			StartTime:        start,
			StopTime:         stop,
			RangeDescription: tr.RangeDescr,
		})
	}

	s.TimeRanges = ranges

	return s, nil
}

func parseScheduleHour(hour string) (string, string) {
	parts := strings.SplitN(hour, "-", 2)
	if len(parts) != 2 {
		return hour, ""
	}

	return parts[0], parts[1]
}

// ====================================================================
// PHP builders
// ====================================================================

func scheduleBuildTimeRange(tr ScheduleTimeRange) string {
	var b strings.Builder

	b.WriteString("$tr = array();")

	if tr.Position != "" {
		fmt.Fprintf(&b, "$tr['position'] = '%s';", phpEscape(tr.Position))
	} else {
		fmt.Fprintf(&b, "$tr['month'] = '%s';", phpEscape(tr.Month))
		fmt.Fprintf(&b, "$tr['day'] = '%s';", phpEscape(tr.Day))
	}

	fmt.Fprintf(&b, "$tr['hour'] = '%s-%s';", phpEscape(tr.StartTime), phpEscape(tr.StopTime))
	fmt.Fprintf(&b, "$tr['rangedescr'] = '%s';", phpEscape(tr.RangeDescription))
	b.WriteString("$schedule['timerange'][] = $tr;")

	return b.String()
}

func scheduleBuild(req Schedule) string {
	var b strings.Builder

	b.WriteString("$schedule = array();")
	fmt.Fprintf(&b, "$schedule['name'] = '%s';", phpEscape(req.Name))
	fmt.Fprintf(&b, "$schedule['descr'] = '%s';", phpEscape(req.Description))

	if req.Label != "" {
		fmt.Fprintf(&b, "$schedule['schedlabel'] = '%s';", phpEscape(req.Label))
	} else {
		b.WriteString("$schedule['schedlabel'] = uniqid();")
	}

	b.WriteString("$schedule['timerange'] = array();")
	for _, tr := range req.TimeRanges {
		b.WriteString(scheduleBuildTimeRange(tr))
	}

	return b.String()
}

// ====================================================================
// CRUD operations
// ====================================================================

func (pf *Client) getSchedules(ctx context.Context) (*Schedules, error) {
	command := "$output = array();" +
		"$schedules = config_get_path('schedules/schedule', array());" +
		"foreach ($schedules as $k => $v) {" +
		"$v['controlID'] = $k;" +
		"array_push($output, $v);" +
		"};" +
		"print(json_encode($output));"

	var schedResp []scheduleResponse
	if err := pf.executePHPCommand(ctx, command, &schedResp); err != nil {
		return nil, err
	}

	schedules := make(Schedules, 0, len(schedResp))
	for _, resp := range schedResp {
		s, err := parseScheduleResponse(resp)
		if err != nil {
			return nil, fmt.Errorf("%w schedule response, %w", ErrUnableToParse, err)
		}

		schedules = append(schedules, s)
	}

	return &schedules, nil
}

func (pf *Client) GetSchedules(ctx context.Context) (*Schedules, error) {
	defer pf.read(&pf.mutexes.FirewallSchedule)()

	schedules, err := pf.getSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w schedules, %w", ErrGetOperationFailed, err)
	}

	return schedules, nil
}

func (pf *Client) GetSchedule(ctx context.Context, name string) (*Schedule, error) {
	defer pf.read(&pf.mutexes.FirewallSchedule)()

	schedules, err := pf.getSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w schedules, %w", ErrGetOperationFailed, err)
	}

	s, err := schedules.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w schedule, %w", ErrGetOperationFailed, err)
	}

	return s, nil
}

func (pf *Client) CreateSchedule(ctx context.Context, req Schedule) (*Schedule, error) {
	defer pf.write(&pf.mutexes.FirewallSchedule)()

	schedules, err := pf.getSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w schedules, %w", ErrCreateOperationFailed, err)
	}

	if _, err := schedules.GetByName(req.Name); err == nil {
		return nil, fmt.Errorf("%w schedule, name %q already exists", ErrCreateOperationFailed, req.Name)
	}

	command := "require_once('filter.inc');" +
		scheduleBuild(req) +
		"$schedules = config_get_path('schedules/schedule', array());" +
		"$schedules[] = $schedule;" +
		"config_set_path('schedules/schedule', $schedules);" +
		"write_config('Terraform: create firewall schedule');" +
		"filter_configure();" +
		"print(json_encode(array('status' => 'ok')));"

	var result map[string]string
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w schedule, %w", ErrCreateOperationFailed, err)
	}

	updated, err := pf.getSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w reading back schedule, %w", ErrCreateOperationFailed, err)
	}

	s, err := updated.GetByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w reading back schedule, %w", ErrCreateOperationFailed, err)
	}

	return s, nil
}

func (pf *Client) UpdateSchedule(ctx context.Context, name string, req Schedule) (*Schedule, error) {
	defer pf.write(&pf.mutexes.FirewallSchedule)()

	schedules, err := pf.getSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w schedules, %w", ErrUpdateOperationFailed, err)
	}

	existing, err := schedules.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w schedule, %w", ErrUpdateOperationFailed, err)
	}

	controlID := existing.controlID

	if req.Label == "" {
		req.Label = existing.Label
	}

	command := "require_once('filter.inc');" +
		scheduleBuild(req) +
		fmt.Sprintf("config_set_path('schedules/schedule/%d', $schedule);", controlID) +
		"write_config('Terraform: update firewall schedule');" +
		"filter_configure();" +
		"print(json_encode(array('status' => 'ok')));"

	var result map[string]string
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w schedule, %w", ErrUpdateOperationFailed, err)
	}

	updated, err := pf.getSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w reading back schedule, %w", ErrUpdateOperationFailed, err)
	}

	s, err := updated.GetByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("%w reading back schedule, %w", ErrUpdateOperationFailed, err)
	}

	return s, nil
}

func (pf *Client) DeleteSchedule(ctx context.Context, name string) error {
	defer pf.write(&pf.mutexes.FirewallSchedule)()

	schedules, err := pf.getSchedules(ctx)
	if err != nil {
		return fmt.Errorf("%w schedules, %w", ErrDeleteOperationFailed, err)
	}

	controlID, err := schedules.GetControlIDByName(name)
	if err != nil {
		return fmt.Errorf("%w schedule, %w", ErrDeleteOperationFailed, err)
	}

	command := "require_once('filter.inc');" +
		fmt.Sprintf("config_del_path('schedules/schedule/%d');", controlID) +
		"write_config('Terraform: delete firewall schedule');" +
		"filter_configure();" +
		"print(json_encode(array('status' => 'ok')));"

	var result map[string]string
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return fmt.Errorf("%w schedule, %w", ErrDeleteOperationFailed, err)
	}

	return nil
}
