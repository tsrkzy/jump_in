package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/tsrkzy/jump_in/helper"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"
)

func Open() (*MyDB, error) {
	return (&MyDB{}).Connect()
}

type MyDB struct {
	db *sql.DB
}

func (d *MyDB) GetDB() *sql.DB {
	return d.db
}

func (d *MyDB) Connect() (*MyDB, error) {
	var (
		/* db名 */
		dbName = helper.MustGetenv("PG_DB_NAME")
		/* ユーザとPW */
		dbUser = helper.MustGetenv("PG_DB_USER")
		dbPwd  = helper.MustGetenv("PG_DB_PASS")
		/* 接続先ホスト */
		dbTCPHost = helper.MustGetenv("PG_INSTANCE_HOST")
		/* port */
		dbPort = helper.MustGetenv("PG_DB_PORT")
	)

	dbURI := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbTCPHost, dbPort, dbName, dbUser, dbPwd,
	)

	JST := time.FixedZone("Local", +9*60*60)
	boil.SetLocation(JST)
	db, err := sql.Open("postgres", dbURI)
	d.db = db
	if err != nil {
		return d, fmt.Errorf("cannot open database connection: %v", err)
	}

	return d, nil
}

func (d *MyDB) Close() {
	err := d.db.Close()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

// Tx
//
// @REF https://qiita.com/fiemon/items/eb38c8d681ed1ae05925
// txFuncの中で
//  - panic: panic
//  - return error: rollback
//  - return nil: commit
func (d *MyDB) Tx(ctx context.Context, txFunc func(*sql.Tx) error) error {
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			/* func Transaction で panic が発生したらそのまま外にpanicで流す */
			log.Error(p)
			err := tx.Rollback()
			if err != nil {
				log.Error("rollback failed")
				log.Error(err)
			}
			panic(p)
		} else if err != nil {
			log.Error("rollback")
			_ = tx.Rollback()
		} else {
			log.Info("commit")
			err = tx.Commit()
		}
	}()

	err = txFunc(tx)

	return err
}
