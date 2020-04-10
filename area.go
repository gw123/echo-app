package echoapp

type AreaService interface {
	GetAreaList(areaId string) (map[string]string, error)
}
