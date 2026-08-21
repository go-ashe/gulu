package mathx

import (
	"math"
)

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Int2Float 把整数类型转换为浮点数类型
func Int2Float[T Integer](v T) float64 {
	return float64(v) / 100
}

// Float2Int 把浮点数类型转换为整数类型
func Float2Int[T Integer, F ~float32 | ~float64](v F) T {
	r := int64(v*100 + 0.5)
	return T(r)
}

// Mm2M 毫米 -> 米，保留 2 位小数
func Mm2M(mm int) float64 {
	r := float64(mm) / 1000.0 // 先转 float64
	return math.Round(r*100) / 100
}

// M2Mm  米 -> 毫米
func M2Mm(m float64) int {
	return int(math.Round(m * 1000))
}

// Deg2Rad 毫度 -> 弧度，保留 2 位小数，支持负数
// mdeg范围：(-180000, 180000]
func Deg2Rad(mdeg int) float64 {
	mdeg = normalizeAngleIn180(mdeg)
	if mdeg < -175000 || mdeg > 175000 {
		mdeg = 180000
	} else if mdeg > -5000 && mdeg < 5000 {
		mdeg = 0
	}

	rad := float64(mdeg) * math.Pi / 180000.0
	if rad >= 0 {
		return math.Round(rad*100) / 100
	}
	// 对负数统一四舍五入
	return -math.Round(-rad*100) / 100
}

// Rad2Deg 弧度 -> 毫度，保留 2 位小数，支持负数
func Rad2Deg(rad float64) int {
	mdeg := int(math.Round(rad * 180000.0 / math.Pi))
	mdeg = normalizeAngleIn180(mdeg)
	if mdeg < -175000 || mdeg > 175000 {
		mdeg = 180000
	} else if mdeg > -5000 && mdeg < 5000 {
		mdeg = 0
	}
	return mdeg
}

func normalizeAngleIn180(v int) int {
	if v%360000 == 0 {
		v = 0
	} else if v%180000 == 0 {
		return 180000
	}
	const base = 180000
	return ((v+base)%(2*base)+(2*base))%(2*base) - base
}

// ReverseDirection 计算反方向角度
func ReverseDirection[T int | int8 | int16 | int32 | int64](angle T) T {
	reverse := int64(angle) - 180000

	// 规范化到 [-180000, 180000) 范围内
	reverse = ((reverse + 180000) % 360000) - 180000

	// 处理负数取模的情况
	if reverse < -180000 {
		reverse += 360000
	} else if reverse >= 180000 {
		reverse -= 360000
	}

	return T(reverse)
}
