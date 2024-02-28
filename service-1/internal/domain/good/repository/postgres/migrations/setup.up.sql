CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
);

INSERT INTO projects(name, created_at)
VALUES('test', NOW());

CREATE TABLE IF NOT EXISTS goods (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    priority INTEGER NOT NULL,
    removed BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

INSERT INTO goods(project_id, name, description, priority, removed, created_at)
VALUES(1, 'toy', 'an expensive toy', 1, 'f', NOW());
