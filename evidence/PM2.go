package evidence

var (
	PM2ARAFThreshold  = 0.005
	PM2ADAFThreshold  = 0.0
	PM2ADAFThreshold2 = 0.00001
)

var PM2AFList = []string{
	"ESP6500 AF",
	"1000G AF",
	"ExAC AF",
	"ExAC EAS AF",
	"GnomAD AF",
	"GnomAD EAS AF",
}

// PM2
func CheckPM2(item map[string]string) string {
	inherit := item["ModeInheritance"]
	if isARDRXLPRDDNA.MatchString(inherit) || inherit == "" {
		if CheckAFAllLowThen(item, PM2AFList, PM2ARAFThreshold, true) {
			return "1"
		} else {
			return "0"
		}
	} else if isADPDYL.MatchString(inherit) {
		if bs2GeneList[item["Gene Symbol"]] {
			if CheckAFAllLowThen(item, PM2AFList, PM2ADAFThreshold2, true) {
				return "1"
			}
		} else {
			if CheckAFAllLowThen(item, PM2AFList, PM2ADAFThreshold, true) {
				return "1"
			}
		}
		return "0"
	}
	return "0"
}

func ComparePM2(item map[string]string) {
	rule := "PM2"
	val := CheckPM2(item)
	if val != item[rule] {
		PrintConflict(item, rule, val, append(PM2AFList, "ModeInheritance")...)
	}
}
