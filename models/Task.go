package models

import (
	"time"
)

type TaskStatus string

const (
	TaskStatusSuccess    TaskStatus = "success"
	TaskStatusInProgress TaskStatus = "inprogress"
	TaskStatusCancel     TaskStatus = "cancel"
	TaskStatusFail       TaskStatus = "fail"
	TaskStatusPending    TaskStatus = "pending"
)

type Task struct {
	Id          int        `gorm:"primaryKey"`
	Status      TaskStatus `gorm:"enum('success','cancel','fail','pending','inprogress');default:pending"`
	Title       string     `gorm:"index;not null"`
	Description string
	LabelId     uint       `gorm:"index;not null"`
	CreateDate  time.Time  `gorm:"index;not null"`
	Deadline    *time.Time `gorm:"index;default:null"`
	HasUserId
}

// type TaskResource struct {
// 	Id          int
// 	Status      TaskStatus
// 	Title       string
// 	Description string
// 	LabelId     int
// 	CreateDate  MyDate
// 	Deadline    MyDate
// }

// func (t Task) Resource() any {
// 	return TaskResource{
// 		Id:          t.Id,
// 		Status:      t.Status,
// 		Title:       t.Title,
// 		Description: t.Description,
// 		LabelId:     int(t.LabelId),
// 		CreateDate:  MyDate{&t.CreateDate},
// 		Deadline:    MyDate{t.Deadline},
// 	}
// }

// func (t *Task) Validate(db *gorm.DB, id int, userId int) services.FormErrors {
// 	m := validate.ValidateInBatch(db,
// 		validate.NewExistsRule("labels", t.LabelId, "label not found", validate.NewFilter("user_id = ?", userId), "label_id"))
// 	if t.Deadline != nil && t.Deadline.Before(t.CreateDate) {
// 		m["deadline"] = errors.New("deadline cannot be before create date")
// 	}
// 	return m
// }
