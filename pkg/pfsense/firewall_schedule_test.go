package pfsense

import (
	"strings"
	"testing"
)

func TestSchedule_SetName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "WorkHours", false},
		{"empty", "", true},
		{"lan reserved", "LAN", true},
		{"wan reserved", "WAN", true},
		{"lan lowercase reserved", "lan", true},
		{"wan mixedcase reserved", "Wan", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Schedule

			err := s.SetName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}

			if err == nil && s.Name != tt.input {
				t.Errorf("SetName(%q) got %q", tt.input, s.Name)
			}
		})
	}
}

func TestSchedule_SetTimeRanges(t *testing.T) {
	tests := []struct {
		name    string
		ranges  []ScheduleTimeRange
		wantErr bool
	}{
		{
			name:    "empty list",
			ranges:  []ScheduleTimeRange{},
			wantErr: true,
		},
		{
			name: "valid position",
			ranges: []ScheduleTimeRange{
				{Position: "1,2,3,4,5", StartTime: "9:00", StopTime: "17:00"},
			},
			wantErr: false,
		},
		{
			name: "valid month and day",
			ranges: []ScheduleTimeRange{
				{Month: "12", Day: "25", StartTime: "0:00", StopTime: "23:59"},
			},
			wantErr: false,
		},
		{
			name: "position and month/day both set",
			ranges: []ScheduleTimeRange{
				{Position: "1", Month: "12", Day: "25", StartTime: "0:00", StopTime: "23:59"},
			},
			wantErr: true,
		},
		{
			name: "neither position nor month/day",
			ranges: []ScheduleTimeRange{
				{StartTime: "0:00", StopTime: "23:59"},
			},
			wantErr: true,
		},
		{
			name: "month without day",
			ranges: []ScheduleTimeRange{
				{Month: "12", StartTime: "0:00", StopTime: "23:59"},
			},
			wantErr: true,
		},
		{
			name: "day without month",
			ranges: []ScheduleTimeRange{
				{Day: "25", StartTime: "0:00", StopTime: "23:59"},
			},
			wantErr: true,
		},
		{
			name: "missing start time",
			ranges: []ScheduleTimeRange{
				{Position: "1", StopTime: "23:59"},
			},
			wantErr: true,
		},
		{
			name: "missing stop time",
			ranges: []ScheduleTimeRange{
				{Position: "1", StartTime: "0:00"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Schedule

			err := s.SetTimeRanges(tt.ranges)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetTimeRanges() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseScheduleHour(t *testing.T) {
	tests := []struct {
		name      string
		hour      string
		wantStart string
		wantStop  string
	}{
		{"normal", "9:00-17:00", "9:00", "17:00"},
		{"full day", "0:00-23:59", "0:00", "23:59"},
		{"no separator", "9:00", "9:00", ""},
		{"empty", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, stop := parseScheduleHour(tt.hour)
			if start != tt.wantStart || stop != tt.wantStop {
				t.Errorf("parseScheduleHour(%q) = (%q, %q), want (%q, %q)", tt.hour, start, stop, tt.wantStart, tt.wantStop)
			}
		})
	}
}

func TestSchedules_GetByName(t *testing.T) {
	schedules := Schedules{
		{Name: "first", controlID: 0},
		{Name: "second", controlID: 1},
	}

	s, err := schedules.GetByName("second")
	if err != nil {
		t.Fatalf("GetByName returned error: %v", err)
	}

	if s.controlID != 1 {
		t.Errorf("GetByName got controlID %d, want 1", s.controlID)
	}

	if _, err := schedules.GetByName("missing"); err == nil {
		t.Error("GetByName(missing) expected error, got nil")
	}
}

func TestSchedules_GetControlIDByName(t *testing.T) {
	schedules := Schedules{
		{Name: "first", controlID: 0},
		{Name: "second", controlID: 1},
	}

	id, err := schedules.GetControlIDByName("first")
	if err != nil {
		t.Fatalf("GetControlIDByName returned error: %v", err)
	}

	if id != 0 {
		t.Errorf("GetControlIDByName got %d, want 0", id)
	}

	if _, err := schedules.GetControlIDByName("missing"); err == nil {
		t.Error("GetControlIDByName(missing) expected error, got nil")
	}
}

func TestScheduleBuild(t *testing.T) {
	tests := []struct {
		name     string
		schedule Schedule
		contains []string
	}{
		{
			name: "position with label",
			schedule: Schedule{
				Name:        "WorkHours",
				Description: "Business hours",
				Label:       "abc123",
				TimeRanges: []ScheduleTimeRange{
					{Position: "1,2,3,4,5", StartTime: "9:00", StopTime: "17:00", RangeDescription: "weekdays"},
				},
			},
			contains: []string{
				"$schedule['name'] = 'WorkHours';",
				"$schedule['descr'] = 'Business hours';",
				"$schedule['schedlabel'] = 'abc123';",
				"$tr['position'] = '1,2,3,4,5';",
				"$tr['hour'] = '9:00-17:00';",
				"$tr['rangedescr'] = 'weekdays';",
				"$schedule['timerange'][] = $tr;",
			},
		},
		{
			name: "month and day without label",
			schedule: Schedule{
				Name:        "Christmas",
				Description: "Holiday",
				TimeRanges: []ScheduleTimeRange{
					{Month: "12", Day: "25", StartTime: "0:00", StopTime: "23:59"},
				},
			},
			contains: []string{
				"$schedule['schedlabel'] = uniqid();",
				"$tr['month'] = '12';",
				"$tr['day'] = '25';",
				"$tr['hour'] = '0:00-23:59';",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := scheduleBuild(tt.schedule)
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("scheduleBuild() = %q, want it to contain %q", got, want)
				}
			}
		})
	}
}
