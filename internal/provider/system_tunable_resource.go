package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"

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
	_ resource.Resource                = (*SystemTunableResource)(nil)
	_ resource.ResourceWithConfigure   = (*SystemTunableResource)(nil)
	_ resource.ResourceWithImportState = (*SystemTunableResource)(nil)
)

type SystemTunableResourceModel struct {
	SystemTunableModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewSystemTunableResource() resource.Resource { //nolint:ireturn
	return &SystemTunableResource{}
}

type SystemTunableResource struct {
	client *pfsense.Client
}

func (r *SystemTunableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_tunable", req.ProviderTypeName)
}

func (r *SystemTunableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "System tunable (sysctl) for kernel parameter configuration. Tunables allow adjusting FreeBSD kernel parameters that control networking, memory, and other system behaviors.",
		MarkdownDescription: "System [tunable](https://docs.netgate.com/pfsense/en/latest/system/advanced-tunables.html) (sysctl) for kernel parameter configuration. Tunables allow adjusting FreeBSD kernel parameters that control networking, memory, and other system behaviors.",
		Attributes: map[string]schema.Attribute{
			"tunable": schema.StringAttribute{
				Description: SystemTunableModel{}.descriptions()["tunable"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9._]+$`),
						"must contain only alphanumeric characters, dots, and underscores",
					),
				},
			},
			"value": schema.StringAttribute{
				Description: SystemTunableModel{}.descriptions()["value"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9.\-_%/]+$`),
						"must contain only alphanumeric characters, dots, hyphens, underscores, percent signs, and forward slashes",
					),
				},
			},
			"description": schema.StringAttribute{
				Description: SystemTunableModel{}.descriptions()["description"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 200),
				},
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

func (r *SystemTunableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *SystemTunableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SystemTunableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var tunableReq pfsense.SystemTunable
	resp.Diagnostics.Append(data.Value(ctx, &tunableReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tunable, err := r.client.CreateSystemTunable(ctx, tunableReq)
	if addError(&resp.Diagnostics, "Error creating system tunable", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *tunable)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplySystemTunableChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system tunable changes", err)
	}
}

func (r *SystemTunableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SystemTunableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tunable, err := r.client.GetSystemTunable(ctx, data.Tunable.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading system tunable", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *tunable)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemTunableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SystemTunableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var tunableReq pfsense.SystemTunable
	resp.Diagnostics.Append(data.Value(ctx, &tunableReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tunable, err := r.client.UpdateSystemTunable(ctx, tunableReq)
	if addError(&resp.Diagnostics, "Error updating system tunable", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *tunable)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplySystemTunableChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system tunable changes", err)
	}
}

func (r *SystemTunableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SystemTunableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteSystemTunable(ctx, data.Tunable.ValueString())
	if addError(&resp.Diagnostics, "Error deleting system tunable", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplySystemTunableChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system tunable changes", err)
	}
}

func (r *SystemTunableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("tunable"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(defaultApply))...)
}
