package raid

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/external/raid"
	"github.com/SURF-Innovatie/MORIS/internal/domain/entities"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/infra/persistence/eventstore"
	"github.com/google/uuid"
)

type Service interface {
	MintRaid(ctx context.Context, userID uuid.UUID, projectID uuid.UUID, req *raid.RAiDCreateRequest) (*raid.RAiDDto, error)
	UpdateRaid(ctx context.Context, userID uuid.UUID, projectID uuid.UUID, handle string, req *raid.RAiDUpdateRequest) (*raid.RAiDDto, error)
	FindRaid(ctx context.Context, userID uuid.UUID, prefix, suffix string) (*raid.RAiDDto, error)
	FindAllRaids(ctx context.Context, userID uuid.UUID) ([]raid.RAiDDto, error)

	// Local RaidInfo management
	GetRaidInfo(ctx context.Context, raidID string) (*entities.RAiDInfo, error)
	GetRaidInfoByProject(ctx context.Context, projectID uuid.UUID) (*entities.RAiDInfo, error)
	SaveRaidInfo(ctx context.Context, info *entities.RAiDInfo) (*entities.RAiDInfo, error)
}

type Client interface {
	MintRaid(ctx context.Context, req *raid.RAiDCreateRequest) (*raid.RAiDDto, error)
	UpdateRaid(ctx context.Context, prefix, suffix string, req *raid.RAiDUpdateRequest) (*raid.RAiDDto, error)
	FindRaid(ctx context.Context, prefix, suffix string) (*raid.RAiDDto, error)
	FindAllRaids(ctx context.Context) ([]raid.RAiDDto, error)
}

type service struct {
	client Client
	repo   Repository
	es     eventstore.Store
}

func NewService(client Client, repo Repository, es eventstore.Store) Service {
	return &service{
		client: client,
		repo:   repo,
		es:     es,
	}
}

func (s *service) MintRaid(ctx context.Context, userID uuid.UUID, projectID uuid.UUID, req *raid.RAiDCreateRequest) (*raid.RAiDDto, error) {
	// Mint via client
	dto, err := s.client.MintRaid(ctx, req)
	if err != nil {
		return nil, err
	}

	// Persist local info
	if dto != nil {
		info := s.mapDtoToEntity(dto)
		info.ProjectID = projectID
		// We use SaveRaidInfo which handles create/update logic (upsert)
		// But mapDtoToEntity needs to be defined
		if _, err := s.SaveRaidInfo(ctx, info); err != nil {
			// Do we fail the request if local save fails?
			// Ideally yes, or we have inconsistency.
			return nil, err
		}

		// Emit project.raid_linked event
		// We use the pointer to the entity we just created/mapped
		err := s.emitEvent(ctx, projectID, userID, func(pid, actor uuid.UUID, cur *entities.Project, status events.Status) (events.Event, error) {
			input := events.RaidLinkedInput{RAiDInfo: info}
			return events.DecideRaidLinked(pid, actor, cur, input, status)
		})
		if err != nil {
			// If event emit fails, we have inconsistency between local RaidInfoRepo and EventStore.
			// Ideally we run this in a Saga or verify consistency later.
			// For now, return error.
			return nil, err
		}
	}

	return dto, nil
}

func (s *service) UpdateRaid(ctx context.Context, userID uuid.UUID, projectID uuid.UUID, handle string, req *raid.RAiDUpdateRequest) (*raid.RAiDDto, error) {
	// Parse handle to prefix/suffix
	// Assuming simple split by "/"
	// TODO: Make robust handle parsing
	// Parse handle to prefix/suffix
	// Assuming simple split by "/"
	// TODO: Make robust handle parsing
	prefix, suffix := s.parseHandle(handle)
	if prefix == "" || suffix == "" {
		// Should we return error?
		// For now let's assume we need 2 parts.
		// If fails, we can't call client.
		// Use dummy valid for compilation if needed or custom error.
		// Let's assume the handle is "prefix/suffix"
		// Return error
		return nil, errors.New("invalid handle format")
	}

	dto, err := s.client.UpdateRaid(ctx, prefix, suffix, req)
	if err != nil {
		return nil, err
	}

	// Persist local info update
	if dto != nil {
		info := s.mapDtoToEntity(dto)
		info.ProjectID = projectID
		if _, err := s.SaveRaidInfo(ctx, info); err != nil {
			return nil, err
		}

		// Emit project.raid_updated event
		err := s.emitEvent(ctx, projectID, userID, func(pid, actor uuid.UUID, cur *entities.Project, status events.Status) (events.Event, error) {
			input := events.RaidUpdatedInput{RAiDInfo: info}
			return events.DecideRaidUpdated(pid, actor, cur, input, status)
		})
		if err != nil {
			return nil, err
		}
	}

	return dto, nil
}

