package database

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMyDB_Connect(t *testing.T) {
	mydb, err := Open()
	assert.NoError(t, err)

	defer mydb.Close()

	ctx := context.Background()
	err = mydb.Tx(ctx, func(tx *sql.Tx) error {
		return nil
	})
	assert.NoError(t, err)
}

// TestMyDB_Connect2 同じDBに対して別の接続をOpenし、それぞれ別のトランザクションを開始できるか
func TestMyDB_Connect2(t *testing.T) {
	mydb1, err := Open()
	assert.NoError(t, err)
	mydb2, err := Open()
	assert.NoError(t, err)

	defer mydb1.Close()
	defer mydb2.Close()

	ctx1 := context.Background()
	ctx2 := context.Background()

	err = mydb1.Tx(ctx1, func(tx *sql.Tx) error {
		return mydb2.Tx(ctx2, func(tx *sql.Tx) error {
			err := mydb1.db.Ping()
			assert.NoError(t, err)

			err = mydb2.db.Ping()
			assert.NoError(t, err)

			q := `select "txid_current"()`

			q1, err := mydb1.db.QueryContext(ctx1, q)
			assert.NoError(t, err)
			defer func(q1 *sql.Rows) {
				err := q1.Close()
				assert.NoError(t, err)
			}(q1)

			q2, err := mydb2.db.QueryContext(ctx2, q)
			assert.NoError(t, err)
			defer func(q2 *sql.Rows) {
				err := q2.Close()
				assert.NoError(t, err)
			}(q2)

			var (
				txId1 int64
				txId2 int64
			)

			for q1.Next() {
				err := q1.Scan(&txId1)
				assert.NoError(t, err)
				fmt.Printf("txId1 is: %d\n", txId1)
			}

			for q2.Next() {
				err := q2.Scan(&txId2)
				assert.NoError(t, err)
				fmt.Printf("txId2 is: %d\n", txId2)
			}

			assert.NotEqual(t, txId1, txId2)

			return nil
		})

	})
	assert.NoError(t, err)
}

// TestMyDB_Connect3 同じDBに対しての接続をOpenし、2つの独立したトランザクションを開始できるか
func TestMyDB_Connect3(t *testing.T) {
	mydb, err := Open()
	assert.NoError(t, err)
	defer mydb.Close()

	ctx := context.Background()

	err = mydb.Tx(ctx, func(tx *sql.Tx) error {
		return mydb.Tx(ctx, func(tx *sql.Tx) error {
			err := mydb.db.Ping()
			assert.NoError(t, err)

			err = mydb.db.Ping()
			assert.NoError(t, err)

			q := `select "txid_current"()`

			q1, err := mydb.db.QueryContext(ctx, q)
			assert.NoError(t, err)
			defer func(q1 *sql.Rows) {
				err := q1.Close()
				assert.NoError(t, err)
			}(q1)

			q2, err := mydb.db.QueryContext(ctx, q)
			assert.NoError(t, err)
			defer func(q2 *sql.Rows) {
				err := q2.Close()
				assert.NoError(t, err)
			}(q2)

			var (
				txId1 int64
				txId2 int64
			)

			for q1.Next() {
				err := q1.Scan(&txId1)
				assert.NoError(t, err)
				fmt.Printf("txId1 is: %d\n", txId1)
			}

			for q2.Next() {
				err := q2.Scan(&txId2)
				assert.NoError(t, err)
				fmt.Printf("txId2 is: %d\n", txId2)
			}

			assert.NotEqual(t, txId1, txId2)

			return nil
		})

	})
	assert.NoError(t, err)
}
