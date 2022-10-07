package color_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed/color"
)

func TestFromIndex(t *testing.T) {
	assert.Equal(t, color.Red, color.FromIndex(0))
	assert.Equal(t, color.Green, color.FromIndex(1))
	assert.Equal(t, color.Grey, color.FromIndex(19))
	assert.Equal(t, color.Red, color.FromIndex(20))
}
