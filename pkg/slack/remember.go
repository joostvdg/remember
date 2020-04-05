package slack

import (
	"fmt"
	"github.com/joostvdg/remember/pkg/remember"
	"github.com/joostvdg/remember/pkg/store"
	"go.uber.org/zap"
)

func handleRemember(slackCommand SlackCommand, log *zap.SugaredLogger, store store.MemoryStore) string {
	log.Infof("Remember invoked for: %v (%v), in %v", slackCommand.UserName, slackCommand.UserId, slackCommand.Team)

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

	flatList := ""
	for _, list := range userLists {
		listLine := fmt.Sprintf("- %v (id: %v)\n", list.Name, list.Id)
		flatList = flatList + listLine
	}
	// - Detective Chimp\n- Bouncing Boy\n- Aqualad
	return flatList
}
