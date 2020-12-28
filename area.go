package echoapp

import "strings"

type AreaOption struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type AreaOptionArray []AreaOption


func (arr AreaOptionArray) Len() int {
	return len(arr)
}

func (arr AreaOptionArray) Swap(i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

func (arr AreaOptionArray) Less(i, j int) bool {
	return strings.Compare(arr[i].Id, arr[j].Id) == -1
}

type AreaService interface {
	GetAreaMap(areaId string) (map[string]string, error)
	GetAreaArray(areaId string) ([]AreaOption, error)
}
