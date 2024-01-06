package auth

import (
	"context"

	"github.com/adriancrafter/todoapp/internal/am"
)

type (
	Service interface {
		am.Core

		// TenantVM

		//GetTenants(ctx context.Context) (vm.TenantsVM, error)
		//GetTenant(ctx context.Context, slug string) (vm.TenantVM, error)
		//CreateTenant(ctx context.Context, table vm.TenantVM) error
		//UpdateTenant(ctx context.Context, table vm.TenantVM) error
		//SoftDeleteTenant(ctx context.Context, critearia vm.CriteriaVM) error
		//DeleteTenant(ctx context.Context, criteria vm.CriteriaVM) error
		//PurgeTenants(ctx context.Context, criteria vm.CriteriaVM) error
		//
		//// AffiliationVM
		//
		//GetAffiliations(ctx context.Context) (vm.AffiliationsVM, error)
		//GetAffiliated(ctx context.Context, criteria vm.CriteriaVM) (vm.UsersVM, error)
		//GetAffiliatedTo(ctx context.Context, criteria vm.CriteriaVM) (vm.TenantsVM, error)
		//
		//// Service
		//
		//GetUsers(ctx context.Context, criteria vm.CriteriaVM) (vm.UsersVM, error)
		//GetUser(ctx context.Context, criteria vm.CriteriaVM) (vm.UserVM, error)
		//CreateUser(ctx context.Context, table vm.UserVM) error
		//UpdateUser(ctx context.Context, table vm.UserVM) error
		//SoftDeleteUser(ctx context.Context, critearia vm.CriteriaVM) error
		//DeleteUser(ctx context.Context, criteria vm.CriteriaVM) error
		//PurgeUsers(ctx context.Context, criteria vm.CriteriaVM) error
		//
		//// Account
		//
		//GetAccounts(ctx context.Context, criteria vm.CriteriaVM) (vm.AccountsVM, error)
		//GetAccount(ctx context.Context, criteria vm.CriteriaVM) (vm.AccountVM, error)
		//CreateAccount(ctx context.Context, table vm.AccountVM) error
		//UpdateAccount(ctx context.Context, table vm.AccountVM) error
		//SoftDeleteAccount(ctx context.Context, criteria vm.CriteriaVM) error
		//DeleteAccount(ctx context.Context, criteria vm.CriteriaVM) error
		//PurgeAccounts(ctx context.Context, criteria vm.CriteriaVM) error
		//
		//// Role
		//
		//GetRoles(ctx context.Context, criteria vm.CriteriaVM) (vm.RolesVM, error)
		//GetRole(ctx context.Context, criteria vm.CriteriaVM) (vm.RoleVM, error)
		//CreateRole(ctx context.Context, table vm.RoleVM) error
		//UpdateRole(ctx context.Context, table vm.RoleVM) error
		//SoftDeleteRole(ctx context.Context, criteria vm.CriteriaVM) error
		//DeleteRole(ctx context.Context, criteria vm.CriteriaVM) error
		//PurgeRoles(ctx context.Context, criteria vm.CriteriaVM) error
		//
		//// Permission
		//
		//GetPermissions(ctx context.Context, criteria vm.CriteriaVM) (vm.PermissionsVM, error)
		//GetPermission(ctx context.Context, criteria vm.CriteriaVM) (vm.PermissionVM, error)
		//CreatePermission(ctx context.Context, table vm.PermissionVM) error
		//UpdatePermission(ctx context.Context, table vm.PermissionVM) error
		//SoftDeletePermission(ctx context.Context, criteria vm.CriteriaVM) error
		//DeletePermission(ctx context.Context, criteria vm.CriteriaVM) error
		//PurgePermissions(ctx context.Context, criteria vm.CriteriaVM) error
		//
		//// Resource
		//
		//GetResources(ctx context.Context, criteria vm.CriteriaVM) (vm.ResourcesVM, error)
		//GetResource(ctx context.Context, criteria vm.CriteriaVM) (vm.ResourceVM, error)
		//CreateResource(ctx context.Context, table vm.ResourceVM) error
		//UpdateResource(ctx context.Context, table vm.ResourceVM) error
		//SoftDeleteResource(ctx context.Context, criteria vm.CriteriaVM) error
		//DeleteResource(ctx context.Context, criteria vm.CriteriaVM) error
		//PurgeResources(ctx context.Context, criteria vm.CriteriaVM) error
		//
		//// ResourcePermission
		//
		////GetResourcePermissions(ctx context.Context, criteria vm.CriteriaVM) (vm.ResourcePermissions, error)
		////GetResourcePermission(ctx context.Context, criteria vm.CriteriaVM) (vm.ResourcePermission, error)
		////CreateResourcePermission(ctx context.Context, table vm.ResourcePermission) error
		////UpdateResourcePermission(ctx context.Context, table vm.ResourcePermission) error
		////SoftDeleteResourcePermission(ctx context.Context, criteria vm.CriteriaVM) error
		////DeleteResourcePermission(ctx context.Context, criteria vm.CriteriaVM) error
		////PurgeResourcePermissions(ctx context.Context, criteria vm.CriteriaVM) error
		//
		//// RolePermission
		//
		////GetRolePermissions(ctx context.Context, criteria vm.CriteriaVM) (vm.RolePermissions, error)
		////GetRolePermission(ctx context.Context, criteria vm.CriteriaVM) (vm.RolePermission, error)
		////CreateRolePermission(ctx context.Context, table vm.RolePermission) error
		////UpdateRolePermission(ctx context.Context, table vm.RolePermission) error
		////SoftDeleteRolePermission(ctx context.Context, criteria vm.CriteriaVM) error
		////DeleteRolePermission(ctx context.Context, criteria vm.CriteriaVM) error
		////PurgeRolePermissions(ctx context.Context, criteria vm.CriteriaVM) error
		//
		//// AccountRole
		//
		////GetAccountRoles(ctx context.Context, criteria vm.CriteriaVM) (vm.AccountRoles, error)
		////GetAccountRole(ctx context.Context, criteria vm.CriteriaVM) (vm.AccountRole, error)
		////CreateAccountRole(ctx context.Context, table vm.AccountRole) error
		////UpdateAccountRole(ctx context.Context, table vm.AccountRole) error
		////SoftDeleteAccountRole(ctx context.Context, criteria vm.CriteriaVM) error
		////DeleteAccountRole(ctx context.Context, criteria vm.CriteriaVM) error
		////PurgeAccountRoles(ctx context.Context, criteria vm.CriteriaVM) error
		//
		//// AccountPermission
		//
		////GetAccountPermissions(ctx context.Context, criteria vm.CriteriaVM) (vm.AccountPermissions, error)
		////GetAccountPermission(ctx context.Context, criteria vm.CriteriaVM) (vm.AccountPermission, error)
		////CreateAccountPermission(ctx context.Context, table vm.AccountPermission) error
		////UpdateAccountPermission(ctx context.Context, table vm.AccountPermission) error
		////SoftDeleteAccountPermission(ctx context.Context, criteria vm.CriteriaVM) error
		////DeleteAccountPermission(ctx context.Context, criteria vm.CriteriaVM) error
		////PurgeAccountPermissions(ctx context.Context, criteria vm.CriteriaVM) error
		//
		//// Service related
		//
		//SignUpUser(user *model.UserDA) (am.ValErrorSet, error)
		//ConfirmUser(slug, token string) error
		//SignInUser(username, password string) (user vm.Service, err error)
		SignInUser(ctc context.Context, si SigninVM) (user UserVM, err error)
	}
)

type (
	Repo interface {
		GetUser(ctx context.Context, userID string) (user User, err error)
		SignIn(ctx context.Context, sivm Signin) (ua UserAuth, err error)
	}
)
