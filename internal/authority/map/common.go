package _map

const (
	DefaultPage     = 1
	DefaultPageSize = 10
)

var DefaultPageRequest = PageList{
	Page:     DefaultPage,
	PageSize: DefaultPageSize,
}

type PageList struct {
	Page     uint32 `json:"page" form:"page" validate:"required,number" label:"页码"`
	PageSize uint32 `json:"page_size" form:"page_size" validate:"required,number" label:"页码大小"`
	Keyword  string `json:"keyword" form:"keyword"`
	IsDelete bool   `json:"is_delete" form:"is_delete"`
}

type Id struct {
	Id uint `validate:"required,number"`
}

type Ids struct {
	List []uint32 `validate:"required"`
}
