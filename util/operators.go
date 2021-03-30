package echoapp_util

import (
	"math"
)

// const (
// 	a = 1.0
// 	b = 2.0
// )

// func WFGHM(a, b float64, k int, theta1, theta2, theta3, theta4 []float64) (float64, error) {
// 	var Ptheta = [][]float64{
// 		theta1, theta2, theta3, theta4,
// 	}
// 	wight := []float64{0.25, 0.1, 0.3, 0.35}
// 	temp0 := 1.0
// 	q := 1.0
// 	for i, thetai := range Ptheta {
// 		for j, thetaj := range Ptheta {
// 			temp0 = temp0 * math.Pow((1-math.Pow((1-math.Pow(thetai[k], q)), a)*math.Pow((1-math.Pow(thetaj[k], q)), b)), wight[j]*wight[i])
// 		}
// 	}
// 	tmp := (math.Pow((1 - math.Pow((1-temp0), 1/(a+b))), 1/q))
// 	return tmp, nil
// }

func WFGHM(a, b float64, Ptheta, wight []float64) (float64, error) {
	// var Ptheta = []float64{
	// 	theta1, theta2, theta3, theta4,
	// }
	//wight := []float64{0.25,1, 0.3, 0.35}
	temp0 := 1.0
	q := 1.0
	for i, thetai := range Ptheta {
		for j, thetaj := range Ptheta {
			temp0 = temp0 * math.Pow((1-math.Pow((1-math.Pow(thetai, q)), a)*math.Pow((1-math.Pow(thetaj, q)), b)), wight[j]*wight[i])
		}
	}
	tmp := (math.Pow((1 - math.Pow((1-temp0), 1/(a+b))), 1/q))
	// if tmp >= 1 || tmp <= 0 {
	// 	return -1, errors.New("wfghmop error")
	// }
	return tmp, nil
}

// func main() {
// 	theta4 := [][]float64{ //5件商品，4个属性，一个用户
// 		{0.5, 0.5, 0.3, 0.4},
// 		{0.7, 0.7, 0.6, 0.6},
// 		{0.5, 0.6, 0.6, 0.6},
// 		{0.8, 0.7, 0.4, 0.5},
// 		{0.4, 0.4, 0.4, 0.4},
// 	}
// 	wight := []float64{0.25, 0.1, 0.3, 0.35}
// 	res := make([]float64, 0)
// 	for i := 0; i < len(theta4); i++ {
// 		temp, _ := WFGHM(a, b, theta4[i], wight)
// 		res = append(res, temp)
// 	}
// 	fmt.Println(res)
// }

// func LinguisticToTFS(value int) []float64 {
// 	if value == 0 {
// 		return []float64{0, 0, 0.25}
// 	} else if value == 1 {
// 		return []float64{0, 0.25, 0.5}
// 	} else if value == 2 {
// 		return []float64{0.25, 0.5, 0.75}
// 	} else if value == 3 {
// 		return []float64{0.5, 0.75, 1}
// 	} else {
// 		return []float64{0.75, 1, 1}
// 	}
// 	return nil
// }
//LinguisticToTFS 评分共0-6七个级别
func LinguisticToTFS(value int) []float64 {

	if value == 0 {
		return []float64{0, 0, 0}
	} else if value == 1 {
		return []float64{0, 0, 0.25}
	} else if value == 2 {
		return []float64{0, 0.25, 0.5}
	} else if value == 3 {
		return []float64{0.25, 0.5, 0.75}
	} else if value == 4 {
		return []float64{0.5, 0.75, 1}
	} else if value == 5 {
		return []float64{0.75, 1, 1}
	}
	return []float64{1, 1, 1}

}
func TFSToFS(tfn []float64) float64 {
	if tfn[2] == 0 {
		return 0
	}
	c_mean := tfn[2]
	return ((tfn[0] + tfn[1] + tfn[2]) / c_mean) / 3
}
