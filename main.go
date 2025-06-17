package main

import (
	"fmt"
	"os"

	"Todo_Cli/todo"

	"github.com/spf13/pflag"
)

const fileName = "./data/todo.json"

func main() {
	// Flag variables
	var addTask string
	var listTasks bool
	var doneTask int
	var removeTask int

	pflag.StringVarP(&addTask, "add", "a", "", "Add a new task")
	pflag.BoolVarP(&listTasks, "list", "l", false, "List all tasks")
	pflag.IntVarP(&doneTask, "done", "d", -1, "Done a task")
	pflag.IntVarP(&removeTask, "remove", "r", -1, "Remove a task")
	pflag.Parse()


	switch {
	case addTask != "":
		// Load todo list
		todos, err := todo.LoadTodos(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load Todo list: %v", err)
			return
		}	

		// Append the task on the todo list
		todos = append(todos, todo.Todo{addTask, false})

		// Recreate the json file
		err = todo.SaveTodos(fileName, todos)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load Todo list: %v", err)
			return
		}

		fmt.Fprintf(os.Stdout, "Success to add Task\n")

	case listTasks != false:
		// Loading Todo list
		todos, err := todo.LoadTodos(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load Todo list: %v", err)
			return
		}

		// Display the list of Todos
		for i, t := range todos {
			fmt.Println(i + 1, t.Task, '\t', t.Done)
		}
	case doneTask != -1:
		// Loading Todo List
		todos, err := todo.LoadTodos(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load Todo list: %v", err)
			return	
		}

		// Get the given todo by line number
		todos[doneTask - 1].Done = true

		// Save todo list
		err = todo.SaveTodos(fileName, todos)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load Todo list: %v", err)
			return
		}

	case removeTask != -1:
		// Load the todo list
		todos, err := todo.LoadTodos(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Falied to load Todo list: %v", err)
			return
		}

		// Remove the target element
		todos = append(todos[:removeTask - 1], todos[removeTask:]...)

		// Save todo list
		err = todo.SaveTodos(fileName, todos)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load Todo list: %v", err)
			return
		}
	}
}