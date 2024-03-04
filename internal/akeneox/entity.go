package akeneox

import goakeneo "github.com/ezifyio/go-akeneo"

// AttributeGroup is the struct for an akeneo attribute group
type AttributeGroup struct {
	Links      *goakeneo.Links   `json:"_links,omitempty" mapstructure:"_links"`
	Code       string            `json:"code,omitempty" mapstructure:"code"`
	Attributes []string          `json:"attributes,omitempty" mapstructure:"attributes"`
	SortOrder  int               `json:"sort_order,omitempty" mapstructure:"sort_order"`
	Labels     map[string]string `json:"labels,omitempty" mapstructure:"labels"`
}
