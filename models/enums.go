package models

import "time"

type ImageReferenceType string

const (
	ImageReferenceTypeGhost ImageReferenceType = "ghost"
)

type MyDateTime struct {
	time.Time
}

func (t MyDateTime) String() string {
	return t.Format("Jan 2 3:04 PM")
}

type MyDate struct {
	time.Time
}

func (t MyDate) String() string {
	return t.Format("Jan 2")
}

func (t MyDate) InputValue() string {
	return t.Format(time.DateOnly)
}

// func (t MyDate) IsZero() bool {
// 	return t.IsZero()
// }
