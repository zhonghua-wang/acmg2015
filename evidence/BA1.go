package evidence

var BA1AFThreshold = 0.05
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
	if !CheckAFAllLowThen(item, BA1AFList, BA1AFThreshold, true) {
		return "1"
	} else {
		return "0"
	}
}

func CompareBA1(item map[string]string, lostOnly bool) {
	rule := "BA1"
	val := CheckBA1(item)
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
