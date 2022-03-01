package model

type Model struct {
	Id          uint  `json:"id" gorm:"column:id;type:int(10) unsigned not null AUTO_INCREMENT;primaryKey"`
	CreatedTime int64 `json:"created_time" gorm:"column:created_time;type:int(11) not null;default:0;index:idx_created_time"`
	UpdatedTime int64 `json:"updated_time" gorm:"column:updated_time;type:int(11) not null;default:0;index:idx_updated_time"`
	DeletedTime int64 `json:"-" gorm:"column:deleted_time;type:int(11) not null;default:0;"`
}
