package repositories

import (
	"fmt"
	"sistem-pengelolaan-bank-sampah/dto"
	"sistem-pengelolaan-bank-sampah/models"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TrashTransactionRepository interface {
	FindAllTrashTransaction(limit, offset int, filter dto.TrashTransactionFilter, searchQuery string) (*[]models.TrxTrashCustomer, int64, error)
	FindTrashTransactionByID(trashTransactionID uuid.UUID) (*models.TrxTrashCustomer, error)
	CreateTrashTransaction(trashTransaction *models.TrxTrashCustomer) (*models.TrxTrashCustomer, error)
	UpdateTrashTransaction(trashTransaction *models.TrxTrashCustomer) (*models.TrxTrashCustomer, error)
	DeleteTrashTransaction(trashTransaction *models.TrxTrashCustomer) (*models.TrxTrashCustomer, error)
	// check trash category
	FindTrashCategoryByID(trashID uint) (*models.MstTrashCategory, error)
	// check user
	FindUserByID(id uuid.UUID) (*models.MstUser, error)
}

func (r *repository) FindAllTrashTransaction(limit, offset int, filter dto.TrashTransactionFilter, searchQuery string) (*[]models.TrxTrashCustomer, int64, error) {
	var (
		transactionTrash      []models.TrxTrashCustomer
		totalTransactionTrash int64
	)

	trx := r.db.Session(&gorm.Session{})

	if filter.UserID != "" {
		userId, _ := uuid.Parse(filter.UserID)
		trx = trx.Where("user_id = ?", userId)
	}

	if filter.TrashID != 0 {
		trx = trx.Where("trash_id = ?", filter.TrashID)
	}

	// join tables, used for complex searching on relation table
	trx = trx.Joins("JOIN mst_trash_categories ON mst_trash_categories.id = trx_trash_customers.trash_id")
	trx = trx.Joins("JOIN mst_users ON mst_users.id = trx_trash_customers.user_id")

	if searchQuery != "" {
		intQuery, _ := strconv.Atoi(searchQuery)

		trx = trx.Where(`mst_trash_categories.category LIKE ? 
			OR mst_trash_categories.price = ? 
			OR mst_users.full_name LIKE ? 
			OR mst_users.email LIKE ? 
			OR mst_users.phone LIKE ? 
			OR mst_users.address LIKE ? 
			OR qty = ?
			OR total_price = ?`,
			fmt.Sprintf("%%%s%%", searchQuery), // category
			fmt.Sprintf("%d", intQuery),        // price
			fmt.Sprintf("%%%s%%", searchQuery), // fullname
			fmt.Sprintf("%%%s%%", searchQuery), // email
			fmt.Sprintf("%%%s%%", searchQuery), // phone
			fmt.Sprintf("%%%s%%", searchQuery), // address
			fmt.Sprintf("%d", intQuery),        // qty
			fmt.Sprintf("%d", intQuery))        // total
	}

	// preloading, used for get relation data for results
	trx = trx.Preload("TrashCategory").
		Preload("User").Preload("User.Role").
		Preload("Creator").Preload("Creator.Role")

	trx.Model(&models.TrxTrashCustomer{}).
		Count(&totalTransactionTrash)

	err := trx.Limit(limit).
		Offset(offset).
		Find(&transactionTrash).Error

	return &transactionTrash, totalTransactionTrash, err
}

func (r *repository) FindTrashTransactionByID(trashTransactionID uuid.UUID) (*models.TrxTrashCustomer, error) {
	var trash models.TrxTrashCustomer

	err := r.db.Preload("TrashCategory").
		Preload("User").Preload("User.Role").
		Preload("Creator").Preload("Creator.Role").
		Where("id = ?", trashTransactionID).
		First(&trash).Error

	return &trash, err
}

func (r *repository) CreateTrashTransaction(trashTransaction *models.TrxTrashCustomer) (*models.TrxTrashCustomer, error) {
	err := r.db.Create(trashTransaction).Error

	return trashTransaction, err
}

func (r *repository) UpdateTrashTransaction(trashTransaction *models.TrxTrashCustomer) (*models.TrxTrashCustomer, error) {
	err := r.db.Model(&trashTransaction).
		Updates(*trashTransaction).Error

	return trashTransaction, err
}

func (r *repository) DeleteTrashTransaction(trashTransaction *models.TrxTrashCustomer) (*models.TrxTrashCustomer, error) {
	err := r.db.Delete(trashTransaction).Error

	return trashTransaction, err
}
