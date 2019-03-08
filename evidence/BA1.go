package evidence

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
