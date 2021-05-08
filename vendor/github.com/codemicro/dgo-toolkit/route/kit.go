package route

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"regexp"
	"strings"
	"sync"
)

// Kit is the core model for command parsing and routing
type Kit struct {
	Session                *discordgo.Session
	ErrorHandler           func(error)
	Prefixes               []string
	IsCaseSensitive        bool
	DebugMode              bool
	AllowBots              bool
	DefaultAllowedMentions discordgo.MessageAllowedMentions
	AllowDirectMessages    bool
	UserErrorFunc          func(string) string

	commandSet            []*Command
	tempMessageHandlerSet map[int]MessageRunFunc
	tempMessageHandlerMux *sync.RWMutex

	reactionSet      []*Reaction
	tempReactionSet  map[int]*Reaction
	tempReactionsMux *sync.RWMutex

	middlewareSet []*Middleware
}

// NewKit creates a new Kit instance
func NewKit(session *discordgo.Session, prefixes []string) *Kit {
	return &Kit{
		Session:  session,
		Prefixes: prefixes,
		UserErrorFunc: func(s string) string {
			return "❌ **Error:** " + s
		},
		tempReactionSet:  make(map[int]*Reaction),
		tempReactionsMux: new(sync.RWMutex),
		tempMessageHandlerSet:  make(map[int]MessageRunFunc),
		tempMessageHandlerMux: new(sync.RWMutex),
	}
}

// HandleError is the internal function used to handle an error that accounts for *kit.ErrorHandler being nil
func (b *Kit) handleError(e error, i ...string) {
	if len(i) >= 1 {
		e = fmt.Errorf("%s - %s", strings.Join(i, " "), e.Error())
	}
	if b.ErrorHandler == nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", e.Error())
	} else {
		b.ErrorHandler(e)
	}
}

// AddCommand adds commands to the command set for this instance of Kit. Command overloading is supported for commands
// that have AllowOverloading set true. If more than one command is matched by an incoming message, the command that was
// added to the Kit first is run.
func (b *Kit) AddCommand(commands ...*Command) {

	for _, c := range commands {
		var rx []string
		for _, x := range c.CommandText {
			rx = append(rx, regexp.QuoteMeta(x))
		}

		var isc string
		if !b.IsCaseSensitive {
			isc = `(?i)`
		}

		c.detectRegexp = regexp.MustCompile(isc + `^` + strings.Join(rx, ` +`))

		b.commandSet = append(b.commandSet, c)
	}

}

// AddTemporaryMessageHandler creates registers a temporary message handler and returns the created handler's ID.
func (b *Kit) AddTemporaryMessageHandler(handler MessageRunFunc) int {
	b.tempMessageHandlerMux.Lock()
	n := len(b.tempMessageHandlerSet)
	b.tempMessageHandlerSet[n] = handler
	b.tempMessageHandlerMux.Unlock()
	return n
}

// RemoveTemporaryMessageHandler removes a temporary message handler based on the provided ID. If the ID does not
// exist, RemoveTemporaryMessageHandler is a no-op.
func (b *Kit) RemoveTemporaryMessageHandler(id int) {
	b.tempMessageHandlerMux.Lock()
	delete(b.tempMessageHandlerSet, id)
	b.tempMessageHandlerMux.Unlock()
}

// AddReaction adds a reaction create handler to the reaction set for this instance of Kit
func (b *Kit) AddReaction(reactions ...*Reaction) {
	b.reactionSet = append(b.reactionSet, reactions...)
}

// AddTemporaryReaction creates registers a temporary reaction handler and returns the created handler's ID.
func (b *Kit) AddTemporaryReaction(reaction *Reaction) int {
	b.tempReactionsMux.Lock()
	n := len(b.tempReactionSet)
	b.tempReactionSet[n] = reaction
	b.tempReactionsMux.Unlock()
	return n
}

// RemoveTemporaryReaction removes a temporary reaction handler based on the provided ID. If the ID does not exist,
// RemoveTemporaryReaction is a no-op.
func (b *Kit) RemoveTemporaryReaction(id int) {
	b.tempReactionsMux.Lock()
	delete(b.tempReactionSet, id)
	b.tempReactionsMux.Unlock()
}

// AddMiddleware adds a middleware to the middleware set for this instance of Kit
func (b *Kit) AddMiddleware(middlewares ...*Middleware) {
	b.middlewareSet = append(b.middlewareSet, middlewares...)
}

func (b *Kit) CreateHandlers() {
	if b.commandSet != nil && len(b.commandSet) > 0 {
		b.Session.AddHandler(b.onMessageCreate)
	}
	if b.reactionSet != nil && len(b.reactionSet) > 0 {
		b.Session.AddHandler(b.onReactionAdd)
		b.Session.AddHandler(b.onReactionRemove)
	}
}

// caseCompare compares two strings either with or without case sensitivity depending on the value set in the parent Kit
func (b *Kit) caseCompare(x, y string) bool {
	if b.IsCaseSensitive {
		return x == y
	}
	return strings.EqualFold(x, y)
}

// hasPrefix is an implementation of strings.HasPrefix that uses caseCompare
func (b *Kit) hasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && b.caseCompare(s[0:len(prefix)], prefix)
}

// trimPrefix is an implementation of strings.TrimPrefix that uses caseCompare
func (b *Kit) trimPrefix(s, prefix string) string {
	if b.hasPrefix(s, prefix) {
		return s[len(prefix):]
	}
	return s

}

func (b *Kit) GetNums() (int, int) {
	return len(b.commandSet), len(b.reactionSet) + len(b.tempReactionSet)
}
