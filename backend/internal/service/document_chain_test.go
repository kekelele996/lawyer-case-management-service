package service

import (
	"testing"

	"cylawcase/internal/constants"
	"cylawcase/internal/model"
	"cylawcase/internal/repository"
)

func TestDocumentUploadAndOrderChain(t *testing.T) {
	db := newLawDB(t)
	client := &model.Client{Name: "王五", IDNumber: "110101197001011234", Contact: "13700000003"}
	if err := db.Create(client).Error; err != nil {
		t.Fatalf("create client: %v", err)
	}
	c := &model.Case{CaseNo: "CY20250003", Title: "交通事故", CaseType: constants.CaseTypeCivil, Status: constants.CaseStatusFiled, ClientID: client.ID, LeadLawyerID: 1, CoLawyerIDs: model.CoLawyerJSON("[]")}
	if err := db.Create(c).Error; err != nil {
		t.Fatalf("create case: %v", err)
	}
	repo := repository.NewDocumentRepository(db)
	svc := NewDocumentService(repo, repository.NewCaseRepository(db), lawLogger())

	d1, err := svc.Create(c.ID, 1, "起诉状", constants.DocTypeComplaint, "/files/1.pdf")
	if err != nil {
		t.Fatalf("Create d1: %v", err)
	}
	d2, err := svc.Create(c.ID, 1, "证据清单", constants.DocTypeEvidence, "/files/2.pdf")
	if err != nil {
		t.Fatalf("Create d2: %v", err)
	}

	list, err := svc.ListByCase(c.ID)
	if err != nil {
		t.Fatalf("ListByCase: %v", err)
	}
	if len(list) != 2 || list[0].ID != d2.ID {
		t.Fatalf("ListByCase order wrong: %+v", list)
	}
	_ = d1
}
