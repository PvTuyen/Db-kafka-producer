package schema

import (
	"db-producer/common"
	"db-producer/logger"
	"errors"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"strconv"
)

func (schema *Schema) getSchemaInfo() error {
	path := "config/input/"
	file, err := excelize.OpenFile(path + "schema.xlsx")
	if err != nil {
		logger.LogError(common.GET_SCHEMA_FAIL, err.Error())

		return err
	}
	sheetNames := file.GetSheetList()
	for _, v := range sheetNames {
		err := schema.getValueSheet(v, file)
		if err != nil {
			logger.LogError("", err.Error())

			return err
		}
	}

	return nil
}

func (schema *Schema) getValueSheet(sheetName string, file *excelize.File) error {
	if sheetName == "info" {
		return nil
	}
	tableName, err := file.GetCellValue(sheetName, common.TABLE_NAME_INDEX)
	backup, err := file.GetCellValue(sheetName, common.TABLE_BACKUP_INDEX)
	if backup != common.BACKUP_VALUE {
		return nil
	}
	if err != nil {

		return err
	}
	if tableName == "" {

		return errors.New(common.TABLE_NAME_NULL)
	}
	indexCol := common.NAME_COLUMNS_INDEX_START
	columns := []string{}
	for {
		nameCol, err := file.GetCellValue(sheetName, common.NAME_COLUMNS_INDEX + strconv.Itoa(indexCol))
		if err != nil {

			return err
		}
		if nameCol == "" {

			break
			//return errors.New(common.NAME_CLOLUMN_NULL)
		}
		columns = append(columns, nameCol)
		indexCol++
	}
	schema.MapInfo[tableName] = columns

	return nil
}
