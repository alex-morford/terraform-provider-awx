package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &WorkflowJobTemplatesResource{}
var _ resource.ResourceWithImportState = &WorkflowJobTemplatesResource{}

func NewWorkflowJobTemplatesResource() resource.Resource {
	return &WorkflowJobTemplatesResource{}
}

type WorkflowJobTemplatesResource struct {
	client *AwxClient
}

type WorkflowJobTemplatesResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Description          types.String `tfsdk:"description"`
	ExtraVars            types.String `tfsdk:"extra_vars"`
	Organization         types.Int32  `tfsdk:"organization"`
	SurveyEnabled        types.Bool   `tfsdk:"survey_enabled"`
	AllowSimultaneous    types.Bool   `tfsdk:"allow_simultaneous"`
	AskVariablesOnLaunch types.Bool   `tfsdk:"ask_variables_on_launch"`
	Inventory            types.Int32  `tfsdk:"inventory"`
	Limit                types.String `tfsdk:"limit"`
	ScmBranch            types.String `tfsdk:"scm_branch"`
	AskInventoryOnLaunch types.Bool   `tfsdk:"ask_inventory_on_launch"`
	AskScmBranchOnLaunch types.Bool   `tfsdk:"ask_scm_branch_on_launch"`
	AskLimitOnLaunch     types.Bool   `tfsdk:"ask_limit_on_launch"`
	WebhookService       types.String `tfsdk:"webhook_service"`
	WebhookCredential    types.String `tfsdk:"webhook_credential"`
	AskLabelsOnLaunch    types.Bool   `tfsdk:"ask_labels_on_launch"`
	AskSkipTagsOnLaunch  types.Bool   `tfsdk:"ask_skip_tags_on_launch"`
	AskTagsOnLaunch      types.Bool   `tfsdk:"ask_tags_on_launch"`
	SkipTags             types.String `tfsdk:"skip_tags"`
	JobTags              types.String `tfsdk:"job_tags"`
}

type WorkflowJobTemplateAPIModel struct {
	Id                   int    `json:"id"`
	Name                 string `json:"name"`
	Description          any    `json:"description,omitempty"`
	ExtraVars            string `json:"extra_vars"`
	Organization         int    `json:"organization"`
	SurveyEnabled        bool   `json:"survey_enabled"`
	AllowSimultaneous    bool   `json:"allow_simultaneous"`
	AskVariablesOnLaunch bool   `json:"ask_variables_on_launch"`
	Inventory            any    `json:"inventory"`
	Limit                string `json:"limit,omitempty"`
	ScmBranch            string `json:"scm_branch"`
	AskInventoryOnLaunch bool   `json:"ask_inventory_on_launch"`
	AskScmBranchOnLaunch bool   `json:"ask_scm_branch_on_launch"`
	AskLimitOnLaunch     bool   `json:"ask_limit_on_launch"`
	WebhookService       string `json:"webhook_service"`
	WebhookCredential    string `json:"webhook_credential"`
	AskLabelsOnLaunch    bool   `json:"ask_labels_on_launch"`
	AskSkipTagsOnLaunch  bool   `json:"ask_skip_tags_on_launch"`
	AskTagsOnLaunch      bool   `json:"ask_tags_on_launch"`
	SkipTags             string `json:"skip_tags"`
	JobTags              string `json:"job_tags"`
}

func (r *WorkflowJobTemplatesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_workflow_job_template"
}

