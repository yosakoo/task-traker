package service

import (
	"context"
	"time"

	"github.com/yosakoo/task-traker/internal/repository"
	"github.com/yosakoo/task-traker/pkg/auth"
	"github.com/yosakoo/task-traker/pkg/hash"
	"github.com/yosakoo/task-traker/pkg/logger"
	"github.com/yosakoo/task-traker/pkg/rabbitmq"
)

type UserSignUpInput struct {
	Name     string
	Email    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type AuthUser struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput) (Tokens, error)
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (Tokens, error)
	GetUserByID(ctx context.Context, userID int) (AuthUser, error)
}

type TaskInput struct {
	Title  string
	Status string
	Text   string
	Time   time.Time
}

type TaskOut struct {
	ID     int       `json:"id"`
	Status string    `json:"status"`
	Title  string    `json:"title"`
	Text   string    `json:"text"`
	Time   time.Time `json:"time"`
}

type Tasks interface {
	GetTaskByID(ctx context.Context, taskID int) (TaskOut, error)
	GetUserTasks(ctx context.Context, userID int) (completedTasks []TaskOut, pendingTasks []TaskOut, err error)
	CreateTask(ctx context.Context, userID int, input TaskInput) (int, error)
	UpdateTask(ctx context.Context, taskID int, input TaskInput) error
	DeleteTask(ctx context.Context, taskID int) error
}


type Email struct {
    Subject string `json:"subject"`
    Body    string `json:"body"`
    To      string `json"to"`
}

type Emails interface{
	SendEmail(ctx context.Context, email *Email) error
}

type Services struct {
    Users  Users
    Tasks  Tasks
    Emails Emails
}

type Deps struct {
    Repos           *repo.Repositories
    QueueConn       *rabbitmq.Connection
    Log             *logger.Logger
    Hasher          hash.PasswordHasher
    TokenManager    auth.TokenManager
    AccessTokenTTL  time.Duration
    RefreshTokenTTL time.Duration
    EmailService    Emails 
}



func NewServices(deps Deps) *Services {
	
    emailService := NewEmailService(deps.QueueConn)
    userService :=  NewUserService(deps.Repos.Users, deps.Log, deps.Hasher, deps.TokenManager, emailService, deps.AccessTokenTTL, deps.RefreshTokenTTL)
    taskService :=  NewTaskService(deps.Repos.Tasks)
    return &Services{Users: userService, Tasks: taskService, Emails: emailService}
}

