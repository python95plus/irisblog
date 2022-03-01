package model

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Admin struct {
	Model
	UserName string `json:"user_name" gorm:"column:user_name;type:varchar(16) not null; default:'';index:idx_user_name"`
	Password string `jsong:"-" gorm:"column:password;type:varchar(128) not null; default:''"`
	Status   uint   `json:"status" gorm:"column:status;type:tinyint(1) unsigned not null;default:0;index:idx_status"`
	Group    string `json:"group" gorm:"column:group;type:varchar(64) not null;default:''"`
}

func (admin *Admin) Save(db *gorm.DB) error {
	if admin.Id == 0 {
		admin.CreatedTime = time.Now().Unix()
	}
	admin.UpdatedTime = time.Now().Unix()
	if err := db.Save(admin).Error; err != nil {
		return err
	}
	return nil
}

func (admin *Admin) EncryptPassword(password string) string {
	if password == "" {
		return ""
	}
	pass := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		return ""
	}
	return string(hash)
}

func (admin *Admin) CheckPassword(password string) bool {
	if password == "" {
		return false
	}
	byteHas := []byte(admin.Password)
	bytePass := []byte(password)

	err := bcrypt.CompareHashAndPassword(byteHas, bytePass)
	if err != nil {
		return false
	} else {
		return true
	}
}
