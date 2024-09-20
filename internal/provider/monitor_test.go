package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"terraform-provider-kuma/internal/kuma"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func setupTest() (*kuma.Monitor, *MonitorResourceModel, diag.Diagnostics) {
	ctx := context.Background()

	plan := kuma.Monitor{
		ID:                  1,
		Name:                "name",
		Description:         "description test",
		PathName:            "name",
		Url:                 "https://example.com",
		Method:              "GET",
		Type:                "http",
		Interval:            1,
		MaxRedirects:        1,
		RetryInterval:       1,
		ResendInterval:      1,
		MaxRetries:          1,
		ExpiryNotification:  true,
		UpsideDown:          true,
		IgnoreTls:           true,
		AcceptedStatusCodes: []string{"200"},
		NotificationIDList:  []int64{5},
		Tags: []kuma.MonitorTag{
			{
				TagId: 5,
				Name:  "tag_name",
				Value: "tag_value",
			},
		},
	}

	mapTag := make(map[string]string)
	for _, tag := range plan.Tags {
		mapTag[tag.Name] = tag.Value
	}

	providerTag, err := types.MapValueFrom(ctx, types.StringType, mapTag)
	if err != nil {
		return nil, nil, err
	}

	providerPlan := MonitorResourceModel{
		ID:                  types.Int64Value(plan.ID),
		Name:                types.StringValue(plan.Name),
		Description:         types.StringValue(plan.PathName),
		Parent:              types.Int64Value(plan.Parent),
		Url:                 types.StringValue(plan.Url),
		Method:              types.StringValue(plan.Method),
		Type:                types.StringValue(plan.Type),
		Interval:            types.Int64Value(plan.Interval),
		MaxRedirects:        types.Int64Value(plan.MaxRedirects),
		RetryInterval:       types.Int64Value(plan.RetryInterval),
		ResendInterval:      types.Int64Value(plan.ResendInterval),
		MaxRetries:          types.Int64Value(plan.MaxRedirects),
		ExpiryNotification:  types.BoolValue(plan.ExpiryNotification),
		UpsideDown:          types.BoolValue(plan.UpsideDown),
		IgnoreTls:           types.BoolValue(plan.IgnoreTls),
		AcceptedStatusCodes: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("200")}),
		NotificationIDList:  types.ListValueMust(types.Int64Type, []attr.Value{types.Int64Value(5)}),
		Tags:                providerTag,
	}
	return &plan, &providerPlan, nil
}

func TestMonitorConvert(t *testing.T) {
	var providerPlan MonitorResourceModel
	plan, _, diag := setupTest()
	if diag != nil {
		t.Fatal(diag)
	}

	diag = providerPlan.ConvertFrom(*plan)
	if diag.HasError() {
		t.Fatal(diag)
	}

	result, diag := providerPlan.Convert()
	if diag.HasError() {
		t.Fatal(diag)
	}

	byteResult, _ := json.Marshal(result)
	bytePlan, _ := json.Marshal(plan)

	t.Logf("%+v", providerPlan)

	if bytes.Equal(byteResult, bytePlan) {
		t.Fatal("failed to convert")
	}
}
