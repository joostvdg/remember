package slack

import (
	"encoding/json"
	"github.com/joostvdg/remember/pkg/context"
	"github.com/labstack/echo/v4"
	"github.com/slack-go/slack"
	"io"

	"io/ioutil"
	"net/http"
)

var signingSecret = ""

func SlackInteractiveHandler(c echo.Context) error {

	return c.JSON(http.StatusOK,"OK")
}

func SlackHoppaHandler(c echo.Context) error {
	cc := c.(*context.CustomContext)
	cc.Log.Info("Received Slack Slash Command")
	params := &slack.Msg{
		Text:         "Hello World!",
		ResponseType: "ephemeral",
	}
	cc.Log.Info("Returning Valid Response (I believe)")
	return c.JSON(http.StatusOK, params)
}

func SlackHandler(c echo.Context) error {
	cc := c.(*context.CustomContext)
	cc.Log.Info("Received Slack Slash Command")

	request := c.Request()
	verifier, err := slack.NewSecretsVerifier(request.Header, signingSecret)
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
			BotID: "A011834NXNY",
		}
		cc.Log.Info("Returning Valid Response (I believe)")
		return c.JSON(http.StatusOK, params)
	case "list":
		response := HandleList(slackCommand, cc.Log, cc.MemoryStore)
		jsonResponse, _  := json.Marshal(response)
		cc.Log.Infof("%v\n", string(jsonResponse))
		return c.JSON(http.StatusOK, response)
	default:
		cc.Log.Warnf("Found error: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Please provide valid request")
	}
}
