package crontab

import (
	"github.com/bamzi/jobrunner"
	"net-http/myapp/crontab/cron_crawler"
)

func ToStartCron() {
	jobrunner.Start()
	jobrunner.Schedule("@every 5s", cron_crawler.Zenn{Url: "https://zenn.dev/api/articles"})
	//
	jobrunner.Schedule("@every 5s", cron_crawler.Qiita{Url: "http://qiita.com/api/v2/items"})
	//jobrunner.Schedule("TZ=Asia/Tokyo 30 04 * * * *", cron_crawler.Zenn{Url: "https://zenn.dev/api/articles"})

}
