package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*VLANResource)(nil)
	_ resource.ResourceWithConfigure   = (*VLANResource)(nil)
	_ resource.ResourceWithImportState = (*VLANResource)(nil)
)

type VLANResourceModel struct {
	VLANModel
}

func NewVLANResource() resource.Resource { //nolint:ireturn
	return &VLANResource{}
}

type VLANResource struct {
	client *pfsense.Client
}

func (r *VLANResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_vlan", req.ProviderTypeName)
}

func (r *VLANResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "VLAN interface. VLANs allow segmenting a physical network into multiple virtual networks using IEEE 802.1Q tagging.",
		MarkdownDescription: "[VLAN](https://docs.netgate.com/pfsense/en/latest/interfaces/vlan.html) interface. VLANs allow segmenting a physical network into multiple virtual networks using IEEE 802.1Q tagging.",
		Attributes: map[string]schema.Attribute{
			"parent_interface": schema.StringAttribute{
				Description: VLANModel{}.descriptions()["parent_interface"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringIsInterface(),
				},
			},
			"tag": schema.Int64Attribute{
				Description: VLANModel{}.descriptions()["tag"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
				Validators: []validator.Int64{
					int64validator.Between(int64(pfsense.MinVLANTag), int64(pfsense.MaxVLANTag)),
				},
			},
			"pcp": schema.Int64Attribute{
				Description: VLANModel{}.descriptions()["pcp"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(int64(pfsense.MinVLANPCP), int64(pfsense.MaxVLANPCP)),
				},
			},
			"description": schema.StringAttribute{
				Description: VLANModel{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
			},
			"vlan_interface": schema.StringAttribute{
				Description: VLANModel{}.descriptions()["vlan_interface"].Description,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *VLANResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *VLANResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *VLANResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var vlanReq pfsense.VLAN
	resp.Diagnostics.Append(data.Value(ctx, &vlanReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vlan, err := r.client.CreateVLAN(ctx, vlanReq)
	if addError(&resp.Diagnostics, "Error creating VLAN", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *vlan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VLANResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *VLANResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vlan, err := r.client.GetVLAN(ctx, data.VLANInterface.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading VLAN", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *vlan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VLANResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *VLANResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var vlanReq pfsense.VLAN
	resp.Diagnostics.Append(data.Value(ctx, &vlanReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	vlan, err := r.client.UpdateVLAN(ctx, vlanReq)
	if addError(&resp.Diagnostics, "Error updating VLAN", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *vlan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *VLANResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *VLANResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteVLAN(ctx, data.VLANInterface.ValueString())
	if addError(&resp.Diagnostics, "Error deleting VLAN", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *VLANResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("vlan_interface"), req, resp)
}
