package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

type IPsecPhase2sModel struct {
	All types.List `tfsdk:"all"`
}

type IPsecPhase2Model struct {
	UniqID                    types.String `tfsdk:"uniq_id"`
	IKEId                     types.String `tfsdk:"ike_id"`
	Mode                      types.String `tfsdk:"mode"`
	ReqID                     types.String `tfsdk:"req_id"`
	LocalID                   types.Object `tfsdk:"local_id"`
	RemoteID                  types.Object `tfsdk:"remote_id"`
	NATLocalID                types.Object `tfsdk:"nat_local_id"`
	Protocol                  types.String `tfsdk:"protocol"`
	EncryptionAlgorithmOption types.List   `tfsdk:"encryption_algorithm_option"`
	HashAlgorithmOption       types.List   `tfsdk:"hash_algorithm_option"`
	PFSGroup                  types.String `tfsdk:"pfs_group"`
	Lifetime                  types.String `tfsdk:"lifetime"`
	RekeyTime                 types.String `tfsdk:"rekey_time"`
	RandTime                  types.String `tfsdk:"rand_time"`
	PingHost                  types.String `tfsdk:"ping_host"`
	Keepalive                 types.String `tfsdk:"keepalive"`
	Description               types.String `tfsdk:"description"`
	Disabled                  types.Bool   `tfsdk:"disabled"`
	Mobile                    types.Bool   `tfsdk:"mobile"`
}

type IPsecPhase2IDModel struct {
	Type    types.String `tfsdk:"type"`
	Address types.String `tfsdk:"address"`
	NetBits types.String `tfsdk:"net_bits"`
}

type IPsecPhase2EncryptionAlgorithmModel struct {
	Name   types.String `tfsdk:"name"`
	KeyLen types.String `tfsdk:"key_length"`
}

func (IPsecPhase2Model) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"uniq_id": {
			Description: "Unique identifier for this phase 2 entry.",
		},
		"ike_id": {
			Description: "IKE identifier linking this phase 2 to a phase 1 entry.",
		},
		"mode": {
			Description: "IPsec mode (tunnel, transport, vti).",
		},
		"req_id": {
			Description: "Request ID for the phase 2 SA.",
		},
		"local_id": {
			Description: "Local network identifier (type, address, netbits).",
		},
		"remote_id": {
			Description: "Remote network identifier (type, address, netbits).",
		},
		"nat_local_id": {
			Description: "NAT/BINAT local network identifier (optional).",
		},
		"protocol": {
			Description: "IPsec protocol (esp, ah).",
		},
		"encryption_algorithm_option": {
			Description: "Encryption algorithm options.",
		},
		"hash_algorithm_option": {
			Description: "Hash algorithm options.",
		},
		"pfs_group": {
			Description: "PFS key group number.",
		},
		"lifetime": {
			Description: "SA lifetime in seconds.",
		},
		"rekey_time": {
			Description: "Rekey time in seconds.",
		},
		"rand_time": {
			Description: "Random time range in seconds.",
		},
		"ping_host": {
			Description: "IP address to ping for keepalive.",
		},
		"keepalive": {
			Description: "Keepalive setting (enabled/disabled).",
		},
		"description": {
			Description: descriptionDescription,
		},
		"disabled": {
			Description: "Whether this phase 2 entry is disabled.",
		},
		"mobile": {
			Description: "Whether this is a mobile client phase 2 entry.",
		},
	}
}

func (IPsecPhase2IDModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"type": {
			Description: "Network type (address, network, lan, etc.).",
		},
		"address": {
			Description: "Network address.",
		},
		"net_bits": {
			Description: "Network prefix length (CIDR bits).",
		},
	}
}

func (IPsecPhase2EncryptionAlgorithmModel) descriptions() map[string]attrDescription {
	return map[string]attrDescription{
		"name": {
			Description: "Algorithm name.",
		},
		"key_length": {
			Description: "Key length in bits.",
		},
	}
}

func (IPsecPhase2IDModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":     types.StringType,
		"address":  types.StringType,
		"net_bits": types.StringType,
	}
}

func (IPsecPhase2EncryptionAlgorithmModel) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":       types.StringType,
		"key_length": types.StringType,
	}
}

