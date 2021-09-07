package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"strconv"

	"github.com/CyberWake/meetmeup/graph/generated"
	"github.com/CyberWake/meetmeup/graph/model"
	"github.com/CyberWake/meetmeup/internal/auth"
	"github.com/CyberWake/meetmeup/internal/links"
	"github.com/CyberWake/meetmeup/internal/pkg/jwt"
	"github.com/CyberWake/meetmeup/internal/users"
)

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	userCreated := user.Create()
	if !userCreated {
		return "", &users.UserAlreadyExistsError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		// 1
		return "", &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Link{}, fmt.Errorf("access denied")
	}
	var link links.Link
	link.Title = input.Title
	link.Address = input.Address
	link.User = user
	linkID := link.Save()
	grahpqlUser := &model.User{
		ID:   user.ID,
		Name: user.Username,
	}
	return &model.Link{ID: strconv.FormatInt(linkID, 10), Title: link.Title, Address: link.Address, User: grahpqlUser}, nil
}

func (r *mutationResolver) UpdateLink(ctx context.Context, input model.UpdateLinkInput) (*model.Link, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return &model.Link{}, fmt.Errorf("access denied")
	}
	var link links.Link
	link.ID = input.ID
	link.Title = *input.Title
	link.Address = *input.Address
	link.User = user
	link, err := links.UpdateLink(link)
	if err != nil {
		return &model.Link{}, err
	}
	return &model.Link{ID: link.ID, Address: link.Address, Title: link.Title, User: &model.User{
		ID:   link.User.ID,
		Name: link.User.Username,
	}}, nil
}

func (r *mutationResolver) DeleteLink(ctx context.Context, input model.LinkID) (string, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return "", fmt.Errorf("access denied")
	}
	result, err := links.DeleteLink(input.ID, user.ID)
	if err != nil {
		return "", err
	}
	return result + input.ID, nil
}

func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	var resultLinks []*model.Link
	var dbLinks = links.GetAll()
	for _, link := range dbLinks {
		grahpqlUser := &model.User{
			ID:   link.User.ID,
			Name: link.User.Username,
		}
		resultLinks = append(resultLinks, &model.Link{ID: link.ID, Title: link.Title, Address: link.Address, User: grahpqlUser})
	}
	return resultLinks, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	var resultUsers []*model.User
	var dbUsers = users.FetchAllUsers()
	for _, user := range dbUsers {
		resultUsers = append(resultUsers, &model.User{ID: user.ID, Name: user.Username})
	}
	return resultUsers, nil
}

func (r *queryResolver) LinkByID(ctx context.Context, input model.LinkID) (*model.Link, error) {
	var id = input.ID
	var resultLink *model.Link
	var link, err = links.GetLink(id)
	if err != nil {
		return nil, &links.LinkNotPresent{}
	}
	resultLink = &model.Link{ID: link.ID, Title: link.Title, Address: link.Address, User: &model.User{
		ID:   link.User.ID,
		Name: link.User.Username,
	}}
	return resultLink, nil
}

func (r *queryResolver) LinksByUserID(ctx context.Context, input model.UserID) ([]*model.Link, error) {
	user := auth.ForContext(ctx)
	var resultLinks []*model.Link
	if user == nil {
		return resultLinks, fmt.Errorf("access denied")
	}
	var dbLinks = links.LinksByUserID(input.ID)
	for _, link := range dbLinks {
		grahpqlUser := &model.User{
			ID:   link.User.ID,
			Name: link.User.Username,
		}
		resultLinks = append(resultLinks, &model.Link{ID: link.ID, Title: link.Title, Address: link.Address, User: grahpqlUser})
	}
	return resultLinks, nil
}

func (r *queryResolver) MyLinks(ctx context.Context) ([]*model.Link, error) {
	user := auth.ForContext(ctx)
	var resultLinks []*model.Link
	if user == nil {
		return resultLinks, fmt.Errorf("access denied")
	}
	var dbLinks = links.LinksByUserID(user.ID)
	for _, link := range dbLinks {
		grahpqlUser := &model.User{
			ID:   link.User.ID,
			Name: link.User.Username,
		}
		resultLinks = append(resultLinks, &model.Link{ID: link.ID, Title: link.Title, Address: link.Address, User: grahpqlUser})
	}
	return resultLinks, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
