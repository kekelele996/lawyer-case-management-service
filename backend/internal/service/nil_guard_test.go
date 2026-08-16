package service

import (
	"errors"
	"testing"

	"cylawcase/internal/constants"
	"cylawcase/internal/model"
	"cylawcase/internal/repository"
)

func expectNoPanicErr(t *testing.T, name string, fn func() error) error {
	t.Helper()
	var err error
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("%s panicked: %v", name, r)
			}
		}()
		err = fn()
	}()
	return err
}

func TestMissingBillingAndCaseReturnNotFound(t *testing.T) {
	db := newLawDB(t)
	client := &model.Client{Name: "孙七", IDNumber: "110101195001011234", Contact: "13500000005"}
	if err := db.Create(client).Error; err != nil {
		t.Fatalf("create client: %v", err)
	}
	billingRepo := repository.NewBillingRepository(db)
	caseRepo := repository.NewCaseRepository(db)
	billingSvc := NewBillingService(billingRepo, caseRepo, repository.NewClientRepository(db), lawLogger())
	caseSvc := NewCaseService(caseRepo, repository.NewClientRepository(db), repository.NewUserRepository(db), lawLogger())

	err := expectNoPanicErr(t, "MarkPaid", func() error {
		_, err := billingSvc.MarkPaid(999)
		return err
	})
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("MarkPaid missing = %v, want ErrNotFound", err)
	}

	err = expectNoPanicErr(t, "ChangeStatus", func() error {
		_, err := caseSvc.ChangeStatus(999, constants.RoleAdmin, constants.CaseStatusClosed)
		return err
	})
	if !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("ChangeStatus missing = %v, want ErrNotFound", err)
	}
}
