package remember

import (
	"fmt"
	// TODO: enable when back online
	//"github.com/google/uuid"
	"math/rand"
	"time"
)

type User struct {
	Id    string
	Email string
	Name  string
	Lists []*MediaList
}

func (u *User) AddList(list *MediaList) {
	// TODO: enable when back online
	//list.Id = uuid.New().String()
	newId := "A"
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10; i++ {
		number := rand.Intn(100)
		newId += fmt.Sprintf("%v", number)
	}

	list.Id = newId
	list.Owner = u.Id
	u.Lists = append(u.Lists, list)
}
