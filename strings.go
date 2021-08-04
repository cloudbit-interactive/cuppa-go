/* REQUIRES
> go get github.com/google/uuid
> go get golang.org/x/crypto/bcrypt
*/

package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"regexp"
	"strings"
	"unicode"
)

func PathData(path string)map[string]interface{}{
	var result = make(map[string]interface{})
	// get base
		var base = path
		if(strings.Index(base, "?") != -1){ base = base[0:strings.Index(base, "?")] }
		if(strings.Index(base, "#") != -1){ base = base[0:strings.Index(base, "#")] }
		result["base"] = base;
		result["baseArray"] = strings.Split(base,"/")
	// domain
		var domain = strings.Replace(base, "https://", "", -1)
			domain = strings.Replace(domain, "http://", "", -1)
			domain = domain[0: strings.Index(domain,"/")]
			result["domain"] = domain
	// protocol
		result["protocol"] = "http"
		if(strings.Index(path, "https://") != -1){ result["protocol"] = "https" }
	// data
		var dataStr = path
		if(strings.Index(dataStr, "?") != -1 || strings.Index(dataStr, "#") != -1) {
			if(strings.Index(dataStr, "?") != -1){ dataStr = dataStr[strings.Index(dataStr, "?")+1: len(dataStr)] }
			if(strings.Index(dataStr, "#") != -1){ dataStr = dataStr[strings.Index(dataStr, "#")+1: len(dataStr)] }
			var data = make(map[string]interface{})
			var dataArray = strings.Split(dataStr,"&")
			for i := 0; i < len(dataArray); i++{
				var parts = strings.Split(dataArray[i],"=")
				if(parts[0] != ""){
					if(parts[1] != ""){
						data[parts[0]] = parts[1]
					}else{
						data[parts[0]] = ""
					}
				}
			}
			result["data"] = data
		}
	return result
}

func JSONEncode(value interface{}, base64Encode bool)string{
	result := ""
	bites, err := json.Marshal(value)
	if err == nil { result = string(bites) }
	if base64Encode == true { result = Base64Encode(result) }
	return result;
}

/*
		data := utils.JSONDecode(string, false)
 */
func JSONDecode(value string, base64Decode bool)interface{}{
	if base64Decode == true { value = Base64Decode(value) }
	var result interface{}
	json.Unmarshal([]byte(value), &result)
	return result
}

func Base64Encode(string string)string{
	result := base64.StdEncoding.EncodeToString([]byte(string))
	return result
}

func Base64Decode(value string)string{
	result := ""
	data, err := base64.StdEncoding.DecodeString(value)
	if err == nil { result = string(data) }
	return result
}

/* Parse any value to String */
func String(value interface{}) string{
	result := fmt.Sprint(value)
	if result == "<nil>" {
		result = ""
	}
	return result
}

func UUID() string {
	uuid, _ := uuid.NewUUID()
	return uuid.String()
}

func Hash(value string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(value), bcrypt.DefaultCost)
	return string(bytes)
}

func CompareHash(hash string, value string) bool {
	byteHash := []byte(hash)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(value))
	if err != nil { return false }
	return true
}

func InterfaceToString(value interface{})string{
	str := fmt.Sprintf("%v", value)
	return str
}

func ReplaceNotCase(value string, search string, replace string)string{
	return strings.Trim(regexp.MustCompile(`(?i)`+search).ReplaceAllString(value, replace), " ")
}

func Camelize(in string) string {
	in = strings.ReplaceAll(in, " ", "_");
	in = strings.ReplaceAll(in, "-", "_");
	runes := []rune(in)
	var out []rune
	for i, r := range runes {
		if r == '_' {
			continue
		}
		if i == 0 || runes[i-1] == '_' {
			out = append(out, unicode.ToUpper(r))
			continue
		}
		out = append(out, r)
	}
	return string(out)
}