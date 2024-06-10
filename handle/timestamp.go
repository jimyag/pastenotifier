package handle

import (
	"strconv"
	"time"
)

type Timestamp struct{}

func (*Timestamp) Handle(content string) (title string, message string, err error) {
	if content == "" {
		return "", "", nil
	}

	if len(content) < 10 {
		return "", "", nil
	}

	timeStamp, err := strconv.Atoi(content)
	if err != nil {
		return "", "", nil
	}
	var ti time.Time

	switch len(content) {
	// 秒
	case 10:
		ti = time.Unix(int64(timeStamp), 0)
	// 毫秒
	case 13:
		ti = time.UnixMilli(int64(timeStamp))
	// 微秒
	case 16:
		ti = time.UnixMicro(int64(timeStamp))
	// 百纳秒
	case 17:
		ti = time.Unix(0, int64(timeStamp)*100)
	// 纳秒
	case 19:
		ti = time.Unix(0, int64(timeStamp))
	default:
		return "", "", nil
	}
	title = content
	message = ti.Format("2006-01-02 15:04:05")
	return title, message, nil
}
