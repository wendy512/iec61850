package main

import (
	"fmt"
	"github.com/wendy512/iec61850"
	"github.com/wendy512/iec61850/client"
	"time"
)

func main() {
	settings := iec61850.NewSettings()
	c, err := client.NewClient(settings)
	if err != nil {
		panic(err)
	}

	value, err := c.Read("CNPDMONT/SFAN1.FJStu.stVal", iec61850.ST)
	if err != nil {
		panic(err)
	}

	fmt.Printf("read value %v\n", value)

	if err = c.Write("CNPDMONT/SFAN1.FJOp.setVal", iec61850.SE, true); err != nil {
		panic(err)
	}

	time.Sleep(time.Millisecond * 200)
	defer c.Close()
}
