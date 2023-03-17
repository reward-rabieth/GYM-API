package api

import "testing"

func TestHandleCreateMember(t *testing.T) {

	s := APIServer{ListenAddress: "9000"}
	go s.Run()

}
