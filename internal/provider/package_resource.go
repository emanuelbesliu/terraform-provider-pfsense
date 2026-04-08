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
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*PackageResource)(nil)
	_ resource.ResourceWithConfigure   = (*PackageResource)(nil)
	_ resource.ResourceWithImportState = (*PackageResource)(nil)
)

type PackageResourceModel struct {
	Name             types.String `tfsdk:"name"`
	PackageURL       types.String `tfsdk:"package_url"`
	InstalledVersion types.String `tfsdk:"installed_version"`
	Description      types.String `tfsdk:"description"`
}

func NewPackageResource() resource.Resource { //nolint:ireturn
	return &PackageResource{}
}

type PackageResource struct {
	client *pfsense.Client
}

func (r *PackageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_package", req.ProviderTypeName)
}

func (r *PackageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages pfSense packages. Installs a package from the official repository by name, or from a URL for third-party packages. Removing the resource uninstalls the package.",
		MarkdownDescription: "Manages pfSense packages. Installs a package from the official repository by name, or from a URL for third-party packages. Removing the resource uninstalls the package.",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: PackageModel{}.descriptions()["name"].Description,
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"package_url": schema.StringAttribute{
				Description:         "URL to a .pkg file for third-party package installation. When set, the package is installed via 'pkg-static add' from this URL instead of the official repository. Changing this forces reinstallation.",
				MarkdownDescription: "URL to a `.pkg` file for third-party package installation. When set, the package is installed via `pkg-static add` from this URL instead of the official repository. Changing this forces reinstallation.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"installed_version": schema.StringAttribute{
				Description: PackageModel{}.descriptions()["installed_version"].Description,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: PackageModel{}.descriptions()["description"].Description,
				Computed:    true,
			},
		},
	}
}

func (r *PackageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *PackageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data PackageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	name := data.Name.ValueString()

	var pkg *pfsense.Package
	var err error

	if !data.PackageURL.IsNull() && data.PackageURL.ValueString() != "" {
		pkg, err = r.client.InstallPackageFromURL(ctx, name, data.PackageURL.ValueString())
	} else {
		pkg, err = r.client.InstallPackage(ctx, name)
	}

	if addError(&resp.Diagnostics, "Error installing package", err) {
		return
	}

	data.Name = types.StringValue(pkg.Name)

	if pkg.InstalledVersion != "" {
		data.InstalledVersion = types.StringValue(pkg.InstalledVersion)
	} else {
		data.InstalledVersion = types.StringNull()
	}

	if pkg.Description != "" {
		data.Description = types.StringValue(pkg.Description)
	} else {
		data.Description = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PackageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data PackageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	pkg, err := r.client.GetPackage(ctx, data.Name.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading package", err) {
		return
	}

	data.Name = types.StringValue(pkg.Name)

	if pkg.InstalledVersion != "" {
		data.InstalledVersion = types.StringValue(pkg.InstalledVersion)
	} else {
		data.InstalledVersion = types.StringNull()
	}

	if pkg.Description != "" {
		data.Description = types.StringValue(pkg.Description)
	} else {
		data.Description = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PackageResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

func (r *PackageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data PackageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePackage(ctx, data.Name.ValueString())
	if addError(&resp.Diagnostics, "Error uninstalling package", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *PackageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
