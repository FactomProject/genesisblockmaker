package GenesisCSVParser

import (
	"encoding/csv"
	"errors"
	"os"
)

type Entry struct {
	FundingTxID string `json:"funding txid,"`

	ConfTime      int64  `json:"conf time unix,"`
	ConfTimeHuman string `json:"conf time human,"`

	ED25519PubKey string `json:"ed25519 pubkey,"`

	Bitcoins  int64 `json:"# bitcoins,"`
	Rate      int64 `json:"rate,"`
	Factoshis int64 `json:"# factoshis,"`

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

	//Parse all of the entries
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
		intVar, err = String2Int64(v[4])
		if err != nil {
			return nil, err
		}
		entry.Bitcoins = intVar
		intVar, err = String2Int64(v[5])
		if err != nil {
			return nil, err
		}
		entry.Rate = intVar
		intVar, err = String2Int64(v[6])
		if err != nil {
			return nil, err
		}
		entry.Factoshis = intVar
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
