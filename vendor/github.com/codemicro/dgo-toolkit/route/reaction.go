package route

// ReactionEvent represents a the type of reaction event that has been triggered
type ReactionEvent uint8

const (
	ReactionAdd ReactionEvent = iota
	ReactionRemove
)

// Reaction represents a new reaction handler
type Reaction struct {
	Name  string
	Run   ReactionRunFunc
	Event ReactionEvent
}
