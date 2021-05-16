package route

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

type MessageRunFunc func(ctx *MessageContext) error
type ReactionRunFunc func(ctx *ReactionContext) error

type ReactionContext struct {
	*CommonContext
	Reaction *discordgo.MessageReaction
	Event    ReactionEvent
}

// onMessageCreate is a callback function to be used with a DiscordGo session that iterates through all registered
// commands and runs the first one that it finds that matches
func (b *Kit) onMessageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {

	// ignore self
	if usr, err := session.User("@me"); err != nil {
		b.ErrorHandler(err)
		return
	} else if usr.ID == message.Author.ID {
		return
	}

	if !b.AllowBots {
		if message.Author.Bot {
			return
		}
	}

	if !b.AllowDirectMessages {
		// no guild ID means this isn't a guild, therefore it's some kind of DM
		// discord please never add a third channel type
		if message.GuildID == "" {
			return
		}
	}

	// check if the message has a given prefix
	var trimmedContent string
	for _, prefix := range b.Prefixes {
		// slightly modified version of strings.HasPrefix
		if b.hasPrefix(message.Content, prefix) {
			trimmedContent = b.trimPrefix(message.Content, prefix)
			break
		}
	}

	ctx := &MessageContext{
		CommonContext: &CommonContext{
			Session: session,
			Kit:     b,
		},
		Message:   message,
	}

	ctx.Raw = message.Content
	b.tempMessageHandlerMux.RLock()
	for n, r := range b.tempMessageHandlerSet {
		x := r // this is a loop var and can change before the function is actually executed in the goroutine, which
		// would be bad
		go func() {
			y := *ctx // clone ctx
			err := x(&y)
			if err != nil {
				b.handleError(err, "error running temporary message handler", strconv.Itoa(n))
			}
		}()
	}
	b.tempMessageHandlerMux.RUnlock()

	if trimmedContent == "" {
		// no command? nothing for us to do
		if err := b.runMiddlewares(MiddlewareTriggerInvalid, ctx); err != nil {
			b.handleError(err, "error running middleware")
		}
		return
	}

	// find commands that match the trimmed command string
	var possibleCommands []*Command
	for _, cmd := range b.commandSet {
		if cmd.detectRegexp.MatchString(trimmedContent) {

			if len(possibleCommands) == 0 || cmd.AllowOverloading {
				possibleCommands = append(possibleCommands, cmd)
			}

		}
	}

	// if there are no matching commands, return
	if len(possibleCommands) == 0 {
		if err := b.runMiddlewares(MiddlewareTriggerInvalid, ctx); err != nil {
			b.handleError(err, "error running middleware")
		}
		return
	}

	// remove commands that don't match their restrictions
	{
		var badCommands []int

		for i, cmd := range possibleCommands {
			ok := true
			if cmd.Restrictions != nil {
				for _, rf := range cmd.Restrictions {
					rfOk, err := rf(session, message)
					if err != nil {
						b.handleError(err)
						return // TODO: could something else be done here?
					}
					ok = ok && (rfOk || b.DebugMode) // use debug mode to ignore restrictions
				}
			}

			if !ok {
				badCommands = append(badCommands, i)
			}
		}

		// actually remove from slice
		for i := len(badCommands) - 1; i >= 0; i -= 1 {

			n := badCommands[i]

			if n < len(possibleCommands)-1 {
				copy(possibleCommands[n:], possibleCommands[n+1:])
			}
			possibleCommands[len(possibleCommands)-1] = nil
			possibleCommands = possibleCommands[:len(possibleCommands)-1]
		}

	}

	// if no commands match restrictions
	if len(possibleCommands) == 0 {
		if err := session.MessageReactionAdd(message.ChannelID, message.ID, "⚠"); err != nil {
			b.handleError(err)
		}

		if err := b.runMiddlewares(MiddlewareTriggerInvalid, ctx); err != nil {
			b.handleError(err, "error running middleware")
		}
		return
	}

	// parse arguments for those commands
	var parseFailures []string
	var runCommand *Command
	var runArguments map[string]interface{}
	for _, cmd := range possibleCommands {
		// remove command text
		tcx := cmd.detectRegexp.Split(trimmedContent, -1)
		cmdTrimmedContent := strings.TrimSpace(tcx[1])

		// parse arguments
		var failMessage string
		argumentMap := make(map[string]interface{})
		if cmd.Arguments != nil {

			for _, arg := range cmd.Arguments {

				var val interface{}

				if len(cmdTrimmedContent) == 0 { // if there's nothing left to parse
					if arg.Default != nil { // if there's a default available

						dv, err := arg.Default(session, message)
						if err != nil {
							b.handleError(err)
							return
						}
						val = dv

					} else {
						failMessage = "argument missing"
					}

				} else { // otherwise parse from the available text
					var err error
					val, err = arg.Type.Parse(&cmdTrimmedContent)
					if err != nil {
						failMessage = err.Error()
					}
				}

				if failMessage != "" {
					x := fmt.Sprintf(" • If you were trying to run the **%s** command, what you entered into "+
						"`%s` was incorrect. The following error message was provided: **%s**",
						strings.ToLower(cmd.Name), arg.Name, failMessage)
					parseFailures = append(parseFailures, x)
					break // move onto the next command
				}

				argumentMap[arg.Name] = val
			}
		}

		if failMessage == "" {
			runCommand = cmd
			runArguments = argumentMap
			break // let's run a command!
		}
	}

	// if there was no command that parsed correctly
	if runCommand == nil {
		var extraText string
		if len(parseFailures) > 1 {
			extraText = " - the problem with it depends on what command you were trying to run."
		}
		errText := fmt.Sprintf("❌ **An error was encountered when running your command**%s\n\n", extraText) +
			strings.Join(parseFailures, "\n") + fmt.Sprintf("\n\nFor more information, run `%shelp`", b.Prefixes[0])
		_, err := session.ChannelMessageSend(message.ChannelID, errText)
		if err != nil {
			b.handleError(err)
		}
		if err := b.runMiddlewares(MiddlewareTriggerInvalid, ctx); err != nil {
			b.handleError(err, "error running middleware")
		}
		return
	}

	ctx.Raw = trimmedContent
	ctx.Arguments = runArguments
	y := *runCommand
	ctx.Command = &y

	if err := b.runMiddlewares(MiddlewareTriggerValid, ctx); err != nil {
		b.handleError(err, "error running middleware")
	}

	err := runCommand.Run(ctx)
	if err != nil {
		b.handleError(err, runCommand.Name)
		_ = b.Session.MessageReactionAdd(message.ChannelID, message.ID, "🚨")
	}

}

