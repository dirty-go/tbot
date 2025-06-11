package tbot

type ClientOptions func(*Client)

func WithBaseURL(baseURL string) ClientOptions {
	return func(client *Client) {
		client.baseURL = baseURL
	}
}