func (r *WorkflowJobTemplatesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "A workflow job template resource.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"extra_vars": schema.StringAttribute{
				Optional: true,
			},
			"organization": schema.Int32Attribute{
				Required: true,
			},
			"survey_enabled": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Defaults to `false`.",
			},
			"allow_simultaneous": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Defaults to `false`.",
			},
			"ask_variables_on_launch": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Defaults to `false`.",
			},
			"inventory": schema.Int32Attribute{
				Optional:    true,
				Description: "Inventory ID of the inventory containing the hosts you want this job to manage.",
			},
			"limit": schema.StringAttribute{
				Optional:    true,
				Description: "Provide a host pattern to further constrain the list of hosts that will be managed or affected by the playbook. Multiple patterns are allowed.",
			},
			"scm_branch": schema.StringAttribute{
				Optional:    true,
				Description: "Select a branch for the workflow. This branch is applied to all job template nodes that prompt for a branch.",
			},
			"ask_inventory_on_launch": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Prompt on launch for Inventory. Defaults to `false`.",
			},
			"ask_scm_branch_on_launch": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Prompt on launch for Source Control Branch. Defaults to `false`.",
			},
			"ask_limit_on_launch": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Prompt on launch for Limit. Defaults to `false`.",
			},
			"webhook_service": schema.StringAttribute{
				Optional:    true,
				Default:     stringdefault.StaticString(""),
				Computed:    true,
				Description: "Defaults to `\"\"`.",
			},
			"webhook_credential": schema.StringAttribute{
				Optional: true,
			},
			"ask_labels_on_launch": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Prompt on launch for Labels. Defaults to `false`.",
			},
			"ask_skip_tags_on_launch": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Prompt on launch for Skip Tags. Defaults to `false`.",
			},
			"ask_tags_on_launch": schema.BoolAttribute{
				Optional:    true,
				Default:     booldefault.StaticBool(false),
				Computed:    true,
				Description: "Prompt on launch for Ask Tags. Defaults to `false`.",
			},
			"skip_tags": schema.StringAttribute{
				Optional:    true,
				Description: "Tags are useful when you have a large playbook, and you want to run a specific part of a play or task. Use commas to separate multiple tags.",
			},
			"job_tags": schema.StringAttribute{
				Optional:    true,
				Description: "Skip tags are useful when you have a large playbook, and you want to skip specific parts of a play or task. Use commas to separate multiple tags.",
			},
		},
	}
}

