package role

type Permission string

const (
	PermissionManageMembers           Permission = "manage_members"
	PermissionManageProjectRoles      Permission = "manage_project_roles"
	PermissionManageOrganisationRoles Permission = "manage_organisation_roles"
	PermissionManageCustomFields      Permission = "manage_custom_fields"
	PermissionManageDetails           Permission = "manage_details"
	PermissionCreateProject           Permission = "create_project"

	// Budget permissions
	PermissionViewBudget    Permission = "budget_view"
	PermissionEditBudget    Permission = "budget_edit"
	PermissionApproveBudget Permission = "budget_approve"
	PermissionRecordActuals Permission = "budget_record_actuals"

	// Analytics permissions
	PermissionViewAnalytics Permission = "analytics_view"
	PermissionAccessOData   Permission = "odata_access"
)

type PermissionDefinition struct {
	Permission  Permission
	Label       string
	Description string
}

var Definitions = []PermissionDefinition{
	{Permission: PermissionManageMembers, Label: "Manage Members", Description: "Can invite, remove, and manage members"},
	{Permission: PermissionManageProjectRoles, Label: "Manage Project Roles", Description: "Can create and manage project-level roles"},
	{Permission: PermissionManageOrganisationRoles, Label: "Manage Organisation Roles", Description: "Can create and manage organisation-level roles"},
	{Permission: PermissionManageCustomFields, Label: "Manage Custom Fields", Description: "Can manage custom fields for projects and people"},
	{Permission: PermissionManageDetails, Label: "Manage Details", Description: "Can update organisation details"},
	{Permission: PermissionCreateProject, Label: "Create Project", Description: "Can create new projects"},
	// Budget permissions
	{Permission: PermissionViewBudget, Label: "View Budget", Description: "Can view project budgets and expenditures"},
	{Permission: PermissionEditBudget, Label: "Edit Budget", Description: "Can create and modify budget line items"},
	{Permission: PermissionApproveBudget, Label: "Approve Budget", Description: "Can approve and lock budgets"},
	{Permission: PermissionRecordActuals, Label: "Record Actuals", Description: "Can record actual expenditures"},
	// Analytics permissions
	{Permission: PermissionViewAnalytics, Label: "View Analytics", Description: "Can view organisation-level analytics dashboards"},
	{Permission: PermissionAccessOData, Label: "OData Access", Description: "Can access OData endpoints for Power BI integration"},
}

var AllPermissions = []Permission{
	PermissionManageMembers,
	PermissionManageProjectRoles,
	PermissionManageOrganisationRoles,
	PermissionManageCustomFields,
	PermissionManageDetails,
	PermissionCreateProject,
	PermissionViewBudget,
	PermissionEditBudget,
	PermissionApproveBudget,
	PermissionRecordActuals,
	PermissionViewAnalytics,
	PermissionAccessOData,
}

func (p Permission) String() string {
	return string(p)
}
