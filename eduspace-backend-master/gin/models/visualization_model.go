package models

import (
	"encoding/json"
	"errors"
)

// CustomSubmitCodeForm
// 提交自定义代码的表单
// 参数:code
type CustomSubmitCodeForm struct {
	Code string `json:"code" binding:"required"`
}

func (r *CustomSubmitCodeForm) UnmarshalJSON(data []byte) error {
	required := struct {
		Code string `json:"code"`
	}{}

	err := json.Unmarshal(data, &required)
	if err != nil {
		return err
	} else if required.Code == "" {
		err = errors.New("缺少必填字段code")
	} else {
		r.Code = required.Code
	}

	return nil
}
