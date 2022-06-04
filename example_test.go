package config_test

import (
	"context"
	"fmt"

	"gojini.dev/config"
)

func ExampleStore() {
	store := config.New()
	cfgStr := `{"log": {"file": "example.log"}, "server": {"port": 8088}}`

	if e := store.LoadFromStr(context.Background(), cfgStr); e != nil {
		panic(e)
	}

	l := &struct {
		File string `json:"file"`
	}{File: "default.log"}

	if e := store.Get("log", l); e == nil {
		fmt.Println("log file:", l.File)
	}

	s := &struct {
		Port int `json:"port"`
	}{Port: 0}
	if e := store.Get("server", s); e == nil {
		fmt.Println("server port:", s.Port)
	}

	// Output:
	// log file: example.log
	// server port: 8088
}
