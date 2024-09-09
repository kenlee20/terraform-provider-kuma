package provider

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MonitorResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	PathName    types.String `tfsdk:"path_name"`
	Url         types.String `tfsdk:"url"`
	Method      types.String `tfsdk:"method"`
	// Active      types.Bool   `tfsdk:"active"`
	Type types.String `tfsdk:"type"`

	// Timeout        types.Int64 `tfsdk:"timeout"`
	Interval       types.Int64 `tfsdk:"interval"`
	RetryInterval  types.Int64 `tfsdk:"retry_interval"`
	ResendInterval types.Int64 `tfsdk:"resend_interval"`
	MaxRetries     types.Int64 `tfsdk:"max_retries"`
	MaxRedirects   types.Int64 `tfsdk:"max_redirects"`

	AcceptedStatusCodes []types.String `tfsdk:"accepted_statuscodes"`
	// NotificationIDList  []types.Int64  `tfsdk:"notification_id_list"`
	ExpiryNotification types.Bool `tfsdk:"expiry_notification"`
	IgnoreTls          types.Bool `tfsdk:"ignore_tls"`
	UpsideDown         types.Bool `tfsdk:"upside_down"`
}

func ConvertStruct(src interface{}, dest interface{}, toModel bool) error {
	srcVal := reflect.ValueOf(src)
	srcValType := reflect.TypeOf(src)
	destVal := reflect.ValueOf(dest).Elem()

	if srcVal.Kind() == reflect.Ptr {
		return fmt.Errorf("source is not a pointer")
	}

	for i := 0; i < srcVal.NumField(); i++ {
		switch srcVal.Kind() {
		case reflect.Struct:
			fieldName := srcValType.Field(i).Name
			sourceField := srcVal.Field(i)
			destField := destVal.FieldByName(fieldName)

			if !destField.IsValid() || !destField.CanSet() {
				continue
			}
			if err := convert(sourceField, destField, toModel); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupport target type: %s", srcVal.Kind())
		}
	}
	return nil
}

func convert(src, dest reflect.Value, toModel bool) error {
	srcField := src.Interface()
	destVal := dest

	switch toModel {
	case true:
		fmt.Println(src.Kind(), src)
		switch src.Kind() {
		case reflect.String:
			if v, ok := srcField.(string); ok {
				destVal.Set(reflect.ValueOf(types.StringValue(v)))
			}
		case reflect.Int:
			if v, ok := srcField.(int); ok {
				destVal.Set(reflect.ValueOf(types.Int64Value(int64(v))))
			}
		case reflect.Bool:
			if v, ok := srcField.(bool); ok {
				destVal.Set(reflect.ValueOf(types.BoolValue(v)))
			}
		case reflect.Slice:
			for i := 0; i < src.Len(); i++ {
				sourceField := src.Index(i)
				elemType := destVal.Type().Elem()
				destElem := reflect.New(elemType).Elem()

				if err := convert(sourceField, destElem, toModel); err != nil {
					return err
				}
				destVal.Set(reflect.Append(destVal, destElem))
			}

		default:
			return fmt.Errorf("unsupport source type: %s", destVal.Kind())
		}
	case false:
		switch destVal.Kind() {
		case reflect.String:
			if v, ok := srcField.(types.String); ok {
				destVal.SetString(v.ValueString())
			}
		case reflect.Int:
			if v, ok := srcField.(types.Int64); ok {
				destVal.SetInt(v.ValueInt64())
			}
		case reflect.Bool:
			if v, ok := srcField.(types.Bool); ok {
				destVal.SetBool(v.ValueBool())
			}
		case reflect.Slice:
			for i := 0; i < src.Len(); i++ {
				sourceField := src.Index(i)
				elemType := destVal.Type().Elem()
				destElem := reflect.New(elemType).Elem()

				if err := convert(sourceField, destElem, toModel); err != nil {
					return err
				}
				destVal.Set(reflect.Append(destVal, destElem))
			}
		default:
			return fmt.Errorf("unsupport target type: %s", destVal.Kind())
		}
	}
	return nil
}
