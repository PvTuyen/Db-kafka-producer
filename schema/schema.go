package schema

import (
	"db-producer/common"
	"db-producer/logger"
)

type Schema struct {
	// key= name table, value= columns
	MapInfo map[string][]string
}

func InitSchema() *Schema {
	infoSchema := Schema{
		MapInfo :make(map[string][]string),
	}
	err := infoSchema.getSchemaInfo()
	if err != nil {
		logger.LogError(common.INIT_SCHEMA_FAIL, err.Error())
	}
	return &infoSchema
}