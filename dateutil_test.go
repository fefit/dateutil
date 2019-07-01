package dateutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrToTime(t *testing.T) {
	result, _ := StrToTime("16-5-12 11:22:33")
	assert.Nil(t, result)
}
