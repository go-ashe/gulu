package mathx

import "math"

// PointToLineDistance 点到直线距离（毫米）
func PointToLineDistance(px, py, ax, ay, bx, by int) float64 {
	dx := float64(bx - ax)
	dy := float64(by - ay)

	if dx == 0 && dy == 0 {
		return math.Hypot(float64(px-ax), float64(py-ay))
	}
	return math.Abs(float64(px-ax)*dy-float64(py-ay)*dx) / math.Hypot(dx, dy)
}

// PointInLine 判断点是否在线上
func PointInLine(px, py, x1, y1, x2, y2, dist int) bool {
	A := float64(y2 - y1)
	B := float64(x1 - x2)
	C := float64(x2*y1 - x1*y2)
	return math.Abs(A*float64(px)+B*float64(py)+C)/math.Sqrt(A*A+B*B) <= float64(dist)
}
