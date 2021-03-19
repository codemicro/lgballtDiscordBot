package pressf

import "github.com/codemicro/dgo-toolkit/route"

func (pf *PressF) Reaction(ctx *route.ReactionContext) error {
	pf.mux.RLock()
	apf, found := pf.active[ctx.Reaction.MessageID]
	pf.mux.RUnlock()
	if found && ctx.Reaction.Emoji.Name == regionalFEmoji {
		apf.Trigger <- ctx.Reaction
	}
	return nil
}
