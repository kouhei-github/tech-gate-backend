package utils

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

func WriteLogFile(text string) {
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	nowJST := time.Now().In(loc)

	timestamp := nowJST.Format(time.RFC3339)
	fileName := strconv.Itoa(nowJST.Year()) + "-" + strconv.Itoa(int(nowJST.Month())) + "-" + strconv.Itoa(nowJST.Day()) + ".log"
	//ロガーの生成
	var logger = logrus.New()
	//ファイル取得
	//ファイルは無ければ生成(os.O_CREATE)、書き込み(os.O_WRONLY)、追記モード(os.O_APPEND)、権限は0666
	file, err := os.OpenFile("log/"+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.Out = file
	} else {
		panic(err)
	}

	logMessage := logFormat(text, timestamp)
	_, err = file.WriteString(logMessage)
	if err != nil {
		panic(err)
	}
}

func logFormat(logMsg string, currentTime string) string {
	msg := currentTime + "\n" + logMsg + "\n"
	return msg
}
