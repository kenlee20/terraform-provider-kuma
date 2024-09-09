// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-kuma/internal/kuma"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &httpMonitorResource{}
var _ resource.ResourceWithImportState = &httpMonitorResource{}

func NewHttpMonitorResource() resource.Resource {
	return &httpMonitorResource{}
}

// ExampleResource defines the resource implementation.
type httpMonitorResource struct {
	client *kuma.Client
}

func (r *httpMonitorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_monitor"
}

func (r *httpMonitorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Example resource",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Example identifier",
			},
			"path_name": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Example identifier",
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Example identifier",
				Default:             stringdefault.StaticString(""),
			},
			"type": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				Default:             stringdefault.StaticString("http"),
			},
			"url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Example identifier",
			},
			"method": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Example identifier",
				Default:             stringdefault.StaticString("GET"),
			},
			// "active": schema.BoolAttribute{
			// 	Optional:            true,
			// 	Computed:            true,
			// 	MarkdownDescription: "Example identifier",
			// 	Default:             booldefault.StaticBool(true),
			// },
			// "timeout": schema.Int64Attribute{
			// 	Optional:            true,
			// 	Computed:            true,
			// 	MarkdownDescription: "Example identifier",
			// 	Default:             int64default.StaticInt64(60),
			// },
			"interval": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Example identifier",
				Default:             int64default.StaticInt64(60),
			},
			"retry_interval": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Example identifier",
				Default:             int64default.StaticInt64(20),
			},
			"resend_interval": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Example identifier",
				Default:             int64default.StaticInt64(0),
			},
			"max_retries": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Example identifier",
				Default:             int64default.StaticInt64(5),
			},
			"max_redirects": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Example identifier",
				Default:             int64default.StaticInt64(10),
			},
			"accepted_statuscodes": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Example identifier",
				ElementType:         types.StringType,
				Default: listdefault.StaticValue(
					types.ListValueMust(types.StringType, []attr.Value{types.StringValue("200-299")}),
				),
			},
			"expiry_notification": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Example identifier",
				Default:             booldefault.StaticBool(true),
			},
			"ignore_tls": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Example identifier",
				Default:             booldefault.StaticBool(true),
			},
			"upside_down": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Example identifier",
				Default:             booldefault.StaticBool(false),
			},
		},
	}

}

func (r *httpMonitorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*kuma.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *httpMonitorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan, o_plan MonitorResourceModel
	var item kuma.Monitor

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "[INPUT_PLAN]"+fmt.Sprintf("%+v", plan))

	err := ConvertStruct(plan, &item, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Convert Struct",
			"Could not Convert terraform struct to api struct, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "[INPUT_ITEM]"+fmt.Sprintf("%+v", item))
	// Create new order and set the ID on the state.

	monitorID, err := r.client.CreateMonitor(item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating monitor",
			"Could not create monitor, unexpected error: "+err.Error(),
		)
		return
	}

	monitor, err := r.client.GetMonitor(*monitorID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kuma Monitor",
			fmt.Sprintf("Could not read Kuma Monitor %s, ID: %d %s", plan.Name.ValueString(), int(plan.ID.ValueInt64()), err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "[OUTPUT_ITEM]"+fmt.Sprintf("%+v", monitor))

	err = ConvertStruct(*monitor, &o_plan, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Convert Struct",
			"Could not Convert terraform struct to api struct, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "[OUTPUT_PLAN]"+fmt.Sprintf("%+v", o_plan))

	// Set state to fully populated data

	diags = resp.State.Set(ctx, o_plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *httpMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MonitorResourceModel
	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := r.client.GetMonitor(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kuma Tag",
			fmt.Sprintf("Could not read Kuma Tag %s, ID: %d %s", state.Name.ValueString(), int(state.ID.ValueInt64()), err.Error()),
		)
		return
	}

	if err := ConvertStruct(*monitor, &state, true); err != nil {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *httpMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, o_plan MonitorResourceModel
	var item kuma.Monitor

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ConvertStruct(plan, &item, false); err != nil {
		resp.Diagnostics.AddError(
			"Error Convert Struct",
			"Could not Convert terraform struct to api struct, unexpected error: "+err.Error(),
		)
		return
	}

	err := r.client.UpdateMonitor(int(plan.ID.ValueInt64()), item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Kuma Monitor",
			fmt.Sprintf("Could not update Kuma Monitor %s, ID: %d %s", plan.Name.ValueString(), int(plan.ID.ValueInt64()), err.Error()),
		)
		return
	}

	monitor, err := r.client.GetMonitor(int(plan.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kuma Monitor",
			fmt.Sprintf("Could not read Kuma Monitor %s, ID: %d %s", plan.Name.ValueString(), int(plan.ID.ValueInt64()), err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "[OUTPUT_ITEM]"+fmt.Sprintf("%+v", monitor))

	if err := ConvertStruct(*monitor, &o_plan, true); err != nil {
		resp.Diagnostics.AddError(
			"Error Convert Struct",
			"Could not Convert api struct to terraform struct , unexpected error: "+err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, o_plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *httpMonitorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MonitorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteMonitor(int(state.ID.ValueInt64()))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Kuma Monitor",
			fmt.Sprintf("Could not delete Kuma Monitor %s, ID: %d %s", state.Name.ValueString(), int(state.ID.ValueInt64()), err.Error()),
		)
		return
	}
}

func (r *httpMonitorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