func (IPsecPhase2Model) AttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"uniq_id":                     types.StringType,
		"ike_id":                      types.StringType,
		"mode":                        types.StringType,
		"req_id":                      types.StringType,
		"local_id":                    types.ObjectType{AttrTypes: IPsecPhase2IDModel{}.AttrTypes()},
		"remote_id":                   types.ObjectType{AttrTypes: IPsecPhase2IDModel{}.AttrTypes()},
		"nat_local_id":                types.ObjectType{AttrTypes: IPsecPhase2IDModel{}.AttrTypes()},
		"protocol":                    types.StringType,
		"encryption_algorithm_option": types.ListType{ElemType: types.ObjectType{AttrTypes: IPsecPhase2EncryptionAlgorithmModel{}.AttrTypes()}},
		"hash_algorithm_option":       types.ListType{ElemType: types.StringType},
		"pfs_group":                   types.StringType,
		"lifetime":                    types.StringType,
		"rekey_time":                  types.StringType,
		"rand_time":                   types.StringType,
		"ping_host":                   types.StringType,
		"keepalive":                   types.StringType,
		"description":                 types.StringType,
		"disabled":                    types.BoolType,
		"mobile":                      types.BoolType,
	}
}

func setIPsecPhase2IDModel(id pfsense.IPsecPhase2ID) IPsecPhase2IDModel {
	return IPsecPhase2IDModel{
		Type:    types.StringValue(id.Type),
		Address: types.StringValue(id.Address),
		NetBits: types.StringValue(id.NetBits),
	}
}

func (m *IPsecPhase2sModel) Set(ctx context.Context, phase2s pfsense.IPsecPhase2s) diag.Diagnostics {
	var diags diag.Diagnostics

	phase2Models := []IPsecPhase2Model{}
	for _, phase2 := range phase2s {
		var phase2Model IPsecPhase2Model
		diags.Append(phase2Model.Set(ctx, phase2)...)
		phase2Models = append(phase2Models, phase2Model)
	}

	phase2sValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: IPsecPhase2Model{}.AttrTypes()}, phase2Models)
	diags.Append(newDiags...)
	m.All = phase2sValue

	return diags
}

func (m *IPsecPhase2Model) Set(ctx context.Context, p2 pfsense.IPsecPhase2) diag.Diagnostics {
	var diags diag.Diagnostics

	m.UniqID = types.StringValue(p2.UniqID)
	m.IKEId = types.StringValue(p2.IKEId)
	m.Mode = types.StringValue(p2.Mode)
	m.Protocol = types.StringValue(p2.Protocol)
	m.PFSGroup = types.StringValue(p2.PFSGroup)
	m.Lifetime = types.StringValue(p2.Lifetime)
	m.Disabled = types.BoolValue(p2.Disabled)
	m.Mobile = types.BoolValue(p2.Mobile)

	if p2.ReqID != "" {
		m.ReqID = types.StringValue(p2.ReqID)
	}
	if p2.RekeyTime != "" {
		m.RekeyTime = types.StringValue(p2.RekeyTime)
	}
	if p2.RandTime != "" {
		m.RandTime = types.StringValue(p2.RandTime)
	}
	if p2.PingHost != "" {
		m.PingHost = types.StringValue(p2.PingHost)
	}
	if p2.Keepalive != "" {
		m.Keepalive = types.StringValue(p2.Keepalive)
	}
	if p2.Description != "" {
		m.Description = types.StringValue(p2.Description)
	}

	// Local ID
	localIDModel := setIPsecPhase2IDModel(p2.LocalID)
	localIDValue, newDiags := types.ObjectValueFrom(ctx, IPsecPhase2IDModel{}.AttrTypes(), localIDModel)
	diags.Append(newDiags...)
	m.LocalID = localIDValue

	// Remote ID
	remoteIDModel := setIPsecPhase2IDModel(p2.RemoteID)
	remoteIDValue, newDiags := types.ObjectValueFrom(ctx, IPsecPhase2IDModel{}.AttrTypes(), remoteIDModel)
	diags.Append(newDiags...)
	m.RemoteID = remoteIDValue

	// NAT Local ID (optional)
	if p2.NATLocalID != nil {
		natLocalIDModel := setIPsecPhase2IDModel(*p2.NATLocalID)
		natLocalIDValue, newDiags := types.ObjectValueFrom(ctx, IPsecPhase2IDModel{}.AttrTypes(), natLocalIDModel)
		diags.Append(newDiags...)
		m.NATLocalID = natLocalIDValue
	} else {
		m.NATLocalID = types.ObjectNull(IPsecPhase2IDModel{}.AttrTypes())
	}

	// Encryption algorithms
	encModels := []IPsecPhase2EncryptionAlgorithmModel{}
	for _, alg := range p2.EncryptionAlgorithmOption {
		encModels = append(encModels, IPsecPhase2EncryptionAlgorithmModel{
			Name:   types.StringValue(alg.Name),
			KeyLen: types.StringValue(alg.KeyLen),
		})
	}

	encValue, newDiags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: IPsecPhase2EncryptionAlgorithmModel{}.AttrTypes()}, encModels)
	diags.Append(newDiags...)
	m.EncryptionAlgorithmOption = encValue

	// Hash algorithms
	hashValue, newDiags := types.ListValueFrom(ctx, types.StringType, p2.HashAlgorithmOption)
	diags.Append(newDiags...)
	m.HashAlgorithmOption = hashValue

	return diags
}

