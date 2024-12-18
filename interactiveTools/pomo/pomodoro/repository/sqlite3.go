package repository

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"haonguyen.tech/interactiveTools/pomo/pomodoro"
)

const createTableInterval string = `CREATE TABLE IF NOT EXISTS "interval" (
	"id" INTEGER,
	"start_time" DATETIME NOT NULL,
	"planned_duration" INTEGER DEFAULT 0,
	"actual_duration" INTEGER DEFAULT 0,
	"category" TEXT NOT NULL,
	"state" INTEGER DEFAULT 1,
	PRIMARY KEY ("id")
	);`

type dbRepo struct {
	db           *sql.DB
	sync.RWMutex // prevent concurrent access to db
}

func NewSqlite3Repo(dbFile string) (*dbRepo, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(createTableInterval); err != nil {
		return nil, err
	}

	return &dbRepo{
		db: db,
	}, nil
}

func (r *dbRepo) Create(i pomodoro.Interval) (int64, error) {
	r.Lock()
	defer r.Unlock()
	insertStm, err := r.db.Prepare("INSERT INTO interval VALUES(NULL,?,?,?,?,?)")
	if err != nil {
		return 0, err
	}
	defer insertStm.Close()

	res, err := insertStm.Exec(i.StartTime, i.PlannedDuration, i.ActualDuration, i.Category, i.State)
	if err != nil {
		return 0, err
	}
	// INSERT results
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *dbRepo) Update(i pomodoro.Interval) error {
	r.Lock()
	defer r.Unlock()
	updateStm, err := r.db.Prepare("UPDATE interval SET start_time=?, actual_duration=?, state=? WHERE id=?")
	if err != nil {
		return err
	}

	defer updateStm.Close()
	res, err := updateStm.Exec(i.StartTime, i.ActualDuration, i.State, i.ID)
	if err != nil {
		return err
	}

	if _, err := res.RowsAffected(); err != nil {
		return err
	}
	return nil
}

func (r *dbRepo) ByID(id int64) (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()
	row := r.db.QueryRow("SELECT * FROM interval WHERE id=?", id)
	i := pomodoro.Interval{}
	err := row.Scan(&i.ID, &i.StartTime, &i.PlannedDuration, &i.ActualDuration, &i.Category, &i.State)
	if err != nil {
		return i, err
	}
	return i, nil
}

func (r *dbRepo) Last() (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	last := pomodoro.Interval{}
	// Query the latest , we sort by id for now
	err := r.db.QueryRow("SELECT * FROM interval ORDER BY id desc LIMIT 1").Scan(&last.ID, &last.StartTime, &last.PlannedDuration, &last.ActualDuration, &last.Category, &last.State)
	if err == sql.ErrNoRows {
		return last, pomodoro.ErrNoIntervals
	}

	if err != nil {
		return last, err
	}

	return last, nil
}

func (r *dbRepo) Breaks(n int) ([]pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()
	stmt := `SELECT * FROM interval WHERE category LIKE '%Break' ORDER BY id DESC LIMIT ?`

	rows, err := r.db.Query(stmt, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := []pomodoro.Interval{}

	for rows.Next() {
		i := pomodoro.Interval{}
		err := rows.Scan(&i.ID, &i.StartTime, &i.PlannedDuration, &i.ActualDuration, &i.Category, &i.State)
		if err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return data, nil
}
