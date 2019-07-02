package dateutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrToTime(t *testing.T) {
	result, _ := StrToTime("16-05-12 11:22:33GMT+07:00")
	assert.Nil(t, result)
}
