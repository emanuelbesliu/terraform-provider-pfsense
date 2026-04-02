package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type IPsecPhase1sModel struct {
	All types.List `tfsdk:"all"`
}

type IPsecPhase1Model struct {
	IKEId                types.String `tfsdk:"ike_id"`
	IKEType              types.String `tfsdk:"ike_type"`
	Interface            types.String `tfsdk:"interface"`
	Protocol             types.String `tfsdk:"protocol"`
	RemoteGateway        types.String `tfsdk:"remote_gateway"`
	AuthenticationMethod types.String `tfsdk:"authentication_method"`
	PreSharedKey         types.String `tfsdk:"pre_shared_key"`
	MyIDType             types.String `tfsdk:"my_id_type"`
	MyIDData             types.String `tfsdk:"my_id_data"`
	PeerIDType           types.String `tfsdk:"peer_id_type"`
	PeerIDData           types.String `tfsdk:"peer_id_data"`
	Description          types.String `tfsdk:"description"`
	NATTraversal         types.String `tfsdk:"nat_traversal"`
	Mobike               types.String `tfsdk:"mobike"`
	DPDDelay             types.String `tfsdk:"dpd_delay"`
	DPDMaxFail           types.String `tfsdk:"dpd_max_fail"`
	Lifetime             types.String `tfsdk:"lifetime"`
	RekeyTime            types.String `tfsdk:"rekey_time"`
	ReauthTime           types.String `tfsdk:"reauth_time"`
	RandTime             types.String `tfsdk:"rand_time"`
	StartAction          types.String `tfsdk:"start_action"`
	CloseAction          types.String `tfsdk:"close_action"`
	Encryption           types.List   `tfsdk:"encryption"`
	Disabled             types.Bool   `tfsdk:"disabled"`
	CertRef              types.String `tfsdk:"cert_ref"`
	CARef                types.String `tfsdk:"ca_ref"`
	PKCS11CertRef        types.String `tfsdk:"pkcs11_cert_ref"`
	PKCS11Pin            types.String `tfsdk:"pkcs11_pin"`
	Mobile               types.Bool   `tfsdk:"mobile"`
	IKEPort              types.String `tfsdk:"ike_port"`
	NATTPort             types.String `tfsdk:"natt_port"`
	GWDuplicates         types.Bool   `tfsdk:"gw_duplicates"`
	PRFSelectEnable      types.Bool   `tfsdk:"prf_select_enable"`
	SplitConn            types.Bool   `tfsdk:"split_conn"`
	TFCEnable            types.Bool   `tfsdk:"tfc_enable"`
	TFCBytes             types.String `tfsdk:"tfc_bytes"`
}

type IPsecPhase1EncryptionModel struct {
	Algorithm    types.String `tfsdk:"algorithm"`
	KeyLen       types.String `tfsdk:"key_length"`
	HashAlgo     types.String `tfsdk:"hash_algorithm"`
	PRFAlgorithm types.String `tfsdk:"prf_algorithm"`
	DHGroup      types.String `tfsdk:"dh_group"`
}

