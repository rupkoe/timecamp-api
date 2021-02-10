# TimeCamp API Parser Package

The parser package allows to work on data retrieved from the TimeCamp API.

As a prerequisite, get the desired data from the API. Use filters to your liking.

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

If you retrived all tasks, you may get the projects (i.e. top-level tasks) by using

	projectList := parser.GetProjectList(tasks)

Walk the task tree starting at a given (root) node, providing a callback function that is called for every task being visited.

    parser.WalkTaskTree(tasks, project, printit)

    func printit(task api.Task, parentIds map[int]string) {
        for i := 1; i < task.LevelParsed(); i++ {
            fmt.Print("--")
        }
        fmt.Println(task.Level, task.TaskID, task.Name, parentIds)
    }

Summarize the times spent on tasks and return a map with total and billable times per task, summarized all the way up to the root node.

		tasktotals := parser.SummarizeTaskTree(tasks, timeEntries, project)

