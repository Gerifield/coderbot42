package jatek

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetRegs(t *testing.T) {
	users := make([]string, 0)
	setRegs(&users, "a", 4)

	assert.Equal(t, 4, len(users))
	assert.Equal(t, []string{"a", "a", "a", "a"}, users)

	setRegs(&users, "a", 2)
	assert.Equal(t, []string{"a", "a", "a", "a"}, users)
}
