package infrastrucure

import (
	"context"

	"github.com/null-channel/eddington/api/app/models"
	"github.com/uptrace/bun"
)

type BunAppDatastore struct {
	db *bun.DB
}

func NewBunAppDatastore(db *bun.DB) *BunAppDatastore {

	_, err := db.NewCreateTable().Model((*models.NullApplication)(nil)).Exec(context.Background())
	if err != nil {
		panic(err)
	}

	_, err = db.NewCreateTable().Model((*models.NullApplicationService)(nil)).Exec(context.Background())
	if err != nil {
		panic(err)
	}

	return &BunAppDatastore{
		db: db,
	}
}

func (datastore *BunAppDatastore) CreateNullApplication(ctx context.Context, nullApplication *models.NullApplication) error {
	_, err := datastore.db.NewInsert().Model(nullApplication).Exec(context.Background())
	return err
}
func (datastore *BunAppDatastore) CreateNullApplicationService(ctx context.Context, nullApplicationSvc *models.NullApplicationService) error {
	_, err := datastore.db.NewInsert().Model(nullApplicationSvc).Exec(context.Background())
	return err
}
func (datastore *BunAppDatastore) GetApplicationsByOrgID(ctx context.Context, orgId int64) ([]*models.NullApplication, error) {
	nullApplications := []*models.NullApplication{}
	err := datastore.db.NewSelect().Model(&nullApplications).Relation("NullApplicationService").Where("org_id = ?", orgId).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return nullApplications, nil
}

func (datastore *BunAppDatastore) GetApplicationByID(ctx context.Context, id int64) (*models.NullApplication, error) {
	nullApplication := &models.NullApplication{}
	err := datastore.db.NewSelect().Model(nullApplication).Relation("NullApplicationService").Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return nullApplication, nil
}
func (datastore *BunAppDatastore) GetAllAppSvc(ctx context.Context) ([]*models.NullApplicationService, error) {
	nullApplicationService := []*models.NullApplicationService{}
	err := datastore.db.NewSelect().Model(&nullApplicationService).Scan(ctx, &nullApplicationService)
	if err != nil {
		return nil, err
	}
	return nullApplicationService, nil
}

func (datastore *BunAppDatastore) GetApplicationServiceByAppID(ctx context.Context, nullApplicationID int64) ([]*models.NullApplicationService, error) {
	nullApplicationService := []*models.NullApplicationService{}
	err := datastore.db.NewSelect().Model(&nullApplicationService).Where("null_application_id = ?", nullApplicationID).Scan(ctx, &nullApplicationService)

	if err != nil {
		return nil, err
	}
	return nullApplicationService, nil
}
