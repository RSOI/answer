SET SYNCHRONOUS_COMMIT = 'off';
CREATE EXTENSION IF NOT EXISTS CITEXT;
CREATE SCHEMA IF NOT EXISTS answer;

DROP TABLE IF EXISTS answer.answer;
DROP TABLE IF EXISTS answer.services;

CREATE TABLE answer.answer (
	id SERIAL PRIMARY KEY,
	question_id INTEGER NOT NULL,
	content CITEXT NULL,
	author_id INTEGER NOT NULL,
	author_nickname CITEXT NOT NULL,
	is_best BOOLEAN DEFAULT FALSE,
	created TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS question_id_index ON answer.answer (question_id);
CREATE INDEX IF NOT EXISTS author_id_index ON answer.answer (author_id);
CREATE INDEX IF NOT EXISTS is_best_index ON answer.answer (is_best);
CREATE INDEX IF NOT EXISTS question_id__is_best_index ON answer.answer (question_id, is_best);

CREATE TABLE answer.services (
	id SERIAL PRIMARY KEY,
	request CITEXT NOT NULL,
	request_time TIMESTAMPTZ DEFAULT NOW(),
	response_status INTEGER NOT NULL,
	response_error_text CITEXT NULL
);
