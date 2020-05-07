package evidence

import (
	"strings"

	"github.com/brentp/bix"
	"github.com/liserjrqlxue/goUtil/stringsUtil"
)

var (
	pm1Ext = 3
)

// PM1
func CheckPM1(item map[string]string, tbx *bix.Bix) string {
	if !isMissenseIndel.MatchString(item["Function"]) {
		return "0"
	}
	var dbNSFP = item["Interpro_domain"]
	var pfam = item["pfamId"]
	var flag bool

	for _, k := range strings.Split(dbNSFP, ";") {
		if pm1InterproDomain[k] {
			flag = true
		}
	}
	for _, k := range strings.Split(pfam, ";") {
		if pm1PfamId[k] {
			flag = true
		}
	}
	if !flag {
		var chr = strings.Replace(item["#Chr"], "chr", "", 1)
		var start = stringsUtil.Atoi(item["Start"])
		var end = stringsUtil.Atoi(item["Stop"])
		n := countBix(tbx, chr, start-pm1Ext, end+pm1Ext)
		if n >= 2 {
			flag = true
		}
	}
	if flag {
		return "1"
	} else {
		return "0"
	}
}

func ComparePM1(item map[string]string, tbx *bix.Bix) {
	rule := "PM1"
	val := CheckPM1(item, tbx)
	if val != item[rule] {
		PrintConflict(item, rule, val, "Interpro_domain", "pfamId")
	}
}
