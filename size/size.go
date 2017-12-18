package size

type unit int

const (
	Point unit = iota
	Percent
	Fraction
)

var Auto = Size{Unit: Fraction, Value: 1}

type Size struct {
	Value uint32
	Unit  unit
}

func WithFraction(value uint32) Size {
	return Size{Value: value, Unit: Fraction}
}

func WithPercent(value uint32) Size {
	return Size{Value: value, Unit: Percent}
}

func WithPoints(value uint32) Size {
	return Size{Value: value, Unit: Point}
}
