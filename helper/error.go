package helper

import "fmt"

// Recovery function for recovering from panic
func ErrorRecovery() {
	if r := recover(); r != nil {
		fmt.Println("recovered from ", r)
	}
}

// Check function for checking error
func ErrorCheck(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
