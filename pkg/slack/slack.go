package slack

import (
	"encoding/json"
	"fmt"
	"github.com/joostvdg/remember/pkg/context"
	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
	"io"
	"io/ioutil"
	"net/http"
)

func InteractiveHandler(c echo.Context) error {
	cc := c.(*context.CustomContext)
	rawPayload := c.FormValue("payload")
	cc.Log.Info("Received Slack Interaction Request")
	cc.Log.Infof("Got a callback, JSON: %v", rawPayload)

	request := c.Request()
	verifier, err := slack.NewSecretsVerifier(request.Header, cc.SigningSecret)
	if err != nil {
		cc.Log.Warn("Could not parse request")
		return echo.NewHTTPError(http.StatusInternalServerError, "Please provide valid request")
	}

	ioutil.NopCloser(io.TeeReader(request.Body, &verifier))
	if err = verifier.Ensure(); err != nil {
		cc.Log.Warnf("Found error: %v", err)
		// return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid signing")
	}
	var payload slack.InteractionCallback
	err = json.Unmarshal([]byte(rawPayload), &payload)
	if err != nil {
		errorMessage := fmt.Sprintf("Could not parse interactive payload: %v", err)
		cc.Log.Warnf(errorMessage)
		return c.JSON(http.StatusInternalServerError, errorMessage)
	}

	var blockAction *slack.BlockAction
	blockActions := payload.ActionCallback.BlockActions
	if blockActions != nil {
		blockAction = blockActions[0]
	}

	actionId := ""
	actionValue := ""
	if blockAction != nil {
		actionValue = blockAction.Value
		actionId = blockAction.ActionID
	}

	callbackId := payload.View.CallbackID
	if callbackId != "" && payload.Type == "view_submission" {
		actionId = callbackId
	}

	cc.Log.Infof("Processing CallBack: [ userId: %v, triggerId: %v, token: %v, callbackId: %v, action: %v, value: %v, responseUrl: %v]",
		payload.User.ID,
		payload.TriggerID,
		payload.Token,
		callbackId,
		actionId,
		actionValue,
		payload.ResponseURL)

	// Send Response after User Interaction
	// See Slack API: https://api.slack.com/interactivity/handling#acknowledgment_response
	switch actionId {
	case "list_new":
		cc.Log.Info("Found List update, sending view update")
		viewSubmitResponse := CreateViewSubmitResponse(payload.View.State)

		viewSubmitResponseJson, _ := json.Marshal(viewSubmitResponse)
		cc.Log.Infof("View Submit Response: %v", string(viewSubmitResponseJson))

		return c.JSON(http.StatusOK, viewSubmitResponse)
	case "new_list":
		view := NewListView()
		err := SendCreateView(payload.TriggerID, view, cc.APIToken)
		if err != nil {
			warnMessage := fmt.Sprintf("Could not write view response to Slack API: %v", err)
			cc.Log.Warnf(warnMessage)
			return c.JSON(http.StatusInternalServerError, warnMessage)
		}
	case "open_list":
		list, err := cc.MemoryStore.GetListForUser(payload.User.ID, actionValue)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "cannot retrieve list for user")
		}
		cc.Log.Infof("Found List %v with %d entries", list.Name, len(list.Entries))
		response, err := SendListResponse(payload.ResponseURL, list, cc.APIToken, cc.BotID)
		if response != nil {
			cc.Log.Infof("Response on Sending Slack Message: %v [%v]", response.Status, response.Header)
			if response.Status != "200 Ok" {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					cc.Log.Warnf("Error reading body: %v", err)
				} else {
					cc.Log.Infof("Response Body: %v", string(body))
				}
			}
		} else {
			cc.Log.Info("Response from Slack was empty: %v", response)
		}
		if err != nil {
			warnMessage := fmt.Sprintf("Could not write message response to Slack API: %v", err)
			cc.Log.Warnf(warnMessage)
			return c.JSON(http.StatusInternalServerError, warnMessage)
		}
	}

	cc.Log.Infof("Send Slack Message response, now sending Confirmation Response")
	return c.JSON(http.StatusOK, "OK")
}

func HoppaHandler(c echo.Context) error {
	cc := c.(*context.CustomContext)
	cc.Log.Info("Received Slack Slash Command")
	params := &slack.Msg{
		Text:         "Hello World!",
		ResponseType: "ephemeral",
	}
	cc.Log.Info("Returning Valid Response (I believe)")
	return c.JSON(http.StatusOK, params)
}

func DefaultHandler(c echo.Context) error {
	cc := c.(*context.CustomContext)
	cc.Log.Info("Received Slack Slash Command")

	request := c.Request()
	verifier, err := slack.NewSecretsVerifier(request.Header, cc.SigningSecret)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Please provide valid request")
	}

	request.Body = ioutil.NopCloser(io.TeeReader(request.Body, &verifier))
	slashCommand, err := slack.SlashCommandParse(request)
	if err != nil {
		cc.Log.Warnf("Found error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Please provide valid request")
	}

	if err = verifier.Ensure(); err != nil {
		cc.Log.Warnf("Found error: %v", err)
		return echo.NewHTTPError(http.StatusUnauthorized, "Please provide valid signing")
	}

	command, tokens := tokenizeSlackCommand(slashCommand.Text)
	slackCommand := SlackCommand{
		Command:       command,
		Text:          slashCommand.Text,
		TokenizedText: tokens,
		UserId:        slashCommand.UserID,
		UserName:      slashCommand.UserName,
		Team:          slashCommand.TeamDomain,
	}

	switch command {
	case "test":
		response := handleRemember(slackCommand, cc.Log, cc.MemoryStore)
		params := &slack.Msg{
			Text:  response,
			BotID: cc.BotID,
		}
		cc.Log.Info("Returning Valid Response (I believe)")
		return c.JSON(http.StatusOK, params)
	case "list":
		msg := HandleList(slackCommand, cc.Log, cc.MemoryStore)
		if msg != nil {
			cc.Log.Info("Returning a Slack Message")
			return c.JSON(http.StatusOK, msg)
		} else {
			cc.Log.Warnf("Did not manage to create a slack message")
			return echo.NewHTTPError(http.StatusInternalServerError, "Could not build response")
		}
	default:
		cc.Log.Warnf("Found error: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Please provide valid request")
	}
}
