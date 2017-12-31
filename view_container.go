package grip

type ViewContainer interface {
	View
	Find(path ...ViewID) (View, error)
}
