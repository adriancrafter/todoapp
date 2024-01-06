package auth

func ToUserVM(m User) UserVM {
	return UserVM{
		Slug:              m.Slug,
		Username:          m.Username,
		Email:             m.Email,
		EmailConfirmation: m.Email,
		IsNew:             m.IsNew(),
	}
}

// View model to model

func ToSigninModel(si SigninVM) Signin {
	return Signin{
		Username: si.Username,
		Password: si.Password,
	}
}

// Model to view model

func ToUserListVM(m []User) (vm []UserVM) {
	for _, user := range m {
		vm = append(vm, ToUserVM(user))
	}
	return vm
}
