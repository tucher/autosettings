package autosettings

import (
	"encoding/json"
	"flag"
	"github.com/kardianos/osext"
	"io/ioutil"
	"log"
	"reflect"
	"strings"
)

type Defaultable interface {
	Default() Defaultable
}

func ReadConfig(tmpl Defaultable) {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	d := tmpl.Default()
	jsonBlob, err := ioutil.ReadFile(folderPath + "/settings.json")
	if err == nil {
		if err := json.Unmarshal(jsonBlob, d); err != nil {
			log.Fatal(err)
		}
	}

	cloneValue(d, tmpl)

	v := reflect.ValueOf(tmpl).Elem()
	// defaultV := reflect.ValueOf(tmpl.Default()).Elem()
	tp := v.Type()
	for i := 0; i < tp.NumField(); i++ {
		fieldName := tp.Field(i).Name
		fieldType := tp.Field(i).Type
		fieldTag := tp.Field(i).Tag.Get("autosettings")
		fieldDefaultValue := v.Field(i)

		switch fieldType.Kind() {
		case reflect.Bool:
			flag.BoolVar(v.Field(i).Addr().Interface().(*bool), strings.ToLower(fieldName), fieldDefaultValue.Bool(), fieldTag)
		case reflect.Int64:
			flag.Int64Var(v.Field(i).Addr().Interface().(*int64), strings.ToLower(fieldName), fieldDefaultValue.Int(), fieldTag)
		case reflect.Int:
			flag.IntVar(v.Field(i).Addr().Interface().(*int), strings.ToLower(fieldName), fieldDefaultValue.Interface().(int), fieldTag)
		case reflect.String:
			flag.StringVar(v.Field(i).Addr().Interface().(*string), strings.ToLower(fieldName), fieldDefaultValue.String(), fieldTag)
		case reflect.Float64:
			flag.Float64Var(v.Field(i).Addr().Interface().(*float64), strings.ToLower(fieldName), fieldDefaultValue.Float(), fieldTag)
		}
	}
	flag.Parse()
}

func cloneValue(source interface{}, destin interface{}) {
	x := reflect.ValueOf(source)
	if reflect.ValueOf(destin).Kind() != reflect.Ptr {
		return
	}
	if x.Kind() == reflect.Ptr {
		reflect.ValueOf(destin).Elem().Set(x.Elem())
	} else {
		reflect.ValueOf(destin).Elem().Set(x)
	}
}
