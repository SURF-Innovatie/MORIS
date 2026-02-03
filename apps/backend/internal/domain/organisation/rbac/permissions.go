package rbac

type Permission string

const (
	PermissionManageMembers           Permission = "manage_members"
	PermissionManageProjectRoles      Permission = "manage_project_roles"
	PermissionManageOrganisationRoles Permission = "manage_organisation_roles"
	PermissionManageCustomFields      Permission = "manage_custom_fields"
	PermissionManageDetails           Permission = "manage_details"
	PermissionCreateProject           Permission = "create_project"
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
}

var AllPermissions = []Permission{
	PermissionManageMembers,
	PermissionManageProjectRoles,
	PermissionManageOrganisationRoles,
	PermissionManageCustomFields,
	PermissionManageDetails,
	PermissionCreateProject,
}

func (p Permission) String() string {
	return string(p)
}
