package entity

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Gender    int
}

func (u User) GetProperName() string {
	var title string

	if u.Gender == 1 {
		title = "Mr. "
	} else if u.Gender == 2 {
		title = "Mrs. "
	}

	return title + u.GetFullName()
}

func (u User) GetFullName() string {
	if u.LastName == "" {
		return u.FirstName
	}
	return u.FirstName + " " + u.LastName
}
