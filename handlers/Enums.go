package handlers

import (
	"encoding/json"
	"strings"
)

// auto trimming space
type inputString string

func (i *inputString) UnmarshalJSON(data []byte) error {
	var temp string
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	temp = strings.TrimSpace(temp)

	*i = inputString(temp)
	return nil
}

func (i inputString) String() string {
	return string(i)
}

type optionalString string

func (i *optionalString) UnmarshalJSON(data []byte) error {
	var temp string
	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}
	temp = strings.TrimSpace(temp)

	*i = optionalString(temp)
	return nil
}

// pointer string
func (i *optionalString) StringPtr() *string {
	if i == nil {
		return nil
	} else {
		if *i == "" {
			return nil
		}
	}
	s := string(*i)
	return &s
}

// returns empty string if nil, safe dereference
func (i *optionalString) String() string {
	if i == nil {
		return ""
	}
	return string(*i)
}

// if field is omitted
func (i *optionalString) IsNil() bool {
	return i == nil
}

func (i *optionalString) IsPresent() bool {
	return i != nil
}

// if field is empty
func (i *optionalString) IsEmpty() bool {
	if i == nil {
		return true
	}
	if *i == "" {
		return true
	}

	return false
}
