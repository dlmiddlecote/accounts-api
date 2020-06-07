package account

// Account represents an Ometria client account.
type Account struct {
	ID    int    `json:"id"`
	Hash  string `json:"hash"`
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Service represents a service for managing accounts.
type Service interface {
	Account(id int) (*Account, error)
}
