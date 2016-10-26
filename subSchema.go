// Copyright 2015 xeipuuv ( https://github.com/xeipuuv )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may Not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless Required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author           xeipuuv
// author-github    https://github.com/xeipuuv
// author-mail      xeipuuv@gmail.com
//
// repository-name  gojsonschema
// repository-desc  An implementation of JSON Schema, based on IETF's draft v4 - Go language.
//
// description      Defines the structure of a sub-subSchema.
//                  A sub-subSchema can contain other sub-schemas.
//
// created          27-02-2013

package gojsonschema

import (
	"errors"
	"regexp"
	"strings"

	"github.com/xeipuuv/gojsonreference"
)

const (
	KEY_SCHEMA                = "$subSchema"
	KEY_ID                    = "$id"
	KEY_REF                   = "$ref"
	KEY_TITLE                 = "title"
	KEY_DESCRIPTION           = "description"
	KEY_TYPE                  = "type"
	KEY_ITEMS                 = "items"
	KEY_ADDITIONAL_ITEMS      = "additionalItems"
	KEY_PROPERTIES            = "properties"
	KEY_PATTERN_PROPERTIES    = "patternProperties"
	KEY_ADDITIONAL_PROPERTIES = "additionalProperties"
	KEY_DEFINITIONS           = "definitions"
	KEY_MULTIPLE_OF           = "multipleOf"
	KEY_MINIMUM               = "minimum"
	KEY_MAXIMUM               = "maximum"
	KEY_EXCLUSIVE_MINIMUM     = "exclusiveMinimum"
	KEY_EXCLUSIVE_MAXIMUM     = "exclusiveMaximum"
	KEY_MIN_LENGTH            = "minLength"
	KEY_MAX_LENGTH            = "maxLength"
	KEY_PATTERN               = "pattern"
	KEY_FORMAT                = "format"
	KEY_MIN_PROPERTIES        = "minProperties"
	KEY_MAX_PROPERTIES        = "maxProperties"
	KEY_DEPENDENCIES          = "dependencies"
	KEY_REQUIRED              = "required"
	KEY_READONLY              = "readonly"
	KEY_MIN_ITEMS             = "minItems"
	KEY_MAX_ITEMS             = "maxItems"
	KEY_UNIQUE_ITEMS          = "uniqueItems"
	KEY_ENUM                  = "enum"
	KEY_ONE_OF                = "oneOf"
	KEY_ANY_OF                = "anyOf"
	KEY_ALL_OF                = "allOf"
	KEY_NOT                   = "not"
)

type SubSchema struct {

	// basic subSchema meta properties
	Id          *string
	Title       *string
	Description *string

	Property string

	// Types associated with the subSchema
	Types jsonSchemaType

	// Reference url
	Ref *gojsonreference.JsonReference
	// Schema referenced
	RefSchema *SubSchema
	// Json reference
	SubSchema *gojsonreference.JsonReference

	// hierarchy
	Parent                      *SubSchema
	Definitions                 map[string]*SubSchema
	DefinitionsChildren         []*SubSchema
	ItemsChildren               []*SubSchema
	ItemsChildrenIsSingleSchema bool
	PropertiesChildren          []*SubSchema

	// validation : number / integer
	MultipleOf       *float64
	Maximum          *float64
	ExclusiveMaximum bool
	Minimum          *float64
	ExclusiveMinimum bool

	// validation : string
	MinLength *int
	MaxLength *int
	Pattern   *regexp.Regexp
	Format    string

	// validation : object
	MinProperties *int
	MaxProperties *int
	Required      []string

	Dependencies         map[string]interface{}
	AdditionalProperties interface{}
	PatternProperties    map[string]*SubSchema

	// validation : array
	MinItems    *int
	MaxItems    *int
	UniqueItems bool

	AdditionalItems interface{}

	// validation : all
	Enum []string

	ReadOnly bool

	// validation : subSchema
	OneOf []*SubSchema
	AnyOf []*SubSchema
	AllOf []*SubSchema
	Not   *SubSchema
}

func (s *SubSchema) AddEnum(i interface{}) error {

	is, err := marshalToJsonString(i)
	if err != nil {
		return err
	}

	if isStringInSlice(s.Enum, *is) {
		return errors.New(formatErrorDescription(
			Locale.KeyItemsMustBeUnique(),
			ErrorDetails{"key": KEY_ENUM},
		))
	}

	s.Enum = append(s.Enum, *is)

	return nil
}

func (s *SubSchema) ContainsEnum(i interface{}) (bool, error) {

	is, err := marshalToJsonString(i)
	if err != nil {
		return false, err
	}

	return isStringInSlice(s.Enum, *is), nil
}

func (s *SubSchema) AddOneOf(subSchema *SubSchema) {
	s.OneOf = append(s.OneOf, subSchema)
}

func (s *SubSchema) AddAllOf(subSchema *SubSchema) {
	s.AllOf = append(s.AllOf, subSchema)
}

func (s *SubSchema) AddAnyOf(subSchema *SubSchema) {
	s.AnyOf = append(s.AnyOf, subSchema)
}

func (s *SubSchema) SetNot(subSchema *SubSchema) {
	s.Not = subSchema
}

func (s *SubSchema) AddRequired(value string) error {

	if isStringInSlice(s.Required, value) {
		return errors.New(formatErrorDescription(
			Locale.KeyItemsMustBeUnique(),
			ErrorDetails{"key": KEY_REQUIRED},
		))
	}

	s.Required = append(s.Required, value)

	return nil
}

func (s *SubSchema) AddDefinitionChild(child *SubSchema) {
	s.DefinitionsChildren = append(s.DefinitionsChildren, child)
}

func (s *SubSchema) AddItemsChild(child *SubSchema) {
	s.ItemsChildren = append(s.ItemsChildren, child)
}

func (s *SubSchema) AddPropertiesChild(child *SubSchema) {
	s.PropertiesChildren = append(s.PropertiesChildren, child)
}

func (s *SubSchema) PatternPropertiesString() string {

	if s.PatternProperties == nil || len(s.PatternProperties) == 0 {
		return STRING_UNDEFINED // should never happen
	}

	patternPropertiesKeySlice := []string{}
	for pk := range s.PatternProperties {
		patternPropertiesKeySlice = append(patternPropertiesKeySlice, `"`+pk+`"`)
	}

	if len(patternPropertiesKeySlice) == 1 {
		return patternPropertiesKeySlice[0]
	}

	return "[" + strings.Join(patternPropertiesKeySlice, ",") + "]"

}
