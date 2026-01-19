package service

import (
	"github.com/google/uuid"
	"github.com/mukundaparajuli/fixr/internal/model"
)

type ServiceAttachment struct {
	model.Base
	Name        string     `json:"name" db:"name"`
	ServiceID   *uuid.UUID `json:"serviceId" db:"service_id"`
	UploadedBy  string     `json:"uploadedBy" db:"uploaded_by"`
	DownloadKey string     `json:"donwloadKey" db:"downloadKey"`
	FileSize    *int64     `json:"fileSize" db:"file_size"`
	MimeType    *string    `json:"mimeType" db:"mime_type"`
}
