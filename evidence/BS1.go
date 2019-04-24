package evidence

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
			PrintConflict(item, rule, val, BS1AFList...)
		}
	}
}
