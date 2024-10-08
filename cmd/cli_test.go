package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Cli(t *testing.T) {
	app := setup()
	assert.NoError(t, app.Run([]string{"h"}))
	assert.EqualValues(t, len(app.Commands), 6)
	assert.Equal(t, app.Name, "Xcm tools")

	names := []string{"send", "parse", "tracker", "trackerEthBridge", "trackerS2SBridge", "help"}
	for i := 0; i < len(names); i++ {
		assert.Equal(t, app.Commands[i].Name, names[i])
	}
}
