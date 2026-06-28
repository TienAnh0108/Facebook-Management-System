package facebook

type PageInfo struct {
	ID             string `jason:"id"`
	Name           string `json:"name"`
	Category       string `json:"category"`
	FollowersCount int    `json:"followers_count"`
}

type PostResponse struct {
	Data []Post `json:"data"`
}

type Post struct {
	ID          string       `json:"id"`
	CreatedTime string       `json:"created_time"`
	Message     string       `json:"message"`
	Likes       ReactionData `json:"Like"`
	Comments    ReactionData `json:"comments"`
	Shares      ShareData    `json:"shares"`
}

type ReactionData struct {
	Summary struct {
		TotalCount int `json:"total_count"`
	} `json:"summary"`
}

type ShareData struct {
	Count int `json:"count"`
}
