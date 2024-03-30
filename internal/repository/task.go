package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/yosakoo/task-traker/internal/domain"
	"github.com/yosakoo/task-traker/internal/domain/models"
	"github.com/yosakoo/task-traker/pkg/postgres"
)

type TaskRepo struct {
	s *postgres.Storage
}

func NewTaskRepo(pg *postgres.Storage) *TaskRepo {
	return &TaskRepo{s: pg}
}

func (r *TaskRepo) GetTaskByID(ctx context.Context, taskID int) (*models.Task, error) {
    var task models.Task
    query := "SELECT id, status, title, text, time FROM tasks WHERE id = $1"
    err := r.s.Pool.QueryRow(ctx, query, taskID).Scan(&task.ID, &task.Status, &task.Title, &task.Text, &task.Time)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domain.ErrTaskNotFound
        }
        return nil, err
    }

    return &task, nil
}

func (r *TaskRepo) GetUserTasks(ctx context.Context, userID int) ([]models.Task, error) {
    var tasks []models.Task
    query := "SELECT id, status, title, text, time FROM tasks WHERE user_id = $1"
    rows, err := r.s.Pool.Query(ctx, query, userID)
    if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
            return nil, domain.ErrTaskNotFound
        }
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var task models.Task
        err := rows.Scan(&task.ID, &task.Status, &task.Title, &task.Text, &task.Time)
        if err != nil {
            return nil, err
        }
        tasks = append(tasks, task)
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return tasks, nil
}

func (r *TaskRepo) CreateTask(ctx context.Context, userID int, task models.Task) (int, error) {
	txOptions := pgx.TxOptions{}

	tx, err := r.s.Pool.BeginTx(ctx, txOptions)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	var taskID int
	err = tx.QueryRow(ctx, "INSERT INTO tasks (user_id, title) VALUES ($1, $2) RETURNING id", userID, task.Title).Scan(&taskID)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, errors.New("error committing database transaction")
	}

	return taskID, nil
}

func (r *TaskRepo) UpdateTask(ctx context.Context, taskID int, task models.Task) error {
	txOptions := pgx.TxOptions{}

	tx, err := r.s.Pool.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	fmt.Println(task.Status)
	_, err = tx.Exec(ctx, "UPDATE tasks SET title = $1,status =$2, text = $3, time = $4 WHERE id = $5", task.Title, task.Status, task.Text, task.Time, taskID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.New("error committing database transaction")
	}

	return nil
}

func (r *TaskRepo) DeleteTask(ctx context.Context, taskID int) error {
	txOptions := pgx.TxOptions{}

	tx, err := r.s.Pool.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "DELETE FROM tasks WHERE id = $1", taskID)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return errors.New("error committing database transaction")
	}

	return nil
}
