package auth

// func TestUserCheck(t *testing.T) {
// 	req, err := http.NewRequest("GET", "/", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	rr := httptest.NewRecorder()
// 	handler := UserCheck(nil)

// 	handler.ServeHTTP(rr, req)

// 	if status := rr.Code; status != http.StatusForbidden {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusForbidden)
// 	}

// 	expected := "User access required\n"
// 	if rr.Body.String() != expected {
// 		t.Errorf("handler returned unexpected body: got %q want %q",
// 			rr.Body.String(), expected)
// 	}

// }
