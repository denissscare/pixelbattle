package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {

	p := Pixel{
		X:         MaxWidth + 10,
		Y:         20,
		Color:     "#FF0000",
		Author:    "user123",
		Timestamp: time.Now(),
	}
	require.Error(t, p.Validate())

	p = Pixel{
		X:         10,
		Y:         20,
		Color:     "#FF0000",
		Author:    "",
		Timestamp: time.Now(),
	}
	require.Error(t, p.Validate())

	p = Pixel{
		X:         10,
		Y:         20,
		Color:     "#FF0000",
		Author:    "artist",
		Timestamp: time.Now().Add(time.Hour),
	}
	require.Error(t, p.Validate())

	assert.Equal(t, 1, 0)
}
