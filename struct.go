package moist

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"reflect"
)

// JSON文件转换成Map
func JsonFileToMap(jsonFile string) (result interface{}, err error) {
	b, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		log.Println("JsonFileToMap: ", err)
		return result, err
	}

	s := string(b)

	// 尝试转换成单个JSON对象
	obj, err := JsonToMap(s)
	if err != nil {
		// 尝试转换成JSON数组
		arr, err2 := JsonArrayToMap(s)
		if err2 != nil {

			log.Println("JsonFileToMap: ", err2)
			return arr, err2
		} else {
			err = nil
		}

		result = arr
	} else {
		result = obj
	}

	if err != nil {
		log.Println("JsonFileToMap: ", err)
	}

	return result, err
}

// JSON字符串转换成Map
func JsonToMap(sJson string) (result map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(sJson), &result)
	if err != nil {
		log.Println("JsonToMap: ", err)
		return result, err
	}

	return result, nil
}

// JSON(数组)字符串转换成Map
func JsonArrayToMap(sJson string) (result []map[string]interface{}, err error) {
	err = json.Unmarshal([]byte(sJson), &result)
	if err != nil {
		log.Println("JsonArrayToMap: ", err)
		return result, err
	}

	return result, nil
}

// 将struct转化为map
// 》使用json
func StrctToMap(s interface{}) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	j, err := json.Marshal(s)
	if err != nil {
		log.Print("Error occurs when convert struct to map. function: StrctToMap() -> json.Marshal()")
		return m, err
	}
	err = json.Unmarshal(j, &m)
	if err != nil {
		log.Print("Error occurs when convert struct to map. function: StrctToMap() -> json.Unmarshal()")
		return m, err
	}

	return m, nil
}

// 将string类型Key的Map转化为interface{}类型Key的Map
// 》interface{}类型Key的Map用于go-flutter插件
func ToIfKeyMap(m map[string]interface{}) (result map[interface{}]interface{}, err error) {

	result = make(map[interface{}]interface{})
	for k, v := range m {
		kindTmp := reflect.ValueOf(v).Kind()

		if kindTmp == reflect.Slice {
			vSlice, ok := v.([]interface{})

			if !ok {
				return result, errors.New("Error occurs when reflect slice.")
			}

			// 用于存放Map的切片
			var valSlice []interface{}
			for _, subVal := range vSlice {

				// 切片元素为Map类型
				if reflect.ValueOf(subVal).Kind() == reflect.Map {
					itemMap, itemOk := subVal.(map[string]interface{})
					if itemOk {
						ifKeyMap, err := ToIfKeyMap(itemMap)

						if err != nil {
							return nil, err
						}

						valSlice = append(valSlice, ifKeyMap)
					}
				}
			}

			result[k] = valSlice
		} else if kindTmp == reflect.Map {

			tmpMap := make(map[interface{}]interface{})

			for k2, v2 := range v.(map[string]interface{}) {
				tmpMap[k2] = v2
			}

			// bTmp, err := json.Marshal(v)
			// if err != nil {
			// 	log.Println(err)
			// }

			// sTmp := string(bTmp)

			// tmp, err := JsonToMap(sTmp)

			// toMap(tmp, tmpMap)
			result[k] = tmpMap
		} else {
			result[k] = v
		}
	}

	return result, nil
}

// 将struct转化为interface{}类型Key的Map
// 》interface{}类型Key的Map用于go-flutter插件
func StrctToIfKeyMap(s interface{}) (result map[interface{}]interface{}, err error) {
	// 先转换成string类型Key的Map
	m, err := StrctToMap(s)
	if err != nil {
		log.Print("Error occurs when convert struct to map with interface key. function: StrctToIfKeyMap() -> StrctToMap()")
		return result, err
	}

	// 再转换成interface类型Key的Map
	result, err = ToIfKeyMap(m)
	if err != nil {
		log.Print("Error occurs when convert struct to map with interface key. function: StrctToIfKeyMap() -> ToIfKeyMap()")
		return result, err
	}

	return result, nil
}

// func ToMap(tests []interface) {

// 	l := make([]map[string]interface{}, 0)
// 	for _, t := range tests {
// 		elem := reflect.ValueOf(&t).Elem()
// 		relType := elem.Type()

// 		m := make(map[string]interface{}, 1)
// 		for i := 0; i < relType.NumField(); i++ {
// 			m[relType.Field(i).Name] = elem.Field(i).Interface()
// 		}

// 		l = append(l, m)
// 	}
// }
