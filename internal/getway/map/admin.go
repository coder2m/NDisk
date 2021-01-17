package _map

type (
	// 请求
	User struct {
		Name     string `validate:"required" json:"name"`
		Alias    string `validate:"required" json:"alias"`
		Tel      string `validate:"required,phone" json:"tel"`
		Email    string `validate:"required" json:"email"`
		Password string `validate:"required" json:"password,email,min=8"`
	}

	UpdateUser struct {
		Uid        uint64 `validate:"required,number" json:"uid"`
		Name       string `validate:"required" json:"name"`
		Alias      string `validate:"required" json:"alias"`
		Tel        string `validate:"required,phone" json:"tel"`
		Email      string `validate:"required,email" json:"email"`
		Password   string `validate:"max=20,min=8" json:"password"`
		RePassword string `validate:"max=20,min=8,eqfield=Password" json:"re_password"`
	}

	UidList struct {
		List []uint32 `validate:"required" json:"list"`
	}

	CreateUser struct {
		Data []User `validate:"required" json:"data"`
	}

	UpdateStatus struct {
		Uid    uint64 `validate:"required,number" json:"uid"`
		Status uint32 `validate:"required,number" json:"status"`
	}
	Uid struct {
		Uid uint64 `uri:"uid" json:"uid" validate:"required,number,min=1" label:"uid"`
	}

	CompetenceReq struct {
		Objective string `json:"objective"  validate:"required"`
		Action    string `json:"action" validate:"required"`
	}

	RoleReq struct {
		Content string `json:"content" validate:"required"`
	}

	RoleCompetenceReq struct {
		Role      string `json:"role" validate:"required"`
		Objective string `json:"objective"  validate:"required"`
		Action    string `json:"action" validate:"required"`
	}

	UserRolesReq struct {
		Uid  uint64   `uri:"uid" json:"uid" validate:"required,number,min=1" label:"uid"`
		Role []string `json:"role" validate:"required"`
	}

	UserRoleReq struct {
		Uid  uint64 `uri:"uid" json:"uid" validate:"required,number,min=1" label:"uid"`
		Role string `json:"role" validate:"required"`
	}

	// 返回
	Batch struct {
		Count uint32 `json:"count,omitempty"`
	}

	UserList struct {
		Count uint32     `json:"count,omitempty"`
		Data  []UserInfo `json:"data,omitempty"`
	}

	Competence struct {
		Objective string `json:"objective,omitempty"`
		Action    string `json:"action,omitempty"`
	}

	CompetenceList struct {
		Data []Competence `json:"data,omitempty"`
	}
)
