package goenv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

func ParseWithFile(target interface{}, path string) error {

	if len(path) == 0 {
		return errors.New("empty file path")
	}

	valueType := reflect.TypeOf(target)

	if valueType.Kind() != reflect.Ptr {
		return fmt.Errorf("pointer to struct required")
	}

	valueType = valueType.Elem()

	if valueType.Kind() != reflect.Struct {
		return fmt.Errorf("pointer to struct required")
	}

	extension, err := getFileExtension(path)

	if err != nil {
		return err
	} else if len(extension) == 0 {
		return fmt.Errorf("invalid file extension[%s]", path)
	}

	switch extension {
	case "json":
		if err = getJsonValue(target, path); err != nil {
			return err
		}
	case "yml":
		if err = getYamlValue(target, path); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid file extension[%s]", path)
	}

	return parseStruct(reflect.ValueOf(target).Elem(), valueType, false)
}

func getFileExtension(path string) (string, error) {

	index := strings.LastIndex(path, ".")

	if index <= 0 {
		return "", fmt.Errorf("invalid file path[%s]", path)
	}

	return path[index+1:], nil
}

func getJsonValue(target interface{}, path string) error {

	file, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	return json.Unmarshal(file, target)
}

func getYamlValue(target interface{}, path string) error {

	file, err := ioutil.ReadFile(path)

	if err != nil {
		return err
	}

	return yaml.Unmarshal(file, target)
}
