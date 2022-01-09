package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cast"
	"log"
	"reflect"
	"strconv"
	"strings"
)

//GetFieldName 获取结构体中字段的名称
func GetFieldName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		log.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		result = append(result, t.Field(i).Name)
	}
	return result
}

//GetTagName 获取结构体中Tag的值，如果没有tag则返回字段值
func GetTagName(structName interface{}) []string {
	t := reflect.TypeOf(structName)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		logrus.Println("Check type error not Struct")
		return nil
	}
	fieldNum := t.NumField()
	result := make([]string, 0, fieldNum)
	for i := 0; i < fieldNum; i++ {
		tagName := t.Field(i).Name
		tags := strings.Split(string(t.Field(i).Tag), "\"")
		if len(tags) > 1 {
			tagName = tags[1]
		}
		result = append(result, tagName)
	}
	return result
}

//StructToSlice 获取结构体中的值列表
func StructToSlice(structName interface{}) (result []string) {
	v := reflect.ValueOf(structName)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)
		switch f.Kind() {
		case reflect.String:
			result = append(result, f.String())
		case reflect.Int32:
			result = append(result, strconv.FormatInt(int64(f.Int()), 10))
		case reflect.Map:
			var str []string
			iter := f.MapRange()
			for iter.Next() {
				src := iter.Key().String() + "=" + iter.Value().String()
				str = append(str, src)
			}
			keys := f.MapKeys()
			for i = 0; i < f.Len(); i++ {
				str = append(str, f.MapIndex(keys[i]).String())
			}
			result = append(result, strings.Join(str, " "))
		case reflect.Slice:
			var str []string
			for i := 0; i < f.Len(); i++ {
				str = append(str, f.Index(i).String())
			}
			result = append(result, strings.Join(str, " "))
		default:
			result = append(result, f.String())
		}
	}
	return
}

// Map2StrSliceE 字典转切片
func Map2StrSliceE(i interface{}) []string {
	kind := reflect.TypeOf(i).Kind()
	if kind != reflect.Map {
		logrus.Error("the input is not a map")
		return nil
	}
	m := reflect.ValueOf(i)
	keys := m.MapKeys()
	res := make([]string, 0, len(keys))
	for _, k := range keys {
		// convert the key to string
		sK, err := cast.ToStringE(k.Interface())
		if err != nil {
			return nil
		}
		// convert the value to string
		v := m.MapIndex(k)
		sV, err := cast.ToStringE(v.Interface())
		if err != nil {
			return nil
		}
		res = append(res, fmt.Sprintf("%s=%s", sK, sV))
	}
	return res
}
