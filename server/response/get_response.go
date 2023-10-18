package response

type DataResponse struct {
	Data       interface{} `json:"data,omitempty"`
	TotalCount int         `json:"totalCount,omitempty"`
}
