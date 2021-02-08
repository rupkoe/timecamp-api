package parser

import (
	"fmt"
	"github.com/rupkoe/timecamp-api"
	"sort"
	"strings"
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

func PrintProjectTree(tasks []api.Task, project api.Task) {
	TraverseTree(tasks, project, printit)
}

func printit(task api.Task) {
	for i := 1; i < task.LevelParsed(); i++ {
		fmt.Print("--")
	}
	fmt.Println(task.Level, task.Name)
}

func TraverseTree(tasks []api.Task, parent api.Task, callback func(api.Task)) {
	for _, task := range tasks {
		if task.ParentID == parent.TaskID {
			callback(task)
			TraverseTree(tasks, task, callback)
		}
	}
}
