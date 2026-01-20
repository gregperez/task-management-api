package domain

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestUser_RoleChecks(t *testing.T) {
	tests := []struct {
		name     string
		role     UserRole
		isAdmin  bool
		isExecutor bool
		isAuditor bool
	}{
		{
			name:     "Admin user",
			role:     RoleAdmin,
			isAdmin:  true,
			isExecutor: false,
			isAuditor: false,
		},
		{
			name:     "Executor user",
			role:     RoleExecutor,
			isAdmin:  false,
			isExecutor: true,
			isAuditor: false,
		},
		{
			name:     "Auditor user",
			role:     RoleAuditor,
			isAdmin:  false,
			isExecutor: false,
			isAuditor: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{Role: tt.role}
			assert.Equal(t, tt.isAdmin, user.IsAdmin())
			assert.Equal(t, tt.isExecutor, user.IsExecutor())
			assert.Equal(t, tt.isAuditor, user.IsAuditor())
		})
	}
}

func TestUser_MustChangePassword(t *testing.T) {
	assert := assert.New(t)

	userWithTempPassword := &User{IsTemporaryPassword: true}
	assert.True(userWithTempPassword.MustChangePassword(), "User with temporary password should require password change")
	userWithPermanentPassword := &User{IsTemporaryPassword: false}
	assert.False(userWithPermanentPassword.MustChangePassword(), "User with permanent password should not require password change")
}