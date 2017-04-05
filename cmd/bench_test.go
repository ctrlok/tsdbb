package cmd

import "testing"
import "github.com/stretchr/testify/assert"

func TestParseListen(t *testing.T) {
	var uri string
	var str string
	var err error

	str = "localhost:80"
	uri = parseListen(str)
	assert.NoError(t, err)
	assert.Equal(t, str, uri)

	str = ":80"
	uri = parseListen(str)
	assert.NoError(t, err)
	assert.Equal(t, str, uri)

	str = "127.1:80"
	uri = parseListen(str)
	assert.NoError(t, err)
	assert.Equal(t, str, uri)

	str = "80"
	uri = parseListen(str)
	assert.NoError(t, err)
	assert.Equal(t, ":80", uri)
}
