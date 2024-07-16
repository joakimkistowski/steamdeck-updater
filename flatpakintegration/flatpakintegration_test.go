package flatpakintegration

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_UpdateFlatpakRemotes_DoesNotPanic(t *testing.T) {
	assert.NotPanics(t, UpdateFlatpakRemotes)
}

func Test_GetUpdateableFlatpaks_ReturnsNotNil(t *testing.T) {
	got := GetUpdateableFlatpaks()
	assert.NotNil(t, got)
}
