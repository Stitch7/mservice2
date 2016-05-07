package main

import "fmt"

var currentTodoId int

var todos Todos

// Give us some seed data
func init() {
	TodoRepoCreate(Todo{Name: "Write presentation"})
	TodoRepoCreate(Todo{Name: "Host meetup"})
}

func TodoRepoFindById(id int) Todo {
	for _, t := range todos {
		if t.Id == id {
			return t
		}
	}
	// return empty Todo if not found
	return Todo{}
}

//this is bad, I don't think it passes race condtions
func TodoRepoCreate(t Todo) Todo {
	currentTodoId += 1
	t.Id = currentTodoId
	todos = append(todos, t)
	return t
}

func TodoRepoDeleteById(id int) error {
	for i, t := range todos {
		if t.Id == id {
			todos = append(todos[:i], todos[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Todo with id of %d to delete", id)
}
