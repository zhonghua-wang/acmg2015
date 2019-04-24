package evidence

import (
	"strconv"
)

var PS4GWASdbORThreshold = 5.0

// ture	:	"1"
// flase:	"0"
// nil	:	""
func CheckPS4(item map[string]string) string {
	or := item["GWASdb_or"]
	if or != "" && or != "." {
		orFloat, err := strconv.ParseFloat(or, 32)
		if err == nil {
			if orFloat > PS4GWASdbORThreshold {
				return "1"
			} else {
				return "0"
			}
		}
	}
	return ""
}

func ComparePS4(item map[string]string) {
	rule := "PS4"
	val := CheckPS4(item)
	if val != item[rule] {
		if item[rule] == "0" && val == "" {
		} else {
			PrintConflict(item, rule, val, "GWASdb_or")
		}
	}
}
