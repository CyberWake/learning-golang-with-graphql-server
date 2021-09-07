package users

type WrongUsernameOrPasswordError struct{}

func (m *WrongUsernameOrPasswordError) Error() string {
	return "wrong username or password"
}

type UserAlreadyExistsError struct{}

func (m *UserAlreadyExistsError) Error() string {
	return "username already registered"
}
