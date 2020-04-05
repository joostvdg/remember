package slack

type SlackCommand struct {
	Command       string
	Text          string
	TokenizedText []string
	UserId        string
	UserName      string
	Team          string
}
