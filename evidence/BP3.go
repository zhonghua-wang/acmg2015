package evidence

import (
	"fmt"
	"os"
	"regexp"
)

var BP3Func = map[string]bool{
	"cds-del":   true,
	"cds-ins":   true,
	"cds-indel": true,
}

var (
	isRepeatSeq  = regexp.MustCompile(`\([ACGT]+\)n`)
	isDeepIntron = regexp.MustCompile(`intron|span|splice[+-]10|splice[+-]20`)
)

// ture	:	"1"
// flase:	"0"
// nil	:	""
func CheckBP3(item map[string]string) string {
	if BP3Func[item["Function"]] {
		if item["RepeatTag"] == "" || item["RepeatTag"] == "." {
			return "0"
		} else {
			return "1"
		}
	} else if isDeepIntron.MatchString(item["Function"]) &&
		isRepeatSeq.MatchString(item["RepeatTag"]) &&
		item["VarType"] != "snv" && item["VarType"] != "ref" {
		return "1"
	} else {
		return "0"
	}
	return ""
}

func CompareBP3(item map[string]string) {
	rule := "BP3"
	val := CheckBP3(item)
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
			for _, key := range []string{"Function", "RepeatTag", "VarType"} {
				fmt.Fprintf(os.Stderr, "\t%30s:[%s]\n", key, item[key])
			}
		}
	}
}
