package gemgen

import "fmt"

type Gemtext struct {
	Buffer string
}

// AddUnformatted adds unformatted content to the Gemtext buffer
func (g *Gemtext) AddUnformatted(content string) {
	g.Buffer += fmt.Sprintf("%s\n", content)
}

// AddBlankLine adds a blank line to the Gemtext buffer
func (g *Gemtext) AddBlankLine() {
	g.Buffer += "\n"
}

// AddLink adds a link to the Gemtext buffer
func (g *Gemtext) AddLink(url, text string) {
	g.Buffer += fmt.Sprintf("=> %s %s\n", url, text)
}

// AddListLink adds as a link as an unordered list item to the Gemtext buffer
func (g *Gemtext) AddListLink(url string) {
	g.Buffer += fmt.Sprintf("=> * %s\n", url)
}

// AddHeading adds a heading to the Gemtext buffer
func (g *Gemtext) AddHeading(content string) {
	g.Buffer += fmt.Sprintf("# %s\n", content)
}

// AddSubHeadding adds a subheading to the Gemtext buffer
func (g *Gemtext) AddSubHeading(content string) {
	g.Buffer += fmt.Sprintf("## %s\n", content)
}

// AddSubSubHeading adds a sub-subheading to the Gemtext buffer
func (g *Gemtext) AddSubSubHeading(content string) {
	g.Buffer += fmt.Sprintf("### %s\n", content)
}

// AddListItem adds an unordered list item to the Gemtext buffer
func (g *Gemtext) AddListItem(content string) {
	g.Buffer += fmt.Sprintf("* %s\n", content)
}

// AddQuote adds a quote to the Gemtext buffer
func (g *Gemtext) AddQuote(content string) {
	g.Buffer += fmt.Sprintf("> %s\n", content)
}

// AddCodeBlock adds a codeblock, with monospace font, to the Gemtext buffer
func (g *Gemtext) AddCodeBlock(content string) {
	g.Buffer += fmt.Sprintf("```%s```\n", content)
}
