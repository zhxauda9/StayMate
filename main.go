package main

import "github.com/zhxauda9/StayMate/cmd"

func main() {
	go cmd.InitMicroservice()
	cmd.InitApp()
}
