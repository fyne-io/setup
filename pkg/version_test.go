package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFyneVersion(t *testing.T) {
	assert.Equal(t, "v2.5.5", parseFyneVersion(
		`fyne cli version: v2.5.5`))
	assert.Equal(t, "v1.6.1", parseFyneVersion(
		`fyne cli version: v1.6.1
fyne library version: v2.6.0`))
}
