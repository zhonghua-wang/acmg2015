package evidence

import "regexp"

// regexp
var (
	isARXRNA = regexp.MustCompile(`AR|XR|NA`)
	isADXD   = regexp.MustCompile(`AD|XD`)
)

var (
	PM2ARAFThreshold = 0.05
	PM2ADAFThreshold = 0.0
)

var PM2AFList = []string{
	"ESP6500 AF",
	"1000G AF",
	"ExAC AF",
	//"ExAC EAS",
	"GnomAD AF",
	//"GnomAD EAS AF",
}

// PM2
func CheckPM2(item map[string]string) string {
	inherit := item["OMIM inheritance"]
	if isARXRNA.MatchString(inherit) {
		if CheckAFAllLowThen(item, PM2AFList, PM2ARAFThreshold, true) {
			return "1"
		} else {
			return "0"
		}
	} else if isADXD.MatchString(inherit) {
		if CheckAFAllLowThen(item, PM2AFList, PM2ADAFThreshold, true) {
			return "1"
		} else {
			return "0"
		}
	}
	return ""
}

func ComparePM2(item map[string]string) {
	rule := "PM2"
	val := CheckPM2(item)
	if val != item[rule] {
		PrintConflict(item, rule, val, append(PM2AFList, "OMIM inheritance")...)
	}
}
