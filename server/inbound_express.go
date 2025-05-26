package server

import (
	"context"
	"errors"
	inbound "hpc-express-service/inbound/express"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type inboundExpressHandler struct {
	s inbound.InboundExpressService
}

func (h *inboundExpressHandler) router() chi.Router {

	r := chi.NewRouter()

	r.Post("/upload", h.uploadManifest)
	r.Get("/download/pre-import", h.downloadPreImport)
	r.Get("/download/raw-pre-import", h.downloadRawPreImport)
	r.Post("/upload-update-manifest", h.uploadUpdateManifest)

	return r
}

func (h *inboundExpressHandler) uploadManifest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	r.ParseMultipartForm(10 << uint32(20)) // 10 * 2^20
	file, handler, err := r.FormFile("file")
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	templateCode := r.FormValue("templateCode")
	userUUID := GetUserUUIDFromContext(r)

	log.Println("#1 ", templateCode)

	// result, err := h.s.UploadManifest(ctx, userUUID, handler.Filename, fileBytes)
	err = h.s.UploadManifest(ctx, userUUID, handler.Filename, templateCode, fileBytes)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, SuccessResponse(nil, "success"))
}

func (h *inboundExpressHandler) downloadPreImport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	uploadLoggingUUID := r.URL.Query().Get("uploadLoggingUUID")
	if len(uploadLoggingUUID) == 0 {
		render.Render(w, r, ErrInvalidRequest(errors.New("required uuid")))
		return
	}

	fileName, zipBuf, err := h.s.DownloadPreImport(ctx, uploadLoggingUUID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// Send ZIP file as response
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", `attachment; filename="`+fileName+`.zip"`)
	w.WriteHeader(http.StatusOK)
	w.Write(zipBuf.Bytes())
}

func (h *inboundExpressHandler) downloadRawPreImport(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	uploadLoggingUUID := r.URL.Query().Get("uploadLoggingUUID")
	if len(uploadLoggingUUID) == 0 {
		render.Render(w, r, ErrInvalidRequest(errors.New("required uuid")))
		return
	}

	fileName, excelBuf, err := h.s.DownloadRawPreImport(ctx, uploadLoggingUUID)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// Send ZIP file as response
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("File-Name", fileName)
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusOK)
	w.Write(excelBuf.Bytes())
}

func (h *inboundExpressHandler) uploadUpdateManifest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	r.ParseMultipartForm(10 << uint32(20)) // 10 * 2^20
	file, handler, err := r.FormFile("file")
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	userUUID := GetUserUUIDFromContext(r)

	err = h.s.UploadUpdateRawPreImport(ctx, userUUID, handler.Filename, fileBytes)
	if err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	render.Respond(w, r, SuccessResponse(nil, "success"))
}
