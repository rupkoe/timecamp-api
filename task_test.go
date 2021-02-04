package api

import "testing"

func Test_taskUrl(t *testing.T) {
	type args struct {
		c Connection
		p TaskParams
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "All Tasks",
			args: args{
				c: Connection{ApiUrl: "http://apiurl", Token: "TOKEN"},
				p: TaskParams{
					OnlyArchivedTasks: false,
					OnlyActiveTasks:   false,
				},
			},
			want:    "http://apiurl/tasks/format/json/api_token/TOKEN?",
			wantErr: false,
		}, {
			name: "Archived Tasks Filter",
			args: args{
				c: Connection{ApiUrl: "http://apiurl", Token: "TOKEN"},
				p: TaskParams{
					OnlyArchivedTasks: true,
					OnlyActiveTasks:   false,
				},
			},
			want:    "http://apiurl/tasks/format/json/api_token/TOKEN?exclude_archived=1",
			wantErr: false,
		}, {
			name: "Active Tasks Filter",
			args: args{
				c: Connection{ApiUrl: "http://apiurl", Token: "TOKEN"},
				p: TaskParams{
					OnlyArchivedTasks: false,
					OnlyActiveTasks:   true,
				},
			},
			want:    "http://apiurl/tasks/format/json/api_token/TOKEN?exclude_archived=0",
			wantErr: false,
		}, {
			name: "Too Many Filters",
			args: args{
				c: Connection{ApiUrl: "http://apiurl", Token: "TOKEN"},
				p: TaskParams{
					OnlyArchivedTasks: true,
					OnlyActiveTasks:   true,
				},
			},
			want:    "",
			wantErr: true,
		}, {
			name: "URL Character Quoting",
			args: args{
				c: Connection{ApiUrl: "apiurl//{query}[#]", Token: "TOKEN"},
				p: TaskParams{
					OnlyArchivedTasks: false,
					OnlyActiveTasks:   false,
				},
			},
			want:    "apiurl//%7Bquery%7D%5B#]/tasks/format/json/api_token/TOKEN?",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := taskUrl(tt.args.c, tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("taskUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("taskUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}
