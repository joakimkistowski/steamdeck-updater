package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_getCurrentUpgradingMessage_CreatesDifferentMessagesAfterWait(t *testing.T) {
	got0 := getCurrentUpgradingMessage()
	time.Sleep(1 * time.Second)
	got1 := getCurrentUpgradingMessage()
	assert.NotEmpty(t, got0)
	assert.NotEmpty(t, got1)
	assert.NotEqual(t, got0, got1)
}
