package pfsense

import (
	"strings"
	"testing"
)

func TestNATNPt_SetInterface(t *testing.T) {
	tests := []struct {
		name    string
		iface   string
		wantErr bool
	}{
		{"wan", "wan", false},
		{"lan", "lan", false},
		{"opt1", "opt1", false},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r NATNPt

			err := r.SetInterface(tt.iface)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetInterface(%q) error = %v, wantErr %v", tt.iface, err, tt.wantErr)
			}

			if err == nil && r.Interface != tt.iface {
				t.Errorf("SetInterface(%q) got %q", tt.iface, r.Interface)
			}
		})
	}
}

func TestNATNPt_SetPrefixes(t *testing.T) {
	var r NATNPt

	if err := r.SetSourcePrefix("fd00:1::/64"); err != nil {
		t.Fatalf("SetSourcePrefix returned error: %v", err)
	}

	if r.SourcePrefix != "fd00:1::/64" {
		t.Errorf("SetSourcePrefix got %q", r.SourcePrefix)
	}

	if err := r.SetDestinationPrefix("2001:db8::/64"); err != nil {
		t.Fatalf("SetDestinationPrefix returned error: %v", err)
	}

	if r.DestinationPrefix != "2001:db8::/64" {
		t.Errorf("SetDestinationPrefix got %q", r.DestinationPrefix)
	}
}

func TestNATNPts_GetByDescription(t *testing.T) {
	rules := NATNPts{
		{Description: "first", ControlID: 0},
		{Description: "second", ControlID: 1},
	}

	r, err := rules.GetByDescription("second")
	if err != nil {
		t.Fatalf("GetByDescription returned error: %v", err)
	}

	if r.ControlID != 1 {
		t.Errorf("GetByDescription got ControlID %d, want 1", r.ControlID)
	}

	if _, err := rules.GetByDescription("missing"); err == nil {
		t.Error("GetByDescription(missing) expected error, got nil")
	}
}

func TestNATNPts_GetControlIDByDescription(t *testing.T) {
	rules := NATNPts{
		{Description: "first", ControlID: 0},
		{Description: "second", ControlID: 1},
	}

	id, err := rules.GetControlIDByDescription("first")
	if err != nil {
		t.Fatalf("GetControlIDByDescription returned error: %v", err)
	}

	if id != 0 {
		t.Errorf("GetControlIDByDescription got %d, want 0", id)
	}

	if _, err := rules.GetControlIDByDescription("missing"); err == nil {
		t.Error("GetControlIDByDescription(missing) expected error, got nil")
	}
}

func TestNATNPt_buildRule(t *testing.T) {
	tests := []struct {
		name     string
		rule     NATNPt
		contains []string
	}{
		{
			name: "basic",
			rule: NATNPt{
				Interface:         "wan",
				SourcePrefix:      "fd00:1::/64",
				DestinationPrefix: "2001:db8::/64",
				Description:       "test",
			},
			contains: []string{
				"$rule['interface'] = 'wan';",
				"$rule['descr'] = 'test';",
				"$rule['source']['address'] = 'fd00:1::/64';",
				"$rule['destination']['address'] = '2001:db8::/64';",
			},
		},
		{
			name: "negated and disabled",
			rule: NATNPt{
				Interface:         "lan",
				SourcePrefix:      "fd00:2::/64",
				SourceNot:         true,
				DestinationPrefix: "2001:db8:1::/64",
				DestinationNot:    true,
				Disabled:          true,
				Description:       "neg",
			},
			contains: []string{
				"$rule['source']['not'] = '';",
				"$rule['destination']['not'] = '';",
				"$rule['disabled'] = '';",
			},
		},
		{
			name: "track6 destination network",
			rule: NATNPt{
				Interface:         "wan",
				SourcePrefix:      "fd00:3::/64",
				DestinationPrefix: "lan",
				Description:       "track6",
			},
			contains: []string{
				"$rule['destination']['network'] = 'lan';",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := natNPtBuildRule(tt.rule)
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("natNPtBuildRule() = %q, want it to contain %q", got, want)
				}
			}
		})
	}
}
