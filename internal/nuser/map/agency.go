package _map

type (
	AgencyReq struct {
		ParentId uint   `validate:"required"` //父id
		Name     string `validate:"required"` //名字
		Remark   string `validate:"required"` //描述信息
	}

	CreateManyAgencyReq struct {
		Uid    uint32      `validate:"required"`
		Agency []AgencyReq `validate:"required"`
	}

	UpdateAgency struct {
		ID       uint   `json:"id" validate:"required"`
		ParentId uint   `json:"parent_id validate:"required""`
		Name     string `json:"name" validate:"required"`
		Remark   string `json:"remark" validate:"required"`
	}

	//----------------------------
	AgencyInf struct {
		AUId       uint     `json:"auid,omitempty"` //关联表id
		ID         uint     `json:"id,omitempty"`
		ParentId   uint     `json:"parent_id,omitempty"`
		Name       string   `json:"name,omitempty"`
		Remark     string   `json:"remark,omitempty"`
		Status     uint     `json:"status,omitempty"`
		CreateUser UserInfo `json:"create_user,omitempty"`
		CreatedAt  int64    `json:"created_at,omitempty"`
		UpdatedAt  int64    `json:"updated_at,omitempty"`
		DeletedAt  int64    `json:"deleted_at,omitempty"`
	}
)
