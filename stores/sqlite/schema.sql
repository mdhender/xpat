--  Copyright (c) 2024 Michael D Henderson. All rights reserved.

-- --------------------------------------------------------------------------
-- this file defines the schema for Sqlite3 data store.

PRAGMA foreign_keys = ON;

-- --------------------------------------------------------------------------
-- Create the users table
CREATE TABLE users
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,                -- Unique identifier for each user
    name       TEXT    NOT NULL,                                 -- User's name
    email      TEXT    NOT NULL UNIQUE,                          -- User's email address
    --
    created    INTEGER NOT NULL DEFAULT (strftime('%s', 'now')), -- Creation timestamp as Unix epoch
    created_by INTEGER NOT NULL                                  -- ID of the user who created the task
);

-- --------------------------------------------------------------------------
-- Create the tasks table
CREATE TABLE tasks
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,                -- Unique identifier for each task
    title        TEXT    NOT NULL,                                 -- Title of the task
    description  TEXT,                                             -- Detailed description of the task
    due_date     DATETIME,                                         -- Unix timestamp, stored UTC
    --
    assigned_to  INTEGER,                                          -- ID of the user assigned to the task
    parent_id    INTEGER,                                          -- ID of the parent task for sub-tasks
    --
    created      INTEGER NOT NULL DEFAULT (strftime('%s', 'now')), -- Creation timestamp as Unix epoch
    created_by   INTEGER NOT NULL,                                 -- ID of the user who created the task
    --
    completed    INTEGER,                                          -- Completion timestamp as Unix epoch
    completed_by INTEGER,                                          -- ID of the user who completed the task
    --
    modified     INTEGER NOT NULL DEFAULT (strftime('%s', 'now')), -- Last modified timestamp as Unix epoch
    modified_by  INTEGER,                                          -- ID of the user who modified the task
    --
    FOREIGN KEY (assigned_to) REFERENCES users (id),               -- Foreign key to the users table
    FOREIGN KEY (completed_by) REFERENCES users (id),              -- Foreign key to the users table
    FOREIGN KEY (created_by) REFERENCES users (id),                -- Foreign key to the users table
    FOREIGN KEY (modified_by) REFERENCES users (id),               -- Foreign key to the users table
    FOREIGN KEY (parent_id) REFERENCES tasks (id)                  -- Self-referencing foreign key
        ON DELETE CASCADE                                          -- Automatically delete children when parent is deleted
);

-- --------------------------------------------------------------------------
-- Create trigger to complete sub-tasks when parent task is completed.
CREATE TRIGGER update_children_on_completion
    AFTER UPDATE OF completed
    ON tasks
    FOR EACH ROW
    WHEN NEW.completed IS NOT NULL -- Trigger only if the parent task is marked as completed
BEGIN
    UPDATE tasks
    SET completed    = NEW.completed,
        completed_by = NEW.completed_by
    WHERE parent_id = NEW.id;
END;

-- --------------------------------------------------------------------------
-- Create trigger to update modified timestamp.
CREATE TRIGGER update_modified
    AFTER UPDATE
    ON tasks
    FOR EACH ROW
BEGIN
    UPDATE tasks
    SET modified = strftime('%s', 'now')
    WHERE id = OLD.id;
END;

-- --------------------------------------------------------------------------
-- Create indexes for common queries
CREATE INDEX idx_todos_parent_id ON tasks (parent_id);
CREATE INDEX idx_todos_assigned_to ON tasks (assigned_to);
CREATE INDEX idx_todos_completed ON tasks (completed);
