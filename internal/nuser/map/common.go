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
	Page     int    `json:"page" form:"page" validate:"required,number" label:"页码"`
	PageSize int    `json:"page_size" form:"page_size" validate:"required,number" label:"页码大小"`
	Keyword  string `json:"keyword" form:"keyword"`
	IsDelete bool   `json:"is_delete" form:"is_delete"`
}

type IdMap struct {
	Id uint `uri:"id" json:"id" validate:"required,number,min=1" label:"id"`
}