func (s *service) FindRaid(ctx context.Context, userID uuid.UUID, prefix, suffix string) (*raid.RAiDDto, error) {
	return s.client.FindRaid(ctx, prefix, suffix)
}

func (s *service) FindAllRaids(ctx context.Context, userID uuid.UUID) ([]raid.RAiDDto, error) {
	return s.client.FindAllRaids(ctx)
}

func (s *service) GetRaidInfo(ctx context.Context, raidID string) (*entities.RAiDInfo, error) {
	return s.repo.Get(ctx, raidID)
}

func (s *service) GetRaidInfoByProject(ctx context.Context, projectID uuid.UUID) (*entities.RAiDInfo, error) {
	return s.repo.GetByProjectID(ctx, projectID)
}

func (s *service) SaveRaidInfo(ctx context.Context, info *entities.RAiDInfo) (*entities.RAiDInfo, error) {
	// Check if exists, if not create, else update
	_, err := s.repo.Get(ctx, info.RAiDId)
	if err != nil {
		if ent.IsNotFound(err) {
			return s.repo.Create(ctx, info)
		}
		return nil, err
	}
	return s.repo.Update(ctx, info)
}

func (s *service) mapDtoToEntity(d *raid.RAiDDto) *entities.RAiDInfo {
	info := &entities.RAiDInfo{
		RAiDId:                      d.Identifier.IdValue,
		SchemaUri:                   d.Identifier.SchemaUri,
		RegistrationAgencyId:        d.Identifier.RegistrationAgency.Id,
		RegistrationAgencySchemaUri: d.Identifier.RegistrationAgency.SchemaUri,
		OwnerId:                     d.Identifier.Owner.Id,
		OwnerSchemaUri:              d.Identifier.Owner.SchemaUri,
		License:                     d.Identifier.License,
		Version:                     d.Identifier.Version,
	}
	if d.Identifier.Owner.ServicePoint != nil {
		// Try parse int64
		// But wait, the Entity has *int64 for OwnerServicePoint?
		// And Model has *string context.
		// Let's check entity definition from file view earlier.
		// Entity: OwnerServicePoint *int64
		// Model: ServicePoint *string
		// We need conversion.
		if sp, err := strconv.ParseInt(*d.Identifier.Owner.ServicePoint, 10, 64); err == nil {
			info.OwnerServicePoint = &sp
		}
	}
	return info
}

func (s *service) parseHandle(handle string) (prefix, suffix string) {
	parts := strings.Split(handle, "/")
	if len(parts) >= 2 {
		return parts[0], parts[1]
	}
	// Fallback? or Error?
	// Client might handle empty strings.
	return "", ""
}

func (s *service) emitEvent(ctx context.Context, projectID uuid.UUID, userID uuid.UUID, decide func(uuid.UUID, uuid.UUID, *entities.Project, events.Status) (events.Event, error)) error {
	// 1. Load event stream
	history, version, err := s.es.Load(ctx, projectID)
	if err != nil {
		return err
	}

	// 2. Replay to build state
	proj := &entities.Project{}
	for _, e := range history {
		if applier, ok := e.(events.Applier); ok {
			applier.Apply(proj)
		}
	}

	// 3. Decide new event
	evt, err := decide(projectID, userID, proj, events.StatusApproved)
	if err != nil {
		return err
	}
	if evt == nil {
		return nil // No op
	}

	// 4. Append
	// Retry logic is omitted for brevity, but optimistic lock checks happen here
	return s.es.Append(ctx, projectID, version, evt)
}
