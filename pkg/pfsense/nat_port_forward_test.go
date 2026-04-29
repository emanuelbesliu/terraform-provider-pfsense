package pfsense

import (
	"testing"
)

func TestNATPortForward_SetInterface(t *testing.T) {
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
			var r NATPortForward

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

func TestNATPortForward_SetIPProtocol(t *testing.T) {
	tests := []struct {
		name    string
		ipp     string
		wantErr bool
	}{
		{"inet", "inet", false},
		{"inet6", "inet6", false},
		{"inet46", "inet46", false},
		{"invalid", "invalid", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r NATPortForward

			err := r.SetIPProtocol(tt.ipp)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetIPProtocol(%q) error = %v, wantErr %v", tt.ipp, err, tt.wantErr)
			}

			if err == nil && r.IPProtocol != tt.ipp {
				t.Errorf("SetIPProtocol(%q) got %q", tt.ipp, r.IPProtocol)
			}
		})
	}
}

func TestNATPortForward_SetProtocol(t *testing.T) {
	tests := []struct {
		name    string
		proto   string
		wantErr bool
	}{
		{"tcp", "tcp", false},
		{"udp", "udp", false},
		{"tcp_udp", "tcp/udp", false},
		{"icmp", "icmp", false},
		{"any", "any", true},
		{"invalid", "invalid", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r NATPortForward

			err := r.SetProtocol(tt.proto)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetProtocol(%q) error = %v, wantErr %v", tt.proto, err, tt.wantErr)
			}

			if err == nil && r.Protocol != tt.proto {
				t.Errorf("SetProtocol(%q) got %q", tt.proto, r.Protocol)
			}
		})
	}
}

func TestNATPortForward_SetTarget(t *testing.T) {
	tests := []struct {
		name    string
		target  string
		wantErr bool
	}{
		{"ip", "10.0.1.100", false},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r NATPortForward

			err := r.SetTarget(tt.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetTarget(%q) error = %v, wantErr %v", tt.target, err, tt.wantErr)
			}

			if err == nil && r.Target != tt.target {
				t.Errorf("SetTarget(%q) got %q", tt.target, r.Target)
			}
		})
	}
}

func TestNATPortForward_SetNATReflection(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		wantErr bool
	}{
		{"empty_default", "", false},
		{"enable", "enable", false},
		{"disable", "disable", false},
		{"purenat", "purenat", false},
		{"invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r NATPortForward

			err := r.SetNATReflection(tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetNATReflection(%q) error = %v, wantErr %v", tt.mode, err, tt.wantErr)
			}

			if err == nil && r.NATReflection != tt.mode {
				t.Errorf("SetNATReflection(%q) got %q", tt.mode, r.NATReflection)
			}
		})
	}
}

func TestNATPortForward_SetAssociatedRuleID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr bool
	}{
		{"pass", "pass", false},
		{"empty", "", false},
		{"invalid", "block", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r NATPortForward

			err := r.SetAssociatedRuleID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetAssociatedRuleID(%q) error = %v, wantErr %v", tt.id, err, tt.wantErr)
			}

			if err == nil && r.AssociatedRuleID != tt.id {
				t.Errorf("SetAssociatedRuleID(%q) got %q", tt.id, r.AssociatedRuleID)
			}
		})
	}
}

func TestParseNATPortForwardResponse(t *testing.T) {
	resp := natPortForwardResponse{
		Interface:  "wan",
		IPProtocol: "inet",
		Protocol:   "tcp",
		Source:     map[string]any{"any": ""},
		Destination: map[string]any{
			"network": "wanip",
			"port":    "8080",
		},
		Target:           "10.0.1.50",
		LocalPort:        "80",
		Description:      "Web Server Forward",
		Disabled:         "",
		NoRDR:            "",
		NATReflection:    "enable",
		AssociatedRuleID: "pass",
		ControlID:        3,
	}

	r, err := parseNATPortForwardResponse(resp)
	if err != nil {
		t.Fatalf("parseNATPortForwardResponse() error = %v", err)
	}

	if r.Interface != "wan" {
		t.Errorf("Interface = %q, want %q", r.Interface, "wan")
	}

	if r.IPProtocol != "inet" {
		t.Errorf("IPProtocol = %q, want %q", r.IPProtocol, "inet")
	}

	if r.Protocol != "tcp" {
		t.Errorf("Protocol = %q, want %q", r.Protocol, "tcp")
	}

	if r.SourceAddress != "any" {
		t.Errorf("SourceAddress = %q, want %q", r.SourceAddress, "any")
	}

	if r.DestAddress != "wanip" {
		t.Errorf("DestAddress = %q, want %q", r.DestAddress, "wanip")
	}

	if r.DestPort != "8080" {
		t.Errorf("DestPort = %q, want %q", r.DestPort, "8080")
	}

	if r.Target != "10.0.1.50" {
		t.Errorf("Target = %q, want %q", r.Target, "10.0.1.50")
	}

	if r.LocalPort != "80" {
		t.Errorf("LocalPort = %q, want %q", r.LocalPort, "80")
	}

	if r.Description != "Web Server Forward" {
		t.Errorf("Description = %q, want %q", r.Description, "Web Server Forward")
	}

	if r.Disabled {
		t.Error("Disabled = true, want false")
	}

	if r.NATReflection != "enable" {
		t.Errorf("NATReflection = %q, want %q", r.NATReflection, "enable")
	}

	if r.AssociatedRuleID != "pass" {
		t.Errorf("AssociatedRuleID = %q, want %q", r.AssociatedRuleID, "pass")
	}

	if r.controlID != 3 {
		t.Errorf("controlID = %d, want 3", r.controlID)
	}
}