func (m IPsecPhase2Model) Value(ctx context.Context, p2 *pfsense.IPsecPhase2) diag.Diagnostics {
	var diags diag.Diagnostics

	p2.IKEId = m.IKEId.ValueString()
	p2.Mode = m.Mode.ValueString()
	p2.Protocol = m.Protocol.ValueString()
	p2.PFSGroup = m.PFSGroup.ValueString()
	p2.Lifetime = m.Lifetime.ValueString()
	p2.Disabled = m.Disabled.ValueBool()
	p2.Mobile = m.Mobile.ValueBool()

	if !m.UniqID.IsNull() && !m.UniqID.IsUnknown() {
		p2.UniqID = m.UniqID.ValueString()
	}
	if !m.ReqID.IsNull() {
		p2.ReqID = m.ReqID.ValueString()
	}
	if !m.RekeyTime.IsNull() {
		p2.RekeyTime = m.RekeyTime.ValueString()
	}
	if !m.RandTime.IsNull() {
		p2.RandTime = m.RandTime.ValueString()
	}
	if !m.PingHost.IsNull() {
		p2.PingHost = m.PingHost.ValueString()
	}
	if !m.Keepalive.IsNull() {
		p2.Keepalive = m.Keepalive.ValueString()
	}
	if !m.Description.IsNull() {
		p2.Description = m.Description.ValueString()
	}

	// Local ID
	if !m.LocalID.IsNull() {
		var localIDModel IPsecPhase2IDModel
		diags.Append(m.LocalID.As(ctx, &localIDModel, basetypes.ObjectAsOptions{})...)
		p2.LocalID = pfsense.IPsecPhase2ID{
			Type:    localIDModel.Type.ValueString(),
			Address: localIDModel.Address.ValueString(),
			NetBits: localIDModel.NetBits.ValueString(),
		}
	}

	// Remote ID
	if !m.RemoteID.IsNull() {
		var remoteIDModel IPsecPhase2IDModel
		diags.Append(m.RemoteID.As(ctx, &remoteIDModel, basetypes.ObjectAsOptions{})...)
		p2.RemoteID = pfsense.IPsecPhase2ID{
			Type:    remoteIDModel.Type.ValueString(),
			Address: remoteIDModel.Address.ValueString(),
			NetBits: remoteIDModel.NetBits.ValueString(),
		}
	}

	// NAT Local ID
	if !m.NATLocalID.IsNull() {
		var natLocalIDModel IPsecPhase2IDModel
		diags.Append(m.NATLocalID.As(ctx, &natLocalIDModel, basetypes.ObjectAsOptions{})...)
		natLocalID := pfsense.IPsecPhase2ID{
			Type:    natLocalIDModel.Type.ValueString(),
			Address: natLocalIDModel.Address.ValueString(),
			NetBits: natLocalIDModel.NetBits.ValueString(),
		}
		p2.NATLocalID = &natLocalID
	}

	// Encryption algorithms
	var encModels []IPsecPhase2EncryptionAlgorithmModel
	if !m.EncryptionAlgorithmOption.IsNull() {
		diags.Append(m.EncryptionAlgorithmOption.ElementsAs(ctx, &encModels, false)...)
	}

	p2.EncryptionAlgorithmOption = make([]pfsense.IPsecPhase2EncryptionAlgorithm, 0, len(encModels))
	for _, encModel := range encModels {
		p2.EncryptionAlgorithmOption = append(p2.EncryptionAlgorithmOption, pfsense.IPsecPhase2EncryptionAlgorithm{
			Name:   encModel.Name.ValueString(),
			KeyLen: encModel.KeyLen.ValueString(),
		})
	}

	// Hash algorithms
	if !m.HashAlgorithmOption.IsNull() {
		diags.Append(m.HashAlgorithmOption.ElementsAs(ctx, &p2.HashAlgorithmOption, false)...)
	}

	return diags
}
