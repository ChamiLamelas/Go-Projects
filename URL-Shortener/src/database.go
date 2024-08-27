package main

const SQL_DRIVER = "sqlite3"

const DATABASE_FOLDER = "../data/"

const DATABASE_FILE = DATABASE_FOLDER + "database.db"

var QUERY_CREATE_TABLE = `
CREATE TABLE IF NOT EXISTS aliases (
	URL TEXT UNIQUE NOT NULL,
	Alias TEXT PRIMARY KEY,
	Expansions INT,
	Automatic BOOL
);
`

const QUERY_GET_NEXT_ALIAS = `
SELECT MAX(CAST(Alias AS INTEGER))
FROM aliases 
WHERE Automatic
`

const QUERY_MAKE_MAPPING_TEMPLATE = `
INSERT INTO aliases (URL, Alias, Expansions, Automatic) 
VALUES (?, ?, ?, ?)
`

const QUERY_GET_ALIAS_BY_URL_TEMPLATE = `
SELECT Alias
FROM aliases 
WHERE URL = ?
`

const QUERY_GET_URL_BY_ALIAS_TEMPLATE = `
SELECT URL
FROM aliases
WHERE Alias = ?
`

const QUERY_UPDATE_ANALYTICS_BY_ALIAS_TEMPLATE = `
UPDATE aliases 
SET Expansions = Expansions + 1
WHERE Alias = ?
`

const QUERY_GET_ANALYTICS_BY_ALIAS_TEMPLATE = `
SELECT URL, Expansions
FROM aliases
WHERE Alias = ?
`

const DUPLICATE_URL_VIOLATION = "UNIQUE constraint failed: aliases.URL"

const DUPLICATE_ALIAS_VIOLATION = "UNIQUE constraint failed: aliases.Alias"
