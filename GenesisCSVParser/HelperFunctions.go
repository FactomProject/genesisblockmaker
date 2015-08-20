package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/FactomProject/factoid"
	"regexp"
	"strconv"
)

func String2Int64(s string) (int64, error) {
	answer, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return answer, nil
}
func String2UInt64(s string) (uint64, error) {
	answer, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return answer, nil
}

func IsHex(s string) bool {
	matched, _ := regexp.MatchString("(?:0[xX])?[0-9a-fA-F]+", s)
	return matched
}

func EncodeJSONString(data interface{}) (string, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(encoded), err
}

func ED25519PubKeyToIAddress(ed25519PubKey string) (factoid.IAddress, error) {
	hex, err := hex.DecodeString(ed25519PubKey)
	if err != nil {
		return nil, err
	}
	ircd := factoid.NewRCD_1(hex)
	return ircd.GetAddress()
}
