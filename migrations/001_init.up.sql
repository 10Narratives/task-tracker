CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL,
    title TEXT NOT NULL,
    comment TEXT,
    repeat TEXT CHECK(LENGTH(repeat) <= 128)
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);