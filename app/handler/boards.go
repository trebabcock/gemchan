package handler

import (
	"gem/app/model"
	"log"
)

var boards = []model.Board{
	{
		Route:       "rand",
		Description: "This board is for random discussion.",
		ID:          "f35f4d09-b49f-4220-b22a-bd0e612d64bd",
	},
	{
		Route:       "phys",
		Description: "This board is for the discusson of the physical sciences.",
		ID:          "2467d7f4-3565-462a-82e4-c51f433aa7e6",
	},
	{
		Route:       "math",
		Description: "This board is for the discusson of mathematics.",
		ID:          "29b01ff2-fcaa-41d6-a675-be1e623f43c1",
	},
	{
		Route:       "cs",
		Description: "This board is for the discusson of computer science.",
		ID:          "cbdbe8fa-745a-4393-a4ce-ca2c40aeed0a",
	},
	{
		Route:       "gem",
		Description: "This board is for the discusson of all things Gemini.",
		ID:          "5feb3921-6a6b-4187-aaec-f003b72436a2",
	},
	{
		Route:       "meta",
		Description: "This board is for the discussion of Gemchan, including board requests, feature requests, or bug reports.",
		ID:          "a3acffa4-3d83-4ff8-a911-411f2a206a8d",
	},
}

func GetBoards() []model.Board {
	return boards
}

func GetBoard(route string) model.Board {
	for _, b := range boards {
		if b.Route == route {
			return b
		}
	}
	log.Println("Couldn't find board")
	return model.Board{}
}
