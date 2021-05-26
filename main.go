package main

import (
	app "gem/app"
)

func main() {
	app := &app.App{}
	app.Init()
	app.Run("gemchan.space.crt", "gemchan.space.key")
}
