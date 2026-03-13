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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*InterfaceResource)(nil)
	_ resource.ResourceWithConfigure   = (*InterfaceResource)(nil)
	_ resource.ResourceWithImportState = (*InterfaceResource)(nil)
)

type InterfaceResourceModel struct {
	InterfaceModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewInterfaceResource() resource.Resource { //nolint:ireturn
	return &InterfaceResource{}
}

type InterfaceResource struct {
	client *pfsense.Client
}

func (r *InterfaceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_interface", req.ProviderTypeName)
}

func (r *InterfaceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Interface assignment. Assigns a physical or VLAN interface and configures its network settings.",
		MarkdownDescription: "[Interface assignment](https://docs.netgate.com/pfsense/en/latest/interfaces/index.html). Assigns a physical or VLAN interface and configures its network settings.",
		Attributes: map[string]schema.Attribute{
			"logical_name": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["logical_name"].Description,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"port": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["port"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: InterfaceModel{}.descriptions()["enabled"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"ipv4_type": schema.StringAttribute{
				Description:         InterfaceModel{}.descriptions()["ipv4_type"].Description,
				MarkdownDescription: InterfaceModel{}.descriptions()["ipv4_type"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("none"),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.InterfaceIPv4Types...),
				},
			},
			"ipv4_address": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv4_address"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ipv4_subnet": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv4_subnet"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ipv4_gateway": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv4_gateway"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ipv6_type": schema.StringAttribute{
				Description:         InterfaceModel{}.descriptions()["ipv6_type"].Description,
				MarkdownDescription: InterfaceModel{}.descriptions()["ipv6_type"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString("none"),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.InterfaceIPv6Types...),
				},
			},
			"ipv6_address": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv6_address"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ipv6_subnet": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv6_subnet"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"ipv6_gateway": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["ipv6_gateway"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"spoof_mac": schema.StringAttribute{
				Description: InterfaceModel{}.descriptions()["spoof_mac"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"mtu": schema.Int64Attribute{
				Description: InterfaceModel{}.descriptions()["mtu"].Description,
				Optional:    true,
			},
			"mss": schema.Int64Attribute{
				Description: InterfaceModel{}.descriptions()["mss"].Description,
				Optional:    true,
			},
			"block_private": schema.BoolAttribute{
				Description: InterfaceModel{}.descriptions()["block_private"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"block_bogons": schema.BoolAttribute{
				Description: InterfaceModel{}.descriptions()["block_bogons"].Description,
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

func (r *InterfaceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *InterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *InterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var ifaceReq pfsense.Interface
	resp.Diagnostics.Append(data.Value(ctx, &ifaceReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	iface, err := r.client.CreateInterface(ctx, ifaceReq)
	if addError(&resp.Diagnostics, "Error creating interface", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *iface)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyInterfaceChanges(ctx, iface.LogicalName)
		addWarning(&resp.Diagnostics, "Error applying interface changes", err)
	}
}

func (r *InterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	iface, err := r.client.GetInterface(ctx, data.LogicalName.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading interface", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *iface)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *InterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *InterfaceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to preserve logical_name.
	var state *InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure logical_name from state is used.
	data.LogicalName = state.LogicalName

	var ifaceReq pfsense.Interface
	resp.Diagnostics.Append(data.Value(ctx, &ifaceReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	iface, err := r.client.UpdateInterface(ctx, ifaceReq)
	if addError(&resp.Diagnostics, "Error updating interface", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *iface)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyInterfaceChanges(ctx, iface.LogicalName)
		addWarning(&resp.Diagnostics, "Error applying interface changes", err)
	}
}

func (r *InterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *InterfaceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	logicalName := data.LogicalName.ValueString()

	err := r.client.DeleteInterface(ctx, logicalName)
	if addError(&resp.Diagnostics, "Error deleting interface", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyInterfaceChanges(ctx, logicalName)
		addWarning(&resp.Diagnostics, "Error applying interface changes", err)
	}
}

func (r *InterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("logical_name"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(defaultApply))...)
}
