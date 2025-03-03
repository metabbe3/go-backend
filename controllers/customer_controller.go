package controllers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/metabbe3/go-backend/models"
	"github.com/metabbe3/go-backend/repositories"
	"github.com/metabbe3/go-backend/utils"
)

type CustomerController struct {
	CustomerRepo repositories.CustomerRepositoryInterface
}

// NewCustomerController returns a new instance of CustomerController
func NewCustomerController(customerRepo repositories.CustomerRepositoryInterface) *CustomerController {
	return &CustomerController{CustomerRepo: customerRepo}
}

// CreateCustomer handles customer creation
func (ctrl *CustomerController) CreateCustomer(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
		Phone string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, "Invalid request data", err.Error())
		return
	}

	customer := models.Customer{
		Name:  req.Name,
		Email: req.Email,
		Phone: req.Phone,
	}

	if err := ctrl.CustomerRepo.CreateCustomer(&customer); err != nil {
		utils.SendInternalServerError(c, "Failed to create customer")
		return
	}

	utils.SendCreated(c, "Customer created successfully", gin.H{"customer": customer})
}

// GetCustomer handles getting customer details
func (ctrl *CustomerController) GetCustomer(c *gin.Context) {
	customerID := c.Param("id")
	id, err := strconv.ParseUint(customerID, 10, 32)
	if err != nil {
		utils.SendValidationError(c, "Invalid customer ID", err.Error())
		return
	}

	customer, err := ctrl.CustomerRepo.FindCustomerByID(uint(id))
	if err != nil {
		utils.SendNotFound(c, "Customer not found")
		return
	}

	utils.SendSuccess(c, "Customer details fetched successfully", gin.H{"customer": customer})
}

// UpdateCustomer handles updating customer details
func (ctrl *CustomerController) UpdateCustomer(c *gin.Context) {
	customerID := c.Param("id")
	id, err := strconv.ParseUint(customerID, 10, 32)
	if err != nil {
		utils.SendValidationError(c, "Invalid customer ID", err.Error())
		return
	}

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email" binding:"email"`
		Phone string `json:"phone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendValidationError(c, "Invalid request data", err.Error())
		return
	}

	customer, err := ctrl.CustomerRepo.FindCustomerByID(uint(id))
	if err != nil {
		utils.SendNotFound(c, "Customer not found")
		return
	}

	// Update customer fields
	if req.Name != "" {
		customer.Name = req.Name
	}
	if req.Email != "" {
		customer.Email = req.Email
	}
	if req.Phone != "" {
		customer.Phone = req.Phone
	}

	if err := ctrl.CustomerRepo.UpdateCustomer(customer); err != nil {
		utils.SendInternalServerError(c, "Failed to update customer")
		return
	}

	utils.SendSuccess(c, "Customer updated successfully", gin.H{"customer": customer})
}

// DeleteCustomer handles deleting a customer
func (ctrl *CustomerController) DeleteCustomer(c *gin.Context) {
	customerID := c.Param("id")
	id, err := strconv.ParseUint(customerID, 10, 32)
	if err != nil {
		utils.SendValidationError(c, "Invalid customer ID", err.Error())
		return
	}

	err = ctrl.CustomerRepo.DeleteCustomer(uint(id))
	if err != nil {
		utils.SendNotFound(c, "Customer not found")
		return
	}

	utils.SendSuccess(c, "Customer deleted successfully", nil)
}

// GetAllCustomers handles fetching all customers with pagination
func (ctrl *CustomerController) GetAllCustomers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	customers, totalCount, err := ctrl.CustomerRepo.GetAllCustomers(limit, offset)
	if err != nil {
		utils.SendInternalServerError(c, "Failed to fetch customers")
		return
	}

	utils.SendSuccess(c, "Customers fetched successfully", gin.H{
		"data":        customers,
		"total_count": totalCount,
	})
}
