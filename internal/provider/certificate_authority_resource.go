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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*CertificateAuthorityResource)(nil)
	_ resource.ResourceWithConfigure   = (*CertificateAuthorityResource)(nil)
	_ resource.ResourceWithImportState = (*CertificateAuthorityResource)(nil)
)

type CertificateAuthorityResourceModel struct {
	CertificateAuthorityModel
}

func NewCertificateAuthorityResource() resource.Resource { //nolint:ireturn
	return &CertificateAuthorityResource{}
}

type CertificateAuthorityResource struct {
	client *pfsense.Client
}

func (r *CertificateAuthorityResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_certificate_authority", req.ProviderTypeName)
}

func (r *CertificateAuthorityResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a certificate authority (CA) in the pfSense trust store. Supports importing external CA certificates for TLS validation.",
		MarkdownDescription: "Manages a [certificate authority](https://docs.netgate.com/pfsense/en/latest/certificates/cas.html) (CA) in the pfSense trust store. Supports importing external CA certificates for TLS validation.",
		Attributes: map[string]schema.Attribute{
			"refid": schema.StringAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["refid"].Description,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"descr": schema.StringAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["descr"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"certificate": schema.StringAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["certificate"].Description,
				Required:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"private_key": schema.StringAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["private_key"].Description,
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"trust": schema.BoolAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["trust"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"random_serial": schema.BoolAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["random_serial"].Description,
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"next_serial": schema.Int64Attribute{
				Description: CertificateAuthorityModel{}.descriptions()["next_serial"].Description,
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"subject": schema.StringAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["subject"].Description,
				Computed:    true,
			},
			"issuer": schema.StringAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["issuer"].Description,
				Computed:    true,
			},
			"serial": schema.StringAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["serial"].Description,
				Computed:    true,
			},
			"has_private_key": schema.BoolAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["has_private_key"].Description,
				Computed:    true,
			},
			"is_self_signed": schema.BoolAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["is_self_signed"].Description,
				Computed:    true,
			},
			"valid_from": schema.StringAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["valid_from"].Description,
				Computed:    true,
			},
			"valid_to": schema.StringAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["valid_to"].Description,
				Computed:    true,
			},
			"in_use": schema.BoolAttribute{
				Description: CertificateAuthorityModel{}.descriptions()["in_use"].Description,
				Computed:    true,
			},
		},
	}
}

func (r *CertificateAuthorityResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *CertificateAuthorityResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CertificateAuthorityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var caReq pfsense.CertificateAuthority
	resp.Diagnostics.Append(data.Value(ctx, &caReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ca, err := r.client.ImportCertificateAuthority(ctx, caReq)
	if addError(&resp.Diagnostics, "Error importing certificate authority", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *ca)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve the user-provided certificate and private_key values since the
	// API read-back may not return the private key and the certificate PEM
	// whitespace may differ.
	if !req.Plan.Raw.IsNull() {
		var planData CertificateAuthorityResourceModel
		resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)

		if !resp.Diagnostics.HasError() {
			data.Certificate = planData.Certificate

			if !planData.PrivateKey.IsNull() {
				data.PrivateKey = planData.PrivateKey
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateAuthorityResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CertificateAuthorityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve user-provided sensitive values before overwriting with API data.
	priorCertificate := data.Certificate
	priorPrivateKey := data.PrivateKey

	ca, err := r.client.GetCertificateAuthority(ctx, data.RefID.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading certificate authority", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *ca)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Restore the original user-provided certificate and private_key to
	// avoid spurious diffs caused by PEM whitespace differences or the API
	// not returning the private key.
	if !priorCertificate.IsNull() && !priorCertificate.IsUnknown() {
		data.Certificate = priorCertificate
	}

	if !priorPrivateKey.IsNull() && !priorPrivateKey.IsUnknown() {
		data.PrivateKey = priorPrivateKey
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateAuthorityResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CertificateAuthorityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var stateData CertificateAuthorityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var caReq pfsense.CertificateAuthority
	resp.Diagnostics.Append(data.Value(ctx, &caReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	ca, err := r.client.UpdateCertificateAuthority(ctx, stateData.RefID.ValueString(), caReq)
	if addError(&resp.Diagnostics, "Error updating certificate authority", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *ca)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Preserve user-provided sensitive values.
	var planData CertificateAuthorityResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)

	if !resp.Diagnostics.HasError() {
		data.Certificate = planData.Certificate

		if !planData.PrivateKey.IsNull() {
			data.PrivateKey = planData.PrivateKey
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateAuthorityResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CertificateAuthorityResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCertificateAuthority(ctx, data.RefID.ValueString())
	if addError(&resp.Diagnostics, "Error deleting certificate authority", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *CertificateAuthorityResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Import by refid — the passthrough puts the import ID into the refid attribute,
	// then Read fetches the full state from pfSense.
	resource.ImportStatePassthroughID(ctx, path.Root("refid"), req, resp)

	// Set sensitive fields to null so Read doesn't fail trying to preserve them.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("certificate"), types.StringNull())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("private_key"), types.StringNull())...)
}
