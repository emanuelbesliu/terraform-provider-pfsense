package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*IPsecPhase1Resource)(nil)
	_ resource.ResourceWithConfigure   = (*IPsecPhase1Resource)(nil)
	_ resource.ResourceWithImportState = (*IPsecPhase1Resource)(nil)
)

type IPsecPhase1ResourceModel struct {
	IPsecPhase1Model
	Apply types.Bool `tfsdk:"apply"`
}

func NewIPsecPhase1Resource() resource.Resource { //nolint:ireturn
	return &IPsecPhase1Resource{}
}

type IPsecPhase1Resource struct {
	client *pfsense.Client
}

func (r *IPsecPhase1Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_ipsec_phase1", req.ProviderTypeName)
}

func (r *IPsecPhase1Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "IPsec Phase 1 tunnel configuration. Defines the IKE SA parameters for establishing a VPN tunnel.",
		MarkdownDescription: "IPsec [Phase 1](https://docs.netgate.com/pfsense/en/latest/vpn/ipsec/configure.html) tunnel configuration. Defines the IKE SA parameters for establishing a VPN tunnel.",
		Attributes: map[string]schema.Attribute{
			"ike_id": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["ike_id"].Description,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ike_type": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["ike_type"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("ikev1", "ikev2", "auto"),
				},
			},
			"interface": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["interface"].Description,
				Required:    true,
			},
			"protocol": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["protocol"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("inet", "inet6", "both"),
				},
			},
			"remote_gateway": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["remote_gateway"].Description,
				Optional:    true,
			},
			"authentication_method": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["authentication_method"].Description,
				Required:    true,
			},
			"pre_shared_key": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["pre_shared_key"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"my_id_type": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["my_id_type"].Description,
				Required:    true,
			},
			"my_id_data": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["my_id_data"].Description,
				Optional:    true,
			},
			"peer_id_type": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["peer_id_type"].Description,
				Required:    true,
			},
			"peer_id_data": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["peer_id_data"].Description,
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
			"nat_traversal": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["nat_traversal"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("on", "force"),
				},
			},
			"mobike": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["mobike"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("on", "off"),
				},
			},
			"dpd_delay": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["dpd_delay"].Description,
				Optional:    true,
			},
			"dpd_max_fail": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["dpd_max_fail"].Description,
				Optional:    true,
			},
			"lifetime": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["lifetime"].Description,
				Required:    true,
			},
			"rekey_time": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["rekey_time"].Description,
				Optional:    true,
			},
			"reauth_time": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["reauth_time"].Description,
				Optional:    true,
			},
			"rand_time": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["rand_time"].Description,
				Optional:    true,
			},
			"start_action": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["start_action"].Description,
				Optional:    true,
			},
			"close_action": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["close_action"].Description,
				Optional:    true,
			},
			"encryption": schema.ListNestedAttribute{
				Description: IPsecPhase1Model{}.descriptions()["encryption"].Description,
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"algorithm": schema.StringAttribute{
							Description: IPsecPhase1EncryptionModel{}.descriptions()["algorithm"].Description,
							Required:    true,
						},
						"key_length": schema.StringAttribute{
							Description: IPsecPhase1EncryptionModel{}.descriptions()["key_length"].Description,
							Optional:    true,
						},
						"hash_algorithm": schema.StringAttribute{
							Description: IPsecPhase1EncryptionModel{}.descriptions()["hash_algorithm"].Description,
							Required:    true,
						},
						"prf_algorithm": schema.StringAttribute{
							Description: IPsecPhase1EncryptionModel{}.descriptions()["prf_algorithm"].Description,
							Optional:    true,
						},
						"dh_group": schema.StringAttribute{
							Description: IPsecPhase1EncryptionModel{}.descriptions()["dh_group"].Description,
							Required:    true,
						},
					},
				},
			},
			"disabled": schema.BoolAttribute{
				Description: IPsecPhase1Model{}.descriptions()["disabled"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"cert_ref": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["cert_ref"].Description,
				Optional:    true,
			},
			"ca_ref": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["ca_ref"].Description,
				Optional:    true,
			},
			"pkcs11_cert_ref": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["pkcs11_cert_ref"].Description,
				Optional:    true,
			},
			"pkcs11_pin": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["pkcs11_pin"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"mobile": schema.BoolAttribute{
				Description: IPsecPhase1Model{}.descriptions()["mobile"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ike_port": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["ike_port"].Description,
				Optional:    true,
			},
			"natt_port": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["natt_port"].Description,
				Optional:    true,
			},
			"gw_duplicates": schema.BoolAttribute{
				Description: IPsecPhase1Model{}.descriptions()["gw_duplicates"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"prf_select_enable": schema.BoolAttribute{
				Description: IPsecPhase1Model{}.descriptions()["prf_select_enable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"split_conn": schema.BoolAttribute{
				Description: IPsecPhase1Model{}.descriptions()["split_conn"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"tfc_enable": schema.BoolAttribute{
				Description: IPsecPhase1Model{}.descriptions()["tfc_enable"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"tfc_bytes": schema.StringAttribute{
				Description: IPsecPhase1Model{}.descriptions()["tfc_bytes"].Description,
				Optional:    true,
			},
			"apply": schema.BoolAttribute{
				Description:         applyDescription,
				MarkdownDescription: applyMarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             booldefault.StaticBool(defaultApply),
			},
		},
	}
}

func (r *IPsecPhase1Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *IPsecPhase1Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *IPsecPhase1ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var p1Req pfsense.IPsecPhase1
	resp.Diagnostics.Append(data.Value(ctx, &p1Req)...)

	if resp.Diagnostics.HasError() {
		return
	}

	p1, err := r.client.CreateIPsecPhase1(ctx, p1Req)
	if addError(&resp.Diagnostics, "Error creating IPsec phase 1", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *p1)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyIPsecChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying IPsec changes", err)
	}
}

func (r *IPsecPhase1Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *IPsecPhase1ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	p1, err := r.client.GetIPsecPhase1(ctx, data.IKEId.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading IPsec phase 1", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *p1)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IPsecPhase1Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *IPsecPhase1ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var p1Req pfsense.IPsecPhase1
	resp.Diagnostics.Append(data.Value(ctx, &p1Req)...)

	if resp.Diagnostics.HasError() {
		return
	}

	p1, err := r.client.UpdateIPsecPhase1(ctx, p1Req)
	if addError(&resp.Diagnostics, "Error updating IPsec phase 1", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *p1)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyIPsecChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying IPsec changes", err)
	}
}

func (r *IPsecPhase1Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *IPsecPhase1ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIPsecPhase1(ctx, data.IKEId.ValueString())
	if addError(&resp.Diagnostics, "Error deleting IPsec phase 1", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyIPsecChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying IPsec changes", err)
	}
}

func (r *IPsecPhase1Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by ikeid
	ikeId := req.ID

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ike_id"), types.StringValue(ikeId))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(defaultApply))...)

	// Set default values for boolean fields to avoid null issues during read
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("disabled"), types.BoolValue(false))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("mobile"), types.BoolValue(false))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("gw_duplicates"), types.BoolValue(false))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("prf_select_enable"), types.BoolValue(false))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("split_conn"), types.BoolValue(false))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tfc_enable"), types.BoolValue(false))...)

	// Set empty encryption list default
	encValue, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: IPsecPhase1EncryptionModel{}.AttrTypes()}, []attr.Value{})
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("encryption"), encValue)...)
}

func init() {
	// Ensure the default list value for encryption is valid
	_ = types.ListValueMust(types.ObjectType{AttrTypes: IPsecPhase1EncryptionModel{}.AttrTypes()}, []attr.Value{})
}
