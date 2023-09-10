package service

import (
	"context"
	"io"

	"github.com/vinicius-gregorio/go_grpc_bff/internal/database"
	"github.com/vinicius-gregorio/go_grpc_bff/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CategoryService struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDB database.Category
}

func NewCategoryService(categoryDB database.Category) *CategoryService {
	return &CategoryService{
		CategoryDB: categoryDB,
	}
}

func (c *CategoryService) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.CategoryResponse, error) {
	category, err := c.CategoryDB.Create(in.Name, in.Description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error creating category: %v", err)
	}

	categoryResp := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return &pb.CategoryResponse{
		Category: categoryResp,
	}, nil
}

func (c *CategoryService) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {
	categories, err := c.CategoryDB.FindAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error listing categories: %v", err)
	}

	var categoriesResp []*pb.Category
	for _, category := range categories {
		categoryResp := &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		}
		categoriesResp = append(categoriesResp, categoryResp)
	}

	return &pb.CategoryList{
		Categories: categoriesResp,
	}, nil
}

func (c *CategoryService) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.CategoryResponse, error) {
	category, err := c.CategoryDB.FindByID(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error getting category: %v", err)
	}

	categoryResp := &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return &pb.CategoryResponse{
		Category: categoryResp,
	}, nil
}

func (c *CategoryService) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := pb.CategoryList{}
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&categories) // aqui
		}
		if err != nil {
			return err
		}
		catResult, err := c.CategoryDB.Create(category.Name, category.Description)
		if err != nil {
			return err
		}
		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          catResult.ID,
			Name:        catResult.Name,
			Description: catResult.Description,
		})
	}
}
