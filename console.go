package cuppago

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"runtime"
	"time"
)

func Log(values ...interface{}) {
	log.SetOutput(os.Stdout)
	for i := 0; i < len(values); i++ {
		log.Print(values[i])
	}
}

func LogFile(values ...interface{}) {
	path := GetRootPath() + "/log/"
	CreateFolder(path)
	fileName := path + "log_" + time.Now().String()[0:10] + ".log"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err)
	}
	log.SetOutput(file)
	for i := 0; i < len(values); i++ {
		log.Print(values[i])
		fmt.Fprintln(os.Stdout, values[i])
	}
}

func Error(values ...interface{}) {
	path := GetRootPath() + "/log/"
	CreateFolder(path)
	fileName := path + "error_" + time.Now().String()[0:10] + ".log"
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Print(err)
	}
	log.SetOutput(file)
	for i := 0; i < len(values); i++ {
		log.Println(values[i])
		fmt.Fprintln(os.Stderr, values[i])
	}
}

func Type(value interface{}) string {
	var result = reflect.TypeOf(value).String()
	return result
}

func Wait(values ...interface{}) {
	Log(values)
	fmt.Scanf("$s")
}

func Goroutines(intervalTime time.Duration) {
	if intervalTime == 0 {
		intervalTime = 5
	}
	Log("GOROUTINES [" + String(runtime.NumGoroutine()) + "]")
	time.Sleep(intervalTime * time.Second)
	Goroutines(intervalTime)
}
