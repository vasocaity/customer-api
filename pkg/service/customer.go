package service

import (
	"customer-api/pkg/model"
	"customer-api/pkg/repository"
	"log"

	"github.com/google/uuid"
)

type CustomerService interface {
	Create(req *CreateCustomerRequest) (*model.Customer, error)
	Get(id uuid.UUID) (*model.Customer, error)
	Update(id uuid.UUID, req *UpdateCustomerRequest) (*model.Customer, error)
	Delete(id uuid.UUID) error
	List(query string, limit, offset int) ([]model.Customer, error)
}

type service struct {
	repo repository.CustomerRepository
}

type CreateCustomerRequest struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
	Phone string `json:"phone"`
}

type UpdateCustomerRequest struct {
	Name  *string `json:"name" validate:"required"`
	Email *string `json:"email" validate:"omitempty,email"`
	Phone *string `json:"phone"`
}

// Create implements CustomerService.
func (s *service) Create(req *CreateCustomerRequest) (*model.Customer, error) {
	c := &model.Customer{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

// Delete implements CustomerService.
func (s *service) Delete(id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

// Get implements CustomerService.
func (s *service) Get(id uuid.UUID) (*model.Customer, error) {
	return s.repo.GetByID(id)
}

// List implements CustomerService.
func (s *service) List(query string, limit int, offset int) ([]model.Customer, error) {
	if limit == 0 {
		limit = 10
	}
	log.Default().Printf("limit: %d", limit)
	return s.repo.List(query, limit, offset)
}

// Update implements CustomerService.
func (s *service) Update(id uuid.UUID, req *UpdateCustomerRequest) (*model.Customer, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Email != nil {
		c.Email = *req.Email
	}
	if req.Phone != nil {
		c.Phone = *req.Phone
	}
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil

}

func NewService(r repository.CustomerRepository) CustomerService {
	return &service{repo: r}
}
