package models

type User struct {
	ID int `json:"id"`
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
