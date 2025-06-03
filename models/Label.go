package models

type Label struct {
	Id          int    `gorm:"primaryKey"`
	Name        string `gorm:"index;not null"`
	Description string
	HasIsActive
	HasUserId
}

// func (l Label) Validate(db *gorm.DB, id int, userId int) map[string]error {
// 	return validate.ValidateInBatch(db,
// 		validate.NewUniqueRule("labels", "name", l.Name, id, "duplicate label name").Filter("user_id = ?", userId),
// 	)
// }
