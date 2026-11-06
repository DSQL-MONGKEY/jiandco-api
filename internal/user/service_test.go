package user

import (
	"testing"

	"jianco-api/internal/models"
)

func TestCanManage(t *testing.T) {
	tests := []struct {
		name        string
		actorLevel  int
		targetLevel int
		want        bool
	}{
		{name: "developer manages developer", actorLevel: 0, targetLevel: 0, want: true},
		{name: "developer manages staff", actorLevel: 0, targetLevel: 3, want: true},
		{name: "super admin manages admin", actorLevel: 1, targetLevel: 2, want: true},
		{name: "super admin cannot manage developer", actorLevel: 1, targetLevel: 0, want: false},
		{name: "admin manages staff", actorLevel: 2, targetLevel: 3, want: true},
		{name: "admin cannot manage admin", actorLevel: 2, targetLevel: 2, want: false},
		{name: "staff cannot manage staff", actorLevel: 3, targetLevel: 3, want: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actor := &models.User{Role: models.Role{Level: test.actorLevel}}
			if got := canManage(actor, test.targetLevel); got != test.want {
				t.Fatalf("canManage() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestCanManageNilActor(t *testing.T) {
	if canManage(nil, 3) {
		t.Fatal("nil actor must not manage users")
	}
}
