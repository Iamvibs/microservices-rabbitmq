package repositories

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/gorm_postgres"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/logger"
	"github.com/meysamhadeli/shop-golang-microservices/internal/pkg/utils"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/config"
	"github.com/meysamhadeli/shop-golang-microservices/internal/services/product-service/product/models"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type postgresProductRepository struct {
	log  logger.ILogger
	cfg  *config.Config
	db   *pgxpool.Pool
	gorm *gorm.DB
}

func NewPostgresProductRepository(log logger.ILogger, cfg *config.Config, gorm *gorm.DB) *postgresProductRepository {
	return &postgresProductRepository{log: log, cfg: cfg, gorm: gorm}
}

func (p *postgresProductRepository) GetAllProducts(ctx context.Context, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {

	result, err := gorm_postgres.Paginate[*models.Product](ctx, listQuery, p.gorm)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *postgresProductRepository) SearchProducts(ctx context.Context, searchText string, listQuery *utils.ListQuery) (*utils.ListResult[*models.Product], error) {

	whereQuery := fmt.Sprintf("%s IN (?)", "Name")
	query := p.gorm.Where(whereQuery, searchText)

	result, err := gorm_postgres.Paginate[*models.Product](ctx, listQuery, query)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (p *postgresProductRepository) GetProductById(ctx context.Context, uuid uuid.UUID) (*models.Product, error) {

	var product models.Product

	if err := p.gorm.First(&product, uuid).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't find the product with id %s into the database.", uuid))
	}

	return &product, nil
}

func (p *postgresProductRepository) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {

	if err := p.gorm.Create(&product).Error; err != nil {
		return nil, errors.Wrap(err, "error in the inserting product into the database.")
	}

	return product, nil
}

func (p *postgresProductRepository) UpdateProduct(ctx context.Context, updateProduct *models.Product) (*models.Product, error) {

	if err := p.gorm.Save(updateProduct).Error; err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("error in updating product with id %s into the database.", updateProduct.ProductId))
	}

	return updateProduct, nil
}

func (p *postgresProductRepository) DeleteProductByID(ctx context.Context, uuid uuid.UUID) error {

	var product models.Product

	if err := p.gorm.First(&product, uuid).Error; err != nil {
		return errors.Wrap(err, fmt.Sprintf("can't find the product with id %s into the database.", uuid))
	}

	if err := p.gorm.Delete(&product).Error; err != nil {
		return errors.Wrap(err, "error in the deleting product into the database.")
	}

	return nil
}