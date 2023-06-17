package rpcclient

import (
	// golang package
	"context"
	"fmt"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/config"

	// external package
	"github.com/afex/hystrix-go/hystrix"
	"github.com/ybbus/jsonrpc/v3"
)

//go:generate mockgen -source=rpcclient.go -destination=rpcclient_mock.go -package=rpcclient

type jsonRpcClient interface {
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

// RpcClient is abstraction of rpc client with circuit breaker
type RpcClient struct {
	cbConfig map[string]bool
	client   jsonRpcClient
	config   config.RpcClientConfig
}

// NewClient create rpc client for future rpc request
func NewClient(rpcClientConfig config.RpcClientConfig) *RpcClient {
	client := &RpcClient{
		cbConfig: map[string]bool{},
		client:   jsonrpc.NewClient(rpcClientConfig.ServiceUrl),
		config:   rpcClientConfig,
	}

	return client
}

// Call will call rpc method with given method name and params
//
// returns rpc response and nil error if success
// otherwise return empty rpc response and non nil error
func (rpcClient *RpcClient) Call(ctx context.Context, method string, params ...interface{}) (*jsonrpc.RPCResponse, error) {
	cbKey := fmt.Sprintf("%s_%s", rpcClient.config.ServiceName, method)
	if _, exists := rpcClient.cbConfig[method]; !exists {
		hystrix.ConfigureCommand(cbKey, hystrix.CommandConfig{
			Timeout:               rpcClient.config.Timeout,
			MaxConcurrentRequests: rpcClient.config.MaxConcurrentRequest,
			ErrorPercentThreshold: rpcClient.config.ErrorPercentageThreshold,
		})

		rpcClient.cbConfig[method] = true
	}

	var result *jsonrpc.RPCResponse
	var err error

	err = hystrix.Do(cbKey, func() error {
		result, err = rpcClient.client.Call(ctx, method, params...)
		return err
	}, nil)

	return result, err
}
