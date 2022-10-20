package common

const (
	ERROR = "ERROR"
	// init system
	GET_ENV_FAIL = "get value in file env fail "
	CONNECT_SQL_FAIL = "connect sql fail"
	QUERY_FAIL = "query fail"
	ROWS_CONVERT_FAIL = "convert record fail"
	LEN_COLUMNS_NOT_EQUALS_LEN_ROWS = "columns length not equal length rows"
	DRIVERR_SQL_NOT_SUPPORT = "driver not support"

	// kafka
	SEND_MESSAGE_FAIL = "send message fail"

	// schema
	GET_SCHEMA_FAIL = "get schema fail"
	INIT_SCHEMA_FAIL = "init schema fail"
	// csv
	TABLE_NAME_NULL = "table name is null"
	NAME_CLOLUMN_NULL = "columns name is null"
	// sql
	DELETE_SNAPSHOT_FAIL = "delete snapshot record fail"
)
