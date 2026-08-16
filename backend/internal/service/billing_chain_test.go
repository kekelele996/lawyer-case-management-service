package service

import (
	"testing"
	"time"

	"cylawcase/internal/constants"
	"cylawcase/internal/model"
	"cylawcase/internal/repository"
)

func TestBillingStatusAndSummaryChain(t *testing.T) {
	db := newLawDB(t)
	client := &model.Client{Name: "张三", IDNumber: "110101199001011234", Contact: "13800000001"}
	if err := db.Create(client).Error; err != nil {
		t.Fatalf("create client: %v", err)
	}
	cs := &model.Case{CaseNo: "CY20250001", Title: "合同纠纷", CaseType: constants.CaseTypeCivil, Status: constants.CaseStatusFiled, ClientID: client.ID, LeadLawyerID: 1, CoLawyerIDs: model.CoLawyerJSON("[]")}
	if err := db.Create(cs).Error; err != nil {
		t.Fatalf("create case: %v", err)
	}

	repo := repository.NewBillingRepository(db)
	svc := NewBillingService(repo, repository.NewCaseRepository(db), repository.NewClientRepository(db), lawLogger())

	b1, err := svc.Create(cs.ID, client.ID, constants.BillingTypeAttorneyFee, 100, "")
	if err != nil {
		t.Fatalf("Create b1: %v", err)
	}
	if _, err := svc.MarkPaid(b1.ID); err != nil {
		t.Fatalf("MarkPaid b1: %v", err)
	}
	if inv, err := svc.MarkInvoiced(b1.ID, "INV-001"); err != nil {
		t.Fatalf("MarkInvoiced b1: %v", err)
	} else if inv.Status != constants.BillingStatusInvoiced {
		t.Fatalf("b1 status = %s, want invoiced", inv.Status)
	}

	b2, err := svc.Create(cs.ID, client.ID, constants.BillingTypeCourtFee, 50, "")
	if err != nil {
		t.Fatalf("Create b2: %v", err)
	}
	if _, err := svc.Void(b2.ID); err != nil {
		t.Fatalf("Void b2: %v", err)
	}
	if _, err := svc.Void(b2.ID); err == nil {
		t.Fatal("double void should fail")
	}

	now := time.Now()
	if err := db.Create(&model.Billing{BillNo: "B3", BillingType: constants.BillingTypeOther, Amount: 200, Status: constants.BillingStatusInvoiced, CaseID: cs.ID, ClientID: client.ID, CreatedAt: now}).Error; err != nil {
		t.Fatalf("create invoiced billing: %v", err)
	}
	sum, err := svc.Summary()
	if err != nil {
		t.Fatalf("Summary: %v", err)
	}
	if sum["received"] != 300 {
		t.Fatalf("received = %v, want 300", sum["received"])
	}
}
