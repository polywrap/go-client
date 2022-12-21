package wasm

type (
	ClientConfig struct {
	}

	Client struct {
	}
)

func NewClient(cfg *ClientConfig) *Client {
	return nil
}

func (client *Client) Invoke(method string, data []byte) ([]byte, error) {
	return nil, nil
}

func Invoke[InvokeArg, InvokeRes any](client *Client, method string, arguments InvokeArg) (*InvokeRes, error) {
	args, err := Encode(arguments)
	if err != nil {
		return nil, err
	}

	resp, err := client.Invoke(method, args)
	if err != nil {
		return nil, err
	}

	res, err := Decode[InvokeRes](resp)
	if err != nil {
		return nil, err
	}

	return &res, nil
}
