package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	infos := listAllCommands()
	assert.Equal(t, len(commands), len(infos))
}
