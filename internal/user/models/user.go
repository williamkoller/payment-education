package user_model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type User struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	Name        string
	Surname     string
	Nickname    string
	Age         int32
	Email       string `gorm:"uniqueIndex"`
	Password    string
	Roles       Roles `gorm:"type:jsonb"`
	Permissions Roles `gorm:"type:jsonb"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Roles []string

func (r Roles) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *Roles) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &r)
}
