package evidence

import (
	"fmt"
	"os"
)

var BA1AFThreshold = 0.5
var BA1AFList = []string{
	"ESP6500 AF",
	"1000G AF",
	"ExAC AF",
	//"ExAC EAS",
	"GnomAD AF",
	//"GnomAD EAS AF",
}

// ture	:	"1"
// flase:	"0"
// nil	:	""
func CheckBA1(item map[string]string) string {
	if CheckAFAllLowThen(item, BA1AFList, BA1AFThreshold, true) {
		return "0"
	} else {
		return "1"
	}
}

func CompareBA1(item map[string]string) {
	rule := "BA1"
	val := CheckBA1(item)
	if val != item[rule] {
		if item[rule] == "0" && val == "" {
		} else {
			fmt.Fprintf(
				os.Stderr,
				"Conflict %s:[%s] vs [%s]\t%s[%s]\n",
				rule,
				val,
				item[rule],
				"MutationName",
				item["MutationName"],
			)
			for _, key := range BA1AFList {
				fmt.Fprintf(os.Stderr, "\t%s:[%s]\n", key, item[key])
			}
		}
	}
}
