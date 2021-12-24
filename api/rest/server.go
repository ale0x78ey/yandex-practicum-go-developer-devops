package rest

import (
	"net/http"
)

type Server interface {
	NewRouter() *http.ServeMux
}

type server struct {
}

func NewServer() *server {
	return &server{}
}

func (s *server) UpdateAllocHandler(value string) {
}

func (s *server) UpdateBuckHashSysHandler(value string) {
}

func (s *server) UpdateFreesHandler(value string) {
}

func (s *server) UpdateGCCPUFractionHandler(value string) {
}

func (s *server) UpdateGCSysHandler(value string) {
}

func (s *server) UpdateHeapAllocHandler(value string) {
}

func (s *server) UpdateHeapIdleHandler(value string) {
}

func (s *server) UpdateHeapInuseHandler(value string) {
}

func (s *server) UpdateHeapObjectsHandler(value string) {
}

func (s *server) UpdateHeapReleasedHandler(value string) {
}

func (s *server) UpdateHeapSysHandler(value string) {
}

func (s *server) UpdateLastGCHandler(value string) {
}

func (s *server) UpdateLookupsHandler(value string) {
}

func (s *server) UpdateMCacheInuseHandler(value string) {
}

func (s *server) UpdateMCacheSysHandler(value string) {
}

func (s *server) UpdateMSpanInuseHandler(value string) {
}

func (s *server) UpdateMSpanSysHandler(value string) {
}

func (s *server) UpdateMallocsHandler(value string) {
}

func (s *server) UpdateNextGCHandler(value string) {
}

func (s *server) UpdateNumForcedGCHandler(value string) {
}

func (s *server) UpdateNumGCHandler(value string) {
}

func (s *server) UpdateOtherSysHandler(value string) {
}

func (s *server) UpdatePauseTotalNsHandler(value string) {
}

func (s *server) UpdateStackInuseHandler(value string) {
}

func (s *server) UpdateStackSysHandler(value string) {
}

func (s *server) UpdateSysHandler(value string) {
}

func (s *server) UpdateRandomValueHandler(value string) {
}

func (s *server) UpdatePollCountHandler(value string) {
}
