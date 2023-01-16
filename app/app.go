package app

import (
	"fmt"
	"gemchan/app/handler"
	"gemchan/db"
	"log"
	"strings"

	gem "github.com/trebabcock/gemgen"

	"github.com/pitr/gig"
	"gorm.io/gorm"
)

var banner = `
  ____                     _                 
 / ___| ___ _ __ ___   ___| |__   __ _ _ __  
| |  _ / _ \ '_ ' _ \ / __| '_ \ / _' | '_ \ 
| |_| |  __/ | | | | | (__| | | | (_| | | | |
 \____|\___|_| |_| |_|\___|_| |_|\__,_|_| |_|

`

func baseURL(path string) string {
	return fmt.Sprintf("gemini://gemchan.space%s", path)
}

type App struct {
	Router *gig.Gig
	DB     *gorm.DB
}

func (a *App) Init() {
	a.DB = db.Init()
	a.Router = gig.Default()
	a.setRoutes()
}

func (a *App) setRoutes() {
	a.handle("/", a.index)
	a.handle("/board/:route", a.board)
	a.handle("/newpost/:board", a.newPost)
	a.handle("/post/:id", a.post)
	a.handle("/addComment/:id", a.addComment)
	a.handle("/replyComment/:id", a.replyComment)
	a.handle("/update", a.updatePage)
	a.handle("/overboard", a.overboard)
	a.handle("*", a.notFound)
}

func (a *App) handle(path string, f func(c gig.Context) error) {
	a.Router.Handle(path, f)
}

func (a *App) index(c gig.Context) error {
	buffer := gem.Gemtext{}
	buffer.AddCodeBlock(banner)
	buffer.AddHeading("Gemchan")
	buffer.AddUnformatted("Welcome to Gemchan, a textboard for Gemini")
	buffer.AddBlankLine()
	buffer.AddLink("/update", "Update")
	buffer.AddBlankLine()
	buffer.AddSubHeading("Boards")
	boards := handler.GetBoards()
	buffer.AddLink("/overboard", "Overboard")
	for _, b := range boards {
		buffer.AddLink(fmt.Sprintf("/board/%s", b.Route), fmt.Sprintf("/%s/", b.Route))
	}
	buffer.AddBlankLine()
	buffer.AddBlankLine()
	return c.Gemini(buffer.Buffer)
}

func (a *App) board(c gig.Context) error {
	r := c.Param("route")
	b, err := handler.GetBoard(r)
	if err != nil {
		return c.NoContent(gig.StatusRedirectTemporary, "/")
	}
	buffer := gem.Gemtext{}
	buffer.AddHeading(fmt.Sprintf("Welcome to /%s/", b.Route))
	buffer.AddUnformatted(b.Description)
	buffer.AddLink("/", "Home")
	buffer.AddBlankLine()
	buffer.AddLink(fmt.Sprintf("/newpost/%s", b.Route), "Create Post")
	buffer.AddSubHeading("Posts")
	buffer.AddBlankLine()
	for _, p := range handler.GetPostsFromBoard(a.DB, b.Route) {
		buffer.AddLink("/post/"+p.ID, p.ID)
		buffer.AddUnformatted(p.Time)
		c := handler.GetComments(a.DB, p.ID)
		if len(c) == 1 {
			buffer.AddUnformatted(fmt.Sprintf("%d comment", len(c)))
		} else {
			buffer.AddUnformatted(fmt.Sprintf("%d comments", len(c)))
		}
		buffer.AddQuote(p.Content)
		buffer.AddBlankLine()
	}
	return c.Gemini(buffer.Buffer)
}

func (a *App) newPost(c gig.Context) error {
	q, err := c.QueryString()
	if err != nil {
		log.Fatal(err)
	}
	if q == "" {
		return c.NoContent(gig.StatusInput, "Post Text")
	} else {
		id := handler.CreatePost(a.DB, q, c.Param("board"))
		return c.NoContent(gig.StatusRedirectTemporary, "/post/"+id)
	}
}

func (a *App) addComment(c gig.Context) error {
	q, err := c.QueryString()
	if err != nil {
		log.Fatal(err)
	}
	if q == "" {
		return c.NoContent(gig.StatusInput, "Add Comment")
	} else {
		post := handler.GetPost(a.DB, c.Param("id"))
		handler.AddComment(a.DB, q, post.ID)
		return c.NoContent(gig.StatusRedirectTemporary, "/post/"+post.ID)
	}
}

func (a *App) replyComment(c gig.Context) error {
	q, err := c.QueryString()
	if err != nil {
		log.Fatal(err)
	}
	ids := strings.Split(c.Param("id"), "&")
	if q == "" {
		return c.NoContent(gig.StatusInput, "Add Comment")
	} else {
		post := handler.GetPost(a.DB, ids[0])
		handler.AddCommentReply(a.DB, q, ids[1], post.ID)
		return c.NoContent(gig.StatusRedirectTemporary, "/post/"+post.ID)
	}
}

