package GenesisCSVParser

import (
	"encoding/csv"
	"errors"
	"github.com/FactomProject/factoid"
	"github.com/FactomProject/factoid/block"
	"github.com/FactomProject/factoid/wallet"
	"os"
)

type Entry struct {
	FundingTxID string `json:"funding txid,"`

	ConfTime      int64  `json:"conf time unix,"`
	ConfTimeHuman string `json:"conf time human,"`

	ED25519PubKey string `json:"ed25519 pubkey,"`

	Bitcoins  uint64 `json:"# bitcoins,"`
	Rate      uint64 `json:"rate,"`
	Factoshis uint64 `json:"# factoshis,"`

	Notes string `json:"notes,omitempty"`
}

func (e *Entry) Validate() error {
	if len(e.FundingTxID) != 64 {
		return errors.New("Invalid FundingTxID length")
	}
	if IsHex(e.FundingTxID) == false {
		return errors.New("FundingTxID is not a valid hexaadecimal number")
	}
	/*if len(e.ED25519PubKey)!=64 {
		return errors.New("Invalid ED25519PubKey length")
	}*/
	if IsHex(e.ED25519PubKey) == false {
		return errors.New("ED25519PubKey is not a valid hexaadecimal number")
	}

	if e.Bitcoins*e.Rate != e.Factoshis {
		return errors.New("Factoshis don't add up")
	}

	return nil
}

//Function to parse the CSV file into a list of entries
func ParseFile(filePath string) ([]Entry, error) {
	//Open the file
	csvfile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = 8

	rawCSVdata, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	//Parse all of the entries and return them
	return ParseEntries(rawCSVdata)
}

func ParseEntries(rawCSVdata [][]string) ([]Entry, error) {
	answer := []Entry{}
	//skipping first entry with field names
	for i := 1; i < len(rawCSVdata); i++ {
		v := rawCSVdata[i]
		entry := Entry{}
		entry.FundingTxID = v[0]
		intVar, err := String2Int64(v[1])
		if err != nil {
			return nil, err
		}
		entry.ConfTime = intVar
		entry.ConfTimeHuman = v[2]
		entry.ED25519PubKey = v[3]
		uintVar, err := String2UInt64(v[4])
		if err != nil {
			return nil, err
		}
		entry.Bitcoins = uintVar
		uintVar, err = String2UInt64(v[5])
		if err != nil {
			return nil, err
		}
		entry.Rate = uintVar
		uintVar, err = String2UInt64(v[6])
		if err != nil {
			return nil, err
		}
		entry.Factoshis = uintVar
		entry.Notes = v[7]

		//Validate the entry to make sure data is correct
		err = entry.Validate()
		if err != nil {
			return nil, err
		}

		answer = append(answer, entry)
	}
	return answer, nil
}

//Function that converts an array of entries into JSON string
func EntriesToJSON(entries []Entry) (string, error) {
	return EncodeJSONString(entries)
}

type Balance struct {
	ED25519PubKey   string
	RCD             factoid.IRCD
	IAddress        factoid.IAddress
	FactoshiBalance uint64
}

func EntriesToBalances(entries []Entry) ([]Balance, error) {
	balanceMap := map[string]Balance{}

	for _, v := range entries {
		//TODO: FIXME: figure out what is the correct way to handle the zero accounts
		if v.ED25519PubKey == "0" {
			v.ED25519PubKey = "0000000000000000000000000000000000000000000000000000000000000000"
		}
		balance, ok := balanceMap[v.ED25519PubKey]
		if ok == false {
			balance.ED25519PubKey = v.ED25519PubKey
			iAddress, err := ED25519PubKeyToIAddress(v.ED25519PubKey)
			if err != nil {
				return nil, err
			}
			balance.IAddress = iAddress
			rcd, err := factoid.NewRCD_2(1, 1, []factoid.IAddress{iAddress})
			if err != nil {
				return nil, err
			}
			balance.RCD = rcd
		}
		balance.FactoshiBalance += v.Factoshis
		balanceMap[v.ED25519PubKey] = balance
	}

	answer := make([]Balance, 0, len(balanceMap))
	for _, v := range balanceMap {
		answer = append(answer, v)
	}

	return answer, nil
}

//TODO: double-check the magic numbers
var MaxOutputsPerTransaction int = 250 //Hot many outputs will be included in each transaction to keep it under the size limit
var FactoshisPerEC uint64 = 1000

//Function that creates a set of transactions from the list of balances the users should receive, as well as the corresponding genesis transaction
func CreateTransactions(balances []Balance) (block.IFBlock, []factoid.ITransaction, error) {
	answer := make([]factoid.ITransaction, 0, len(balances)/MaxOutputsPerTransaction+1)
	w := new(wallet.SCWallet)
	w.Init()
	inputAddress, err := w.GenerateFctAddress([]byte("Genesis"), 1, 1)
	if err != nil {
		return nil, nil, err
	}
	for i := 0; i < len(balances); i += MaxOutputsPerTransaction {
		max := i + MaxOutputsPerTransaction
		if max > len(balances) {
			max = len(balances)
		}
		t, err := CreateTransaction(balances[i:max], w, inputAddress)
		if err != nil {
			return nil, nil, err
		}
		answer = append(answer, t)
	}
	genesis, err := GetGenesisBlock(0, answer, w, inputAddress)
	if err != nil {
		return nil, nil, err
	}
	return genesis, answer, nil
}

//Creates a transaction crediting the given users
func CreateTransaction(balances []Balance, w *wallet.SCWallet, address factoid.IAddress) (factoid.ITransaction, error) {
	t := w.CreateTransaction(0)
	for _, v := range balances {
		t.AddOutput(v.IAddress, v.FactoshiBalance)
	}
	outputTotal, err := t.TotalOutputs()
	if err != nil {
		return nil, err
	}

	w.AddInput(t, address, outputTotal)

	fees, err := t.CalculateFee(FactoshisPerEC)
	if err != nil {
		return nil, err
	}

	w.UpdateInput(t, 0, address, outputTotal+fees)

	ok, err := w.SignInputs(t)
	if ok == false {
		return nil, errors.New("Unable to sign inputs")
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

/*
//A placeholder function for creating the genesis transaction creating as many factoshis as are needed to credit the various accounts
func CreateGenesisTransaction(transactions []factoid.ITransaction, w *wallet.SCWallet, address factoid.IAddress) (factoid.ITransaction, error) {
	//TODO: update for proper genesis transaction generation before launch
	t := w.CreateTransaction(0)
	var sum uint64 = 0
	for _, v := range transactions {
		input, ok := v.TotalInputs()
		if ok == false {
			return nil, errors.New("TotalInputs returned false")
		}
		sum += input
	}
	w.AddOutput(t, address, sum)

	genesisAddress, err := ED25519PubKeyToIAddress("0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		return nil, err
	}

	w.AddInput(t, genesisAddress, sum)

	return t, nil
}*/

func GetGenesisBlock(ftime uint64, transactions []factoid.ITransaction, w *wallet.SCWallet, address factoid.IAddress) (block.IFBlock, error) {
	genesisBlock := block.NewFBlock(1000000, uint32(0))

	t := w.CreateTransaction(ftime)
	var sum uint64 = 0
	for _, v := range transactions {
		input, err := v.TotalInputs()
		if err != nil {
			return nil, err
		}
		sum += input
	}
	w.AddOutput(t, address, sum)

	err := genesisBlock.AddCoinbase(t)
	if err != nil {
		return nil, err
	}
	genesisBlock.GetHash()

	return genesisBlock, nil
}
