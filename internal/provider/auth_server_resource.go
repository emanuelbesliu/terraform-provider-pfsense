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
	_ resource.Resource                = (*AuthServerResource)(nil)
	_ resource.ResourceWithConfigure   = (*AuthServerResource)(nil)
	_ resource.ResourceWithImportState = (*AuthServerResource)(nil)
)

type AuthServerResourceModel struct {
	AuthServerModel
}

func NewAuthServerResource() resource.Resource { //nolint:ireturn
	return &AuthServerResource{}
}

type AuthServerResource struct {
	client *pfsense.Client
}

func (r *AuthServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_auth_server", req.ProviderTypeName)
}

func (r *AuthServerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	descriptions := AuthServerModel{}.descriptions()
	resp.Schema = schema.Schema{
		Description:         "Manages a pfSense authentication server (LDAP or RADIUS). Auth servers are used to authenticate users against external directory services.",
		MarkdownDescription: "Manages a pfSense [authentication server](https://docs.netgate.com/pfsense/en/latest/usermanager/authservers.html) (LDAP or RADIUS). Auth servers are used to authenticate users against external directory services.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: descriptions["name"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, pfsense.AuthServerMaxNameLength),
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[a-zA-Z0-9.\-_ ]+$`),
						"must contain only alphanumeric characters, dots, hyphens, underscores, and spaces",
					),
				},
			},
			"type": schema.StringAttribute{
				Description: descriptions["type"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("ldap", "radius"),
				},
			},
			"host": schema.StringAttribute{
				Description: descriptions["host"].Description,
				Required:    true,
			},
			"refid": schema.StringAttribute{
				Description: descriptions["refid"].Description,
				Computed:    true,
			},
			// LDAP fields.
			"ldap_port": schema.StringAttribute{
				Description: descriptions["ldap_port"].Description,
				Optional:    true,
				Computed:    true,
			},
			"ldap_urltype": schema.StringAttribute{
				Description: descriptions["ldap_urltype"].Description,
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("Standard TCP", "STARTTLS Encrypted", "SSL/TLS Encrypted"),
				},
			},
			"ldap_protver": schema.StringAttribute{
				Description: descriptions["ldap_protver"].Description,
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("2", "3"),
				},
			},
			"ldap_scope": schema.StringAttribute{
				Description: descriptions["ldap_scope"].Description,
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("one", "subtree"),
				},
			},
			"ldap_basedn": schema.StringAttribute{
				Description: descriptions["ldap_basedn"].Description,
				Optional:    true,
			},
			"ldap_authcn": schema.StringAttribute{
				Description: descriptions["ldap_authcn"].Description,
				Optional:    true,
			},
			"ldap_binddn": schema.StringAttribute{
				Description: descriptions["ldap_binddn"].Description,
				Optional:    true,
			},
			"ldap_bindpw": schema.StringAttribute{
				Description: descriptions["ldap_bindpw"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"ldap_caref": schema.StringAttribute{
				Description: descriptions["ldap_caref"].Description,
				Optional:    true,
			},
			"ldap_timeout": schema.StringAttribute{
				Description: descriptions["ldap_timeout"].Description,
				Optional:    true,
				Computed:    true,
			},
			"ldap_extended_enabled": schema.BoolAttribute{
				Description: descriptions["ldap_extended_enabled"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ldap_extended_query": schema.StringAttribute{
				Description: descriptions["ldap_extended_query"].Description,
				Optional:    true,
			},
			"ldap_attr_user": schema.StringAttribute{
				Description: descriptions["ldap_attr_user"].Description,
				Optional:    true,
				Computed:    true,
			},
			"ldap_attr_group": schema.StringAttribute{
				Description: descriptions["ldap_attr_group"].Description,
				Optional:    true,
				Computed:    true,
			},
			"ldap_attr_member": schema.StringAttribute{
				Description: descriptions["ldap_attr_member"].Description,
				Optional:    true,
				Computed:    true,
			},
			"ldap_attr_groupobj": schema.StringAttribute{
				Description: descriptions["ldap_attr_groupobj"].Description,
				Optional:    true,
				Computed:    true,
			},
			"ldap_pam_groupdn": schema.StringAttribute{
				Description: descriptions["ldap_pam_groupdn"].Description,
				Optional:    true,
			},
			"ldap_utf8": schema.BoolAttribute{
				Description: descriptions["ldap_utf8"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ldap_nostrip_at": schema.BoolAttribute{
				Description: descriptions["ldap_nostrip_at"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ldap_allow_unauthenticated": schema.BoolAttribute{
				Description: descriptions["ldap_allow_unauthenticated"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ldap_rfc2307": schema.BoolAttribute{
				Description: descriptions["ldap_rfc2307"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ldap_rfc2307_userdn": schema.BoolAttribute{
				Description: descriptions["ldap_rfc2307_userdn"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			// RADIUS fields.
			"radius_protocol": schema.StringAttribute{
				Description: descriptions["radius_protocol"].Description,
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("PAP", "CHAP_MD5", "MSCHAPv1", "MSCHAPv2"),
				},
			},
			"radius_auth_port": schema.StringAttribute{
				Description: descriptions["radius_auth_port"].Description,
				Optional:    true,
				Computed:    true,
			},
			"radius_acct_port": schema.StringAttribute{
				Description: descriptions["radius_acct_port"].Description,
				Optional:    true,
				Computed:    true,
			},
			"radius_secret": schema.StringAttribute{
				Description: descriptions["radius_secret"].Description,
				Optional:    true,
				Sensitive:   true,
			},
			"radius_timeout": schema.StringAttribute{
				Description: descriptions["radius_timeout"].Description,
				Optional:    true,
				Computed:    true,
			},
			"radius_nasip_attribute": schema.StringAttribute{
				Description: descriptions["radius_nasip_attribute"].Description,
				Optional:    true,
			},
			"radius_srvcs": schema.StringAttribute{
				Description: descriptions["radius_srvcs"].Description,
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("both", "auth", "acct"),
				},
			},
		},
	}
}

func (r *AuthServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *AuthServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *AuthServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sensitive fields from plan (write-only).
	planBindPW := data.LDAPBindPW
	planRadiusSecret := data.RadiusSecret

	var serverReq pfsense.AuthServer
	resp.Diagnostics.Append(data.Value(ctx, &serverReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.client.CreateAuthServer(ctx, serverReq)
	if addError(&resp.Diagnostics, "Error creating auth server", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *server)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Restore sensitive write-only fields that may not be readable.
	preserveAuthServerSensitiveFields(data, planBindPW, planRadiusSecret)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AuthServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *AuthServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sensitive fields from prior state.
	priorBindPW := data.LDAPBindPW
	priorRadiusSecret := data.RadiusSecret

	server, err := r.client.GetAuthServer(ctx, data.Name.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading auth server", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *server)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Restore sensitive fields from prior state.
	preserveAuthServerSensitiveFields(data, priorBindPW, priorRadiusSecret)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AuthServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *AuthServerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve sensitive fields from plan (write-only).
	planBindPW := data.LDAPBindPW
	planRadiusSecret := data.RadiusSecret

	var serverReq pfsense.AuthServer
	resp.Diagnostics.Append(data.Value(ctx, &serverReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	server, err := r.client.UpdateAuthServer(ctx, serverReq)
	if addError(&resp.Diagnostics, "Error updating auth server", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *server)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Restore sensitive write-only fields.
	preserveAuthServerSensitiveFields(data, planBindPW, planRadiusSecret)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *AuthServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *AuthServerResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteAuthServer(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Error deleting auth server", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *AuthServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

// preserveAuthServerSensitiveFields restores write-only sensitive fields.
// pfSense can read back ldap_bindpw and radius_secret, but we preserve the
// user-provided values from plan/state to maintain consistency with sensitive
// field handling patterns.
func preserveAuthServerSensitiveFields(data *AuthServerResourceModel, priorBindPW, priorRadiusSecret types.String) {
	// Only preserve if the prior value was non-null (user explicitly set it).
	if !priorBindPW.IsNull() && data.LDAPBindPW.IsNull() {
		data.LDAPBindPW = priorBindPW
	}

	if !priorRadiusSecret.IsNull() && data.RadiusSecret.IsNull() {
		data.RadiusSecret = priorRadiusSecret
	}
}
