package v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/yosakoo/task-traker/internal/domain"
	"github.com/yosakoo/task-traker/internal/service"
)

func (h *Handler) initTasksRoutes(router chi.Router) {
	router.Route("/tasks", func(r chi.Router) {
		r.Use(h.AuthMiddleware)

		r.Post("/", h.createTask)
		r.Get("/{taskID}", h.getTaskByID)
		r.Get("/", h.getUserTasks)
		r.Put("/{taskID}", h.updateTask)
		r.Delete("/{taskID}", h.deleteTask)
	})
}

type taskInput struct {
	Title string  `json:"title" validate:"required"`
	Status string `json:"status"`
	Text  string  `json:"text"`
}

type getUserTasksResponse struct {
	CompletedTasks []service.TaskOut `json:"completed"`
	PendingTasks   []service.TaskOut `json:"pending"`
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	var input taskInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request body"))
		return
	}

	if err := h.validate.Struct(input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	userId := r.Context().Value("user_id").(int)
	taskID, err := h.services.Tasks.CreateTask(r.Context(), userId, service.TaskInput{
		Title: input.Title,
	})
	
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		w.Write([]byte("could not create task"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(strconv.Itoa(taskID)))
}

func (h *Handler) getTaskByID(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid task ID"))
		return
	}

	task, err := h.services.Tasks.GetTaskByID(r.Context(), taskID)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("task not found"))
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not get task"))
		return
	}

	jsonResponse, err := json.Marshal(task)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *Handler) getUserTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user_id").(int)
	completedTasks, pendingTasks, err := h.services.Tasks.GetUserTasks(r.Context(), userId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not get tasks"))
		return
	}
	
	if completedTasks == nil {
		completedTasks = []service.TaskOut{}
	}
	if pendingTasks == nil {
		pendingTasks = []service.TaskOut{}
	}

	response := getUserTasksResponse{
		CompletedTasks: completedTasks,
		PendingTasks:   pendingTasks,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not marshal response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid task ID"))
		return
	}

	var input taskInput
	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request body"))
		return
	}
	if input.Title == ""{
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid title"))
		return
	}
	err = h.services.Tasks.UpdateTask(r.Context(), taskID, service.TaskInput{
		Title: input.Title,
		Status: input.Status,
		Text:  input.Text,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not update task"))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskID"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid task ID"))
		return
	}

	err = h.services.Tasks.DeleteTask(r.Context(), taskID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not delete task"))
		return
	}

	w.WriteHeader(http.StatusOK)
}
