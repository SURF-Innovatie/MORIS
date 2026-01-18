package event

import (
	"context"
	"fmt"

	"reflect"
	"strings"

	"github.com/SURF-Innovatie/MORIS/ent"
	"github.com/SURF-Innovatie/MORIS/ent/notification"
	notifservice "github.com/SURF-Innovatie/MORIS/internal/app/notification"
	"github.com/SURF-Innovatie/MORIS/internal/app/project/queries"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events"
	"github.com/SURF-Innovatie/MORIS/internal/domain/events/hydrator"
	"github.com/google/uuid"
)

type StatusUpdateNotificationHandler struct {
	Notifier       notifservice.Service
	Cli            *ent.Client
	Hydrator       *hydrator.Hydrator
	ProjectService queries.Service
}

func (h *StatusUpdateNotificationHandler) Handle(ctx context.Context, e events.Event) error {
	status := e.GetStatus()
	if status != "approved" && status != "rejected" {
		return nil
	}

	u, err := ResolveUser(ctx, h.Cli, e.CreatedByID())
	if err != nil || u == nil {
		return err
	}

	// status is already retrieved above
	eventType := e.Type()

	meta := events.GetMeta(eventType)
	friendlyName := meta.FriendlyName
	if friendlyName == "" {
		friendlyName = eventType
	}

	msg := ""
	template := ""
	vars := make(map[string]string)

	if rn, ok := e.(events.RichNotifier); ok {
		if status == "approved" {
			template = rn.ApprovedTemplate()
		} else if status == "rejected" {
			template = rn.RejectedTemplate()
		}
		// Copy event vars
		for k, v := range rn.NotificationVariables() {
			vars[k] = v
		}
	}

	if template != "" {
		// Hydrate
		detailed := h.Hydrator.HydrateOne(ctx, e)

		// Fetch project for vars
		if e.AggregateID() != (uuid.UUID{}) {
			if details, err := h.ProjectService.GetProject(ctx, e.AggregateID()); err == nil && details != nil {
				h.flattenToMap("project", details.Project, vars)
			}
		}

		h.flattenToMap("creator", detailed.Creator, vars)
		h.flattenToMap("product", detailed.Product, vars)
		h.flattenToMap("person", detailed.Person, vars)
		h.flattenToMap("org_node", detailed.OrgNode, vars)
		h.flattenToMap("role", detailed.ProjectRole, vars)

		msg = h.renderTemplate(template, vars)
	} else {
		msg = fmt.Sprintf("Your request '%s' has been %s.", friendlyName, status)
	}

	_, err = h.Cli.Notification.
		Create().
		SetMessage(msg).
		SetUser(u).
		SetEventID(e.GetID()).
		SetType(notification.TypeStatusUpdate).
		Save(ctx)

	return err
}

func (h *StatusUpdateNotificationHandler) renderTemplate(template string, vars map[string]string) string {
	result := template
	for k, v := range vars {
		result = strings.ReplaceAll(result, fmt.Sprintf("{{%s}}", k), v)
	}
	return result
}

func (h *StatusUpdateNotificationHandler) flattenToMap(prefix string, obj any, vars map[string]string) {
	if obj == nil {
		return
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		key := fmt.Sprintf("%s.%s", prefix, f.Name)
		val := v.Field(i).Interface()

		// Dereference pointers for cleaner output
		valV := reflect.ValueOf(val)
		if valV.Kind() == reflect.Ptr {
			if !valV.IsNil() {
				val = valV.Elem().Interface()
			} else {
				val = ""
			}
		}

		vars[key] = fmt.Sprint(val)
	}
}
