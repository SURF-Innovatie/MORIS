package role

type Permission string

const (
	PermissionManageMembers           Permission = "manage_members"
	PermissionManageProjectRoles      Permission = "manage_project_roles"
	PermissionManageOrganisationRoles Permission = "manage_organisation_roles"
	PermissionManageCustomFields      Permission = "manage_custom_fields"
	PermissionManageDetails           Permission = "manage_details"
)

var AllPermissions = []Permission{
	PermissionManageMembers,
	PermissionManageProjectRoles,
	PermissionManageOrganisationRoles,
	PermissionManageCustomFields,
	PermissionManageDetails,
}

func (p Permission) String() string {
	return string(p)
}
