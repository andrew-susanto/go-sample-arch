package cxpigw

import (
	// golang package
	"context"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/repository/jsonrpc"
)

//go:generate mockgen -source=./resource.go -destination=./resource_mock.go -package=user

type JSONRpcRepository interface {
	// GetTripItinerary is a function to get trip itinerary by booking id from jsonrpc
	//
	// Returns TripItinerary and nil error if sucess
	// Otherwise return empty TripItinerary and non nil error
	GetTripItinerary(ctx context.Context, bookingID int64) (jsonrpc.TripItinerary, error)
}

type resource struct {
	jsonrpc JSONRpcRepository
}

// NewResource creates new resource for cxpigw jsonrpc service
func NewResource(jsonrpc JSONRpcRepository) resource {
	return resource{
		jsonrpc: jsonrpc,
	}
}
