package parser

import (
	"errors"
	"fmt"
	"github.com/rupkoe/timecamp-api"
	"sort"
	"strings"
	"time"
)

type ProjectNode struct {
	Project api.Task
	Tasks   []api.Task
}

// GetProjectList return an array of projects - in TimeCamp, projects are simply tasks at the top level
func GetProjectList(tasks []api.Task) []api.Task {
	var result []api.Task
	for _, task := range tasks {
		if task.IsProject() {
			result = append(result, task)
		}
	}
	return result
}

func GetTaskById(tasks []api.Task, id string) (*api.Task, error) {
	for _, task := range tasks {
		if task.TaskID == id {
			return &task, nil
		}
	}
	return nil, errors.New("task with ID " + id + " not found")
}

type Totals struct {
	TotalTime    time.Duration
	BillableTime time.Duration
}

func (t Totals) add(total Totals) Totals {
	t.BillableTime = t.BillableTime + total.BillableTime
	t.TotalTime = t.TotalTime + total.TotalTime
	return t
}

type TaskTotals map[string]Totals

func (tt TaskTotals) add(taskId string, totals Totals) {
	currentTotals, prs := tt[taskId]
	if !prs {
		tt[taskId] = totals
	} else {
		tt[taskId] = currentTotals.add(totals)
	}
}

// GetProjects returns projects and corresponding tasks, sorted alphabetically.
func GetProjectTree(tasks []api.Task) ([]ProjectNode, error) {
	var result []ProjectNode
	for _, proj := range tasks {
		if proj.IsProject() {
			p := ProjectNode{
				Project: proj,
			}
			for _, task := range tasks {
				if task.ParentID == proj.TaskID {
					p.Tasks = append(p.Tasks, task)
				}
			}
			sort.Slice(p.Tasks, func(i, j int) bool {
				return strings.ToUpper(p.Tasks[i].Name) < strings.ToUpper(p.Tasks[j].Name)
			})
			result = append(result, p)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return strings.ToUpper(result[i].Project.Name) < strings.ToUpper(result[j].Project.Name)
	})

	return result, nil
}

// https://stackoverflow.com/questions/22957638/make-a-tree-from-a-table-using-golang

var parentIds = make(map[int]string)

func TraverseTree(tasks []api.Task, parent api.Task, callback func(api.Task, map[int]string)) {
	for _, task := range tasks {
		if task.ParentID == parent.TaskID {
			parentIds[task.LevelParsed()-1] = task.ParentID
			callback(task, parentIds)
			TraverseTree(tasks, task, callback)
			delete(parentIds, task.LevelParsed()-1)
		}
	}
}

func Summarize(tasks []api.Task, entries []api.TimeEntry, root api.Task) {
	TraverseTree(tasks, root, func(task api.Task, parentIds map[int]string) {
		timeEntries := GetEntriesForTask(entries, task.TaskID)
		var taskTimes Totals
		for _, timeEntry := range timeEntries {
			duration, _ := timeEntry.DurationParsed()
			if timeEntry.IsBillable() {
				taskTimes.BillableTime += duration
			} else {
				taskTimes.TotalTime += duration
			}

		}
		fmt.Printf("Total for Task %v: Total %v, Billable %v", task.TaskID, taskTimes.TotalTime, taskTimes.BillableTime)
		fmt.Println()
		// todo: add to TaskTotals based on the provided parentIds array
	})
}

func GetEntriesForTask(entries []api.TimeEntry, taskID string) []api.TimeEntry {
	var result []api.TimeEntry
	for _, entry := range entries {
		if entry.TaskID == taskID {
			result = append(result, entry)
		}
	}
	return result
}
