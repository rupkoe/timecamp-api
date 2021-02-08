package parser

import (
	api "github.com/rupkoe/timecamp-api"
	"reflect"
	"testing"
)

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
