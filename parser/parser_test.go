package parser

import (
	"reflect"
	"testing"
	"time"

	api "github.com/rupkoe/timecamp-api"
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
						ParentID: 0,
					},
					{
						Name:     "Task A-B",
						ParentID: 1,
					},
					{
						Name:     "Task C",
						ParentID: 0,
					},
				},
			},
			want: []api.Task{
				{
					Name:     "Task A",
					ParentID: 0,
				},
				{
					Name:     "Task C",
					ParentID: 0,
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
		expectedIDs   []int
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
						Name:     "Task 1",
						TaskID:   1,
						ParentID: 0,
						Level:    1,
					},
					{
						Name:     "Task 1-2",
						TaskID:   12,
						ParentID: 1,
						Level:    2,
					},
					{
						Name:     "Task 1-2-3",
						TaskID:   123,
						ParentID: 12,
						Level:    3,
					},
					{
						Name:     "Task C",
						TaskID:   3,
						ParentID: 0,
						Level:    1,
					}},
				api.Task{
					Name:     "Task A",
					TaskID:   1,
					ParentID: 0,
					Level:    1,
				},
				false,
				[]int{12, 123},
			},
		}, {
			name: "Without parent task",
			args: args{
				[]api.Task{
					{
						Name:     "Task A",
						TaskID:   1,
						ParentID: 0,
						Level:    1,
					},
					{
						Name:     "Task 1-2",
						TaskID:   12,
						ParentID: 1,
						Level:    2,
					},
					{
						Name:     "Task 1-2-3",
						TaskID:   123,
						ParentID: 12,
						Level:    3,
					},
					{
						Name:     "Task C",
						TaskID:   3,
						ParentID: 0,
						Level:    1,
					}},
				api.Task{
					Name:     "Task A",
					TaskID:   1,
					ParentID: 0,
					Level:    1,
				},
				true,
				[]int{1, 12, 123},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskIds := map[int]struct{}{}
			traverseTree(tt.args.tasks, tt.args.parent, tt.args.includeParent, func(task api.Task, m map[int]int) {
				_, exists := taskIds[task.TaskID]
				if !exists {
					taskIds[task.TaskID] = struct{}{}
				} else {
					t.Errorf("Each TaskID should only show up once, but found %d more often", task.TaskID)
				}
			})
			for _, expected := range tt.args.expectedIDs {
				if _, exists := taskIds[expected]; !exists {
					t.Errorf("Expected task %d was not returned", expected)
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
						TaskID:   1,
						ParentID: 0,
						Level:    1,
					},
				},
				0, //"Task A"
				[]api.TimeEntry{
					{
						ID:          1,
						Duration:    "600",
						TaskID:      "1",
						Billable:    1,
						Description: "",
					}, {
						ID:          2,
						Duration:    "300",
						TaskID:      "1",
						Billable:    0,
						Description: "",
					},
				},
				TaskTotals{
					1: {
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
						Name:     "Task 1",
						TaskID:   1,
						ParentID: 0,
						Level:    1,
					}, {
						Name:     "Task 1-2",
						TaskID:   11,
						ParentID: 1,
						Level:    2,
					}, {
						Name:     "Task 1-1-1",
						TaskID:   111,
						ParentID: 11,
						Level:    3,
					}, {
						Name:     "Task 1-1-1-1",
						TaskID:   1111,
						ParentID: 111,
						Level:    4,
					}, {
						Name:     "Task 1-1-2",
						TaskID:   112,
						ParentID: 11,
						Level:    3,
					}, {
						Name:     "Task 1-2",
						TaskID:   12,
						ParentID: 1,
						Level:    2,
					},
				},
				1, // "Task A-A"
				[]api.TimeEntry{
					{
						ID:          1,
						Duration:    "600",
						TaskID:      "1",
						Billable:    1,
						Description: "",
					}, {
						ID:          2,
						Duration:    "600",
						TaskID:      "11",
						Billable:    1,
						Description: "",
					}, {
						ID:          2,
						Duration:    "600",
						TaskID:      "111",
						Billable:    1,
						Description: "",
					}, {
						ID:          2,
						Duration:    "600",
						TaskID:      "1111",
						Billable:    1,
						Description: "",
					}, {
						ID:          2,
						Duration:    "600",
						TaskID:      "112",
						Billable:    1,
						Description: "",
					}, {
						ID:          3,
						Duration:    "600",
						TaskID:      "12",
						Billable:    1,
						Description: "",
					},
				},
				TaskTotals{
					11: {
						TotalTime:    2400 * time.Second,
						BillableTime: 2400 * time.Second,
					},
					111: {
						TotalTime:    1200 * time.Second,
						BillableTime: 1200 * time.Second,
					},
					1111: {
						TotalTime:    600 * time.Second,
						BillableTime: 600 * time.Second,
					},
					112: {
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

func TestSummarizeTask(t *testing.T) {
	type args struct {
		task    api.Task
		entries []api.TimeEntry
	}
	tests := []struct {
		name         string
		args         args
		wantBillable time.Duration
		wantTotal    time.Duration
		wantErr      bool
	}{
		{
			name: "Basic test",
			args: args{
				task: api.Task{
					Name:     "Task A",
					TaskID:   1,
					ParentID: 0,
					Level:    1,
				},
				entries: []api.TimeEntry{
					{
						ID:          1,
						Duration:    "1800",
						TaskID:      "1",
						Billable:    1,
						Description: "",
					}, {
						ID:          2,
						Duration:    "3600",
						TaskID:      "1",
						Billable:    0,
						Description: "",
					},
				},
			},
			wantBillable: 30 * time.Minute,
			wantTotal:    90 * time.Minute,
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotBillable, gotTotal, err := SummarizeTask(tt.args.task, tt.args.entries)
			if (err != nil) != tt.wantErr {
				t.Errorf("SummarizeTask() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotBillable != tt.wantBillable {
				t.Errorf("SummarizeTask() gotBillable = %v, want %v", gotBillable, tt.wantBillable)
			}
			if gotTotal != tt.wantTotal {
				t.Errorf("SummarizeTask() gotTotal = %v, want %v", gotTotal, tt.wantTotal)
			}
		})
	}
}
