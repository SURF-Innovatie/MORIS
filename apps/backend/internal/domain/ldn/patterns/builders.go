// Package patterns provides COAR Notify pattern builders.
// See: https://coar-notify.net/specification/1.0.1/
package patterns

import (
	"os"

	"github.com/SURF-Innovatie/MORIS/internal/domain/ldn"
)

// getOriginURL returns the configured origin URL for this MORIS instance.
func getOriginURL() string {
	originURL := os.Getenv("LDN_ORIGIN_URL")
	if originURL == "" {
		originURL = "http://localhost:8080"
	}
	return originURL
}

// RequestReview builds an "Offer" activity for requesting a review.
// Pattern: https://coar-notify.net/specification/1.0.1/request-review/
func RequestReview(object *ldn.Object, actor *ldn.Actor, target *ldn.Service) *ldn.Activity {
	activity := ldn.NewActivity(ldn.ActivityOffer, getOriginURL())
	activity.Type = append(activity.Type, ldn.COARReviewAction)
	activity.Object = object
	activity.Actor = actor
	activity.Target = target
	return activity
}

// RequestEndorsement builds an "Offer" activity for requesting endorsement.
// Pattern: https://coar-notify.net/specification/1.0.1/request-endorsement/
func RequestEndorsement(object *ldn.Object, actor *ldn.Actor, target *ldn.Service) *ldn.Activity {
	activity := ldn.NewActivity(ldn.ActivityOffer, getOriginURL())
	activity.Type = append(activity.Type, ldn.COAREndorsementAction)
	activity.Object = object
	activity.Actor = actor
	activity.Target = target
	return activity
}

// AcceptRequest builds an "Accept" activity in response to an Offer.
// Pattern: https://coar-notify.net/specification/1.0.1/accept/
func AcceptRequest(inReplyTo string, object *ldn.Object, target *ldn.Service) *ldn.Activity {
	activity := ldn.NewActivity(ldn.ActivityAccept, getOriginURL())
	activity.InReplyTo = inReplyTo
	activity.Object = object
	activity.Target = target
	return activity
}

// RejectRequest builds a "Reject" activity in response to an Offer.
// Pattern: https://coar-notify.net/specification/1.0.1/reject/
func RejectRequest(inReplyTo string, object *ldn.Object, target *ldn.Service) *ldn.Activity {
	activity := ldn.NewActivity(ldn.ActivityReject, getOriginURL())
	activity.InReplyTo = inReplyTo
	activity.Object = object
	activity.Target = target
	return activity
}

// AnnounceEndorsement builds an "Announce" with EndorsementAction.
// Pattern: https://coar-notify.net/specification/1.0.1/announce-endorsement/
func AnnounceEndorsement(object *ldn.Object, actor *ldn.Actor, target *ldn.Service) *ldn.Activity {
	activity := ldn.NewActivity(ldn.ActivityAnnounce, getOriginURL())
	activity.Type = append(activity.Type, ldn.COAREndorsementAction)
	activity.Object = object
	activity.Actor = actor
	activity.Target = target
	return activity
}

// AnnounceReview builds an "Announce" with ReviewAction.
// Pattern: https://coar-notify.net/specification/1.0.1/announce-review/
func AnnounceReview(object *ldn.Object, actor *ldn.Actor, target *ldn.Service) *ldn.Activity {
	activity := ldn.NewActivity(ldn.ActivityAnnounce, getOriginURL())
	activity.Type = append(activity.Type, ldn.COARReviewAction)
	activity.Object = object
	activity.Actor = actor
	activity.Target = target
	return activity
}

// AnnounceRelationship builds an "Announce" with RelationshipAction.
// Pattern: https://coar-notify.net/specification/1.0.1/announce-relationship/
func AnnounceRelationship(subject, object *ldn.Object, relationship string, actor *ldn.Actor, target *ldn.Service) *ldn.Activity {
	activity := ldn.NewActivity(ldn.ActivityAnnounce, getOriginURL())
	activity.Type = append(activity.Type, ldn.COARRelationshipAction)
	activity.Object = object
	activity.Actor = actor
	activity.Target = target
	// Note: relationship context would typically be added to the object
	return activity
}

// AnnounceIngest builds an "Announce" with IngestAction for archival workflows.
// Pattern: https://coar-notify.net/specification/1.0.1/announce-resource/
func AnnounceIngest(object *ldn.Object, actor *ldn.Actor, target *ldn.Service) *ldn.Activity {
	activity := ldn.NewActivity(ldn.ActivityAnnounce, getOriginURL())
	activity.Type = append(activity.Type, ldn.COARIngestAction)
	activity.Object = object
	activity.Actor = actor
	activity.Target = target
	return activity
}

// UnprocessableNotification builds an error response for invalid notifications.
// Pattern: https://coar-notify.net/specification/1.0.1/unprocessable/
func UnprocessableNotification(inReplyTo string, summary string, target *ldn.Service) *ldn.Activity {
	activity := ldn.NewActivity(ldn.ActivityFlag, getOriginURL())
	activity.Type = append(activity.Type, ldn.COARUnprocessableNotification)
	activity.InReplyTo = inReplyTo
	activity.Summary = summary
	activity.Target = target
	return activity
}
