package app

import (
	"fmt"
	"gem/app/handler"
	"gem/app/model"
	gem "gem/gemgen"
	"log"
	"strings"

	"github.com/pitr/gig"
	"gorm.io/driver/sqlite"
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
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Migrating database...")
	a.DB = model.DBMigrate(db)
	log.Println("Database migrated")
	a.Router = gig.Default()
	log.Println("Defining routes...")
	a.setRoutes()
	log.Println("Routes defined")
}

func (a *App) setRoutes() {
	a.handle("/", a.index)
	a.handle("/board/:route", a.board)
	a.handle("/newpost/:board", a.newPost)
	a.handle("/post/:id", a.post)
	a.handle("/addComment/:id", a.addComment)
	a.handle("/replyComment/:id", a.replyComment)
	a.handle("*", a.notFound)
}

func (a *App) handle(path string, f func(c gig.Context) error) {
	a.Router.Handle(path, f)
}

func (a *App) index(c gig.Context) error {
	buffer := gem.Gemtext{}
	buffer.AddCodeBlock(banner)
	buffer.AddHeading("Gemchan")
	buffer.AddUnformatted("Welcome to Gemchan, a textboard for Gemini!")
	buffer.AddBlankLine()
	buffer.AddSubHeading("Boards")
	boards := handler.GetBoards()
	for _, b := range boards {
		buffer.AddLink(baseURL(fmt.Sprintf("/board/%s", b.Route)), fmt.Sprintf("/%s/", b.Route))
	}
	buffer.AddBlankLine()
	buffer.AddBlankLine()
	buffer.AddLink("https://ko-fi.com/gemchan", "Buy me a coffee! Or don't, that works too!")
	return c.Gemini(buffer.Buffer)
}

func (a *App) board(c gig.Context) error {
	r := c.Param("route")
	b, err := handler.GetBoard(r)
	if err != nil {
		return c.NoContent(gig.StatusRedirectTemporary, baseURL("/"))
	}
	buffer := gem.Gemtext{}
	buffer.AddHeading(fmt.Sprintf("Welcome to /%s/", b.Route))
	buffer.AddUnformatted(b.Description)
	buffer.AddLink(baseURL("/"), "Home")
	buffer.AddBlankLine()
	buffer.AddLink(baseURL(fmt.Sprintf("/newpost/%s", b.Route)), "Create Post")
	buffer.AddSubHeading("Posts")
	buffer.AddBlankLine()
	for _, p := range handler.GetPostsFromBoard(a.DB, b.Route) {
		buffer.AddLink(baseURL("/post/"+p.ID), p.ID)
		buffer.AddUnformatted(p.Time)
		c := handler.GetComments(a.DB, p.ID)
		if len(p.Comments) == 1 {
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
		return c.NoContent(gig.StatusRedirectTemporary, baseURL("/post/"+id))
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
		return c.NoContent(gig.StatusRedirectTemporary, baseURL("/post/"+post.ID))
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
		return c.NoContent(gig.StatusRedirectTemporary, baseURL("/post/"+post.ID))
	}
}

func (a *App) post(c gig.Context) error {
	post := handler.GetPost(a.DB, c.Param("id"))
	buffer := gem.Gemtext{}
	buffer.AddHeading("Post")
	buffer.AddLink(baseURL("/"), "Home")
	buffer.AddLink(baseURL(fmt.Sprintf("/board/%s", post.Board)), fmt.Sprintf("/%s/", post.Board))
	buffer.AddBlankLine()
	buffer.AddLink(baseURL("/post/"+post.ID), post.ID)
	buffer.AddUnformatted(fmt.Sprintf("%s UTC", post.Time))
	buffer.AddBlankLine()
	buffer.AddQuote(post.Content)
	buffer.AddBlankLine()
	buffer.AddLink(baseURL("/addComment/"+post.ID), "Add Comment")
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
		buffer.AddLink(baseURL(fmt.Sprintf("/replyComment/%s&%s", post.ID, c.ID)), "Reply")
		buffer.AddBlankLine()
	}
	return c.Gemini(buffer.Buffer)
}

func (a *App) notFound(c gig.Context) error {
	buffer := gem.Gemtext{}
	buffer.AddHeading("Not Found")
	buffer.AddUnformatted("Damn, you broke it.")
	buffer.AddBlankLine()
	buffer.AddLink(baseURL("/"), "Return Home")
	return c.Gemini(buffer.Buffer)
}

func (a *App) Run(crt, key string) {
	log.Println("Server running at gemini://gemchan.space")
	a.Router.Run(crt, key)
}
