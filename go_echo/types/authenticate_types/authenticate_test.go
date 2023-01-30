package authenticate_types

import (
	"github.com/stretchr/testify/assert"
	"github.com/tsrkzy/jump_in/helper"
	"testing"
)

func TestMaskMailAddress(t *testing.T) {
	m := "tsrmix-_@gmail.com"
	masked := helper.MaskMailAddress(m)
	assert.Equal(t, "t******_@gmail.com", masked)
}
