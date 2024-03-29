package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	larkBitable "github.com/larksuite/oapi-sdk-go/service/bitable/v1"
	"github.com/sirupsen/logrus"
)

func (b *db) Create(ctx context.Context, database, table string, record map[string]interface{}) (id string, err error) {
	c := core.WrapContext(ctx)

	fs, _ := b.ListFields(ctx, database, table)
	mfs := make(map[string]string, 0)
	for _, it := range fs {
		mfs[it.FieldName] = it.FieldId
	}
	if _, exist := mfs[ID]; exist {
		record[ID] = ""
		defer func() {
			_ = b.Update(ctx, database, table, id, map[string]interface{}{})
		}()
	}
	body := &larkBitable.AppTableRecord{Fields: record}

	logrus.WithContext(c).Infof("request: %s", tools.Prettify(body))
	reqCall := b.service.AppTableRecords.Create(c, body)
	reqCall.SetAppToken(database)
	reqCall.SetTableId(table)
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
	switch len(filters) {
	case 0:
	case 1:
		filter = "AND(" + strings.Join(filters, "") + ")"
	default:
		filter = "AND(" + strings.Join(filters, ",") + ")"
	}

	reqCall := b.service.AppTableRecords.List(c)
	reqCall.SetAppToken(database)
	reqCall.SetTableId(table)
	reqCall.SetFilter(filter)
	reqCall.SetPageSize(1000)
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
		record[ID] = it.RecordId
		res = append(res, record)
	}

	return res
}

func (b *db) Update(ctx context.Context, database, table, id string, record map[string]interface{}) error {
	c := core.WrapContext(ctx)

	fs, _ := b.ListFields(ctx, database, table)
	mfs := make(map[string]string, 0)
	for _, it := range fs {
		mfs[it.FieldName] = it.FieldId
	}
	if _, exist := mfs[ID]; exist {
		record[ID] = id
	}

	body := &larkBitable.AppTableRecord{
		Fields: record,
	}

	reqCall := b.service.AppTableRecords.Update(c, body)
	reqCall.SetAppToken(database)
	reqCall.SetTableId(table)
	reqCall.SetRecordId(id)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("BatchUpdateRecord fail! appToken:%s, tableId:%s, response:%s", database, table, tools.Prettify(message))
		return err
	}
	logrus.WithContext(c).Debugf("response:%s", tools.Prettify(message))
	return nil
}

func (b *db) Delete(ctx context.Context, database, table, id string) error {
	c := core.WrapContext(ctx)

	reqCall := b.service.AppTableRecords.Delete(c)
	reqCall.SetAppToken(database)
	reqCall.SetTableId(table)
	reqCall.SetRecordId(id)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("BatchUpdateRecord fail! appToken:%s, tableId:%s, response:%s", database, table, tools.Prettify(message))
		return err
	}
	logrus.WithContext(c).Debugf("response:%s", tools.Prettify(message))
	return nil
}
