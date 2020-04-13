package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joostvdg/remember/pkg/remember"
	"github.com/slack-go/slack"
	"net/http"
)

func createSlackMessageForResponse(list remember.MediaList, botId string) *slack.Msg {
	numberOfElements := 2 + len(list.Entries)
	elementsAdded := 0
	var blockSet []slack.Block
	blockSet = make([]slack.Block, numberOfElements, numberOfElements)

	mainText := fmt.Sprintf("List: %v", list.Name)
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

	fmt.Printf("Processing %v entries\n", len(list.Entries))
	for _, entry := range list.Entries {
		done := "[ ]"
		if entry.Finished {
			done = "[*]"
		}
		text := fmt.Sprintf("*<%v|%v>* (%v) %v", entry.Item.GetURL(), entry.Item.ItemName(), entry.Item.Type(), done)
		textBlock := slack.NewTextBlockObject("mrkdwn", text, false, false)
		section := slack.SectionBlock{
			Type: "section",
			Text: textBlock,
		}
		blockSet[elementsAdded] = section
		elementsAdded++
	}

	blocks := slack.Blocks{
		BlockSet: blockSet,
	}

	return &slack.Msg{
		Text:   mainText,
		BotID:  botId,
		Blocks: blocks,
	}
}

func SendListResponse(responseUrl string, list remember.MediaList, apiToken string, botId string) (*http.Response, error) {
	slackMessage := createSlackMessageForResponse(list, botId)

	slackMessageJson, err := json.Marshal(slackMessage)
	if err != nil {
		fmt.Println("Could not parse SendListResponse")
		return nil, err
	}
	req, err := http.NewRequest("POST", responseUrl, bytes.NewBuffer(slackMessageJson))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	client := &http.Client{}
	reponse, err := client.Do(req)
	return reponse, err
}

func CreateViewSubmitResponse(update *slack.ViewState) *slack.ViewSubmissionResponse {
	fmt.Printf("Found view submission, name input value: %v\n", update.Values["name_input_block"]["name_input"].Value)
	modalTile := slack.NewTextBlockObject("plain_text", ":memo:New List", true, false)

	updateText := fmt.Sprintf("Created a new list with name *%v*", update.Values["name_input_block"]["name_input"].Value)
	modalSectionText := slack.NewTextBlockObject("plain_text", updateText, false, false)
	modalSection := slack.NewSectionBlock(modalSectionText, nil, nil)

	numberOfElements := 1
	elementsAdded := 0
	var blockSet []slack.Block
	blockSet = make([]slack.Block, numberOfElements, numberOfElements)
	elementsAdded = addBlock(blockSet, modalSection, elementsAdded)

	blocks := slack.Blocks{
		BlockSet: blockSet,
	}

	view := &slack.ModalViewRequest{
		Type:            "modal",
		Title:           modalTile,
		Blocks:          blocks,
	}
	return slack.NewUpdateViewSubmissionResponse(view)
}
