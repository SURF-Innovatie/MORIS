package raid

import (
	"context"
	"errors"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/raidinfo"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/google/uuid"
)

type EntRepo struct {
	cli *ent.Client
}

func NewEntRepo(cli *ent.Client) *EntRepo {
	return &EntRepo{cli: cli}
}

func (r *EntRepo) Get(ctx context.Context, raidID string) (*entities.RAiDInfo, error) {
	e, err := r.cli.RaidInfo.Query().
		Where(raidinfo.RaidIDEQ(raidID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapToEntity(e), nil
}

func (r *EntRepo) GetByProjectID(ctx context.Context, projectID uuid.UUID) (*entities.RAiDInfo, error) {
	e, err := r.cli.RaidInfo.Query().
		Where(raidinfo.ProjectIDEQ(projectID)).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return r.mapToEntity(e), nil
}

func (r *EntRepo) Create(ctx context.Context, info *entities.RAiDInfo) (*entities.RAiDInfo, error) {
	_, err := r.cli.RaidInfo.Create().
		SetRaidID(info.RAiDId).
		SetSchemaURI(info.SchemaUri).
		SetRegistrationAgencyID(info.RegistrationAgencyId).
		SetRegistrationAgencySchemaURI(info.RegistrationAgencySchemaUri).
		SetOwnerID(info.OwnerId).
		SetOwnerSchemaURI(info.OwnerSchemaUri).
		SetNillableOwnerServicePoint(info.OwnerServicePoint).
		SetProjectID(info.ProjectID).
		SetLicense(info.License).
		SetVersion(info.Version).
		SetNillableLatestSync(info.LatestSync).
		SetDirty(info.Dirty).
		SetNillableChecksum(info.Checksum).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	// Reload to get edges
	return r.Get(ctx, info.RAiDId)
}

func (r *EntRepo) Update(ctx context.Context, info *entities.RAiDInfo) (*entities.RAiDInfo, error) {
	n, err := r.cli.RaidInfo.Update().
		Where(raidinfo.RaidIDEQ(info.RAiDId)).
		SetSchemaURI(info.SchemaUri).
		SetRegistrationAgencyID(info.RegistrationAgencyId).
		SetRegistrationAgencySchemaURI(info.RegistrationAgencySchemaUri).
		SetOwnerID(info.OwnerId).
		SetOwnerSchemaURI(info.OwnerSchemaUri).
		SetNillableOwnerServicePoint(info.OwnerServicePoint).
		SetProjectID(info.ProjectID).
		SetLicense(info.License).
		SetVersion(info.Version).
		SetNillableLatestSync(info.LatestSync).
		SetDirty(info.Dirty).
		SetNillableChecksum(info.Checksum).
		Save(ctx)

	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, errors.New("raid info not found")
	}

	// Fetch updated to return full entity
	return r.Get(ctx, info.RAiDId)
}

func (r *EntRepo) Delete(ctx context.Context, raidID string) error {
	_, err := r.cli.RaidInfo.Delete().Where(raidinfo.RaidIDEQ(raidID)).Exec(ctx)
	return err
}

func (r *EntRepo) mapToEntity(e *ent.RaidInfo) *entities.RAiDInfo {
	if e == nil {
		return nil
	}

	return &entities.RAiDInfo{
		RAiDId:                      e.RaidID,
		SchemaUri:                   e.SchemaURI,
		RegistrationAgencyId:        e.RegistrationAgencyID,
		RegistrationAgencySchemaUri: e.RegistrationAgencySchemaURI,
		OwnerId:                     e.OwnerID,
		OwnerSchemaUri:              e.OwnerSchemaURI,
		OwnerServicePoint:           e.OwnerServicePoint,
		ProjectID:                   e.ProjectID,
		License:                     e.License,
		Version:                     e.Version,
		LatestSync:                  e.LatestSync,
		Dirty:                       e.Dirty,
		Checksum:                    e.Checksum,
	}
}
