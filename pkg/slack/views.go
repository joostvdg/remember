package slack

import (
	"github.com/slack-go/slack"
)

func SendCreateView(triggerId string, view *slack.ModalViewRequest, apiToken string) error {
	//viewApiEndpoint := "https://slack.com/api/views.open"
	api := slack.New(apiToken)
	_, err := api.OpenView(triggerId, *view)
	return err
}
