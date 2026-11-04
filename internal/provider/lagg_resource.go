package provider

import (
	"context"
	"errors"

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
	_ resource.Resource                = &LaggResource{}
	_ resource.ResourceWithConfigure   = &LaggResource{}
	_ resource.ResourceWithImportState = &LaggResource{}
)

type LaggResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Members     types.List   `tfsdk:"members"`
	Protocol    types.String `tfsdk:"protocol"`
	Description types.String `tfsdk:"description"`
	Apply       types.Bool   `tfsdk:"apply"`
}

func NewLaggResource() resource.Resource { //nolint:ireturn
	return &LaggResource{}
}

type LaggResource struct {
	client *pfsense.Client
}

func (r *LaggResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lagg"
}

func (r *LaggResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "LAGG (link aggregation group). The name is assigned server-side sequentially (lagg0, lagg1, ...).",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Assigned LAGG name (laggN).",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"members": schema.ListAttribute{
				Description: "Member interfaces (physical NIC ids as pfSense knows them, e.g. em0, igc1).",
				Required:    true,
				ElementType: types.StringType,
			},
			"protocol": schema.StringAttribute{
				Description: "Aggregation protocol: none, failover, lacp, loadbalance, protect, prefixlacp.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("lacp"),
				Validators: []validator.String{
					stringvalidator.OneOf("none", "failover", "lacp", "loadbalance", "protect", "prefixlacp"),
				},
			},
			"description": schema.StringAttribute{
				Description: "For administrative reference (not parsed).",
				Optional:    true,
			},
			"apply": schema.BoolAttribute{
				Description: "Run interface_lagg_configure immediately after writing. FALSE by default: config-plane only, so a remote change can never cut the management link by surprise.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
		},
	}
}

func (r *LaggResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}
	r.client = client
}

func (r *LaggResource) toModel(l pfsense.Lagg, apply types.Bool) LaggResourceModel {
	members, diags := types.ListValueFrom(context.Background(), types.StringType, l.Members)
	_ = diags // element conversion errors impossible for []string
	return LaggResourceModel{
		Name:        types.StringValue(l.Name),
		Members:     members,
		Protocol:    types.StringValue(l.Protocol),
		Description: types.StringValue(l.Description),
		Apply:       apply,
	}
}

func (r *LaggResource) fromModel(ctx context.Context, data LaggResourceModel) (pfsense.Lagg, bool) {
	var lagg pfsense.Lagg
	var members []string
	if diags := data.Members.ElementsAs(ctx, &members, false); diags.HasError() {
		return lagg, false
	}
	if err := lagg.SetMembers(members); err != nil {
		return lagg, false
	}
	if err := lagg.SetProtocol(data.Protocol.ValueString()); err != nil {
		return lagg, false
	}
	if err := lagg.SetDescription(data.Description.ValueString()); err != nil {
		return lagg, false
	}
	return lagg, true
}

func (r *LaggResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data LaggResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	lagg, ok := r.fromModel(ctx, data)
	if !ok {
		resp.Diagnostics.AddError("Invalid LAGG configuration", "members/protocol validation failed")
		return
	}
	created, err := r.client.CreateLagg(ctx, lagg, data.Apply.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Error creating LAGG", err.Error())
		return
	}
	resp.State.Set(ctx, r.toModel(*created, data.Apply))
}

func (r *LaggResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data LaggResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	created, err := r.client.GetLagg(ctx, data.Name.ValueString())
	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error reading LAGG", err.Error())
		return
	}
	resp.State.Set(ctx, r.toModel(*created, data.Apply))
}

func (r *LaggResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data LaggResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	var state LaggResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	lagg, ok := r.fromModel(ctx, data)
	if !ok {
		resp.Diagnostics.AddError("Invalid LAGG configuration", "members/protocol validation failed")
		return
	}
	lagg.Name = state.Name.ValueString()
	updated, err := r.client.UpdateLagg(ctx, lagg, data.Apply.ValueBool())
	if err != nil {
		resp.Diagnostics.AddError("Error updating LAGG", err.Error())
		return
	}
	resp.State.Set(ctx, r.toModel(*updated, data.Apply))
}

func (r *LaggResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data LaggResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.client.DeleteLagg(ctx, data.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting LAGG", err.Error())
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *LaggResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(false))...)
}
