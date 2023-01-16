package main

import (
	"embed"
	app "gemchan/app"
	"gemchan/keys"
)

var (
	//go:embed assets
	fs embed.FS
)

func main() {
	app := &app.App{}
	app.Init()

	crt, key := keys.GetKeys(&fs)

	app.Run(crt, key)

}
