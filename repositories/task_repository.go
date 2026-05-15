package repositories

import (
	"trithong.com/task-golang/entities"
	"errors"
)

type TaskRepository interface {
	Create(task entities.Task) entities.Task
	GetByID(id int) (entities.Task, error)
	GetAll() []entities.Task
	Update(id int, task entities.Task) (entities.Task, error)
	Delete(id int) error
}

type taskRepository struct {
	tasks  map[int]entities.Task
	nextID int
}

func NewTaskRepository() TaskRepository {
	return &taskRepository{
		tasks:  make(map[int]entities.Task),
		nextID: 1,
	}
}

func (r *taskRepository) Create(task entities.Task) entities.Task {
	task.ID = r.nextID
	r.tasks[task.ID] = task
	r.nextID++
	return task
}

func (r *taskRepository) GetByID(id int) (entities.Task, error) {
	task, exists := r.tasks[id]

	if !exists {
		return entities.Task{}, errors.New("task not found")
	}

	return task, nil
}

func (r *taskRepository) GetAll() []entities.Task {
	var list []entities.Task

	for _, task := range r.tasks {
		list = append(list, task)
	}

	return list
}

func (r *taskRepository) Update(id int, task entities.Task) (entities.Task, error) {
	_, exists := r.tasks[id]

	if !exists {
		return entities.Task{}, errors.New("task not found")
	}

	task.ID = id
	r.tasks[id] = task

	return task, nil
}

func (r *taskRepository) Delete(id int) error {
	_, exists := r.tasks[id]

	if !exists {
		return errors.New("task not found")
	}

	delete(r.tasks, id)
	return nil
}