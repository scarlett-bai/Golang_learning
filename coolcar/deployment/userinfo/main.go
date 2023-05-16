package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

func EncryptedMD5(file []byte) string {
	hash := md5.New()
	hash.Write(file)
	return hex.EncodeToString(hash.Sum(nil))
}

func GetUserWsCf(userID string) string {
	return EncryptedMD5([]byte("itm" + userID + "1a9b3c"))
}

const (
	// Default 默认时间格式
	Default = "2006-01-02 15:04:05"
	// Hour 只精确到小时
	Hour = "2006-01-02 15"
	// Day 只精确到天
	Day = "2006-01-02"
	// ShortDay 短日期格式
	ShortDay = "20060102"
	// ShortDateTime 短时间格式
	ShortDateTime = "2006-01-02 15:04"
	// Time 只有时间
	Time = "15:04:05"
	// String time.String() 方法格式
	String = "2006-01-02 15:04:05.999999999 -0700 MST"
)

func GetTimeFormat(t time.Time) string {
	ti, _ := time.Parse(Time, t.String())
	return ti.String()
}

func main() {
	fmt.Println(GetTimeFormat(time.Now()))
}
