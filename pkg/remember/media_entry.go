package remember

type MediaEntry struct {
	Item     MediaItem
	Order    int
	Comment  string
	Finished bool
	Progression []Progression
}

type Progression struct {
	Min int
	Max int
	Current int
}