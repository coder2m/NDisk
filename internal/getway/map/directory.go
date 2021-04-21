package _map

type (
	DirList struct {
		Count uint32          `json:"count,omitempty"`
		Data  []DirectoryInfo `json:"data,omitempty"`
	}

	DirectoryInfo struct {
		Id        uint   `json:"id,omitempty"`
		Uid       uint   `json:"uid,omitempty"`
		FileId    uint   `json:"file_id,omitempty"`
		IsDir     bool   `json:"is_dir"`
		Name      string `json:"name,omitempty"`
		ParentID  uint   `json:"parent_id"`
		Type      uint   `json:"type"`
		CreatedAt int64  `json:"created_at,omitempty"`
		UpdatedAt int64  `json:"updated_at,omitempty"`
		DeletedAt int64  `json:"deleted_at,omitempty"`
	}

	DirectoryPost struct {
		FileId   uint   `json:"file_id"`
		IsDir    bool   `validate:"required" json:"is_dir"`
		Name     string `validate:"required" json:"name"`
		ParentID uint   `json:"parent_id"`
		Uid      uint
		Type     uint `validate:"required,number" json:"type"`
	}

	DirectoryUpdate struct {
		Id       uint
		Name     string `validate:"required" json:"name"`
		ParentID uint   `json:"parent_id"`
		Uid      uint
	}
)
