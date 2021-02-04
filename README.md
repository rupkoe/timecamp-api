# TimeCamp API Client

A partial client for the [TimeCamp API](https://github.com/timecamp/timecamp-api).

It currently supports fetching tasks and time entries. 

## Sample Usage

Minimal example:

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


## TimeCamp API Oddness 

Documents unexpected behaviour of the TimeCamp API for further reference / future development.

- API returns different `id` value types:
    - string for Task
    - number for TimeEntry
- Task JSON has variable, redundant keys (TaskID)
