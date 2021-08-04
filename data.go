package cuppago

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func SendAndLoad(url string, data interface{}) string {
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	text := string(body)
	return text
}

func MergeObjects(values ...map[string]interface{}) map[string]interface{} {
	var result map[string]interface{}
	for i := 0; i < len(values); i++ {
		var temp = values[i]
		for k, v := range temp {
			if _, ok := temp[k]; ok {
				result[k] = v
			}
		}
	}
	return result
}

func FormToMap(r *http.Request) map[string]interface{} {
	r.ParseForm()
	urlValues := r.Form
	data := make(map[string]interface{})
	for key, value := range urlValues {
		data[key] = value[0]
	}
	return data
}

func URLToMap(urlString string) map[string]interface{} {
	m, _ := url.ParseQuery(urlString)
	data := make(map[string]interface{})
	for key, value := range m {
		data[key] = value[0]
	}
	return data
}

func MapListToStructure(list []map[string]interface{}, s interface{}) {
	for i := 0; i < len(list); i++ {
		MapToStructure(list[i], s)
	}
}

func MapToStructure(m map[string]interface{}, s interface{}) {

}

func StructureToMap(element interface{}) map[string]interface{} {
	if IsNil(element) {
		return make(map[string]interface{})
	}
	tmpMap := JSONDecode(JSONEncode(element, false), false).(map[string]interface{})
	return tmpMap
}

func JSONToMap(element interface{}) map[string]interface{} { return StructureToMap(element) }

func Value(data interface{}, path string, defaultValue interface{}) interface{} {
	if IsNil(data) {
		return defaultValue
	}
	if data != nil && path == "" {
		return data
	}
	pathArray := strings.Split(strings.Trim(path, ""), ".")
	element := data
	for i := 0; i < len(pathArray); i++ {
		if element == nil {
			break
		}
		pathSlide := pathArray[i]
		if reflect.TypeOf(element).String() == "map[string]interface {}" {
			element = element.(map[string]interface{})[pathSlide]
		} else if reflect.TypeOf(element).String() == "[]interface {}" {
			index, _ := strconv.Atoi(pathSlide)
			len := len(element.([]interface{}))
			if index > len-1 {
				element = nil
				break
			}
			element = element.([]interface{})[index]
		} else {
			element = StructureToMap(element)[pathSlide]
		}
	}
	if element != nil && element != "" {
		tmp := JSONEncode(element, false)
		if tmp == "{}" || tmp == "[]" || tmp == "" {
			element = ""
		}
	}
	if (element == nil || element == "") && defaultValue != nil {
		element = defaultValue
	}
	return element
}

func IsNil(a interface{}) bool {
	defer func() { recover() }()
	return a == nil || reflect.ValueOf(a).IsNil()
}

func GetRequestBody(r *http.Request) string {
	body, _ := ioutil.ReadAll(r.Body)
	return string(body)
}

// Example: var user User; utils.MapToStruct(result, &user)
func MapToStruct(data map[string]interface{}, structPointer interface{}) {
	s := reflect.ValueOf(structPointer).Elem()
	for key, value := range data {
		sf := s.FieldByName(Camelize(key))
		defer func() {
			if recover() != nil {
			}
		}()
		sfKind := sf.Kind()
		if sfKind == reflect.String {
			sf.SetString(value.(string))
		} else if sfKind == reflect.Int || sfKind == reflect.Int8 || sfKind == reflect.Int16 || sfKind == reflect.Int32 || sfKind == reflect.Int64 {
			valConv, err := strconv.ParseInt(value.(string), 10, 64)
			if err != nil {
				break
			}
			sf.SetInt(valConv)
		} else if sfKind == reflect.Uint {
			valConv, err := strconv.ParseUint(value.(string), 10, 64)
			if err != nil {
				break
			}
			sf.SetUint(valConv)
		} else if sfKind == reflect.Float32 || sfKind == reflect.Float64 {
			valConv, err := strconv.ParseFloat(value.(string), 64)
			if err != nil {
				break
			}
			sf.SetFloat(valConv)
		} else if sfKind == reflect.Bool {
			valConv, err := strconv.ParseBool(value.(string))
			if err != nil {
				break
			}
			sf.SetBool(valConv)
		}
	}
}
