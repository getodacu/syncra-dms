// Package api Syncra DMS API
//
// Syncra DMS backend API.
//
// Schemes: http
// Host: localhost:8080
// BasePath: /
// Version: 0.1
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// swagger:meta
package api

func swaggerOperations() {
	// swagger:operation GET /healthz system getHealth
	//
	// Health check.
	//
	// ---
	// responses:
	//   "200":
	//     description: API process is running.

	// swagger:operation GET /readyz system getReadiness
	//
	// Readiness check.
	//
	// ---
	// responses:
	//   "200":
	//     description: API dependencies are ready.
	//   "503":
	//     description: API dependencies are not ready.

	// swagger:operation GET /version system getVersion
	//
	// Version metadata.
	//
	// ---
	// responses:
	//   "200":
	//     description: API version metadata.

	// swagger:operation POST /api/auth/sign-up/email auth signUpEmail
	//
	// Sign up with email and password.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: User account was created or already exists and email verification is required.
	//   "400":
	//     description: Invalid sign-up request.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/sign-in/email auth signInEmail
	//
	// Sign in with email and password.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Authenticated session.
	//   "401":
	//     description: Invalid email or password, or trusted internal request required.
	//   "403":
	//     description: Email is not verified.

	// swagger:operation GET /api/auth/get-session auth getSession
	//
	// Load the current session from the auth.session_token cookie.
	//
	// Trusted SvelteKit server hook endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Current authenticated session, or null when no valid session exists.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/sign-out auth signOut
	//
	// Sign out the current session.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Session was deleted when present.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/email-otp/send-verification-otp auth sendVerificationOTP
	//
	// Send or rotate a six-digit email verification OTP.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: OTP was created when the account exists and is unverified.
	//   "400":
	//     description: Invalid OTP request.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/email-otp/verify-email auth verifyEmailOTP
	//
	// Confirm an email address with a six-digit OTP.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Email was verified.
	//   "400":
	//     description: Invalid or expired verification code.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/password-reset/request auth requestPasswordReset
	//
	// Request a password reset email.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Reset request accepted.
	//   "400":
	//     description: Invalid password reset request.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/password-reset/confirm auth confirmPasswordReset
	//
	// Confirm password reset with emailed token.
	//
	// Trusted SvelteKit server action endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Password was reset and existing sessions were revoked.
	//   "400":
	//     description: Invalid or expired password reset token.
	//   "401":
	//     description: Trusted internal request required.

	// swagger:operation POST /api/auth/oauth/google/start auth startGoogleOAuth
	//
	// Start Google OAuth.
	//
	// Trusted SvelteKit server endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Google authorization URL and state.
	//   "401":
	//     description: Trusted internal request required.
	//   "503":
	//     description: Google OAuth is not configured.

	// swagger:operation POST /api/auth/oauth/google/callback auth signInGoogleOAuth
	//
	// Complete Google OAuth sign-in.
	//
	// Trusted SvelteKit server endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Authenticated session.
	//   "401":
	//     description: OAuth sign-in failed or trusted internal request required.

	// swagger:operation POST /api/auth/oauth/github/start auth startGitHubOAuth
	//
	// Start GitHub OAuth.
	//
	// Trusted SvelteKit server endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: GitHub authorization URL and state.
	//   "401":
	//     description: Trusted internal request required.
	//   "503":
	//     description: GitHub OAuth is not configured.

	// swagger:operation POST /api/auth/oauth/github/callback auth signInGitHubOAuth
	//
	// Complete GitHub OAuth sign-in.
	//
	// Trusted SvelteKit server endpoint.
	//
	// ---
	// responses:
	//   "200":
	//     description: Authenticated session.
	//   "401":
	//     description: OAuth sign-in failed or trusted internal request required.

	// swagger:operation GET /api/organization-units/tree organizationUnits listOrganizationUnits
	//
	// List active organization units as a hierarchy.
	//
	// Trusted SvelteKit server endpoint. Requires an authenticated session.
	//
	// ---
	// responses:
	//   "200":
	//     description: Active organization unit hierarchy.
	//   "401":
	//     description: Authenticated session or trusted internal request required.

	// swagger:operation GET /api/organization-units/archived organizationUnits listArchivedOrganizationUnits
	//
	// List archived organization units.
	//
	// Trusted SvelteKit server endpoint. Requires an admin session.
	//
	// ---
	// responses:
	//   "200":
	//     description: Archived organization units.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: Admin role required.

	// swagger:operation POST /api/organization-units organizationUnits createOrganizationUnit
	//
	// Create a root or child organization unit.
	//
	// Trusted SvelteKit server endpoint. Requires an admin session.
	//
	// ---
	// responses:
	//   "201":
	//     description: Organization unit created.
	//   "400":
	//     description: Invalid organization unit request.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: Admin role required.
	//   "404":
	//     description: Parent organization unit not found.
	//   "409":
	//     description: Active organization unit code already exists.

	// swagger:operation PATCH /api/organization-units/{id} organizationUnits updateOrganizationUnit
	//
	// Update an active organization unit's details.
	//
	// Trusted SvelteKit server endpoint. Requires an admin session.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: Organization unit updated.
	//   "400":
	//     description: Invalid organization unit request.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: Admin role required.
	//   "404":
	//     description: Organization unit not found.
	//   "409":
	//     description: Active organization unit code already exists.

	// swagger:operation PATCH /api/organization-units/{id}/parent organizationUnits moveOrganizationUnit
	//
	// Move an active organization unit to another parent or root.
	//
	// Trusted SvelteKit server endpoint. Requires an admin session.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: Organization unit moved.
	//   "400":
	//     description: Invalid move request.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: Admin role required.
	//   "404":
	//     description: Organization unit or parent not found.
	//   "409":
	//     description: Move would create a cycle.

	// swagger:operation POST /api/organization-units/{id}/archive organizationUnits archiveOrganizationUnit
	//
	// Archive an active organization unit and its descendants.
	//
	// Trusted SvelteKit server endpoint. Requires an admin session.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: Organization unit subtree archived.
	//   "400":
	//     description: Invalid organization unit id.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: Admin role required.
	//   "404":
	//     description: Organization unit not found.

	// swagger:operation GET /api/users users listUsers
	//
	// List active users.
	//
	// Trusted SvelteKit server endpoint. Requires user.view.
	//
	// ---
	// responses:
	//   "200":
	//     description: Users were listed.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.view permission required.

	// swagger:operation GET /api/users/{id} users getUser
	//
	// Get an active user.
	//
	// Trusted SvelteKit server endpoint. Requires user.view.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: User was loaded.
	//   "400":
	//     description: Invalid user id.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.view permission required.
	//   "404":
	//     description: User not found.

	// swagger:operation POST /api/users users createUser
	//
	// Create an invited or active user.
	//
	// Trusted SvelteKit server endpoint. Requires user.create.
	//
	// ---
	// responses:
	//   "201":
	//     description: User was created.
	//   "400":
	//     description: Invalid user request.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.create permission required.
	//   "409":
	//     description: User already exists.

	// swagger:operation PATCH /api/users/{id} users updateUser
	//
	// Update user profile fields.
	//
	// Trusted SvelteKit server endpoint. Requires user.update.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: User was updated.
	//   "400":
	//     description: Invalid user request.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.update permission required.
	//   "404":
	//     description: User not found.

	// swagger:operation POST /api/users/{id}/activate users activateUser
	//
	// Activate a user.
	//
	// Trusted SvelteKit server endpoint. Requires user.activate.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: User was activated.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.activate permission required.
	//   "404":
	//     description: User not found.

	// swagger:operation POST /api/users/{id}/deactivate users deactivateUser
	//
	// Deactivate a user.
	//
	// Trusted SvelteKit server endpoint. Requires user.update.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: User was deactivated.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.update permission required.
	//   "404":
	//     description: User not found.

	// swagger:operation POST /api/users/{id}/suspend users suspendUser
	//
	// Suspend a user and revoke sessions.
	//
	// Trusted SvelteKit server endpoint. Requires user.suspend.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: User was suspended.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.suspend permission required.
	//   "404":
	//     description: User not found.

	// swagger:operation DELETE /api/users/{id} users deleteUser
	//
	// Soft-delete a user and revoke sessions.
	//
	// Trusted SvelteKit server endpoint. Requires user.delete.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: User was deleted.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.delete permission required.
	//   "404":
	//     description: User not found.

	// swagger:operation POST /api/users/{id}/primary-organization-unit users setUserPrimaryOrganizationUnit
	//
	// Assign or clear a user's primary organization unit.
	//
	// Trusted SvelteKit server endpoint. Requires user.assign_unit.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: Primary organization unit was updated.
	//   "400":
	//     description: Invalid organization unit request.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.assign_unit permission required.
	//   "404":
	//     description: User or organization unit not found.

	// swagger:operation POST /api/users/{id}/roles users assignUserRole
	//
	// Assign a scoped role to a user.
	//
	// Trusted SvelteKit server endpoint. Requires user.assign_role.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "201":
	//     description: Role assignment was created.
	//   "400":
	//     description: Invalid role assignment request.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.assign_role permission required.
	//   "404":
	//     description: User, role, or organization unit not found.

	// swagger:operation DELETE /api/users/{id}/roles/{assignmentId} users removeUserRole
	//
	// Remove a user role assignment.
	//
	// Trusted SvelteKit server endpoint. Requires user.assign_role.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// - name: assignmentId
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: Role assignment was removed.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.assign_role permission required.
	//   "404":
	//     description: Role assignment not found.

	// swagger:operation POST /api/users/{id}/groups users addUserGroup
	//
	// Add a user to a group.
	//
	// Trusted SvelteKit server endpoint. Requires user.assign_group.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "201":
	//     description: Group membership was created.
	//   "400":
	//     description: Invalid group assignment request.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.assign_group permission required.
	//   "404":
	//     description: User or group not found.

	// swagger:operation DELETE /api/users/{id}/groups/{groupId} users removeUserGroup
	//
	// Remove a user from a group.
	//
	// Trusted SvelteKit server endpoint. Requires user.assign_group.
	//
	// ---
	// parameters:
	// - name: id
	//   in: path
	//   required: true
	//   type: string
	// - name: groupId
	//   in: path
	//   required: true
	//   type: string
	// responses:
	//   "200":
	//     description: Group membership was removed.
	//   "401":
	//     description: Authenticated session or trusted internal request required.
	//   "403":
	//     description: user.assign_group permission required.
	//   "404":
	//     description: Group membership not found.
}
