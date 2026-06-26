package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	_ resource.Resource                = (*InterfaceBridgeResource)(nil)
	_ resource.ResourceWithConfigure   = (*InterfaceBridgeResource)(nil)
	_ resource.ResourceWithImportState = (*InterfaceBridgeResource)(nil)
)

type InterfaceBridgeResourceModel struct {
	InterfaceBridgeModel
}

func NewInterfaceBridgeResource() resource.Resource { //nolint:ireturn
	return &InterfaceBridgeResource{}
}

type InterfaceBridgeResource struct {
	client *pfsense.Client
}

func bridgeMemberSubsetListAttribute(description string) schema.ListAttribute {
	return schema.ListAttribute{
		Description: description,
		ElementType: types.StringType,
		Optional:    true,
		Validators: []validator.List{
			listvalidator.UniqueValues(),
			listvalidator.ValueStringsAre(stringIsInterface()),
		},
	}
}

func (r *InterfaceBridgeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_interface_bridge", req.ProviderTypeName)
}

func (r *InterfaceBridgeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	descriptions := InterfaceBridgeModel{}.descriptions()

	resp.Schema = schema.Schema{
		Description:         "Bridge interface joining two or more interfaces together at layer 2.",
		MarkdownDescription: "[Bridge interface](https://docs.netgate.com/pfsense/en/latest/interfaces/bridges.html) joining two or more interfaces together at layer 2.",
		Attributes: map[string]schema.Attribute{
			"bridge_if": schema.StringAttribute{
				Description: descriptions["bridge_if"].Description,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"members": schema.ListAttribute{
				Description: descriptions["members"].Description,
				ElementType: types.StringType,
				Required:    true,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.UniqueValues(),
					listvalidator.ValueStringsAre(stringIsInterface()),
				},
			},
			"description": schema.StringAttribute{
				Description: descriptions["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 255),
				},
			},
			"enable_stp": schema.BoolAttribute{
				Description: descriptions["enable_stp"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ip6_link_local": schema.BoolAttribute{
				Description: descriptions["ip6_link_local"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"protocol": schema.StringAttribute{
				Description: descriptions["protocol"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.ValidBridgeProtocols...),
				},
			},
			"priority": schema.Int64Attribute{
				Description: descriptions["priority"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 61440),
				},
			},
			"hello_time": schema.Int64Attribute{
				Description: descriptions["hello_time"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 2),
				},
			},
			"forward_delay": schema.Int64Attribute{
				Description: descriptions["forward_delay"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(4, 30),
				},
			},
			"max_age": schema.Int64Attribute{
				Description: descriptions["max_age"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(6, 40),
				},
			},
			"hold_count": schema.Int64Attribute{
				Description: descriptions["hold_count"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(1, 10),
				},
			},
			"max_addresses": schema.Int64Attribute{
				Description: descriptions["max_addresses"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"cache_expire": schema.Int64Attribute{
				Description: descriptions["cache_expire"].Description,
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.Between(0, 3600),
				},
			},
			"stp_interfaces":       bridgeMemberSubsetListAttribute(descriptions["stp_interfaces"].Description),
			"static_interfaces":    bridgeMemberSubsetListAttribute(descriptions["static_interfaces"].Description),
			"private_interfaces":   bridgeMemberSubsetListAttribute(descriptions["private_interfaces"].Description),
			"span_interfaces":      bridgeMemberSubsetListAttribute(descriptions["span_interfaces"].Description),
			"edge_interfaces":      bridgeMemberSubsetListAttribute(descriptions["edge_interfaces"].Description),
			"auto_edge_interfaces": bridgeMemberSubsetListAttribute(descriptions["auto_edge_interfaces"].Description),
			"ptp_interfaces":       bridgeMemberSubsetListAttribute(descriptions["ptp_interfaces"].Description),
			"auto_ptp_interfaces":  bridgeMemberSubsetListAttribute(descriptions["auto_ptp_interfaces"].Description),
		},
	}
}

func (r *InterfaceBridgeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *InterfaceBridgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *InterfaceBridgeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var bridgeReq pfsense.Bridge
	resp.Diagnostics.Append(data.Value(ctx, &bridgeReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bridge, err := r.client.CreateBridge(ctx, bridgeReq)
	if addError(&resp.Diagnostics, "Error creating bridge", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *bridge)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceBridgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *InterfaceBridgeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bridge, err := r.client.GetBridge(ctx, data.BridgeIf.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading bridge", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *bridge)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceBridgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *InterfaceBridgeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var state *InterfaceBridgeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var bridgeReq pfsense.Bridge
	resp.Diagnostics.Append(data.Value(ctx, &bridgeReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	bridge, err := r.client.UpdateBridge(ctx, state.BridgeIf.ValueString(), bridgeReq)
	if addError(&resp.Diagnostics, "Error updating bridge", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *bridge)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceBridgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *InterfaceBridgeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteBridge(ctx, data.BridgeIf.ValueString())
	if addError(&resp.Diagnostics, "Error deleting bridge", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *InterfaceBridgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("bridge_if"), req, resp)
}
