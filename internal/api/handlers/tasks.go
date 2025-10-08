package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"queueit/internal/db"
	"queueit/internal/helper"
	"queueit/internal/models"
	"queueit/pkg/logger"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// GetAllTasks godoc
// @Summary      Get all tasks
// @Description  Fetch all tasks, optionally filtering by status and/or priority
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        status   query     string  false  "Comma-separated task statuses to filter (1=pending, 2=wip, 3=done, 4=archived)"
// @Param        priority query     string  false  "Comma-separated priority values to filter (1=high,2=medium,3=low,0=default(medium))"
// @Success      200  {array}   models.GetTasksResponse
// @Failure      400  {string}  string "Bad request"
// @Failure      500  {string}  string "Fetching tasks failed"
// @Router       /v1/tasks [get]
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	helper.SetJSONHeader(w)

	status := r.URL.Query().Get("status")
	priority := r.URL.Query().Get("priority")

	query := `
		SELECT
		task_id, title, description, priority, status, created_at, deadline_at 
		FROM tasksmaster WHERE 1=1
	`
	if status != "" {
		query = fmt.Sprintf("%s AND status IN (%s)", query, strings.Join(strings.Split(status, ","), ","))
	}

	if priority != "" {
		query = fmt.Sprintf("%s AND priority IN (%s)", query, strings.Join(strings.Split(priority, ","), ","))
	}

	rows, err := db.GetDBInfo().Q(query)
	if err != nil {
		logger.Error(err, "GetAllTasks ~ db query failed")
		http.Error(w, "fetching tasks failed", http.StatusInternalServerError)
		return
	}

	var deadline sql.NullString
	defer rows.Close()
	var tasks []models.GetTasksResponse
	for rows.Next() {
		var t models.GetTasksResponse
		if err := rows.Scan(
			&t.TaskID,
			&t.Title,
			&t.Description,
			&t.Priority,
			&t.Status,
			&t.CreatedAt,
			&deadline,
		); err != nil {
			logger.Error(err, "GetAllTasks ~ row scan failed")
			http.Error(w, "fetching tasks failed", http.StatusInternalServerError)
			return
		}

		// validate deadline (else NIL)
		if deadline.Valid {
			ct, err := time.Parse(time.RFC3339, deadline.String)
			if err != nil {
				logger.Error(err, "GetAllTasks ~ deadline validation failed")
				http.Error(w, "fetching tasks failed", http.StatusInternalServerError)
				return
			}
			t.DeadlineAt = &ct
		} else {
			t.DeadlineAt = nil
		}

		tasks = append(tasks, t)
	}

	buffer := new(bytes.Buffer)
	if err = json.NewEncoder(buffer).Encode(tasks); err != nil {
		logger.Error(err, "GetAllTasks ~ JSON encoding failed")
		http.Error(w, "fetching tasks failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buffer.Bytes())
}

// CreateTask godoc
// @Summary      Create a new task
// @Description  Create a new task with title, description, priority, and optional deadline
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        task  body      models.CreateTaskRequest  true  "Task to create"
// @Success      200  {object}  models.GenricTaskResponse  "Task created successfully"
// @Failure      400  {string}  string "Invalid JSON"
// @Failure      422  {string}  string "Unprocessable entity (blank title)"
// @Failure      500  {string}  string "Creating task failed"
// @Router       /v1/tasks [post]
func CreateTask(w http.ResponseWriter, r *http.Request) {
	helper.SetJSONHeader(w)

	var ctr models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&ctr); err != nil {
		logger.Error(err, "CreateTask ~ JSON decoding failed")
		http.Error(w, "creating task failed", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if ctr.Title == "" {
		logger.Error("CreateTask ~ blank title")
		http.Error(w, "creating task failed", http.StatusUnprocessableEntity)
		return
	}
	if !(helper.IsValidPriority(ctr.Priority)) {
		ctr.Priority = 2 // medium
	}

	query := `
		INSERT into tasksmaster
		(
			title,
			description,
			priority,
			status,
			deadline_at
		)
		VALUES
		(
			?,?,?,?,?
		);
	`

	exec_result, err := db.GetDBInfo().E(
		query,
		ctr.Title,
		ctr.Description,
		ctr.Priority,
		models.STATUS_PENDING, // default status
		ctr.DeadlineAt,
	)
	if err != nil {
		logger.Error(err, "CreateTask ~ execution failed")
		http.Error(w, "creating task failed", http.StatusInternalServerError)
		return
	}
	taskID, _ := exec_result.LastInsertId()
	resp := models.GenricTaskResponse{
		TaskID:  taskID,
		Message: "Task created",
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error(err, "CreateTask ~ json encoding failed")
		http.Error(w, "creating task failed", http.StatusInternalServerError)
		return
	}
}

// GetTaskByID godoc
// @Summary      Get task details by ID
// @Description  Fetch a single task record from the database using its unique ID.
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        id   path      int     true  "Task ID"
// @Success      200  {object}  models.GetTasksResponse  "Task details fetched successfully"
// @Failure      400  {string}  string  "Invalid or missing task ID"
// @Failure      404  {string}  string  "Task not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /v1/tasks/{id} [get]
func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	helper.SetJSONHeader(w)

	idstr, exists := mux.Vars(r)["id"]
	if !exists {
		logger.Error("GetTaskByID ~ id missing")
		http.Error(w, "id missing", http.StatusBadRequest)
		return
	}

	var deadline sql.NullString
	var t models.GetTasksResponse
	query := fmt.Sprintf(`
		SELECT
		task_id, title, description, priority, status, created_at, deadline_at 
		FROM tasksmaster WHERE task_id = %s
	`, idstr)

	rows, err := db.GetDBInfo().Q(query)
	if err != nil {
		logger.Error(err, "GetTaskByID ~ query execution failed")
		http.Error(w, "invalid task ID or internal error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		logger.Error("GetTaskByID ~ rows.Next-false, no rows present")
		http.Error(w, "task id not found", http.StatusBadRequest)
		return
	}

	if err := rows.Scan(
		&t.TaskID,
		&t.Title,
		&t.Description,
		&t.Priority,
		&t.Status,
		&t.CreatedAt,
		&deadline,
	); err != nil {
		logger.Error(err, "GetTaskByID ~ row scan failed")
		http.Error(w, "fetching task failed", http.StatusInternalServerError)
		return
	}

	// validate deadline (else NIL)
	if deadline.Valid {
		ct, err := time.Parse(time.RFC3339, deadline.String)
		if err != nil {
			logger.Error(err, "GetTaskByID ~ deadline validation failed")
			http.Error(w, "fetching task failed", http.StatusInternalServerError)
			return
		}
		t.DeadlineAt = &ct
	} else {
		t.DeadlineAt = nil
	}

	buffer := new(bytes.Buffer)
	if err = json.NewEncoder(buffer).Encode(t); err != nil {
		logger.Error(err, "GetTaskByID ~ JSON encoding failed")
		http.Error(w, "fetching tasks failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(buffer.Bytes())
}

// UpdateTask godoc
// @Summary      Update task fields by ID
// @Description  Partially update one or more fields of a task (title, description, status, priority, deadline).
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        id   path      int                     true  "Task ID"
// @Param        task body      models.UpdateTaskRequest  true  "Fields to update"
// @Success      200  {object}  models.GenricTaskResponse  "Task updated successfully"
// @Failure      400  {string}  string  "Invalid input or missing ID"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /v1/tasks/{id} [patch]
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	helper.SetJSONHeader(w)

	idstr, exists := mux.Vars(r)["id"]
	if !exists {
		logger.Error("UpdateTask ~ id missing")
		http.Error(w, "id missing", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	var t models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		logger.Error(err, "UpdateTask ~ request json decoding failed")
		http.Error(w, "invalid JOSN payload in request", http.StatusBadRequest)
		return
	}
	// UPDATE tasksmaster SET title = ?, status = ? WHERE task_id = ?
	var fields []string
	var args []interface{}

	if t.Title != nil {
		if *t.Title == "" {
			http.Error(w, "invalid/empty title", http.StatusBadRequest)
			return
		}
		fields = append(fields, "title=?")
		args = append(args, *t.Title)
	}

	if t.Description != nil {
		fields = append(fields, "description = ?")
		args = append(args, *t.Description)
	}

	if t.Status != nil {
		if !(helper.IsValidStatus(*t.Status)) {
			http.Error(w, "invalid status", http.StatusBadRequest)
			return
		}
		fields = append(fields, "status = ?")
		args = append(args, *t.Status)
	}

	if t.Priority != nil {
		if !(helper.IsValidPriority(*t.Priority)) {
			http.Error(w, "invalid priority", http.StatusBadRequest)
			return
		}
		fields = append(fields, "priority = ?")
		args = append(args, *t.Priority)
	}

	if t.DeadlineAt != nil {
		fields = append(fields, "deadline_at = ?")
		args = append(args, t.DeadlineAt.Format(time.RFC3339))
	}

	if len(fields) == 0 {
		http.Error(w, "no fields to update", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf(`
		UPDATE tasksmaster
		SET %s 
		WHERE task_id = ?`,
		strings.Join(fields, ", "),
	)
	args = append(args, idstr)

	if _, err := db.GetDBInfo().E(query, args...); err != nil {
		logger.Error(err, "UpdateTask ~ query execution failed")
		http.Error(w, "task updation failed", http.StatusInternalServerError)
		return
	}

	resp := models.GenricTaskResponse{
		TaskID:  int64(id),
		Message: "Task updated",
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error(err, "UpdateTask ~ json encoding failed")
		http.Error(w, "updating task failed", http.StatusInternalServerError)
		return
	}
}

// DeleteTask godoc
// @Summary      Delete a task by ID
// @Description  Deletes a task from the database using its unique ID. Returns 404 if the task does not exist.
// @Tags         Tasks
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  models.GenricTaskResponse  "Task deleted successfully"
// @Failure      400  {string}  string  "Invalid or missing task ID"
// @Failure      404  {string}  string  "Task not found"
// @Failure      500  {string}  string  "Internal server error"
// @Router       /v1/tasks/{id} [delete]
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	helper.SetJSONHeader(w)

	idstr, exists := mux.Vars(r)["id"]
	if !exists {
		logger.Error("DeleteTask ~ id missing")
		http.Error(w, "id missing", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	query := `
		DELETE FROM tasksmaster
		WHERE task_id = ?
	`

	result, err := db.GetDBInfo().E(query, id)
	if err != nil {
		logger.Error(err, "DeleteTask ~ delete query failed")
		http.Error(w, "deleting task failed", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	resp := models.GenricTaskResponse{
		TaskID:  int64(id),
		Message: "Task deleted",
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.Error(err, "DeleteTask ~ json encoding failed")
		http.Error(w, "deleting task failed", http.StatusInternalServerError)
		return
	}
}
