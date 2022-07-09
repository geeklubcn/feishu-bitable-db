package bitable

import (
	"context"

	"github.com/larksuite/oapi-sdk-go/core"
	"github.com/larksuite/oapi-sdk-go/core/config"
	larkDriveExplorer "github.com/larksuite/oapi-sdk-go/service/drive_explorer/v2"
)

const Type = "bitable"

type Bitable interface {
	CreateApp(ctx context.Context, name, folderToken string) (string, error)
	QueryByName(ctx context.Context, name, folderToken string) (string, bool)
}

type bitable struct {
	drives        DriveExt
	driveExplorer *larkDriveExplorer.Service
}

func NewBitable(conf *config.Config) Bitable {
	return &bitable{
		drives:        NewDriveExt(conf),
		driveExplorer: larkDriveExplorer.NewService(conf),
	}
}

func (b *bitable) CreateApp(ctx context.Context, title, folderToken string) (string, error) {
	c := core.WrapContext(ctx)

	req := &larkDriveExplorer.FileCreateReqBody{
		Title: title,
		Type:  Type,
	}
	caller := b.driveExplorer.Files.Create(c, req)
	caller.SetFolderToken(folderToken)
	resp, err := caller.Do()
	if err != nil {
		return "", err
	}
	return resp.Token, nil
}

func (b *bitable) QueryByName(ctx context.Context, name, folderToken string) (string, bool) {
	c := core.WrapContext(ctx)

	dfs, err := b.drives.GetDriveFiles(c, folderToken)
	if err != nil {
		return "", false
	}
	for _, df := range dfs.Files {
		if df.Name == name && df.Type == Type {
			return df.Token, true
		}
	}
	return "", false
}
