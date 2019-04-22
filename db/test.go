package main

import (
	"encoding/json"
	"fmt"
	"github.com/liserjrqlxue/acmg2015/evidence"
	"github.com/liserjrqlxue/simple-util"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// os
var (
	ex, _  = os.Executable()
	exPath = filepath.Dir(ex)
	pSep   = string(os.PathSeparator)
	//dbPath       = exPath + pSep + "db" + pSep
	//templatePath = exPath + pSep + "template" + pSep
)

var (
	clinvarAnno = "clinvar.2019-02-19.vcf.gz.2019-02-26.anno.bed.update"
	hgmdAnno    = "hgmd_pro_2018.3_hg19.vcf.2019-02-27.anno.bed.update"
	exacLOF     = "LOF_Intolerance.ExAC.lst"
)

// colname
var (
	clinvarCol = "ClinVar Significance"
	hgmdCol    = "HGMD Pred"
)

func main() {

	// spec.var.list anno clinvar hgmd
	if false {
		mutList := simple_util.File2Array("spec.var.list")
		var annotation = make(map[string]map[string]string)
		var annoLite = make(map[string]map[string]string)
		loadClinVar, _ := simple_util.File2MapArray(clinvarAnno, "\t", nil)
		loadHGMD, _ := simple_util.File2MapArray(hgmdAnno, "\t", nil)
		for _, item := range loadClinVar {
			key := item["MutationName"]
			annotation[key] = item
		}
		for _, item := range loadHGMD {
			key := item["MutationName"]
			info, ok := annotation[key]
			if ok {
				info["HGMD Pred"] = item[hgmdCol]
			} else {
				annotation[key] = item
			}
		}
		var columns = []string{
			"#Chr",
			"Start",
			"Stop",
			"Ref",
			"Call",
			"MutationName",
			"VarType",
			clinvarCol,
			hgmdCol,
		}
		file, err := os.Create("sepc.var.list.txt")
		simple_util.CheckErr(err)
		defer file.Close()

		fmt.Fprintln(file, strings.Join(columns, "\t"))

		for _, key := range mutList {
			var item = make(map[string]string)
			var array []string
			for _, col := range columns {
				item[col] = annotation[key][col]
				array = append(array, item[col])
			}
			fmt.Fprintln(file, strings.Join(array, "\t"))
			annoLite[key] = item
		}
		jsonByte, err := simple_util.JsonIndent(annoLite, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "spec.var.list.json")
	}

	// build PVS1 db
	if false {
		var LOFIntoleranceGeneList = make(map[string]int)

		// load ClinVar
		var ClinVarGeneList = evidence.FindLOFIntoleranceGeneList(clinvarAnno, clinvarCol, evidence.IsClinVarPLP)
		for key, val := range ClinVarGeneList {
			if val > 0 {
				LOFIntoleranceGeneList[key] += val
			}
		}
		jsonByte, err := json.MarshalIndent(ClinVarGeneList, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarLOFIntoleranceGeneList.json")

		// load HGMD
		var HGMDGeneList = evidence.FindLOFIntoleranceGeneList(hgmdAnno, hgmdCol, evidence.IsHgmdDM)
		for key, val := range HGMDGeneList {
			if val > 0 {
				LOFIntoleranceGeneList[key] += val
			}
		}
		jsonByte, err = json.MarshalIndent(HGMDGeneList, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDLOFIntoleranceGeneList.json")

		var exacList = simple_util.File2Array(exacLOF)
		for _, key := range exacList {
			LOFIntoleranceGeneList[key]++
		}
		jsonByte, err = json.MarshalIndent(LOFIntoleranceGeneList, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "LOFIntoleranceGeneList.json")

		// load LOFIntoleranceGeneList
		b, err := ioutil.ReadFile("LOFIntoleranceGeneList.json")
		simple_util.CheckErr(err)
		err = json.Unmarshal(b, &LOFIntoleranceGeneList)
		simple_util.CheckErr(err)
	}
	// test PVS1
	if false {
		var item = map[string]string{
			"Function":    "splice-3",
			"Gene Symbol": "ACD",
		}
		LOFIntoleranceGeneList := simple_util.JsonFile2MapInt("LOFIntoleranceGeneList.json")
		fmt.Println(evidence.CheckPVS1(item, LOFIntoleranceGeneList))
	}

	// build PS1 db
	if false {
		// load ClinVar
		var ClinVarMissense, ClinVarPHGVSlist = evidence.FindPathogenicMissense(clinvarAnno, clinvarCol, evidence.IsClinVarPLP)
		jsonByte, err := simple_util.JsonIndent(ClinVarMissense, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarPathogenicMissense.json")
		jsonByte, err = simple_util.JsonIndent(ClinVarPHGVSlist, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "ClinVarPHGVSlist.json")

		// load HGMD
		var HGMDMissense, HGMDPHGVSlist = evidence.FindPathogenicMissense(clinvarAnno, hgmdCol, evidence.IsHgmdDM)
		jsonByte, err = simple_util.JsonIndent(HGMDMissense, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDPathogenicMissense.json")
		jsonByte, err = simple_util.JsonIndent(HGMDPHGVSlist, "", "\t")
		simple_util.CheckErr(err)
		simple_util.Json2file(jsonByte, "HGMDPHGVSlist.json")
	}
	// test PS1
	if true {
		var item = map[string]string{
			"MutationName": "NM_000142.4(FGFR3): c.1138G>A (p.Gly380Arg)",
			"Transcript":   "NM_000142.4",
			"pHGVS":        "p.G380R | p.Gly380Arg",
		}
		var ClinVarMissense = simple_util.JsonFile2MapInt("ClinVarPathogenicMissense.json")
		var ClinVarPHGVSlist = simple_util.JsonFile2MapInt("ClinVarPHGVSlist.json")
		var HGMDMissense = simple_util.JsonFile2MapInt("HGMDPathogenicMissense.json")
		var HGMDPHGVSlist = simple_util.JsonFile2MapInt("HGMDPHGVSlist.json")
		fmt.Println(evidence.CheckPS1(item, ClinVarMissense, ClinVarPHGVSlist, HGMDMissense, HGMDPHGVSlist))
	}
}
