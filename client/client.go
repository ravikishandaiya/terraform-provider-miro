package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"io/ioutil"
)

type createUserRequest struct {
	Emails	[]string	`json:"emails"`
}

type user struct {
	Type	string	`json:"type"`
	Name	string	`json:"name"`
	ID		string	`json:"id"`
} 

type createUserResponce struct {
	Type	string	`json:"type"`
	ID		string	`json:"id"`
	User	user	`josn:"user"`
}

type Client struct {
	authToken  string
	team_id	   string
	httpClient *http.Client 
}

type team struct {
	Type	string	`json:"type"`
	Name	string	`json:"name"`
	ID		string	`json:"id"`
}
type data struct {
	Type	string	`json:"type"`
	ID		string	`json:"id"`
	User	user	`josn:"user"`
	Team 	team	`json:"team"`
	Role	string	`json:"role"`
}

type listAllUserResponce struct {
	Type	string	`json:"type"`
	Limit	int		`json:"limit"`
	offset	int		`json:"offset"`
	size	int		`json:"size"`
	prevLink	string	`json:"prevLink"`
	nextLink	string	`json:"nextLink"`
	Data	[]data	`josn:"data"`
}

type getUserStruct struct {
	Type	string	`json:"type"`
	ID		string	`json:"id"`
	Name	string	`json:"name"`
	CreatedAt	string	`json:"createdAt"`
	Role	string	`json:"role"`
	TeamName	string
	Email	string	`json:"email"`
	State	string	`json:"state"`
}

type update struct {
	Role	string	`json:"role"`
}

var (
    Errors = make(map[int]string)
)

func init() {
	Errors[400] = "Bad Request, StatusCode = 400"
	Errors[404] = "User Does Not Exist , StatusCode = 404"
	Errors[409] = "User Already Exist, StatusCode = 409"
	Errors[401] = "Unautharized Access, StatusCode = 401"
	Errors[429] = "User Has Sent Too Many Request, StatusCode = 429"
	Errors[500] = "Internal Server Error"
	
	Errors[501] = "Not Implemented"
	Errors[502] = "Bad Gateway"
	Errors[503] = "Service Unavailable"
	Errors[504] = "Gateway Timeout"
	Errors[505] = "HTTP Version Not Supported"
	Errors[506] = "Variant Also Negotiates"
	Errors[507] = "Insufficient Storage"
	Errors[508] = "Loop Detected"
	Errors[510] = "Not Extended"
	Errors[511] = "Network Authentication Required"
}

func NewClient(token string,team_id string) *Client {
	return &Client{
		authToken:  token,
		team_id:	team_id,
		httpClient: &http.Client{},
	}
}

func (c *Client) handleRequest(httpMethod string,url string, body []byte) (responce *http.Response, err error) {
	httpClient  := &http.Client{}
	var req 	*http.Request
	req, err 	= http.NewRequest(httpMethod, url, bytes.NewBuffer(body))
	if err != nil {
		return 
	}
	req.Header.Add("Authorization", "Bearer "+c.authToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	responce, err = httpClient.Do(req)
	if err != nil {
		return 
	}
	return
}

func (c *Client) CreateUser(email string) (error) {
	url 	  := fmt.Sprintf("https://api.miro.com/v1/teams/%s/invite", c.team_id)
	payload   := createUserRequest {
		Emails: []string {email},
	}
	body, err 	  := json.Marshal(payload)
	responce ,err := c.handleRequest(http.MethodPost, url, body)
	if err != nil {
		return err
	}
	resp, err := ioutil.ReadAll(responce.Body)
	var createResponceStruct []createUserResponce
	err = json.Unmarshal(resp, &createResponceStruct)
	if responce.StatusCode >= 200 && responce.StatusCode <= 299 {
		if len(createResponceStruct) == 0 {
			return fmt.Errorf("Error : User Already Exist.")
		}
		return nil
    } else {
		return fmt.Errorf("Error : %v",Errors[responce.StatusCode] )
    }
}

func (c *Client) getAllTeamMembers() ([] data, []string, error) {
	var list []string
	var ResponceStruct listAllUserResponce
	url := fmt.Sprintf("https://api.miro.com/v1/teams/%s/user-connections?limit=10&offset=0", c.team_id)
	resp,err := c.handleRequest(http.MethodGet, url, nil)
	if err != nil {
		return ResponceStruct.Data,list,err
	}
	responce, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ResponceStruct.Data,list,err
	}
	err = json.Unmarshal(responce, &ResponceStruct)
	if err != nil {
		return ResponceStruct.Data,list,err
	}
	for a := range(ResponceStruct.Data) {
		list = append(list, (ResponceStruct.Data[a].User.ID))
	}
	return ResponceStruct.Data,list,err
}

func (c* Client) checkUserExist(Data []data,userIDs []string, email string) (getUserStruct, error) {
	base_url := "https://api.miro.com/v1/users/"
	var responceStruct getUserStruct
	for a := range(userIDs) {
		resp, err := c.handleRequest(http.MethodGet, base_url+userIDs[a], nil)
		if err != nil {
			return responceStruct,err
		}
		responce, err := ioutil.ReadAll(resp.Body)
		err = json.Unmarshal(responce, &responceStruct)
		if err != nil {
			return responceStruct,err
		}
		if responceStruct.Email == email {
			for b := range(Data) {
				if responceStruct.ID == Data[b].User.ID {
					responceStruct.ID = Data[b].ID
					responceStruct.Role = Data[b].Role
					responceStruct.TeamName = Data[b].Team.Name
					break
				}else if b == len(Data)-1 {
					return responceStruct, fmt.Errorf("Undefined Error Encountered. IDs doesn't match")
				}
			}
			return responceStruct, nil
		}
	}
	return responceStruct, fmt.Errorf("User Not Found")
}

func (c *Client) GetUser(email string) (getUserStruct, error) {
	Data,userIds,err := c.getAllTeamMembers()
	if err != nil {
		var returnStruct getUserStruct
		return returnStruct,err
	}
	return c.checkUserExist(Data,userIds, email)
}

func (c *Client) Get_User_ID(email string) (user_id string, err error) {
	Data,userIds,err := c.getAllTeamMembers()
	if err != nil {
		return
	}
	resp, err := c.checkUserExist(Data,userIds, email)
	if err!= nil {
		return	
	}
	user_id = resp.ID
	return
}

func (c *Client) UpdateUser(email string, role string) (error) {
	payload := update{
		Role: role,
	}
	user_id, err := c.Get_User_ID(email)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.miro.com/v1/team-user-connections/%s",user_id)
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp ,err := c.handleRequest(http.MethodPatch, url, body)
	if resp.StatusCode != 200 {
		return  fmt.Errorf("%s",Errors[resp.StatusCode])
	}
	return err
}

func (c *Client) DeleteUser(email string) (error) {
	user_id, err := c.Get_User_ID(email)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("https://api.miro.com/v1/team-user-connections/%s",user_id)
	resp ,err := c.handleRequest(http.MethodDelete, url, nil)
	if resp.StatusCode != 204 {
		return  fmt.Errorf("%s",Errors[resp.StatusCode])
	}
	return err
}
