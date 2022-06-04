package store_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gojini.dev/config"
)

func TestStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	store := config.New()
	cfgStr := `{"log": {"file": "example.log"}, "server": {"port": 8088}}`

	assert.Nil(store.LoadFromStr(context.Background(), cfgStr))

	l := &struct {
		File string `json:"file"`
	}{File: "default.log"}

	assert.Nil(store.Get("blah", nil))
	assert.Nil(store.Get("", l))
	assert.Equal("default.log", l.File)
	assert.Nil(store.Get("log", l))
	assert.Equal("example.log", l.File)

	s := &struct {
		Port int `json:"port"`
	}{Port: 0}

	assert.Nil(store.Get("server", s))
	assert.Equal(s.Port, 8088)
	assert.NotNil(store.Get("blah", l))

	j := 0
	assert.Nil(store.LoadFromStr(context.Background(), `{"log": "hello"}`))
	assert.NotNil(store.Get("server", &j))
}

func TestBadJsonStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	store := config.New()
	cfgStr := `{{"file": "example.log"}, "server": {"port": 8088}}`

	assert.NotNil(store.LoadFromStr(context.Background(), cfgStr))
	assert.NotNil(store.Load(context.Background(), badReader{}))
}

type badReader struct {
}

var brErr = errors.New("br error")

func (br badReader) Read([]byte) (int, error) {
	return 0, brErr
}
