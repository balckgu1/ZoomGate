package auth

import "zoomgate/internal/model"

// Permission defines what actions a role can perform.
type Permission struct {
	ManageUsers     bool
	ManageProviders bool
	ManagePolicies  bool
	ViewAuditLogs   bool
	ViewDashboard   bool
	UseProxy        bool
}

var rolePermissions = map[model.Role]Permission{
	model.RoleAdmin: {
		ManageUsers:     true,
		ManageProviders: true,
		ManagePolicies:  true,
		ViewAuditLogs:   true,
		ViewDashboard:   true,
		UseProxy:        true,
	},
	model.RoleUser: {
		ManageUsers:     false,
		ManageProviders: false,
		ManagePolicies:  false,
		ViewAuditLogs:   false,
		ViewDashboard:   true,
		UseProxy:        true,
	},
	model.RoleViewer: {
		ManageUsers:     false,
		ManageProviders: false,
		ManagePolicies:  false,
		ViewAuditLogs:   true,
		ViewDashboard:   true,
		UseProxy:        false,
	},
}

// GetPermissions returns the permissions for a given user role.
func GetPermissions(role model.Role) Permission {
	userPermission, ok := rolePermissions[role]
	if ok {
		return userPermission
	}
	return Permission{}
}

// HasRole checks if a user has a required role.
func HasRole(userRole string, required string) bool {
	// admin has all roles
	if userRole == string(model.RoleAdmin) {
		return true
	}
	// only userRole matches required role, return true
	return userRole == required
}
