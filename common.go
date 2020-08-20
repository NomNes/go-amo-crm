package amo

import "time"

type IdNameItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type IdLink struct {
	Id int `json:"id"`
}

type EntityTime struct {
	CreatedAt int64 `json:"created_at"`
	UpdatedAt int64 `json:"updated_at"`
}

func (et *EntityTime) SetCreatedAtTime(time time.Time) {
	et.CreatedAt = time.Unix()
}

func (et *EntityTime) GetCreatedAtTime() time.Time {
	return time.Unix(et.CreatedAt, 0)
}

func (et *EntityTime) SetUpdatedAtTime(time time.Time) {
	et.UpdatedAt = time.Unix()
}

func (et *EntityTime) GetUpdatedAtTime() time.Time {
	return time.Unix(et.UpdatedAt, 0)
}
