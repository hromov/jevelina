package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type Test struct {
	name           string
	server         httptest.Server
	expectedError  string
	expectedStatus int
}

func TestUserCheck(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(UserCheck(nil))

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusForbidden)
	}

	// Check the response body is what we expect.
	expected := "User access required\n"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %q want %q",
			rr.Body.String(), expected)
	}

}

// func TestUserCheckHandler(t *testing.T) {
// 	tests := []Test{
// 		{
// 			name: "Not a user test",
// 			server: *httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 				w.WriteHeader(http.StatusForbidden)
// 			})),
// 			expectedStatus: http.StatusForbidden,
// 			expectedError:  "Authorization error: Authorization required\n",
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			defer test.server.Close()
// 			request, _ := http.NewRequest(http.MethodGet, "/authcheck", nil)
// 			response := httptest.NewRecorder()

// 			UserCheckHandler(response, request)

// 			if response.Result().StatusCode != test.expectedStatus {
// 				t.Errorf("Status code error. Got %q, want %q", response.Result().StatusCode, test.expectedStatus)
// 			}

// 			if response.Body.String() != test.expectedError {
// 				t.Errorf("got %q, want %q", response.Body.String(), test.expectedError)
// 			}
// 		})
// 	}
// }
