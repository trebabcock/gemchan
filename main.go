package main

import (
	app "gem/app"
)

func main() {
	app := &app.App{}
	app.Init()
	app.Run("server.crt", "server.key")
}