func (a *App) post(c gig.Context) error {
	post := handler.GetPost(a.DB, c.Param("id"))
	buffer := gem.Gemtext{}
	buffer.AddLink("/", "Home")
	buffer.AddLink(fmt.Sprintf("/board/%s", post.Board), fmt.Sprintf("/%s/", post.Board))
	buffer.AddHeading("Post")
	buffer.AddBlankLine()
	buffer.AddLink("/post/"+post.ID, post.ID)
	buffer.AddUnformatted(fmt.Sprintf("%s UTC", post.Time))
	buffer.AddBlankLine()
	buffer.AddQuote(post.Content)
	buffer.AddBlankLine()
	buffer.AddLink("/addComment/"+post.ID, "Add Comment")
	buffer.AddBlankLine()
	buffer.AddSubHeading("Comments")
	buffer.AddBlankLine()
	for _, c := range handler.GetComments(a.DB, post.ID) {
		buffer.AddUnformatted(c.ID)
		buffer.AddUnformatted(fmt.Sprintf("%s UTC", c.Time))
		if c.ReplyTo != "" {
			buffer.AddUnformatted(fmt.Sprintf("Reply to %s", c.ReplyTo))
		}
		buffer.AddQuote(c.Content)
		buffer.AddLink(fmt.Sprintf("/replyComment/%s&%s", post.ID, c.ID), "Reply")
		buffer.AddBlankLine()
	}
	return c.Gemini(buffer.Buffer)
}

func (a *App) notFound(c gig.Context) error {
	buffer := gem.Gemtext{}
	buffer.AddHeading("Not Found")
	buffer.AddUnformatted("You broke it.")
	buffer.AddBlankLine()
	buffer.AddLink("/", "Return Home")
	return c.Gemini(buffer.Buffer)
}

func (a *App) updatePage(c gig.Context) error {
	buffer := gem.Gemtext{}
	buffer.AddHeading("Update")
	buffer.AddBlankLine()
	buffer.AddUnformatted("So, you may have noticed that all posts are gone. That's my bad. I lost the SSH key for the server back in 2021, so Gemchan has just been stagnant since. I finally decided it was worth deleting all previous posts so that I could get some changes pushed.")
	buffer.AddBlankLine()
	buffer.AddUnformatted("I wrote a small gemini crawler in Go so that I could grab all posts and comments for archiving purposes. You can find the backup here:")
	buffer.AddLink("https://gist.github.com/trebabcock/4560d34fae1b62f305728ea8e57b1cf9", "Gemini Backup 01/16/2023")
	buffer.AddBlankLine()
	buffer.AddSubHeading("Changes")
	buffer.AddListItem("Now using shortuuid for more concise post and comment IDs.")
	buffer.AddListItem("Dates now have years. I guess I didn't expect gemchan to still exist almost two years later.")
	buffer.AddListItem("Comments are now stored in their own table, with the PostID as the foreign key. No idea why I thought storing them in an array in the Posts table made any sense at all.")
	buffer.AddListItem("Posts are now sorted by bump date, however posts never \"expire\" or get removed. I'll revisit this in the future if it becomes a problem.")
	buffer.AddListItem("Updated boards list to topics that more people will care about. I was really into theoretical physics and programming when I originally made Gemchan.")
	buffer.AddListItem("Added overboard. Contains posts from all boards, sorted by dump date.")
	buffer.AddBlankLine()
	buffer.AddUnformatted("Hopefully Gemchan can be a little less boring now, and possibly get some more traction. The only complaint I've heard is that it's dead, which is absolutely true.")
	buffer.AddBlankLine()
	buffer.AddUnformatted("Feel free to report bugs, or request features and boards in /meta/.")
	buffer.AddBlankLine()
	buffer.AddLink("/", "Return Home")
	return c.Gemini(buffer.Buffer)
}

func (a *App) overboard(c gig.Context) error {
	buffer := gem.Gemtext{}
	buffer.AddHeading("Overboard")
	buffer.AddUnformatted("Posts from all boards, sorted by bump date.")
	buffer.AddLink("/", "Home")
	buffer.AddBlankLine()
	buffer.AddSubHeading("Posts")
	buffer.AddBlankLine()
	for _, p := range handler.GetAllPosts(a.DB) {
		buffer.AddLink("/post/"+p.ID, p.ID)
		buffer.AddUnformatted(p.Time)
		c := handler.GetComments(a.DB, p.ID)
		if len(c) == 1 {
			buffer.AddUnformatted(fmt.Sprintf("%d comment", len(c)))
		} else {
			buffer.AddUnformatted(fmt.Sprintf("%d comments", len(c)))
		}
		buffer.AddQuote(p.Content)
		buffer.AddBlankLine()
	}
	return c.Gemini(buffer.Buffer)
}

func (a *App) Run(crt, key []byte) {
	log.Println("Server running at gemini://gemchan.space")
	log.Fatal(a.Router.Run(crt, key))
}
