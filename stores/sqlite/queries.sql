--  Copyright (c) 2024 Michael D Henderson. All rights reserved.

-- --------------------------------------------------------------------------
-- InsertUser inserts a new User.
-- Since it allows the caller to specify the ID, it should only be used by system tasks.
-- name: InsertUser :exec
INSERT INTO users (id, name, email, created_by)
VALUES (:id, :name, :email, :created_by);

-- --------------------------------------------------------------------------
-- CreateUser creates a new User.
-- name: CreateUser :one
INSERT INTO users (name, email, created_by)
VALUES (:name, :email, :created_by)
RETURNING id;

-- --------------------------------------------------------------------------
-- CreateTask creates a new Task.
-- name: CreateTask :one
INSERT INTO tasks (title, description, due_date, assigned_to, created, created_by)
VALUES (:title, :description, :due_date, :assigned_to, strftime('%s', 'now'), :created_by)
RETURNING id;

-- --------------------------------------------------------------------------
-- CreateSubTask creates a new Task as a sub-task of an existing Task.
-- name: CreateSubTask :one
INSERT INTO tasks (title, description, due_date, assigned_to, parent_id,
                   created, created_by)
VALUES (:title, :description, :due_date, :assigned_to, :task_id,
        strftime('%s', 'now'), :created_by)
RETURNING id;

-- --------------------------------------------------------------------------
-- CompleteTask marks a task as completed.
-- name: CompleteTask :exec
UPDATE tasks
SET completed    = strftime('%s', 'now'),
    completed_by = :completed_by,
    modified     = strftime('%s', 'now'),
    modified_by  = :modified_by
WHERE id = :task_id;

-- --------------------------------------------------------------------------
-- CompleteSubTasks marks all sub-tasks of a given task as completed.
-- name: CompleteSubTasks :exec
UPDATE tasks
SET completed    = strftime('%s', 'now'),
    completed_by = :completed_by,
    modified     = strftime('%s', 'now'),
    modified_by  = :modified_by
WHERE parent_id = :task_id;