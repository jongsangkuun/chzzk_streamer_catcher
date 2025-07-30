package chzzk

import "fmt"

const BaseURL = "https://openapi.chzzk.naver.com/"

var ChannelUrl = fmt.Sprintf("%s/open/v1/channels", BaseURL)

var CategoriyUrl = fmt.Sprintf("%s/open/v1/categories", BaseURL)
