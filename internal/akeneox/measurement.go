package akeneox

import (
	goakeneo "github.com/ezifyio/go-akeneo"
)

const (
	measurementFamilyPath = "/api/rest/v1/measurement-families"
)

type MeasurementFamilyService struct {
	client *goakeneo.Client
}

func NewMeasurementFamilyClient(client *goakeneo.Client) *MeasurementFamilyService {
	return &MeasurementFamilyService{
		client: client,
	}
}

func (a *MeasurementFamilyService) GetMeasurementFamily(code string) (*MeasurementFamily, error) {
	data := new(MeasurementFamily)
	response := new([]MeasurementFamily)
	err := a.client.GET(
		measurementFamilyPath,
		nil,
		nil,
		response,
	)
	if err != nil {
		return nil, err
	}

	for _, r := range *response {
		if r.Code == code {
			selected := r
			data = &selected
			break
		}
	}

	return data, nil
}

func (a *MeasurementFamilyService) UpdateMeasurementFamilies(families []MeasurementFamily) (*[]MeasurementFamilyPatchResponse, error) {
	response := new([]MeasurementFamilyPatchResponse)
	err := a.client.PATCH(
		measurementFamilyPath,
		nil,
		families,
		response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}
