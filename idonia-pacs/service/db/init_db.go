package db

var INIT_DB = []string{
	"CREATE TABLE IF NOT EXISTS task (uuid TEXT,created_date DATETIME,auth TEXT,steps TEXT, step INTEGER,status TEXT,error TEXT,study_uid TEXT,container_id INTEGER,study_id TEXT,is_reported BOOL,internal_data TEXT)",
	"CREATE TABLE IF NOT EXISTS session (uuid TEXT, expire_date DATETIME, idonia_account_id INTEGER, idonia_token TEXT, idonia_api_key TEXT, idonia_api_secret TEXT)",
	"CREATE TABLE IF NOT EXISTS study (uuid TEXT, created_date DATETIME, study_uid TEXT, study_id TEXT, container_id INTEGER, dicom_files TEXT, is_reported BOOL, report_id INTEGER, magic_link TEXT)",
	"CREATE TABLE IF NOT EXISTS task_additional_field (task_uuid TEXT, key TEXT, value TEXT)",
	"CREATE TABLE IF NOT EXISTS task_dicom(uuid TEXT,task_uuid TEXT,series_uid TEXT,study_uid TEXT,container_id INTEGER,sop_uid TEXT,status TEXT,error TEXT,is_pending BOOL)",
	"CREATE TABLE IF NOT EXISTS info (key TEXT, value TEXT)",
	"INSERT INTO info (key, value) VALUES ('version', '1.0.0')",
}

var MIGRATIONS_DB = map[string][]string{}
