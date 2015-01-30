package db

import ()

const ()

type (
	Crowdfund struct {
		Id         int64    `json:"id"`
		OwnerId    int64    `json:"owner"`
		CreateTime int64    `json:"ctime"`
		Goal       int64    `json:"goal"`
		Сollected  int64    `json:"сollected"`
		Name       string   `json:"name"`
		Cover      string   `json:"cover"`
		Images     []string `json:"images"`
	}
)
