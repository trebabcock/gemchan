package handler

import (
	"errors"
	"gemchan/app/model"
	"log"
)

var boards = []model.Board{
	{
		Route:       "b",
		Description: "Random discussion.",
		ID:          "ZfWHz7HRN7b7MEd2fvsAA5",
	},
	{
		Route:       "tech",
		Description: "Technology discussion.",
		ID:          "Y9xyaYeptExQL8jqrtJBYd",
	},
	{
		Route:       "sci",
		Description: "Science discussion.",
		ID:          "LALPgrCnazTUNmp7AwuRRc",
	},
	{
		Route:       "hobby",
		Description: "Hobby discussion.",
		ID:          "j3UK7DFUu9wqCmKjJuo2sc",
	},
	{
		Route:       "game",
		Description: "Gaming discussion.",
		ID:          "mq7CKKbDsFwikbuiNLGoLr",
	},
	{
		Route:       "gem",
		Description: "Gemini discussion.",
		ID:          "kASdYvnyDhnUzGQAYoVoAx",
	},
	{
		Route:       "meta",
		Description: "This board is for the discussion of Gemchan, including board requests, feature requests, or bug reports.",
		ID:          "hzjoYkDki8yJGfaL6gH2No",
	},
}

func GetBoards() []model.Board {
	return boards
}

func GetBoard(route string) (model.Board, error) {
	for _, b := range boards {
		if b.Route == route {
			return b, nil
		}
	}
	log.Println("Couldn't find board")
	return model.Board{}, errors.New("BOARD NOT FOUND")
}
