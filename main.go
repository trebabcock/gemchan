package main

import (
	app "gem/app"
)

func main() {
	app := &app.App{}
	app.Init()
	app.Run("cert.pem", "privkey.pem")
}
