package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// Store represents a map of string keys and []byte values. In this case, the
// key is the configuration name and the value is the related config.
type Store struct {
	kb map[string][]byte
}

var ErrKeyNotFound = errors.New("key not found")

// New returns a config store that can hold a map of key and arbitarary config
// blob. The key is the name of the configuration that can be retrieved using
// the Get method once configuration is loaded.
func New() *Store {
	return &Store{kb: map[string][]byte{}}
}

// Load loads the configuration using a io.Reader. The reader must point to a
// valid json stream of maps. An example configuration is:
//
//  {"log": {"file": "example.log"}, "server": {"port": 8088}}
//
// this config has two keys, log and server with corresponding configuration
// which can be retreived using Get method.
//
// Returns error if the configuration stream is not parsable in the above format.
func (s *Store) Load(ctx context.Context, src io.Reader) error {
	d := json.NewDecoder(src)

	for {
		data := map[string]*cfgData{}
		if err := d.Decode(&data); errors.Is(err, io.EOF) {
			return nil
		} else if err != nil {
			return fmt.Errorf("load error: %w", err)
		}

		// Cache the key and its corresponding json data
		for k, v := range data {
			s.kb[k] = v.b
		}
	}
}

// LoadFromStr loads configuration from a string.
func (s *Store) LoadFromStr(ctx context.Context, cfg string) error {
	return s.Load(ctx, strings.NewReader(cfg))
}

// LoadFromFile loads configuration from a file.
func (s *Store) LoadFromFile(ctx context.Context, cfgFile string) error {
	cfg, err := os.Open(cfgFile)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}

	defer cfg.Close()

	return s.Load(ctx, cfg)
}

// Get retrieves the config specified by the key if exists. If the config
// cannot be parsed into the provided structure it returns an error.
func (s *Store) Get(key string, config interface{}) error {
	if config == nil || key == "" {
		// Empty key or no config so just return the default config back
		return nil
	}

	if b, ok := s.kb[key]; ok {
		// Bad buffer for the current type but lets keep it around
		// in case the registry is modified with a new type
		// and we can process it in future Get calls
		if e := json.Unmarshal(b, config); e != nil {
			return fmt.Errorf("config format error %w", e)
		}

		return nil
	}

	return fmt.Errorf("%s %w", key, ErrKeyNotFound)
}

// cfgData represents a byte slice where we store information about a specific
// configuration.
type cfgData struct {
	b []byte
}

// UnmarshalJSON sets the cfgData to the specified byte slice, returning nil
// if no error occurs.
func (d *cfgData) UnmarshalJSON(b []byte) error {
	d.b = b

	return nil
}
