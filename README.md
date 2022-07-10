# feishu-bitable-db

把飞书多维表格作为底层存储, 实现数据库的基本操作(CRUD).

DDL
- 新建数据库 `SaveDatabase(ctx context.Context, name string) (id string, err error)`
- 新建数据表 `SaveTable(ctx context.Context, database string, table Table) (string, error)`
- 查看所有数据表 `ListTables(ctx context.Context, database string) []string`
- 删除数据表 `DropTable(ctx context.Context, database, table string) error`

DML
- 新建记录 `Create(ctx context.Context, database, table string, record map[string]interface{}) (id string, err error)`
- 查询记录 `Read(ctx context.Context, database, table string, ss []SearchCmd) []map[string]interface{}`
- 修改记录 `Update(ctx context.Context, database, table, id string, record map[string]interface{}) error`
- 删除记录 `Delete(ctx context.Context, database, table, id string) error`

## 使用方式

需要确认飞书app已经申请了相关权限, 详情参考 https://open.feishu.cn/document/ukTMukTMukTM/uYTM5UjL2ETO14iNxkTN/scope-list

- drive:drive
- bitable:app

```go
package main

import (
	"context"
	"fmt"
	"github.com/geeklubcn/feishu-bitable-db/db"
)

func main() {
	it, err := db.NewDB("${appId}", "${appSecret}")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	_, _ = it.SaveDatabase(ctx, "MyDB")
	_, _ = it.SaveTable(ctx, "MyDB", db.Table{
		Name: "MyTable",
		Fields: []db.Field{
			{Name: "username", Type: db.String},
			{Name: "passport", Type: db.String},
		},
	})
	// alter table fields
	_, _ = it.SaveTable(ctx, "MyDB", db.Table{
		Name: "MyTable",
		Fields: []db.Field{
			{Name: "username", Type: db.String},
			{Name: "age", Type: db.Int},
		},
	})
	rid, _ := it.Create(ctx, "MyDB", "MyTable", map[string]interface{}{
		"username": "zhangsan",
		"age":      12,
	})
	_ = it.Update(ctx, "MyDB", "MyTable", rid, map[string]interface{}{
		"age": 13,
	})
	res := it.Read(ctx, "MyDB", "MyTable", []db.SearchCmd{
		{"age", "=", 13},
	})

	fmt.Println(res) // [map[age:13 id:rec3gAKlML username:zhangsan]]
}

```
