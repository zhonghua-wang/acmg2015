package evidence

import (
	simple_util "github.com/liserjrqlxue/simple-util"
	"regexp"
	"strconv"
)

var BP3Func = map[string]bool{
	"cds-del":   true,
	"cds-ins":   true,
	"cds-indel": true,
}

var (
	isRepeatSeq  = regexp.MustCompile(`\([ACGT]+\)n`)
	isDeepIntron = regexp.MustCompile(`intron|span|splice[+-]10|splice[+-]20`)
	repeatSeq    = regexp.MustCompile(`c\..*\[(\d+)>\d+]`)
)

// ture	:	"1"
// flase:	"0"
// nil	:	""
func CheckBP3(item map[string]string) string {
	if BP3Func[item["Function"]] {
		if item["RepeatTag"] == "" || item["RepeatTag"] == "." {
			return "0"
		} else {
			subMatch := repeatSeq.FindStringSubmatch(item["cHGVS"])
			if len(subMatch) > 1 {
				dupCount, err := strconv.Atoi(subMatch[1])
				simple_util.CheckErr(err)
				if dupCount < 10 {
					return "0"
				} else {
					return "1"
				}
			} else {
				return "1"
			}
		}
	} else if isDeepIntron.MatchString(item["Function"]) &&
		isRepeatSeq.MatchString(item["RepeatTag"]) &&
		item["VarType"] != "snv" && item["VarType"] != "ref" {
		subMatch := repeatSeq.FindStringSubmatch(item["cHGVS"])
		if len(subMatch) > 1 {
			dupCount, err := strconv.Atoi(subMatch[1])
			simple_util.CheckErr(err)
			if dupCount < 10 {
				return "0"
			} else {
				return "1"
			}
		} else {
			return "1"
		}
	} else {
		return "0"
	}
}

func CompareBP3(item map[string]string) {
	rule := "BP3"
	val := CheckBP3(item)
	if val != item[rule] {
		if item[rule] == "0" && val == "" {
		} else {
			PrintConflict(item, rule, val, "Function", "RepeatTag", "VarType")
		}
	}
}
