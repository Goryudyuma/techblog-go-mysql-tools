package modelutil

import "reflect"

func ListDBTagString(variable any) (ret []string) {
	v := reflect.Indirect(reflect.ValueOf(variable))
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := range t.NumField() {
			ti := t.Field(i)
			tag := ti.Tag.Get("db")
			if tag != "" && tag != "-" {
				tags := ListDBTagString(reflect.New(ti.Type).Interface())

				if len(tags) > 0 {
					for _, retTag := range tags {
						ret = append(ret, tag+"."+retTag)
					}
				} else {
					ret = append(ret, tag)
				}
			}
		}
	}

	return
}
