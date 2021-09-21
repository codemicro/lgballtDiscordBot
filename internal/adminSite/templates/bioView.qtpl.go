// Code generated by qtc from "bioView.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

package templates

import "github.com/codemicro/lgballtDiscordBot/internal/db"

import "github.com/codemicro/lgballtDiscordBot/internal/config"

import "strings"

import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

type BioViewPage struct {
	BasePage
	Bio db.UserBio
}

func (p *BioViewPage) StreamBody(qw422016 *qt422016.Writer) {
	qw422016.N().S(`
<div>

    <a href="/bio">Return to search</a>
    <br>
    
    <h2>Bio for `)
	qw422016.E().S(p.Bio.UserId)
	qw422016.N().S(`
        `)
	if p.Bio.SysMemberID != "" {
		qw422016.N().S(`
            - `)
		qw422016.E().S(p.Bio.SysMemberID)
		qw422016.N().S(`
        `)
	}
	qw422016.N().S(`
    </h2>

    <p class="text-secondary"><i>Click a category name to edit that category's value.</i></p>

    <table class="table">
        <thead>
        <tr>
            <th scope="col">Category</th>
            <th scope="col">Value</th>
        </tr>
        </thead>
        <tbody>

            `)
	for _, fieldName := range config.BioFields {
		qw422016.N().S(`
                <tr>
                    <th scope="row"><a href="/bio/edit/field?user=`)
		qw422016.N().U(p.Bio.UserId)
		qw422016.N().S(`&member=`)
		qw422016.N().U(p.Bio.SysMemberID)
		qw422016.N().S(`&field=`)
		qw422016.N().U(fieldName)
		qw422016.N().S(`">`)
		qw422016.E().S(fieldName)
		qw422016.N().S(`</a></th>
                    <td>
                        `)
		for _, line := range strings.Split(p.Bio.BioData[fieldName], "\n") {
			qw422016.N().S(`
                            `)
			qw422016.E().S(line)
			qw422016.N().S(`<br>
                        `)
		}
		qw422016.N().S(`
                    </td>
                </tr>
            `)
	}
	qw422016.N().S(`

            <tr class="border-top">
                <th scope="row"><a href="/bio/edit/image?user=`)
	qw422016.N().U(p.Bio.UserId)
	qw422016.N().S(`&member=`)
	qw422016.N().U(p.Bio.SysMemberID)
	qw422016.N().S(`">Image</a></th>
                <td>
                    `)
	qw422016.E().S(p.Bio.ImageURL)
	qw422016.N().S(`
                    <br><img class="border pt-1" src="`)
	qw422016.E().S(p.Bio.ImageURL)
	qw422016.N().S(`" alt="Bio image preview" style="max-width: 500px">
                </td>
            </tr>

        </tbody>
    </table>

</div>
`)
}

func (p *BioViewPage) WriteBody(qq422016 qtio422016.Writer) {
	qw422016 := qt422016.AcquireWriter(qq422016)
	p.StreamBody(qw422016)
	qt422016.ReleaseWriter(qw422016)
}

func (p *BioViewPage) Body() string {
	qb422016 := qt422016.AcquireByteBuffer()
	p.WriteBody(qb422016)
	qs422016 := string(qb422016.B)
	qt422016.ReleaseByteBuffer(qb422016)
	return qs422016
}
