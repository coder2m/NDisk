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
)