func (IPsecPhase1Model) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"ike_id": {
			Description: "Unique IKE identifier for this phase 1 entry.",
		},
		"ike_type": {
			Description: "IKE version (ikev1, ikev2, auto).",
		},
		"interface": {
			Description: "Interface for the IPsec tunnel endpoint.",
		},
		"protocol": {
			Description: "IP protocol version (inet, inet6, both).",
		},
		"remote_gateway": {
			Description: "IP address or hostname of the remote gateway.",
		},
		"authentication_method": {
			Description: "Authentication method (pre_shared_key, cert, etc.).",
		},
		"pre_shared_key": {
			Description: "Pre-shared key for authentication.",
		},
		"my_id_type": {
			Description: "Local identifier type (myaddress, fqdn, user_fqdn, dyn_dns, etc.).",
		},
		"my_id_data": {
			Description: "Local identifier data.",
		},
		"peer_id_type": {
			Description: "Peer identifier type (peeraddress, fqdn, user_fqdn, etc.).",
		},
		"peer_id_data": {
			Description: "Peer identifier data.",
		},
		"description": {
			Description: descriptionDescription,
		},
		"nat_traversal": {
			Description: "NAT traversal setting (on, force).",
		},
		"mobike": {
			Description: "MOBIKE protocol setting (on, off).",
		},
		"dpd_delay": {
			Description: "Dead peer detection delay in seconds.",
		},
		"dpd_max_fail": {
			Description: "Dead peer detection maximum failures.",
		},
		"lifetime": {
			Description: "SA lifetime in seconds.",
		},
		"rekey_time": {
			Description: "Rekey time in seconds.",
		},
		"reauth_time": {
			Description: "Re-authentication time in seconds.",
		},
		"rand_time": {
			Description: "Random time range in seconds.",
		},
		"start_action": {
			Description: "Start action (none, start, trap).",
		},
		"close_action": {
			Description: "Close action (none, start, trap).",
		},
		"encryption": {
			Description: "Encryption algorithm configuration.",
		},
		"disabled": {
			Description: "Whether this phase 1 entry is disabled.",
		},
		"cert_ref": {
			Description: "Certificate reference for certificate-based authentication.",
		},
		"ca_ref": {
			Description: "CA reference for certificate-based authentication.",
		},
		"pkcs11_cert_ref": {
			Description: "PKCS#11 certificate reference.",
		},
		"pkcs11_pin": {
			Description: "PKCS#11 PIN.",
		},
		"mobile": {
			Description: "Whether this is a mobile client tunnel.",
		},
		"ike_port": {
			Description: "Custom IKE port (default 500).",
		},
		"natt_port": {
			Description: "Custom NAT-T port (default 4500).",
		},
		"gw_duplicates": {
			Description: "Allow multiple phase 1 entries with the same remote gateway.",
		},
		"prf_select_enable": {
			Description: "Enable PRF selection.",
		},
		"split_conn": {
			Description: "Split connections.",
		},
		"tfc_enable": {
			Description: "Enable Traffic Flow Confidentiality.",
		},
		"tfc_bytes": {
			Description: "TFC bytes value.",
		},
	}
}

func (IPsecPhase1EncryptionModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"algorithm": {
			Description: "Encryption algorithm name.",
		},
		"key_length": {
			Description: "Key length in bits.",
		},
		"hash_algorithm": {
			Description: "Hash algorithm.",
		},
		"prf_algorithm": {
			Description: "PRF algorithm.",
		},
		"dh_group": {
			Description: "Diffie-Hellman group number.",
		},
	}
}

func (IPsecPhase1EncryptionModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"algorithm":      types.StringType,
		"key_length":     types.StringType,
		"hash_algorithm": types.StringType,
		"prf_algorithm":  types.StringType,
		"dh_group":       types.StringType,
	}
}

func (IPsecPhase1Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"ike_id":                types.StringType,
		"ike_type":              types.StringType,
		"interface":             types.StringType,
		"protocol":              types.StringType,
		"remote_gateway":        types.StringType,
		"authentication_method": types.StringType,
		"pre_shared_key":        types.StringType,
		"my_id_type":            types.StringType,
		"my_id_data":            types.StringType,
		"peer_id_type":          types.StringType,
		"peer_id_data":          types.StringType,
		"description":           types.StringType,
		"nat_traversal":         types.StringType,
		"mobike":                types.StringType,
		"dpd_delay":             types.StringType,
		"dpd_max_fail":          types.StringType,
		"lifetime":              types.StringType,
		"rekey_time":            types.StringType,
		"reauth_time":           types.StringType,
		"rand_time":             types.StringType,
		"start_action":          types.StringType,
		"close_action":          types.StringType,
		"encryption":            types.ListType{ElemType: types.ObjectType{AttrTypes: IPsecPhase1EncryptionModel{}.AttrTypes()}},
		"disabled":              types.BoolType,
		"cert_ref":              types.StringType,
		"ca_ref":                types.StringType,
		"pkcs11_cert_ref":       types.StringType,
		"pkcs11_pin":            types.StringType,
		"mobile":                types.BoolType,
		"ike_port":              types.StringType,
		"natt_port":             types.StringType,
		"gw_duplicates":         types.BoolType,
		"prf_select_enable":     types.BoolType,
		"split_conn":            types.BoolType,
		"tfc_enable":            types.BoolType,
		"tfc_bytes":             types.StringType,
	}
}

