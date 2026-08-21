package uuidx

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

var (
	snowflakeNode *snowflake.Node
	initErr       error
)

func init() {
	snowflakeNode, initErr = snowflake.NewNode(1)
	if initErr != nil {
		log.Fatalf("初始化雪花ID生成器失败: %v", initErr)
	}
}

func GetSnowflakeId() snowflake.ID {
	if initErr != nil {
		log.Println("生成雪花id失败")
		return 0
	}

	return snowflakeNode.Generate()
}
