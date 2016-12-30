package lib

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestResolveOneRedirect(t *testing.T) {
	res := Resolve("http://bit.ly/2hcXx5Z", "http://bit.ly/", 10, false, "Googlebot")
	assert.Equal(t, "https://maddevs.io/", res)
}

func TestResolveTwoRedirects(t *testing.T) {
	res := Resolve("http://ow.ly/Pkxu306YmRs", "http://ow.ly/", 10, false, "Googlebot")
	assert.Equal(t, "https://maddevs.io/", res)
}

func TestResolveSSL(t *testing.T) {
	res := Resolve("https://bit.ly/2hcXx5Z", "http://bit.ly/", 10, false, "Googlebot")
	assert.Equal(t, "https://maddevs.io/", res)

	res = Resolve("https://bit.ly/ID7AM5", "http://bit.ly/", 10, false, "Googlebot")
	assert.Equal(t, "https://www.youtube.com/", res)
}
