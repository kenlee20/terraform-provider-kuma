// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"terraform-provider-upkuapi/internal/kuma"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
		MarkdownDescription: "Provides an Moniotr resource. This allows monitors to be created, updated, and deleted.",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Options for monitor display name.",
			},
			"path_name": schema.StringAttribute{
				Computed: true,
			},
			"description": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Describes the monitor.",
			},
			"type": schema.StringAttribute{
				Computed: true,
				Default:  stringdefault.StaticString("http"),
			},
			"url": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Options for monitoring url.",
			},
			"interval": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for heartbeat Interval. default to `60`.",
				Default:             int64default.StaticInt64(60),
			},
			"retry_interval": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for Retry every secend. default to `20`.",
				Default:             int64default.StaticInt64(20),
			},
			"resend_interval": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for resend every times. defaults to `0`",
				Default:             int64default.StaticInt64(0),
			},
			"max_retries": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for maximum retries before the service is marked as down and a notification is sent. default to `5`.",
				Default:             int64default.StaticInt64(5),
			},
			"max_redirects": schema.Int64Attribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for maximum number of redirects to follow. Set to 0 to disable redirects. defaults to `10`",
				Default:             int64default.StaticInt64(10),
			},
			"http_option_method": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for http monitor method. default to `GET`.",
				Validators: []validator.String{
					stringvalidator.OneOf("GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"),
				},
				Default: stringdefault.StaticString("GET"),
			},
			"http_option_body_encoding": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for body encoding. default to `none`.",
				Validators: []validator.String{
					stringvalidator.OneOf("json", "xml"),
				},
				Default: stringdefault.StaticString("json"),
			},
			"http_option_body": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for body content. default to `none`.",
			},
			"http_option_headers": schema.StringAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for http headers.",
			},
			"notification_list": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for notification id list, automatically enable default notifications.",
				ElementType:         types.Int64Type,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplaceIfConfigured(),
				},
			},
			"accepted_statuscodes": schema.ListAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for Accepted Status Codes. Select status codes which are considered as a successful response., defaults to `[\"200-299\"]`",
				ElementType:         types.StringType,
				Default: listdefault.StaticValue(
					types.ListValueMust(types.StringType, []attr.Value{types.StringValue("200-299")}),
				),
			},
			"expiry_notification": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for Certificate Expiry Notification. defaults to `true`.",
				Default:             booldefault.StaticBool(true),
			},
			"ignore_tls": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for ignore TLS/SSL error for HTTPS websites, defaults to `false`.",
				Default:             booldefault.StaticBool(false),
			},
			"upside_down": schema.BoolAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Options for Upside Down Mode. Flip the status upside down. If the service is reachable, it is DOWN. defaults to `false`",
				Default:             booldefault.StaticBool(false),
			},
			"tags": schema.SetNestedAttribute{
				MarkdownDescription: "Options for monitor tag",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Options for name of tag.",
						},
						"value": schema.StringAttribute{
							Required:            true,
							MarkdownDescription: "Options for value of tag.",
						},
						"tag_id": schema.Int64Attribute{
							Computed: true,
						},
					},
				},
				Optional: true,
				Computed: true,
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
	var plan MonitorResourceModel

	tflog.Debug(ctx, "[INIT_PLAN]"+fmt.Sprintf("%+v", req.Plan))

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Debug(ctx, "[INPUT_PLAN]"+fmt.Sprintf("%+v", plan))

	item, diags := plan.Convert()
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Debug(ctx, "[INPUT_ITEM]"+fmt.Sprintf("%+v", item))

	// Create new order and set the ID on the state.
	monitorID, err := r.client.CreateMonitor(*item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating monitor",
			"Could not create monitor, unexpected error: "+err.Error(),
		)
		return
	}

	var monitor *kuma.Monitor

	for _, tag := range item.Tags {
		if err = r.client.CreateMonitorTag(*monitorID, tag); err != nil {
			resp.Diagnostics.AddError(
				"Error Ceating Kuma Monitor Tag",
				fmt.Sprintf("Could not ceate Kuma Monitor tag %s, tags: %+v %s", plan.Name.ValueString(), tag, err.Error()),
			)
			return
		}
	}

	monitor, err = r.client.GetMonitor(*monitorID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kuma Tag",
			"Could not read Kuma Tag",
		)
		return
	}

	tflog.Debug(ctx, "[OUTPUT_ITEM]"+fmt.Sprintf("%+v", monitor))

	diags = plan.ConvertFrom(*monitor)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Debug(ctx, "[OUTPUT_PLAN]"+fmt.Sprintf("%+v", plan))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *httpMonitorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state, outputPlan MonitorResourceModel
	// Read Terraform prior state data into the model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	monitor, err := r.client.GetMonitor(state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kuma Tag",
			fmt.Sprintf("Could not read Kuma Tag %s, ID: %d %s", state.Name.ValueString(), int(state.ID.ValueInt64()), err.Error()),
		)
		return
	}

	diags = outputPlan.ConvertFrom(*monitor)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &outputPlan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *httpMonitorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan, outputPlan MonitorResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, diags := plan.Convert()
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	if err := r.client.UpdateMonitor(plan.ID.ValueInt64(), *item); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Kuma Monitor",
			fmt.Sprintf("Could not update Kuma Monitor %s, ID: %d %s", plan.Name.ValueString(), int(plan.ID.ValueInt64()), err.Error()),
		)
		return
	}

	monitor, err := r.client.GetMonitor(plan.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kuma Monitor",
			fmt.Sprintf("Could not read Kuma Monitor %s, ID: %d %s", plan.Name.ValueString(), int(plan.ID.ValueInt64()), err.Error()),
		)
		return
	}
	tflog.Debug(ctx, "[INPUT_PLAN]"+fmt.Sprintf("%+v", plan))
	tflog.Debug(ctx, "[INPUT_ITEM]"+fmt.Sprintf("%+v", monitor))

	curTag := make(map[string]kuma.MonitorTag)
	planTag := make(map[string]kuma.MonitorTag)

	for _, tag := range monitor.Tags {
		curTag[tag.Name] = tag
	}
	for _, tag := range item.Tags {
		planTag[tag.Name] = tag
	}

	for name, tag := range curTag {
		// check current tag isn't in plan tag
		if _, ok := planTag[name]; !ok {
			if err := r.client.DeleteMonitorTag(plan.ID.ValueInt64(), tag); err != nil {
				resp.Diagnostics.AddError(
					"Error Updating Kuma Tag",
					fmt.Sprintf("Could not Delete Kuma Tag %s, ID: %d %s", name, tag.TagId, err.Error()),
				)
				return
			}
		} else if ok && planTag[name].Value == tag.Value {
			// delete unchange tag
			delete(planTag, name)
		} else if ok && planTag[name].Value != tag.Value {
			// update changed tag
			if err := r.client.DeleteMonitorTag(plan.ID.ValueInt64(), tag); err != nil {
				resp.Diagnostics.AddError(
					"Error Updating Kuma Tag",
					fmt.Sprintf("Could not Delete Kuma Tag %s, ID: %d %s", name, tag.TagId, err.Error()),
				)
				return
			}
		}
	}

	for tag := range planTag {
		if err := r.client.CreateMonitorTag(plan.ID.ValueInt64(), planTag[tag]); err != nil {
			resp.Diagnostics.AddError(
				"Error Updating Kuma Tag",
				fmt.Sprintf("Could not create Kuma Tag %s, ID: %d %s", tag, planTag[tag].TagId, err.Error()),
			)
			return
		}
	}

	// for _, tag := range item.Tags {
	// 	if tag.TagId != 0 {
	// 		continue
	// 	}

	// 	for _, stateTag := range monitor.Tags {
	// 		if stateTag.Name == tag.Name {
	// 			tag.TagId = state.ID
	// 			tag.Value = stateTag.Value

	// 			if err := r.client.DeleteMonitorTag(item.ID, kuma.MonitorTag{
	// 				TagId: state.ID,
	// 				Value: stateTag.Value,
	// 			}); err != nil {
	// 				resp.Diagnostics.AddError(
	// 					"Error Updating Kuma Tag",
	// 					fmt.Sprintf("Could not Delete Kuma Tag %s, ID: %d %s", tag.Name, tag.TagId, err.Error()),
	// 				)
	// 				return
	// 			}
	// 			break
	// 		}
	// 	}
	// 	tflog.Debug(ctx, "[INPUT_ITEM]"+fmt.Sprintf("%+v", tag))

	// 	if err := r.client.CreateMonitorTag(item.ID, tag); err != nil {
	// 		resp.Diagnostics.AddError(
	// 			"Error Updating Kuma Tag",
	// 			fmt.Sprintf("Could not create Kuma Tag %s, ID: %d %s", tag.Name, tag.TagId, err.Error()),
	// 		)
	// 		return
	// 	}
	// }

	monitor, err = r.client.GetMonitor(plan.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kuma Monitor",
			fmt.Sprintf("Could not read Kuma Monitor %s, ID: %d %s", plan.Name.ValueString(), int(plan.ID.ValueInt64()), err.Error()),
		)
		return
	}

	tflog.Debug(ctx, "[OUTPUT_ITEM]"+fmt.Sprintf("%+v", monitor))

	diags = outputPlan.ConvertFrom(*monitor)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	diags = resp.State.Set(ctx, outputPlan)
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
	err := r.client.DeleteMonitor(state.ID.ValueInt64())
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
