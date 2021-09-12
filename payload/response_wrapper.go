package payload

type RespWrap struct {
	Data  interface{} `json:"data" extensions:"x-nullable"`
	Error interface{} `json:"error" extensions:"x-nullable"`
}