func (r *WorkflowJobTemplatesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	configureData, ok := req.ProviderData.(*AwxClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = configureData
}

func (r *WorkflowJobTemplatesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data WorkflowJobTemplatesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var bodyData WorkflowJobTemplateAPIModel

	if !data.Name.IsNull() {
		bodyData.Name = data.Name.ValueString()
	}
	if !data.Description.IsNull() {
		bodyData.Description = data.Description.ValueString()
	}
	if !data.ExtraVars.IsNull() {
		bodyData.ExtraVars = data.ExtraVars.ValueString()
	}
	if !data.Organization.IsNull() {
		bodyData.Organization = int(data.Organization.ValueInt32())
	}
	if !data.SurveyEnabled.IsNull() {
		bodyData.SurveyEnabled = data.SurveyEnabled.ValueBool()
	}
	if !data.AllowSimultaneous.IsNull() {
		bodyData.AllowSimultaneous = data.AllowSimultaneous.ValueBool()
	}
	if !data.AskVariablesOnLaunch.IsNull() {
		bodyData.AskVariablesOnLaunch = data.AskVariablesOnLaunch.ValueBool()
	}
	if !data.Inventory.IsNull() {
		bodyData.Inventory = int(data.Inventory.ValueInt32())
	}
	if !data.Limit.IsNull() {
		bodyData.Limit = data.Limit.ValueString()
	}
	if !data.ScmBranch.IsNull() {
		bodyData.ScmBranch = data.ScmBranch.ValueString()
	}
	if !data.AskInventoryOnLaunch.IsNull() {
		bodyData.AskInventoryOnLaunch = data.AskInventoryOnLaunch.ValueBool()
	}
	if !data.AskScmBranchOnLaunch.IsNull() {
		bodyData.AskScmBranchOnLaunch = data.AskScmBranchOnLaunch.ValueBool()
	}
	if !data.AskLimitOnLaunch.IsNull() {
		bodyData.AskLimitOnLaunch = data.AskLimitOnLaunch.ValueBool()
	}
	if !data.WebhookService.IsNull() {
		bodyData.WebhookService = data.WebhookService.ValueString()
	}
	if !data.WebhookCredential.IsNull() {
		bodyData.WebhookCredential = data.WebhookCredential.ValueString()
	}
	if !data.AskLabelsOnLaunch.IsNull() {
		bodyData.AskLabelsOnLaunch = data.AskLabelsOnLaunch.ValueBool()
	}
	if !data.AskSkipTagsOnLaunch.IsNull() {
		bodyData.AskSkipTagsOnLaunch = data.AskSkipTagsOnLaunch.ValueBool()
	}
	if !data.AskTagsOnLaunch.IsNull() {
		bodyData.AskTagsOnLaunch = data.AskTagsOnLaunch.ValueBool()
	}
	if !data.SkipTags.IsNull() {
		bodyData.SkipTags = data.SkipTags.ValueString()
	}
	if !data.JobTags.IsNull() {
		bodyData.JobTags = data.JobTags.ValueString()
	}

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable marshal json",
			fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
	}

	url := r.client.endpoint + "/api/v2/workflow_job_templates/"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, strings.NewReader(string(jsonData)))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", r.client.auth)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}
	if httpResp.StatusCode != 201 {
		defer httpResp.Body.Close()
		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable read http request response body.",
				err.Error())
			return
		}

		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 201, got %v with message %s. ", httpResp.StatusCode, body))
		return
	}

	tmp := struct {
		Id int `json:"id"`
	}{}

	defer httpResp.Body.Close()
	httpRespBodyData, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get http response body",
			fmt.Sprintf("Error was %v", err))
		return
	}
	err = json.Unmarshal(httpRespBodyData, &tmp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get unmarshal http response to grab ID",
			fmt.Sprintf("error was %v", err))
		return
	}

	idAsString := strconv.Itoa(tmp.Id)

	data.Id = types.StringValue(idAsString)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowJobTemplatesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data WorkflowJobTemplatesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
	}
	url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_templates/%d/", id)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", r.client.auth)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}
	if httpResp.StatusCode != 200 && httpResp.StatusCode != 404 {
		defer httpResp.Body.Close()
		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable read http request response body.",
				err.Error())
			return
		}

		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v with message %s. ", httpResp.StatusCode, body))
		return
	}

	if httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	var responseData WorkflowJobTemplateAPIModel

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to get all data out of the http response data body",
			fmt.Sprintf("Body got %v. ", body))
		return
	}

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable unmarshal response body into object",
			fmt.Sprintf("Error =  %v. ", err.Error()))
		return
	}

	if !(data.Name.IsNull() && responseData.Name == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), responseData.Name)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !(data.Description.IsNull() && responseData.Description == nil) {
		castString, ok := responseData.Description.(string)
		if !ok {
			resp.Diagnostics.AddError("Unable to cast descr as string.", "unable to cast as string.")
			return
		}
		if castString != "" {
			resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), castString)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
	}

	if !(data.ExtraVars.IsNull() && responseData.ExtraVars == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("extra_vars"), responseData.ExtraVars)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if !(data.Organization.IsNull() && responseData.Organization == 0) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization"), responseData.Organization)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("survey_enabled"), responseData.SurveyEnabled)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("allow_simultaneous"), responseData.AllowSimultaneous)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ask_variables_on_launch"), responseData.AskVariablesOnLaunch)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !(data.Inventory.IsNull() && responseData.Inventory == nil) {
		convertInt, ok := responseData.Inventory.(float64)
		if !ok {
			resp.Diagnostics.AddError("Unable to convert any to int32.", "Unable to convert any to int32.")
			return
		}

		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("inventory"), int(convertInt))...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !(data.Limit.IsNull() && responseData.Limit == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("limit"), responseData.Limit)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !(data.ScmBranch.IsNull() && responseData.ScmBranch == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("scm_branch"), responseData.ScmBranch)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ask_inventory_on_launch"), responseData.AskInventoryOnLaunch)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ask_scm_branch_on_launch"), responseData.AskScmBranchOnLaunch)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ask_limit_on_launch"), responseData.AskLimitOnLaunch)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("webhook_service"), responseData.WebhookService)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !(data.WebhookCredential.IsNull() && responseData.WebhookCredential == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("webhook_credential"), responseData.WebhookCredential)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ask_labels_on_launch"), responseData.AskLabelsOnLaunch)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ask_skip_tags_on_launch"), responseData.AskSkipTagsOnLaunch)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("ask_tags_on_launch"), responseData.AskTagsOnLaunch)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !(data.SkipTags.IsNull() && responseData.SkipTags == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("skip_tags"), responseData.SkipTags)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}
	if !(data.JobTags.IsNull() && responseData.JobTags == "") {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("job_tags"), responseData.JobTags)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}

