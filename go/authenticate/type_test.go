package authenticate

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMaskMailAddress(t *testing.T) {
	m := "tsrmix-_@gmail.com"
	masked := maskMailAddress(m)
	assert.Equal(t, "t******_@gmail.com", masked)
}
