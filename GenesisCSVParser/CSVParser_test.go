package GenesisCSVParser

import (
	"testing"
)

func TestParseFile(t *testing.T) {
	resp, err := ParseFile("../genesis.csv")
	if err != nil {
		t.Error(err)
	} else {
		//t.Log(EntriesToJSON(resp[0:10]))
		entries, err := EntriesToBalanceMap(resp[:40])
		if err != nil {
			t.Error(err)
		} else {
			genesis, txs, err := CreateTransactions(entries)
			if err != nil {
				t.Error(err)
			} else {
				t.Log(EncodeJSONString(txs))
				t.Log(EncodeJSONString(genesis))
			}
		}
	}

}
