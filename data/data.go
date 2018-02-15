package data

type Data struct {
	ID         string
	ProviderID string
	Metadata   map[string]string
	Contents   []byte
}
