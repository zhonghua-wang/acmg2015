package evidence

import "regexp"

// regexp
var (
	isARXLNA = regexp.MustCompile(`AR|XL|NA`)
	isAD     = regexp.MustCompile(`AD`)
)

var (
	PM2ARAFThreshold = 0.005
	PM2ADAFThreshold = 0.0
)

var PM2AFList = []string{
	"ESP6500 AF",
	"1000G AF",
	"ExAC AF",
	"ExAC EAS",
	"GnomAD AF",
	"GnomAD EAS AF",
}

// PM2
func CheckPM2(item map[string]string) string {
	inherit := item["ModeInheritance"]
	if isARXLNA.MatchString(inherit) || inherit == "" {
		if CheckAFAllLowThen(item, PM2AFList, PM2ARAFThreshold, true) {
			return "1"
		} else {
			return "0"
		}
	} else if isAD.MatchString(inherit) {
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
		PrintConflict(item, rule, val, append(PM2AFList, "ModeInheritance")...)
	}
}
