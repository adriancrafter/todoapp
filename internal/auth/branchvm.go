package auth

type (
	Branch struct {
		TenantID    string `json:"tenantID" schema:"tenant-id"`
		Slug        string `json:"slug" schema:"slug"`
		Name        string `json:"name" schema:"name"`
		Description string `json:"description" schema:"description"`
		Status      string `json:"status" schema:"status"`
	}

	Branches struct {
		List []Branch `json:"branches"`
	}
)

func (branches *Branches) Add(branch Branch) {
	branches.List = append(branches.List, branch)
}
