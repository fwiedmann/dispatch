package code

type Info struct {
	OwnerId      string `json:"owner_id"`
	LocationName string `json:"location_name"`
}

type Owner struct {
	Id      string   `json:"owner_id"`
	Name    string   `json:"name"`
	Members []string `json:"members"`
}
