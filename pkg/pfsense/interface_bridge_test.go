package pfsense

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBridge_SetMembers(t *testing.T) {
	tests := []struct {
		name    string
		members []string
		wantErr bool
	}{
		{"valid", []string{"lan", "opt1"}, false},
		{"single", []string{"lan"}, false},
		{"empty", []string{}, true},
		{"nil", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Bridge

			err := b.SetMembers(tt.members)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetMembers(%v) error = %v, wantErr %v", tt.members, err, tt.wantErr)
			}
		})
	}
}

func TestBridge_SetProtocol(t *testing.T) {
	tests := []struct {
		name    string
		proto   string
		wantErr bool
	}{
		{"rstp", "rstp", false},
		{"stp", "stp", false},
		{"empty allowed", "", false},
		{"invalid", "mstp", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Bridge

			err := b.SetProtocol(tt.proto)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetProtocol(%q) error = %v, wantErr %v", tt.proto, err, tt.wantErr)
			}

			if err == nil && tt.proto != "" && b.Protocol != tt.proto {
				t.Errorf("SetProtocol(%q) got %q", tt.proto, b.Protocol)
			}
		})
	}
}

func TestBridge_validateSubset(t *testing.T) {
	b := Bridge{Members: []string{"lan", "opt1"}}

	if err := b.validateSubset("stp", []string{"lan"}); err != nil {
		t.Errorf("validateSubset valid subset returned error: %v", err)
	}

	if err := b.validateSubset("stp", []string{"lan", "opt1"}); err != nil {
		t.Errorf("validateSubset full subset returned error: %v", err)
	}

	if err := b.validateSubset("stp", []string{"opt2"}); err == nil {
		t.Error("validateSubset(opt2) expected error, got nil")
	}
}

func TestBridge_validateDisjoint(t *testing.T) {
	b := Bridge{Members: []string{"lan", "opt1"}}

	if err := b.validateDisjoint("span", []string{"opt2"}); err != nil {
		t.Errorf("validateDisjoint non-member returned error: %v", err)
	}

	if err := b.validateDisjoint("span", nil); err != nil {
		t.Errorf("validateDisjoint empty returned error: %v", err)
	}

	if err := b.validateDisjoint("span", []string{"lan"}); err == nil {
		t.Error("validateDisjoint(lan) expected error (span cannot be a member), got nil")
	}
}

func TestParseBridgeMemberList(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  []string
	}{
		{"two", "lan,opt1", []string{"lan", "opt1"}},
		{"one", "lan", []string{"lan"}},
		{"empty", "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseBridgeMemberList(tt.value)
			if len(got) != len(tt.want) {
				t.Fatalf("parseBridgeMemberList(%q) = %v, want %v", tt.value, got, tt.want)
			}

			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("parseBridgeMemberList(%q)[%d] = %q, want %q", tt.value, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseBridgeInt(t *testing.T) {
	parsed, err := parseBridgeInt("priority", "32768")
	if err != nil {
		t.Fatalf("parseBridgeInt returned error: %v", err)
	}

	if parsed == nil || *parsed != 32768 {
		t.Errorf("parseBridgeInt(32768) = %v, want 32768", parsed)
	}

	empty, err := parseBridgeInt("priority", "")
	if err != nil {
		t.Fatalf("parseBridgeInt(empty) returned error: %v", err)
	}

	if empty != nil {
		t.Errorf("parseBridgeInt(empty) = %v, want nil", empty)
	}

	if _, err := parseBridgeInt("priority", "notanumber"); err == nil {
		t.Error("parseBridgeInt(notanumber) expected error, got nil")
	}
}

func TestParseBridgeResponse_Booleans(t *testing.T) {
	resp := bridgeResponse{
		BridgeIf:     "bridge0",
		Members:      "lan,opt1",
		EnableSTP:    json.RawMessage(`true`),
		IP6LinkLocal: nil,
		Priority:     "32768",
	}

	bridge, err := parseBridgeResponse(resp)
	if err != nil {
		t.Fatalf("parseBridgeResponse returned error: %v", err)
	}

	if !bridge.EnableSTP {
		t.Error("expected EnableSTP true")
	}

	if bridge.IP6LinkLocal {
		t.Error("expected IP6LinkLocal false")
	}

	if bridge.Priority == nil || *bridge.Priority != 32768 {
		t.Errorf("expected Priority 32768, got %v", bridge.Priority)
	}

	if len(bridge.Members) != 2 {
		t.Errorf("expected 2 members, got %v", bridge.Members)
	}
}

func TestBridges_GetByBridgeIf(t *testing.T) {
	bridges := Bridges{
		{BridgeIf: "bridge0", controlID: 0},
		{BridgeIf: "bridge1", controlID: 1},
	}

	b, err := bridges.GetByBridgeIf("bridge1")
	if err != nil {
		t.Fatalf("GetByBridgeIf returned error: %v", err)
	}

	if b.controlID != 1 {
		t.Errorf("GetByBridgeIf got controlID %d, want 1", b.controlID)
	}

	id, err := bridges.GetControlIDByBridgeIf("bridge0")
	if err != nil {
		t.Fatalf("GetControlIDByBridgeIf returned error: %v", err)
	}

	if id != 0 {
		t.Errorf("GetControlIDByBridgeIf got %d, want 0", id)
	}

	if _, err := bridges.GetByBridgeIf("bridge9"); err == nil {
		t.Error("GetByBridgeIf(bridge9) expected error, got nil")
	}
}

func TestBridgeBuild(t *testing.T) {
	tests := []struct {
		name        string
		bridge      Bridge
		contains    []string
		notContains []string
	}{
		{
			name: "full stp",
			bridge: Bridge{
				Members:       []string{"lan", "opt1"},
				Description:   "test bridge",
				EnableSTP:     true,
				IP6LinkLocal:  true,
				Protocol:      "rstp",
				Priority:      intPtr(32768),
				STPInterfaces: []string{"lan"},
			},
			contains: []string{
				"$bridge['members'] = 'lan,opt1';",
				"$bridge['descr'] = 'test bridge';",
				"$bridge['enablestp'] = true;",
				"$bridge['ip6linklocal'] = true;",
				"$bridge['proto'] = 'rstp';",
				"$bridge['priority'] = '32768';",
				"$bridge['stp'] = 'lan';",
			},
		},
		{
			name: "minimal no stp",
			bridge: Bridge{
				Members: []string{"lan"},
			},
			contains: []string{
				"$bridge['members'] = 'lan';",
			},
			notContains: []string{
				"$bridge['enablestp']",
				"$bridge['ip6linklocal']",
				"$bridge['proto']",
				"$bridge['priority']",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bridgeBuild(tt.bridge)
			for _, want := range tt.contains {
				if !strings.Contains(got, want) {
					t.Errorf("bridgeBuild() = %q, want it to contain %q", got, want)
				}
			}

			for _, notWant := range tt.notContains {
				if strings.Contains(got, notWant) {
					t.Errorf("bridgeBuild() = %q, want it to NOT contain %q", got, notWant)
				}
			}
		})
	}
}
