package model

import (
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Attachment struct {
	Model
	FileName     string `json:"file_name" gorm:"column:file_name;type:varchar(255) not null;default:''"`
	FileLocation string `json:"file_location" gorm:"column:file_location;type:varchar(255) not null;default:''"`
	FileSize     int64  `json:"file_size" gorm:"column:file_size;type:bigint(20) unsigned not null;default:0"`
	FileMd5      string `json:"file_md5" gorm:"column:file_md5;type:varchar(32) unique not null;default:''"`
	Width        uint   `json:"width" gorm:"column:width;type:int(10) unsigned not null;default:0"`
	Height       uint   `json:"height" gorm:"column:height;type:int(10) unsigned not null;default:0"`
	Status       uint   `json:"status" gorm:"column:status;type:tinyint(1) unsigned not null;default:0;index:idx_status"`
	Logo         string `json:"logo" gorm:"-"`
	Thumb        string `json:"thumb" gorm:"-"`
}

func (attachment *Attachment) Save(db *gorm.DB) error {
	if attachment.Id == 0 {
		attachment.CreatedTime = time.Now().Unix()
	}
	attachment.UpdatedTime = time.Now().Unix()
	if err := db.Save(attachment).Error; err != nil {
		return err
	}
	attachment.GetThumb()
	return nil
}

func (attachment *Attachment) GetThumb() {
	if strings.HasPrefix(attachment.FileLocation, "http") {
		attachment.Logo = attachment.FileLocation
		attachment.Thumb = attachment.FileLocation
	} else {
		pfx := "/static/uploads/"
		attachment.Logo = pfx + attachment.FileLocation
		paths, fileName := filepath.Split(attachment.FileLocation)
		attachment.Thumb = pfx + paths + "thumb_" + fileName
	}
}
