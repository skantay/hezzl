CREATE TABLE IF NOT EXISTS projects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

INSERT INTO projects(name, created_at)
VALUES('Первая запись', NOW());

CREATE TABLE IF NOT EXISTS goods (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    priority INTEGER NOT NULL,
    removed BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id)
);

CREATE INDEX idx_goods_id_projectid_name ON goods(id, project_id, name);