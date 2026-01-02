package router

// This file contains Swagger annotations for all proxied routes

// Auth Routes

// Login (Proxied) godoc
// @Summary User login (Proxied to Auth Service)
// @Description Authenticate user via Auth Service
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body object true "Login credentials"
// @Success 200 {object} map[string]interface{} "Login successful"
// @Failure 401 {object} map[string]interface{} "Authentication failed"
// @Router /api/v1/auth/login [post]
func swaggerLogin() {}

// Register (Proxied) godoc
// @Summary Register user (Proxied to Auth Service)
// @Description Register new user via Auth Service
// @Tags auth
// @Accept json
// @Produce json
// @Param user body object true "User registration"
// @Success 201 {object} map[string]interface{} "User registered"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Router /api/v1/auth/register [post]
func swaggerRegister() {}

// Refresh Token (Proxied) godoc
// @Summary Refresh access token (Proxied to Auth Service)
// @Description Refresh JWT access token via Auth Service
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh body object true "Refresh token"
// @Success 200 {object} map[string]interface{} "Token refreshed"
// @Failure 401 {object} map[string]interface{} "Invalid refresh token"
// @Router /api/v1/auth/refresh [post]
func swaggerRefreshToken() {}

// Logout (Proxied) godoc
// @Summary User logout (Proxied to Auth Service)
// @Description Logout user via Auth Service
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Logout successful"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/auth/logout [post]
func swaggerLogout() {}

// User Routes

// Get Users (Proxied) godoc
// @Summary List users (Proxied to User Service)
// @Description Get list of users via User Service
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} map[string]interface{} "List of users"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/users [get]
func swaggerGetUsers() {}

// Create User (Proxied) godoc
// @Summary Create user (Proxied to User Service)
// @Description Create new user via User Service
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param user body object true "User creation request"
// @Success 201 {object} map[string]interface{} "User created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/users [post]
func swaggerCreateUser() {}

// Get User (Proxied) godoc
// @Summary Get user by ID (Proxied to User Service)
// @Description Get user details via User Service
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /api/v1/users/{id} [get]
func swaggerGetUser() {}

// Update User (Proxied) godoc
// @Summary Update user (Proxied to User Service)
// @Description Update user details via User Service
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "User ID"
// @Param user body object true "User update request"
// @Success 200 {object} map[string]interface{} "User updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /api/v1/users/{id} [put]
func swaggerUpdateUser() {}

// Delete User (Proxied) godoc
// @Summary Delete user (Proxied to User Service)
// @Description Delete user via User Service
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Router /api/v1/users/{id} [delete]
func swaggerDeleteUser() {}

// Search Users (Proxied) godoc
// @Summary Search users (Proxied to User Service)
// @Description Search users via User Service
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param q query string false "Search query"
// @Success 200 {object} map[string]interface{} "Search results"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/users/search [get]
func swaggerSearchUsers() {}

// Tenant Routes

// Get Tenants (Proxied) godoc
// @Summary List tenants (Proxied to Tenant Service)
// @Description Get list of tenants via Tenant Service
// @Tags tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "List of tenants"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/tenants [get]
func swaggerGetTenants() {}

// Create Tenant (Proxied) godoc
// @Summary Create tenant (Proxied to Tenant Service)
// @Description Create new tenant via Tenant Service
// @Tags tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tenant body object true "Tenant creation request"
// @Success 201 {object} map[string]interface{} "Tenant created"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/tenants [post]
func swaggerCreateTenant() {}

// Get Tenant (Proxied) godoc
// @Summary Get tenant by ID (Proxied to Tenant Service)
// @Description Get tenant details via Tenant Service
// @Tags tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Success 200 {object} map[string]interface{} "Tenant details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Tenant not found"
// @Router /api/v1/tenants/{id} [get]
func swaggerGetTenant() {}

// Update Tenant (Proxied) godoc
// @Summary Update tenant (Proxied to Tenant Service)
// @Description Update tenant details via Tenant Service
// @Tags tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Param tenant body object true "Tenant update request"
// @Success 200 {object} map[string]interface{} "Tenant updated"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Tenant not found"
// @Router /api/v1/tenants/{id} [put]
func swaggerUpdateTenant() {}

// Delete Tenant (Proxied) godoc
// @Summary Delete tenant (Proxied to Tenant Service)
// @Description Delete tenant via Tenant Service
// @Tags tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Success 200 {object} map[string]interface{} "Tenant deleted"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Tenant not found"
// @Router /api/v1/tenants/{id} [delete]
func swaggerDeleteTenant() {}

// Add User to Tenant (Proxied) godoc
// @Summary Add user to tenant (Proxied to Tenant Service)
// @Description Add user to tenant via Tenant Service
// @Tags tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Param user body object true "User assignment request"
// @Success 200 {object} map[string]interface{} "User added to tenant"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Tenant or user not found"
// @Router /api/v1/tenants/{id}/users [post]
func swaggerAddUserToTenant() {}

// Remove User from Tenant (Proxied) godoc
// @Summary Remove user from tenant (Proxied to Tenant Service)
// @Description Remove user from tenant via Tenant Service
// @Tags tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID"
// @Param user_id path string true "User ID"
// @Success 200 {object} map[string]interface{} "User removed from tenant"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Tenant or user not found"
// @Router /api/v1/tenants/{id}/users/{user_id} [delete]
func swaggerRemoveUserFromTenant() {}

// Notification Routes

// Send Email (Proxied) godoc
// @Summary Send email (Proxied to Notification Service)
// @Description Send email notification via Notification Service
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param email body object true "Email request"
// @Success 200 {object} map[string]interface{} "Email sent"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/notifications/email [post]
func swaggerSendEmail() {}

// Send Webhook (Proxied) godoc
// @Summary Send webhook (Proxied to Notification Service)
// @Description Send webhook notification via Notification Service
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param webhook body object true "Webhook request"
// @Success 200 {object} map[string]interface{} "Webhook sent"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/notifications/webhook [post]
func swaggerSendWebhook() {}

// Get Notifications (Proxied) godoc
// @Summary List notifications (Proxied to Notification Service)
// @Description Get list of notifications via Notification Service
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} map[string]interface{} "List of notifications"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Router /api/v1/notifications [get]
func swaggerGetNotifications() {}

// Get Notification (Proxied) godoc
// @Summary Get notification by ID (Proxied to Notification Service)
// @Description Get notification details via Notification Service
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Notification ID"
// @Success 200 {object} map[string]interface{} "Notification details"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 404 {object} map[string]interface{} "Notification not found"
// @Router /api/v1/notifications/{id} [get]
func swaggerGetNotification() {}
