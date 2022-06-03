# upv-server

[![GPLv3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

```bash
# tidy
go mod tidy

# run dev
env ENV=development go run main.go

# run prod
env ENV=production go run main.go

# mysql dump
mysqldump --databases upv -u root -p > upv.sql
scp -rp root@upv.life:/~/upv.sql ./

# scp download files
scp -rp root@upv.life:/path/filename /path #将远程文件从服务器下载到本地

# scp upload files
scp -rp /Users/g/code/web/dist/admin/* root@upv.life:/var/www/admin
scp -rp /Users/g/code/web/dist/index/* root@upv.life:/var/www/html


# docker up (todo)
docker-compose up -d -f docker-compose.yml --env-file .env
docker-compose --env-file .env config
```

# ref:

- https://www.digitalocean.com/community/tutorials/how-to-deploy-a-go-web-application-using-nginx-on-ubuntu-18-04

- https://www.digitalocean.com/community/tutorials/how-to-install-go-on-ubuntu-20-04
