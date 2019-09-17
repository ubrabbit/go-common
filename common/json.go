package common

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
)

func JsonEncode(source interface{}) []byte {
	data, err := json.Marshal(source)
	CheckPanic(err)
	return data
}

func JsonDecode(source string, data interface{}) {
	err := json.Unmarshal([]byte(source), data)
	CheckPanic(err)
}

func JsonDecodeSimple(source string) *simplejson.Json {
	js_obj, err := simplejson.NewJson([]byte(source))
	CheckPanic(err)
	return js_obj
}

func DeepCopyJSON(src map[string]interface{}, dest map[string]interface{}) {
	for key, value := range src {
		switch src[key].(type) {
		case map[string]interface{}:
			dest[key] = map[string]interface{}{}
			DeepCopyJSON(src[key].(map[string]interface{}), dest[key].(map[string]interface{}))
		default:
			dest[key] = value
		}
	}
}
