package mathx

import "math"

// IsAngleWithinTolerance 判断两个角度之间的差值是否在误差范围内
// angle1 和 angle2 是角度（以1/1000度为单位）
// tolerance 是允许的误差范围（以1/1000度为单位）
func IsAngleWithinTolerance(angle1, angle2, tolerance int) bool {
	angle1 = NormalAngle(angle1)
	angle2 = NormalAngle(angle2)

	// 计算两个角度之间的差值
	diff := int(math.Abs(float64(angle1 - angle2)))

	// 考虑角度的循环性质
	if diff > 180000 {
		diff = 360000 - diff
	}

	// 检查差值是否在误差范围内
	return diff <= tolerance
}

// NormalAngle 归一化到0到360000范围内
func NormalAngle(angle int) int {
	// 归一化到0到360000范围内
	normalized := angle % 360000 // 360000对应360度（360 * 1000）

	// 处理特殊情况
	if normalized > 180000 {
		normalized -= 360000
	} else if normalized < -180000 {
		normalized += 360000
	}

	// 特殊处理等于-180000的情况
	if normalized == -180000 {
		normalized = 180000
	}
	return normalized
}

// IsAngleNear 判断两个角度之间的差值是否在误差范围内
func IsAngleNear(angle1, angle2, tolerance int) bool {
	angle1 = NormalAngle(angle1)
	angle2 = NormalAngle(angle2)

	// 计算两个角度之间的差值
	diff := math.Abs(float64(angle1 - angle2))

	// 考虑角度的循环性质
	if diff > 180000 {
		diff = 360000 - diff
	}

	// 检查差值是否在误差范围内
	return int(diff) <= tolerance
}

// CalcRotateAngle 从 angle1 到 angle2 的最小旋转角
func CalcRotateAngle(a1, a2 int) int {
	a1 = normalizeAngleIn180(a1)
	a2 = normalizeAngleIn180(a2)

	diff := a2 - a1

	if diff > 180000 {
		diff -= 360000
	}

	if diff <= -180000 {
		diff += 360000
	}

	return diff
}
