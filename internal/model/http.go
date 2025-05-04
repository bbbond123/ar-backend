package model

type ResponseMessage struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Error404 struct {
}

type Error400 struct {
}

// BaseResponse 通用返回（无data时使用）
type BaseResponse struct {
	Success    bool   `json:"success"`              // 请求是否成功
	ErrCode    string `json:"errCode,omitempty"`    // 错误码
	ErrMessage string `json:"errMessage,omitempty"` // 错误信息
}

// Response 通用返回（带data时使用）
type Response[T any] struct {
	Data       T      `json:"data"`                 // 数据
	Success    bool   `json:"success"`              // 请求是否成功
	ErrCode    string `json:"errCode,omitempty"`    // 错误码
	ErrMessage string `json:"errMessage,omitempty"` // 错误信息
}

// ListResponse 列表返回
type ListResponse[T any] struct {
	Total      int64  `json:"total"`                // 总条数
	List       []T    `json:"list"`                 // 列表
	Success    bool   `json:"success"`              // 请求是否成功
	ErrCode    string `json:"errCode,omitempty"`    // 错误码
	ErrMessage string `json:"errMessage,omitempty"` // 错误信息
}
