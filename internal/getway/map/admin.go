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

	RoleInfoReq struct {
		Name        string ` json:"name" validate:"required"`
		Description string ` json:"description" validate:"required"`
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

	MenuReq struct {
		ParentId    uint32 `json:"parent_id" validate:"required"`
		Path        string `json:"path" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description" validate:"required"`
		IconClass   string `json:"icon_class" validate:"required"`
	}

	ResourcesReq struct {
		Name        string `json:"name" validate:"required"`
		Path        string `json:"path" validate:"required"`
		Action      string `json:"action" validate:"required"`
		Description string `json:"description" validate:"required"`
	}

	UpdateRolesMenuAndResourcesReq struct {
		Id        uint32   ` json:"id" validate:"required"`
		Menus     []uint32 ` json:"menus" validate:"required"`
		Resources []uint32 `json:"resources" validate:"required"`
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

	RolesListRes struct {
		Count uint32         `json:"count,omitempty"`
		Data  []RolesInfoRes `json:"data,omitempty"`
	}

	MenuListRes struct {
		Count uint32        `json:"count,omitempty"`
		Data  []MenuInfoRes `json:"data,omitempty"`
	}

	ResourcesListRes struct {
		Count uint32             `json:"count,omitempty"`
		Data  []ResourcesInfoRes `json:"data,omitempty"`
	}

	RolesInfoRes struct {
		Id          uint32             `json:"id,omitempty"`
		Name        string             ` json:"name,omitempty"`
		Description string             ` json:"description,omitempty"`
		Menus       []MenuInfoRes      ` json:"menus,omitempty"`
		Resources   []ResourcesInfoRes ` json:"resources,omitempty"`
		CreatedAt   uint64             ` json:"created_at,omitempty"`
		UpdatedAt   uint64             ` json:"updated_at,omitempty"`
		DeletedAt   uint64             ` json:"deleted_at,omitempty"`
	}

	MenuInfoRes struct {
		Id          uint32 `json:"id,omitempty"`
		ParentId    uint32 `json:"parent_id,omitempty"`
		Path        string `json:"path,omitempty"`
		Name        string ` json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		IconClass   string `json:"icon_class,omitempty"`
		CreatedAt   uint64 `json:"created_at,omitempty"`
		UpdatedAt   uint64 ` json:"updated_at,omitempty"`
		DeletedAt   uint64 `json:"deleted_at,omitempty"`
	}

	ResourcesInfoRes struct {
		Id          uint32 `json:"id,omitempty"`
		Name        string `json:"name,omitempty"`
		Path        string `json:"path,omitempty"`
		Action      string `json:"action,omitempty"`
		Description string `json:"description,omitempty"`
		CreatedAt   uint64 `json:"created_at,omitempty"`
		UpdatedAt   uint64 `json:"updated_at,omitempty"`
		DeletedAt   uint64 `json:"deleted_at,omitempty"`
	}

	CompetenceList struct {
		Data []Competence `json:"data,omitempty"`
	}
)
