package models

//	type User struct {
//		Id       int    `gorm:"primaryKey"`
//		Username string `gorm:"unique;not null"`
//		Email    string `gorm:"unique"`
//		PhoneNo  string `gorm:"unique"`
//		Password string `gorm:"index;not null"`
//		HasIsActive
//		HasShopId
//	}
type User struct {
	Id       int    `gorm:"primaryKey"`
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"unique"`
	Password string `gorm:"index;not null"`
	HasIsActive
}

type UserRes struct {
	Id       int
	Username string
	Email    string
}

func (r *UserRes) CacheID() any {
	return r.Id
}
func (r *UserRes) GetUserId() int {
	return r.Id
}
