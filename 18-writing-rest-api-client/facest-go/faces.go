package facest

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type FacesList struct {
	Count      int    `json:"count"`
	PagesCount int    `json:"pages_count"`
	Faces      []Face `json:"faces"`
}

type Face struct {
	FaceToken  string      `json:"face_token"`
	FaceID     string      `json:"face_id"`
	FaceImages []FaceImage `json:"face_images"`
	CreatedAt  time.Time   `json:"created_at"`
}

type FaceImage struct {
	ImageToken string    `json:"image_token"`
	ImageURL   string    `json:"image_url"`
	CreatedAt  time.Time `json:"created_at"`
}

type FacesListOptions struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}

func (c *Client) GetFaces(ctx context.Context, options *FacesListOptions) (*FacesList, error) {
	limit := 100
	page := 1
	if options != nil {
		limit = options.Limit
		page = options.Page
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/faces?limit=%d&page=%d", c.BaseURL, limit, page), nil)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)

	res := FacesList{}
	if err := c.sendRequest(req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
