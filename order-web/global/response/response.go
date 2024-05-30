package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

// 实现MarshalJSON方法可以自动将此类型转换为json
func (j JsonTime) MarshalJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2006-01-02"))
	return []byte(stmp), nil
}

type GoodsResponse struct {
	Id       int32    `json:"id"`
	NickName string   `json:"name"`
	Birthday JsonTime `json:"birthday"`
	//Birthday time.Time `json:"birthday"`
	//Birthday string `json:"birthday"`
	Gender string `json:"gender"`
	Mobile string `json:"mobile"`
}
