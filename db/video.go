package db

import (
	"database/sql"
	"time"

	"github.com/cliclitv/go-clicli/def"
)

func AddVideos(videos []*def.Video, uid int, pid int) ([]*def.Video, error) {
	var v []*def.Video
	for _, video := range videos {
		if resp, err := AddVideo(video.Oid, video.Title, video.Content, pid, uid); err != nil {
			v = append(v, resp)
			return nil, err
		}
	}
	return v, nil
}

func AddVideo(oid int, title string, content string, pid int, uid int) (*def.Video, error) {
	t := time.Now()
	ctime := t.Format("2006-01-02 15:04")
	stmtIns, err := db.Prepare("INSERT INTO videos (oid,title,content,create_time,pid,uid) VALUES (?,?,?,?,?,?)")
	if err != nil {
		return nil, err
	}
	_, err = stmtIns.Exec(oid, title, content, ctime, pid, uid)
	if err != nil {
		return nil, err
	}
	defer stmtIns.Close()

	res := &def.Video{Oid: oid, Title: title, Content: content, Uid: uid, Pid: pid}
	defer stmtIns.Close()
	return res, err
}

func GetVideos(pid int, page int, pageSize int) (*[]def.Video, error) {
	start := pageSize * (page - 1)
	videos := []def.Video{}

	err := db.Select(&videos, "SELECT * FROM videos WHERE pid=?  ORDER BY oid DESC LIMIT ?,?", pid, start, pageSize)

	print(err)
	return &videos, nil

}

func GetVideo(id int) (*def.Video, error) {
	stmtOut, err := db.Prepare(`SELECT videos.id,videos.oid,videos.title,videos.content,videos.create_time,posts.id,posts.title,users.id,users.name,users.qq FROM (videos INNER JOIN posts ON videos.pid=posts.id)
INNER JOIN users ON videos.uid = users.id WHERE videos.id = ?`)
	if err != nil {
		return nil, err
	}
	var vid, uid, oid, pid int
	var title, content, ctime, uname, uqq, ptitle string

	err = stmtOut.QueryRow(id).Scan(&vid, &oid, &title, &content, &ctime, &pid, &ptitle, &uid, &uname, &uqq)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	defer stmtOut.Close()

	res := &def.Video{Id: vid, Oid: oid, Title: title, Content: content, Pid: pid, Uid: uid}

	return res, nil
}

func UpdateVideo(id int, oid int, title string, content string, pid int, uid int) (*def.Video, error) {
	stmtIns, err := db.Prepare("UPDATE videos SET oid=?,title=?,content=?,pid=?,uid=? WHERE id =?")
	if err != nil {
		return nil, err

	}
	_, err = stmtIns.Exec(&oid, &title, &content, &pid, &uid, &id)
	if err != nil {
		return nil, err
	}
	defer stmtIns.Close()

	res := &def.Video{Id: id, Oid: oid, Title: title, Content: content, Pid: pid, Uid: uid}
	return res, err
}

func DeleteVideo(id int, pid int) error {
	stmtDel, err := db.Prepare("DELETE FROM videos WHERE id=? OR pid=?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(id, pid)
	if err != nil {
		return err
	}
	defer stmtDel.Close()

	return nil

}
