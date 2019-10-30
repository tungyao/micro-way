package util

import (
	"reflect"
)

// This package is used to store common tools.
func CheckConfig(nw interface{}, deflt interface{}) {
	switch reflect.TypeOf(nw).Kind() {
	case reflect.Struct:
		t := reflect.TypeOf(nw).Elem()
		v := reflect.ValueOf(nw).Elem()
		for i := 0; i < t.NumField(); i++ {
			n := v.Field(i)
			switch n.Kind() {
			case reflect.String:
				if n.IsZero() {
					n.SetString(reflect.ValueOf(deflt).Field(i).String())
				}
			case reflect.Int:
				if n.IsZero() {
					n.SetInt(reflect.ValueOf(deflt).Field(i).Int())
				}
			case reflect.Int64:
				if n.IsZero() {
					n.SetInt(reflect.ValueOf(deflt).Field(i).Int())
				}
			case reflect.Bool:
				if n.IsZero() {
					n.SetBool(reflect.ValueOf(deflt).Field(i).Bool())
				}
			case reflect.Float64:
				if n.IsZero() {
					n.SetFloat(reflect.ValueOf(deflt).Field(i).Float())
				}
			}
		}
	case reflect.Ptr:
		n := reflect.ValueOf(nw).Elem()
		switch n.Kind() {
		case reflect.Int:
			if n.IsZero() && n.CanSet() {
				n.SetInt(reflect.ValueOf(deflt).Int())
			}
		case reflect.String:
			if n.IsZero() && n.CanSet() {
				n.SetString(reflect.ValueOf(deflt).String())
			}
		case reflect.Struct:
			t := reflect.TypeOf(nw).Elem()
			v := reflect.ValueOf(nw).Elem()
			for i := 0; i < t.NumField(); i++ {
				n := v.Field(i)
				switch n.Kind() {
				case reflect.String:
					if n.IsZero() {
						n.SetString(reflect.ValueOf(deflt).Field(i).String())
					}
				case reflect.Int:
					if n.IsZero() {
						n.SetInt(reflect.ValueOf(deflt).Field(i).Int())
					}
				case reflect.Int64:
					if n.IsZero() {
						n.SetInt(reflect.ValueOf(deflt).Field(i).Int())
					}
				case reflect.Bool:
					if n.IsZero() {
						n.SetBool(reflect.ValueOf(deflt).Field(i).Bool())
					}
				case reflect.Float64:
					if n.IsZero() {
						n.SetFloat(reflect.ValueOf(deflt).Field(i).Float())
					}
				}
			}
		}

	}
}
func SplitString(str []byte, p []byte) [][]byte {
	group := make([][]byte, 0)
	ps := make([]int, 0)
	for i := 0; i < len(str); i++ {
		ist := make([]int, len(p))
		for k, v := range p {
			if i < len(str)-len(p) && str[i+k] == v {
				ist[k] = 1
			}
		}
		st := true
		for _, v := range ist {
			if v != 1 {
				st = false
			}
		}
		if st {
			ps = append(ps, i)
		}
	}
	ps = append(ps, len(str))
	sto := 0
	for i := 0; i < len(ps); i++ {
		group = append(group, str[sto:ps[i]])
		sto = ps[i] + 2
	}
	return group
}
