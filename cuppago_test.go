package cuppago

import "testing"

func Test_PathData(t *testing.T) {
	result1 := PathData("")
	result2 := PathData("path1/path2")
	result3 := PathData("path1/path2/")
	result4 := PathData("path1/path2/?var1=1&var2=male&var3")
	result5 := PathData("path1/path2/?&var1=1&var2=male&var3=")
	Log("result1", result1)
	Log("result2", result2)
	Log("result3", result3)
	Log("result4", result4)
	Log("result5", result5)
}