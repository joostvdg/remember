package store

import (
	"errors"
	"github.com/joostvdg/remember/pkg/remember"
)

type MemoryStore struct {
	Users []*remember.User
	Lists []*remember.MediaList
}

func (m MemoryStore) AddUser(user remember.User) bool {
	if user.Id == "" {
		return false
	}
	_, userExists := m.FindUser(user.Id)
	if userExists {
		return false
	}
	m.Users = append(m.Users, &user)
	return true
}

func (m MemoryStore) FindUser(userId string) (remember.User, bool) {
	userIsFound := false
	var foundUser remember.User
	for _, user := range m.Users {
		if userId == user.Id {
			foundUser = *user
			userIsFound = true
		}
	}
	return foundUser, userIsFound
}

func (m MemoryStore) GetListForUser(userId string, listId string) (remember.MediaList, error) {
	var list = new(remember.MediaList)

	user, found := m.FindUser(userId)
	if !found {
		return *list, errors.New("Could not find user")
	}

	for _, listItem := range user.Lists {
		if listItem.Id == listId {
			list = listItem
		}
	}

	return *list, nil
}
