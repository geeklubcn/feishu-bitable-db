package db

import (
	"context"

	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/tools"
	larkBitable "github.com/larksuite/oapi-sdk-go/service/bitable/v1"
	"github.com/sirupsen/logrus"
)

type FieldManager interface {
	ListFields(ctx context.Context, appToken string, tableId string) ([]*larkBitable.AppTableField, error)
	CreateField(ctx context.Context, appToken string, tableId string, body *larkBitable.AppTableField) (string, error)
	UpdateField(ctx context.Context, appToken string, tableId string, body *larkBitable.AppTableField) error
	DeleteField(ctx context.Context, appToken string, tableId, fieldId string) error
}

func (b *db) ListFields(ctx context.Context, appToken string, tableId string) ([]*larkBitable.AppTableField, error) {
	c := core.WrapContext(ctx)

	reqCall := b.service.AppTableFields.List(c)
	reqCall.SetAppToken(appToken)
	reqCall.SetTableId(tableId)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("ListFields fail! appToken:%s,tableId:%s,response:%s", appToken, tableId, tools.Prettify(message))
		return nil, err
	}
	logrus.WithContext(c).Debugf("response:%s", tools.Prettify(message))
	return message.Items, nil
}

func (b *db) CreateField(ctx context.Context, appToken string, tableId string, body *larkBitable.AppTableField) (string, error) {
	c := core.WrapContext(ctx)

	reqCall := b.service.AppTableFields.Create(c, body)
	reqCall.SetAppToken(appToken)
	reqCall.SetTableId(tableId)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("CreateField fail! appToken:%s,response:%s", appToken, tools.Prettify(message))
		return "", err
	}
	logrus.WithContext(c).Debugf("response:%s", tools.Prettify(message))
	return message.Field.FieldId, nil
}

func (b *db) UpdateField(ctx context.Context, appToken string, tableId string, body *larkBitable.AppTableField) error {
	c := core.WrapContext(ctx)

	reqCall := b.service.AppTableFields.Update(c, body)
	reqCall.SetAppToken(appToken)
	reqCall.SetTableId(tableId)
	reqCall.SetFieldId(body.FieldId)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithFields(logrus.Fields{
			"appToken": appToken,
			"tableId":  tableId,
			"body":     body,
			"response": tools.Prettify(message),
		}).Errorf("UpdateField fail!")
		return err
	}
	logrus.WithContext(ctx).Debugf("UpdateField response:%s", tools.Prettify(message))
	return nil
}

func (b *db) DeleteField(ctx context.Context, appToken string, tableId, fieldId string) error {
	c := core.WrapContext(ctx)

	reqCall := b.service.AppTableFields.Delete(c)
	reqCall.SetAppToken(appToken)
	reqCall.SetTableId(tableId)
	reqCall.SetFieldId(fieldId)
	message, err := reqCall.Do()
	if err != nil {
		logrus.WithContext(c).WithError(err).Errorf("DeleteField fail! appToken:%s,response:%s", appToken, tools.Prettify(message))
		return err
	}
	logrus.WithContext(c).Debugf("response:%s", tools.Prettify(message))
	return nil
}
