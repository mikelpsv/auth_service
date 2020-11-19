package models

import (
	"github.com/mikelpsv/auth_service/app"
	"testing"
)

func TestValidPassword(t *testing.T) {
	res, err := app.ValidPassword("$2a$10$/ui7v1gRNVLSRtfHOib/muwP5TAr7e33c9y7LPpfdUHmCIWJSO8ny", "mypassword")
	if !res {
		t.Error(!res)
	}
	if err != nil {
		t.Error(err != nil)
	}

	res, err = app.ValidPassword("$2a$10$/ui7v1gRNVLSRtfHOib/muwP5TAr7e33c9y7LPpfdUHmCIWJSO8ny", "mypassword_incorrect")
	if res {
		t.Error(res)
	}
	if err == nil {
		t.Error(err == nil)
	}

	res, err = app.ValidPassword("$2a$10$/ui7v1gRNVLSRtfHOib/muwP5TAr7e33c9y7LPpfdUHmCIWJSO8ny", "")
	if res {
		t.Error(res)
	}
	if err == nil {
		t.Error(err == nil)
	}
}
