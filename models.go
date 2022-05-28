package main

import (
	"fmt"
	"math/rand"
)

type Todo struct {
	Id        int    `json:"id"`
	Desc      string `json:"description"`
	Completed bool   `json:"completed"`
}

type TodoStore struct {
	todos []Todo
}

func (t *TodoStore) GetAll() []Todo {
	return t.todos
}

func (t *TodoStore) GetByStatus(completed bool) []Todo {
	var newTodos []Todo = make([]Todo, 0)
	for _, todo := range t.todos {
		if todo.Completed == completed {
			newTodos = append(newTodos, todo)
		}
	}
	return newTodos
}

func (t *TodoStore) GetTodoDetail(id int) (*Todo, error) {
	for _, todo := range t.todos {
		if todo.Id == id {
			return &todo, nil
		}
	}
	return nil, fmt.Errorf("Id not found")
}

func (t *TodoStore) AddTodo(desc string) {
	t.todos = append(t.todos, Todo{
		Id:        rand.Int(),
		Desc:      desc,
		Completed: false,
	})
}

func (t *TodoStore) SetCompleted(id int, completed bool) error {
	for i, todo := range t.todos {
		if todo.Id == id {
			t.todos[i].Completed = completed
			return nil
		}
	}
	return fmt.Errorf("%d is not found", id)
}

func (t *TodoStore) Update(id int, desc string) error {
	for i, todo := range t.todos {
		if todo.Id == id {
			t.todos[i].Desc = desc
			return nil
		}
	}
	return fmt.Errorf("%d is not found", id)
}

func (t *TodoStore) Delete(id int) error {
	var index int = -1
	for i, todo := range t.todos {
		if todo.Id == id {
			index = i
		}
	}
	if index == -1 {
		return fmt.Errorf("%d is not found", id)
	}
	t.todos = append(t.todos[:index], t.todos[index+1:]...)
	return nil
}
