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
	_ resource.Resource                = (*UserResource)(nil)
	_ resource.ResourceWithConfigure   = (*UserResource)(nil)
	_ resource.ResourceWithImportState = (*UserResource)(nil)
)

type UserResourceModel struct {
	UserModel
	Password types.String `tfsdk:"password"`
	Apply    types.Bool   `tfsdk:"apply"`
}

func NewUserResource() resource.Resource { //nolint:ireturn
	return &UserResource{}
}

type UserResource struct {
	client *pfsense.Client
}

func (r *UserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_user", req.ProviderTypeName)
}

func (r *UserResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	descriptions := UserModel{}.descriptions()
	resp.Schema = schema.Schema{
		Description:         "Manages a pfSense local user account. Users can be assigned privileges and group memberships for controlling access to the web GUI and system services.",
		MarkdownDescription: "Manages a pfSense local [user](https://docs.netgate.com/pfsense/en/latest/usermanager/index.html) account. Users can be assigned privileges and group memberships for controlling access to the web GUI and system services.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: descriptions["name"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, pfsense.UserMaxNameLength),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9.\-_]+$`),
						"must contain only alphanumeric characters, dots, hyphens, and underscores",
					),
				},
			},
			"password": schema.StringAttribute{
				Description: descriptions["password"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"description": schema.StringAttribute{
				Description: descriptions["description"].Description,
				Optional:    true,
			},
			"scope": schema.StringAttribute{
				Description: descriptions["scope"].Description,
				Computed:    true,
			},
			"uid": schema.StringAttribute{
				Description: descriptions["uid"].Description,
				Computed:    true,
			},
			"disabled": schema.BoolAttribute{
				Description: descriptions["disabled"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"expires": schema.StringAttribute{
				Description: descriptions["expires"].Description,
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(0[1-9]|1[0-2])/(0[1-9]|[12]\d|3[01])/\d{4}$`),
						"must be in MM/DD/YYYY format",
					),
				},
			},
			"authorized_keys": schema.StringAttribute{
				Description: descriptions["authorized_keys"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"ipsec_psk": schema.StringAttribute{
				Description: descriptions["ipsec_psk"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"privileges": schema.ListAttribute{
				Description: descriptions["privileges"].Description,
				Optional:    true,
				ElementType: types.StringType,
			},
			"groups": schema.ListAttribute{
				Description: descriptions["groups"].Description,
				Optional:    true,
				ElementType: types.StringType,
			},
			"custom_settings": schema.BoolAttribute{
				Description: descriptions["custom_settings"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"webgui_css": schema.StringAttribute{
				Description: descriptions["webgui_css"].Description,
				Optional:    true,
				Computed:    true,
			},
			"dashboard_columns": schema.StringAttribute{
				Description: descriptions["dashboard_columns"].Description,
				Optional:    true,
				Computed:    true,
			},
			"keep_history": schema.BoolAttribute{
				Description: descriptions["keep_history"].Description,
				Optional:    true,
				Computed:    true,
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

func (r *UserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Password is required for creating a new user.
	if data.Password.IsNull() || data.Password.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Password is required",
			"A password must be provided when creating a new user.",
		)

		return
	}

	var userReq pfsense.User
	resp.Diagnostics.Append(data.Value(ctx, &userReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.client.CreateUser(ctx, userReq, data.Password.ValueString())
	if addError(&resp.Diagnostics, "Error creating user", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *user)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve password in state (write-only, cannot be read back).
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyUserChanges(ctx, user.Name)
		addWarning(&resp.Diagnostics, "Error applying user changes", err)
	}
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the password from state before reading (write-only field).
	previousPassword := data.Password

	user, err := r.client.GetUser(ctx, data.Name.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading user", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *user)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Restore password from state (cannot be read from pfSense).
	data.Password = previousPassword

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var userReq pfsense.User
	resp.Diagnostics.Append(data.Value(ctx, &userReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Password is optional for updates.
	password := ""
	if !data.Password.IsNull() {
		password = data.Password.ValueString()
	}

	user, err := r.client.UpdateUser(ctx, userReq, password)
	if addError(&resp.Diagnostics, "Error updating user", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *user)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve password in state (write-only, cannot be read back).
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyUserChanges(ctx, user.Name)
		addWarning(&resp.Diagnostics, "Error applying user changes", err)
	}
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteUser(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Error deleting user", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		// Apply is a no-op after delete since local_user_del already cleaned up,
		// but we keep the pattern consistent.
		err = r.client.ApplyUserChanges(ctx, data.Name.ValueString())
		// Ignore apply errors after delete since the user no longer exists.
		_ = err
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("apply"), types.BoolValue(defaultApply))...)
}
