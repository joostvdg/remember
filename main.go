package main

import (
	"net/http"

    `github.com/joostvdg/remember/pkg/remember`
    "github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
    e.GET("/users/:id", getUser)
	e.Logger.Fatal(e.Start(":1323"))
}


// e.GET("/users/:id", getUser)
func getUser(c echo.Context) error {
    // User ID from path `users/:id`
    id := c.Param("id")
    user := &remember.User{
        Id:    "ABC" + id,
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
        Item:     book1,
        Order:    0,
        Comment:  "No Comment",
        Finished: false,
    }
    bookEntry2 := remember.MediaEntry{
        Item:     book2,
        Order:    1,
        Comment:  "No Comment",
        Finished: false,
    }
    var bookEntries = []remember.MediaEntry {bookEntry1, bookEntry2}

    bookList := remember.MediaList{
        Owner:        user.Id,
        Name:         "MyBookList",
        Description:  "Where I keep track of the books I want to read",
        Public:       true,
        Entries:      bookEntries,
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
        Title:         "The Expanse",
        URL:           "https://www.imdb.com/378382",
        Comment:       "Cool Space Drama",
        Producer:      "IDK",
        Director:      "IDK",
        Studio:        "IDK",
        Distributor:   "IDK",
        Seasons:        4,
        Labels:        []string {"space"},
    }
    movieEntry1 := remember.MediaEntry{
        Item:     movie1,
        Order:    0,
        Comment:  "No Comment",
        Finished: false,
    }
    movieEntry2 := remember.MediaEntry{
        Item:     movie2,
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
    progressions := []remember.Progression {s1, s2, s3, s4}
    tvSeriesEntry := remember.MediaEntry{
        Item:        tvSeries1,
        Order:       2,
        Comment:     "No Comment",
        Finished:    false,
        Progression: progressions,
    }
    var movieEntries = []remember.MediaEntry {movieEntry1, movieEntry2, tvSeriesEntry}
    movieList := remember.MediaList{
        Owner:        user.Id,
        Contributors: nil,
        Public:       true,
        Entries:      movieEntries,
    }

    var lists = []remember.MediaList { bookList, movieList}
    user.Lists = lists

    return c.JSON(http.StatusOK, user)
}
