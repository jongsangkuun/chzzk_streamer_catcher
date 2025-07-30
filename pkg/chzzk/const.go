package chzzk

import "fmt"

const (
	BaseURL         = "https://openapi.chzzk.naver.com/open"
	ChannelVersion  = "v1"
	CategoryVersion = "v1"
)

var ChannelUrl = fmt.Sprintf("%s/%s/channels", BaseURL, ChannelVersion)

var CategoryUrl = fmt.Sprintf("%s/%s/categories", BaseURL, CategoryVersion)
