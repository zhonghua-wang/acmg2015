package evidence

import (
	"fmt"
	"os"
)

var pm4Func = map[string]bool{
	"cds-del":   true,
	"cds-ins":   true,
	"cds-indel": true,
	"stop-loss": true,
}

// ture	:	"1"
// flase:	"0"
// nil	:	""
func CheckPM4(item map[string]string) string {
	if pm4Func[item["Function"]] {
		if item["RepeatTag"] == "" || item["RepeatTag"] == "." {
			return "1"
		}
	} else {
		return "0"
	}
	return ""
}

func ComparePM4(item map[string]string) {
	rule := "PM4"
	val := CheckPM4(item)
	if val != item[rule] {
		if item[rule] == "0" && val == "" {
		} else {
			fmt.Fprintf(
				os.Stderr,
				"Conflict %s:[%s] vs [%s]\tMutationName[%s]\tFunction[%s]\tRepeatTag[%s]\n",
				rule,
				val,
				item[rule],
				item["MutationName"],
				item["Function"],
				item["RepeatTag"],
			)
		}
	}
}
