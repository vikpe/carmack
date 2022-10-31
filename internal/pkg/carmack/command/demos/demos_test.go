package demos_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	demoCmd "github.com/vikpe/carmack/internal/pkg/carmack/command/demos"
	"github.com/vikpe/qw-hub-api/pkg/qtvscraper"
)

func TestContentFromDemos(t *testing.T) {
	timeStamp, _ := time.Parse("2006-04-02 15:01", "2020-01-01 12:00")
	demos := []qtvscraper.Demo{
		{
			QtvAddress:  "beta:28000",
			Time:        timeStamp,
			Filename:    "duel_xantom_vs_bro[dm6]221028-0355.mvd",
			DownloadUrl: "http://beta:28000/dl/demos/duel_xantom_vs_bro[dm6]221028-0355.mvd",
			QtvplayUrl:  "file:duel_xantom_vs_bro[dm6]221028-0355.mvd@beta:28000",
		},
	}

	expect := strings.Join([]string{
		"**Demo search**: found 1 demo(s)",
		"`0001-01-01 00:00`: [duel - **xantom vs bro** \\[dm6\\]](http://beta:28000/dl/demos/duel_xantom_vs_bro[dm6]221028-0355.mvd)",
	}, "\n")

	assert.Equal(t, expect, demoCmd.ContentFromDemos(demos))
}
