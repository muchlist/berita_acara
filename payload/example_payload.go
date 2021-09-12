package payload

// swag init -g app/app.go

type ErrorExample400 struct {
	Status  int      `json:"status" example:"401"`
	Message string   `json:"message" example:"Unauthorized, memerlukan hak akses [ADMIN]"`
	Error   string   `json:"error" example:"unauthorized"`
	Causes  []string `json:"causes" example:"causes 1,causes 2"`
}

type ErrorExample500 struct {
	Status  int      `json:"status" example:"500"`
	Message string   `json:"message" example:"gagal saat penghapusan item"`
	Error   string   `json:"error" example:"internal_server_error"`
	Causes  []string `json:"causes" example:"ERROR: argument of WHERE must be type boolean. not type integer (SQLSTATE 42804)"`
}

type RespMsgExample struct {
	Data  string      `json:"data" extensions:"x-nullable" example:"Data dengan ID xxx berhasil di [Create/Delete]"`
	Error interface{} `json:"error" extensions:"x-nullable"`
}

type RespFileExample struct {
	Data  string      `json:"data" extensions:"x-nullable" example:"static/image/example.jpg"`
	Error interface{} `json:"error" extensions:"x-nullable"`
}
