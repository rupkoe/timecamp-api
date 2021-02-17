package parser

import (
	"errors"
	"github.com/rupkoe/timecamp-api"
	"time"
)

// Totals handles totals for spent times
type Totals struct {
	TotalTime    time.Duration
	BillableTime time.Duration
}

func (t Totals) add(total Totals) Totals {
	t.BillableTime = t.BillableTime + total.BillableTime
	t.TotalTime = t.TotalTime + total.TotalTime
	return t
}

// TaskTotals keeps totals for spent times for tasks
type TaskTotals map[string]Totals

// add adds totals to a task total. Create map item if needed.
func (t TaskTotals) add(taskId string, totals Totals) {
	currentTotals, ok := t[taskId]
	if !ok {
		t[taskId] = totals
	} else {
		t[taskId] = currentTotals.add(totals)
	}
}

// Get returns the totals for task
func (t TaskTotals) Get(taskId string) Totals {
	totals, ok := t[taskId]
	if ok {
		return totals
	} else {
		return Totals{0, 0}
	}
}

// GetProjectList return an array of projects - in TimeCamp, projects are simply tasks at the top level.
func GetProjectList(tasks []api.Task) []api.Task {
	var result []api.Task
	for _, task := range tasks {
		if task.IsProject() {
			result = append(result, task)
		}
	}
	return result
}

// GetTaskById returns a task identified by its ID.
func GetTaskById(tasks []api.Task, id string) (*api.Task, error) {
	for _, task := range tasks {
		if task.TaskID == id {
			return &task, nil
		}
	}
	return nil, errors.New("task with ID " + id + " not found")
}

// GetEntriesForTask returns an array with time entries for the given task.
func GetEntriesForTask(entries []api.TimeEntry, taskId string) []api.TimeEntry {
	var result []api.TimeEntry
	for _, entry := range entries {
		if entry.TaskID == taskId {
			result = append(result, entry)
		}
	}
	return result
}

// SummarizeTask summarizes the entries directly related to given task.
func SummarizeTask(task api.Task, entries []api.TimeEntry) (billable time.Duration, total time.Duration, err error) {
	taskEntries := GetEntriesForTask(entries, task.TaskID)
	for _, entry := range taskEntries {
		duration, err := entry.DurationParsed()
		if err != nil {
			return 0, 0, err
		}
		if entry.IsBillable() {
			billable += duration
		}
		total += duration
	}
	return
}

// WalkTaskTree recursively walks down the task tree, starting at a given root task, calling a callback function for every node.
// includeRoot controls if callback is also executed with root task.
func WalkTaskTree(tasks []api.Task, root api.Task, includeRoot bool, callback func(api.Task, map[int]string)) {
	traverseTree(tasks, root, includeRoot, callback)
}

// SummarizeTaskTree recursively walks down the task tree, starting at a given root task, summarizing all recorded times
func SummarizeTaskTree(tasks []api.Task, entries []api.TimeEntry, root api.Task) TaskTotals {
	var taskTotals = make(TaskTotals)
	traverseTree(tasks, root, true, func(task api.Task, parentIds map[int]string) {
		timeEntries := GetEntriesForTask(entries, task.TaskID)
		var taskTimes Totals
		for _, timeEntry := range timeEntries {
			duration, _ := timeEntry.DurationParsed()
			taskTimes.TotalTime += duration
			if timeEntry.IsBillable() {
				taskTimes.BillableTime += duration
			}
		}
		taskTotals.add(task.TaskID, taskTimes)
		for _, taskId := range parentIds {
			taskTotals.add(taskId, taskTimes)
		}
	})
	return taskTotals
}

var parentIds = make(map[int]string)

// traverseTree recursively walks down the task tree, starting at a given root task, calling a callback function for every node
func traverseTree(tasks []api.Task, parent api.Task, includeParent bool, callback func(api.Task, map[int]string)) {
	if includeParent && len(parentIds) == 0 {
		callback(parent, parentIds)
	}
	for _, task := range tasks {
		if task.ParentID == parent.TaskID {
			parentIds[task.LevelParsed()-1] = task.ParentID
			callback(task, parentIds)
			traverseTree(tasks, task, includeParent, callback)
			delete(parentIds, task.LevelParsed()-1)
		}
	}
}
