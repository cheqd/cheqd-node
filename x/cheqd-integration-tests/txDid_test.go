package cheqd_integration_tests

import (
	"github.com/rs/zerolog/log"
	"testing"
)

func TestDidCli(t *testing.T) {
	_, err := Setup()
	if err != nil {
		log.Err(err)
		t.Fail()
	}
}
