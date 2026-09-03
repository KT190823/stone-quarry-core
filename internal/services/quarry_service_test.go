package services

import (
	"encoding/json"
	"math"
	"testing"

	"mo-da-backend/internal/models"
)

func TestReconciliationFormula(t *testing.T) {
	extractedVolume := 35000.0 // m3
	density := 1.55            // ton/m3
	actualWeight := 53820.0    // ton

	expectedWeight := extractedVolume * density
	if expectedWeight != 54250.0 {
		t.Fatalf("Expected weight = 54250.0, got %f", expectedWeight)
	}

	differenceTon := expectedWeight - actualWeight
	if differenceTon != 430.0 {
		t.Fatalf("Difference ton = 430.0, got %f", differenceTon)
	}

	diffPct := (differenceTon / expectedWeight) * 100.0
	if math.Abs(diffPct-0.7926) > 0.01 {
		t.Fatalf("Difference pct should be ~0.79%%, got %f%%", diffPct)
	}
}

func TestWebhookSurveyPayloadUnmarshal(t *testing.T) {
	rawJSON := `{
		"quarryCode": "MO-PT-01",
		"cycleCode": "CYCLE-2026-Q3-01",
		"surveyDate": "2026-08-26",
		"surveyMethod": "drone",
		"operatorName": "Nguyen Van A",
		"areas": [
			{
				"areaCode": "KV-01",
				"areaName": "Moong Tang 1",
				"areaM2": 45000,
				"extractedVolumeM3": 43500,
				"materialCode": "DA-1X2",
				"models": [
					{
						"modelType": "mesh",
						"fileFormat": "glb",
						"storagePath": "r2://models/cycle-01/mesh.glb",
						"fileSize": 52400000
					}
				]
			}
		]
	}`

	var payload models.WebhookSurveyPayload
	err := json.Unmarshal([]byte(rawJSON), &payload)
	if err != nil {
		t.Fatalf("Failed to unmarshal webhook payload: %v", err)
	}

	if payload.QuarryCode != "MO-PT-01" {
		t.Errorf("Expected quarryCode MO-PT-01, got %s", payload.QuarryCode)
	}
	if len(payload.Areas) != 1 {
		t.Fatalf("Expected 1 area, got %d", len(payload.Areas))
	}
	if payload.Areas[0].ExtractedVolumeM3 != 43500 {
		t.Errorf("Expected extracted volume 43500, got %f", payload.Areas[0].ExtractedVolumeM3)
	}
	if len(payload.Areas[0].Models) != 1 {
		t.Fatalf("Expected 1 model, got %d", len(payload.Areas[0].Models))
	}
	if payload.Areas[0].Models[0].FileFormat != "glb" {
		t.Errorf("Expected file format glb, got %s", payload.Areas[0].Models[0].FileFormat)
	}
}
