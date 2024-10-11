package akeneox

import (
	"fmt"

	goakeneo "github.com/ezifyio/go-akeneo"
)

const (
	categoryPath       = "/api/rest/v1/categories"
	categorySinglePath = "/api/rest/v1/categories/%s"
)

type CategoryService struct {
	goakeneo.CategoryService
	client *goakeneo.Client
}

func NewCategoryClient(client *goakeneo.Client) *CategoryService {
	return &CategoryService{
		CategoryService: client.Category,
		client:          client,
	}
}

func (a *CategoryService) CreateCategory(category goakeneo.Category) error {
	return a.client.POST(
		categoryPath,
		nil,
		category,
		nil,
	)
}

func (a *CategoryService) UpdateCategory(category goakeneo.Category) (*goakeneo.Category, error) {
	response := new(goakeneo.Category)
	err := a.client.PATCH(
		fmt.Sprintf(categorySinglePath, category.Code),
		nil,
		category,
		response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (a *CategoryService) GetCategory(code string) (*goakeneo.Category, error) {
	return a.client.Category.Get(code)
}
