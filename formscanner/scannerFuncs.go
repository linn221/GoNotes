package formscanner

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Int = func(r *http.Request, name string) (result int, omitted bool, err error) {
	v := r.PostFormValue(name)
	if v == "" {
		omitted = true
		return
	}

	i, err := strconv.Atoi(v)
	if err != nil {
		return
	}
	result = i
	return
}

var IntRequired = func(r *http.Request, name string) (result int, omitted bool, err error) {
	v := r.PostFormValue(name)
	if v == "" {
		err = fmt.Errorf("%s is required", name)
		return
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return
	}
	result = i
	return
}

var Date = func(r *http.Request, name string) (result time.Time, omitted bool, err error) {
	v := r.PostFormValue(name)
	if v == "" {
		omitted = true
		return
	}
	t, err := time.Parse(time.DateOnly, v)
	if err != nil {
		return
	}
	result = t
	return
}

var DateRequired = func(r *http.Request, name string) (result time.Time, omitted bool, err error) {
	v := r.PostFormValue(name)
	if v == "" {
		err = fmt.Errorf("%s is required", name)
		return
	}
	t, err := time.Parse(time.DateOnly, v)
	if err != nil {
		return
	}
	result = t
	return
}

var UIntRequired = func(r *http.Request, name string) (result uint, omitted bool, err error) {
	v := r.PostFormValue(name)
	if v == "" {
		err = fmt.Errorf("%s is required", name)
		return
	}
	i, err := strconv.ParseUint(v, 10, 10)
	if err != nil {
		return
	}
	result = uint(i)
	return
}

// required
var String = func(r *http.Request, name string) (result string, omitted bool, err error) {
	v := r.PostFormValue(name)
	if v == "" {
		omitted = true
		return
	}
	result = strings.TrimSpace(v)
	return
}

var StringRequired = func(r *http.Request, name string) (result string, omitted bool, err error) {
	v := r.PostFormValue(name)
	if v == "" {
		err = fmt.Errorf("%s is required", name)
		return
	}
	result = strings.TrimSpace(v)
	return
}
