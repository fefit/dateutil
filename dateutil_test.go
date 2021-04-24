package dateutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrToTime(t *testing.T) {
	result, _ := DateTime("IV")
	assert.NotNil(t, result)
}
