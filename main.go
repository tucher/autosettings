package autosettings

import (
	"encoding/json"
	"fmt"
	"github.com/kardianos/osext"
	"io/ioutil"
	"log"
	"reflect"
)

type Defaultable interface {
	Default() Defaultable
}

func ReadConfig(tmpl Defaultable) {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}

	jsonBlob, err := ioutil.ReadFile(folderPath + "/settings.json")
	if err != nil {
		log.Fatal(err)
	}

	d := tmpl.Default()
	if err := json.Unmarshal(jsonBlob, d); err != nil {
		log.Fatal(err)
	}
	cloneValue(d, tmpl)
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
