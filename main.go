package main

import (
	"github.com/joostvdg/remember/pkg/slack"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"strconv"

	"github.com/joostvdg/remember/pkg/context"
	"github.com/joostvdg/remember/pkg/oauth"
	"github.com/joostvdg/remember/pkg/remember"
	"github.com/joostvdg/remember/pkg/store"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func main() {
	memoryStore := initMemoryStore()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Use(middleware.Logger())
	//e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
	//	TokenLookup: "header:X-XSRF-TOKEN",
	//}))

	e.GET("/auth/google/login", oauth.OauthGoogleLogin)
	// ITS A GET?
	//e.POST("/auth/google/callback", oauth.OauthGoogleCallback)
	//e.PUT("/auth/google/callback", oauth.OauthGoogleCallback)
	e.GET("/auth/google/callback", oauth.OauthGoogleCallback)

	e.POST("/hoppa", slack.SlackHoppaHandler)
	e.POST("/slack", slack.SlackInteractiveHandler)
	e.POST("/rmb", slack.SlackHandler)
	e.GET("/users/:id", getUser)
	e.GET("/users/:id/lists/:listId", getUserList)
	e.PUT("/users/:id/lists", newList)
	e.POST("/users/:id/lists/:listId/entries/:entryId/progression/:progressionKey/:current", updateProgression)

	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &context.CustomContext{
				c,
				memoryStore,
				sugar,
			}
			return h(cc)
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "1323"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

func newList(c echo.Context) (err error) {
	cc := c.(*context.CustomContext)
	// User ID from path `users/:id`
	id := cc.Param("id")

	userIsFound := false
	var foundUser remember.User
	for _, user := range cc.MemoryStore.Users {
		if id == user.Id {
			foundUser = *user
			userIsFound = true
		}
	}
	if !userIsFound {
		return c.String(http.StatusNotFound, "No user by that name")
	}

	list := new(remember.MediaList)
	if err = c.Bind(list); err != nil {
		return c.String(http.StatusBadRequest, "Not a valid MediaList")
	}
	foundUser.AddList(list)
	return c.String(http.StatusCreated, "Created new MediaList")
}

// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
	cc := c.(*context.CustomContext)
	// User ID from path `users/:id`
	id := cc.Param("id")

	for _, user := range cc.MemoryStore.Users {
		if id == user.Id {
			return c.JSON(http.StatusOK, user)
		}
	}

	return c.String(http.StatusNotFound, "No user by that name")
}

func getUserList(c echo.Context) error {
	cc := c.(*context.CustomContext)
	// User ID from path `users/:id`
	id := cc.Param("id")
	listId := cc.Param("listId")

	userFound := false
	for _, user := range cc.MemoryStore.Users {
		if id == user.Id {
			userFound = true
			for _, list := range user.Lists {
				if list.Id == listId {
					return c.JSON(http.StatusOK, list)
				}
			}
		}
	}

	if !userFound {
		return c.String(http.StatusNotFound, "No user by that name")
	}
	return c.String(http.StatusNotFound, "No list found by that id for user")

}

// updateProgression updates the progression of a users's entry
// /users/:id/lists/:listId/entries/:entryId/progression/:progressionKey/:current
// for example, I've just finished watching Ep5 of Season 4 of The Expanse
// /users/ABC/lists/MovieListOne/entries/EntryOne/progression/3/5
func updateProgression(c echo.Context) error {
	cc := c.(*context.CustomContext)
	id := cc.Param("id")

	userIsFound := false
	var foundUser remember.User
	for _, user := range cc.MemoryStore.Users {
		if id == user.Id {
			foundUser = *user
			userIsFound = true
		}
	}
	if !userIsFound {
		return c.String(http.StatusNotFound, "No user by that name")
	}
	listId := cc.Param("listId")
	entryId := cc.Param("entryId")
	progressionKeyString := cc.Param("progressionKey")
	current := cc.Param("current")

	progressionKey, progressionErr := strconv.Atoi(progressionKeyString)
	currentUpdate, currentErr := strconv.Atoi(current)

	if progressionErr != nil || currentErr != nil {
		return c.String(http.StatusBadRequest, "Supplied ids were not correct")
	}

	found := false
	updated := false
	for _, list := range foundUser.Lists {
		if list.Id == listId {
			for _, entry := range list.Entries {
				if entry.Id == entryId && entry.Progression != nil && len(entry.Progression) >= progressionKey {
					found = true
					entry.Progression[progressionKey].Current = currentUpdate
					updated = true
				}
			}
		}
	}

	if !found {
		return c.String(http.StatusNotFound, "Could not find the progression to update")
	}

	if updated {
		return c.String(http.StatusAccepted, "Update processed")
	} else {
		return c.String(http.StatusNoContent, "Nothing to update")
	}
}

func initMemoryStore() store.MemoryStore {
	user := &remember.User{
		Id:    "U3N251DH9",
		Email: "user@example.com",
		Name:  "User",
		Lists: nil,
	}

	var book1 remember.MediaItem
	book1 = remember.Book{
		Title:     "A Title",
		URL:       "https://something.com/books/123",
		ISBN:      "483989435879",
		Author:    "Pietje",
		Publisher: "Heintje",
		Comment:   "No",
		Series:    false,
	}
	var book2 remember.MediaItem
	book2 = remember.Book{
		Title:     "Some Other Title",
		URL:       "https://something.com/books/124",
		ISBN:      "483989435879",
		Author:    "Pietje",
		Publisher: "Heintje",
		Comment:   "No",
		Series:    true,
	}

	bookEntry1 := remember.MediaEntry{
		Id:       "1",
		Item:     book1,
		Order:    0,
		Comment:  "No Comment",
		Finished: false,
	}
	bookEntry2 := remember.MediaEntry{
		Id:       "2",
		Item:     book2,
		Order:    1,
		Comment:  "No Comment",
		Finished: false,
	}
	var bookEntries = []remember.MediaEntry{bookEntry1, bookEntry2}

	bookList := remember.MediaList{
		Id:          "2",
		Owner:       user.Id,
		Name:        "MyBookList",
		Description: "Where I keep track of the books I want to read",
		Public:      true,
		Entries:     bookEntries,
	}

	movie1 := remember.Movie{
		Title:         "Parasite",
		URL:           "https://www.imdb.com/378383",
		Comment:       "Oscar winner",
		NotableActors: nil,
		Producer:      "bong",
		Director:      "bong",
		Studio:        "?",
		Series:        false,
	}
	movie2 := remember.Movie{
		Title:         "1917",
		URL:           "https://www.imdb.com/378382",
		Comment:       "Oscar winner",
		NotableActors: nil,
		Producer:      "?",
		Director:      "?",
		Studio:        "?",
		Series:        false,
	}
	tvSeries1 := remember.TVSeries{
		Title:       "The Expanse",
		URL:         "https://www.imdb.com/378382",
		Comment:     "Cool Space Drama",
		Producer:    "IDK",
		Director:    "IDK",
		Studio:      "IDK",
		Distributor: "IDK",
		Seasons:     4,
		Labels:      []string{"space"},
	}
	movieEntry1 := remember.MediaEntry{
		Item:     movie1,
		Id:       "1",
		Order:    0,
		Comment:  "No Comment",
		Finished: false,
	}
	movieEntry2 := remember.MediaEntry{
		Item:     movie2,
		Id:       "2",
		Order:    1,
		Comment:  "No Comment",
		Finished: false,
	}
	s1 := remember.Progression{
		Min:     0,
		Max:     10,
		Current: 10,
	}
	s2 := remember.Progression{
		Min:     0,
		Max:     10,
		Current: 10,
	}
	s3 := remember.Progression{
		Min:     0,
		Max:     10,
		Current: 10,
	}
	s4 := remember.Progression{
		Min:     0,
		Max:     10,
		Current: 4,
	}
	progressions := []remember.Progression{s1, s2, s3, s4}
	tvSeriesEntry := remember.MediaEntry{
		Id:          "3",
		Item:        tvSeries1,
		Order:       2,
		Comment:     "No Comment",
		Finished:    false,
		Progression: progressions,
	}
	var movieEntries = []remember.MediaEntry{movieEntry1, movieEntry2, tvSeriesEntry}
	movieList := remember.MediaList{
		Id:           "1",
		Owner:        user.Id,
		Name:         "Movies",
		Contributors: nil,
		Public:       true,
		Entries:      movieEntries,
	}

	var lists = []*remember.MediaList{&bookList, &movieList}
	user.Lists = lists

	store := store.MemoryStore{
		Users: []*remember.User{user},
		Lists: []*remember.MediaList{&bookList, &movieList},
	}
	return store
}
