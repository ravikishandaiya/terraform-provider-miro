package client

import(
	"github.com/stretchr/testify/assert"
	"testing"
	"os"
)

func init(){
	token := "{api token}"
	team_id := "{team ID}"
	os.Setenv("MIRO_TOKEN", token)
	os.Setenv("TEAM_ID",team_id)
}

func TestClient_GetUser(t *testing.T) {
	testCases := []struct {
		testName     string
		userName     string
		seedData     map[string]getUserStruct
		expectErr    bool
		expectedResp getUserStruct
	}{
		{
			testName: "user exists",
			userName: "{email}",
			seedData: map[string]getUserStruct{
				"user1": {
					Type:		"user",
					ID:			"{USER ID}",
					Name:		"{USER NAME}",
					CreatedAt:	"{DATE}",
					Role:		"{ROLE}",
					TeamName:	"{TEAM NAME}",
					Email:		"{EMAIL}",
					State:		"registered",
				},
			},
			expectErr: false,
			expectedResp: getUserStruct{
				Type:		"user",
				ID:			"{USER ID}",
				Name:		"{USER NAME}",
				CreatedAt:	"{DATE}",
				Role:		"{ROLE}",
				TeamName:	"{TEAM NAME}",
				Email:		"{EMAIL}",
				State:		"registered",
			},
		},
			{
			testName:     "user does not exist",
			userName:     "{email of user wwho  doesn't exist in team}",
			seedData:     nil,
			expectErr:    true,
			expectedResp: getUserStruct{},
			},
		}
	

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client 		:= NewClient(os.Getenv("MIRO_TOKEN"),os.Getenv("TEAM_ID"))
			item, err 	:= client.GetUser(tc.userName)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResp, item)
		})
	}
}
func TestClient_CreateUser(t *testing.T) {
	testCases :=  []struct {
		testName  string
		newUser   string
		seedData  map[string]string
		expectErr bool
	}{
		{
			testName: "success",
			newUser:  "{mail of new user}",
			seedData:  nil,
			expectErr: false,
		},
		{
			testName: "user already exists",
			newUser:  "{user's email, who exist in the team}",
			seedData:  map[string]string {"user1": "{email}"},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client 	:= NewClient( os.Getenv("MIRO_TOKEN"), os.Getenv("TEAM_ID"))
			err 	:= client.CreateUser(tc.newUser)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestClient_UpdateUser(t *testing.T) {
	testCases := []struct {
		testName    string
		role		string
		seedData    map[string]string
		expectErr   bool
		email		string
	}{
		{
			testName: "user exists",
			role:	  "{role}",
			email:	  "{email}",
			seedData: map[string]string{
				"user1": "{user ID}",
			},
			expectErr: false,
		},
		{
			testName: "user does not exist",
			role:	  "{role}",
			email:	  "{email}",
			seedData:  nil,
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client := NewClient(os.Getenv("MIRO_TOKEN"),os.Getenv("TEAM_ID"))
			err    := client.UpdateUser(tc.email,tc.role)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			user, err := client.GetUser(tc.email)
			assert.NoError(t, err)
			assert.Equal(t, tc.role, user.Role)
		})
	}
}

func TestClient_DeleteUser(t *testing.T) {
	testCases := []struct {
		testName  string
		user	  string
		seedData  map[string]string
		expectErr bool
	}{
		{
			testName: "user exists",
			user:     "{email}",
			seedData: map[string]string{
				"user1": "{user id}",
			},
			expectErr: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client 		:= NewClient(os.Getenv("MIRO_TOKEN"),os.Getenv("TEAM_ID"))
			err := client.DeleteUser(tc.user)
			if err != nil {
				assert.NoError(t,err)
			}
		})
	}
}
