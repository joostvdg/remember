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

	log.Infof("Processing Slack List Command '%v'", command)
	switch command {
	case "view":
		return GetLists(slackCommand, log, store)
	case "new":
		return newList(slackCommand, log)
	default:
		return GetLists(slackCommand, log, store)
	}
}

func addBlock(blockSet []slack.Block, block slack.Block, counter int) int {
	blockSet[counter] = block
	return counter + 1
}

func newList(slackCommand SlackCommand, log *zap.SugaredLogger) *slack.Msg {
	log.Infof("Remember GetLists invoked for: %v (%v), in %v", slackCommand.UserName, slackCommand.UserId, slackCommand.Team)

	mainText := slack.NewTextBlockObject("plain_text", "Create New List", false, false)
	modalSection := slack.NewSectionBlock(mainText, nil, nil)

	textBlock := slack.NewTextBlockObject("plain_text", "Create New List", false, false)
	buttonText := slack.NewTextBlockObject("plain_text", "Create", false, false)
	button := slack.NewButtonBlockElement("new_list", "abc", buttonText)
	accessory := slack.NewAccessory(button)
	section := slack.SectionBlock{
		Type:      "section",
		Text:      textBlock,
		Accessory: accessory,
	}

	numberOfElements := 3
	elementsAdded := 0
	var blockSet []slack.Block
	blockSet = make([]slack.Block, numberOfElements, numberOfElements)
	elementsAdded = addBlock(blockSet, modalSection, elementsAdded)
	elementsAdded = addBlock(blockSet, slack.NewDividerBlock(), elementsAdded)
	elementsAdded = addBlock(blockSet, section, elementsAdded)

	blocks := slack.Blocks{
		BlockSet: blockSet,
	}
	return &slack.Msg{
		Text:   "New List",
		Blocks: blocks,
	}
}

func NewListView() *slack.ModalViewRequest {
	modalTile := slack.NewTextBlockObject("plain_text", ":memo:New List", true, false)
	modalSectionText := slack.NewTextBlockObject("plain_text", "Enter a name to create a new list", false, false)
	modalSection := slack.NewSectionBlock(modalSectionText, nil, nil)

	inputPlaceHolder := slack.NewTextBlockObject("plain_text", "please enter a name", false, false)
	inputElement := slack.NewPlainTextInputBlockElement(inputPlaceHolder, "name_input")
	inputLabel := slack.NewTextBlockObject("plain_text", "Name", false, false)
	nameInputBlock := slack.NewInputBlock("name_input_block", inputLabel, inputElement)

	numberOfElements := 3
	elementsAdded := 0
	var blockSet []slack.Block
	blockSet = make([]slack.Block, numberOfElements, numberOfElements)
	elementsAdded = addBlock(blockSet, modalSection, elementsAdded)
	elementsAdded = addBlock(blockSet, slack.NewDividerBlock(), elementsAdded)
	elementsAdded = addBlock(blockSet, nameInputBlock, elementsAdded)

	blocks := slack.Blocks{
		BlockSet: blockSet,
	}

	closeText := slack.NewTextBlockObject("plain_text", "Cancel", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Save", false, false)

	return &slack.ModalViewRequest{
		Type:            "modal",
		Title:           modalTile,
		Close:           closeText,
		Submit:          submitText,
		PrivateMetadata: "ssshhhhhhh",
		CallbackID:      "list_new",
		Blocks:          blocks,
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
		text := fmt.Sprintf("- *%v* - *%v* entries", list.Name, len(list.Entries))
		textBlock := slack.NewTextBlockObject("mrkdwn", text, false, false)
		buttonText := slack.NewTextBlockObject("plain_text", "Open", false, false)
		button := slack.NewButtonBlockElement("open_list", list.Id, buttonText)
		accessory := slack.NewAccessory(button)
		section := slack.SectionBlock{
			Type:      "section",
			Text:      textBlock,
			Accessory: accessory,
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
