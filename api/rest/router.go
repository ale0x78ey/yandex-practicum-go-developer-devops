package rest

import (
	io "io/ioutil"
	"net/http"
	"path"
)

type updateHandler func(value string)

func (f updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		code := http.StatusMethodNotAllowed
		http.Error(w, "Only POST requests are allowed!", code)
		return
	}

	_, err := io.ReadAll(r.Body)
	if err != nil {
		code := http.StatusInternalServerError
		http.Error(w, err.Error(), code)
		return
	}

	w.WriteHeader(http.StatusOK)
	f(path.Base(r.URL.Path))
}

func (s *server) NewRouter() (mux *http.ServeMux) {
	mux = http.NewServeMux()

	mux.Handle("/update/gauge/Alloc/", updateHandler(s.UpdateAllocHandler))
	mux.Handle("/update/gauge/BuckHashSys/", updateHandler(s.UpdateBuckHashSysHandler))
	mux.Handle("/update/gauge/Frees/", updateHandler(s.UpdateFreesHandler))
	mux.Handle("/update/gauge/GCCPUFraction/", updateHandler(s.UpdateGCCPUFractionHandler))
	mux.Handle("/update/gauge/GCSys/", updateHandler(s.UpdateGCSysHandler))
	mux.Handle("/update/gauge/HeapAlloc/", updateHandler(s.UpdateHeapAllocHandler))
	mux.Handle("/update/gauge/HeapIdle/", updateHandler(s.UpdateHeapIdleHandler))
	mux.Handle("/update/gauge/HeapInuse/", updateHandler(s.UpdateHeapInuseHandler))
	mux.Handle("/update/gauge/HeapObjects/", updateHandler(s.UpdateHeapObjectsHandler))
	mux.Handle("/update/gauge/HeapReleased/", updateHandler(s.UpdateHeapReleasedHandler))
	mux.Handle("/update/gauge/HeapSys/", updateHandler(s.UpdateHeapSysHandler))
	mux.Handle("/update/gauge/LastGC/", updateHandler(s.UpdateLastGCHandler))
	mux.Handle("/update/gauge/Lookups/", updateHandler(s.UpdateLookupsHandler))
	mux.Handle("/update/gauge/MCacheInuse/", updateHandler(s.UpdateMCacheInuseHandler))
	mux.Handle("/update/gauge/MCacheSys/", updateHandler(s.UpdateMCacheSysHandler))
	mux.Handle("/update/gauge/MSpanInuse/", updateHandler(s.UpdateMSpanInuseHandler))
	mux.Handle("/update/gauge/MSpanSys/", updateHandler(s.UpdateMSpanSysHandler))
	mux.Handle("/update/gauge/Mallocs/", updateHandler(s.UpdateMallocsHandler))
	mux.Handle("/update/gauge/NextGC/", updateHandler(s.UpdateNextGCHandler))
	mux.Handle("/update/gauge/NumForcedGC/", updateHandler(s.UpdateNumForcedGCHandler))
	mux.Handle("/update/gauge/NumGC/", updateHandler(s.UpdateNumGCHandler))
	mux.Handle("/update/gauge/OtherSys/", updateHandler(s.UpdateOtherSysHandler))
	mux.Handle("/update/gauge/PauseTotalNs/", updateHandler(s.UpdatePauseTotalNsHandler))
	mux.Handle("/update/gauge/StackInuse/", updateHandler(s.UpdateStackInuseHandler))
	mux.Handle("/update/gauge/StackSys/", updateHandler(s.UpdateStackSysHandler))
	mux.Handle("/update/gauge/Sys/", updateHandler(s.UpdateSysHandler))
	mux.Handle("/update/gauge/RandomValue/", updateHandler(s.UpdateRandomValueHandler))
	mux.Handle("/update/counter/PollCount/", updateHandler(s.UpdatePollCountHandler))
	return
}
