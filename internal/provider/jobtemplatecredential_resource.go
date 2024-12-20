// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &JobTemplateCredentialResource{}
var _ resource.ResourceWithImportState = &JobTemplateCredentialResource{}

func NewJobTemplateCredentialResource() resource.Resource {
	return &JobTemplateCredentialResource{}
}

// JobTemplateCredentialResource defines the resource implementation.
type JobTemplateCredentialResource struct {
	client *AwxClient
}

// JobTemplateCredentialResourceModel describes the resource data model.
type JobTemplateCredentialResourceModel struct {
	Id           types.Int32 `tfsdk:"id"`
	CredentialId types.Int32 `tfsdk:"credential_id"`
}

type JTCredentialAPIRead struct {
	Results []Result `json:"results"`
}

type Result struct {
	Id int `json:"id"`
}

func (r *JobTemplateCredentialResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jobtemplate_credential"
}

func (r *JobTemplateCredentialResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		//TODO fix description on schema and markdown descr
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int32Attribute{
				Required:            true,
				Description:         "The ID of the containing Job Template.",
				MarkdownDescription: "The ID of the containing Job Template",
			},
			"credential_id": schema.Int32Attribute{
				Required:            true,
				Description:         "The ID of the credential to be attached to the job template.",
				MarkdownDescription: "The ID of the credential to be attached to the job template.",
			},
		},
	}
}

func (r *JobTemplateCredentialResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	configureData := req.ProviderData.(*AwxClient)

	r.client = configureData
}

func (r *JobTemplateCredentialResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data JobTemplateCredentialResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// // set url for create HTTP request
	// id, err := strconv.Atoi(data.Id.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable convert id from string to int",
	// 		fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
	// }

	// url := r.client.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/survey_spec", id)

	// // get body data for HTTP request
	// var bodyData JobTemplateSurvey
	// bodyData.Name = data.Name.ValueString()
	// bodyData.Description = data.Description.ValueString()

	// var specs []SurveySpec
	// for _, spec := range data.Spec {

	// 	// convert choices to slice of strings
	// 	stringSlice := make([]string, 0, len(spec.Choices.Elements()))
	// 	diag := spec.Choices.ElementsAs(ctx, &stringSlice, true)
	// 	resp.Diagnostics.Append(diag...)

	// 	if resp.Diagnostics.HasError() {
	// 		return
	// 	}

	// 	// convert to interface{} type
	// 	var finalList interface{} = stringSlice

	// 	specs = append(specs, SurveySpec{
	// 		Type:                spec.Type.ValueString(),
	// 		QuestionName:        spec.QuestionName.ValueString(),
	// 		QuestionDescription: spec.QuestionDescription.ValueString(),
	// 		Variable:            spec.Variable.ValueString(),
	// 		Required:            spec.Required.ValueBool(),
	// 		Max:                 int(spec.Max.ValueInt32()),
	// 		Min:                 int(spec.Min.ValueInt32()),
	// 		Choices:             finalList,
	// 		Default:             spec.Default,
	// 	})
	// }

	// bodyData.Spec = specs

	// jsonData, err := json.Marshal(bodyData)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable marshal json",
	// 		fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
	// }

	// // create HTTP request
	// httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to generate request",
	// 		fmt.Sprintf("Unable to gen url: %v. ", url))
	// }

	// httpReq.Header.Add("Content-Type", "application/json")
	// httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	// httpResp, err := r.client.client.Do(httpReq)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	// }
	// if httpResp.StatusCode != 200 {
	// 	resp.Diagnostics.AddError(
	// 		"Bad request status code.",
	// 		fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))

	// }

	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobTemplateCredentialResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data JobTemplateCredentialResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// set url for create HTTP request
	id := data.Id.ValueInt32()

	url := r.client.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/credentials/", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))

	}

	var responseData JTCredentialAPIRead

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Uanble to get all data out of the http response data body",
			fmt.Sprintf("Body got %v. ", body))
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Uanble unmarshall response body into object",
			fmt.Sprintf("Error =  %v. ", err.Error()))
	}
	//TODO START BELOW

	data.Name = types.StringValue(responseData.Name)
	data.Description = types.StringValue(responseData.Description)

	var dataSpecs []SurveySpecModel
	for _, item := range responseData.Spec {
		specModel := SurveySpecModel{}
		specModel.Max = types.Int32Value(int32(item.Max))
		specModel.Min = types.Int32Value(int32(item.Min))
		specModel.Type = types.StringValue(item.Type)

		itemChoiceKind := reflect.TypeOf(item.Choices).Kind()

		if itemChoiceKind == reflect.Slice {

			elements := make([]string, 0, len(item.Choices.([]any)))

			for _, v := range item.Choices.([]any) {
				elements = append(elements, v.(string))
			}

			listValue, diags := types.ListValueFrom(ctx, types.StringType, elements)
			if diags.HasError() {
				return
			}

			specModel.Choices = listValue
		} else {
			specModel.Choices = types.ListNull(types.StringType)
		}

		itemDefaultKind := reflect.TypeOf(item.Default).Kind()
		switch itemDefaultKind {
		case reflect.Float64:
			specModel.Default = types.StringValue(fmt.Sprint(item.Default.(float64)))
		default:
			specModel.Default = types.StringValue(item.Default.(string))
		}

		specModel.Required = types.BoolValue(item.Required)
		specModel.QuestionName = types.StringValue(item.QuestionName)
		specModel.QuestionDescription = types.StringValue(item.QuestionDescription)
		specModel.Variable = types.StringValue(item.Variable)
		dataSpecs = append(dataSpecs, specModel)
	}

	data.Spec = dataSpecs

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

// Left intentinally "blank" (as initialized by clone of template scaffold) as these resources is replace by schema plan modifiers
func (r *JobTemplateCredentialResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data JobTemplateCredentialResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update example, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JobTemplateCredentialResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data JobTemplateCredentialResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// // set url for create HTTP request
	// id, err := strconv.Atoi(data.Id.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable convert id from string to int",
	// 		fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
	// }

	// url := r.client.endpoint + fmt.Sprintf("/api/v2/job_templates/%d/survey_spec", id)

	// // create HTTP request
	// httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Unable to generate delete request",
	// 		fmt.Sprintf("Unable to gen url: %v. ", url))
	// }

	// httpReq.Header.Add("Content-Type", "application/json")
	// httpReq.Header.Add("Authorization", "Bearer"+" "+r.client.token)

	// httpResp, err := r.client.client.Do(httpReq)
	// if err != nil {
	// 	resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete got error: %s", err))
	// }
	// if httpResp.StatusCode != 200 {
	// 	resp.Diagnostics.AddError(
	// 		"Bad request status code.",
	// 		fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))

	// }
}

func (r *JobTemplateCredentialResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
