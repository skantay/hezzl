CREATE TABLE IF NOT EXISTS goods (
    ID Int32,
    ProjectID Int32,
    Name String,
    Description String,
    Priority Int32,
    Removed UInt8,
    CreatedAt DateTime('UTC'),
    INDEX idx_id ID TYPE minmax GRANULARITY 1,
    INDEX idx_project_id ProjectID TYPE minmax GRANULARITY 1,
    INDEX idx_name Name TYPE set(0) GRANULARITY 1
) ENGINE = MergeTree()
ORDER BY ID;
