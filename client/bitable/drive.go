package bitable

import (
	"fmt"

	"github.com/larksuite/oapi-sdk-go/api"
	"github.com/larksuite/oapi-sdk-go/api/core/request"
	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
)

type DriveExt interface {
	GetDriveFiles(ctx *core.Context, folderToken string) (*DriveFiles, error)
}

type driveExt struct {
	conf *config.Config
}

func NewDriveExt(conf *config.Config) DriveExt {
	return &driveExt{conf: conf}
}

func (d *driveExt) GetDriveFiles(ctx *core.Context, folderToken string) (*DriveFiles, error) {
	var result = &DriveFiles{}
	req := request.NewRequest(fmt.Sprintf("/open-apis/drive/v1/files?folder_token=%s", folderToken), "GET",
		[]request.AccessTokenType{request.AccessTokenTypeTenant, request.AccessTokenTypeUser}, nil, result)
	err := api.Send(ctx, d.conf, req)
	return result, err
}

type DriveFiles struct {
	HasMore   bool         `json:"has_more,omitempty"`
	PageToken string       `json:"page_token,omitempty"`
	Total     int          `json:"total,omitempty"`
	Files     []*DriveFile `json:"files,omitempty"`
}

type DriveFile struct {
	Name        string `json:"name"`
	ParentToken string `json:"parent_token"`
	Token       string `json:"token"`
	Type        string `json:"type"`
	Url         string `json:"url"`
}
