package GenesisCSVParser

import(
	"encoding/json"
	"strconv"
	"regexp"
)


func String2Int64(s string) (int64, error) {
	answer, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return answer, nil
}

func IsHex(s string) bool {
	matched, _:=regexp.MatchString("(?:0[xX])?[0-9a-fA-F]+", s)
	return matched
}

func EncodeJSONString(data interface{}) (string, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(encoded), err
}