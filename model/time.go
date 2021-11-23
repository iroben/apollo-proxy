package model

import (
	"database/sql/driver"
	"errors"
	"time"
)

const (
	YMD    = "2006-01-02"
	YMDHI  = "2006-01-02 15:04"
	YMDHIS = "2006-01-02 15:04:05"
)

// 年月日
type TimeYmd time.Time

func (t TimeYmd) MarshalJSON() ([]byte, error) {
	_time := time.Time(t)
	return []byte(`"` + _time.Format(YMD) + `"`), nil
}

func (t *TimeYmd) UnmarshalJSON(data []byte) error {
	timeStr := string(data[1 : len(data)-1])
	_t, _err := time.Parse(YMD, timeStr)
	*t = TimeYmd(_t)
	return _err
}
func (t *TimeYmd) Scan(v interface{}) error {
	_t, ok := v.(time.Time)
	if ok {
		*t = TimeYmd(_t)
		return nil
	}
	return errors.New("fail")
}

func (t TimeYmd) Value() (driver.Value, error) {
	_t := time.Time(t)
	return _t.Format(YMD), nil
}

// 年月日时分
type TimeYmdHi time.Time

func (t TimeYmdHi) MarshalJSON() ([]byte, error) {
	_time := time.Time(t)
	return []byte(`"` + _time.Format(YMDHI) + `"`), nil
}

func (t *TimeYmdHi) UnmarshalJSON(data []byte) error {
	timeStr := string(data[1 : len(data)-1])
	_t, _err := time.Parse(YMDHI, timeStr)
	*t = TimeYmdHi(_t)
	return _err
}
func (t TimeYmdHi) Value() (driver.Value, error) {
	_t := time.Time(t)
	return _t.Format(YMDHI), nil
}

func (t *TimeYmdHi) Scan(v interface{}) error {
	_t, ok := v.(time.Time)
	if ok {
		*t = TimeYmdHi(_t)
		return nil
	}
	return errors.New("fail")
}

// 年月日时分秒
type TimeYmdHis time.Time

func (t *TimeYmdHis) MarshalJSON() ([]byte, error) {
	_time := time.Time(*t)
	return []byte(`"` + _time.Format(YMDHIS) + `"`), nil
}

func (t *TimeYmdHis) UnmarshalJSON(data []byte) error {
	timeStr := string(data[1 : len(data)-1])
	_t, _err := time.Parse(YMDHIS, timeStr)
	*t = TimeYmdHis(_t)
	return _err
}

func (t TimeYmdHis) Value() (driver.Value, error) {
	_t := time.Time(t)
	return _t.Format(YMDHI), nil
}

func (t *TimeYmdHis) Scan(v interface{}) error {
	_t, ok := v.(time.Time)
	if ok {
		*t = TimeYmdHis(_t)
		return nil
	}
	return errors.New("fail")
}
