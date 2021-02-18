package _map

type (
	UploadStartReq struct {
		Name string `json:"name" xml:"name" form:"name"`
		Size uint64 `json:"size" xml:"size" form:"size"`
		Type string `json:"type" xml:"type" form:"type"`
	}
	UploadStartRes struct {
		FileId    uint   `json:"file_id" xml:"file_id" form:"file_id"`
		SecTrans  bool   `json:"sec_trans" xml:"sec_trans" form:"sec_trans"`
		BlockSize uint64 `json:"block_size" xml:"block_size" form:"block_size"`
		SliceSize uint64 `json:"slice_size" xml:"slice_size" form:"slice_size"`
	}
)
