package model

import "github.com/jinzhu/gorm"

type Task struct {
	gorm.Model
	Name 			string
	Description 	string `gorm:"type:text"`
	WorkDurations 	[]WorkDuration `gorm:"type:json"`

}

type WorkDuration struct {
	StartTimestamp 		int64
	EndTimestamp 		int64
}