package config

import (
	"log"
	"os"
	"path"
	"testing"
)

func TestConfig(t *testing.T) {
	var fail bool
	ap := path.Join(os.Getenv("GOPATH"), "src/github.com/alextonkonogov/ascii_counter")
	cfg := NewConfig()
	err := cfg.SetConfigFromJson(path.Join(ap, "ftp.json"))
	if err != nil {
		log.Println(err)
		fail = true
	}
	if fail {
		t.Fail()
	}
}
