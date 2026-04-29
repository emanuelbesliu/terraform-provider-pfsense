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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/marshallford/terraform-provider-pfsense/pkg/pfsense"
)

var (
	_ resource.Resource                = (*CertificateResource)(nil)
	_ resource.ResourceWithConfigure   = (*CertificateResource)(nil)
	_ resource.ResourceWithImportState = (*CertificateResource)(nil)
)

type CertificateResourceModel struct {
	CertificateModel
}

func NewCertificateResource() resource.Resource { //nolint:ireturn
	return &CertificateResource{}
}

type CertificateResource struct {
	client *pfsense.Client
}

func (r *CertificateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_certificate", req.ProviderTypeName)
}

func (r *CertificateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a certificate in the pfSense certificate manager. Supports importing external certificates for use with OpenVPN, IPsec, web GUI, and other services.",
		MarkdownDescription: "Manages a [certificate](https://docs.netgate.com/pfsense/en/latest/certificates/index.html) in the pfSense certificate manager. Supports importing external certificates for use with OpenVPN, IPsec, web GUI, and other services.",
		Attributes: map[string]schema.Attribute{
			"refid": schema.StringAttribute{
				Description: CertificateModel{}.descriptions()["refid"].Description,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"descr": schema.StringAttribute{
				Description: CertificateModel{}.descriptions()["descr"].Description,
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"type": schema.StringAttribute{
				Description:         CertificateModel{}.descriptions()["type"].Description,
				MarkdownDescription: CertificateModel{}.descriptions()["type"].MarkdownDescription,
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("server"),
				Validators: []validator.String{
					stringvalidator.OneOf(pfsense.Certificate{}.CertTypes()...),
				},
			},
			"caref": schema.StringAttribute{
				Description: CertificateModel{}.descriptions()["caref"].Description,
				Optional:    true,
			},
			"certificate": schema.StringAttribute{
				Description: CertificateModel{}.descriptions()["certificate"].Description,
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
				Description: CertificateModel{}.descriptions()["private_key"].Description,
				Optional:    true,
				Sensitive:   true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subject": schema.StringAttribute{
				Description: CertificateModel{}.descriptions()["subject"].Description,
				Computed:    true,
			},
			"issuer": schema.StringAttribute{
				Description: CertificateModel{}.descriptions()["issuer"].Description,
				Computed:    true,
			},
			"serial": schema.StringAttribute{
				Description: CertificateModel{}.descriptions()["serial"].Description,
				Computed:    true,
			},
			"has_private_key": schema.BoolAttribute{
				Description: CertificateModel{}.descriptions()["has_private_key"].Description,
				Computed:    true,
			},
			"is_self_signed": schema.BoolAttribute{
				Description: CertificateModel{}.descriptions()["is_self_signed"].Description,
				Computed:    true,
			},
			"valid_from": schema.StringAttribute{
				Description: CertificateModel{}.descriptions()["valid_from"].Description,
				Computed:    true,
			},
			"valid_to": schema.StringAttribute{
				Description: CertificateModel{}.descriptions()["valid_to"].Description,
				Computed:    true,
			},
			"in_use": schema.BoolAttribute{
				Description: CertificateModel{}.descriptions()["in_use"].Description,
				Computed:    true,
			},
		},
	}
}

func (r *CertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	client, ok := configureResourceClient(req, resp)
	if !ok {
		return
	}

	r.client = client
}

func (r *CertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *CertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var certReq pfsense.Certificate
	resp.Diagnostics.Append(data.Value(ctx, &certReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cert, err := r.client.ImportCertificate(ctx, certReq)
	if addError(&resp.Diagnostics, "Error importing certificate", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *cert)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !req.Plan.Raw.IsNull() {
		var planData CertificateResourceModel
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

func (r *CertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *CertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	priorCertificate := data.Certificate
	priorPrivateKey := data.PrivateKey

	cert, err := r.client.GetCertificate(ctx, data.RefID.ValueString())

	if errors.Is(err, pfsense.ErrNotFound) {
		resp.State.RemoveResource(ctx)

		return
	}

	if addError(&resp.Diagnostics, "Error reading certificate", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *cert)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !priorCertificate.IsNull() && !priorCertificate.IsUnknown() {
		data.Certificate = priorCertificate
	}

	if !priorPrivateKey.IsNull() && !priorPrivateKey.IsUnknown() {
		data.PrivateKey = priorPrivateKey
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var stateData CertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var certReq pfsense.Certificate
	resp.Diagnostics.Append(data.Value(ctx, &certReq)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cert, err := r.client.UpdateCertificate(ctx, stateData.RefID.ValueString(), certReq)
	if addError(&resp.Diagnostics, "Error updating certificate", err) {
		return
	}

	resp.Diagnostics.Append(data.Set(ctx, *cert)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var planData CertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)

	if !resp.Diagnostics.HasError() {
		data.Certificate = planData.Certificate

		if !planData.PrivateKey.IsNull() {
			data.PrivateKey = planData.PrivateKey
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *CertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *CertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCertificate(ctx, data.RefID.ValueString())
	if addError(&resp.Diagnostics, "Error deleting certificate", err) {
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *CertificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("refid"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("certificate"), types.StringNull())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("private_key"), types.StringNull())...)
}
