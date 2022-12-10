package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"

	"reflect"
)

const (
	UpdateClosedSql = `UPDATE tasks SET closed = extract(epoch from now()) WHERE id = $1`
	UpdateTitleSql  = `UPDATE tasks SET title = $1 WHERE id = $2`
	UpdateContent   = `UPDATE tasks SET content = $1 WHERE id = $2`
)

type Storage struct {
	db *pgxpool.Pool
}

func New(connstr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), connstr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

type Task struct {
	ID         int
	Opened     int64
	Closed     int64
	AuthorID   int
	AssignedID int
	Title      string
	Content    string
}

// Создание новой задачи
func (s *Storage) NewTask(t Task) (int, error) {
	var id int
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks(title, content)
		VALUES ($1,$2) RETURNING id;

		`,
		t.Title,
		t.Content,
	).Scan(&id)
	return id, err
}

// Получение всех задач
func (s *Storage) Tasks() ([]Task, error) {
	var tasks []Task
	rows, err := s.db.Query(context.Background(), `
		SELECT *
		FROM tasks
		ORDER BY id;
	`,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (s *Storage) TasksByAuthor(authorID int) ([]Task, error) {
	var tasks []Task
	rows, err := s.db.Query(context.Background(), `
		SELECT *
		FROM tasks
		WHERE author_id = $1
		ORDER BY id;
	`,
		authorID,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (s *Storage) TasksByLable(lableID int) ([]Task, error) {
	var tasks []Task
	rows, err := s.db.Query(context.Background(), `
		SELECT *
		FROM tasks
		WHERE author_id = $1
		ORDER BY id;
	`,
		lableID,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t Task
		err = rows.Scan(
			&t.ID,
			&t.Opened,
			&t.Closed,
			&t.AuthorID,
			&t.AssignedID,
			&t.Title,
			&t.Content,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (s *Storage) DeleteTask(taskID int) (bool, error) {
	_, err := s.db.Exec(context.Background(), `
	DELETE FROM tasks 
	WHERE id = $1;
	`,
		taskID)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Создает транзакцию и обновляет каждое поле отдельным запросом. Т.к у нас бэкенд, то в контексте задачи подразумевается,
// что приходит, например, json. Из него создается структура.
// Обновляются только 3 поля, остальные отрабатываются по default
func (s *Storage) UpdateTask(t Task) error {
	tx, err := s.db.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	e := reflect.ValueOf(&t).Elem()
	// iterate through all the fields
	for i := 0; i < e.NumField(); i++ {
		f := e.Type().Field(i).Name
		switch f {
		case "Closed":
			err = s.updateClosed(t)
			if err != nil {
				return err
			}
		case "Title":
			err = s.updateTitle(t)
			if err != nil {
				return err
			}
		case "Content":
			err = s.updateContent(t)
			if err != nil {
				return err
			}
		default:
			continue
		}
	}

	tx.Commit(context.Background())
	return nil
}

func (s *Storage) updateClosed(t Task) error {
	_, err := s.db.Exec(context.Background(), UpdateClosedSql, t.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) updateTitle(t Task) error {
	_, err := s.db.Exec(context.Background(), UpdateTitleSql, t.Title, t.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) updateContent(t Task) error {
	_, err := s.db.Exec(context.Background(), UpdateContent, t.Content, t.ID)
	if err != nil {
		return err
	}
	return nil
}
