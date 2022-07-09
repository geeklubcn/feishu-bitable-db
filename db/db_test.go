package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/geeklubcn/feishu-bitable-db/internal/faker"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	db, _ := NewDB(faker.AppID, faker.AppSecret)
	ctx := context.Background()

	dbName := "test-database"
	tableName := "test-table"

	t.Run("create database", func(t *testing.T) {
		did, err := db.SaveDatabase(ctx, dbName)
		assert.NoError(t, err)
		assert.NotEmpty(t, did)
		did2, err := db.SaveDatabase(ctx, dbName)
		assert.NoError(t, err)
		assert.Equal(t, did, did2)
	})

	t.Run("create table", func(t *testing.T) {
		tid, err := db.SaveTable(ctx, dbName, Table{
			Name: tableName,
			Fields: []Field{
				{Name: "username", Type: String},
				{Name: "age", Type: Int},
			},
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, tid)
		tid2, err := db.SaveTable(ctx, dbName, Table{
			Name: tableName,
			Fields: []Field{
				{Name: "username", Type: String},
				{Name: "passport", Type: String},
				{Name: "age", Type: Int},
			},
		})
		assert.NoError(t, err)
		assert.Equal(t, tid, tid2)
	})

	t.Run("create record", func(t *testing.T) {
		_, _ = db.SaveTable(ctx, dbName, Table{
			Name: tableName,
			Fields: []Field{
				{Name: "username", Type: String},
				{Name: "age", Type: Int},
			},
		})
		id, err := db.Create(ctx, dbName, tableName, map[string]interface{}{
			"username": "zhangsan",
			"age":      12,
		})
		err = db.Update(ctx, dbName, tableName, id, map[string]interface{}{
			"username": "zhangsan13",
			"age":      13,
		})
		assert.NoError(t, err)
		err = db.Delete(ctx, dbName, tableName, id)
		assert.NoError(t, err)
	})

	t.Run("read record", func(t *testing.T) {
		res := db.Read(ctx, dbName, tableName, []SearchCmd{
			{
				"age",
				"=",
				12,
			},
		})
		fmt.Println(res)
	})

}
