package api

import (
	"testing"
	"time"
)

func Test_timeEntryUrl(t *testing.T) {
	type args struct {
		connection Connection
		params     TimeEntryParams
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Minimal Data",
			args: args{
				Connection{
					ApiUrl: "URL",
					Token:  "TOKEN",
				},
				TimeEntryParams{
					From:  time.Time{},
					To:    time.Time{}.Add(time.Hour),
					Tasks: nil,
				},
			},
			want:    "URL/entries/format/json/api_token/TOKEN/from/0001-01-01/to/0001-01-01/task_ids/",
			wantErr: false,
		}, {
			name: "To Date Before From Date",
			args: args{
				Connection{
					ApiUrl: "URL",
					Token:  "TOKEN",
				},
				TimeEntryParams{
					From:  time.Date(2021, 01, 01, 0, 0, 0, 0, time.Local),
					To:    time.Date(2020, 01, 01, 0, 0, 0, 0, time.Local),
					Tasks: nil,
				},
			},
			want:    "",
			wantErr: true,
		}, {
			name: "With Tasks",
			args: args{
				Connection{
					ApiUrl: "URL",
					Token:  "TOKEN",
				},
				TimeEntryParams{
					From: time.Time{},
					To:   time.Time{}.Add(time.Hour),
					Tasks: []Task{
						{TaskID: 1},
						{TaskID: 2},
					},
				},
			},
			want:    "URL/entries/format/json/api_token/TOKEN/from/0001-01-01/to/0001-01-01/task_ids/1,2",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := timeEntryUrl(tt.args.connection, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("timeEntryUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("timeEntryUrl() got = %v, want %v", got, tt.want)
			}
		})
	}
}
