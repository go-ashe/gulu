package mathx

import "math"

type AxisAngle int

const (
	AxisAngle0     AxisAngle = 0
	AxisAngle90    AxisAngle = 90000
	AxisAngle180   AxisAngle = 180000
	AxisAngleNeg90 AxisAngle = -90000
)

type AxisDir int

const (
	AxisDirX AxisDir = iota + 1
	AxisDirNegX
	AxisDirY
	AxisDirNegY
)

// IsAxisAngle 判断角度是否为[0,90,180,-90],误差范围为 diff
func IsAxisAngle(angle, tolerance int) bool {
	angle = normalizeAngleIn180(angle)
	return Abs(angle-90000) <= tolerance ||
		Abs(angle-0) <= tolerance ||
		Abs(Abs(angle)-180000) <= tolerance ||
		Abs(angle+90000) <= tolerance
}

func GetAxisAngle(angle int) AxisAngle {
	switch {
	case angle >= -45000 && angle < 45000:
		return AxisAngle0

	case angle >= 45000 && angle < 135000:
		return AxisAngle90

	case angle >= -135000 && angle < -45000:
		return AxisAngleNeg90

	default:
		// 包含：
		// axis >= 135000
		// axis < -135000
		return AxisAngle180
	}
}

func GetLineAngle(x1, y1, x2, y2 int) int {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	return Rad2Deg(math.Atan2(dy, dx))
}

func Angle2Axis(angle int) AxisDir {
	angle = normalizeAngleIn180(angle)
	if angle >= -45000 && angle < 45000 {
		return AxisDirX
	} else if angle >= 45000 && angle < 135000 {
		return AxisDirY
	} else if angle >= 135000 || angle < -135000 {
		return AxisDirNegX
	} else {
		return AxisDirNegY
	}
}

func DetectDir(dx, dy int) AxisDir {
	if math.Abs(float64(dx)) >= math.Abs(float64(dy)) {
		if dx >= 0 {
			return AxisDirX
		}
		return AxisDirNegX
	}

	if dy >= 0 {
		return AxisDirY
	}
	return AxisDirNegY
}

func DetectAxisDir(x1, y1, x2, y2 int) (ux, uy int, dir AxisDir) {
	dx, dy := x2-x1, y2-y1
	if Abs(dx) >= Abs(dy) {
		if dx >= 0 {
			ux, uy, dir = 1, 0, AxisDirX
		} else {
			dir = AxisDirNegX
			ux, uy, dir = -1, 0, AxisDirNegX
		}
		return
	}

	if dy >= 0 {
		ux, uy, dir = 0, 1, AxisDirY
	} else {
		dir = AxisDirNegX
		ux, uy, dir = 0, -1, AxisDirNegY
	}
	return
}
