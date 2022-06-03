package db

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/cliclitv/go-clicli/def"
	"github.com/cliclitv/go-clicli/util"
)

func CreateSimpleUser(name string, pwd string) error {
	pwd = util.Cipher(pwd)
	stmtIns, err := db.Prepare("INSERT INTO users (name, pwd, level) VALUES(?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(&name, &pwd, 3)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

func UpdateUser(id int, nickname string, pwd string, level int, qq string, bio string, status string, avatar string) (*def.User, error) {
	if pwd == "" {
		stmtIns, err := db.Prepare("UPDATE users SET nickname=?,level=?,qq=?,bio=?,status=?,avatar=? WHERE id =?")
		if err != nil {
			fmt.Println("status", err)
			return nil, err
		}
		_, err = stmtIns.Exec(&nickname, &level, &qq, &bio, &status, &avatar, &id)
		if err != nil {
			return nil, err
		}
		defer stmtIns.Close()

		res := &def.User{Id: id, Nickname: nickname, QQ: qq, Level: level, Status: status, Avatar: avatar}
		defer stmtIns.Close()
		return res, err
	} else {
		pwd = util.Cipher(pwd)
		stmtIns, err := db.Prepare("UPDATE users SET nickname=?,pwd=?,level=?,qq=?,bio=?,status=?,avatar=? WHERE id =?")
		if err != nil {
			return nil, err
		}
		_, err = stmtIns.Exec(&nickname, &pwd, &level, &qq, &bio, &status, &avatar, &id)
		if err != nil {
			return nil, err
		}
		defer stmtIns.Close()

		res := &def.User{Id: id, Nickname: nickname, Pwd: pwd, QQ: qq, Level: level, Status: status, Avatar: avatar}
		return res, err
	}
}

func GetUserByName(name string) (*def.User, error) {
	user := &def.User{}
	err = db.Get(user, "SELECT * FROM users WHERE name = ?", name)
	if user.Id == 0 {
		return nil, errors.New("user not found")
	}
	return user, err
}

func GetUser(name string, id int) (*def.User, error) {
	query := `SELECT users.id,
		users.name,
		users.nickname,
		users.level,
		users.qq,
		users.bio,
		users.status,
	 	users.create_time,
		users.update_time,
		users.avatar
		FROM users WHERE `
	if name != "" {
		query += `name = ?`
	} else if id != 0 {
		query += `id = ?`
	}
	user := &def.User{}

	if name != "" {
		err = db.Get(user, query, name)
	} else if id != 0 {
		err = db.Get(user, query, id)
	}
	if user.Id == 0 {
		return nil, errors.New("user not found")
	}

	return user, err
}

func GetUsers(name string, level int, page int, pageSize int, order string, status string) (*[]def.User, error) {
	start := pageSize * (page - 1)
	var slice []interface{}
	query := `SELECT
		users.id,
		users.name,
		users.nickname,
		users.level,
		users.qq,
		users.bio,
		users.status,
	 	users.create_time,
		users.update_time,
		users.avatar
		FROM users WHERE 1 `

	if level != 10 {
		query += " AND level = ? "
		slice = append(slice, level)
	}

	if name != "" {
		query += " AND nickname LIKE ? "
		slice = append(slice, "%"+name+"%")
	}

	if status == "" {
		//query += " AND NOT status = 0  "
	} else {
		s, _ := strconv.Atoi(status)
		if s != 10 { // 10 查询所有 status
			query += " AND status = ?  "
			slice = append(slice, status)
		}
	}

	if order == "" || (strings.ToUpper(order) != "DESC" && strings.ToUpper(order) != "ASC") {
		order = " DESC "

	}

	query += "ORDER BY create_time " + order
	query += " limit ?,?"
	slice = append(slice, start, pageSize)

	fmt.Println(query, slice)

	res := []def.User{}
	var err error

	if len(slice) > 0 {
		err = db.Select(&res, query, slice...)
	} else {
		err = db.Select(&res, query)
	}

	fmt.Println(err)
	return &res, err
}

func SearchUsers(key string) ([]*def.User, error) {
	key = string("%" + key + "%")
	stmt, _ := db.Prepare("SELECT * FROM users WHERE name LIKE ?")

	var res []*def.User

	rows, err := stmt.Query(key)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var id, level int
		var name, bio, qq string
		if err := rows.Scan(&id, &name, &level, &qq, &bio); err != nil {
			return res, err
		}

		c := &def.User{Id: id, Name: name, Level: level, QQ: qq}
		res = append(res, c)
	}
	defer stmt.Close()

	return res, nil
}

func DeleteUser(id int) error {
	stmtDel, err := db.Prepare("UPDATE users SET status = '0' WHERE id =?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}
	_, err = stmtDel.Exec(id)
	if err != nil {
		return err
	}
	defer stmtDel.Close()

	return nil
}

func DeleteUserByIds(ids []interface{}) error {
	sql := fmt.Sprintf("UPDATE users SET status = '0' WHERE id in (%s)", util.SqlPlaceholders(len(ids)))
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

func GetUserStat(id int) (interface{}, error) {
	userStat := def.UserStat{}
	err = db.Get(&userStat, `
	SELECT
	u.id,
	u.name,
	u.nickname,
	u.level,
	u.qq,
	u.bio,
	u.status,
	u.create_time,
	u.update_time,
	u.avatar,
	(SELECT COUNT(*) FROM posts p WHERE p.uid = u.id and p.type='post') AS post_count,
	(SELECT COUNT(*) FROM posts v WHERE v.uid = u.id and v.type='video') AS video_count,
	(SELECT COUNT(*) FROM likes l WHERE l.uid = u.id) as liked_count,
	(SELECT COUNT(*) FROM comments c WHERE c.uid = u.id) as comment_count,
	(SELECT COUNT(*) FROM collects co WHERE co.uid = u.id) as collect_count,
	IFNULL((SELECT SUM(pv.pv) FROM pv pv WHERE pv.pid in (SELECT id from posts WHERE uid=u.id)),0) as pv_count
	FROM users u
	WHERE u.id = ? limit 1`, id)
	return &userStat, err
}
