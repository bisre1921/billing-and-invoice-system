package main

import (
	"encoding/json"
	"fmt"
	"github.com/bisre1921/billing-and-invoice-system/models"
)

func main() {
	// Test invalid ObjectID JSON unmarshaling
	invalidJSON := `{"company_id": "invalid-id", "email": "test@example.com", "name": "Test", "position": "Dev"}`

	var req models.EmployeeInvitationRequest
	err := json.Unmarshal([]byte(invalidJSON), &req)

	if err != nil {
		fmt.Printf("JSON Unmarshal Error: %v\n", err)
	} else {
		fmt.Printf("JSON Unmarshal Success: %+v\n", req)
		fmt.Printf("CompanyID IsZero: %v\n", req.CompanyID.IsZero())
	}
}
