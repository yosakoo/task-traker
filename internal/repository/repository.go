package repo

import (
	"context"
	"time"

	"github.com/yosakoo/task-traker/internal/domain/models"
	"github.com/yosakoo/task-traker/pkg/postgres"
)

type Users interface {
	GetUserByCredentials(ctx context.Context, email string, password []byte) (*models.User, error)
	AddUser(ctx context.Context, user models.User) (int, error)
	SetSession(ctx context.Context, userId int, refresh string, expiresAt time.Time) error
	GetUserByRefresh(ctx context.Context, refresh string) (int, error)
	GetUserByID(ctx context.Context, userID int) (*models.User, error) 
}

type Tasks interface {
	GetTaskByID(ctx context.Context, taskID int) (*models.Task, error)
	CreateTask(ctx context.Context, userID int, task models.Task) (int, error)
	UpdateTask(ctx context.Context, taskID int, task models.Task) error
	DeleteTask(ctx context.Context, taskID int) error
	GetUserTasks(ctx context.Context, userID int) ([]models.Task, error)
}

type Repositories struct{
	Users Users
	Tasks Tasks
}

func NewRepositories(pool *postgres.Storage) *Repositories{
	return &Repositories{
		Users: NewUserRepo(pool),
		Tasks: NewTaskRepo(pool),
	}
}
