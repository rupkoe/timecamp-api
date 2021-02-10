package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"
)

// TimeEntry maps the JSON  returned by TimeCamp API /entries.
//
// API docs: https://github.com/timecamp/timecamp-api/blob/master/sections/time-entries.md
// Created with https://mholt.github.io/json-to-go/
type TimeEntry struct {
	ID               int    `json:"id"`
	Duration         string `json:"duration"`
	UserID           string `json:"user_id"`
	UserName         string `json:"user_name"`
	TaskID           string `json:"task_id"`
	LastModify       string `json:"last_modify"`
	Date             string `json:"date"`
	StartTime        string `json:"start_time"`
	EndTime          string `json:"end_time"`
	Locked           string `json:"locked"`
	Name             string `json:"name"`
	AddonsExternalID string `json:"addons_external_id"`
	Billable         int    `json:"billable"`
	InvoiceID        string `json:"invoiceId"`
	Color            string `json:"color"`
	Description      string `json:"description"`
}

func (e TimeEntry) DurationParsed() (time.Duration, error) {
	secs := strings.Join([]string{e.Duration, "s"}, "")
	dur, err := time.ParseDuration(secs)
	if err != nil {
		return 0, err
	}
	return dur, nil
}

func (e TimeEntry) DateParsed() time.Time {
	date, err := time.Parse(DateFormat, e.Date)
	if err != nil {
		log.Fatal("could not parse entry date:", e.Date)
	}
	return date
}

func (e TimeEntry) HasDescription() bool {
	return len(strings.Trim(e.Description, " ")) > 0
}

func (e TimeEntry) IsBillable() bool {
	return e.Billable > 0
}

// TimeEntryParams query parameters.
type TimeEntryParams struct {
	From  time.Time
	To    time.Time
	Tasks []Task
}

// GetTimeEntries wraps the "GET /entries" api endpoint.
// If params.Tasks is nil / empty, all tasks' entries are returned.
func GetTimeEntries(con Connection, params TimeEntryParams) ([]TimeEntry, error) {
	queryUrl, err := timeEntryUrl(con, params)
	if err != nil {
		return nil, err
	}

	data, err := httpGet(queryUrl)
	if err != nil {
		return nil, err
	}

	var result []TimeEntry
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func timeEntryUrl(connection Connection, params TimeEntryParams) (string, error) {
	if !params.From.Before(params.To) {
		return "", fmt.Errorf("GetTimeEntries: From date must be before To date")
	}

	var taskIds []string
	for _, task := range params.Tasks {
		taskIds = append(taskIds, task.TaskID)
	}

	queryUrl, err := url.Parse(connection.ApiUrl + "/entries/format/json/api_token/" + connection.Token + "/from/" +
		params.From.Format(DateFormat) + "/to/" + params.To.Format(DateFormat) +
		"/task_ids/" + strings.Join(taskIds, ","))
	if err != nil {
		return "", err
	}

	return queryUrl.String(), err
}
