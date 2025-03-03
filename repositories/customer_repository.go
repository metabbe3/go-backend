package repositories

import (
	"github.com/metabbe3/go-backend/models"
	"gorm.io/gorm"
)

// CustomerRepository is a concrete implementation of the CustomerRepositoryInterface
type CustomerRepository struct {
	DB *gorm.DB
}

// CustomerRepositoryInterface defines the methods to interact with the Customer model
type CustomerRepositoryInterface interface {
	CreateCustomer(customer *models.Customer) error
	FindCustomerByID(id uint) (*models.Customer, error)
	FindCustomerByPhone(phone string) (*models.Customer, error)
	UpdateCustomer(customer *models.Customer) error
	DeleteCustomer(id uint) error
	GetAllCustomers(limit, offset int) ([]models.Customer, int64, error) // Returning totalCount as int64
}

// NewCustomerRepository creates and returns a new instance of CustomerRepository
func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{DB: db}
}

// CreateCustomer saves a new customer in the database
func (r *CustomerRepository) CreateCustomer(customer *models.Customer) error {
	return r.DB.Create(customer).Error
}

// FindCustomerByID retrieves a customer by their ID
func (r *CustomerRepository) FindCustomerByID(id uint) (*models.Customer, error) {
	var customer models.Customer
	if err := r.DB.First(&customer, id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// FindCustomerByPhone retrieves a customer by their phone number
func (r *CustomerRepository) FindCustomerByPhone(phone string) (*models.Customer, error) {
	var customer models.Customer
	if err := r.DB.Where("phone = ?", phone).First(&customer).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

// UpdateCustomer updates the customer's details in the database
func (r *CustomerRepository) UpdateCustomer(customer *models.Customer) error {
	return r.DB.Save(customer).Error
}

// DeleteCustomer deletes a customer by their ID
func (r *CustomerRepository) DeleteCustomer(id uint) error {
	return r.DB.Delete(&models.Customer{}, id).Error
}

// GetAllCustomers retrieves all customers from the database with limit, offset, and total count
func (r *CustomerRepository) GetAllCustomers(limit, offset int) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var totalCount int64

	// Get the total count of customers
	if err := r.DB.Model(&models.Customer{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// Get the customers with limit and offset
	if err := r.DB.Limit(limit).Offset(offset).Find(&customers).Error; err != nil {
		return nil, 0, err
	}

	return customers, totalCount, nil
}
