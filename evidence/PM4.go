package evidence

var PM4Func = map[string]bool{
	"cds-del":   true,
	"cds-ins":   true,
	"cds-indel": true,
	"stop-loss": true,
}

// ture	:	"1"
// flase:	"0"
// nil	:	""
func CheckPM4(item map[string]string) string {
	if PM4Func[item["Function"]] {
		if item["RepeatTag"] == "" || item["RepeatTag"] == "." {
			return "1"
		} else {
			return "0"
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
		PrintConflict(item, rule, val, []string{"Function", "RepeatTag"})
	}
}
