package api

import (
	"errors"
	"sort"
	"strings"
)

type Project struct {
	Project Task
	Tasks   []Task
}

// GetProjects returns projects and corresponding tasks, sorted alphabetically.
func GetProjects(connection Connection, params TaskParams) ([]Project, error) {
	tasks, err := GetTasks(connection, params)
	if err != nil {
		return nil, errors.New("Unable to get tasks/result")
	}

	var result []Project
	for _, proj := range tasks {
		if proj.IsProject() {
			p := Project{
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
