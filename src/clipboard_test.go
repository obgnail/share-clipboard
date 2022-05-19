package src

import (
	"testing"
)

func TestClipboard(t *testing.T) {
	res, err := GetBytesFromClip()
	if err != nil {
		t.Log(err)
	}
	t.Log("======== ", string(res))
}
