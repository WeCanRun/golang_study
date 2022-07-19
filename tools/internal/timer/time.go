package timer

import "time"

const (
	Time_Zone = "Asia/Shanghai"
)

func GetNowTime() time.Time {
	location, _ := time.LoadLocation(Time_Zone)
	return time.Now().In(location)
}

func GetCalculateTime(curTime time.Time, d string) (time.Time, error) {
	duration, err := time.ParseDuration(d)
	if err != nil {
		return time.Time{}, err
	}
	return curTime.Add(duration), nil
}
