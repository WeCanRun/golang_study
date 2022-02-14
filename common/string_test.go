package common

import "testing"

func TestMapWithNil(t *testing.T) {
	var hash map[string]*string = nil
	t.Logf(":%v", len(hash))
}
