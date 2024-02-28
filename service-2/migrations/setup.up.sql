CREATE TABLE IF NOT EXISTS goods (
    ID Int32,
    ProjectID Int32,
    Name String,
    Description String,
    Priority Int32,
    Removed UInt8,
    CreatedAt Timestamp
) ENGINE = MergeTree()
ORDER BY ID;
