// First Go program
package main

import (
	"time"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
//	"errors"
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
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type temperature struct {
	Celsius    float32 `json:"celsius"`
	Fahrenheit float32 `json:"fahrenheit"`
}

type zone_state_setting struct {
	Type        string      `json:"type"`
	Power       string      `json:"power"`
	Temperature temperature `json:"temperature"`
}

type zone_state_setting_wrapper struct {
	Setting zone_state_setting `json:"setting"`
}

type zone_state_overlay_termination struct {
	Type              string `json:"type"`
	TypeSkellBasedApp string `json:"typeSkellBasedApp"`
	ProjectedExpiry   string `json:"projectedExpiry"`
}

type zone_state_overlay struct {
	Type        string                         `json:"type"`
	Setting     zone_state_setting             `json:"setting"`
	Termination zone_state_overlay_termination `json:"termination"`
}

type zone_state_sdp_inside_temperature struct {
	Celsius    float32     `json:"celsius"`
	Fahrenheit float32     `json:"fahrenheit"`
	Timestamp  string      `json:"timestamp"`
	Type       string      `json:"type"`
	Precision  temperature `json:"precision"`
}

type zone_state_sdp_humidity struct {
	Type       string  `json:"type"`
	Percentage float32 `json:"percentage"`
	Timestamp  string  `json:"timestamp"`
}

type zone_state_sensor_data_points struct {
	InsideTemperature zone_state_sdp_inside_temperature `json:"insideTemperature"`
	Humidity          zone_state_sdp_humidity           `json:"humidity"`
}

type zone_state_adp_heating_power struct {
	Type       string  `json:"type"`
	Percentage float32 `json:"percentage"`
	Timestamp  string  `json:"timestamp"`
}

type zone_state_activity_data_points struct {
	HeatingPower zone_state_adp_heating_power `json:"heatingPower"`
}

type zone_state_next_schedule_change struct {
	Start   string             `json:"start"`
	Setting zone_state_setting `json:"setting"`
}

type zone_state struct {
	Setting            zone_state_setting              `json:"setting"`
	OverlayType        string                          `json:"overlayType"`
	Overlay            zone_state_overlay              `json:"overlay"`
	SensorDataPoints   zone_state_sensor_data_points   `json:"sensorDataPoints"`
	ActivityDataPoints zone_state_activity_data_points `json:"activityDataPoints"`
	NextScheduleChange zone_state_next_schedule_change `json:"nextScheduleChange"`
}

// Main function
func main2() {

	fmt.Println("!... Hello Heatcontroller 2 ...!")

	client := http.Client{}

	fmt.Println("-----------------------------------------------------------------------------------")
	token_obj, err := retrier(getTokenI, client, 15, 5)
	if(err != nil){
		log.Println(err)
		panic(err)
	}
	printToken(token_obj)
	/*
		fmt.Println("\n\nAccessToken:", token_obj.AccessToken)
		fmt.Println("\n\nRefreshToken:", token_obj.RefreshToken)
		fmt.Println("\n\nTokenType:", token_obj.TokenType)
		fmt.Println("\n\nExpiresIn:", token_obj.ExpiresIn)
		fmt.Println("\n\nScope:", token_obj.Scope)
		fmt.Println("\n\nJTI:", token_obj.JTI)
	*/

	fmt.Println("-----------------------------------------------------------------------------------")
	me_obj := getMe(client, token_obj)
	printMe(me_obj)
	/*
		fmt.Println("\n\nName:", me_obj.Name)
		fmt.Println("\n\nUserName:", me_obj.UserName)
		fmt.Println("\n\nId:", me_obj.Id)
		fmt.Println("\n\nEMail:", me_obj.EMail)
		fmt.Println("\n\nHomes:", me_obj.Homes)
		fmt.Println("\n\nHomes[0].Name:", me_obj.Homes[0].Name)
		fmt.Println("\n\nHomes[0].Id:", me_obj.Homes[0].Id)
	*/

	my_home := me_obj.Homes[0]

	zones_obj := getZones(client, token_obj, my_home.Id)
	// fmt.Println("\n\nZones:", zones_obj)

	fmt.Println("-----------------------------------------------------------------------------------")
	wohnzimmer_zone := getZone("Wohnzimmer", zones_obj)
	printZone(wohnzimmer_zone)

	fmt.Println("-----------------------------------------------------------------------------------")
	wohnzimmer_zone_state := getZoneState(client, token_obj, wohnzimmer_zone.Id, my_home.Id)
	printZoneState(wohnzimmer_zone.Name, wohnzimmer_zone_state)

	fmt.Println("-----------------------------------------------------------------------------------")
	kueche_zone := getZone("Küche", zones_obj)
	printZone(kueche_zone)

	fmt.Println("-----------------------------------------------------------------------------------")
	kueche_zone_state := getZoneState(client, token_obj, kueche_zone.Id, my_home.Id)
	printZoneState(kueche_zone.Name, kueche_zone_state)

	fmt.Println("-----------------------------------------------------------------------------------")
	kueche_zone_overlay := getZoneOverlay(client, token_obj, kueche_zone.Id, my_home.Id)
	printOverlay(kueche_zone.Name, kueche_zone_overlay)

	fmt.Println("-----------------------------------------------------------------------------------")
	setting := zone_state_setting{}
	setting.Power = "ON"
	setting.Type = "HEATING"
	setting.Temperature.Celsius = 20.0
	putZoneOverlay(client, token_obj, kueche_zone.Id, my_home.Id, setting)

	fmt.Println("-----------------------------------------------------------------------------------")
	kueche_zone_overlay2 := getZoneOverlay(client, token_obj, kueche_zone.Id, my_home.Id)
	printOverlay(kueche_zone.Name, kueche_zone_overlay2)


	my_home_obj := getMyHome(client, token_obj)

	fmt.Println("My Home :")
	printMyHome(&my_home_obj)
	fmt.Println("Done")

}

type my_home_zone struct {
	Id int
	Name string
//	typ string
	ZoneState zone_state
}

type my_home struct {
	Id int
	Name string
	PickTimeSecs int64
	Zones map[string]my_home_zone
}

func getMyHome(client http.Client, token_obj token) my_home {
	fmt.Println("-----------------------------------------------------------------------------------")
	fmt.Println("Get my home:")

	me_obj := getMe(client, token_obj)
	printMe(me_obj)
//	home := me_obj.Homes[0]

	my_home := my_home{
		Id : me_obj.Homes[0].Id,
		Name : me_obj.Homes[0].Name,
		PickTimeSecs : time.Now().Unix(),
		Zones : make(map[string]my_home_zone),
	}

	zones_obj := getZones(client, token_obj, my_home.Id)
	// fmt.Println("\n\nZones:", zones_obj)

	for _, e := range zones_obj {
		zone_obj := getZone(e.Name, zones_obj)
		zone_state_obj := getZoneState(client, token_obj, zone_obj.Id, my_home.Id)
		my_home_zone_obj := my_home_zone{
			Id : zone_obj.Id,
			Name : zone_obj.Name,
			ZoneState : zone_state_obj,
		}
		my_home.Zones[e.Name] = my_home_zone_obj
	}

	return my_home
}

func printMyHome(mh *my_home) {
	js, _ := json.MarshalIndent(mh, "", "   ")
	fmt.Println(string(js))
}

func printZone(zs *zone) {
	js, _ := json.MarshalIndent(zs, "", "   ")
	fmt.Println(string(js))
}

func printZoneState(name string, zs zone_state) {
	fmt.Println("\nzone-state '", name, "':")
	js, _ := json.MarshalIndent(zs, "", "   ")
	fmt.Println(string(js))
}

func printMe(m me) {
	js, _ := json.MarshalIndent(m, "", "   ")
	fmt.Println(string(js))
}

func printToken(m token) {
	js, _ := json.MarshalIndent(m, "", "   ")
	fmt.Println("Token : " + string(js))
}

func printOverlay(name string, m zone_state_overlay) {
	fmt.Println("\nzone-overlay '", name, "':")
	js, _ := json.MarshalIndent(m, "", "   ")
	fmt.Println(string(js))
}

func printSettingWrapper(m zone_state_setting_wrapper) {
	js, _ := json.MarshalIndent(m, "", "   ")
	fmt.Println(string(js))
}

func printSetting(m zone_state_setting) {
	js, _ := json.MarshalIndent(m, "", "   ")
	fmt.Println(string(js))
}

func putZoneOverlay(client http.Client, token_obj token, zone_id int, home_id int, setting zone_state_setting) {
	zone_id_str := strconv.Itoa(zone_id)
	home_id_str := strconv.Itoa(home_id)

	url := "https://my.tado.com/api/v2/homes/" + home_id_str + "/zones/" + zone_id_str + "/overlay"
	fmt.Println("\nURL(PUT):" + url)
	fmt.Println("Setting:")
	//	printSetting(setting)

	wrapper := zone_state_setting_wrapper{}
	wrapper.Setting = setting
	printSettingWrapper(wrapper)

	//str := "{\"setting\": {	\"type\": \"HEATING\",\"power\": \"ON\",\"temperature\": {\"celsius\": 20,\"fahrenheit\": 0}}}"
	//reader := strings.NewReader(str)

	/*
		buf := new(strings.Builder)
		_, err := io.Copy(buf, reader)
		// check errors
		fmt.Println(buf.String())
	*/

	js, _ := json.Marshal(wrapper)
	//fmt.Println("js:\n", js)
	reader := bytes.NewReader(js)

	req, getErr2 := http.NewRequest("PUT", url, reader)
	if getErr2 != nil {
		log.Println(getErr2)
	}

	req.Header.Set("Authorization", "Bearer "+token_obj.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	//	req.Header = http.Header{"Authorization": {"Bearer " + token_obj.AccessToken}, "Content-Type": {"application/json"}}
	_, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	//fmt.Println("PUT resp:\n", resp)
}

func getZoneOverlay(client http.Client, token_obj token, zone_id int, home_id int) zone_state_overlay {
	zone_id_str := strconv.Itoa(zone_id)
	home_id_str := strconv.Itoa(home_id)

	url := "https://my.tado.com/api/v2/homes/" + home_id_str + "/zones/" + zone_id_str + "/overlay"
	fmt.Println("\nURL:" + url)

	req, getErr2 := http.NewRequest("GET", url, nil)
	if getErr2 != nil {
		log.Println(getErr2)
	}

	req.Header = http.Header{"Authorization": {"Bearer " + token_obj.AccessToken}}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
	}

	zone_state_overlay_obj := zone_state_overlay{}
	jsonErr := json.Unmarshal(body, &zone_state_overlay_obj)
	if jsonErr != nil {
		log.Println(jsonErr)
	}

	return zone_state_overlay_obj
}

func getZoneState(client http.Client, token_obj token, zone_id int, home_id int) zone_state {
	zone_id_str := strconv.Itoa(zone_id)
	// fmt.Println("\nzone_id:", zone_id_str)

	home_id_str := strconv.Itoa(home_id)
	// fmt.Println("\nhome_id:", home_id_str)

	url := "https://my.tado.com/api/v2/homes/" + home_id_str + "/zones/" + zone_id_str + "/state"
	fmt.Println("\nURL:" + url)

	req, getErr2 := http.NewRequest("GET", url, nil)
	if getErr2 != nil {
		log.Println(getErr2)
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + token_obj.AccessToken},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
	}
	// bodyStr := string(body)
	// fmt.Println(bodyStr)

	zone_state_obj := zone_state{}
	jsonErr := json.Unmarshal(body, &zone_state_obj)
	if jsonErr != nil {
		log.Println(jsonErr)
	}

	return zone_state_obj
}

