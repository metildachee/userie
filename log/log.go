package log

import "fmt"

func Print(msg, level string) {
	fmt.Println(level, msg)
}
