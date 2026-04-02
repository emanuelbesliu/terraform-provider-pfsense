package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*IPsecPhase2Resource)(nil)
	_ resource.ResourceWithConfigure   = (*IPsecPhase2Resource)(nil)
	_ resource.ResourceWithImportState = (*IPsecPhase2Resource)(nil)
)

type IPsecPhase2ResourceModel struct {
	IPsecPhase2Model
	Apply types.Bool `tfsdk:"apply"`
}

func NewIPsecPhase2Resource() resource.Resource { //nolint:ireturn
	return &IPsecPhase2Resource{}
}

type IPsecPhase2Resource struct {
	client *pfsense.Client
}

func (r *IPsecPhase2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_ipsec_phase2", req.ProviderTypeName)
}

func ipsecPhase2IDSchemaAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: IPsecPhase2IDModel{}.descriptions()["type"].Description,
			Required:    true,
		},
		"address": schema.StringAttribute{
			Description: IPsecPhase2IDModel{}.descriptions()["address"].Description,
			Optional:    true,
		},
		"net_bits": schema.StringAttribute{
			Description: IPsecPhase2IDModel{}.descriptions()["net_bits"].Description,
			Optional:    true,
		},
	}
}

func ipsecPhase2IDSchemaAttributesComputed() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"type": schema.StringAttribute{
			Description: IPsecPhase2IDModel{}.descriptions()["type"].Description,
			Computed:    true,
		},
		"address": schema.StringAttribute{
			Description: IPsecPhase2IDModel{}.descriptions()["address"].Description,
			Computed:    true,
		},
		"net_bits": schema.StringAttribute{
			Description: IPsecPhase2IDModel{}.descriptions()["net_bits"].Description,
			Computed:    true,
		},
	}
}

func (r *IPsecPhase2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "IPsec Phase 2 tunnel configuration. Defines the child SA / traffic selector parameters for an IPsec VPN tunnel.",
		MarkdownDescription: "IPsec [Phase 2](https://docs.netgate.com/pfsense/en/latest/vpn/ipsec/configure.html) tunnel configuration. Defines the child SA / traffic selector parameters for an IPsec VPN tunnel.",
		Attributes: map[string]schema.Attribute{
			"uniq_id": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["uniq_id"].Description,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ike_id": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["ike_id"].Description,
				Required:    true,
			},
			"mode": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["mode"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("tunnel", "tunnel6", "transport", "vti"),
				},
			},
			"req_id": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["req_id"].Description,
				Optional:    true,
			},
			"local_id": schema.SingleNestedAttribute{
				Description: IPsecPhase2Model{}.descriptions()["local_id"].Description,
				Required:    true,
				Attributes:  ipsecPhase2IDSchemaAttributes(),
			},
			"remote_id": schema.SingleNestedAttribute{
				Description: IPsecPhase2Model{}.descriptions()["remote_id"].Description,
				Required:    true,
				Attributes:  ipsecPhase2IDSchemaAttributes(),
			},
			"nat_local_id": schema.SingleNestedAttribute{
				Description: IPsecPhase2Model{}.descriptions()["nat_local_id"].Description,
				Optional:    true,
				Attributes:  ipsecPhase2IDSchemaAttributes(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"protocol": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["protocol"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("esp", "ah"),
				},
			},
			"encryption_algorithm_option": schema.ListNestedAttribute{
				Description: IPsecPhase2Model{}.descriptions()["encryption_algorithm_option"].Description,
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: IPsecPhase2EncryptionAlgorithmModel{}.descriptions()["name"].Description,
							Required:    true,
						},
						"key_length": schema.StringAttribute{
							Description: IPsecPhase2EncryptionAlgorithmModel{}.descriptions()["key_length"].Description,
							Optional:    true,
						},
					},
				},
			},
			"hash_algorithm_option": schema.ListAttribute{
				Description: IPsecPhase2Model{}.descriptions()["hash_algorithm_option"].Description,
				Optional:    true,
				ElementType: types.StringType,
			},
			"pfs_group": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["pfs_group"].Description,
				Required:    true,
			},
			"lifetime": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["lifetime"].Description,
				Required:    true,
			},
			"rekey_time": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["rekey_time"].Description,
				Optional:    true,
			},
			"rand_time": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["rand_time"].Description,
				Optional:    true,
			},
			"ping_host": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["ping_host"].Description,
				Optional:    true,
			},
			"keepalive": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["keepalive"].Description,
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: IPsecPhase2Model{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
			"disabled": schema.BoolAttribute{
				Description: IPsecPhase2Model{}.descriptions()["disabled"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"mobile": schema.BoolAttribute{
				Description: IPsecPhase2Model{}.descriptions()["mobile"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
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

func (r *IPsecPhase2Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *IPsecPhase2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *IPsecPhase2ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var p2Req pfsense.IPsecPhase2
	resp.Diagnostics.Append(data.Value(ctx, &p2Req)...)

	if resp.Diagnostics.HasError() {
		return
	}

	p2, err := r.client.CreateIPsecPhase2(ctx, p2Req)
	if addError(&resp.Diagnostics, "Error creating IPsec phase 2", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *p2)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyIPsecChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying IPsec changes", err)
	}
}

func (r *IPsecPhase2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *IPsecPhase2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	p2, err := r.client.GetIPsecPhase2(ctx, data.UniqID.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading IPsec phase 2", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *p2)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *IPsecPhase2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *IPsecPhase2ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var p2Req pfsense.IPsecPhase2
	resp.Diagnostics.Append(data.Value(ctx, &p2Req)...)

	if resp.Diagnostics.HasError() {
		return
	}

	p2, err := r.client.UpdateIPsecPhase2(ctx, p2Req)
	if addError(&resp.Diagnostics, "Error updating IPsec phase 2", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *p2)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyIPsecChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying IPsec changes", err)
	}
}

func (r *IPsecPhase2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *IPsecPhase2ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIPsecPhase2(ctx, data.UniqID.ValueString())
	if addError(&resp.Diagnostics, "Error deleting IPsec phase 2", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyIPsecChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying IPsec changes", err)
	}
}

func (r *IPsecPhase2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by uniqid
	uniqID := req.ID

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uniq_id"), types.StringValue(uniqID))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(defaultApply))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("disabled"), types.BoolValue(false))...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("mobile"), types.BoolValue(false))...)
}
