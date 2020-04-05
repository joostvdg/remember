package slack

import (
	"fmt"
	"github.com/joostvdg/remember/pkg/remember"
	"github.com/joostvdg/remember/pkg/store"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

func HandleList(slackCommand SlackCommand, log *zap.SugaredLogger, store store.MemoryStore) *slack.Msg {
	command := ""
	if len(slackCommand.TokenizedText) > 0 {
		command = slackCommand.TokenizedText[0]
	}

	switch command {
	case "view":
		return GetLists(slackCommand, log, store)
	default:
		return GetLists(slackCommand, log, store)
	}

}

func GetLists(slackCommand SlackCommand, log *zap.SugaredLogger, store store.MemoryStore) *slack.Msg {
	log.Infof("Remember GetLists invoked for: %v (%v), in %v", slackCommand.UserName, slackCommand.UserId, slackCommand.Team)

	userLists := []*remember.MediaList{}
	userFound, userExists := store.FindUser(slackCommand.UserId)
	if !userExists {
		user := remember.User{
			Id:    slackCommand.UserId,
			Name:  slackCommand.UserName,
			Lists: nil,
		}
		store.AddUser(user)
	} else {
		userLists = userFound.Lists
	}

	numberOfElements := 2 + len(userLists)
	elementsAdded := 0
	var blockSet []slack.Block
	blockSet = make([]slack.Block, numberOfElements, numberOfElements)

	mainText := fmt.Sprintf("We found %v of lists for you:", len(userLists))
	mainSectionText := slack.TextBlockObject{
		Type: "plain_text",
		Text: mainText,
	}
	mainSection := slack.SectionBlock{
		Type: "section",
		Text: &mainSectionText,
	}
	blockSet[elementsAdded] = mainSection
	elementsAdded++

	dividerBlock := slack.DividerBlock{
		Type: "divider",
	}
	blockSet[elementsAdded] = dividerBlock
	elementsAdded++

	for _, list := range userLists {
		text := slack.TextBlockObject{
			Type: "mrkdwn",
			Text: fmt.Sprintf("- *%v* (id: %v)", list.Name, list.Id),
		}
		buttonText := slack.TextBlockObject{
			Type: "plain_text",
			Text: "Open",
		}
		button := slack.ButtonBlockElement{
			Type:     "button",
			Text:     &buttonText,
			ActionID: "open_list",
			Value:    list.Id,
		}
		accessory := slack.Accessory{
			ButtonElement: &button,
		}
		section := slack.SectionBlock{
			Type:      "section",
			Text:      &text,
			Accessory: &accessory,
		}
		blockSet[elementsAdded] = section
		elementsAdded++
	}

	blocks := slack.Blocks{
		BlockSet: blockSet,
	}

	return &slack.Msg{
		Text:   mainText,
		BotID:  "A011834NXNY",
		Blocks: blocks,
	}
}
