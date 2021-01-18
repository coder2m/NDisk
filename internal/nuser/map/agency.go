package _map

type (
	AgencyReq struct {
		ParentId uint   //父id
		Name     string //名字
		Remark   string //描述信息
	}

	CreateManyAgencyReq struct {
		Uid    uint32
		Agency []AgencyReq
	}

	UpdateAgency struct {
		ID       uint   `json:"id"`
		ParentId uint   `json:"parent_id"`
		Name     string `json:"name"`
		Remark   string `json:"remark"`
	}

	//----------------------------
	AgencyInf struct {
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
