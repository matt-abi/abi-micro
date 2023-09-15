package runtime

import (
	"bytes"
	"log"
	"reflect"
	"strings"

	"github.com/ability-sh/abi-lib/dynamic"
	"github.com/ability-sh/abi-lib/errors"
	"github.com/ability-sh/abi-micro/micro"
)

type reflectExecutor struct {
	s      interface{}
	exec   map[string]reflect.Value
	scheme *micro.Scheme
}

var contextType = reflect.TypeOf((*micro.Context)(nil)).Elem()
var errorType = reflect.TypeOf((*error)(nil)).Elem()

func getName(name string, b *bytes.Buffer) string {
	b.Reset()
	for i, r := range name {
		if r >= 'A' && r <= 'Z' {
			if i != 0 {
				b.WriteRune('/')
			}
			b.WriteRune(r + 32)
		} else {
			b.WriteRune(r)
		}
	}
	b.WriteString(".json")
	return b.String()
}

func eachField(tp reflect.Type, fn func(name string, required bool, t reflect.StructField) bool) {
	count := tp.NumField()

	for i := 0; i < count; i++ {

		tf := tp.Field(i)

		name := tf.Tag.Get("name")

		if name == "" {
			name = tf.Tag.Get("json")
		}

		if name == "-" {
			continue
		}

		if tf.Type.Kind() == reflect.Struct {
			eachField(tf.Type, fn)
			continue
		}

		if name == "" {
			continue
		}

		ns := strings.Split(name, ",")
		name = ns[0]
		required := true

		if len(ns) > 1 && ns[1] == "omitempty" {
			required = false
		}

		if !fn(name, required, tf) {
			break
		}

	}
}

func typeToSchemeType(t reflect.Type, scheme *micro.Scheme) string {
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return "int32"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "uint32"
	case reflect.Int64:
		return "int64"
	case reflect.Uint64:
		return "uint64"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	case reflect.Slice:
		return typeToSchemeType(t.Elem(), scheme) + "[]"
	case reflect.Ptr:
		return typeToSchemeType(t.Elem(), scheme)
	case reflect.Struct:
		_, ok := scheme.Objects[t.Name()]
		if !ok {
			dst := &micro.SchemeObject{}
			typeToSchemeObject(t, dst, scheme)
			scheme.Objects[t.Name()] = dst
		}
		return t.Name()
	}
	return "any"
}

func typeFieldToSchemeItem(name string, required bool, tf reflect.StructField, scheme *micro.Scheme) *micro.SchemeField {
	rs := micro.SchemeField{Name: name, Required: required}
	rs.Type = typeToSchemeType(tf.Type, scheme)
	rs.Title = tf.Tag.Get("title")
	return &rs
}

func typeToSchemeObject(t reflect.Type, dst *micro.SchemeObject, scheme *micro.Scheme) {

	switch t.Kind() {
	case reflect.Ptr:
		typeToSchemeObject(t.Elem(), dst, scheme)
	case reflect.Struct:
		eachField(t, func(name string, required bool, tf reflect.StructField) bool {
			dst.Fields = append(dst.Fields, typeFieldToSchemeItem(name, required, tf, scheme))
			return true
		})
	}

}

func NewReflectExecutor(s interface{}) micro.Executor {
	rs := &reflectExecutor{s: s, exec: map[string]reflect.Value{}, scheme: &micro.Scheme{Objects: map[string]*micro.SchemeObject{}}}

	v := reflect.ValueOf(s)
	t := v.Type()

	num := v.NumMethod()

	rs.scheme.Name = t.Name()

	b := bytes.NewBuffer(nil)

	for i := 0; i < num; i++ {

		m := v.Method(i)

		inCount := m.Type().NumIn()

		if inCount != 2 {
			continue
		}

		outCount := m.Type().NumOut()

		if outCount != 2 {
			continue
		}

		if !m.Type().In(0).AssignableTo(contextType) {
			continue
		}

		inType := m.Type().In(1)

		if inType.Kind() != reflect.Ptr || inType.Elem().Kind() != reflect.Struct || !strings.HasSuffix(inType.Elem().Name(), "Task") {
			continue
		}

		if !m.Type().Out(1).AssignableTo(errorType) {
			continue
		}

		name := t.Method(i).Name

		n := getName(name, b)

		rs.exec[n] = m

		schemeItem := micro.SchemeItem{Name: n, Task: &micro.SchemeObject{}, Result: &micro.SchemeObject{}}

		typeToSchemeObject(inType, schemeItem.Task, rs.scheme)
		typeToSchemeObject(m.Type().Out(0), schemeItem.Result, rs.scheme)

		rs.scheme.Items = append(rs.scheme.Items, &schemeItem)

		log.Println("Executor", "=>", n, "=>", name)

	}

	return rs
}

func (r *reflectExecutor) Exec(ctx micro.Context, name string, data interface{}) (interface{}, error) {

	m, ok := r.exec[name]

	if ok {

		task := reflect.New(m.Type().In(1).Elem())

		dynamic.SetReflectValue(task, data)

		rs := m.Call([]reflect.Value{reflect.ValueOf(ctx), task})

		if len(rs) > 0 {

			if rs[1].CanInterface() && !rs[1].IsNil() {
				return nil, rs[1].Interface().(error)
			}

			if rs[0].CanInterface() {
				return rs[0].Interface(), nil
			}
		}

		return nil, errors.Errorf(404, "Not Return %s", name)

	} else {
		return nil, errors.Errorf(404, "Not Found %s", name)
	}
}

func (r *reflectExecutor) Scheme(ctx micro.Context) micro.IScheme {
	return r.scheme
}
