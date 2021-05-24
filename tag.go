package goenv

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

type FieldTag struct {
	Type         string
	Source       string
	Value        string
	DefaultValue string
	IsOmitempty  bool
}

func getFieldTag(fieldDescription reflect.StructField) (*FieldTag, error) {

	var tag FieldTag

	if value, ok := fieldDescription.Tag.Lookup(tagDefault); ok {
		tag.DefaultValue = value
	}

	if value, ok := fieldDescription.Tag.Lookup(tagEnv); ok {

		tag.Type = tagEnv
		tag.Source = value

		err := tag.getTagOptions()

		if err != nil {
			return nil, err
		}

		if len(tag.Source) == 0 {
			return nil, fmt.Errorf("empty field[%s] tag", fieldDescription.Name)
		}

		return &tag, nil
	}

	if value, ok := fieldDescription.Tag.Lookup(tagFile); ok {

		tag.Type = tagFile
		tag.Source = value

		err := tag.getTagOptions()

		if err != nil {
			return nil, err
		}

		if len(tag.Source) == 0 {
			return nil, fmt.Errorf("empty field[%s] tag", fieldDescription.Name)
		}

		return &tag, nil
	}

	if value, ok := fieldDescription.Tag.Lookup(tagCall); ok {

		tag.Type = tagCall
		tag.Source = value

		err := tag.getTagOptions()

		if err != nil {
			return nil, err
		}

		if len(tag.Source) == 0 {
			return nil, fmt.Errorf("empty field[%s] tag", fieldDescription.Name)
		}

		return &tag, nil
	}

	return nil, nil
}

func (t *FieldTag) getTagOptions() error {

	index := strings.Index(t.Source, ",")

	if index <= 0 {
		return nil
	}

	options := strings.Split(t.Source[index+1:], ",")

	t.Source = t.Source[:index]

	for _, option := range options {

		if len(option) == 0 {
			continue
		}

		switch option {
		case optionOmitempty:
			t.IsOmitempty = true
		default:
			return fmt.Errorf("undefined tag[%s] option[%s]", t.Source, option)
		}
	}

	return nil
}

func (t *FieldTag) getEnvValue() error {

	if value, ok := os.LookupEnv(t.Source); ok {
		t.Value = value
		return nil
	}

	if len(t.DefaultValue) != 0 {
		t.Value = t.DefaultValue
		return nil
	}

	if t.IsOmitempty {
		return nil
	}

	return fmt.Errorf("empty required environment variable[%s]", t.Source)
}