func (r *WorkflowJobTemplatesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data WorkflowJobTemplatesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
		return
	}

	var bodyData WorkflowJobTemplateAPIModel
	bodyData.Name = data.Name.ValueString()
	if !data.Description.IsNull() {
		bodyData.Description = data.Description.ValueString()
	}
	if !data.ExtraVars.IsNull() {
		bodyData.ExtraVars = data.ExtraVars.ValueString()
	}
	bodyData.Organization = int(data.Organization.ValueInt32())
	bodyData.SurveyEnabled = data.SurveyEnabled.ValueBool()
	bodyData.AllowSimultaneous = data.AllowSimultaneous.ValueBool()
	bodyData.AskVariablesOnLaunch = data.AskVariablesOnLaunch.ValueBool()
	if !data.Inventory.IsNull() {
		bodyData.Inventory = int(data.Inventory.ValueInt32())
	}
	bodyData.Limit = data.Limit.ValueString()
	bodyData.ScmBranch = data.ScmBranch.ValueString()
	bodyData.AskInventoryOnLaunch = data.AskInventoryOnLaunch.ValueBool()
	bodyData.AskScmBranchOnLaunch = data.AskScmBranchOnLaunch.ValueBool()
	bodyData.AskLimitOnLaunch = data.AskLimitOnLaunch.ValueBool()
	bodyData.WebhookService = data.WebhookService.ValueString()
	bodyData.WebhookCredential = data.WebhookCredential.ValueString()
	bodyData.AskLabelsOnLaunch = data.AskLabelsOnLaunch.ValueBool()
	bodyData.AskSkipTagsOnLaunch = data.AskSkipTagsOnLaunch.ValueBool()
	bodyData.AskTagsOnLaunch = data.AskTagsOnLaunch.ValueBool()
	bodyData.SkipTags = data.SkipTags.ValueString()
	bodyData.JobTags = data.JobTags.ValueString()

	jsonData, err := json.Marshal(bodyData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable marshal json",
			fmt.Sprintf("Unable to convert id: %+v. ", bodyData))
		return
	}

	url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_templates/%d/", id)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPut, url, strings.NewReader(string(jsonData)))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", r.client.auth)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create example, got error: %s", err))
		return
	}
	if httpResp.StatusCode != 200 {
		defer httpResp.Body.Close()
		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable read http request response body.",
				err.Error())
			return
		}

		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v with message %s. ", httpResp.StatusCode, body))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *WorkflowJobTemplatesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data WorkflowJobTemplatesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int",
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
	}
	url := r.client.endpoint + fmt.Sprintf("/api/v2/workflow_job_templates/%d/", id)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate delete request",
			fmt.Sprintf("Unable to gen url: %v. ", url))
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", r.client.auth)

	httpResp, err := r.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete got error: %s", err))
	}
	if httpResp.StatusCode != 204 {
		defer httpResp.Body.Close()
		body, err := io.ReadAll(httpResp.Body)
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable read http request response body.",
				err.Error())
			return
		}

		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 204, got %v with message %s. ", httpResp.StatusCode, body))
		return

	}
}

func (r *WorkflowJobTemplatesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
