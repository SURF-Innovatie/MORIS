package odrl_test

import (
	"testing"

	"github.com/SURF-Innovatie/MORIS/pkg/odrl"
	"github.com/stretchr/testify/assert"
)

func TestSimpleSetPolicy(t *testing.T) {
	// Example from ODRL Spec 2.1.1 Set Class
	policy := odrl.NewSet("http://example.com/policy:1010").
		AddPermission(odrl.NewPermission(odrl.ActionUse).
			WithTarget("http://example.com/asset:9898.movie"))

	jsonBytes, err := policy.ToJSON()
	assert.NoError(t, err)

	expectedJSON := `{"@context":"http://www.w3.org/ns/odrl.jsonld","@type":"Set","uid":"http://example.com/policy:1010","permission":[{"action":"use","target":"http://example.com/asset:9898.movie"}]}`
	assert.JSONEq(t, expectedJSON, string(jsonBytes))
}

func TestOfferPolicyWithConstraint(t *testing.T) {
	// Example from ODRL Spec 2.6.1 Permission Class
	rightOp := map[string]string{
		"@value": "2017-12-31",
		"@type":  "xsd:date",
	}

	policy := odrl.NewOffer("http://example.com/policy:9090").
		WithProfile("http://example.com/odrl:profile:07").
		AddPermission(odrl.NewPermission(odrl.ActionPlay).
			WithTarget("http://example.com/game:9090").
			WithAssigner("http://example.com/org:xyz").
			WithConstraint(odrl.NewConstraint("dateTime", odrl.OperatorLteq, rightOp)))

	jsonBytes, err := policy.ToJSON()
	assert.NoError(t, err)

	expectedJSON := `{
		"@context": "http://www.w3.org/ns/odrl.jsonld",
		"@type": "Offer",
		"uid": "http://example.com/policy:9090",
		"profile": "http://example.com/odrl:profile:07",
		"permission": [{
			"target": "http://example.com/game:9090",
			"assigner": "http://example.com/org:xyz",
			"action": "play",
			"constraint": [{
				"leftOperand": "dateTime",
				"operator": "lteq",
				"rightOperand": { "@value": "2017-12-31", "@type": "xsd:date" }
			}]
		}]
	}`
	assert.JSONEq(t, expectedJSON, string(jsonBytes))
}

func TestAgreementWithPartyObjects(t *testing.T) {
	// Example with Party objects instead of IRIs
	assigner := odrl.NewParty("http://example.com/org:sony-music")
	assignee := odrl.NewParty("http://example.com/people/billie")

	policy := odrl.NewAgreement("http://example.com/policy:8888").
		AddPermission(odrl.NewPermission(odrl.ActionPlay).
			WithTarget("http://example.com/music/1999.mp3").
			WithAssignerParty(assigner).
			WithAssigneeParty(assignee))

	jsonBytes, err := policy.ToJSON()
	assert.NoError(t, err)

	expectedJSON := `{
		"@context": "http://www.w3.org/ns/odrl.jsonld",
		"@type": "Agreement",
		"uid": "http://example.com/policy:8888",
		"permission": [{
			"action": "play",
			"target": "http://example.com/music/1999.mp3",
			"assigner": { "uid": "http://example.com/org:sony-music", "@type": "Party" },
			"assignee": { "uid": "http://example.com/people/billie", "@type": "Party" }
		}]
	}`
	assert.JSONEq(t, expectedJSON, string(jsonBytes))
}

func TestProhibitionWithDuty(t *testing.T) {
	// Test a Prohibition rule
	policy := odrl.NewSet("http://example.com/policy:2020").
		AddProhibition(odrl.NewProhibition(odrl.ActionDistribute).
			WithTarget("http://example.com/document:secret"))

	jsonBytes, err := policy.ToJSON()
	assert.NoError(t, err)

	expectedJSON := `{
		"@context": "http://www.w3.org/ns/odrl.jsonld",
		"@type": "Set",
		"uid": "http://example.com/policy:2020",
		"prohibition": [{
			"action": "distribute",
			"target": "http://example.com/document:secret"
		}]
	}`
	assert.JSONEq(t, expectedJSON, string(jsonBytes))
}

func TestObligation(t *testing.T) {
	// Test policy-level duties (obligations)
	policy := odrl.NewSet("http://example.com/policy:3030").
		AddObligation(odrl.NewDuty(odrl.ActionAttribution).
			WithTarget("http://example.com/creator:john"))

	jsonBytes, err := policy.ToJSON()
	assert.NoError(t, err)

	expectedJSON := `{
		"@context": "http://www.w3.org/ns/odrl.jsonld",
		"@type": "Set",
		"uid": "http://example.com/policy:3030",
		"obligation": [{
			"action": "attribution",
			"target": "http://example.com/creator:john"
		}]
	}`
	assert.JSONEq(t, expectedJSON, string(jsonBytes))
}
