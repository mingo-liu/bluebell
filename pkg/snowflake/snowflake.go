package  snowflake

import (
	// "fmt"
	"time"
	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)	// 解析开始时间
	if err != nil {
		return 
	}
	snowflake.Epoch = st.UnixNano() / 1000000
	node, err = snowflake.NewNode(machineID)
	return 
}

func GenIDInt64() int64 {
	return node.Generate().Int64()
}

func GenIDString() string {
	return node.Generate().String()
}

func GenIDBase2() string {
	return node.Generate().Base2()
}

// func main() {
// 	if err := Init("2025-01-03", 1); err != nil {
// 		fmt.Printf("init failed, err: %v\n", err)
// 		return
// 	}
	
// 	id := GenIDInt64()
// 	idStr := GenIDString()
// 	idBase2 := GenIDBase2()

// 	fmt.Printf("Int64 ID: %d\n", id)
// 	fmt.Printf("String ID: %s\n", idStr)
// 	fmt.Printf("Base2 ID: %s %d\n", idBase2, len(idBase2))
// }