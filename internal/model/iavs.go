package model

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type (
	IAV struct {
		db *sql.DB
	}
	IAVModel struct {
		Id      int64  `db:"id" json:"id"`
		Name    string `db:"name" json:"name"`
		Path    string `db:"path" json:"path"`
		Md5     string `db:"md5" json:"md5"`
		Type    string `db:"type" json:"type"`
		Suffix  string `db:"suffix" json:"suffix"`
		Tags    string `db:"tags" json:"tags"`
		ShootAt string `db:"shoot_at" json:"shoot_at"`
		Lng     string `db:"lng" json:"lng"`
		Lat     string `db:"lat" json:"lat"`
	}
)

func (p *IAV) Insert(pm *IAVModel) error {
	stmt, err := p.db.Prepare("INSERT INTO iavs(name, path, md5,type,suffix,tags,shoot_at,lng,lat)" +
		"values(?,?,?,?,?,?,?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(pm.Name, pm.Path, pm.Md5, pm.Type, pm.Suffix, pm.Tags, pm.ShootAt, pm.Lng, pm.Lat)
	return err
}

func (p *IAV) Query(pm *IAVModel) (*IAVModel, error) {
	row := p.db.QueryRow("SELECT * FROM iavs where md5 = ? limit 1", pm.Md5)
	r := &IAVModel{}
	err := row.Scan(&r.Id, &r.Name, &r.Path, &r.Md5, &r.Type, &r.Suffix, &r.Tags, &r.ShootAt, &r.Lng, &r.Lat)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return r, err
}

func NewIAVSModel(datasource string) *IAV {
	db, err := sql.Open("sqlite3", datasource)
	if err != nil {
		panic(err.Error())
	}
	dq := `drop table if exists "iavs"`
	db.Exec(dq)

	cq := `
CREATE TABLE IF NOT EXISTS "iavs" (
  "id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  "name" text(255) NOT NULL,
  "path" text NOT NULL DEFAULT '',
  "md5" text(255) NOT NULL DEFAULT '',
  "type" text(50) NOT NULL,
  "suffix" text(20) NOT NULL,
  "tags" text NOT NULL,
  "shoot_at" TEXT(50) NOT NULL,
  "lng" TEXT(50) NOT NULL,
  "lat" TEXT(50) NOT NULL
);
`
	db.Exec(cq)
	return &IAV{
		db: db,
	}
}
