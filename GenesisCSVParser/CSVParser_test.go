package GenesisCSVParser

import (
	"testing"
)

func TestParseFile(t *testing.T) {
	resp, err := ParseFile("../genesis.csv")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(EntriesToJSON(resp[0:10]))
	}
}
