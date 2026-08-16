package service

import (
	"testing"

	"cylawcase/internal/constants"
	"cylawcase/internal/model"
	"cylawcase/internal/repository"
)

func TestCaseStatusAndListChain(t *testing.T) {
	db := newLawDB(t)
	client := &model.Client{Name: "李四", IDNumber: "110101198001011234", Contact: "13900000002"}
	if err := db.Create(client).Error; err != nil {
		t.Fatalf("create client: %v", err)
	}
	c := &model.Case{CaseNo: "CY20250002", Title: "劳动仲裁", CaseType: constants.CaseTypeLabor, Status: constants.CaseStatusFiled, ClientID: client.ID, LeadLawyerID: 2, CoLawyerIDs: model.CoLawyerJSON("[]")}
	if err := db.Create(c).Error; err != nil {
		t.Fatalf("create case: %v", err)
	}
	repo := repository.NewCaseRepository(db)
	svc := NewCaseService(repo, repository.NewClientRepository(db), repository.NewUserRepository(db), lawLogger())

	if _, err := svc.ChangeStatus(c.ID, "lawyer", constants.CaseStatusInvestigating); err != nil {
		t.Fatalf("forward transition: %v", err)
	}
	if _, err := svc.ChangeStatus(c.ID, "lawyer", constants.CaseStatusFiled); err != nil {
		t.Fatalf("backward transition: %v", err)
	}
	if _, err := svc.ChangeStatus(c.ID, constants.RoleAdmin, constants.CaseStatusClosed); err != nil {
		t.Fatalf("admin transition: %v", err)
	}

	list, err := repo.ListByLawyer(2)
	if err != nil {
		t.Fatalf("ListByLawyer: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("ListByLawyer = %d, want 1", len(list))
	}
}
