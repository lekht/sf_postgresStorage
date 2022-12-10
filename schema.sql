--DROP DATABASE IF EXISTS tasks;

--CREATE DATABASE tasks;

DROP TABLE IF EXISTS users, labels, tasks, tasks_labels;

CREATE TABLE users (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE labels (
	id BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE tasks (
	id BIGSERIAL PRIMARY KEY,
	opened BIGINT NOT NULL DEFAULT extract(epoch from now()),
	closed BIGINT DEFAULT 0,
	author_id INTEGER REFERENCES users(id) DEFAULT 0,
	assigned_id INTEGER REFERENCES users(id) DEFAULT 0,
	title TEXT,
	content TEXT NOT NULL
);

CREATE TABLE tasks_labels (
	task_id INTEGER REFERENCES tasks(id),
	label_id INTEGER REFERENCES labels(id)
);

INSERT INTO users(id, name) values (0, 'default')