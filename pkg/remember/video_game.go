package remember

type VideoGame struct {
    Title string
    URL string
    Comment string
    Publisher string
    Studio string
    Franchise string
    Multiplayer bool
    Coop bool
    MinPlayers int
    MaxPlayers int
    Platforms []string
    Expansion bool
    Labels []string
}

func (v VideoGame) GetLabels() []string {
    return v.Labels
}

func (v VideoGame) ItemName() string {
    return v.Title
}

func (v VideoGame) GetURL() string {
    return v.URL
}

func (v VideoGame) Type() string {
    return "VideoGame"
}
