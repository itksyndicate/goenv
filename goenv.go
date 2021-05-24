package goenv

import (
	"fmt"
	"log"
	"reflect"
	"strings"
)

func Parse(value interface{}) error {

	valueType := reflect.TypeOf(value)

	if valueType.Kind() != reflect.Ptr {
		return fmt.Errorf("pointer to struct required")
	}

	valueType = valueType.Elem()

	if valueType.Kind() != reflect.Struct {
		return fmt.Errorf("pointer to struct required")
	}

	return parseStruct(reflect.ValueOf(value).Elem(), valueType, true)
}

func parseStruct(value reflect.Value, structType reflect.Type, overwrite bool) error {

	for index := 0; index < value.NumField(); index++ {

		field := value.Field(index)

		if field.CanSet() == false {
			continue
		}

		err := parseField(field, structType.Field(index), overwrite)

		if err != nil {
			return err
		}
	}

	return nil
}

func parseField(field reflect.Value, fieldDescription reflect.StructField, overwrite bool) error {

	tag, err := getFieldTag(fieldDescription)

	if err != nil {
		return err
	}

	fieldType := field.Type()

	if tag == nil {

		if fieldType.Kind() == reflect.Struct {
			return parseStruct(field, field.Type(), overwrite)
		}

		if fieldType.Kind() != reflect.Ptr {
			return nil
		}

		fieldType = fieldType.Elem()

		if fieldType.Kind() != reflect.Struct {
			return nil
		}

		if field.IsNil() {
			field.Set(reflect.New(fieldType))
		}

		return parseStruct(field.Elem(), fieldType, overwrite)
	}

	if tag.Type == tagCall {

		if field.Kind() == reflect.Ptr {

			fieldType = fieldType.Elem()

			if field.IsNil() {
				field.Set(reflect.New(fieldType))
			}
		}

		item := reflect.New(field.Type())

		method := item.Elem().MethodByName(tag.Source)

		if method.IsValid() == false {
			return fmt.Errorf("invalid field[%s] method[%s]", fieldDescription.Name, tag.Source)
		}

		if method.Type().NumOut() != 1 {
			return fmt.Errorf("invalid field[%s] method[%s] result count", fieldDescription.Name, tag.Source)
		}

		if method.Type().Out(0).AssignableTo(field.Type()) == false {
			return fmt.Errorf("invalid field[%s] method[%s] result type", fieldDescription.Name, tag.Source)
		}

		data := item.Elem().MethodByName(tag.Source).Call(nil)

		if len(data) == 0 {
			return fmt.Errorf("empty field[%s] method[%s] call result", fieldDescription.Name, tag.Source)
		}

		field.Set(data[0])

		return nil
	}

	if field.Kind() == reflect.Ptr {

		fieldType = fieldType.Elem()

		if field.IsNil() {
			field.Set(reflect.New(fieldType))
		}

		field = field.Elem()
	}

	if field.IsZero() == false && overwrite == false {
		return nil
	}

	if tag.Type == tagFile {

		item := reflect.New(field.Type())

		err = ParseWithFile(item.Interface(), tag.Source)

		if err != nil {
			return err
		}

		field.Set(item.Elem())

		return nil
	}

	err = tag.getEnvValue()

	if err != nil {
		return err
	}

	if len(tag.Value) == 0 && tag.IsOmitempty {
		return nil
	}

	switch field.Kind() {
	case reflect.Slice:
		return setSliceFieldValue(field, tag.Value)
	case reflect.Map:
		return setMapFieldValue(field, tag.Value)
	default:
		return setFieldValue(field, tag.Value)
	}
}

func setFieldValue(field reflect.Value, value string) error {

	parser, ok := valueParsers[field.Kind()]

	if ok == false {
		return fmt.Errorf("undefined parser for field kind[%s]", field.Kind().String())
	}

	data, err := parser(value)

	if err != nil {
		return err
	}

	field.Set(reflect.ValueOf(data).Convert(field.Type()))

	return nil
}

func setSliceFieldValue(field reflect.Value, value string) error {

	sliceType := field.Type()

	sliceItemType := sliceType.Elem()

	if sliceItemType.Kind() == reflect.Ptr {
		sliceItemType = sliceItemType.Elem()
	}

	parser, ok := valueParsers[sliceItemType.Kind()]

	if ok == false {
		return fmt.Errorf("undefined parser for slice field kind[%s]", sliceItemType.Kind().String())
	}

	values := strings.Split(value, ",")

	result := reflect.MakeSlice(sliceType, 0, len(values))

	for _, value := range values {

		data, err := parser(value)

		if err != nil {
			return err
		}

		var dataItem reflect.Value

		if sliceType.Elem().Kind() == reflect.Ptr {
			dataItem = reflect.New(sliceItemType)
			dataItem.Elem().Set(reflect.ValueOf(data).Convert(sliceItemType))
		} else {
			dataItem = reflect.ValueOf(data).Convert(sliceItemType)
		}

		result = reflect.Append(result, dataItem)
	}

	field.Set(result)

	return nil
}

func setMapFieldValue(field reflect.Value, value string) error {

	mapType := field.Type()

	log.Println("MAP", mapType.Key())

	return nil
}
