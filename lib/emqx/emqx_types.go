package emqx

import "log"

type EMQX struct {
	logger  *log.Logger
	emqxUrl string
	apiName string
	apiPass string
}

type EmqxExportDataResp struct {
	Node     string `json:"node"`
	Filename string `json:"filename"`
}

type EmqxErrorResp struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
