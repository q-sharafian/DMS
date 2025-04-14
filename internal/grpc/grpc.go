package grpcserver

import (
	l "DMS/internal/logger"
	m "DMS/internal/models"
	service "DMS/internal/services"
	"context"
	"runtime/debug"
	"time"

	pbAuth "github.com/q-sharafian/file-transfer/pkg/pb/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	logger    l.Logger
	fpService service.FilePermissionService
	pbAuth.UnimplementedAuthServer
}

func NewGRPCServer(fpService service.FilePermissionService, logger l.Logger) GRPCServer {
	return GRPCServer{logger: logger, fpService: fpService}
}

func (s *GRPCServer) IsAllowedDownload(c context.Context, dar *pbAuth.DownloadAccessReq) (*pbAuth.AllowDownloadResult, error) {
	objectTokens := make([]m.Token, 0)
	for _, token := range dar.ObjectTokens {
		objectTokens = append(objectTokens, m.Str2Token(token))
	}
	accessInfo := m.DownloadReq{
		AuthToken:    m.Str2Token(dar.AuthToken),
		ObjectTokens: objectTokens,
	}

	result, err := s.fpService.IsAllowedDownload(&accessInfo)
	if err != nil {
		s.logger.Debugf("Error in checking download permission: %s", err.Error())
		switch err.GetCode() {
		case service.SEInternal:
			return &pbAuth.AllowDownloadResult{StatusCode: pbAuth.StatusCode_ErrInternal,
				Errmsg: err.Error()}, nil
		case service.SEForbidden:
			return &pbAuth.AllowDownloadResult{StatusCode: pbAuth.StatusCode_ErrForbidden,
				Errmsg: err.Error()}, nil
		case service.SEAuthFailed:
			return &pbAuth.AllowDownloadResult{StatusCode: pbAuth.StatusCode_ErrUnauthorized,
				Errmsg: err.Error()}, nil
		default:
			s.logger.Infof("Unexpected error in checking download permission (err code: %d): %s", err.GetCode(), err.Error())
			return &pbAuth.AllowDownloadResult{StatusCode: pbAuth.StatusCode_ErrInternal,
				Errmsg: err.Error()}, nil
		}
	}
	downloadResult := pbAuth.AllowDownloadResult{StatusCode: pbAuth.StatusCode_OK, Files: make(map[string]bool)}
	for token, isAllowed := range result {
		downloadResult.Files[token.String()] = isAllowed
	}
	return &downloadResult, nil
}

func (s *GRPCServer) IsAllowedUpload(c context.Context, uar *pbAuth.UploadAccessReq) (*pbAuth.AllowUploadResult, error) {
	accessInfo := m.UploadReq{
		AuthToken:   m.Str2Token(uar.AuthToken),
		ObjectTypes: make(map[m.FileExtension]uint, 0),
	}
	for k, count := range uar.ObjectTypes {
		accessInfo.ObjectTypes[m.FileExtension(k)] = uint(count)
	}

	result, err := s.fpService.IsAllowedUpload(&accessInfo)
	if err != nil {
		s.logger.Debugf("Error in checking upload permission: %s", err.Error())
		switch err.GetCode() {
		case service.SEInternal:
			return &pbAuth.AllowUploadResult{StatusCode: pbAuth.StatusCode_ErrInternal,
				Errmsg: err.Error()}, nil
		case service.SEForbidden:
			return &pbAuth.AllowUploadResult{StatusCode: pbAuth.StatusCode_ErrForbidden,
				Errmsg: err.Error()}, nil
		case service.SEAuthFailed:
			return &pbAuth.AllowUploadResult{StatusCode: pbAuth.StatusCode_ErrUnauthorized,
				Errmsg: err.Error()}, nil
		default:
			s.logger.Infof("Unexpected error in checking download permission (err code: %d): %s", err.GetCode(), err.Error())
			return &pbAuth.AllowUploadResult{StatusCode: pbAuth.StatusCode_ErrInternal,
				Errmsg: err.Error()}, nil
		}
	}
	uploadResult := pbAuth.AllowUploadResult{StatusCode: pbAuth.StatusCode_OK, FileTypes: make([]*pbAuth.AcceptableType, 0)}
	for _, t := range result {
		uploadResult.FileTypes = append(uploadResult.FileTypes, &pbAuth.AcceptableType{
			FileType: string(t.FileType),
			MaxSize:  t.MaxSize,
			IsAllow:  t.IsAllow,
		})
	}
	return &uploadResult, nil
}

func (s *GRPCServer) LoggerInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	s.logger.Debug("--------------------------------------------------------------------------------")
	startTime := time.Now()
	resp, err = handler(ctx, req)
	duration := time.Since(startTime)
	s.logger.Debugf("Received request via gRPC. Method: %s, Process Duration: %s", info.FullMethod, duration)
	return resp, err
}
func (s *GRPCServer) ErrorInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			s.logger.Errorf("Panic recovered: %+v\n%s", r, debug.Stack())
			err = status.Errorf(codes.Internal, "internal server panic %+v", r)
		}
	}()
	resp, err = handler(ctx, req)
	if err != nil {
		// Log and convert the error to a gRPC error
		s.logger.Debugf("Error: %v", err)
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	return resp, nil
}
