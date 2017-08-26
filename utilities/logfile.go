package utilities

import (
	"fmt"
	"time"
	"os"
)

type LogFile struct {

}

func (lf *LogFile) CurrentFile()(*os.File){
	fileName := fmt.Sprintf("log.%s.txt", time.Now().Format("2006_01_02"))
	file,_ := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND|os.O_CREATE,0666)
	return file
}

func (lf *LogFile) Write(p []byte) (n int, err error) {
	file := lf.CurrentFile()
	return file.Write(p)
}
