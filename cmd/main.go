package main

import (
	"fmt"
	"github.com/wendy512/iec61850"
	"time"
)

func main() {
	settings := iec61850.NewSettings()
	c, err := iec61850.NewClient(settings)
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

	mmsValues, err := c.ReadDataSet("CNPDMONT/LLN0.dsState")
	if err != nil {
		panic(err)
	}
	for _, mmsValue := range mmsValues {
		fmt.Printf("dataset value %v\n", mmsValue.Type)
	}
	time.Sleep(time.Millisecond * 200)
	defer c.Close()
}
