package evidence

import "github.com/liserjrqlxue/goUtil/textUtil"

var BA1AFThreshold = 0.05
var BA1AFList = []string{
	"ESP6500 AF",
	"1000G AF",
	"ExAC AF",
	"GnomAD AF",
	"ExAC EAS AF",
	"GnomAD EAS AF",
}

func CreateBA1Exception(fileName string) (BA1Exception map[string]bool) {
	BA1Exception = make(map[string]bool)
	for _, key := range textUtil.File2Array(fileName) {
		BA1Exception[key] = true
	}
	return
}

// ture	:	"1"
// flase:	"0"
func CheckBA1(item map[string]string, BA1Exception map[string]bool) string {
	var key = item["Transcript"] + ":" + item["cHGVS"]
	if BA1Exception[key] {
		return "0"
	}
	if CheckAFAllLowThen(item, BA1AFList, BA1AFThreshold, true) {
		return "0"
	} else {
		return "1"
	}
}

func CompareBA1(item map[string]string, BA1Exception map[string]bool, lostOnly bool) {
	rule := "BA1"
	val := CheckBA1(item, BA1Exception)
	if val != item[rule] {
		if item[rule] == "0" && val == "" {
		} else {
			if lostOnly && val != "1" {
				return
			}
			PrintConflict(item, rule, val, BA1AFList...)
		}
	}
}
