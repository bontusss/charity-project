package services

import (
	db "charity/db/sqlc"
	"context"
	"database/sql"
)

type BlogService struct {
	queries *db.Queries
}

func NewBlogService(queries *db.Queries) *BlogService {
	return &BlogService{
		queries: queries,
	}
}

func (s *BlogService) CreateBlogPost(ctx context.Context, title, body, imageLink string) (*db.BlogPost, error) {
	params := db.CreateBlogPostParams{
		Title:     title,
		Body:      body,
		ImageLink: sql.NullString{String: imageLink, Valid: imageLink != ""},
	}

	post, err := s.queries.CreateBlogPost(ctx, params)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *BlogService) GetBlogPost(ctx context.Context, id int32) (*db.BlogPost, error) {
	post, err := s.queries.GetBlogPost(ctx, id)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *BlogService) ListBlogPosts(ctx context.Context) ([]db.BlogPost, error) {
	return s.queries.ListBlogPosts(ctx)
}

func (s *BlogService) UpdateBlogPost(ctx context.Context, id int32, title, body, imageLink string) (*db.BlogPost, error) {
	params := db.UpdateBlogPostParams{
		ID:        id,
		Title:     title,
		Body:      body,
		ImageLink: sql.NullString{String: imageLink, Valid: imageLink != ""},
	}

	post, err := s.queries.UpdateBlogPost(ctx, params)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (s *BlogService) DeleteBlogPost(ctx context.Context, id int32) error {
	return s.queries.DeleteBlogPost(ctx, id)
}
