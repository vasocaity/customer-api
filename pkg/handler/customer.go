package handler

import (
	"customer-api/pkg/repository"
	"customer-api/pkg/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	svc         service.CustomerService
	validate    *validator.Validate
	productRepo repository.ProductRepository
}

func NewCustomerHandler(svc service.CustomerService, productRepo repository.ProductRepository) *CustomerHandler {
	return &CustomerHandler{
		svc:         svc,
		validate:    validator.New(),
		productRepo: productRepo,
	}
}

func (h *CustomerHandler) CreateCustomer(c *gin.Context) {
	var req service.CreateCustomerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.svc.Create(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *CustomerHandler) UpdateByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req service.UpdateCustomerRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cust, err := h.svc.Update(id, &req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cust)
}

func (h *CustomerHandler) Get(c *gin.Context) {
	keyword := c.Query("keyword")
	limit, _ := strconv.Atoi(c.Query("limit"))
	offset, _ := strconv.Atoi(c.Query("offset"))

	cuts, err := h.svc.List(keyword, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var responses []service.CustomerResponse

	for _, v := range cuts {
		var feedbackResponse []service.FeedbackResponse
		uniqueProduct := make(map[uuid.UUID][]service.CommentResponse)
		productName := make(map[uuid.UUID]string)
		for _, f := range v.Feedbacks {
			p, productErr := h.productRepo.GetByID(f.ProductID)
			if productErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": productErr.Error()})
				return
			}
			productName[p.ID] = p.Name
			uniqueProduct[p.ID] = append(uniqueProduct[p.ID], service.CommentResponse{
				Comment: f.Comment,
				Rating:  f.Rating,
			})
		}

		for i, v := range uniqueProduct {
			feedbackResponse = append(feedbackResponse, service.FeedbackResponse{
				ProductName: productName[i],
				Comments:    v,
			})
		}

		customerResponse := service.CustomerResponse{
			Name:      v.Name,
			Email:     v.Email,
			Phone:     v.Phone,
			Feedbacks: feedbackResponse,
		}

		responses = append(responses, customerResponse)
	}

	c.JSON(http.StatusOK, responses)
}

func (h *CustomerHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	cust, err := h.svc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cust)
}

func (h *CustomerHandler) DeleteByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.svc.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete successfully"})

}
