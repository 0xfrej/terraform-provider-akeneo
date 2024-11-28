package akeneox

import goakeneo "github.com/ezifyio/go-akeneo"

// AttributeGroup is the struct for an akeneo attribute group.
type AttributeGroup struct {
	//Links      *goakeneo.Links   `json:"_links,omitempty" mapstructure:"_links"`
	Code       string            `json:"code,omitempty" mapstructure:"code"`
	Attributes []string          `json:"attributes,omitempty" mapstructure:"attributes"`
	SortOrder  int               `json:"sort_order,omitempty" mapstructure:"sort_order"`
	Labels     map[string]string `json:"labels,omitempty" mapstructure:"labels"`
}

type MeasurementUnitConversion struct {
	Operator string `json:"operator,omitempty" mapstructure:"operator"`
	Value    string `json:"value,omitempty" mapstructure:"value"`
}

type MeasurementUnit struct {
	Code                string                      `json:"code,omitempty" mapstructure:"code"`
	Labels              map[string]string           `json:"labels,omitempty" mapstructure:"labels"`
	ConvertFromStandard []MeasurementUnitConversion `json:"convert_from_standard" mapstructure:"convert_from_standard"`
	Symbol              string                      `json:"symbol,omitempty" mapstructure:"symbol"`
}

type MeasurementFamily struct {
	Code             string                     `json:"code,omitempty" mapstructure:"code"`
	Labels           map[string]string          `json:"labels,omitempty" mapstructure:"labels"`
	StandardUnitCode string                     `json:"standard_unit_code,omitempty" mapstructure:"standard_unit_code"`
	Units            map[string]MeasurementUnit `json:"units,omitempty" mapstructure:"units"`
}

type MeasurementFamilyPatchResponse struct {
	Code       string                     `json:"code"`
	StatusCode int                        `json:"status_code"`
	Message    string                     `json:"message,omitempty"`
	Errors     []goakeneo.ValidationError `json:"errors,omitempty"`
}

type AssociationType struct {
	Code         string            `json:"code,omitempty" mapstructure:"code"`
	Labels       map[string]string `json:"labels,omitempty" mapstructure:"labels"`
	IsQuantified *bool             `json:"is_quantified,omitempty" mapstructure:"is_quantified"`
	IsTwoWay     *bool             `json:"is_two_way,omitempty" mapstructure:"is_two_way"`
}
