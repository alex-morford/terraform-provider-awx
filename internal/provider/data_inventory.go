package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &InventoryDataSource{}

func NewInventoryDataSource() datasource.DataSource {
	return &InventoryDataSource{}
}

// InventoryDataSource defines the data source implementation.
type InventoryDataSource struct {
	client *AwxClient
}

func (d *InventoryDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_inventory"
}

func (d *InventoryDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Get inventory datasource",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Inventory ID.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Inventory name.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Inventory description.",
				Computed:    true,
			},
			"organization": schema.Int32Attribute{
				Description: "Organization ID for the inventory to live in.",
				Computed:    true,
			},
			"variables": schema.StringAttribute{
				Description: "Enter inventory variables using either JSON or YAML syntax.",
				Computed:    true,
			},
			"kind": schema.StringAttribute{
				Description: "Set to `smart` for smart inventories",
				Computed:    true,
			},
			"host_filter": schema.StringAttribute{
				Description: "Populate the hosts for this inventory by using a search filter. Example: ansible_facts__ansible_distribution:\"RedHat\".",
				Computed:    true,
			},
		},
	}
}

func (d *InventoryDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	configureData, ok := req.ProviderData.(*AwxClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = configureData
}

func (d *InventoryDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data InventoryModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var url string

	// set url for read by id HTTP request
	id, err := strconv.Atoi(data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable convert id from string to int.",
			fmt.Sprintf("Unable to convert id: %v. ", data.Id.ValueString()))
		return
	}
	url = d.client.endpoint + fmt.Sprintf("/api/v2/inventories/%d/", id)

	// create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to generate request",
			fmt.Sprintf("Unable to gen url: %v.", url))
		return
	}

	httpReq.Header.Add("Content-Type", "application/json")
	httpReq.Header.Add("Authorization", "Bearer"+" "+d.client.token)

	httpResp, err := d.client.client.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Client Error",
			fmt.Sprintf("Unable to read inventory, got error: %s", err))
		return
	}
	if httpResp.StatusCode != 200 && httpResp.StatusCode != 404 {
		resp.Diagnostics.AddError(
			"Bad request status code.",
			fmt.Sprintf("Expected 200, got %v. ", httpResp.StatusCode))
		return
	}

	if httpResp.StatusCode == 404 {
		resp.State.RemoveResource(ctx)
		return
	}

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read the http response data body",
			fmt.Sprintf("Body: %v.", body))
		return
	}

	var responseData InventoryAPIModel

	err = json.Unmarshal(body, &responseData)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to unmarshall response body into object",
			fmt.Sprintf("Error =  %v.", err.Error()))
		return
	}

	idAsString := strconv.Itoa(responseData.Id)
	data.Id = types.StringValue(idAsString)

	data.Name = types.StringValue(responseData.Name)
	data.Organization = types.Int32Value(int32(responseData.Organization))

	if responseData.Description != "" {
		data.Description = types.StringValue(responseData.Description)
	}

	if responseData.Variables != "" {
		data.Variables = types.StringValue(responseData.Variables)
	}

	if responseData.Kind != "" {
		data.Kind = types.StringValue(responseData.Kind)
	}

	if responseData.HostFilter != "" {
		data.HostFilter = types.StringValue(responseData.HostFilter)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
