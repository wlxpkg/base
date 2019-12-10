/*
 * @Author: qiuling
 * @Date: 2019-09-05 16:01:26
 * @Last Modified by: qiuling
 * @Last Modified time: 2019-12-05 11:10:17
 */
package model

import (
	. "github.com/wlxpkg/base"
)

type AdminOperationLog struct {
	ID        uint      `json:"id"`
	UserId    string    `json:"user_id"`
	Path      string    `json:"path"`
	Method    string    `json:"method"`
	Ip        string    `json:"ip"`
	Request   []byte    `json:"request"`
	CreatedAt JSONTime  `json:"created_at"`
	UpdatedAt JSONTime  `json:"updated_at"`
	DeletedAt *JSONTime `json:"-"`
}

func (AdminOperationLog) TableName() string {
	return "user_admin_operation_log"
}
