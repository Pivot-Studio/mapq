package mapq

import "encoding/json"

// QueryJson 查询json是否满足条件
func QueryJson(jsons string, query string) (bool, error) {
	data, err := JsonToMap(jsons)
	if err != nil {
		return false, err
	}
	return QueryMap(data, query)
}

// JsonToMap json转map
func JsonToMap(jsons string) (map[string]interface{}, error) {
	var data map[string]interface{}
	err := json.Unmarshal([]byte(jsons), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// JsonToList json转list
func JsonToList(jsons string) ([]map[string]interface{}, error) {
	var data []map[string]interface{}
	err := json.Unmarshal([]byte(jsons), &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// QueryListJsonContains 查询json数组，只要一个元素符合条件就返回true
func QueryListJsonContains(jsons string, query string) (bool, error) {
	data, err := JsonToList(jsons)
	if err != nil {
		return false, err
	}
	for _, v := range data {
		if ok, err := QueryMap(v, query); err == nil && ok {
			return true, nil
		}
	}
	return false, nil
}

// QueryMap 查询map
func QueryMap(data map[string]interface{}, query string) (bool, error) {
	p := &Parser{}
	n, err := p.Parse(query)
	if err != nil {
		return false, err
	}
	return n.Eval(data).(bool), nil
}

// RunQuery 查询
func RunQuery(root Node, data map[string]interface{}) (bool, error) {
	return root.Eval(data).(bool), nil
}
