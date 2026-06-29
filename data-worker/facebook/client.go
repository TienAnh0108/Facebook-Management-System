package facebook

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client quản lý kết nối HTTP
type Client struct {
	HttpClient *http.Client
}

// NewClient khởi tạo một Facebook API Client với thời gian timeout
func NewClient() *Client {
	return &Client{
		HttpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// FetchPageInfo gọi API lấy thông tin cơ bản (Tên, Hạng mục, Followers)
func (c *Client) FetchPageInfo(pageID string, accessToken string) (*PageInfo, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s?fields=name,category,followers_count&access_token=%s", pageID, accessToken)

	resp, err := c.HttpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var info PageInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

// FetchDailyInteractions quét các bài viết trong ngày để tính tổng React, Comment, Share
func (c *Client) FetchDailyInteractions(pageID string, accessToken string) (int, int, int, error) {
	// Lấy mốc thời gian từ 00:00:00 đến 23:59:59 của ngày hôm nay dưới dạng Unix Timestamp
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location()).Unix()

	url := fmt.Sprintf("https://graph.facebook.com/v19.0/%s/posts?fields=likes.summary(true),comments.summary(true),shares&since=%d&until=%d&access_token=%s",
		pageID, startOfDay, endOfDay, accessToken)

	resp, err := c.HttpClient.Get(url)
	if err != nil {
		return 0, 0, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, 0, err
	}

	var postResp PostResponse
	if err := json.Unmarshal(body, &postResp); err != nil {
		return 0, 0, 0, err
	}

	totalReacts := 0
	totalComments := 0
	totalShares := 0

	// Vòng lặp cộng dồn các chỉ số tương tác của từng bài viết trong ngày
	for _, post := range postResp.Data {
		totalReacts += post.Likes.Summary.TotalCount
		totalComments += post.Comments.Summary.TotalCount
		totalShares += post.Shares.Count
	}

	return totalReacts, totalComments, totalShares, nil
}
