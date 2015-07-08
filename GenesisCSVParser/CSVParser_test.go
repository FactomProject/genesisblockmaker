package GenesisCSVParser

import (
	"testing"
)

func TestEverything(t *testing.T) {
	resp, err := ParseFile("../genesis.csv")
	if err != nil {
		t.Error(err)
	} else {
		//t.Log(EntriesToJSON(resp[0:10]))
		entries, err := EntriesToBalances(resp[:40])
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

func TestParseFile(t *testing.T) {
	_, err := ParseFile("testInvalid.csv")
	if err == nil {
		t.Error("ParseFile did not fail on invalid CSV file")
	}
	resp, err := ParseFile("testValid.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(resp) != 9 {
		t.Error("Invalid number of entries returned")
	}
	if resp[0].FundingTxID != "001fa3b6f6a9be524d182815cbd8c6f59859e1e6c74c4bec5f44816a52d5bc4d" {
		t.Error("FundingTxID not filled properly")
	}
	if resp[1].ConfTime != 1428420172 {
		t.Error("ConfTime not filled properly")
	}
	if resp[2].ConfTimeHuman != "2015-Apr-01 17:30:20" {
		t.Error("ConfTimeHuman not filled properly")
	}
	if resp[3].ED25519PubKey != "0" {
		t.Error("ED25519PubKey not filled properly")
	}
	if resp[4].Bitcoins != 35200000 {
		t.Error("Bitcoins not filled properly")
	}
	if resp[5].Rate != 2000 {
		t.Error("Rate not filled properly")
	}
	if resp[6].Factoshis != 100000000000 {
		t.Error("Factoshis not filled properly")
	}
	if resp[1].Notes != "sent before cutoff per koinify" {
		t.Error("Notes not filled properly")
	}
}

func TestEntriesToBalances(t *testing.T) {
	resp, err := ParseFile("testValid.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	balances, err := EntriesToBalances(resp)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(EncodeJSONString(balances))
	if len(balances) != 5 {
		t.Error("Invalid number of balances returned")
	}
	if !findAndValidateBalance(balances, "9c3e8ec1786f08cc7d3fc10ba2865dea1c849f9ad5318872055b60eaf93d4753", 379800000000) {
		t.Error("Invalid balance calculated")
	}
	if !findAndValidateBalance(balances, "154bfe5238638f92b74bcb36e08d8fc06c638f66b907f68edfed5dbd39c35323", 8000000000) {
		t.Error("Invalid balance calculated")
	}
	if !findAndValidateBalance(balances, "0000000000000000000000000000000000000000000000000000000000000000", 110000000000) {
		t.Error("Invalid balance calculated")
	}
	if !findAndValidateBalance(balances, "89b170ef56e7e3137f79041f0e1ea1875fdec489f4590145b6adc83c8850321a", 56320000000) {
		t.Error("Invalid balance calculated")
	}
	if !findAndValidateBalance(balances, "b8bb22b157633d18a368bbe217425f148166961bf27a02c5aadd63bbaf7ba990", 160047500000) {
		t.Error("Invalid balance calculated")
	}
}

func findAndValidateBalance(balances []Balance, key string, balance uint64) bool {
	for _, v := range balances {
		if v.ED25519PubKey == key {
			return v.FactoshiBalance == balance
		}
	}
	return false
}

func TestCreateTransactions(t *testing.T) {
	resp, err := ParseFile("testValid.csv")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	balances, err := EntriesToBalances(resp)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	MaxOutputsPerTransaction = 2
	genesis, txs, err := CreateTransactions(balances)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	if len(txs) != 3 {
		t.Error("Invalid number of transactions returned")
	}
	var outSum uint64
	var inSum uint64
	for _, v := range txs {
		out, ok := v.TotalOutputs()
		if ok == false {
			t.Error("ok == false")
		}
		outSum += out

		in, ok := v.TotalInputs()
		if ok == false {
			t.Error("ok == false")
		}
		inSum += in
	}
	if outSum != 714167500000 {
		t.Error("Invalid output sum")
	}
	genesisOut, ok := genesis.TotalOutputs()
	if ok == false {
		t.Error("ok == false")
	}
	if genesisOut != inSum {
		t.Error("Genesis output doesn't add to input sum for transactions")
	}
}
