package cuppago

import "testing"

func Test_PathData(t *testing.T) {
	result := PathData("path1/path2/?var1=1&var2=male")
	Log("RESULT", result)
}