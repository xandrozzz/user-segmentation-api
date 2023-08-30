package models

import "time"

type User struct {
	ID  int `json:"id"`
	TTL int `json:"ttl"`
}

type Segment struct {
	Name       string `json:"name"`
	Percentage int    `json:"percentage"`
}

type NameReq struct {
	Name string `json:"name"`
}

type IDReq struct {
	ID int `json:"id"`
}

type ChangeReq struct {
	Add    []string `json:"add"`
	Remove []string `json:"remove"`
}

type StatsReq struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

type StatString struct {
	UserID      int       `json:"id"`
	SegmentName string    `json:"name"`
	Operation   string    `json:"operation"`
	Timestamp   time.Time `json:"timestamp"`
}

type SegmentCounter struct {
	Segment    Segment
	Proportion int
	Count      int
}

func NewCounter(segment Segment) SegmentCounter {
	var counter SegmentCounter
	counter.Count = 0
	counter.Segment = segment
	counter.Proportion = segment.Percentage
	return counter
}