func (b *Kit) onReactionAdd(session *discordgo.Session, reaction *discordgo.MessageReactionAdd) {

	// ignore self
	if usr, err := session.User("@me"); err != nil {
		b.ErrorHandler(err)
		return
	} else if usr.ID == reaction.UserID {
		return
	}

	mCtx := ReactionContext{
		CommonContext: &CommonContext{
			Session: session,
			Kit:     b,
		},
		Reaction: reaction.MessageReaction,
		Event:    ReactionAdd,
	}

	if err := b.runMiddlewares(MiddlewareTriggerReactionAdd, &mCtx); err != nil {
		b.handleError(err, "error running middleware")
	}

	f := func(r *Reaction) {

		if r.Event != ReactionAdd {
			return
		}

		ctx := mCtx
		err := r.Run(&ctx)
		if err != nil {
			b.handleError(err, r.Name)
		}
	}

	for _, r := range b.reactionSet {
		f(r)
	}

	b.tempReactionsMux.RLock()
	for _, r := range b.tempReactionSet {
		f(r)
	}
	b.tempReactionsMux.RUnlock()

}

func (b *Kit) onReactionRemove(session *discordgo.Session, reaction *discordgo.MessageReactionRemove) {

	// ignore self
	if usr, err := session.User("@me"); err != nil {
		b.ErrorHandler(err)
		return
	} else if usr.ID == reaction.UserID {
		return
	}

	mCtx := ReactionContext{
		CommonContext: &CommonContext{
			Session: session,
			Kit:     b,
		},
		Reaction: reaction.MessageReaction,
		Event:    ReactionRemove,
	}

	if err := b.runMiddlewares(MiddlewareTriggerReactionRemove, &mCtx); err != nil {
		b.handleError(err, "error running middleware")
	}

	f := func(r *Reaction) {
		if r.Event != ReactionRemove {
			return
		}

		ctx := mCtx
		err := r.Run(&ctx)
		if err != nil {
			b.handleError(err, r.Name)
		}
	}

	for _, r := range b.reactionSet {
		f(r)
	}

	b.tempReactionsMux.RLock()
	for _, r := range b.tempReactionSet {
		f(r)
	}
	b.tempReactionsMux.RUnlock()

}
