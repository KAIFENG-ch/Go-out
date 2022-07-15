package Go_out

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNew(t *testing.T) {
	router := New()
	assert.Len(t, router.Handlers, 2)
}
