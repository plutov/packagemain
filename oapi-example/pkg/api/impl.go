package api

import "context"

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (Server) ListDevices(ctx context.Context, request ListDevicesRequestObject) (ListDevicesResponseObject, error) {
	// actual implementation
	return ListDevices200JSONResponse{}, nil
}

func (Server) GetDevice(ctx context.Context, request GetDeviceRequestObject) (GetDeviceResponseObject, error) {
	// actual implementation
	return GetDevice200JSONResponse{}, nil
}
