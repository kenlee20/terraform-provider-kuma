package provider

import (
	"context"
	"fmt"
	"terraform-provider-upkuapi/internal/kuma"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource              = &NotificationsDataSource{}
	_ datasource.DataSourceWithConfigure = &NotificationsDataSource{}
)

func NewNotificationsDataSource() datasource.DataSource {
	return &NotificationsDataSource{}
}

type Notifications struct {
	Notifications []Notification `tfsdk:"notifications"`
}

type NotificationsDataSource struct {
	client *kuma.Client
}

func (d *NotificationsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_notifications"
}

// Schema defines the schema for the data source.
func (d *NotificationsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"notifications": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.Int64Attribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"user_id": schema.Int64Attribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"default": schema.BoolAttribute{
							Computed: true,
						},
						"active": schema.BoolAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *NotificationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var states Notifications

	notifications, err := d.client.GetNotifications()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Notifications",
			err.Error(),
		)
		return
	}

	for _, notification := range notifications {
		var state Notification

		state.ConvertFrom(notification)

		states.Notifications = append(states.Notifications, state)
	}

	// Set state
	diags := resp.State.Set(ctx, &states)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Configure adds the provider configured client to the data source.
func (d *NotificationsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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

	d.client = client
}
