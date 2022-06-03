package util

import (
	"os"
	"os/exec"
	"time"

	"github.com/joho/godotenv"
)

func BackupDB() {
	godotenv.Load(".env")
	user := os.Getenv("MYSQL_USERNAME")
	dbName := os.Getenv("MYSQL_DATABASE")
	pwd := os.Getenv("MYSQL_PASSWORD")

	backupDir := "backup"
	backupFile := dbName + "-" + time.Now().Format("2006-01-02-15-04-05") + ".sql"
	backupPath := backupDir + "/" + backupFile
	backupCmd := "mysqldump -u " + user + " -p" + pwd + " " + dbName + " > " + backupPath
	exec.Command("mkdir", backupDir).Run()
	exec.Command("/bin/bash", "-c", backupCmd).Run()
}
