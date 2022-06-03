package db

import (
	"database/sql"

	"github.com/cliclitv/go-clicli/def"
)

func CreateComment(uid int, pid int, content string) (*def.Comment, error) {
	comment := &def.Comment{}
	if err := db.Get(comment, `INSERT INTO comments (uid, pid, content) VALUES (?, ?, ?)`, uid, pid, content); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return comment, nil
}

func GetPostComments(id int) (*[]def.Comment, error) {
	comments := &[]def.Comment{}
	err := db.Select(comments, `SELECT c.*,
		u.nickname as creator_nickname,
		u.name as creator_name,
		u.avatar as creator_avatar,
		u.id as creator_id,
		u.bio as creator_bio,
		u.qq as creator_qq
		FROM comments c LEFT JOIN users u ON  c.uid = u.id WHERE c.pid = ?
		ORDER BY c.create_time DESC`, id)
	return comments, err
}
