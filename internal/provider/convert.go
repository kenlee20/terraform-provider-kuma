package provider

import (
	"fmt"
	"reflect"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
			var destSlice []attr.Value
			var elemType attr.Type

			for i := 0; i < src.Len(); i++ {
				switch src.Index(i).Kind() {
				case reflect.String:
					elemType = types.StringType
					destSlice = append(destSlice, types.StringValue(src.Index(i).String()))
				case reflect.Int:
					elemType = types.Int64Type
					destSlice = append(destSlice, types.Int64Value(int64(src.Index(i).Int())))
				case reflect.Bool:
					elemType = types.BoolType
					destSlice = append(destSlice, types.BoolValue(src.Index(i).Bool()))

				default:
					return fmt.Errorf("unsupport source type: %s", src.Index(i).Kind())
				}
			}
			destVal.Set(reflect.ValueOf(types.ListValueMust(elemType, destSlice)))

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
			source := src.Interface().(types.List)

			for _, v := range source.Elements() {
				elemType := destVal.Type().Elem()
				destElem := reflect.New(elemType).Elem()

				if err := convert(reflect.ValueOf(v), destElem, toModel); err != nil {
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
