package GenesisCSVParser

import (
	"encoding/csv"
	"errors"
	"github.com/FactomProject/factoid"
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

func EntriesToBalanceMap(entries []Entry) ([]Balance, error) {
	balanceMap := map[string]Balance{}

	for _, v := range entries {
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

var MaxOutputsPerTransaction int = 20

func CreateTransactions(balances []Balance) []*factoid.Transaction {
	answer := make([]*factoid.Transaction, 0, len(balances)/MaxOutputsPerTransaction+1)
	for i := 0; i < len(balances); i += MaxOutputsPerTransaction {
		max := i + MaxOutputsPerTransaction
		if max > len(balances) {
			max = len(balances)
		}
		t := CreateTransaction(balances[i:max])
		answer = append(answer, t)
	}
	return answer
}

func CreateTransaction(balances []Balance) *factoid.Transaction {
	t := new(factoid.Transaction)
	for _, v := range balances {
		t.AddOutput(v.IAddress, v.FactoshiBalance)
	}
	return t
}
