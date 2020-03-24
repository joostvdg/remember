package remember

type MediaList struct {
	Id           string
	Owner        string   // user Id
	Contributors []string // user Ids
	Name         string
	Description  string
	Public       bool
	Entries      []MediaEntry
}
