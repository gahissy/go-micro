package schemas

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type FileInfo struct {
	Id   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty" validate:"required"`
	Mime string `json:"mime,omitempty"`
}

func (a *FileInfo) Value() (driver.Value, error) {
	return json.Marshal(a)
}

type FileInfoList struct {
	Data []FileInfo `json:"-"`
}

func (f *FileInfoList) Value() (driver.Value, error) {
	return json.Marshal(f.Data)
}

func (f *FileInfoList) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Data)
}

func (f *FileInfoList) Scan(src interface{}) error {
	if src == nil {
		return nil
	}
	switch src := src.(type) {
	case []byte:
		return json.Unmarshal(src, &f.Data)
	case string:
		return json.Unmarshal([]byte(src), &f.Data)
	default:
		return fmt.Errorf("cannot scan type %T into MyData", src)
	}
}
