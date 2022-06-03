package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/cliclitv/go-clicli/def"
	"github.com/cliclitv/go-clicli/util"
)

func GetPostById(id int) (*def.Post, error) {
	post := &def.Post{}
	err := db.Get(post, `SELECT p.*,
		u.nickname as creator_nickname,
		u.name as creator_name,
		u.avatar as creator_avatar,
		u.id as creator_id,
		u.bio as creator_bio,
		u.qq as creator_qq
		FROM posts p
		LEFT JOIN users u ON  p.uid = u.id
		WHERE p.id = ? limit 1`, id)

	if sql.ErrNoRows == err {
		return nil, err
	}
	return post, err
}

func GetPostByIdWithUser(id int, uid int) (*def.Post, error) {
	print(uid, id)
	if uid == -1 {
		return GetPostById(id)
	}

	post := &def.Post{}
	err := db.Get(post, `SELECT p.*,
		u.nickname as creator_nickname,
		u.name as creator_name,
		u.avatar as creator_avatar,
		u.id as creator_id,
		u.bio as creator_bio,
		u.qq as creator_qq,
		COUNT(l.id) as liked_count,
		(EXISTS( SELECT l.id where l.uid = `+fmt.Sprintf("%d", uid)+`)) AS liked
		FROM posts p
		LEFT JOIN users u ON  p.uid = u.id
		LEFT JOIN likes l ON p.id = l.pid
		WHERE p.id = ?
		GROUP BY p.id,u.id,l.id
		limit 1`, id)

	if post.Id == 0 {
		return nil, err
	}
	return post, err
}

func GetPosts(page int, size int, status string, sort string, tag string, uid int, key string, t string, realUid int) (*[]def.Post, error) {
	posts := []def.Post{}
	start := size * (page - 1)
	tags := strings.Fields(tag)

	var query string
	var slice []interface{}
	if status != "10" && status != "" {
		query += ` AND p.status = ` + status
	}

	if sort != "" {
		if sort == "bgm" { //非原创
			query += ` AND NOT p.sort='原创'`
		} else {
			query += ` AND p.sort =?`
			slice = append(slice, sort)
		}
	}

	if uid != 0 {
		query += ` AND p.uid = ?`
		slice = append(slice, uid)
	}

	if t != "" {
		query += ` AND p.type = ?`
		slice = append(slice, t)
	}

	if len(tags) != 0 {
		query += ` AND (1=2 `
		for i := 0; i < len(tags); i++ {
			query += `OR p.tag LIKE ?`
			slice = append(slice, string("%"+tags[i]+"%"))
		}
		query += `)`
	}

	if key != "" {
		query += ` AND (title LIKE ? ) `
		slice = append(slice, "%"+key+"%")
	}

	sqlRaw := `SELECT p.*,
		u.nickname as creator_nickname,
		u.name as creator_name,
		u.avatar as creator_avatar,
		u.id as creator_id,
		u.bio as creator_bio,
		u.qq as creator_qq,
        (select count(lk.id) from likes lk where lk.pid = p.id) as liked_count `

	if realUid > 0 {
		sqlRaw += `,IF(ISNULL((select 1 from likes lk2
				where lk2.pid = p.id and lk2.uid = ` + fmt.Sprintf("%d", realUid) + ` limit 1)),0,1) as liked `
	}

	sqlRaw += `	FROM posts p LEFT JOIN users u ON p.uid = u.id WHERE 1 ` + query +
		` ORDER BY create_time DESC limit ` + fmt.Sprintf("%d,%d", start, size)

	var error error
	if len(slice) == 0 {
		error = db.Select(&posts, sqlRaw)
	} else {
		error = db.Select(&posts, sqlRaw, slice...)
	}

	// print(error.Error())
	return &posts, error
}

func GetRankList() (*[]def.Post, error) {
	posts := []def.Post{}
	err := db.Select(&posts, `SELECT
		p.*,
		u.nickname as creator_nickname,
		u.name as creator_name,
		u.avatar as creator_avatar,
		u.id as creator_id,
		u.bio as creator_bio,
		u.qq as creator_qq,
		IFNULL(pv, 0) as pv
		FROM posts p
		LEFT JOIN pv ON p.id = pv.pid
	 	LEFT JOIN users u ON p.uid = u.id
		where p.type = 'video' AND p.status = 3
		ORDER BY pv.pv DESC limit 0,10`)
	return &posts, err
}

func CreatePost(title string, content string, status string, sort string, tag string, uid int, t string) (*def.Post, error) {
	stmtIns, err := db.Prepare("INSERT INTO posts (title,content,status,sort,tag,uid,type) VALUES (?,?,?,?,?,?,?)")
	if err != nil {
		return nil, err
	}
	resp, _ := stmtIns.Exec(title, content, status, sort, tag, uid, t)
	id, err := resp.LastInsertId()
	if err != nil {
		return nil, err
	}
	res := &def.Post{Title: title, Content: content, Status: status, Sort: sort, Tag: tag, Id: int(id)}
	defer stmtIns.Close()
	return res, err
}

func UpdatePost(id int, title string, content string, status string, sort string, tag string, t string) (*def.Post, error) {
	stmtIns, err := db.Prepare("UPDATE posts SET title=?,content=?,status=?,sort=?,tag=?,type=? WHERE id =?")
	if err != nil {
		return nil, err
	}
	_, err = stmtIns.Exec(&title, &content, &status, &sort, &tag, &t, &id)
	if err != nil {
		return nil, err
	}
	res := &def.Post{Id: id, Title: title, Content: content, Status: status, Sort: sort, Tag: tag, Type: t}
	defer stmtIns.Close()
	return res, err
}

