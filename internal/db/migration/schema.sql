CREATE TABLE IF NOT EXISTS tasksmaster (
    task_id INTEGER PRIMARY KEY AUTOINCREMENT,     -- unique incremental ID
    title TEXT NOT NULL,                           -- task title
    description TEXT,                              -- detailed text
    priority INTEGER CHECK(priority BETWEEN 1 AND 3) DEFAULT 2, -- 1=high,2=med,3=low
    status INTEGER CHECK(status BETWEEN 1 AND 4) DEFAULT 1,     -- 1=pending,2=wip,3=done,4=archived
    parent_task_id INTEGER,                        -- reference to parent task
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- creation time
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- last updated time
    deadline_at DATETIME,                          -- deadline (optional)
    -- children relationship handled separately (via self-reference or mapping table)
    FOREIGN KEY (parent_task_id) REFERENCES tasksmaster(task_id)
);

-- optional helper table for multiple child relationships
CREATE TABLE IF NOT EXISTS task_children (
    parent_id INTEGER NOT NULL,
    child_id INTEGER NOT NULL,
    PRIMARY KEY (parent_id, child_id),
    FOREIGN KEY (parent_id) REFERENCES tasksmaster(task_id),
    FOREIGN KEY (child_id) REFERENCES tasksmaster(task_id)
);
