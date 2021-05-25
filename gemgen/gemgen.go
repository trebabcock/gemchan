package gemgen

import "fmt"

type Gemtext struct {
	Content string
}

func (g *Gemtext) AddContent(content string) {
	g.Content += fmt.Sprintf("%s\n", content)
}

func (g *Gemtext) AddBlankLine() {
	g.Content += fmt.Sprintf("\n")
}

func (g *Gemtext) AddLink(url, text string) {
	g.Content += fmt.Sprintf("=> %s %s\n", url, text)
}

func (g *Gemtext) AddListLink(url string) {
	g.Content += fmt.Sprintf("=> * %s\n", url)
}

func (g *Gemtext) AddHeading(content string) {
	g.Content += fmt.Sprintf("# %s\n", content)
}

func (g *Gemtext) AddSubHeading(content string) {
	g.Content += fmt.Sprintf("## %s\n", content)
}

func (g *Gemtext) AddSubSubHeading(content string) {
	g.Content += fmt.Sprintf("### %s\n", content)
}

func (g *Gemtext) AddListItem(content string) {
	g.Content += fmt.Sprintf("* %s\n", content)
}

func (g *Gemtext) AddQuote(content string) {
	g.Content += fmt.Sprintf("> %s\n", content)
}

func (g *Gemtext) AddCodeBlock(content string) {
	g.Content += fmt.Sprintf("```%s```\n", content)
}
