package db

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/geeklubcn/feishu-bitable-db/client/bitable"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	larkBitable "github.com/larksuite/oapi-sdk-go/service/bitable/v1"
	larkDriveExplorer "github.com/larksuite/oapi-sdk-go/service/drive_explorer/v2"
	"github.com/sirupsen/logrus"
)

type DB interface {
	// SaveDatabase create database if not exists.
	SaveDatabase(ctx context.Context, name string) (id string, err error)
	// SaveTable create or update table.
	SaveTable(ctx context.Context, database string, table Table) (string, error)
	// ListTables list table by database.
	ListTables(ctx context.Context, database string) []string
	// DropTable should return err if database or table not exists.
	DropTable(ctx context.Context, database, table string) error

	// Create table record. id is generated.
	Create(ctx context.Context, database, table string, record map[string]interface{}) (id string, err error)
	// Read table records.
	Read(ctx context.Context, database, table string, ss []SearchCmd) []map[string]interface{}
	// Update record by id.
	Update(ctx context.Context, database, table, id string, record map[string]interface{}) error
	// Delete record by id.
	Delete(ctx context.Context, database, table, id string) error
}

type db struct {
	conf    *config.Config
	service *larkBitable.Service
	bitable bitable.Bitable
	cache   sync.Map
	token   string
	userId  string
}

func NewDB(appId string, appSecret string) (DB, error) {
	logrus.Debugf("init lark cli. appId:%s, appSecret:%s", appId, appSecret)
	appSettings := core.NewInternalAppSettings(
		core.SetAppCredentials(appId, appSecret),
	)
	conf := core.NewConfig(core.DomainFeiShu, appSettings, core.SetLoggerLevel(core.LoggerLevelError))

	ctx := context.Background()
	c := core.WrapContext(ctx)

	res, err := larkDriveExplorer.NewService(conf).Folders.RootMeta(c).Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("Get RootMeta fail! response:%s", tools.Prettify(res))
		return nil, err
	}

	return &db{
		conf:    conf,
		token:   res.Token,
		userId:  res.UserId,
		service: larkBitable.NewService(conf),
		bitable: bitable.NewBitable(conf),
	}, nil
}

func (b *db) SaveDatabase(ctx context.Context, database string) (id string, err error) {
	if did, exist := b.getDid(ctx, database); exist {
		return did, nil
	}
	did, err := b.bitable.CreateApp(ctx, database, b.token)
	if err == nil {
		b.cache.Store(fmt.Sprintf("db-%s", database), did)
	}
	return did, err
}

func (b *db) SaveTable(ctx context.Context, database string, table Table) (string, error) {
	var tableId string
	c := core.WrapContext(ctx)

	did, _ := b.SaveDatabase(ctx, database)

	tables := b.listTables(ctx, database)
	if tid, exist := tables[table.Name]; exist {
		tableId = tid
	} else {
		tableId, _ = b.saveTable(c, did, table.Name)
	}

	fs, _ := b.ListFields(ctx, did, tableId)
	ofm := make(map[string]*larkBitable.AppTableField, 0)
	for _, it := range fs {
		ofm[it.FieldName] = it
	}

	// id
	if fs[0].FieldName != "id" {
		_ = b.UpdateField(ctx, did, tableId, &larkBitable.AppTableField{
			FieldId:   fs[0].FieldId,
			FieldName: "id",
			Type:      int(String),
		})
		delete(ofm, fs[0].FieldName)
	}
	delete(ofm, "id")

	for _, f := range table.Fields {
		of, exist := ofm[f.Name]
		if exist {
			if f.Name == of.FieldName && int(f.Type) == of.Type {
				delete(ofm, f.Name)
				continue
			}
			_ = b.UpdateField(ctx, did, tableId, &larkBitable.AppTableField{
				FieldId:   of.FieldId,
				FieldName: f.Name,
				Type:      int(f.Type),
			})
			delete(ofm, f.Name)
		} else {
			_, _ = b.CreateField(ctx, did, tableId, &larkBitable.AppTableField{
				FieldName: f.Name,
				Type:      int(f.Type),
			})
		}
	}

	for _, f := range ofm {
		_ = b.DeleteField(ctx, did, tableId, f.FieldId)
	}
	return tableId, nil
}

func (b *db) ListTables(ctx context.Context, database string) []string {
	m := b.listTables(ctx, database)
	res := make([]string, len(m))
	for k := range m {
		res = append(res, k)
	}
	return res
}

func (b *db) DropTable(ctx context.Context, database, table string) error {
	c := core.WrapContext(ctx)
	did, exist := b.getDid(ctx, database)
	if !exist {
		return errors.New(fmt.Sprintf("database[%s] not exists", database))
	}
	tid, exist := b.getTid(ctx, database, table)
	if !exist {
		return errors.New(fmt.Sprintf("table[%s.%s] not exists", database, table))
	}

	reqCall := b.service.AppTables.Delete(c)
	reqCall.SetAppToken(did)
	reqCall.SetTableId(tid)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("ListTables fail! database:%s,table:%s,response:%s", did, tid, tools.Prettify(message))
		return err
	}
	logrus.WithContext(c).Debugf("response:%s", tools.Prettify(message))
	return nil
}

func (b *db) listTables(ctx context.Context, database string) map[string]string {
	c := core.WrapContext(ctx)

	did, _ := b.SaveDatabase(ctx, database)
	reqCall := b.service.AppTables.List(c)
	reqCall.SetAppToken(did)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("ListTables fail! database: %s, appToken:%s, response:%s", database, did, tools.Prettify(message))
		return nil
	}
	logrus.WithContext(c).Debugf("response:%s", tools.Prettify(message))
	res := make(map[string]string, len(message.Items))
	for _, it := range message.Items {
		res[it.Name] = it.TableId
	}
	return res
}

func (b *db) saveTable(ctx context.Context, appToken, name string) (string, error) {
	c := core.WrapContext(ctx)
	body := &larkBitable.AppTableCreateReqBody{
		Table: &larkBitable.ReqTable{
			Name: name,
		},
	}
	reqCall := b.service.AppTables.Create(c, body)
	reqCall.SetAppToken(appToken)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("CreateTable fail! appToken:%s,response:%s", appToken, tools.Prettify(message))
		return "", err
	}
	logrus.WithContext(c).Debugf("response:%s", tools.Prettify(message))
	return message.TableId, nil
}

func (b *db) getDid(ctx context.Context, database string) (string, bool) {
	if v, ok := b.cache.Load(fmt.Sprintf("db-%s", database)); ok {
		return v.(string), true
	}
	if did, exist := b.bitable.QueryByName(ctx, database, b.token); exist {
		b.cache.Store(fmt.Sprintf("db-%s", database), did)
		return did, true
	}
	return "", false
}

func (b *db) getTid(ctx context.Context, database, table string) (string, bool) {
	if v, ok := b.cache.Load(fmt.Sprintf("table-%s-%s", database, table)); ok {
		return v.(string), true
	}
	tables := b.listTables(ctx, database)
	for name, id := range tables {
		b.cache.Store(fmt.Sprintf("table-%s-%s", database, name), id)
	}
	if tid, exists := tables[table]; exists {
		return tid, true
	}
	return "", false
}
