package jenkins

type Client struct {
	URL   string
	User  string
	Token string
}

func NewClient(url, user, token string) *Client {
	return &Client{
		URL:   url,
		User:  user,
		Token: token,
	}
}
