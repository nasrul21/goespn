package espn

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

//ICoreGateway .
type ICoreGateway interface {
	GetStandings() (interface{}, error)
}

//CoreGateway .
type CoreGateway struct {
	Client   Client
	SportID  string
	LeagueID string
}

// NewCoreGateway .
func NewCoreGateway(client Client, sportID, leagueID string) ICoreGateway {
	return &CoreGateway{
		Client:   client,
		SportID:  sportID,
		LeagueID: leagueID,
	}
}

// Call : base method to call Core API
func (gateway *CoreGateway) Call(method, path string, body io.Reader, v interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return gateway.Client.Call(method, path, body, v)
}

// GetStandings : get league standings
func (gateway *CoreGateway) GetStandings() (interface{}, error) {
	path := fmt.Sprintf("/v2/sports/%s/%s/standings", gateway.SportID, gateway.LeagueID)

	var resp interface{}
	err := gateway.Call(http.MethodGet, path, nil, &resp)
	if err != nil {
		fmt.Println("Error GetStandings: ", err)
		return resp, err
	}

	// fmt.Println("RESPONSE: ", resp)

	return resp, nil
}
