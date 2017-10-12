package slackmessaging

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/linkedin/goavro"
)

type Codecs struct {
	codecs map[string]*goavro.Codec
}

func (c *Codecs) ByName(name string) (*goavro.Codec, error) {
	codec, ok := c.codecs[name]
	if !ok {
		return nil, fmt.Errorf("No such codec %s", name)
	}
	return codec, nil
}

func LoadCodecsFromPath(path string) (*Codecs, error) {
	codecs := &Codecs{make(map[string]*goavro.Codec)}

	err := codecs.loadCodecFromFile(path, "PuzzleNotification.avsc", "nona_PuzzleNotification")
	if err != nil {
		return nil, err
	}

	err = codecs.loadCodecFromFile(path, "Chat.avsc", "Chat")
	if err != nil {
		return nil, err
	}

	return codecs, nil
}

func (c *Codecs) loadCodecFromFile(path, file, name string) error {
	schema, err := ioutil.ReadFile(filepath.Join(path, file))
	if err != nil {
		return err
	}
	codec, err := goavro.NewCodec(string(schema))
	if err != nil {
		return err
	}
	c.codecs[name] = codec
	return nil
}
