// Code generated by qtc from "index.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

package templates

import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

type IndexPage struct {
	BasePage
	DiscordLoginURL string
}

func (p *IndexPage) StreamBody(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
<div>
    <a href="`)
	qw422016.E().S(p.DiscordLoginURL)
	qw422016.N().S(`">Click here to login with Discord</a>
</div>
`)
}

func (p *IndexPage) WriteBody(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamBody(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *IndexPage) Body() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteBody(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}
