package remember

type Movie struct {
    Title string
    URL string
    Comment string
    Franchise string
    NotableActors []string
    Producer string
    Director string
    Studio string
    Series bool
    Labels []string
}

func (m Movie) GetLabels() []string {
    return m.Labels
}

func (m Movie) ItemName() string {
    return m.Title
}

func (m Movie) GetURL() string {
    return m.URL
}

func (m Movie) Type() string {
    return "Movie"
}