func (m *IPsecPhase1sModel) Set(ctx context.Context, phase1s pfsense.IPsecPhase1s) diag.Diagnostics {
	var diags diag.Diagnostics

	phase1Models := []IPsecPhase1Model{}
	for _, phase1 := range phase1s {
		var phase1Model IPsecPhase1Model
		diags.Append(phase1Model.Set(ctx, phase1)...)
		phase1Models = append(phase1Models, phase1Model)
	}

	phase1sValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: IPsecPhase1Model{}.AttrTypes()}, phase1Models)
	diags.Append(newDiags...)
	m.All = phase1sValue

	return diags
}

func (m *IPsecPhase1Model) Set(ctx context.Context, p1 pfsense.IPsecPhase1) diag.Diagnostics {
	var diags diag.Diagnostics

	m.IKEId = types.StringValue(p1.IKEId)
	m.IKEType = types.StringValue(p1.IKEType)
	m.Interface = types.StringValue(p1.Interface)
	m.Protocol = types.StringValue(p1.Protocol)
	m.RemoteGateway = types.StringValue(p1.RemoteGateway)
	m.AuthenticationMethod = types.StringValue(p1.AuthenticationMethod)
	m.PreSharedKey = types.StringValue(p1.PreSharedKey)
	m.MyIDType = types.StringValue(p1.MyIDType)
	m.MyIDData = types.StringValue(p1.MyIDData)
	m.PeerIDType = types.StringValue(p1.PeerIDType)
	m.PeerIDData = types.StringValue(p1.PeerIDData)
	m.NATTraversal = types.StringValue(p1.NATTraversal)
	m.Mobike = types.StringValue(p1.Mobike)
	m.Lifetime = types.StringValue(p1.Lifetime)
	m.Disabled = types.BoolValue(p1.Disabled)
	m.Mobile = types.BoolValue(p1.Mobile)
	m.GWDuplicates = types.BoolValue(p1.GWDuplicates)
	m.PRFSelectEnable = types.BoolValue(p1.PRFSelectEnable)
	m.SplitConn = types.BoolValue(p1.SplitConn)
	m.TFCEnable = types.BoolValue(p1.TFCEnable)

	if p1.Description != "" {
		m.Description = types.StringValue(p1.Description)
	}
	if p1.DPDDelay != "" {
		m.DPDDelay = types.StringValue(p1.DPDDelay)
	}
	if p1.DPDMaxFail != "" {
		m.DPDMaxFail = types.StringValue(p1.DPDMaxFail)
	}
	if p1.RekeyTime != "" {
		m.RekeyTime = types.StringValue(p1.RekeyTime)
	}
	if p1.ReauthTime != "" {
		m.ReauthTime = types.StringValue(p1.ReauthTime)
	}
	if p1.RandTime != "" {
		m.RandTime = types.StringValue(p1.RandTime)
	}
	if p1.StartAction != "" {
		m.StartAction = types.StringValue(p1.StartAction)
	}
	if p1.CloseAction != "" {
		m.CloseAction = types.StringValue(p1.CloseAction)
	}
	if p1.CertRef != "" {
		m.CertRef = types.StringValue(p1.CertRef)
	}
	if p1.CARef != "" {
		m.CARef = types.StringValue(p1.CARef)
	}
	if p1.PKCS11CertRef != "" {
		m.PKCS11CertRef = types.StringValue(p1.PKCS11CertRef)
	}
	if p1.PKCS11Pin != "" {
		m.PKCS11Pin = types.StringValue(p1.PKCS11Pin)
	}
	if p1.IKEPort != "" {
		m.IKEPort = types.StringValue(p1.IKEPort)
	}
	if p1.NATTPort != "" {
		m.NATTPort = types.StringValue(p1.NATTPort)
	}
	if p1.TFCBytes != "" {
		m.TFCBytes = types.StringValue(p1.TFCBytes)
	}

	// Encryption
	encModels := []IPsecPhase1EncryptionModel{}
	for _, enc := range p1.Encryption {
		encModels = append(encModels, IPsecPhase1EncryptionModel{
			Algorithm:    types.StringValue(enc.Algorithm),
			KeyLen:       types.StringValue(enc.KeyLen),
			HashAlgo:     types.StringValue(enc.HashAlgo),
			PRFAlgorithm: types.StringValue(enc.PRFAlgorithm),
			DHGroup:      types.StringValue(enc.DHGroup),
		})
	}

	encValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: IPsecPhase1EncryptionModel{}.AttrTypes()}, encModels)
	diags.Append(newDiags...)
	m.Encryption = encValue

	return diags
}

