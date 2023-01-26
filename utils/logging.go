package utils

import (
	"io"
	"log"
	"os"
)

func Logging(filename string) {
	//    　　　　　　　　　　読み書き作成、追記 開く際のモード|ファイルが存在しない場合は新たに作成|ファイルが既に存在する場合末尾に書き足す
	logfile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	// 　　　　　　　　　　　　　　　ログの書き込み先を標準出力とログファイルに指定
	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	// フォーマット指定　「logとして何を出力するか」が書かれています。今回は日付・時刻・ファイル名をそれぞれ詳細表示しています。
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	// ログの出力先
	log.SetOutput(multiLogFile)
}
