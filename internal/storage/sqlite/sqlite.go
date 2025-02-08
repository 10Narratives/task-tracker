package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/10Narratives/task-tracker/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

type TaskStorage struct {
	Limit uint    // Maximum number of tasks to fetch in a group.
	DB    *sql.DB // Database connection used to interact with the scheduler table.
}

// New creates a new TaskStorage instance with a given database connection and task limit.
//
// Parameters:
// - db: Pointer to an active SQL database connection.
// - limit: Maximum number of tasks to retrieve in queries.
//
// Returns:
// - TaskStorage: A struct for managing task storage operations.
func New(db *sql.DB, limit uint) TaskStorage {
	return TaskStorage{DB: db, Limit: limit}
}

// Prepare initializes the database by creating the 'scheduler' table and an index on the 'date' column if they do not exist.
//
// Returns:
// - error: An error if the table creation or index setup fails.
func (s TaskStorage) Prepare() error {
	query := `
		CREATE TABLE IF NOT EXISTS scheduler (
    	id INTEGER PRIMARY KEY AUTOINCREMENT,
    	date TEXT NOT NULL,
    	title TEXT NOT NULL,
    	comment TEXT,
    	repeat TEXT CHECK(LENGTH(repeat) <= 128)
	);

	CREATE INDEX IF NOR EXISTS idx_scheduler_date ON scheduler(date);
	`

	stmt, err := s.DB.Prepare(query)

	if err != nil {
		return fmt.Errorf("can not prepare statement: %w", err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return fmt.Errorf("can not prepare database: %w", err)
	}

	return nil
}

func (s TaskStorage) Create(ctx context.Context, date, title, comment, repeat string) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	result, err := s.DB.ExecContext(ctx, query, date, title, comment, repeat)
	if err != nil {
		return 0, fmt.Errorf("cannot insert task in database: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("cannot take last insert id: %w", err)
	}

	return lastID, nil
}

// Read retrieves a task from the scheduler database by its ID.
//
// Parameters:
// - ctx: Context for request cancellation and timeout control.
// - id: Unique identifier of the task to be retrieved.
//
// Returns:
// - models.Task: The retrieved task if found.
// - error: Returns nil if no task is found, or a wrapped error if a database operation fails.
func (s TaskStorage) Read(ctx context.Context, id int64) (models.Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	row := s.DB.QueryRow(query, id)
	task := models.Task{}

	err := row.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Task{}, nil
	}

	if err != nil {
		return models.Task{}, fmt.Errorf("cannot read task from database: %w", err)
	}

	return task, nil
}

func (s TaskStorage) queryTasks(ctx context.Context, query string, args ...interface{}) ([]models.Task, error) {
	rows, err := s.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return make([]models.Task, 0), fmt.Errorf("cannot execute query: %w", err)
	}
	defer rows.Close()

	tasks := make([]models.Task, 0)
	for rows.Next() {
		task := models.Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return make([]models.Task, 0), fmt.Errorf("cannot read row: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return make([]models.Task, 0), fmt.Errorf("cannot read task group: %w", err)
	}

	return tasks, nil
}

// ReadGroup retrieves a limited number of tasks ordered by date.
//
// Parameters:
// - ctx: Context for request cancellation and timeout control.
//
// Returns:
// - []models.Task: A slice of retrieved tasks, ordered by date.
// - error: Wrapped error if the query fails.
func (s TaskStorage) ReadGroup(ctx context.Context) ([]models.Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT (?)`
	return s.queryTasks(ctx, query, s.Limit)
}

// ReadByDate retrieves tasks that match a specific date.
//
// Parameters:
// - ctx: Context for request cancellation and timeout control.
// - date: The date to filter tasks by (must not be empty).
//
// Returns:
// - []models.Task: A slice of tasks that match the given date
// - error: Returns ErrEmptyDate if the date is empty or a wrapped error if the query fails.
func (s TaskStorage) ReadByDate(ctx context.Context, date string) ([]models.Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE date = ? LIMIT ?`
	return s.queryTasks(ctx, query, date, s.Limit)
}

// ReadByPayload retrieves tasks where the title or comment matches the given payload.
//
// Parameters:
// - ctx: Context for request cancellation and timeout control.
// - payload: The search keyword (must not be empty).
//
// Returns:
// - []models.Task: A slice of matching tasks, ordered by date.
// - error: Returns ErrEmptyPayload if the payload is empty or a wrapped error if the query fails.
func (s TaskStorage) ReadByPayload(ctx context.Context, payload string) ([]models.Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT ?`
	return s.queryTasks(ctx, query, payload, payload, s.Limit)
}

// Update modifies an existing task in the scheduler database.
//
// Parameters:
// - ctx: Context for request cancellation and timeout control.
// - t: Pointer to the Task model containing updated values (must not be nil).
//
// Returns:
// - error: Returns ErrNilTaskUpdate if t is nil or a wrapped error if the update fails.
func (s TaskStorage) Update(ctx context.Context, t *models.Task) error {
	if t == nil {
		return fmt.Errorf("cannot update task using nil pointer")
	}

	query := `
		UPDATE scheduler
		SET date = ?, title = ?, comment = ?, repeat = ?
		WHERE id = ?
	`

	_, err := s.DB.ExecContext(ctx, query, t.Date, t.Title, t.Comment, t.Repeat, t.ID)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// Delete removes a task from the scheduler database by its ID.
//
// Parameters:
// - ctx: Context for request cancellation and timeout control.
// - id: Unique identifier of the task to be deleted.
//
// Returns:
// - error: Wrapped error if the deletion fails.
func (s TaskStorage) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM scheduler
		WHERE id = ?
	`
	_, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}
	return nil
}
