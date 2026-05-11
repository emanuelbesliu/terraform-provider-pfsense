package pfsense

import (
	"testing"
)

func TestNATOutboundRule_SetInterface(t *testing.T) {
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
			var r NATOutboundRule

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

func TestNATOutboundRule_SetProtocol(t *testing.T) {
	tests := []struct {
		name    string
		proto   string
		wantErr bool
	}{
		{"empty_any", "", false},
		{"tcp", "tcp", false},
		{"udp", "udp", false},
		{"tcp_udp", "tcp/udp", false},
		{"icmp", "icmp", false},
		{"esp", "esp", false},
		{"ah", "ah", false},
		{"gre", "gre", false},
		{"ipv6", "ipv6", false},
		{"igmp", "igmp", false},
		{"pim", "pim", false},
		{"ospf", "ospf", false},
		{"invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r NATOutboundRule

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

func TestNATOutboundRule_SetPoolOpts(t *testing.T) {
	tests := []struct {
		name    string
		opts    string
		wantErr bool
	}{
		{"empty", "", false},
		{"round_robin", "round-robin", false},
		{"round_robin_sticky", "round-robin sticky-address", false},
		{"random", "random", false},
		{"random_sticky", "random sticky-address", false},
		{"source_hash", "source-hash", false},
		{"bitmask", "bitmask", false},
		{"invalid", "invalid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var r NATOutboundRule

			err := r.SetPoolOpts(tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetPoolOpts(%q) error = %v, wantErr %v", tt.opts, err, tt.wantErr)
			}

			if err == nil && r.PoolOpts != tt.opts {
				t.Errorf("SetPoolOpts(%q) got %q", tt.opts, r.PoolOpts)
			}
		})
	}
}

func TestNATOutboundRule_SetSourceAddress(t *testing.T) {
	var r NATOutboundRule
	if err := r.SetSourceAddress("192.168.1.0/24"); err != nil {
		t.Errorf("SetSourceAddress() unexpected error: %v", err)
	}
	if r.SourceAddress != "192.168.1.0/24" {
		t.Errorf("got %q, want %q", r.SourceAddress, "192.168.1.0/24")
	}
}

func TestNATOutboundRule_SetSourcePort(t *testing.T) {
	var r NATOutboundRule
	if err := r.SetSourcePort("443"); err != nil {
		t.Errorf("SetSourcePort() unexpected error: %v", err)
	}
	if r.SourcePort != "443" {
		t.Errorf("got %q, want %q", r.SourcePort, "443")
	}
}

func TestNATOutboundRule_SetDestAddress(t *testing.T) {
	var r NATOutboundRule
	if err := r.SetDestAddress("10.0.0.0/8"); err != nil {
		t.Errorf("SetDestAddress() unexpected error: %v", err)
	}
	if r.DestAddress != "10.0.0.0/8" {
		t.Errorf("got %q, want %q", r.DestAddress, "10.0.0.0/8")
	}
}

func TestNATOutboundRule_SetDestPort(t *testing.T) {
	var r NATOutboundRule
	if err := r.SetDestPort("80"); err != nil {
		t.Errorf("SetDestPort() unexpected error: %v", err)
	}
	if r.DestPort != "80" {
		t.Errorf("got %q, want %q", r.DestPort, "80")
	}
}

func TestNATOutboundRule_SetTarget(t *testing.T) {
	var r NATOutboundRule
	if err := r.SetTarget("(self)"); err != nil {
		t.Errorf("SetTarget() unexpected error: %v", err)
	}
	if r.Target != "(self)" {
		t.Errorf("got %q, want %q", r.Target, "(self)")
	}
}

func TestNATOutboundRule_SetTargetIP(t *testing.T) {
	var r NATOutboundRule
	if err := r.SetTargetIP("203.0.113.1"); err != nil {
		t.Errorf("SetTargetIP() unexpected error: %v", err)
	}
	if r.TargetIP != "203.0.113.1" {
		t.Errorf("got %q, want %q", r.TargetIP, "203.0.113.1")
	}
}

func TestNATOutboundRule_SetTargetIPSubnet(t *testing.T) {
	var r NATOutboundRule
	if err := r.SetTargetIPSubnet("24"); err != nil {
		t.Errorf("SetTargetIPSubnet() unexpected error: %v", err)
	}
	if r.TargetIPSubnet != "24" {
		t.Errorf("got %q, want %q", r.TargetIPSubnet, "24")
	}
}

func TestNATOutboundRule_SetNATPort(t *testing.T) {
	var r NATOutboundRule
	if err := r.SetNATPort("1024"); err != nil {
		t.Errorf("SetNATPort() unexpected error: %v", err)
	}
	if r.NATPort != "1024" {
		t.Errorf("got %q, want %q", r.NATPort, "1024")
	}
}

func TestNATOutboundRule_SetSourceHashKey(t *testing.T) {
	var r NATOutboundRule
	if err := r.SetSourceHashKey("0x12345678"); err != nil {
		t.Errorf("SetSourceHashKey() unexpected error: %v", err)
	}
	if r.SourceHashKey != "0x12345678" {
		t.Errorf("got %q, want %q", r.SourceHashKey, "0x12345678")
	}
}

func TestNATOutboundRule_SetDescription(t *testing.T) {
	var r NATOutboundRule
	if err := r.SetDescription("test rule"); err != nil {
		t.Errorf("SetDescription() unexpected error: %v", err)
	}
	if r.Description != "test rule" {
		t.Errorf("got %q, want %q", r.Description, "test rule")
	}
}

func TestNATOutboundRule_SetBooleans(t *testing.T) {
	boolSetters := []struct {
		name   string
		setter func(*NATOutboundRule, bool) error
		getter func(*NATOutboundRule) bool
	}{
		{"SourceNot", (*NATOutboundRule).SetSourceNot, func(r *NATOutboundRule) bool { return r.SourceNot }},
		{"DestNot", (*NATOutboundRule).SetDestNot, func(r *NATOutboundRule) bool { return r.DestNot }},
		{"StaticNATPort", (*NATOutboundRule).SetStaticNATPort, func(r *NATOutboundRule) bool { return r.StaticNATPort }},
		{"NoSync", (*NATOutboundRule).SetNoSync, func(r *NATOutboundRule) bool { return r.NoSync }},
		{"NoNAT", (*NATOutboundRule).SetNoNAT, func(r *NATOutboundRule) bool { return r.NoNAT }},
		{"Disabled", (*NATOutboundRule).SetDisabled, func(r *NATOutboundRule) bool { return r.Disabled }},
	}

	for _, bs := range boolSetters {
		t.Run(bs.name+"_true", func(t *testing.T) {
			var r NATOutboundRule
			if err := bs.setter(&r, true); err != nil {
				t.Errorf("Set%s(true) unexpected error: %v", bs.name, err)
			}
			if !bs.getter(&r) {
				t.Errorf("Set%s(true) got false", bs.name)
			}
		})
		t.Run(bs.name+"_false", func(t *testing.T) {
			var r NATOutboundRule
			if err := bs.setter(&r, false); err != nil {
				t.Errorf("Set%s(false) unexpected error: %v", bs.name, err)
			}
			if bs.getter(&r) {
				t.Errorf("Set%s(false) got true", bs.name)
			}
		})
	}
}
