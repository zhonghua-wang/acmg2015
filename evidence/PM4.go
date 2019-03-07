package evidence

var pm4Func = map[string]bool{
	"cds-del":   true,
	"cds-ins":   true,
	"cds-indel": true,
	"stop-loss": true,
}

// ture	:	"1"
// flase:	"0"
// nil	:	""
func CheckPM4(item map[string]string) string {
	if pm4Func[item["Function"]] {
		if item["RepeatTag"] == "" || item["RepeatTag"] == "." {
			return "1"
		}
	} else {
		return "0"
	}
	return ""
}
