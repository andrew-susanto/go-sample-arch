package jsonrpc

import (
	// golang package
	"context"

	// extermal package
	"github.com/ybbus/jsonrpc/v3"
)

//go:generate mockgen -source=./repo.go -destination=./repo_mock.go -package=jsonrpc

type RPCClient interface {
	// Call is used to send a JSON-RPC request to the server endpoint.
	//
	// The spec states, that params can only be an array or an object, no primitive values.
	// So there are a few simple rules to notice:
	//
	// 1. no params: params field is omitted. e.g. Call(ctx, "getinfo")
	//
	// 2. single params primitive value: value is wrapped in array. e.g. Call(ctx, "getByID", 1423)
	//
	// 3. single params value array or object: value is unchanged. e.g. Call(ctx, "storePerson", &Person{Name: "Alex"})
	//
	// 4. multiple params values: always wrapped in array. e.g. Call(ctx, "setDetails", "Alex, 35, "Germany", true)
	//
	// Examples:
	//   Call(ctx, "getinfo") -> {"method": "getinfo"}
	//   Call(ctx, "getPersonId", 123) -> {"method": "getPersonId", "params": [123]}
	//   Call(ctx, "setName", "Alex") -> {"method": "setName", "params": ["Alex"]}
	//   Call(ctx, "setMale", true) -> {"method": "setMale", "params": [true]}
	//   Call(ctx, "setNumbers", []int{1, 2, 3}) -> {"method": "setNumbers", "params": [1, 2, 3]}
	//   Call(ctx, "setNumbers", 1, 2, 3) -> {"method": "setNumbers", "params": [1, 2, 3]}
	//   Call(ctx, "savePerson", &Person{Name: "Alex", Age: 35}) -> {"method": "savePerson", "params": {"name": "Alex", "age": 35}}
	//   Call(ctx, "setPersonDetails", "Alex", 35, "Germany") -> {"method": "setPersonDetails", "params": ["Alex", 35, "Germany"}}
	//
	// for more information, see the examples or the unit tests
	Call(ctx context.Context, method string, params ...interface{}) (*jsonrpc.RPCResponse, error)
}

type Repository struct {
	cxpigw RPCClient
}

// NewRepository is a function to initialize jsonrpc repository
func NewRepository(cxpigw RPCClient) Repository {
	return Repository{
		cxpigw: cxpigw,
	}
}
