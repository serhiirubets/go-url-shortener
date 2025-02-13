package link

type CreateLinkRequest struct {
	Url string `json:"url" validate:"required,url"`
}

type UpdateLinkRequest struct {
	Url  string `json:"url" validate:"required,url"`
	Hash string `json:"hash"`
}

type GetAllLinksResponse struct {
	Links []Link `json:"links"`
	Count int64  `json:"count"`
}
