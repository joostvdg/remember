package remember

type TVSeries struct {
	Title         string
	URL           string
	NotableActors []string
	Comment       string
	Producer      string
	Director      string
	Studio        string
	Distributor   string
	Seasons       int
	Labels        []string
}

func (t TVSeries) GetLabels() []string {
	return t.Labels
}

func (t TVSeries) ItemName() string {
	return t.Title
}

func (t TVSeries) GetURL() string {
	return t.URL
}

func (t TVSeries) Type() string {
	return "TVSeries"
}
