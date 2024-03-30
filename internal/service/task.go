package service

import (
	"context"
	"time"

	"github.com/yosakoo/task-traker/internal/domain"
	"github.com/yosakoo/task-traker/internal/domain/models"
	"github.com/yosakoo/task-traker/internal/repository"
)

type TaskService struct {
	repo repo.Tasks
}

func NewTaskService(repo repo.Tasks) *TaskService {
	return &TaskService{
		repo: repo,
	}
}

func (s *TaskService) GetTaskByID(ctx context.Context, taskID int) (TaskOut, error) {
	task, err := s.repo.GetTaskByID(ctx, taskID)
	if err != nil {
		return TaskOut{}, err
	}
	taskOut := TaskOut{
		ID:    task.ID,
		Status: task.Status,
		Title: task.Title,
		Text:  *task.Text,
		Time:  *task.Time,
	}
	return taskOut, nil
}


func (s *TaskService) GetUserTasks(ctx context.Context, userID int) (completedTasks []TaskOut, pendingTasks []TaskOut, err error) {
    tasks, err := s.repo.GetUserTasks(ctx, userID)
    if err != nil {
        return nil, nil, err
    }
    if len(tasks) == 0 {
        return nil, nil, domain.ErrTaskNotFound
    }

    for _, task := range tasks {
        taskOut := TaskOut{
            ID:     task.ID,
			Title: task.Title,
            Status: task.Status,
        }
        if task.Text != nil {
            taskOut.Text = *task.Text
        }

        if task.Time != nil {
            taskOut.Time = *task.Time
        }

        if task.Status == "completed" {
            completedTasks = append(completedTasks, taskOut)
        } else if task.Status == "pending" {
            pendingTasks = append(pendingTasks, taskOut)
        }
    }

    return completedTasks, pendingTasks, nil
}

func (s *TaskService) CreateTask(ctx context.Context, userID int, input TaskInput) (int, error) {
	task := models.Task{
		UserID: userID,
		Title:  input.Title,
		Status: "pending",
	}
	taskID, err := s.repo.CreateTask(ctx, userID, task)
	if err != nil {
		return 0, err
	}
	return taskID, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, taskID int, input TaskInput) error {
    var currentTime *time.Time

    if input.Status == "pending" {
        currentTime = nil
    } else {
        now := time.Now()
        currentTime = &now
    }

    task := models.Task{
        Title:  input.Title,
        Status: input.Status,
        Text:   &input.Text,
        Time:   currentTime, 
    }

    err := s.repo.UpdateTask(ctx, taskID, task)
    if err != nil {
        return err
    }
    return nil
}





func (s *TaskService) DeleteTask(ctx context.Context, taskID int) error {
	err := s.repo.DeleteTask(ctx, taskID)
	if err != nil {
		return err
	}
	return nil
}
