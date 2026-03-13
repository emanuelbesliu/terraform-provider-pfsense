package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*SystemAdvancedAdminResource)(nil)
	_ resource.ResourceWithConfigure   = (*SystemAdvancedAdminResource)(nil)
	_ resource.ResourceWithImportState = (*SystemAdvancedAdminResource)(nil)
)

type SystemAdvancedAdminResourceModel struct {
	SystemAdvancedAdminModel
	Apply types.Bool `tfsdk:"apply"`
}

func NewSystemAdvancedAdminResource() resource.Resource { //nolint:ireturn
	return &SystemAdvancedAdminResource{}
}

type SystemAdvancedAdminResource struct {
	client *pfsense.Client
}

func (r *SystemAdvancedAdminResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_system_advanced_admin", req.ProviderTypeName)
}

func (r *SystemAdvancedAdminResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	descriptions := SystemAdvancedAdminModel{}.descriptions()

	resp.Schema = schema.Schema{
		Description:         "System advanced admin access configuration including webConfigurator, SSH, login protection, serial console, and console settings. This is a singleton resource — only one instance can exist per pfSense installation.",
		MarkdownDescription: "[System advanced admin access](https://docs.netgate.com/pfsense/en/latest/config/advanced/admin.html) configuration including webConfigurator, SSH, login protection, serial console, and console settings. This is a **singleton** resource — only one instance can exist per pfSense installation.",
		Attributes: map[string]schema.Attribute{
			// webConfigurator
			"webgui_protocol": schema.StringAttribute{
				Description:         descriptions["webgui_protocol"].Description,
				MarkdownDescription: descriptions["webgui_protocol"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedAdminWebGUIProto),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedAdmin{}.WebGUIProtoOptions()...),
				},
			},
			"ssl_certificate": schema.StringAttribute{
				Description: descriptions["ssl_certificate"].Description,
				Optional:    true,
			},
			"webgui_port": schema.Int64Attribute{
				Description: descriptions["webgui_port"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultAdvancedAdminWebGUIPort)),
				Validators: []validator.Int64{
					int64validator.Between(0, 65535),
				},
			},
			"max_processes": schema.Int64Attribute{
				Description: descriptions["max_processes"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultAdvancedAdminMaxProcs)),
				Validators: []validator.Int64{
					int64validator.Between(1, 500),
				},
			},
			"disable_http_redirect": schema.BoolAttribute{
				Description: descriptions["disable_http_redirect"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disable_hsts": schema.BoolAttribute{
				Description: descriptions["disable_hsts"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ocsp_staple": schema.BoolAttribute{
				Description: descriptions["ocsp_staple"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"login_autocomplete": schema.BoolAttribute{
				Description: descriptions["login_autocomplete"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"quiet_login": schema.BoolAttribute{
				Description: descriptions["quiet_login"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"roaming": schema.BoolAttribute{
				Description: descriptions["roaming"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(true),
			},
			"disable_anti_lockout": schema.BoolAttribute{
				Description: descriptions["disable_anti_lockout"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disable_dns_rebind_check": schema.BoolAttribute{
				Description: descriptions["disable_dns_rebind_check"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"disable_http_referer_check": schema.BoolAttribute{
				Description: descriptions["disable_http_referer_check"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"alternate_hostnames": schema.StringAttribute{
				Description: descriptions["alternate_hostnames"].Description,
				Optional:    true,
			},
			"page_name_first": schema.BoolAttribute{
				Description: descriptions["page_name_first"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},

			// SSH
			"ssh_enabled": schema.BoolAttribute{
				Description: descriptions["ssh_enabled"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"sshd_key_only": schema.StringAttribute{
				Description:         descriptions["sshd_key_only"].Description,
				MarkdownDescription: descriptions["sshd_key_only"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedAdminSSHdKeyOnly),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedAdmin{}.SSHdKeyOnlyOptions()...),
				},
			},
			"sshd_agent_forwarding": schema.BoolAttribute{
				Description: descriptions["sshd_agent_forwarding"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"ssh_port": schema.Int64Attribute{
				Description: descriptions["ssh_port"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultAdvancedAdminSSHPort)),
				Validators: []validator.Int64{
					int64validator.Between(0, 65535),
				},
			},

			// Login Protection
			"login_protection_threshold": schema.Int64Attribute{
				Description: descriptions["login_protection_threshold"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultAdvancedAdminSshguardThreshold)),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"login_protection_blocktime": schema.Int64Attribute{
				Description: descriptions["login_protection_blocktime"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultAdvancedAdminSshguardBlocktime)),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"login_protection_detection_time": schema.Int64Attribute{
				Description: descriptions["login_protection_detection_time"].Description,
				Computed:    true,
				Optional:    true,
				Default:     int64default.StaticInt64(int64(pfsense.DefaultAdvancedAdminSshguardDetectionTime)),
				Validators: []validator.Int64{
					int64validator.AtLeast(0),
				},
			},
			"login_protection_pass_list": schema.StringAttribute{
				Description: descriptions["login_protection_pass_list"].Description,
				Optional:    true,
			},

			// Serial
			"serial_terminal": schema.BoolAttribute{
				Description: descriptions["serial_terminal"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},
			"serial_speed": schema.Int64Attribute{
				Description:         descriptions["serial_speed"].Description,
				MarkdownDescription: descriptions["serial_speed"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             int64default.StaticInt64(int64(pfsense.DefaultAdvancedAdminSerialSpeed)),
				Validators: []validator.Int64{
					int64validator.OneOf(115200, 57600, 38400, 19200, 14400, 9600),
				},
			},
			"primary_console": schema.StringAttribute{
				Description:         descriptions["primary_console"].Description,
				MarkdownDescription: descriptions["primary_console"].MarkdownDescription,
				Computed:            true,
				Optional:            true,
				Default:             stringdefault.StaticString(pfsense.DefaultAdvancedAdminPrimaryConsole),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.AdvancedAdmin{}.PrimaryConsoleOptions()...),
				},
			},

			// Console
			"disable_console_menu": schema.BoolAttribute{
				Description: descriptions["disable_console_menu"].Description,
				Computed:    true,
				Optional:    true,
				Default:     booldefault.StaticBool(false),
			},

			// Apply
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

func (r *SystemAdvancedAdminResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *SystemAdvancedAdminResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SystemAdvancedAdminResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var aReq pfsense.AdvancedAdmin
	resp.Diagnostics.Append(data.Value(ctx, &aReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.UpdateAdvancedAdmin(ctx, aReq)
	if addError(&resp.Diagnostics, "Error creating system advanced admin settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedAdminChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced admin changes", err)
	}
}

func (r *SystemAdvancedAdminResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SystemAdvancedAdminResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.GetAdvancedAdmin(ctx)
	if addError(&resp.Diagnostics, "Error reading system advanced admin settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SystemAdvancedAdminResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *SystemAdvancedAdminResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var aReq pfsense.AdvancedAdmin
	resp.Diagnostics.Append(data.Value(ctx, &aReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	a, err := r.client.UpdateAdvancedAdmin(ctx, aReq)
	if addError(&resp.Diagnostics, "Error updating system advanced admin settings", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedAdminChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced admin changes", err)
	}
}

func (r *SystemAdvancedAdminResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SystemAdvancedAdminResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Reset to pfSense defaults
	defaultAdmin := pfsense.AdvancedAdmin{
		WebGUIProto:    pfsense.DefaultAdvancedAdminWebGUIProto,
		MaxProcs:       pfsense.DefaultAdvancedAdminMaxProcs,
		Roaming:        true,
		SSHdKeyOnly:    pfsense.DefaultAdvancedAdminSSHdKeyOnly,
		SerialSpeed:    pfsense.DefaultAdvancedAdminSerialSpeed,
		PrimaryConsole: pfsense.DefaultAdvancedAdminPrimaryConsole,
	}

	_, err := r.client.UpdateAdvancedAdmin(ctx, defaultAdmin)
	if addError(&resp.Diagnostics, "Error resetting system advanced admin settings to defaults", err) {
		return
	}

	resp.State.RemoveResource(ctx)

	if data.Apply.ValueBool() {
		err = r.client.ApplyAdvancedAdminChanges(ctx)
		addWarning(&resp.Diagnostics, "Error applying system advanced admin changes", err)
	}
}

func (r *SystemAdvancedAdminResource) ImportState(ctx context.Context, _ resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	a, err := r.client.GetAdvancedAdmin(ctx)
	if addError(&resp.Diagnostics, "Error importing system advanced admin settings", err) {
		return
	}

	var data SystemAdvancedAdminResourceModel
	resp.Diagnostics.Append(data.Set(ctx, *a)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.Apply = types.BoolValue(defaultApply)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
