package diff

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	// ErrTypeMismatch : Compared types do not match
	ErrTypeMismatch = errors.New("types do not match")
)

// Changelog : stores a list of changed items
type Changelog []Change

// Change : stores information about a changed item
type Change struct {
	Type string      `json:"type"`
	Path []string    `json:"path"`
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}

// Changed ...
func Changed(a, b interface{}) bool {
	cl, _ := Diff(a, b)
	return len(cl) > 0
}

// Diff ...
func Diff(a, b interface{}) (Changelog, error) {
	var cl Changelog

	return cl, cl.diff([]string{}, reflect.ValueOf(a), reflect.ValueOf(b))
}

func (cl *Changelog) diff(path []string, a, b reflect.Value) error {
	var err error

	if a.Kind() != b.Kind() {
		return errors.New("types do not match")
	}

	switch a.Kind() {
	case reflect.Struct:
		err = cl.diffStruct(path, a, b)
	case reflect.Array:
		err = cl.diffArray(path, a, b)
	case reflect.Slice:
		err = cl.diffSlice(path, a, b)
	case reflect.String:
		err = cl.diffString(path, a, b)
	case reflect.Bool:
		err = cl.diffBool(path, a, b)
	case reflect.Int:
		err = cl.diffInt(path, a, b)
	default:
		err = errors.New("unsupported type: " + a.Kind().String())
	}

	return err
}

func (cl *Changelog) diffStruct(path []string, a, b reflect.Value) error {
	if a.Kind() != b.Kind() {
		return ErrTypeMismatch
	}

	fmt.Println(a.NumField())

	for i := 0; i < a.NumField(); i++ {
		name := a.Type().Field(i).Name

		af := a.Field(i)
		bf := b.FieldByName(name)

		fpath := append(path, tagName(a, i))

		err := cl.diff(fpath, af, bf)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cl *Changelog) diffArray(path []string, a, b reflect.Value) error {
	if a.Kind() != b.Kind() {
		return ErrTypeMismatch
	}

	fmt.Println("ARRAY")

	return nil
}

func (cl *Changelog) diffSlice(path []string, a, b reflect.Value) error {
	if a.Kind() != b.Kind() {
		return ErrTypeMismatch
	}

	for i := 0; i < a.Len(); i++ {
		ae := a.Index(i)

		switch ae.Kind() {
		case reflect.Struct:
			id := identifier(ae)
			if id != nil {
				x := c[id]
				x.A = &ae
			}
		default:
			fmt.Println(ae.Interface())
			//cl.diff(path, a, b)
		}
	}

	for i := 0; i < b.Len(); i++ {
		be := b.Index(i)

		switch be.Kind() {
		case reflect.Struct:
			id := identifier(be)
			if id != nil {
				x := c[id]
				x.B = &be
			}
		default:
			fmt.Println(be.Interface())
			//cl.diff(path, a, b)
		}
	}

	for k, v := range c {

	}

	return nil
}

func (cl *Changelog) diffString(path []string, a, b reflect.Value) error {
	if a.Kind() != b.Kind() {
		return ErrTypeMismatch
	}

	if a.String() != b.String() {
		(*cl) = append((*cl), Change{
			Type: "update",
			Path: path,
			From: a.Interface(),
			To:   b.Interface(),
		})
	}

	return nil
}

func (cl *Changelog) diffBool(path []string, a, b reflect.Value) error {
	if a.Kind() != b.Kind() {
		return ErrTypeMismatch
	}

	return nil
}

func (cl *Changelog) diffInt(path []string, a, b reflect.Value) error {
	if a.Kind() != b.Kind() {
		return ErrTypeMismatch
	}

	return nil
}

func tag(v reflect.Value, i int) string {
	return v.Type().Field(i).Tag.Get("diff")
}

func tagName(v reflect.Value, i int) string {
	t := tag(v, i)

	parts := strings.Split(t, ",")
	if len(parts) < 1 {
		return ""
	}

	return parts[0]
}

func identifier(v reflect.Value) interface{} {
	for i := 0; i < v.NumField(); i++ {
		t := tag(v, i)

		parts := strings.Split(t, ",")
		if len(parts) < 1 {
			continue
		}

		if parts[1] == "identifier" {
			return v.Field(i).Interface()
		}
	}

	return nil
}