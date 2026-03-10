package task

import "time"

type TaskStruct struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	UserID    uint `gorm:"not null;index"` 
	Task      string
	IsDone    bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (TaskStruct) TableName() string {
    return "task_structs"
}