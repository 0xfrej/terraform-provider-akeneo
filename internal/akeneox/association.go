package akeneox

import (
	"fmt"

	goakeneo "github.com/ezifyio/go-akeneo"
)

const (
	associationTypesSinglePath = "/api/rest/v1/association-types/%s"
)

type AssociationTypeService struct {
	client *goakeneo.Client
}

func NewAssociationTypeClient(client *goakeneo.Client) *AssociationTypeService {
	return &AssociationTypeService{
		client: client,
	}
}

func (a *AssociationTypeService) UpdateAssociationTypes(association AssociationType) error {
	err := a.client.PATCH(
		fmt.Sprintf(associationTypesSinglePath, association.Code),
		nil,
		association,
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (a *AssociationTypeService) GetAssociationType(code string) (*AssociationType, error) {
	response := new(AssociationType)
	err := a.client.GET(
		fmt.Sprintf(associationTypesSinglePath, code),
		nil,
		nil,
		response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}
