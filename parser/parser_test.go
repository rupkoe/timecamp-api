package parser

import (
	api "github.com/rupkoe/timecamp-api"
	"reflect"
	"testing"
	"time"
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

func TestTraverseTree(t *testing.T) {
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

func TestSummarizeTaskTree(t *testing.T) {
	type args struct {
		tasks         []api.Task
		parentTaskIdx int
		entries       []api.TimeEntry
		expected      TaskTotals
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Billable and total time",
			args: args{
				[]api.Task{
					{
						Name:     "Task A",
						TaskID:   "A",
						ParentID: "0",
						Level:    "1",
					},
				},
				0, //"Task A"
				[]api.TimeEntry{
					{
						ID:          1,
						Duration:    "600",
						TaskID:      "A",
						Billable:    1,
						Description: "",
					}, {
						ID:          2,
						Duration:    "300",
						TaskID:      "A",
						Billable:    0,
						Description: "",
					},
				},
				TaskTotals{
					"A": {
						TotalTime:    900 * time.Second,
						BillableTime: 600 * time.Second,
					},
				},
			},
		}, {
			name: "Task tree totals",
			args: args{
				[]api.Task{
					{
						Name:     "Task A",
						TaskID:   "A",
						ParentID: "0",
						Level:    "1",
					}, {
						Name:     "Task A-A",
						TaskID:   "A-A",
						ParentID: "A",
						Level:    "2",
					}, {
						Name:     "Task A-A-A",
						TaskID:   "A-A-A",
						ParentID: "A-A",
						Level:    "3",
					}, {
						Name:     "Task A-A-A-A",
						TaskID:   "A-A-A-A",
						ParentID: "A-A-A",
						Level:    "4",
					}, {
						Name:     "Task A-A-B",
						TaskID:   "A-A-B",
						ParentID: "A-A",
						Level:    "3",
					}, {
						Name:     "Task A-B",
						TaskID:   "A-B",
						ParentID: "A",
						Level:    "2",
					},
				},
				1, // "Task A-A"
				[]api.TimeEntry{
					{
						ID:          1,
						Duration:    "600",
						TaskID:      "A",
						Billable:    1,
						Description: "",
					}, {
						ID:          2,
						Duration:    "600",
						TaskID:      "A-A",
						Billable:    1,
						Description: "",
					}, {
						ID:          2,
						Duration:    "600",
						TaskID:      "A-A-A",
						Billable:    1,
						Description: "",
					}, {
						ID:          2,
						Duration:    "600",
						TaskID:      "A-A-A-A",
						Billable:    1,
						Description: "",
					}, {
						ID:          2,
						Duration:    "600",
						TaskID:      "A-A-B",
						Billable:    1,
						Description: "",
					}, {
						ID:          3,
						Duration:    "600",
						TaskID:      "A-B",
						Billable:    1,
						Description: "",
					},
				},
				TaskTotals{
					"A-A": {
						TotalTime:    2400 * time.Second,
						BillableTime: 2400 * time.Second,
					},
					"A-A-A": {
						TotalTime:    1200 * time.Second,
						BillableTime: 1200 * time.Second,
					},
					"A-A-A-A": {
						TotalTime:    600 * time.Second,
						BillableTime: 600 * time.Second,
					},
					"A-A-B": {
						TotalTime:    600 * time.Second,
						BillableTime: 600 * time.Second,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			totals := SummarizeTaskTree(tt.args.tasks, tt.args.entries, tt.args.tasks[tt.args.parentTaskIdx])
			if !reflect.DeepEqual(tt.args.expected, totals) {
				t.Errorf("calculated and expected totals do not match")
			}
		})

	}
}
