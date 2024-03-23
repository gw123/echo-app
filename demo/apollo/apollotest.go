package main

import (
	"fmt"

	"github.com/philchia/agollo/v4"
)

func main() {
	c := agollo.NewClient(&agollo.Conf{
		AppID:           "SampleApp",
		Cluster:         "default",
		NameSpaceNames:  []string{"application.properties"},
		MetaAddr:        "http://127.0.0.1:8080",
		AccesskeySecret: "061825b59ff44624b40d26caa5380367",
	})

	err := c.Start()
	if err != nil {
		panic(err)
	}
	arr := c.GetAllKeys(agollo.WithNamespace("application"))
	fmt.Println(arr)
}
