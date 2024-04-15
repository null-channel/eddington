package infrastrucure

import (
	"context"

	"github.com/uptrace/bun"

	"github.com/null-channel/eddington/api/users/models"
)

type BunMemberDatastore struct {
	Database *bun.DB
}

func NewBunMemberDatastore(db *bun.DB) *BunMemberDatastore {
	_, err := db.NewCreateTable().Model((*models.User)(nil)).Exec(context.Background())
	if err != nil {
		panic(err)
	}
	_, err = db.NewCreateTable().Model((*models.ResourceGroup)(nil)).Exec(context.Background())
	if err != nil {
		panic(err)
	}
	_, err = db.NewCreateTable().Model((*models.Org)(nil)).Exec(context.Background())
	if err != nil {
		panic(err)
	}

	return &BunMemberDatastore{
		Database: db,
	}
}

func (bmd *BunMemberDatastore) CreateOrUpdateUser(ctx context.Context, user *models.User) error {
	_, err := bmd.Database.NewInsert().
		Model(user).
		On("CONFLICT (id) DO UPDATE").
		Set("name = EXCLUDED.name, email = EXCLUDED.email, newsletterConsent = EXCLUDED.newsletterConsent, dob = EXCLUDED.dob").
		Exec(ctx)

	if err != nil {
		return err
	}
	return nil
}

func (bmd *BunMemberDatastore) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	err := bmd.Database.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (bmd *BunMemberDatastore) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := bmd.Database.NewSelect().
		Model(&user).
		Where("email = ?", email).
		Scan(ctx, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (bmd *BunMemberDatastore) CreateOrg(ctx context.Context, org *models.Org) error {
	_, err := bmd.Database.NewInsert().
		Model(org).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (bmd *BunMemberDatastore) GetOrgByID(ctx context.Context, id int64) (*models.Org, error) {
	var org models.Org
	err := bmd.Database.NewSelect().
		Model(&org).
		Where("id = ?", id).
		Scan(ctx, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (bmd *BunMemberDatastore) UpdateOrg(ctx context.Context, org *models.Org) error {
	_, err := bmd.Database.NewUpdate().
		Model(org).
		Where("id = ?", org.ID).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (bmd *BunMemberDatastore) GetOrgByOwnerId(ctx context.Context, ownerId string) (*models.Org, error) {
	var org models.Org
	err := bmd.Database.NewSelect().
		Model(&org).
		Where("owner_id = ?", ownerId).
		Limit(1).
		Scan(ctx, &org)
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (bmd *BunMemberDatastore) CreateResourceGroup(ctx context.Context, resourceGroup *models.ResourceGroup) error {
	_, err := bmd.Database.NewInsert().
		Model(resourceGroup).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (bmd *BunMemberDatastore) GetResourceGroupByID(ctx context.Context, id int64) (*models.ResourceGroup, error) {
	var resourceGroup models.ResourceGroup
	err := bmd.Database.NewSelect().
		Model(&resourceGroup).
		Where("id = ?", id).
		Scan(ctx, &resourceGroup)
	if err != nil {
		return nil, err
	}
	return &resourceGroup, nil
}

func (bmd *BunMemberDatastore) UpdateResourceGroup(ctx context.Context, resourceGroup *models.ResourceGroup) error {
	_, err := bmd.Database.NewUpdate().
		Model(resourceGroup).
		Where("id = ?", resourceGroup.ID).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
