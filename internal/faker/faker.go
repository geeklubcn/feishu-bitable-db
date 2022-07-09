package faker

import "os"

var AppID string
var AppSecret string

func init() {
	AppID = os.Getenv("app_id")
	AppSecret = os.Getenv("app_secret")
}