func (m IPsecPhase1Model) Value(ctx context.Context, p1 *pfsense.IPsecPhase1) diag.Diagnostics {
	var diags diag.Diagnostics

	p1.IKEType = m.IKEType.ValueString()
	p1.Interface = m.Interface.ValueString()
	p1.Protocol = m.Protocol.ValueString()
	p1.RemoteGateway = m.RemoteGateway.ValueString()
	p1.AuthenticationMethod = m.AuthenticationMethod.ValueString()
	p1.PreSharedKey = m.PreSharedKey.ValueString()
	p1.MyIDType = m.MyIDType.ValueString()
	p1.MyIDData = m.MyIDData.ValueString()
	p1.PeerIDType = m.PeerIDType.ValueString()
	p1.PeerIDData = m.PeerIDData.ValueString()
	p1.NATTraversal = m.NATTraversal.ValueString()
	p1.Mobike = m.Mobike.ValueString()
	p1.Lifetime = m.Lifetime.ValueString()
	p1.Disabled = m.Disabled.ValueBool()
	p1.Mobile = m.Mobile.ValueBool()
	p1.GWDuplicates = m.GWDuplicates.ValueBool()
	p1.PRFSelectEnable = m.PRFSelectEnable.ValueBool()
	p1.SplitConn = m.SplitConn.ValueBool()
	p1.TFCEnable = m.TFCEnable.ValueBool()

	if !m.IKEId.IsNull() && !m.IKEId.IsUnknown() {
		p1.IKEId = m.IKEId.ValueString()
	}
	if !m.Description.IsNull() {
		p1.Description = m.Description.ValueString()
	}
	if !m.DPDDelay.IsNull() {
		p1.DPDDelay = m.DPDDelay.ValueString()
	}
	if !m.DPDMaxFail.IsNull() {
		p1.DPDMaxFail = m.DPDMaxFail.ValueString()
	}
	if !m.RekeyTime.IsNull() {
		p1.RekeyTime = m.RekeyTime.ValueString()
	}
	if !m.ReauthTime.IsNull() {
		p1.ReauthTime = m.ReauthTime.ValueString()
	}
	if !m.RandTime.IsNull() {
		p1.RandTime = m.RandTime.ValueString()
	}
	if !m.StartAction.IsNull() {
		p1.StartAction = m.StartAction.ValueString()
	}
	if !m.CloseAction.IsNull() {
		p1.CloseAction = m.CloseAction.ValueString()
	}
	if !m.CertRef.IsNull() {
		p1.CertRef = m.CertRef.ValueString()
	}
	if !m.CARef.IsNull() {
		p1.CARef = m.CARef.ValueString()
	}
	if !m.PKCS11CertRef.IsNull() {
		p1.PKCS11CertRef = m.PKCS11CertRef.ValueString()
	}
	if !m.PKCS11Pin.IsNull() {
		p1.PKCS11Pin = m.PKCS11Pin.ValueString()
	}
	if !m.IKEPort.IsNull() {
		p1.IKEPort = m.IKEPort.ValueString()
	}
	if !m.NATTPort.IsNull() {
		p1.NATTPort = m.NATTPort.ValueString()
	}
	if !m.TFCBytes.IsNull() {
		p1.TFCBytes = m.TFCBytes.ValueString()
	}

	// Encryption
	var encModels []IPsecPhase1EncryptionModel
	if !m.Encryption.IsNull() {
		diags.Append(m.Encryption.ElementsAs(ctx, &encModels, false)...)
	}

	p1.Encryption = make([]pfsense.IPsecPhase1EncryptionItem, 0, len(encModels))
	for _, encModel := range encModels {
		p1.Encryption = append(p1.Encryption, pfsense.IPsecPhase1EncryptionItem{
			Algorithm:    encModel.Algorithm.ValueString(),
			KeyLen:       encModel.KeyLen.ValueString(),
			HashAlgo:     encModel.HashAlgo.ValueString(),
			PRFAlgorithm: encModel.PRFAlgorithm.ValueString(),
			DHGroup:      encModel.DHGroup.ValueString(),
		})
	}

	return diags
}
