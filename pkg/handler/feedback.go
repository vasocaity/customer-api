package handler

import (
	"customer-api/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FeedbackHandler struct {
	db *gorm.DB
}

func NewFeedbackHandler(db *gorm.DB) *FeedbackHandler {
	return &FeedbackHandler{
		db: db,
	}
}

type FeedbackRequest struct {
	CustomerID string `json:"customerId"`
	ProductID  string `json:"productId"`
	Rating     int    `json:"rating"`
	Comment    string `json:"comment"`
}

type FeedbackResponse struct {
	Customer  model.Customer `json:"customer"`
	ProductID model.Product  `json:"product"`
	Rating    int            `json:"rating"`
	Comment   string         `json:"comment"`
}

// สร้าง feedback ใหม่
func (h *FeedbackHandler) CreateFeedback(c *gin.Context) {
	var req FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customerId, _ := uuid.Parse(req.CustomerID)
	productId, _ := uuid.Parse(req.ProductID)

	feedback := model.Feedback{
		ID:         uuid.New(),
		CustomerID: customerId,
		ProductID:  productId,
		Rating:     req.Rating,
		Comment:    req.Comment,
	}

	if err := h.db.Create(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var customer model.Customer
	if err := h.db.First(&customer, customerId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var product model.Product
	if err := h.db.First(&product, productId).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := FeedbackResponse{
		Customer:  customer,
		ProductID: product,
		Rating:    feedback.Rating,
		Comment:   feedback.Comment,
	}

	c.JSON(http.StatusCreated, response)
}

// อ่าน feedback ด้วย id
func (h *FeedbackHandler) GetFeedback(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var feedback model.Feedback
	if err := h.db.First(&feedback, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "feedback not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feedback)
}

// อัพเดต feedback
func (h *FeedbackHandler) UpdateFeedback(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	var req FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var feedback model.Feedback
	if err := h.db.First(&feedback, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "feedback not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	customerId, _ := uuid.Parse(req.CustomerID)
	productId, _ := uuid.Parse(req.ProductID)

	feedback.CustomerID = customerId
	feedback.ProductID = productId
	feedback.Rating = req.Rating
	feedback.Comment = req.Comment

	if err := h.db.Save(&feedback).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, feedback)
}

// ลบ feedback
func (h *FeedbackHandler) DeleteFeedback(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID"})
		return
	}

	if err := h.db.Delete(&model.Feedback{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// list feedback ทั้งหมด (optionally filter by customer or product)
func (h *FeedbackHandler) ListFeedbacks(c *gin.Context) {
	var feedbacks []model.Feedback

	customerID := c.Query("customer_id")
	productID := c.Query("product_id")

	query := h.db.Model(&model.Feedback{})

	if customerID != "" {
		cid, err := uuid.Parse(customerID)
		if err == nil {
			query = query.Where("customer_id = ?", cid)
		}
	}

	if productID != "" {
		pid, err := uuid.Parse(productID)
		if err == nil {
			query = query.Where("product_id = ?", pid)
		}
	}

	if err := query.Find(&feedbacks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// var res []FeedbackResponse
	// for _, f := range feedbacks {
	// 	res = append(res, f)
	// }

	c.JSON(http.StatusOK, feedbacks)
}
