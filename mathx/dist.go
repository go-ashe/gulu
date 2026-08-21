package mathx

import "math"

// distance 欧氏距离
func distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}

// IsPointsNear 判断 2 个点是否相近
func IsPointsNear(x1, x2, y1, y2 int, tolerance float64) bool {
	return DistLessThanThreshold(x1, x2, tolerance) &&
		DistLessThanThreshold(y1, y2, tolerance)
}

// DistLessThanThreshold 判断2个值的距离是否在允许的误差范围内
func DistLessThanThreshold(num1, num2 int, tolerance float64) bool {
	diff := math.Abs(float64(num1 - num2))
	return diff <= tolerance
}

// 点到直线距离（毫米）
func pointToLineDistance(px, py, x1, y1, x2, y2 int) float64 {
	A := float64(y2 - y1)
	B := float64(x1 - x2)
	C := float64(x2*y1 - x1*y2)

	return math.Abs(A*float64(px)+B*float64(py)+C) / math.Sqrt(A*A+B*B)
}

func GetAxisLineAngle(x1, y1, x2, y2 int) (ok bool, angle int) {
	ok = true
	if x1 == x2 {
		if y1 < y2 {
			angle = 90000
		} else {
			angle = -90000
		}
	} else if y1 == y2 {
		if x1 < x2 {
			angle = 0
		} else {
			angle = 180000
		}
	} else {
		ok = false
	}
	return
}
