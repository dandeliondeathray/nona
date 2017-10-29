package plumber

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/linkedin/goavro"
)

// Codecs maintains a collection of all Avro codecs found at a given directory path.
type Codecs struct {
	codecs map[string]*goavro.Codec
}

// ByName finds a given codec by its schema name.
func (c *Codecs) ByName(name string) (*goavro.Codec, error) {
	codec, ok := c.codecs[name]
	if !ok {
		return nil, fmt.Errorf("No such codec %s", name)
	}
	return codec, nil
}

// LoadCodecsFromPath loads all codecs files found in a given path.
func LoadCodecsFromPath(path string) (*Codecs, error) {
	codecs := &Codecs{make(map[string]*goavro.Codec)}

	pattern := filepath.Join(path, "*.avsc")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	log.Println("Schema files:", matches)

	for i := range matches {
		match := matches[i]
		if err = codecs.loadCodecFromFile(path, match); err != nil {
			log.Println("Failed to load schema in file", match, "because", err)
		}
	}

	return codecs, nil
}

type schemaName struct {
	Name string `json:"name"`
}

func (c *Codecs) loadCodecFromFile(path, file string) error {
	schema, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	codec, err := goavro.NewCodec(string(schema))
	if err != nil {
		return err
	}

	var schemaName schemaName
	err = json.Unmarshal(schema, &schemaName)
	if err != nil {
		return err
	}
	log.Println("Loaded codec:", schemaName.Name)

	c.codecs[schemaName.Name] = codec
	return nil
}
