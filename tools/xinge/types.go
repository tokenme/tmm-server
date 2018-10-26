package xinge

const (
	DOMAIN           = "openapi.xg.qq.com"
	VERSION          = "v2"
	DefaultValidTime = 600
)

type CalType = int

const (
	OfflineCalType  CalType = 0
	RealtimeCalType CalType = 1
)

type HttpMethod = string

const (
	GET  HttpMethod = "GET"
	POST HttpMethod = "POST"
)

const (
	AND = "AND"
	OR  = "OR"
)
