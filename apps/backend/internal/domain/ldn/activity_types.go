package ldn

// Activity Streams 2.0 Activity Types
// See: https://www.w3.org/TR/activitystreams-vocabulary/#activity-types
const (
	// Core AS2 activity types used by COAR Notify
	ActivityAccept          = "Accept"
	ActivityAnnounce        = "Announce"
	ActivityCreate          = "Create"
	ActivityDelete          = "Delete"
	ActivityFlag            = "Flag"
	ActivityOffer           = "Offer"
	ActivityReject          = "Reject"
	ActivityTentativeAccept = "TentativeAccept"
	ActivityTentativeReject = "TentativeReject"
	ActivityUndo            = "Undo"
	ActivityUpdate          = "Update"
)

// COAR Notify Activity Types (extended vocabulary)
// See: https://coar-notify.net/specification/vocabulary/
const (
	// EndorsementAction - for endorsement/review publishing workflows
	COAREndorsementAction = "coar-notify:EndorsementAction"

	// IngestAction - for archival/ingest workflows
	COARIngestAction = "coar-notify:IngestAction"

	// RelationshipAction - for linking resources
	COARRelationshipAction = "coar-notify:RelationshipAction"

	// ReviewAction - for peer review workflows
	COARReviewAction = "coar-notify:ReviewAction"

	// UnprocessableNotification - error response
	COARUnprocessableNotification = "coar-notify:UnprocessableNotification"
)

// NotificationTypeMapping maps legacy MORIS notification types to AS2 activity types.
var NotificationTypeMapping = map[string]string{
	"info":             ActivityAnnounce,
	"approval_request": ActivityOffer,
	"status_update":    ActivityUpdate,
}

// AS2TypeFromLegacy converts a legacy notification type to AS2 activity type.
func AS2TypeFromLegacy(legacyType string) string {
	if as2Type, ok := NotificationTypeMapping[legacyType]; ok {
		return as2Type
	}
	return ActivityAnnounce
}

// LegacyTypeFromAS2 converts an AS2 activity type back to legacy type.
func LegacyTypeFromAS2(as2Type string) string {
	for legacy, as2 := range NotificationTypeMapping {
		if as2 == as2Type {
			return legacy
		}
	}
	return "info"
}

// Relationship types from scholarly ontologies
// See: https://purl.org/vocab/relationship/
const (
	RelationshipCites     = "https://purl.org/vocab/relationship/cites"
	RelationshipIsPartOf  = "https://purl.org/vocab/relationship/isPartOf"
	RelationshipHasPart   = "https://purl.org/vocab/relationship/hasPart"
	RelationshipIsBasedOn = "https://purl.org/vocab/relationship/isBasedOn"
)

// Object types commonly used in scholarly communication
const (
	ObjectTypeDocument     = "Document"
	ObjectTypeDataset      = "Dataset"
	ObjectTypeSoftware     = "Software"
	ObjectTypeArticle      = "Article"
	ObjectTypePreprint     = "sorg:ScholarlyArticle"
	ObjectTypeReview       = "sorg:Review"
	ObjectTypeOrganization = "Organization"
	ObjectTypeProject      = "Project"
)
