package remember

type MediaItem interface {
	ItemName() string
	GetURL() string
	Type() string
	GetLabels() []string
}
