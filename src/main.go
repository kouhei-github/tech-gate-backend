package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"net-http/myapp/crontab"
	"net-http/myapp/route"
	"net/http"
	"os"
)

func main() {
	// crontabでジョブの実行
	crontab.ToStartCron()

	// 環境変数の読み込み
	err := godotenv.Load(".env")
	// もし err がnilではないなら、"読み込み出来ませんでした"が出力されます。
	if err != nil {
		panic(err)
	}
	router := route.Router{Mutex: http.NewServeMux()}
	router.GetRouter()
	router.GetAuthRouter()
	// corsについて https://maku77.github.io/p/goruwy4/
	corsOrigin := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("ALLOW_ORIGIN")},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	////全てを許可する Access-Control-Allow-Origin: *
	//corsOrigin := cors.Default()
	handler := corsOrigin.Handler(router.Mutex)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		panic(err)
	}
}
