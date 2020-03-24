package remember

type BoardGame struct {
	Title      string
	URL        string
	Comment    string
	Publisher  string
	MinPlayers int
	MaxPlayers int
	Expansion  bool
	Labels     []string
}

func (b BoardGame) GetLabels() []string {
	return b.Labels
}

func (b BoardGame) ItemName() string {
	return b.Title
}

func (b BoardGame) GetURL() string {
	return b.URL
}

func (b BoardGame) Type() string {
	return "BoardGame"
}
