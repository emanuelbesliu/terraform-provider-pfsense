package pfsense

import (
	"testing"
)

func TestVirtualIP_SetMode(t *testing.T) {
	tests := []struct {
		name    string
		mode    string
		wantErr bool
	}{
		{"ipalias", VirtualIPModeIPAlias, false},
		{"carp", VirtualIPModeCarp, false},
		{"proxyarp", VirtualIPModeProxyARP, false},
		{"other", VirtualIPModeOther, false},
		{"invalid", "invalid", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var vip VirtualIP

			err := vip.SetMode(tt.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetMode(%q) error = %v, wantErr %v", tt.mode, err, tt.wantErr)
			}

			if err == nil && vip.Mode != tt.mode {
				t.Errorf("SetMode(%q) got %q", tt.mode, vip.Mode)
			}
		})
	}
}

func TestVirtualIP_SetInterface(t *testing.T) {
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
			var vip VirtualIP

			err := vip.SetInterface(tt.iface)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetInterface(%q) error = %v, wantErr %v", tt.iface, err, tt.wantErr)
			}

			if err == nil && vip.Interface != tt.iface {
				t.Errorf("SetInterface(%q) got %q", tt.iface, vip.Interface)
			}
		})
	}
}

func TestVirtualIP_SetVHID(t *testing.T) {
	tests := []struct {
		name    string
		vhid    *int
		wantErr bool
	}{
		{"nil", nil, false},
		{"min", intPtr(MinVHID), false},
		{"max", intPtr(MaxVHID), false},
		{"mid", intPtr(100), false},
		{"below_min", intPtr(MinVHID - 1), true},
		{"above_max", intPtr(MaxVHID + 1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var vip VirtualIP

			err := vip.SetVHID(tt.vhid)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetVHID(%v) error = %v, wantErr %v", tt.vhid, err, tt.wantErr)
			}
		})
	}
}

func TestVirtualIP_SetAdvSkew(t *testing.T) {
	tests := []struct {
		name    string
		advskew *int
		wantErr bool
	}{
		{"nil", nil, false},
		{"min", intPtr(MinAdvSkew), false},
		{"max", intPtr(MaxAdvSkew), false},
		{"below_min", intPtr(MinAdvSkew - 1), true},
		{"above_max", intPtr(MaxAdvSkew + 1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var vip VirtualIP

			err := vip.SetAdvSkew(tt.advskew)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetAdvSkew(%v) error = %v, wantErr %v", tt.advskew, err, tt.wantErr)
			}
		})
	}
}

func TestVirtualIP_SetAdvBase(t *testing.T) {
	tests := []struct {
		name    string
		advbase *int
		wantErr bool
	}{
		{"nil", nil, false},
		{"min", intPtr(MinAdvBase), false},
		{"max", intPtr(MaxAdvBase), false},
		{"below_min", intPtr(MinAdvBase - 1), true},
		{"above_max", intPtr(MaxAdvBase + 1), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var vip VirtualIP

			err := vip.SetAdvBase(tt.advbase)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetAdvBase(%v) error = %v, wantErr %v", tt.advbase, err, tt.wantErr)
			}
		})
	}
}

func TestVirtualIP_SetSubnet(t *testing.T) {
	tests := []struct {
		name    string
		subnet  string
		wantErr bool
	}{
		{"ipv4", "10.0.1.100", false},
		{"ipv6", "fd00::1", false},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var vip VirtualIP

			err := vip.SetSubnet(tt.subnet)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetSubnet(%q) error = %v, wantErr %v", tt.subnet, err, tt.wantErr)
			}

			if err == nil && vip.Subnet != tt.subnet {
				t.Errorf("SetSubnet(%q) got %q", tt.subnet, vip.Subnet)
			}
		})
	}
}

func TestVirtualIP_SetSubnetBits(t *testing.T) {
	tests := []struct {
		name    string
		bits    int
		wantErr bool
	}{
		{"min", 1, false},
		{"max", 128, false},
		{"ipv4_host", 32, false},
		{"ipv4_net", 24, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"above_max", 129, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var vip VirtualIP

			err := vip.SetSubnetBits(tt.bits)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetSubnetBits(%d) error = %v, wantErr %v", tt.bits, err, tt.wantErr)
			}

			if err == nil && vip.SubnetBits != tt.bits {
				t.Errorf("SetSubnetBits(%d) got %d", tt.bits, vip.SubnetBits)
			}
		})
	}
}

