package provider

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*WakeOnLanResource)(nil)
	_ resource.ResourceWithConfigure   = (*WakeOnLanResource)(nil)
	_ resource.ResourceWithImportState = (*WakeOnLanResource)(nil)
)

type WakeOnLanResourceModel struct {
	WakeOnLanModel
}

func NewWakeOnLanResource() resource.Resource { //nolint:ireturn
	return &WakeOnLanResource{}
}

type WakeOnLanResource struct {
	client *pfsense.Client
}

func (r *WakeOnLanResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_wake_on_lan", req.ProviderTypeName)
}

func (r *WakeOnLanResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a Wake-on-LAN entry on pfSense.",
		MarkdownDescription: "Manages a [Wake-on-LAN](https://docs.netgate.com/pfsense/en/latest/services/wake-on-lan.html) entry on pfSense.",
		Attributes: map[string]schema.Attribute{
			"interface": schema.StringAttribute{
				Description: WakeOnLanModel{}.descriptions()["interface"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"mac": schema.StringAttribute{
				Description: WakeOnLanModel{}.descriptions()["mac"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: WakeOnLanModel{}.descriptions()["description"].Description,
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *WakeOnLanResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *WakeOnLanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *WakeOnLanResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var entryReq pfsense.WakeOnLanEntry
	resp.Diagnostics.Append(data.Value(ctx, &entryReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := r.client.CreateWakeOnLanEntry(ctx, entryReq)
	if addError(&resp.Diagnostics, "Error creating wake on lan entry", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *entry)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WakeOnLanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *WakeOnLanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := r.client.GetWakeOnLanEntry(ctx, data.MAC.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading wake on lan entry", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *entry)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WakeOnLanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *WakeOnLanResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var entryReq pfsense.WakeOnLanEntry
	resp.Diagnostics.Append(data.Value(ctx, &entryReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	entry, err := r.client.UpdateWakeOnLanEntry(ctx, entryReq)
	if addError(&resp.Diagnostics, "Error updating wake on lan entry", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *entry)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WakeOnLanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *WakeOnLanResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteWakeOnLanEntry(ctx, data.MAC.ValueString())
	if addError(&resp.Diagnostics, "Error deleting wake on lan entry", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *WakeOnLanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("mac"), req, resp)
}
