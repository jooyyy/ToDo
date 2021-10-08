package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

const (
	TaskStatusPause = "pause"
	TaskStatusDoing = "doing"
	TaskStatusDone = "done"
)

type Task struct {
	gorm.Model
	ProjectID 		uint
	Project 		Project
	Date 			string
	Name 			string
	Description 	string `gorm:"type:text"`
	Durations 		Durations `gorm:"type:text"`
	Status  		string
}

type Durations []Duration

type Duration struct {
	StartTime 		string
	EndTime 		string
}

func (j *Durations) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal Durations value:", value))
	}

	err := json.Unmarshal(bytes, &j)
	return err
}

// Value return json value, implement driver.Valuer interface
func (j Durations) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(j)
	if err != nil {
		return "[]", nil
	}
	return string(b), nil
}