package main

import (
	cont "github.com/expproletariy/iot-edge-system-process/container"
	"os"
)

func main() {
	proc := cont.NewIotEdgeSystemProcess()
	proc.SetConfigPath("./config/config.yml")
	proc.Run(os.Args)
}
