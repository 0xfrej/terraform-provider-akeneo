package akeneox

import (
	"fmt"
	goakeneo "github.com/ezifyio/go-akeneo"
)

const (
	attributePath             = "/api/rest/v1/attributes"
	attributeSinglePath       = "/api/rest/v1/attributes/%s"
	attributeOptionPath       = "/api/rest/v1/attributes/%s/options"
	attributeOptionSinglePath = "/api/rest/v1/attributes/%s/options/%s"
	attributeGroupPath        = "/api/rest/v1/attribute-groups"
	attributeGroupSinglePath  = "/api/rest/v1/attribute-groups/%s"
)

type AttributeService struct {
	goakeneo.AttributeService
	client *goakeneo.Client
}

func NewAttributeClient(client *goakeneo.Client) *AttributeService {
	return &AttributeService{
		AttributeService: client.Attribute,
		client:           client,
	}
}

func (a *AttributeService) CreateAttribute(attribute goakeneo.Attribute) error {
	return a.client.POST(
		attributePath,
		nil,
		attribute,
		nil,
	)
}

func (a *AttributeService) UpdateAttribute(attribute goakeneo.Attribute) (*goakeneo.Attribute, error) {
	response := new(goakeneo.Attribute)
	err := a.client.PATCH(
		fmt.Sprintf(attributeSinglePath, attribute.Code),
		nil,
		attribute,
		response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (a *AttributeService) CreateAttributeOption(option goakeneo.AttributeOption) error {
	return a.client.POST(
		fmt.Sprintf(attributeOptionPath, option.Attribute),
		nil,
		option,
		nil,
	)
}

func (a *AttributeService) UpdateAttributeOption(option goakeneo.AttributeOption) (*goakeneo.AttributeOption, error) {
	response := new(goakeneo.AttributeOption)
	err := a.client.PATCH(
		fmt.Sprintf(attributeOptionSinglePath, option.Attribute, option.Code),
		nil,
		option,
		response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (a *AttributeService) GetAttributeOption(attribute string, code string) (*goakeneo.AttributeOption, error) {
	response := new(goakeneo.AttributeOption)
	err := a.client.GET(
		fmt.Sprintf(attributeOptionSinglePath, attribute, code),
		nil,
		nil,
		response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (a *AttributeService) GetAttributeGroup(code string) (*AttributeGroup, error) {
	response := new(AttributeGroup)
	err := a.client.GET(
		fmt.Sprintf(attributeGroupSinglePath, code),
		nil,
		nil,
		response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (a *AttributeService) CreateAttributeGroup(group AttributeGroup) error {
	return a.client.POST(
		attributeGroupPath,
		nil,
		group,
		nil,
	)
}

func (a *AttributeService) UpdateAttributeGroup(group AttributeGroup) (*AttributeGroup, error) {
	response := new(AttributeGroup)
	err := a.client.PATCH(
		fmt.Sprintf(attributeGroupSinglePath, group.Code),
		nil,
		group,
		response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}
