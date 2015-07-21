package main

import (
	"testing"
)

func TestIsHex(t *testing.T) {
	if IsHex("") == true {
		t.Error("Empty string shouldn't be a valid hex")
	}
	if IsHex("0") == false {
		t.Error("0 should be a valid hex")
	}
	if IsHex("0x00") == false {
		t.Error("0x00 should be a valid hex")
	}
	if IsHex("0xFF") == false {
		t.Error("0xFF should be a valid hex")
	}
	if IsHex("FF") == false {
		t.Error("FF should be a valid hex")
	}
	if IsHex("0xff") == false {
		t.Error("0xff should be a valid hex")
	}
	if IsHex("ff") == false {
		t.Error("ff should be a valid hex")
	}
	if IsHex("0123456789abcdef") == false {
		t.Error("0123456789abcdef should be a valid hex")
	}
	if IsHex("0123456789ABCDEF") == false {
		t.Error("0123456789ABCDEF should be a valid hex")
	}
	if IsHex("g") == true {
		t.Error("g should not be a valid hex")
	}
	if IsHex("G") == true {
		t.Error("g should not be a valid hex")
	}
}

func TestString2Int64(t *testing.T) {
	answer, err := String2Int64("")
	if err == nil {
		t.Error("Empty string should return an error")
	}
	if answer != 0 {
		t.Error("Empty string should return zero value")
	}

	answer, err = String2Int64("0")
	if err != nil {
		t.Error("Unexpected error")
	}
	if answer != 0 {
		t.Error("Received wrong output")
	}

	answer, err = String2Int64("1")
	if err != nil {
		t.Error("Unexpected error")
	}
	if answer != 1 {
		t.Error("Received wrong output")
	}

	answer, err = String2Int64("9223372036854775807")
	if err != nil {
		t.Error("Unexpected error")
	}
	if answer != 9223372036854775807 {
		t.Error("Received wrong output")
	}

	answer, err = String2Int64("-9223372036854775808")
	if err != nil {
		t.Error("Unexpected error")
	}
	if answer != -9223372036854775808 {
		t.Error("Received wrong output")
	}
}
