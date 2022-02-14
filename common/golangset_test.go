package common

import (
	"encoding/json"
	"github.com/deckarep/golang-set"
	"testing"
)

type s struct {
	Set   mapset.Set `json:"set"`
	Other string     `json:"other"`
}

func TestGolangSet(t *testing.T) {
	s1 := &s{
		//Set:   mapset.NewSet(),
		Other: "other",
	}
	//s1.Set.Add(1)
	marshal, err := json.Marshal(s1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(marshal)

	if err := json.Unmarshal(marshal, &s1); err != nil {
		t.Fatal(err)
	}
	t.Log(s1)

	s3 := &s{
		Set:   mapset.NewSet(),
		Other: "",
	}

	stringData := `{
		"set": [1,2],
        "other":"other"
	}`
	if err := json.Unmarshal([]byte(stringData), s3); err != nil {
		t.Fatal(err)
	}
	t.Log(s3)

}
