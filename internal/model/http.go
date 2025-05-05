package model

type ResponseMessage struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

// Error400 错误请求
// @Description 请求参数错误
type Error400 struct {
	Success    bool   `json:"success" example:"false"`
	ErrCode    string `json:"errCode" example:"400"`
	ErrMessage string `json:"errMessage" example:"请求参数错误"`
}

// Error401 未授权
// @Description 未授权或token无效
type Error401 struct {
	Success    bool   `json:"success" example:"false"`
	ErrCode    string `json:"errCode" example:"401"`
	ErrMessage string `json:"errMessage" example:"未授权或token无效"`
}

// Error403 禁止访问
// @Description 权限不足
type Error403 struct {
	Success    bool   `json:"success" example:"false"`
	ErrCode    string `json:"errCode" example:"403"`
	ErrMessage string `json:"errMessage" example:"权限不足"`
}

// Error404 资源未找到
// @Description 资源未找到
type Error404 struct {
	Success    bool   `json:"success" example:"false"`
	ErrCode    string `json:"errCode" example:"404"`
	ErrMessage string `json:"errMessage" example:"资源未找到"`
}

// Error409 冲突
// @Description 资源冲突
type Error409 struct {
	Success    bool   `json:"success" example:"false"`
	ErrCode    string `json:"errCode" example:"409"`
	ErrMessage string `json:"errMessage" example:"资源冲突"`
}

// Error422 参数校验失败
// @Description 参数校验失败
type Error422 struct {
	Success    bool   `json:"success" example:"false"`
	ErrCode    string `json:"errCode" example:"422"`
	ErrMessage string `json:"errMessage" example:"参数校验失败"`
}

// Error429 请求过多
// @Description 请求过于频繁
type Error429 struct {
	Success    bool   `json:"success" example:"false"`
	ErrCode    string `json:"errCode" example:"429"`
	ErrMessage string `json:"errMessage" example:"请求过于频繁"`
}

// Error500 服务器内部错误
// @Description 服务器内部错误
type Error500 struct {
	Success    bool   `json:"success" example:"false"`
	ErrCode    string `json:"errCode" example:"500"`
	ErrMessage string `json:"errMessage" example:"服务器内部错误"`
}

// Error502 网关错误
// @Description 网关错误
type Error502 struct {
	Success    bool   `json:"success" example:"false"`
	ErrCode    string `json:"errCode" example:"502"`
	ErrMessage string `json:"errMessage" example:"网关错误"`
}

// Error503 服务不可用
// @Description 服务不可用
type Error503 struct {
	Success    bool   `json:"success" example:"false"`
	ErrCode    string `json:"errCode" example:"503"`
	ErrMessage string `json:"errMessage" example:"服务不可用"`
}

// Error504 网关超时
// @Description 网关超时
type Error504 struct {
	Success    bool   `json:"success" example:"false"`
	ErrCode    string `json:"errCode" example:"504"`
	ErrMessage string `json:"errMessage" example:"网关超时"`
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
