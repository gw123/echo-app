package services

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
)

type Area struct {
	Name     string          `json:"n"`
	Children map[string]Area `json:"c"`
}

type AreaService struct {
	areaMap   map[string]Area
	areasPath string
}

func NewAreaService(araesPath string) (*AreaService, error) {
	areaFormatObjFilePath := araesPath + "/area_format_object.json"
	data, err := ioutil.ReadFile(areaFormatObjFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "newAreaService")
	}
	areaMap := make(map[string]Area)
	err = json.Unmarshal(data, &areaMap)
	if err != nil {
		return nil, errors.Wrap(err, "newAreaService")
	}
	return &AreaService{areaMap: areaMap, areasPath: araesPath}, nil
}

func (t *AreaService) GetAreaList(areaId string) (map[string]string, error) {
	length := len(areaId)
	areaList := make(map[string]string)
	if length == 0 {
		for key, p := range t.areaMap {
			areaList[key] = p.Name
		}
	} else if length == 2 {
		cityList, ok := t.areaMap[areaId]
		if !ok {
			return nil, errors.New("参数错误")
		}
		for index, item := range cityList.Children {
			areaList[index] = item.Name
		}
	} else if length == 4 {
		cityList, ok := t.areaMap[areaId[0:2]]
		if !ok {
			return nil, errors.New("参数错误")
		}
		areaList1, ok := cityList.Children[areaId]
		if !ok {
			return nil, errors.New("参数错误")
		}
		for index, item := range areaList1.Children {
			areaList[index] = item.Name
		}
	} else if length == 6 {
		data, err := ioutil.ReadFile(t.areasPath + "/town/" + areaId + ".json")
		if err != nil {
			return nil, errors.Wrap(err, "解析失败")
		}
		err = json.Unmarshal(data, &areaList)
		if err != nil {
			return nil, errors.Wrap(err, "解析失败")
		}
	}
	return areaList, nil
}
