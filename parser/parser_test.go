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
		tasks         []api.Task
		parent        api.Task
		includeParent bool
		expectedIDs   []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "With parent task",
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
				false,
				[]string{"A-B", "A-B-C"},
			},
		}, {
			name: "Without parent task",
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
				true,
				[]string{"A", "A-B", "A-B-C"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskIds := map[string]struct{}{}
			traverseTree(tt.args.tasks, tt.args.parent, tt.args.includeParent, func(task api.Task, m map[int]string) {
				_, exists := taskIds[task.TaskID]
				if !exists {
					taskIds[task.TaskID] = struct{}{}
				} else {
					t.Errorf("Each TaskID should only show up once, but found %s more often", task.TaskID)
				}
			})
			for _, expected := range tt.args.expectedIDs {
				if _, exists := taskIds[expected]; !exists {
					t.Errorf("Expected task %s was not returned", expected)
				}
			}
			if len(tt.args.expectedIDs) != len(taskIds) {
				t.Errorf("Number of returned task IDs (%d) is not equal number of exptected task IDs (%d)",
					len(taskIds), len(tt.args.expectedIDs))
			}
		})
	}
}
