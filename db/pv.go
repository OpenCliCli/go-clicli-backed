package db

import (
	"database/sql"

	"github.com/cliclitv/go-clicli/def"
)

func GetPv(pid int) (*def.Pv, error) {
	stmtCount, _ := db.Prepare("SELECT pv FROM pv WHERE pid = ?")
	var pv int
	err = stmtCount.QueryRow(pid).Scan(&pv)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	res := &def.Pv{Pid: pid, Pv: pv}

	defer stmtCount.Close()
	return res, nil
}

func ReplacePv(pid int, pv int) (*def.Pv, error) {
	stmtIns, err := db.Prepare("REPLACE INTO pv (pid,pv) VALUES (?,?)")
	if err != nil {
		return nil, err
	}
	_, err = stmtIns.Exec(pid, pv)
	if err != nil {
		return nil, err
	}
	defer stmtIns.Close()

	res := &def.Pv{Pid: pid, Pv: pv}
	return res, nil
}
