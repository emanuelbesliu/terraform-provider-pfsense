package provider

import (
	"context"
	"errors"
	"fmt"
	"regexp"

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
	_ resource.Resource                = (*GroupResource)(nil)
	_ resource.ResourceWithConfigure   = (*GroupResource)(nil)
	_ resource.ResourceWithImportState = (*GroupResource)(nil)
)

type GroupResourceModel struct {
	GroupModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewGroupResource() resource.Resource { //nolint:ireturn
	return &GroupResource{}
}

type GroupResource struct {
	client *pfsense.Client
}

func (r *GroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_group", req.ProviderTypeName)
}

func (r *GroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	descriptions := GroupModel{}.descriptions()
	resp.Schema = schema.Schema{
		Description:         "Manages a pfSense local group. Groups can be assigned privileges and user memberships for controlling access to the web GUI and system services.",
		MarkdownDescription: "Manages a pfSense local [group](https://docs.netgate.com/pfsense/en/latest/usermanager/groups.html). Groups can be assigned privileges and user memberships for controlling access to the web GUI and system services.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: descriptions["name"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, pfsense.GroupMaxNameLength),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9.\-_]+$`),
						"must contain only alphanumeric characters, dots, hyphens, and underscores",
					),
				},
			},
			"description": schema.StringAttribute{
				Description: descriptions["description"].Description,
				Optional:    true,
			},
			"scope": schema.StringAttribute{
				Description: descriptions["scope"].Description,
				Computed:    true,
			},
			"gid": schema.StringAttribute{
				Description: descriptions["gid"].Description,
				Computed:    true,
			},
			"members": schema.ListAttribute{
				Description: descriptions["members"].Description,
				Optional:    true,
				ElementType: types.StringType,
			},
			"privileges": schema.ListAttribute{
				Description: descriptions["privileges"].Description,
				Optional:    true,
				ElementType: types.StringType,
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

func (r *GroupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *GroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the plan's null/empty distinction for list attributes.
	planMembers := data.Members
	planPrivileges := data.Privileges

	var groupReq pfsense.Group
	resp.Diagnostics.Append(data.Value(ctx, &groupReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.CreateGroup(ctx, groupReq)
	if addError(&resp.Diagnostics, "Error creating group", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *group)...)

	if resp.Diagnostics.HasError() {
		return
	}

	preserveGroupListSemantics(data, planMembers, planPrivileges)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyGroupChanges(ctx, group.Name)
		addWarning(&resp.Diagnostics, "Error applying group changes", err)
	}
}

func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *GroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the prior state's null/empty distinction for list attributes.
	priorMembers := data.Members
	priorPrivileges := data.Privileges

	group, err := r.client.GetGroup(ctx, data.Name.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading group", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *group)...)

	if resp.Diagnostics.HasError() {
		return
	}

	preserveGroupListSemantics(data, priorMembers, priorPrivileges)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *GroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the plan's null/empty distinction for list attributes.
	planMembers := data.Members
	planPrivileges := data.Privileges

	var groupReq pfsense.Group
	resp.Diagnostics.Append(data.Value(ctx, &groupReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.client.UpdateGroup(ctx, groupReq)
	if addError(&resp.Diagnostics, "Error updating group", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *group)...)

	if resp.Diagnostics.HasError() {
		return
	}

	preserveGroupListSemantics(data, planMembers, planPrivileges)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyGroupChanges(ctx, group.Name)
		addWarning(&resp.Diagnostics, "Error applying group changes", err)
	}
}

func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *GroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteGroup(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Error deleting group", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *GroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(defaultApply))...)
}

// preserveGroupListSemantics ensures that if the user configured an empty list
// (e.g. members = []), we return an empty list rather than null. This prevents
// the "Provider produced inconsistent result" error when Terraform planned with
// an empty list but the read-back from pfSense returns no items (which Set()
// converts to null).
func preserveGroupListSemantics(data *GroupResourceModel, priorMembers, priorPrivileges types.List) {
	if !priorMembers.IsNull() && data.Members.IsNull() {
		data.Members, _ = types.ListValue(types.StringType, []attr.Value{})
	}

	if !priorPrivileges.IsNull() && data.Privileges.IsNull() {
		data.Privileges, _ = types.ListValue(types.StringType, []attr.Value{})
	}
}
