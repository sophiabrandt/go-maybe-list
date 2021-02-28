package handlers

import (
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

const BASE_URL = "https://gomaybelist.app.localhost"

func TestDebug(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Get(BASE_URL + "/debug/health")
	if err != nil {
		t.Fail()
	}
	assert.Equal(t, 200, resp.StatusCode())
}
