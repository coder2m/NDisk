/**
 * @Author: yangon
 * @Description
 * @Date: 2021/1/15 17:35
 **/
package _map

type (
	Target struct {
		To string `validate:"required"`
	}

	Batch struct {
		To      string   `validate:"required"`
		Operate []string `validate:"required"`
	}

	Single struct {
		To      string `validate:"required"`
		Operate string `validate:"required"`
	}

	Resources struct {
		Role   string `validate:"required"`
		Obj    string `validate:"required"`
		Action string `validate:"required"`
	}

	Array struct {
		Data []string `validate:"required"`
	}

	RolesReq struct {
		Name        string `validate:"required"`
		Description string `validate:"required"`
	}

	ResourcesReq struct {
		Name        string `validate:"required"`
		Path        string `validate:"required"`
		Action      string `validate:"required"`
		Description string `validate:"required"`
	}

	MenuReq struct {
		ParentId    uint   `validate:"required"`
		Path        string `validate:"required"`
		Name        string `validate:"required"`
		Description string `validate:"required"`
		IconClass   string `validate:"required"`
	}

	UpdateRolesMenuAndResourcesReq struct {
		ID        uint   `validate:"required"`
		Menus     []uint `validate:"required"`
		Resources []uint `validate:"required"`
	}

	//------------------------------
	MenuResInfo struct {
		ID          uint   `json:"id,omitempty"`
		ParentId    uint   `json:"parent_id,omitempty"`
		Path        string `json:"path,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		IconClass   string `json:"icon_class,omitempty"`
		CreatedAt   uint64 `json:"created_at,omitempty"`
		UpdatedAt   uint64 `json:"updated_at,omitempty"`
		DeletedAt   uint64 `json:"deleted_at,omitempty"`
	}

	ResourcesResInfo struct {
		ID          uint   `json:"id,omitempty"`
		Name        string `json:"name,omitempty"`
		Path        string `json:"path,omitempty"`
		Action      string `json:"action,omitempty"`
		Description string `json:"description,omitempty"`
		CreatedAt   uint64 `json:"created_at,omitempty"`
		UpdatedAt   uint64 `json:"updated_at,omitempty"`
		DeletedAt   uint64 `json:"deleted_at,omitempty"`
	}

	RolesResInfo struct {
		ID          uint               `json:"id,omitempty"`
		Name        string             `json:"name,omitempty"`
		Description string             `json:"description,omitempty"`
		Menus       []MenuResInfo      `json:"menus,omitempty"`
		Resources   []ResourcesResInfo `json:"resources,omitempty"`
		CreatedAt   uint64             `json:"created_at,omitempty"`
		UpdatedAt   uint64             `json:"updated_at,omitempty"`
		DeletedAt   uint64             `json:"deleted_at,omitempty"`
	}
)
