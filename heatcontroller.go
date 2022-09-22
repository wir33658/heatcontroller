// First Go program
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	JTI          string `json:"jti"`
}

type home struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type me struct { // Only a part, check Postman for the complete returned data object.
	Name     string `json:"name"`
	EMail    string `json:"email"`
	UserName string `json:"username"`
	Id       string `json:"id"`
	Homes    []home `json:"homes"`
}

type zone struct { // Only a part, check Postman for the complete returned data object.
	Id   int // no explicit json description needed as long as the name is the same (doesn't matter lower/upper case)
	Name string
	Type string
}

type temperature struct {
	Celsius    float32
	Fahrenheit float32
}

type zone_state_setting struct {
	Type        string
	Power       string
	Temperature temperature
}

type zone_state_overlay_termination struct {
	Type              string
	TypeSkellBasedApp string
	ProjectedExpiry   string
}

type zone_state_overlay struct {
	Type        string
	Setting     zone_state_setting
	Termination zone_state_overlay_termination
}

type zone_state_sdp_precision struct {
}
type zone_state_sdp_inside_temperature struct {
	Celsius    float32
	Fahrenheit float32
	Timestamp  string
	Type       string
	// Precision
}
type zone_state_sensor_data_points struct {
}

type zone_state struct {
	Setting     zone_state_setting
	OverlayType string
	Overlay     zone_state_overlay
}

/*

    "sensorDataPoints": {
        "insideTemperature": {
            "celsius": 18.28,
            "fahrenheit": 64.90,
            "timestamp": "2022-09-21T04:48:07.906Z",
            "type": "TEMPERATURE",
            "precision": {
                "celsius": 0.1,
                "fahrenheit": 0.1
            }
        },
        "humidity": {
            "type": "PERCENTAGE",
            "percentage": 69.60,
            "timestamp": "2022-09-21T04:48:07.906Z"
        }
    }
{
    "tadoMode": "HOME",
    "geolocationOverride": false,
    "geolocationOverrideDisableTime": null,
    "preparation": null,
    "setting": {
        "type": "HEATING",
        "power": "ON",
        "temperature": {
            "celsius": 20.00,
            "fahrenheit": 68.00
        }
    },
    "overlayType": "MANUAL",
    "overlay": {
        "type": "MANUAL",
        "setting": {
            "type": "HEATING",
            "power": "ON",
            "temperature": {
                "celsius": 20.00,
                "fahrenheit": 68.00
            }
        },
        "termination": {
            "type": "MANUAL",
            "typeSkillBasedApp": "MANUAL",
            "projectedExpiry": null
        }
    },
    "openWindow": null,
    "nextScheduleChange": {
        "start": "2022-09-21T06:30:00Z",
        "setting": {
            "type": "HEATING",
            "power": "ON",
            "temperature": {
                "celsius": 18.00,
                "fahrenheit": 64.40
            }
        }
    },
    "nextTimeBlock": {
        "start": "2022-09-21T06:30:00.000Z"
    },
    "link": {
        "state": "ONLINE"
    },
    "activityDataPoints": {
        "heatingPower": {
            "type": "PERCENTAGE",
            "percentage": 100.00,
            "timestamp": "2022-09-21T04:47:40.386Z"
        }
    },
    "sensorDataPoints": {
        "insideTemperature": {
            "celsius": 18.28,
            "fahrenheit": 64.90,
            "timestamp": "2022-09-21T04:48:07.906Z",
            "type": "TEMPERATURE",
            "precision": {
                "celsius": 0.1,
                "fahrenheit": 0.1
            }
        },
        "humidity": {
            "type": "PERCENTAGE",
            "percentage": 69.60,
            "timestamp": "2022-09-21T04:48:07.906Z"
        }
    }
}
*/

// Main function
func main() {

	fmt.Println("!... Hello World ...!")

	client := http.Client{}

	token_obj := getToken(client)
	fmt.Println("\n\nAccessToken:", token_obj.AccessToken)
	fmt.Println("\n\nRefreshToken:", token_obj.RefreshToken)
	fmt.Println("\n\nTokenType:", token_obj.TokenType)
	fmt.Println("\n\nExpiresIn:", token_obj.ExpiresIn)
	fmt.Println("\n\nScope:", token_obj.Scope)
	fmt.Println("\n\nJTI:", token_obj.JTI)

	me_obj := getMe(client, token_obj)
	fmt.Println("\n\nName:", me_obj.Name)
	fmt.Println("\n\nUserName:", me_obj.UserName)
	fmt.Println("\n\nId:", me_obj.Id)
	fmt.Println("\n\nEMail:", me_obj.EMail)
	fmt.Println("\n\nHomes:", me_obj.Homes)
	fmt.Println("\n\nHomes[0].Name:", me_obj.Homes[0].Name)
	fmt.Println("\n\nHomes[0].Id:", me_obj.Homes[0].Id)

	zones_obj := getZones(client, token_obj, me_obj.Homes[0].Id)
	fmt.Println("\n\nZones:", zones_obj)

	fmt.Println("Done")

}

func getZones(client http.Client, token_obj token, home_id int) []zone {
	home_id_str := strconv.Itoa(home_id)
	fmt.Println("\nhome_id:", home_id_str)

	url := "https://my.tado.com/api/v2/homes/" + home_id_str + "/zones"
	fmt.Println("\nURL:" + url)

	req, getErr2 := http.NewRequest("GET", url, nil)
	if getErr2 != nil {
		log.Fatal(getErr2)
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + token_obj.AccessToken},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	bodyStr := string(body)
	fmt.Println(bodyStr)

	zone_objs := []zone{}
	jsonErr := json.Unmarshal(body, &zone_objs)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return zone_objs
}

func getMe(client http.Client, token_obj token) me {

	req, getErr2 := http.NewRequest("GET", "https://my.tado.com/api/v2/me", nil)
	if getErr2 != nil {
		log.Fatal(getErr2)
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + token_obj.AccessToken},
	}

	//	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Body)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	bodyStr := string(body)
	fmt.Println(bodyStr)

	me_obj := me{}
	jsonErr := json.Unmarshal(body, &me_obj)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return me_obj
}

func getToken(client http.Client) token {

	req, getErr2 := http.NewRequest("POST", "https://auth.tado.com/oauth/token?client_id=tado-web-app&grant_type=password&scope=home.user&username=robert.weissmann@web.de&password=127SushiRob721&client_secret=wZaRN7rpjn3FoNyF5IFuxg9uMzYJcvOoQ8QWiIqS3hfk6gLhVlG57j5YNoZL2Rtc", nil)
	if getErr2 != nil {
		log.Fatal(getErr2)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.Body)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	fmt.Println(string(body))

	token_obj := token{}
	jsonErr := json.Unmarshal(body, &token_obj)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return token_obj
}
