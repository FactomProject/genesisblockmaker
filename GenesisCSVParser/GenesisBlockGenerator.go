package main

import (
	"encoding/hex"
	"io/ioutil"
	"log"
)

func main() {
	resp, err := ParseFile("../genesis.csv")
	if err != nil {
		panic(err)
	}
	entries, err := EntriesToBalances(resp)
	if err != nil {
		panic(err)
	}
	genesis, _, _, err := CreateTransactions(entries)
	if err != nil {
		panic(err)
	}
	bin, err := genesis.MarshalBinary()
	if err != nil {
		panic(err)
	}
	log.Printf("Genesis Hash - %X", genesis.GetHash().Bytes())
	WriteToFile("Genesis.txt", hex.EncodeToString(bin))
	WriteToFile("GenesisBrokenDown.txt", BreakTextDown(hex.EncodeToString(bin), 120))
	WriteToFile("Wallet.txt", "TODO: dump wallet")

	//Genesis block for testing
	resp2, err := ParseFile("../testing.csv")
	if err != nil {
		panic(err)
	}

	entries2, err := EntriesToBalances(resp2)
	if err != nil {
		panic(err)
	}
	genesis2, _, _, err := CreateTransactions(append(entries, entries2...))
	if err != nil {
		panic(err)
	}
	bin2, err := genesis2.MarshalBinary()
	if err != nil {
		panic(err)
	}

	log.Printf("Test genesis Hash - %X", genesis2.GetHash().Bytes())
	WriteToFile("TestGenesis.txt", hex.EncodeToString(bin2))
	WriteToFile("TestGenesisBrokenDown.txt", BreakTextDown(hex.EncodeToString(bin2), 120))
	WriteToFile("TestWallet.txt", "TODO: dump wallet")
}

func WriteToFile(filename, content string) {
	ioutil.WriteFile(filename, []byte(content), 0777)
}

func BreakTextDown(text string, charactersPerLine int) string {
	var answer string = ""
	max := 0
	for {
		min := max
		max += charactersPerLine
		if max >= len(text) {
			max = len(text)
			answer = answer + "\"" + text[min:max] + "\"\n"
			break
		} else {
			answer = answer + "\"" + text[min:max] + "\" +\n"
		}
	}
	return answer
}
