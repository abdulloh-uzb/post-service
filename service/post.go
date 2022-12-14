package service

import (
	"context"
	"database/sql"
	pbc "exam/post-service/genproto/customer"
	pbp "exam/post-service/genproto/post"
	pbr "exam/post-service/genproto/reyting"
	l "exam/post-service/pkg/logger"
	"exam/post-service/service/grpcClient"
	"exam/post-service/storage"
	"fmt"
	"strconv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jmoiron/sqlx"
)

type PostService struct {
	storage storage.IStorage
	logger  l.Logger
	client  grpcClient.GrpcClientI
}

func NewPostService(db *sqlx.DB, log l.Logger, client grpcClient.GrpcClientI) *PostService {
	return &PostService{
		storage: storage.NewStorage(db),
		logger:  log,
		client:  client,
	}
}

func (p *PostService) CreatePost(ctx context.Context, req *pbp.PostReq) (*pbp.Post, error) {

	post, err := p.storage.Post().Create(req)
	if err != nil {
		p.logger.Error("error insert", l.Any("error insert post", err))
		return &pbp.Post{}, status.Error(codes.Internal, "something went wrong, please check create post")
	}
	return post, nil
}

func (p *PostService) DeletePost(ctx context.Context, req *pbp.Id) (*pbp.Empty, error) {
	post, err := p.storage.Post().DeletePost(int(req.Id))
	p.client.Ranking().DeleteRankingByPostId(ctx, &pbr.Id{Id: req.Id})
	if err != nil {
		p.logger.Error("error delete", l.Any("error delete post", err))
		return &pbp.Empty{}, status.Error(codes.Internal, "something went wrong, please check delete post")
	}
	return post, nil
}

func (p *PostService) GetPost(ctx context.Context, req *pbp.Id) (*pbp.GetPostResponse, error) {
	post, err := p.storage.Post().GetPost(int(req.Id))
	if err == sql.ErrNoRows {
		return &pbp.GetPostResponse{Message: "afsuski post topilmadi"}, nil
	}
	if err != nil {
		p.logger.Error("error get", l.Any("error get post", err))
		return &pbp.GetPostResponse{}, status.Error(codes.Internal, "something went wrong, please check get post")
	}
	c_id, err := strconv.Atoi(post.CustomerId)
	if err != nil {
		p.logger.Error("error while strconv", l.Any("error while converting customer_id from string to int", err))
		return &pbp.GetPostResponse{}, status.Error(codes.Internal, "something went wrong, please check get post")
	}
	customerInfo, err := p.client.Customer().GetCustomer(ctx, &pbc.CustomerId{Id: int64(c_id)})
	if err != nil {
		p.logger.Error("error get", l.Any("error get post", err))
		return &pbp.GetPostResponse{}, status.Error(codes.Internal, "something went wrong, please check get post")
	}

	rankings, err := p.client.Ranking().GetRankings(ctx, &pbr.Id{
		Id: req.Id,
	})
	if err != nil {
		p.logger.Error("error get", l.Any("error get post", err))
		return &pbp.GetPostResponse{}, status.Error(codes.Internal, "something went wrong, please check get post")
	}
	post.CustomerInfo = append(post.CustomerInfo, &pbp.Customer{
		FirstName:   customerInfo.FirstName,
		LastName:    customerInfo.LastName,
		Bio:         customerInfo.Bio,
		Email:       customerInfo.Email,
		PhoneNumber: customerInfo.PhoneNumber,
	})
	for _, r := range rankings.Rankings {
		post.Rankings = append(post.Rankings, &pbp.Ranking{
			Name:        r.Name,
			Description: r.Description,
			Ranking:     r.Ranking,
			PostId:      r.PostId,
			CustomerId:  r.CustomerId,
		})
	}
	return post, nil
}

func (p *PostService) ListPost(ctx context.Context, req *pbp.Empty) (*pbp.Posts, error) {
	posts, err := p.storage.Post().ListPost()
	if err != nil {
		p.logger.Error("error list", l.Any("error list post", err))
		return &pbp.Posts{}, status.Error(codes.Internal, "something went wrong, please check list post")

	}

	for _, post := range posts.Posts {
		customerInfo, err := p.client.Customer().GetCustomer(ctx, &pbc.CustomerId{Id: post.CustomerId})
		fmt.Println(err)
		if err != nil {
			p.logger.Error("error get", l.Any("error get post", err))
			return &pbp.Posts{}, status.Error(codes.Internal, "something went wrong, please check get post")
		}

		rankings, err := p.client.Ranking().GetRankings(ctx, &pbr.Id{
			Id: post.Id,
		})
		if err != nil {
			p.logger.Error("error delete", l.Any("error delete post", err))
			return &pbp.Posts{}, status.Error(codes.Internal, "something went wrong, please check get post")
		}
		post.Customer = append(post.Customer, &pbp.Customer{
			FirstName:   customerInfo.FirstName,
			LastName:    customerInfo.LastName,
			Bio:         customerInfo.Bio,
			Email:       customerInfo.Email,
			PhoneNumber: customerInfo.PhoneNumber,
		})
		for _, r := range rankings.Rankings {
			post.Ranking = append(post.Ranking, &pbp.Ranking{
				Name:        r.Name,
				Description: r.Description,
				Ranking:     r.Ranking,
				PostId:      r.PostId,
				CustomerId:  r.CustomerId,
			})
		}
	}

	return posts, nil
}
func (p *PostService) UpdatePost(ctx context.Context, req *pbp.Post) (*pbp.Post, error) {
	post, err := p.storage.Post().UpdatePost(req)
	if err != nil {
		p.logger.Error("error update", l.Any("error update post", err))
		return &pbp.Post{}, status.Error(codes.Internal, "something went wrong, please check update post")
	}
	return post, nil
}

func (p *PostService) GetPostByCustomerId(ctx context.Context, req *pbp.Id) (*pbp.Posts, error) {

	posts, err := p.storage.Post().GetPostByCustomerId(int(req.Id))
	if err != nil {
		p.logger.Error("error update", l.Any("error get post", err))
		return &pbp.Posts{}, status.Error(codes.Internal, "something went wrong, please check get post")
	}
	return posts, err
}

func (p *PostService) DeleteByCustomerId(ctx context.Context, req *pbp.Id) (*pbp.Empty, error) {
	post, err := p.storage.Post().DeletePostByCustomerId(int(req.Id))
	p.client.Ranking().DeleteRankingByPostId(ctx, &pbr.Id{Id: req.Id})
	if err != nil {
		p.logger.Error("error delete", l.Any("error delete post", err))
		return &pbp.Empty{}, status.Error(codes.Internal, "something went wrong, please check delete post")
	}
	return post, nil
}
