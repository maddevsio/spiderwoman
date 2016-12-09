package lib

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestResolveOneRedirect(t *testing.T) {
	res := Resolve("http://bit.ly/2hcXx5Z", "http://bit.ly/", 10, false, "Googlebot")
	assert.Equal(t, "http://maddevs.io/", res)
}

func TestResolveTwoRedirects(t *testing.T) {
	res := Resolve("http://ow.ly/Pkxu306YmRs", "http://ow.ly/", 10, false, "Googlebot")
	assert.Equal(t, "http://maddevs.io/", res)
}