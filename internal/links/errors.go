package links

type LinkNotPresent struct{}

func (m *LinkNotPresent) Error() string {
	return "Link doesn't exist"
}

type LinkUpdationRightMissing struct{}

func (m *LinkUpdationRightMissing) Error() string {
	return "Link updation right missing"
}
