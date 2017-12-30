package size

type unit int

const (
	Point unit = iota
	Percent
	Fraction
)

var Auto = Size{Unit: Fraction, Value: 1}

type Size struct {
	Value int
	Unit  unit
}

func WithFraction(value int) Size {
	return Size{Value: value, Unit: Fraction}
}

func WithPercent(value int) Size {
	return Size{Value: value, Unit: Percent}
}

func WithPoints(value int) Size {
	return Size{Value: value, Unit: Point}
}
