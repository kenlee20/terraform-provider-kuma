package provider

import (
	"context"
	"fmt"
	"terraform-provider-upkuapi/internal/kuma"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &groupResource{}
	_ resource.ResourceWithConfigure   = &groupResource{}
	_ resource.ResourceWithImportState = &groupResource{}
)

// NewGroupResource is a helper function to simplify the provider implementation.
func NewGroupResource() resource.Resource {
	return &groupResource{}
}

type groupResource struct {
	client *kuma.Client
}

// Metadata returns the resource type name.
func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

// Schema defines the schema for the resource.
func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"type": schema.StringAttribute{
				Computed: true,
				Default:  stringdefault.StaticString("group"),
			},
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"tags": schema.MapAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "Options for monitor tag",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan groupModel

	diags := req.Plan.Get(ctx, &plan)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	item, diags := plan.Convert(ctx)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	tflog.Debug(ctx, "[INPUT_PLAN]"+fmt.Sprintf("%+v", plan))
	tflog.Debug(ctx, "[INPUT_ITEM]"+fmt.Sprintf("%+v", item))

	monitorID, err := r.client.CreateMonitor(*item)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating monitor",
			"Could not create monitor, unexpected error: "+err.Error(),
		)
		return
	}

	tflog.Debug(ctx, "Monitor create done")

	for _, tag := range item.Tags {
		if err = r.client.CreateMonitorTag(*monitorID, tag); err != nil {
			resp.Diagnostics.AddError(
				"Error Ceating Kuma Monitor Tag",
				fmt.Sprintf("Could not ceate Kuma Monitor tag %s, tags: %+v %s", plan.Name.ValueString(), tag, err.Error()),
			)
			return
		}
	}

	monitor, err := r.client.GetMonitor(*monitorID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kuma Tag",
			"Could not read Kuma Tag",
		)
		return
	}

	tflog.Debug(ctx, "[OUTPUT_ITEM]"+fmt.Sprintf("%+v", monitor))

	diags = plan.ConvertFrom(ctx, monitor)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state groupModel
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

	diags = state.ConvertFrom(ctx, monitor)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan groupModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	item, diags := plan.Convert(ctx)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Update existing tag
	if err := r.client.UpdateMonitor(plan.ID.ValueInt64(), *item); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Kuma Monitor",
			fmt.Sprintf("Could not update Kuma Monitor %s, ID: %d %s", plan.Name.ValueString(), int(plan.ID.ValueInt64()), err.Error()),
		)
		return
	}

	// Fetch updated tag
	monitor, err := r.client.GetMonitor(plan.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kuma Monitor",
			fmt.Sprintf("Could not read Kuma Monitor %s, ID: %d %s", plan.Name.ValueString(), int(plan.ID.ValueInt64()), err.Error()),
		)
		return
	}

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

	monitor, err = r.client.GetMonitor(plan.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Kuma Monitor",
			fmt.Sprintf("Could not read Kuma Monitor %s, ID: %d %s", plan.Name.ValueString(), int(plan.ID.ValueInt64()), err.Error()),
		)
		return
	}

	diags = plan.ConvertFrom(ctx, monitor)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state groupModel
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

func (r *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*kuma.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
