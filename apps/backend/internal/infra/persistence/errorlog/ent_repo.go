package errorlog

import (
	"context"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/internal/app/errorlog"
	errorlog2 "github.com/SURF-Innovatie/MORIS/internal/domain/errorlog"
	"github.com/google/uuid"
)

type entRepository struct {
	cli *ent.Client
}

func NewRepository(cli *ent.Client) errorlog.Repository {
	return &entRepository{cli: cli}
}

func (r *entRepository) Create(ctx context.Context, in errorlog2.ErrorLogCreateInput) error {
	create := r.cli.ErrorLog.Create().
		SetHTTPMethod(in.HTTPMethod).
		SetRoute(in.Route).
		SetStatusCode(in.StatusCode).
		SetErrorMessage(in.Message)

	if in.StackTrace != nil {
		create.SetStackTrace(*in.StackTrace)
	}

	if in.UserID != nil && *in.UserID != uuid.Nil {
		create.SetUserID(in.UserID.String())
	}

	return create.Exec(ctx)
}
