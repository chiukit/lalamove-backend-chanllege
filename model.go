package main

import "encoding/json"

type Path [][]string

type Record struct {
	Status        string `json:"status"`
	Path          Path   `json:"path,omitempty"`
	TotalDistance int    `json:"total_distance,omitempty"`
	TotalTime     int    `json:"total_time,omitempty"`
}

func (r Record) MarshalBinary() ([]byte, error) {
	return json.Marshal(r)
}
