package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTime(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {

		rw.Write([]byte(`{"week_number":41,"utc_offset":"+03:00","utc_datetime":"2019-10-11T12:48:58.263816+00:00","unixtime":1570798138,"timezone":"Europe/Moscow","raw_offset":10800,"dst_until":null,"dst_offset":0,"dst_from":null,"dst":false,"day_of_year":284,"day_of_week":5,"datetime":"2019-10-11T15:48:58.263816+03:00","client_ip":"83.139.159.134","abbreviation":"MSK"}`))
	}))
	defer server.Close()

	want := int64(1570798138)
	tc := newTimeClient(server.URL)
	tt, err := tc.getTime()

	if err != nil {
		t.Errorf("unexpected error %v", err)
		return
	}

	if got := tt.Unix(); got != want {
		t.Errorf("want %d but got %d", want, got)
	}
}
