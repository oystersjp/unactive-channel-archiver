package main

import (
	"os"
	"strconv"
	"archiver/archiver"
)

func main() {
	token := os.Getenv("SLACK_TOKEN")
	expireDate, _ := strconv.Atoi(os.Getenv("EXPIRE_DATE"))
	summaryCh := os.Getenv("SUMMARY_CHANNEL")

	service := archiver.NewArchiver(token, []string{}, expireDate)
	service.Exec(summaryCh)
}
