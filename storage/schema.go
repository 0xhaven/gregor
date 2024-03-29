package storage

var schema = []string{

	`DROP TABLE IF EXISTS dismissals_by_time`,
	`DROP TABLE IF EXISTS dismissals_by_id`,
	`DROP TABLE IF EXISTS reminders`,
	`DROP TABLE IF EXISTS items`,
	`DROP TABLE IF EXISTS messages`,

	`CREATE TABLE messages (
		uid   CHAR(16) NOT NULL,
		msgid CHAR(16) NOT NULL,
		ctime DATETIME(6) NOT NULL,
		devid CHAR(16),
		mtype INTEGER UNSIGNED NOT NULL, -- "specify for 'Update' or 'Sync' types",
		PRIMARY KEY(uid, msgid)
	)`,

	`CREATE TABLE items (
		uid   CHAR(16) NOT NULL,
		msgid CHAR(16) NOT NULL,
		category VARCHAR(128) NOT NULL,
		dtime DATETIME(6),
		body BLOB,
		FOREIGN KEY(uid, msgid) REFERENCES messages (uid, msgid),
		PRIMARY KEY(uid, msgid)
	)`,

	`CREATE INDEX user_order ON items (uid, category)`,

	`CREATE INDEX cleanup_order ON items (uid, dtime)`,

	`CREATE TABLE reminders (
		uid   CHAR(16) NOT NULL,
		msgid CHAR(16) NOT NULL,
		ntime DATETIME(6) NOT NULL,
		PRIMARY KEY(uid, msgid, ntime)
	)`,

	`CREATE TABLE dismissals_by_id (
		uid   CHAR(16) NOT NULL,
		msgid CHAR(16) NOT NULL,
		dmsgid CHAR(16) NOT NULL, -- "the message IDs to dismiss",
		FOREIGN KEY(uid, msgid) REFERENCES messages (uid, msgid),
		PRIMARY KEY(uid, msgid, dmsgid)
	)`,

	`CREATE TABLE dismissals_by_time (
		uid   CHAR(16) NOT NULL,
		msgid CHAR(16) NOT NULL,
		category VARCHAR(128) NOT NULL,
		dtime DATETIME(6) NOT NULL, -- "throw out matching events before dtime",
		FOREIGN KEY(uid, msgid) REFERENCES messages (uid, msgid),
		PRIMARY KEY(uid, msgid, category, dtime)
	)`,
}

func Schema(engine string) []string {
	return schema
}
