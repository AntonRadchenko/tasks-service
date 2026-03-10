package task

import (
	"strings"
	"time"

	"github.com/AntonRadchenko/tasks-service/internal/database"
	"gorm.io/gorm"
)

type TaskRepoInterface interface {
	Create(task *TaskStruct) (*TaskStruct, error)
	GetAll() ([]TaskStruct, error)
	GetByID(id uint) (TaskStruct, error)
	GetByUserID(userID uint) ([]TaskStruct, error) // для ListTasksByUser
	Update(task *TaskStruct) (*TaskStruct, error)
	Delete(task *TaskStruct) error
}

type TaskRepo struct{}

func (r *TaskRepo) Create(task *TaskStruct) (*TaskStruct, error) {
	err := database.DB.Create(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *TaskRepo) GetAll() ([]TaskStruct, error) {
	var tasks []TaskStruct

	err := database.DB.Find(&tasks).Error
	if err != nil {
		if strings.Contains(err.Error(), "relation") {
			return []TaskStruct{}, nil
		}
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepo) GetByID(id uint) (TaskStruct, error) {
	var task TaskStruct
	err := database.DB.First(&task, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return TaskStruct{}, err
		}
		return task, err
	}
	return task, nil
}

func (r *TaskRepo) GetByUserID(userID uint) ([]TaskStruct, error) {
	var tasks []TaskStruct
	err := database.DB.Where("user_id = ?", userID).Find(&tasks).Error
	if err != nil {
		if strings.Contains(err.Error(), "relation") {
			return []TaskStruct{}, nil
		}
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepo) Update(task *TaskStruct) (*TaskStruct, error) {
	task.UpdatedAt = time.Now()
	err := database.DB.Save(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (r *TaskRepo) Delete(task *TaskStruct) error {
	now := time.Now()
	task.DeletedAt = &now
	err := database.DB.Delete(task).Error
	if err != nil {
		return err
	}
	return nil
}