func getZone(name string, zones []zone) *zone {
	for i, e := range zones {
		if strings.ToLower(e.Name) == strings.ToLower(name) {
			return &zones[i]
		}
	}
	return nil
}

func getZones(client http.Client, token_obj token, home_id int) []zone {
	home_id_str := strconv.Itoa(home_id)
	// fmt.Println("\nhome_id:", home_id_str)

	url := "https://my.tado.com/api/v2/homes/" + home_id_str + "/zones"
	fmt.Println("\nURL:" + url)

	req, getErr2 := http.NewRequest("GET", url, nil)
	if getErr2 != nil {
		log.Println(getErr2)
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + token_obj.AccessToken},
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
	}
	// bodyStr := string(body)
	// fmt.Println(bodyStr)

	zone_objs := []zone{}
	jsonErr := json.Unmarshal(body, &zone_objs)
	if jsonErr != nil {
		log.Println(jsonErr)
	}

	return zone_objs
}

func getMe(client http.Client, token_obj token) me {

	req, getErr2 := http.NewRequest("GET", "https://my.tado.com/api/v2/me", nil)
	if getErr2 != nil {
		log.Println(getErr2)
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + token_obj.AccessToken},
	}

	//	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	// fmt.Println(resp.Body)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
	}
	// bodyStr := string(body)
	// fmt.Println(bodyStr)

	me_obj := me{}
	jsonErr := json.Unmarshal(body, &me_obj)
	if jsonErr != nil {
		log.Println(jsonErr)
	}

	return me_obj
}

