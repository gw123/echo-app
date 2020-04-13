package services

import (
	"encoding/json"
	echoapp "github.com/gw123/echo-app"
	"github.com/pkg/errors"
	"io/ioutil"
	"sort"
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

func (aSvr *AreaService) GetAreaMap(areaId string) (map[string]string, error) {
	areaArray, err := aSvr.GetAreaArray(areaId)
	if err != nil {
		return nil, err
	}
	areaMap := make(map[string]string)
	for _, item := range areaArray {
		areaMap[item.Id] = item.Name
	}

	return areaMap, nil
}

func (aSvr *AreaService) GetAreaArray(areaId string) ([]echoapp.AreaOption, error) {
	length := len(areaId)
	areaList := make(echoapp.AreaOptionArray, 0)
	if length == 0 {
		for key, p := range aSvr.areaMap {
			areaList = append(areaList, echoapp.AreaOption{
				Id:   key,
				Name: p.Name,
			})
		}
	} else if length == 2 {
		cityList, ok := aSvr.areaMap[areaId]
		if !ok {
			return nil, errors.New("参数错误")
		}
		for index, item := range cityList.Children {
			areaList = append(areaList, echoapp.AreaOption{
				Id:   index,
				Name: item.Name,
			})
		}
	} else if length == 4 {
		cityList, ok := aSvr.areaMap[areaId[0:2]]
		if !ok {
			return nil, errors.New("参数错误")
		}
		areaList1, ok := cityList.Children[areaId]
		if !ok {
			return nil, errors.New("参数错误")
		}
		for index, item := range areaList1.Children {
			areaList = append(areaList, echoapp.AreaOption{
				Id:   index,
				Name: item.Name,
			})
		}
	} else if length == 6 {
		data, err := ioutil.ReadFile(aSvr.areasPath + "/town/" + areaId + ".json")
		if err != nil {
			return nil, errors.Wrap(err, "解析失败")
		}
		areaMap := make(map[string]string)
		err = json.Unmarshal(data, &areaMap)
		if err != nil {
			return nil, errors.Wrap(err, "解析失败")
		}
		for key, name := range areaMap {
			areaList = append(areaList, echoapp.AreaOption{
				Id:   key,
				Name: name,
			})
		}
	}
	sort.Sort(areaList)
	return areaList, nil
}
