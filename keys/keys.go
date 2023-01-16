package keys

import (
	"embed"
	"log"
)

func GetKeys(fs *embed.FS) ([]byte, []byte) {
	crt, err := fs.ReadFile("assets/gemchan.space.crt")
	if err != nil {
		log.Fatal(err)
	}

	key, err := fs.ReadFile("assets/gemchan.space.key")
	if err != nil {
		log.Fatal(err)
	}

	return crt, key
}
