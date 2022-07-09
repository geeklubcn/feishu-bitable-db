package db

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	larkBitable "github.com/larksuite/oapi-sdk-go/service/bitable/v1"
	"github.com/sirupsen/logrus"
)

func (b *db) Create(ctx context.Context, database, table string, record map[string]interface{}) (id string, err error) {
	c := core.WrapContext(ctx)
	did, exist := b.getDid(ctx, database)
	if !exist {
		return "", errors.New(fmt.Sprintf("database[%s] not exists", database))
	}
	tid, exist := b.getTid(ctx, database, table)
	if !exist {
		return "", errors.New(fmt.Sprintf("table[%s.%s] not exists", database, table))
	}

	fs, _ := b.ListFields(ctx, did, tid)
	mfs := make(map[string]string, 0)
	for _, it := range fs {
		mfs[it.FieldName] = it.FieldId
	}
	if id, exist := mfs["id"]; exist {
		record["id"] = ""
		defer func() {
			_ = b.Update(ctx, database, table, id, map[string]interface{}{})
		}()
	}
	body := &larkBitable.AppTableRecord{Fields: record}

	logrus.WithContext(c).Infof("request: %s", tools.Prettify(body))
	reqCall := b.service.AppTableRecords.Create(c, body)
	reqCall.SetAppToken(did)
	reqCall.SetTableId(tid)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("BatchCreateRecord fail! database:%s,table:%s,response:%s", database, table, tools.Prettify(message))
		return "", err
	}
	logrus.WithContext(c).Debugf("response:%s", tools.Prettify(message))
	id = message.Record.RecordId
	return id, nil
}

func (b *db) Read(ctx context.Context, database, table string, ss []SearchCmd) []map[string]interface{} {
	c := core.WrapContext(ctx)
	did, exist := b.getDid(ctx, database)
	if !exist {
		return nil
	}
	tid, exist := b.getTid(ctx, database, table)
	if !exist {
		return nil
	}
	filters := make([]string, 0)
	for _, s := range ss {
		switch v := s.Val.(type) {
		case string:
			filters = append(filters, fmt.Sprintf("CurrentValue.[%s]%s\"%s\"", s.Key, s.Operator, v))
		case int, int8, int16, int32, int64:
			filters = append(filters, fmt.Sprintf("CurrentValue.[%s]%s%d", s.Key, s.Operator, v))
		default:
			filters = append(filters, fmt.Sprintf("CurrentValue.[%s]%s%s", s.Key, s.Operator, v))
		}
	}
	filter := ""
	if len(filters) > 1 {
		filter = "AND(" + strings.Join(filters, ",") + ")"
	} else {
		filter = "AND(" + strings.Join(filters, "") + ")"
	}
	reqCall := b.service.AppTableRecords.List(c)
	reqCall.SetAppToken(did)
	reqCall.SetTableId(tid)
	reqCall.SetFilter(filter)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("ListRecords fail! database:%s, table:%s, filter:%s, response:%s", database, table, filter, tools.Prettify(message))
		return nil
	}
	logrus.WithContext(c).Debugf("response:%s", tools.Prettify(message))
	res := make([]map[string]interface{}, 0)
	for _, it := range message.Items {
		record := make(map[string]interface{}, 0)
		for k, v := range it.Fields {
			record[k] = v
		}
		res = append(res, record)
	}

	return res
}

func (b *db) Update(ctx context.Context, database, table, id string, record map[string]interface{}) error {
	c := core.WrapContext(ctx)
	did, exist := b.getDid(ctx, database)
	if !exist {
		return errors.New(fmt.Sprintf("database[%s] not exists", database))
	}
	tid, exist := b.getTid(ctx, database, table)
	if !exist {
		return errors.New(fmt.Sprintf("table[%s.%s] not exists", database, table))
	}

	fs, _ := b.ListFields(ctx, did, tid)
	mfs := make(map[string]string, 0)
	for _, it := range fs {
		mfs[it.FieldName] = it.FieldId
	}
	if id, exist := mfs["id"]; exist {
		record["id"] = id
	}

	body := &larkBitable.AppTableRecord{
		Fields: record,
	}

	reqCall := b.service.AppTableRecords.Update(c, body)
	reqCall.SetAppToken(did)
	reqCall.SetTableId(tid)
	reqCall.SetRecordId(id)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("BatchUpdateRecord fail! appToken:%s, tableId:%s, response:%s", did, tid, tools.Prettify(message))
		return err
	}
	logrus.WithContext(c).Debugf("response:%s", tools.Prettify(message))
	return nil
}

func (b *db) Delete(ctx context.Context, database, table, id string) error {
	c := core.WrapContext(ctx)
	did, exist := b.getDid(ctx, database)
	if !exist {
		return errors.New(fmt.Sprintf("database[%s] not exists", database))
	}
	tid, exist := b.getTid(ctx, database, table)
	if !exist {
		return errors.New(fmt.Sprintf("table[%s.%s] not exists", database, table))
	}

	reqCall := b.service.AppTableRecords.Delete(c)
	reqCall.SetAppToken(did)
	reqCall.SetTableId(tid)
	reqCall.SetRecordId(id)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("BatchUpdateRecord fail! appToken:%s, tableId:%s, response:%s", did, tid, tools.Prettify(message))
		return err
	}
	logrus.WithContext(c).Debugf("response:%s", tools.Prettify(message))
	return nil
}
