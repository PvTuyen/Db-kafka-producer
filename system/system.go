package system

import (
	"db-producer/common"
	"db-producer/kafkaProducer"
	"db-producer/logger"
	"db-producer/rdb"
	"db-producer/schema"
)

type System struct {
	Rdb    *rdb.Rdb
	Pro    *kafkaProducer.Producer
	Schema *schema.Schema
}

func InitSystem() *System {
	logger.LogInfo(common.START_SERVER, "")
	sys := System{
		Schema: schema.InitSchema(),
		Rdb:    rdb.InitRdbInfo(),
		Pro:    kafkaProducer.InitKafka(),
	}

	sys.Pro.ProcessSendMsg(sys.Schema, 1000)
	return &sys
}
