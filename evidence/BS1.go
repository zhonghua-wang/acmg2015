package evidence

import (
	"fmt"
	"os"
)

var BS1MorbidityThreshold = 0.01
var BS1AFList = []string{
	"1000G AF",
	"ExAC AF",
	"GnomAD AF",
}

// ture	:	"1"
// flase:	"0"
// nil	:	""
func CheckBS1(item map[string]string) string {
	if CheckAFAllLowThen(item, BS1AFList, BS1MorbidityThreshold, true) {
		return "0"
	} else {
		return "1"
	}
}

func CompareBS1(item map[string]string, lostOnly bool) {
	rule := "BS1"
	val := CheckBS1(item)
	if val != item[rule] {
		if item[rule] == "0" && val == "" {
		} else {
			if lostOnly && val != "1" {
				return
			}
			fmt.Fprintf(
				os.Stderr,
				"Conflict %s:[%s] vs [%s]\t%s[%s]\n",
				rule,
				val,
				item[rule],
				"MutationName",
				item["MutationName"],
			)
			for _, key := range BS1AFList {
				fmt.Fprintf(os.Stderr, "\t%s:[%s]\n", key, item[key])
			}
		}
	}
}