type Credentials struct {
	Name string `json:"Name"`
	PW  string `json:"PW"`
}

func readCredentials() Credentials {
	js, err := ioutil.ReadFile("credentials/my_secrets.txt")  
    if err != nil {
        fmt.Println("Err")
    }
	var cred Credentials
	json.Unmarshal(js, &cred)
	return cred
}

func retrier[G any](fn func(http.Client)(G,error), client http.Client, retry int, start int) (G, error) {
	var token G
	for i:=0; i<retry; i++ {               
		token, err := fn(client)
		if err == nil {
			return token, nil
		}
		log.Println(err)
		waiting_time := time.Second * time.Duration(start * i)
        fmt.Println("Error: retrying, waiting ", waiting_time, " seconds ...")    
		time.Sleep(waiting_time)   
    }  
	log.Fatal("Tried many times, but something is not working, check the logs!")
	return token, nil
}

func getTokenI(client http.Client) (token,error) {

	//var js = []byte("{\"Name\":\"Dinesh Krishnan\",\"PW\":\"apassword\"}")
	var cred Credentials = readCredentials()

	url := "https://auth.tado.com/oauth/token?client_id=tado-web-app&grant_type=password&scope=home.user&username=" + cred.Name + "&password=" + cred.PW + "&client_secret=wZaRN7rpjn3FoNyF5IFuxg9uMzYJcvOoQ8QWiIqS3hfk6gLhVlG57j5YNoZL2Rtc"
	req, getErr2 := http.NewRequest("POST", url, nil)
	if getErr2 != nil {
		log.Println(getErr2)
		return token{}, getErr2
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return token{}, err
	}
	// fmt.Println(resp.Body)

	body, readErr := ioutil.ReadAll(resp.Body)
	if readErr != nil {
		log.Println(readErr)
		return token{}, readErr
	}
	// fmt.Println(string(body))

	token_obj := token{}
	jsonErr := json.Unmarshal(body, &token_obj)
	if jsonErr != nil {
		log.Println(jsonErr)
		return token{}, jsonErr
	}
	// fmt.Println("Token : " + string(token_obj))

	return token_obj,nil
}
