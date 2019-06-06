package acmg2015

func AddACMG2015(inputData map[string]string) map[string]string {
	//var LOFIntoleranceGeneList = simple_util.JsonFile2MapInt("db/LOFIntoleranceGeneList.json")
	var info = make(map[string]string)
	//info["PVS1"] = evidence.CheckPVS1(inputData, LOFIntoleranceGeneList)
	inputData["ACMG"] = PredACMG2015(info)
	return inputData
}
