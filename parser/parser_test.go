package parser

import (
	api "github.com/rupkoe/timecamp-api"
	"reflect"
	"testing"
)

// todo: use data from rupert's fixture file

func TestGetProjectList(t *testing.T) {
	type args struct {
		tasks []api.Task
	}
	tests := []struct {
		name string
		args args
		want []api.Task
	}{
		{
			name: "Minimal Test",
			args: args{
				[]api.Task{
					{
						Name:     "Task A",
						ParentID: "0",
					},
					{
						Name:     "Task A-B",
						ParentID: "A",
					},
					{
						Name:     "Task C",
						ParentID: "0",
					},
				},
			},
			want: []api.Task{
				{
					Name:     "Task A",
					ParentID: "0",
				},
				{
					Name:     "Task C",
					ParentID: "0",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetProjectList(tt.args.tasks); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProjects() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_TraverseTree(t *testing.T) {
	type args struct {
		tasks    []api.Task
		parent   api.Task
		callback func(api.Task, map[int]string)
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Minimal Test",
			args: args{
				[]api.Task{
					{
						Name:     "Task A",
						TaskID:   "A",
						ParentID: "0",
						Level:    "1",
					},
					{
						Name:     "Task A-B",
						TaskID:   "A-B",
						ParentID: "A",
						Level:    "2",
					},
					{
						Name:     "Task A-B-C",
						TaskID:   "A-B-C",
						ParentID: "A-B",
						Level:    "3",
					},
					{
						Name:     "Task C",
						TaskID:   "C",
						ParentID: "0",
						Level:    "1",
					}},
				api.Task{
					Name:     "Task A",
					TaskID:   "A",
					ParentID: "0",
					Level:    "1",
				},
				func(task api.Task, parentIds map[int]string) {
					// todo: check for parentIds here as well?
					ok := task.TaskID == "A-B" || task.TaskID == "A-B-C"
					if !ok {
						t.Errorf("invalid (sub-)task %v", task.TaskID)
					}
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traverseTree(tt.args.tasks, tt.args.parent, tt.args.callback)
		})
	}
}