func DeletePost(id int) error {
	stmtDel, err := db.Prepare("UPDATE posts SET status = 0 WHERE id=?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(id)
	if err != nil {
		return err
	}
	defer stmtDel.Close()

	return nil
}

func FindPostByTitleOrContent(key string, status string) (*[]def.Post, error) {
	post := &[]def.Post{}
	var err error
	key = "%" + key + "%"
	if status != "" {
		err = db.Select(post, `SELECT p.*,
			u.nickname as creator_nickname,
			u.name as creator_name,
			u.avatar as creator_avatar,
			u.id as creator_id,
			u.bio as creator_bio,
			u.qq as creator_qq
			FROM posts p LEFT JOIN users u ON  p.uid = u.id WHERE (p.title LIKE ? OR p.content LIKE ?) AND p.status = ?`, key, key, status)
	} else {
		err = db.Select(post, `SELECT p.*,
			u.nickname as creator_nickname,
			u.name as creator_name,
			u.avatar as creator_avatar,
			u.id as creator_id,
			u.bio as creator_bio,
			u.qq as creator_qq
			FROM posts p LEFT JOIN users u ON  p.uid = u.id WHERE p.title LIKE ? OR p.content LIKE ?`, key, key)

	}
	return post, err
}

func DeletePostByIds(ids []interface{}) error {
	sql := fmt.Sprintf("UPDATE posts SET status = '0' WHERE id in (%s)", util.SqlPlaceholders(len(ids)))
	stmtDel, err := db.Prepare(sql)
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	_, err = stmtDel.Exec(ids...)
	if err != nil {
		return err
	}
	defer stmtDel.Close()

	return nil
}

func FindPostByTag(tag string, type_ string) (*[]def.Post, error) {
	post := &[]def.Post{}
	tags := strings.Split(tag, " ")

	var like string
	var slice []interface{}
	for i := 0; i < len(tags); i++ {
		like += "p.tag LIKE ? OR "
		slice = append(slice, "%"+tags[i]+"%")
	}
	like = strings.TrimSuffix(like, "OR ")
	slice = append(slice, 3, type_)

	err := db.Select(post, `SELECT p.*,
		u.nickname as creator_nickname,
		u.name as creator_name,
		u.avatar as creator_avatar,
		u.id as creator_id,
		u.bio as creator_bio,
		u.qq as creator_qq
		FROM posts p LEFT JOIN users u ON  p.uid = u.id WHERE (`+like+`) AND p.status = ? AND p.type=?`, slice...)

	return post, err
}

//获取用户点赞的文章
func FindPostByUserLike(uid int) (*[]def.Post, error) {
	post := &[]def.Post{}
	err := db.Select(post, `SELECT p.*,
		u.nickname as creator_nickname,
		u.name as creator_name,
		u.avatar as creator_avatar,
		u.id as creator_id,
		u.bio as creator_bio,
		u.qq as creator_qq
		FROM posts p LEFT JOIN users u ON  p.uid = u.id WHERE p.id in (SELECT pid FROM likes WHERE uid = ?)`, uid)

	return post, err
}

// 点赞
func LikePost(pid int, uid int) error {
	var like int
	if err := db.Get(&like, `SELECT id FROM likes WHERE pid = ? AND uid = ?`, pid, uid); err != nil && err != sql.ErrNoRows {
		return err
	}

	stmtIns, err := db.Prepare("INSERT INTO likes (pid,uid) VALUES (?,?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(pid, uid)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

// 取消点赞
func UnLikePost(pid int, uid int) error {
	stmtDel, err := db.Prepare("DELETE FROM likes WHERE pid=? AND uid=?")
	if err != nil {
		return err
	}
	_, err = stmtDel.Exec(pid, uid)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

//收藏文章
func CollectPost(pid int, uid int) error {
	if collect, err := IsCollectedPost(pid, uid); err != nil {
		return err
	} else if collect {
		return errors.New("已经收藏")
	}

	stmtIns, err := db.Prepare("INSERT INTO collects (pid,uid) VALUES (?,?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(pid, uid)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func IsCollectedPost(pid int, uid int) (bool, error) {
	var collect int
	err := db.Get(&collect, `SELECT id FROM collects WHERE pid = ? AND uid = ? limit 1`, pid, uid)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return collect > 0, err
}

//获取用户收藏的文章
func FindUserCollectPost(uid int) (*[]def.Post, error) {
	post := &[]def.Post{}
	err := db.Select(post, `SELECT p.*,
		u.nickname as creator_nickname,
		u.name as creator_name,
		u.avatar as creator_avatar,
		u.id as creator_id,
		u.bio as creator_bio,
		u.qq as creator_qq,
		(select count(lk.id) from likes lk where lk.pid = p.id) as liked_count,
		IF(ISNULL((select 1 from likes lk2 where lk2.pid = p.id and lk2.uid = `+fmt.Sprintf("%d", uid)+` limit 1)),0,1) as liked
		FROM posts p LEFT JOIN users u ON  p.uid = u.id
		WHERE p.id in (SELECT pid FROM collects WHERE uid = ?)`, uid)
	return post, err
}

//取消收藏
func UnCollectPost(pid int, uid int) error {
	stmtDel, err := db.Prepare("DELETE FROM collects WHERE pid=? AND uid=?")
	if err != nil {
		return err
	}
	_, err = stmtDel.Exec(pid, uid)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}