func TestParseNATPortForwardResponse_Disabled(t *testing.T) {
	resp := natPortForwardResponse{
		Interface:  "wan",
		IPProtocol: "inet",
		Protocol:   "tcp",
		Source:     map[string]any{"any": ""},
		Destination: map[string]any{
			"address": "203.0.113.1",
			"port":    "443",
		},
		Target:    "10.0.1.50",
		LocalPort: "443",
		Disabled:  "yes",
		NoRDR:     "yes",
	}

	r, err := parseNATPortForwardResponse(resp)
	if err != nil {
		t.Fatalf("parseNATPortForwardResponse() error = %v", err)
	}

	if !r.Disabled {
		t.Error("Disabled = false, want true")
	}

	if !r.NoRDR {
		t.Error("NoRDR = false, want true")
	}
}

func TestParseNATPortForwardResponse_SourceWithNot(t *testing.T) {
	resp := natPortForwardResponse{
		Interface:  "wan",
		IPProtocol: "inet",
		Protocol:   "tcp",
		Source: map[string]any{
			"address": "192.168.1.0/24",
			"not":     "",
		},
		Destination: map[string]any{
			"any": "",
		},
		Target: "10.0.1.1",
	}

	r, err := parseNATPortForwardResponse(resp)
	if err != nil {
		t.Fatalf("parseNATPortForwardResponse() error = %v", err)
	}

	if !r.SourceNot {
		t.Error("SourceNot = false, want true")
	}

	if r.SourceAddress != "192.168.1.0/24" {
		t.Errorf("SourceAddress = %q, want %q", r.SourceAddress, "192.168.1.0/24")
	}
}

func TestNATPortForwards_GetByDescription(t *testing.T) {
	rules := NATPortForwards{
		{Description: "rule1", Interface: "wan", Target: "10.0.1.1"},
		{Description: "rule2", Interface: "lan", Target: "10.0.1.2"},
	}

	r, err := rules.GetByDescription("rule2")
	if err != nil {
		t.Fatalf("GetByDescription() error = %v", err)
	}

	if r.Target != "10.0.1.2" {
		t.Errorf("Target = %q, want %q", r.Target, "10.0.1.2")
	}

	_, err = rules.GetByDescription("nonexistent")
	if err == nil {
		t.Fatal("GetByDescription() expected error for nonexistent description")
	}
}

func TestNATPortForwards_GetControlIDByDescription(t *testing.T) {
	rules := NATPortForwards{
		{Description: "rule1", controlID: 0},
		{Description: "rule2", controlID: 7},
	}

	cid, err := rules.GetControlIDByDescription("rule2")
	if err != nil {
		t.Fatalf("GetControlIDByDescription() error = %v", err)
	}

	if *cid != 7 {
		t.Errorf("controlID = %d, want 7", *cid)
	}

	_, err = rules.GetControlIDByDescription("nonexistent")
	if err == nil {
		t.Fatal("GetControlIDByDescription() expected error for nonexistent description")
	}
}

func TestNATPortForwardPHPSource(t *testing.T) {
	tests := []struct {
		name     string
		req      NATPortForward
		contains []string
	}{
		{
			"any",
			NATPortForward{SourceAddress: "any"},
			[]string{"$rule['source']['any']"},
		},
		{
			"cidr",
			NATPortForward{SourceAddress: "10.0.0.0/24"},
			[]string{"$rule['source']['address'] = '10.0.0.0/24'"},
		},
		{
			"single_ip",
			NATPortForward{SourceAddress: "10.0.0.1"},
			[]string{"$rule['source']['address'] = '10.0.0.1'"},
		},
		{
			"special_wanip",
			NATPortForward{SourceAddress: "wanip"},
			[]string{"$rule['source']['network'] = 'wanip'"},
		},
		{
			"with_not",
			NATPortForward{SourceAddress: "any", SourceNot: true},
			[]string{"$rule['source']['not']", "$rule['source']['any']"},
		},
		{
			"with_port",
			NATPortForward{SourceAddress: "any", SourcePort: "8080"},
			[]string{"$rule['source']['port'] = '8080'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := natPortForwardPHPSource(tt.req)
			for _, substr := range tt.contains {
				if !containsStr(result, substr) {
					t.Errorf("natPortForwardPHPSource() = %q, missing %q", result, substr)
				}
			}
		})
	}
}

func TestNATPortForwardPHPDest(t *testing.T) {
	tests := []struct {
		name     string
		req      NATPortForward
		contains []string
	}{
		{
			"any",
			NATPortForward{DestAddress: "any"},
			[]string{"$rule['destination']['any']"},
		},
		{
			"special_wanip",
			NATPortForward{DestAddress: "wanip"},
			[]string{"$rule['destination']['network'] = 'wanip'"},
		},
		{
			"with_port",
			NATPortForward{DestAddress: "any", DestPort: "443"},
			[]string{"$rule['destination']['port'] = '443'"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := natPortForwardPHPDest(tt.req)
			for _, substr := range tt.contains {
				if !containsStr(result, substr) {
					t.Errorf("natPortForwardPHPDest() = %q, missing %q", result, substr)
				}
			}
		})
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
