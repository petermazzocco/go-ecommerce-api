package user

type Customer interface {
	GetUser() (*User, error)
	NewUser(u User) (*User, error)
}

func (u *User) GetUser() (*User, error) {
	return &User{}, nil
}

func (u *User) NewUser(newU User) (*User, error)	{
	return &User{
		ID: newU.ID,
		FirstName: newU.FirstName,
		LastName: newU.LastName,
		Address: newU.Address,
		Email: newU.Email,
	}, nil
}


