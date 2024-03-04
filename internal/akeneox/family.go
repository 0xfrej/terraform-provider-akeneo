package akeneox

import (
	"fmt"
	goakeneo "github.com/ezifyio/go-akeneo"
)

const (
	familyPath       = "/api/rest/v1/families"
	familySinglePath = "/api/rest/v1/families/%s"
)

type FamilyService struct {
	goakeneo.FamilyService
	client *goakeneo.Client
}

func NewFamilyClient(client *goakeneo.Client) *FamilyService {
	return &FamilyService{
		FamilyService: client.Family,
		client:        client,
	}
}

func (a *FamilyService) UpdateFamily(family goakeneo.Family) (*goakeneo.Family, error) {
	response := new(goakeneo.Family)
	err := a.client.PATCH(
		fmt.Sprintf(familySinglePath, family.Code),
		nil,
		family,
		response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}
