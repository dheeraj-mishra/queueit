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
	"strings"
	"time"
)

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
	json.NewEncoder(w).Encode(tasks)
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	helper.SetTextHeader(w)

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
	if ctr.Priority == 0 || !(ctr.Priority == 1 || ctr.Priority == 2 || ctr.Priority == 3) {
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

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Task created: %d", taskID)))
}
