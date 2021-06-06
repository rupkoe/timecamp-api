package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
)

// Task maps the JSON returned by TimeCamp API /tasks.
//
// API Docs: https://github.com/timecamp/timecamp-api/blob/master/sections/tasks.md
// Created with https://mholt.github.io/json-to-go/
type Task struct {
	TaskID           int         `json:"task_id"`
	ParentID         int         `json:"parent_id"`
	AssignedBy       string      `json:"assigned_by"`
	Name             string      `json:"name"`
	ExternalTaskID   string      `json:"external_task_id"`
	ExternalParentID string      `json:"external_parent_id"`
	Level            string      `json:"level"`
	Archived         string      `json:"archived"`
	Tags             string      `json:"tags"`
	Budgeted         string      `json:"budgeted"`
	BudgetUnit       string      `json:"budget_unit"`
	RootGroupID      string      `json:"root_group_id"`
	Billable         string      `json:"billable"`
	Note             string      `json:"note"`
	PublicHash       string      `json:"public_hash"`
	AddDate          string      `json:"add_date"`
	ModifyTime       string      `json:"modify_time"`
	Color            string      `json:"color"`
	Users            interface{} `json:"users"`
	UserAccessType   int         `json:"user_access_type"`
}

// IsProject is true if task is a project (=top-level task) in TimeCamp
func (t Task) IsProject() bool {
	// Alternatively `t.Level == "1"` could be used to identify project tasks
	return t.ParentID == 0
}

// LevelParsed converts the api's string value into a numeric value.
func (t Task) LevelParsed() int {
	level, err := strconv.Atoi(t.Level)
	if err != nil {
		log.Fatal("malformed task - error converting level to int")
	}
	return level
}

type TaskParams struct {
	OnlyArchivedTasks bool
	OnlyActiveTasks   bool
}

// GetTasks wraps the "GET /tasks" api endpoint.
// Both "Projects" and "Tasks" in TimeCamp's UI are tasks.
func GetTasks(c Connection, params TaskParams) ([]Task, error) {
	queryUrl, err := taskUrl(c, params)
	if err != nil {
		return nil, err
	}

	data, err := httpGet(queryUrl)
	if err != nil {
		return nil, err
	}

	// The returned json contains dynamic keys for the task object.
	// It therefore can only be unmarshalled into a map.
	var result map[string]Task
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	// The key is redundant, strip it.
	var tasks []Task
	for _, t := range result {
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func taskUrl(connection Connection, params TaskParams) (string, error) {
	var exclude string
	if params.OnlyActiveTasks && params.OnlyArchivedTasks {
		return "", fmt.Errorf("at least one of active or archived tasks must be included")
	} else if params.OnlyActiveTasks {
		exclude = "exclude_archived=0"
	} else if params.OnlyArchivedTasks {
		exclude = "exclude_archived=1"
	} else {
		exclude = "" //nothing excluded
	}

	queryUrl, err := url.Parse(connection.ApiUrl + "/tasks/format/json/api_token/" + connection.Token + "?" + exclude)
	if err != nil {
		return "", err
	}
	return queryUrl.String(), err
}