func TestParseVirtualIPResponse_IPAlias(t *testing.T) {
	resp := virtualIPResponse{
		Mode:        "ipalias",
		Interface:   "lan",
		Subnet:      "10.0.1.100",
		SubnetBits:  "32",
		Description: "Test VIP",
		Type:        "single",
		UniqueID:    "abc123",
		ControlID:   0,
	}

	vip, err := parseVirtualIPResponse(resp)
	if err != nil {
		t.Fatalf("parseVirtualIPResponse() error = %v", err)
	}

	if vip.Mode != VirtualIPModeIPAlias {
		t.Errorf("Mode = %q, want %q", vip.Mode, VirtualIPModeIPAlias)
	}

	if vip.Interface != "lan" {
		t.Errorf("Interface = %q, want %q", vip.Interface, "lan")
	}

	if vip.Subnet != "10.0.1.100" {
		t.Errorf("Subnet = %q, want %q", vip.Subnet, "10.0.1.100")
	}

	if vip.SubnetBits != 32 {
		t.Errorf("SubnetBits = %d, want 32", vip.SubnetBits)
	}

	if vip.Description != "Test VIP" {
		t.Errorf("Description = %q, want %q", vip.Description, "Test VIP")
	}

	if vip.UniqueID != "abc123" {
		t.Errorf("UniqueID = %q, want %q", vip.UniqueID, "abc123")
	}

	if vip.VHID != nil {
		t.Errorf("VHID = %v, want nil", vip.VHID)
	}
}

func TestParseVirtualIPResponse_CARP(t *testing.T) {
	resp := virtualIPResponse{
		Mode:        "carp",
		Interface:   "wan",
		VHID:        "10",
		AdvSkew:     "0",
		AdvBase:     "1",
		Password:    "secret",
		Subnet:      "203.0.113.1",
		SubnetBits:  "24",
		Description: "CARP VIP",
		Type:        "single",
		UniqueID:    "def456",
		ControlID:   1,
	}

	vip, err := parseVirtualIPResponse(resp)
	if err != nil {
		t.Fatalf("parseVirtualIPResponse() error = %v", err)
	}

	if vip.Mode != VirtualIPModeCarp {
		t.Errorf("Mode = %q, want %q", vip.Mode, VirtualIPModeCarp)
	}

	if vip.VHID == nil || *vip.VHID != 10 {
		t.Errorf("VHID = %v, want 10", vip.VHID)
	}

	if vip.AdvSkew == nil || *vip.AdvSkew != 0 {
		t.Errorf("AdvSkew = %v, want 0", vip.AdvSkew)
	}

	if vip.AdvBase == nil || *vip.AdvBase != 1 {
		t.Errorf("AdvBase = %v, want 1", vip.AdvBase)
	}

	if vip.Password != "secret" {
		t.Errorf("Password = %q, want %q", vip.Password, "secret")
	}
}

func TestParseVirtualIPResponse_InvalidSubnetBits(t *testing.T) {
	resp := virtualIPResponse{
		Mode:       "ipalias",
		Interface:  "lan",
		Subnet:     "10.0.1.1",
		SubnetBits: "notanumber",
		UniqueID:   "xyz",
	}

	_, err := parseVirtualIPResponse(resp)
	if err == nil {
		t.Fatal("parseVirtualIPResponse() expected error for invalid subnet_bits")
	}
}

func TestVirtualIPs_GetByUniqueID(t *testing.T) {
	vips := VirtualIPs{
		{UniqueID: "aaa", Mode: VirtualIPModeIPAlias, Interface: "lan", Subnet: "10.0.1.1", SubnetBits: 32},
		{UniqueID: "bbb", Mode: VirtualIPModeCarp, Interface: "wan", Subnet: "10.0.2.1", SubnetBits: 24},
	}

	vip, err := vips.GetByUniqueID("bbb")
	if err != nil {
		t.Fatalf("GetByUniqueID() error = %v", err)
	}

	if vip.Mode != VirtualIPModeCarp {
		t.Errorf("Mode = %q, want %q", vip.Mode, VirtualIPModeCarp)
	}

	_, err = vips.GetByUniqueID("nonexistent")
	if err == nil {
		t.Fatal("GetByUniqueID() expected error for nonexistent ID")
	}
}

func TestVirtualIPs_GetControlIDByUniqueID(t *testing.T) {
	vips := VirtualIPs{
		{UniqueID: "aaa", controlID: 0},
		{UniqueID: "bbb", controlID: 5},
	}

	cid, err := vips.GetControlIDByUniqueID("bbb")
	if err != nil {
		t.Fatalf("GetControlIDByUniqueID() error = %v", err)
	}

	if *cid != 5 {
		t.Errorf("controlID = %d, want 5", *cid)
	}

	_, err = vips.GetControlIDByUniqueID("nonexistent")
	if err == nil {
		t.Fatal("GetControlIDByUniqueID() expected error for nonexistent ID")
	}
}

func intPtr(i int) *int {
	return &i
}
