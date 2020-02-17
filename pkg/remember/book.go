package remember

type Book struct {
	Title     string
	URL       string
	ISBN      string
	Author    string
	Publisher string
	Comment   string
	Series    bool
	Labels 	  []string
}

func (b Book) GetLabels() []string {
	return b.Labels
}

func (b Book) ItemName() string {
	return b.Title
}

func (b Book) GetURL() string {
	return b.URL
}

func (b Book) Type() string {
	return "Book"
}
