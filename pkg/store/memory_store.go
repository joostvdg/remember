package store

import (
	"github.com/joostvdg/remember/pkg/remember"
)

type MemoryStore struct {
	Users []*remember.User
	Lists []*remember.MediaList
}
