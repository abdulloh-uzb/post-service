package postgres

import (
	"fmt"
	pbp "post-service/genproto/post"
	"time"

	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}

func (p *postRepo) Create(req *pbp.PostReq) (*pbp.Post, error) {
	post := pbp.Post{}
	err := p.db.QueryRow(`insert into posts(name, description, customer_id) values($1,$2,$3) 
		returning id, name, description, customer_id`,
		req.Name, req.Description, req.CustomerId).
		Scan(&post.Id, &post.Name, &post.Description, &post.CustomerId)

	fmt.Println(err)
	if err != nil {
		return &pbp.Post{}, err
	}
	fmt.Println(err)

	for _, media := range req.Medias {
		mediaResp := &pbp.Media{}
		err := p.db.QueryRow(`insert into medias(post_id, name, link, type) values($1,$2,$3, $4) returning name, link, type`,
			post.Id, media.Name, media.Link, media.Type).Scan(&mediaResp.Name, &mediaResp.Link, &mediaResp.Type)
		if err != nil {
			return &pbp.Post{}, err
		}
		post.Medias = append(post.Medias, mediaResp)
	}

	return &post, nil

}

func (p *postRepo) DeletePost(id int) (*pbp.Empty, error) {

	_, err := p.db.Exec(`update posts set deleted_at = $1 where id=$2 and deleted_at is null`, time.Now(), id)
	if err != nil {
		return &pbp.Empty{}, err
	}
	return &pbp.Empty{}, nil

}

func (p *postRepo) GetPost(id int) (*pbp.GetPostResponse, error) {
	post := &pbp.GetPostResponse{}
	err := p.db.QueryRow(`select id, name, description from posts where id=$1 and deleted_at is null`, id).
		Scan(&post.Id, &post.Name, &post.Description)
	if err != nil {
		return &pbp.GetPostResponse{}, err
	}

	rows, err := p.db.Query(`select name, link, type from medias where post_id=$1`, post.Id)
	for rows.Next() {
		media := &pbp.Media{}
		err := rows.Scan(&media.Name, &media.Link, &media.Type)
		if err != nil {
			return &pbp.GetPostResponse{}, err
		}
		post.Medias = append(post.Medias, media)
	}

	return post, nil
}
func (p *postRepo) ListPost() (*pbp.Posts, error) {
	posts := pbp.Posts{}
	rows, err := p.db.Query(`select id, name, description, customer_id from posts where deleted_at is null`)
	if err != nil {
		return &pbp.Posts{}, err
	}

	for rows.Next() {
		post := &pbp.Post{}
		err := rows.Scan(&post.Id, &post.Name, &post.Description, &post.CustomerId)

		rows, err := p.db.Query(`select name, link, type from medias where post_id=$1`, post.Id)
		for rows.Next() {
			media := &pbp.Media{}
			err := rows.Scan(&media.Name, &media.Link, &media.Type)
			if err != nil {
				return &pbp.Posts{}, err
			}
			post.Medias = append(post.Medias, media)
		}
		if err != nil {
			return &pbp.Posts{}, err
		}
		posts.Posts = append(posts.Posts, post)
	}

	return &posts, nil
}
func (p *postRepo) UpdatePost(req *pbp.Post) (*pbp.Post, error) {
	post := &pbp.Post{}
	fmt.Println(req, post)
	return post, nil
}

func (p *postRepo) GetPostByCustomerId(id int) (*pbp.Posts, error) {
	posts := &pbp.Posts{}
	rows, err := p.db.Query(`select id, name, description, customer_id from posts where customer_id = $1 and deleted_at is null`, id)

	if err != nil {
		return &pbp.Posts{}, err
	}

	for rows.Next() {
		post := &pbp.Post{}
		err := rows.Scan(&post.Id, &post.Name, &post.Description, &post.CustomerId)
		if err != nil {
			return &pbp.Posts{}, err
		}
		posts.Posts = append(posts.Posts, post)
	}

	return posts, nil
}
