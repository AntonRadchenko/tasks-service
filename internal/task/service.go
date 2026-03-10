package task

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type CreateTaskParams struct {
	Task   string
	IsDone *bool
	UserID uint
}

type UpdateTaskParams struct {
	Task   *string
	IsDone *bool
	UserID *uint
}

type Task struct {
	ID     uint
	Task   string
	IsDone *bool
	UserID uint
}

type TaskService struct {
	repo TaskRepoInterface
}

func NewTaskService(r TaskRepoInterface) *TaskService {
	return &TaskService{repo: r}
}

func (s *TaskService) CreateTask(params CreateTaskParams) (*Task, error) {
	if strings.TrimSpace(params.Task) == "" {
		return nil, errors.New("task is empty")
	}

	if params.UserID == 0 {
		return nil, errors.New("user_id is required")
	}

	isDone := false
	if params.IsDone != nil {
		isDone = *params.IsDone
	}

	dbTask := &TaskStruct{
		Task:   params.Task,
		IsDone: isDone,
		UserID: params.UserID,
	}

	createdTask, err := s.repo.Create(dbTask)
	if err != nil {
		return nil, err
	}

	return &Task{
		ID:     createdTask.ID,
		Task:   createdTask.Task,
		IsDone: &createdTask.IsDone,
		UserID: createdTask.UserID,
	}, nil
}

func (s *TaskService) GetTasks() ([]Task, error) {
	dbTasks, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0, len(dbTasks))
	for _, dbTask := range dbTasks {
		tasks = append(tasks, Task{
			ID:     dbTask.ID,
			Task:   dbTask.Task,
			IsDone: &dbTask.IsDone,
			UserID: dbTask.UserID,
		})
	}
	return tasks, nil
}

func (s *TaskService) GetTaskByID(id uint) (*Task, error) {
	dbTask, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return &Task{
		ID:     dbTask.ID,
		Task:   dbTask.Task,
		IsDone: &dbTask.IsDone,
		UserID: dbTask.UserID,
	}, nil
}

func (s *TaskService) GetTasksByUser(userID uint) ([]Task, error) {
	if userID == 0 {
		return nil, errors.New("user_id is required")
	}

	dbTasks, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0, len(dbTasks))
	for _, dbTask := range dbTasks {
		tasks = append(tasks, Task{
			ID:     dbTask.ID,
			Task:   dbTask.Task,
			IsDone: &dbTask.IsDone,
			UserID: dbTask.UserID,
		})
	}
	return tasks, nil
}

func (s *TaskService) UpdateTask(id uint, params UpdateTaskParams) (*Task, error) {
	dbTask, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	updated := false

	if params.Task != nil {
		if strings.TrimSpace(*params.Task) == "" {
			return nil, errors.New("task is empty")
		}
		dbTask.Task = *params.Task
		updated = true
	}

	if params.IsDone != nil {
		dbTask.IsDone = *params.IsDone
		updated = true
	}

	if params.UserID != nil {
		if *params.UserID == 0 {
			return nil, errors.New("user_id cannot be 0")
		}
		
		dbTask.UserID = *params.UserID
		updated = true
	}

	if !updated {
		return nil, errors.New("no fields to update")
	}

	updatedTask, err := s.repo.Update(&dbTask)
	if err != nil {
		return nil, err
	}

	return &Task{
		ID:     updatedTask.ID,
		Task:   updatedTask.Task,
		IsDone: &updatedTask.IsDone,
		UserID: updatedTask.UserID,
	}, nil
}

func (s *TaskService) DeleteTask(id uint) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("task not found")
		}
		return err
	}

	err = s.repo.Delete(&task)
	if err != nil {
		return err
	}
	return nil
}