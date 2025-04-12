package dto

type UserDataResponse struct {
	PoorUserDataResponse
	FirstName  string         `json:"first_name"`
	LastName   string         `json:"last_name"`
	SecondName string         `json:"second_name"`
	Birthday   string         `json:"birthday"`
	Region     RegionResponse `json:"region"`
}

type PoorUserDataResponse struct {
	Email      string `json:"email"`
	SyncGoogle bool   `json:"sync_google"`
	SyncApple  bool   `json:"sync_apple"`
}
