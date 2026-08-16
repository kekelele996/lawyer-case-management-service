package service

import (
	"errors"
	"testing"

	"cylawcase/internal/constants"
	"cylawcase/internal/model"
	"cylawcase/internal/repository"
)

func TestErrorChainSentinels(t *testing.T) {
	db := newLawDB(t)
	client := &model.Client{Name: "赵六", IDNumber: "110101196001011234", Contact: "13600000004"}
	if err := db.Create(client).Error; err != nil {
		t.Fatalf("create client: %v", err)
	}
	billingRepo := repository.NewBillingRepository(db)
	caseRepo := repository.NewCaseRepository(db)
	clientRepo := repository.NewClientRepository(db)

	if _, err := billingRepo.FindByID(999); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("billing FindByID missing = %v, want ErrNotFound", err)
	}
	if _, err := caseRepo.FindByID(999); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("case FindByID missing = %v, want ErrNotFound", err)
	}
	if _, err := clientRepo.FindByID(999); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("client FindByID missing = %v, want ErrNotFound", err)
	}

	svc := NewBillingService(billingRepo, caseRepo, clientRepo, lawLogger())
	if _, err := svc.Create(999, client.ID, constants.BillingTypeAttorneyFee, 100, ""); !errors.Is(err, repository.ErrNotFound) {
		t.Fatalf("Create with missing case = %v, want ErrNotFound", err)
	}
}
