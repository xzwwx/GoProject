package common

import (
	"strconv"
	"time"
)

func GenerateKey(uid uint64) string {
	return time.Now().String() + strconv.FormatInt(int64(uid), 10) + "ppt"
}
