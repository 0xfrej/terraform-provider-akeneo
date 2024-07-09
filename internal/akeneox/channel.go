package akeneox

import (
	"fmt"

	goakeneo "github.com/ezifyio/go-akeneo"
)

const (
	channelPath       = "/api/rest/v1/channels"
	channelSinglePath = "/api/rest/v1/channels/%s"
)

type ChannelService struct {
	goakeneo.ChannelService
	client *goakeneo.Client
}

func NewChannelClient(client *goakeneo.Client) *ChannelService {
	return &ChannelService{
		ChannelService: client.Channel,
		client:         client,
	}
}

func (a *ChannelService) CreateChannel(channel goakeneo.Channel) error {
	return a.client.POST(
		channelPath,
		nil,
		channel,
		nil,
	)
}

func (a *ChannelService) UpdateChannel(channel goakeneo.Channel) (*goakeneo.Channel, error) {
	response := new(goakeneo.Channel)
	err := a.client.PATCH(
		fmt.Sprintf(channelSinglePath, channel.Code),
		nil,
		channel,
		response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (a *ChannelService) GetChannel(code string) (*goakeneo.Channel, error) {
	response := new(goakeneo.Channel)
	err := a.client.GET(
		fmt.Sprintf(channelSinglePath, code),
		nil,
		nil,
		response,
	)
	if err != nil {
		return nil, err
	}
	return response, nil
}
