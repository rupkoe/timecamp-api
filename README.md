# Go TimeCamp API Client

![Go](https://github.com/rupkoe/timecamp-api/workflows/Go/badge.svg)

A partial client for the [TimeCamp API](https://github.com/timecamp/timecamp-api).

The **api** package currently supports fetching tasks and time entries from the api. 
The optional **parser** package helps with processing the retrieved data.

---

⚠️ *As I do not use Timecamp API right now, the repo is not seeing any fixes or feature updates.*

---

## API

Provides raw data from the TimeCamp API. Supports currently only fetching tasks and time entries.

### Sample Usage

```Go
package main

import 	(
	"github.com/rupkoe/timecamp-api"
	"time"
)

func main() {
    connection := api.Connection{
        ApiUrl: "https://www.timecamp.com/third_party/api",
        Token:  "secret-token",
    }

    tasks, err := api.GetTasks(connection, api.TaskParams{})
    //...
    
    from, _ := time.Parse(time.RFC822, "01 Jan 21 00:00 CET")
    to, _ := time.Parse(time.RFC822, "31 Jan 21 00:00 CET")
    entries, err := api.GetTimeEntries(connection, api.TimeEntryParams{From: from, To: to})
    //...
}
```

## Parser

The parser package allows to work on data retrieved from the TimeCamp API.

As a prerequisite, get the desired data from the API. Use filters to your liking.

```go
// get time entries from api
timeEntries, err := api.GetTimeEntries(connection, api.TimeEntryParams{})
if err != nil {
    ...
}

// get tasks from api
tasks, err := api.GetTasks(connection, api.TaskParams{})
if err != nil {
    ...
}
```

If you retrieved all tasks, you may get the projects (=top-level tasks) by using

	projectList := parser.GetProjectList(tasks)

Walk the task tree starting at a given (root) node, providing a callback function that is called for every task being visited.

```go
parser.WalkTaskTree(tasks, project, true, printit)

func printit(task api.Task, parentIds map[int]string) {
    for i := 1; i < task.LevelParsed(); i++ {
        fmt.Print("-")
    }
    fmt.Println(task.Level, task.TaskID, task.Name, parentIds)
}
```
    

Summarize the times spent on tasks and return a map with total and billable times per task, summarized all the way up to the root node.

    tasktotals := parser.SummarizeTaskTree(tasks, timeEntries, project)


## TimeCamp API Oddness 

Documents unexpected behaviour of the TimeCamp API for further reference / future development.

- API returns different `id` value types:
    - string for Task
    - number for TimeEntry
- Task JSON has variable, redundant keys (TaskID)
