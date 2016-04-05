package estimoteslack

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/urlfetch"
	"html/template"
	"net/http"
	"net/url"
	"time"
)

type WebhookURL struct {
  PropName string
  PropValue string
}

type Asset struct {
	BeaconId   string
	BeaconName string
	AssetId    string
	AssetName  string
	AppType    string
	AppData    string
	Status     string
}

type Issue struct {
	AssetId       string
	AssetName     string
	IssueType     string
	IssueDateTime time.Time
	IssueDetails  string
	IssueRaisedBy string
	Status        string
}

type APIResponse struct {
	ErrorCode string
	ErrorMsg  string
}

type Message struct {
	Channel    string `json:"channel"`
	Text       string `json:"text"`
	Username   string `json:"username"`
	Icon_emoji string `json:"icon_emoji"`
}

func init() {
	http.HandleFunc("/", homepagehandler)
	http.HandleFunc("/assets", assetshandler)
	http.HandleFunc("/issues", issueshandler)
	http.HandleFunc("/listissues", issuespagehandler)
	http.HandleFunc("/addAsset", addAsset)
	http.HandleFunc("/raiseIssue", raiseIssue)
  http.HandleFunc("/addSlackTeamWebhook",addSlackTeamWebhook)
  http.HandleFunc("/SlackTeamWebhook",getSlackTeamWebhook)
  http.HandleFunc("/clearData",clearData)
}

//This handler is for the home page. It shows the list of assets
func homepagehandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Asset").Filter("Status=", "ACTIVE")
	var results []Asset
	_, err := q.GetAll(c, &results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = assetsListTemplate.Execute(w, results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

//Clear datastore
func clearData(w http.ResponseWriter, r *http.Request) {
   c := appengine.NewContext(r)
   q1 := datastore.NewQuery("Asset")
  var results1 []Asset
  keys1,err := q1.GetAll(c,&results1)
  if err != nil {
               http.Error(w, err.Error(), http.StatusInternalServerError)
               return
  }
  datastore.DeleteMulti(c,keys1)

  q2 := datastore.NewQuery("Issue")
  var results2 []Issue
  keys2,err := q2.GetAll(c,&results2)
  if err != nil {
               http.Error(w, err.Error(), http.StatusInternalServerError)
               return
  }
  datastore.DeleteMulti(c,keys2)

}

//Show list of Issues
func issuespagehandler(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	q := datastore.NewQuery("Issue").Filter("Status=", "ACTIVE")
	var results []Issue
	_, err := q.GetAll(c, &results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = issuesListTemplate.Execute(w, results); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

//API - List of Assets
func assetshandler(w http.ResponseWriter, r *http.Request) {

	c := appengine.NewContext(r)
	q := datastore.NewQuery("Asset")
	beaconname := r.FormValue("BeaconName")

	if beaconname != "" {
		q = datastore.NewQuery("Asset").Filter("Status=", "ACTIVE").Filter("BeaconName=", beaconname)
	} else {
		q = datastore.NewQuery("Asset").Filter("Status=", "ACTIVE")
	}

	var results []Asset
	_, err := q.GetAll(c, &results)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, e := json.Marshal(results)
	if e != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(b))
}

//API - List of Issues
func issueshandler(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
      q := datastore.NewQuery("Issue").Filter("Status=","ACTIVE")
    var results []Issue
    _,err := q.GetAll(c,&results)
    if err != nil {
                  http.Error(w, err.Error(), http.StatusInternalServerError)
                  return
          }

    b,e := json.Marshal(results)
    if e != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
                  return
    }
    fmt.Fprint(w,string(b))
}

func webhookURLKey(c context.Context, url string) *datastore.Key {
    return datastore.NewKey(c, "WebhookURL", "webhookurlkey", 0, nil)
}

//Add Slack Team Webhook url
func addSlackTeamWebhook(w http.ResponseWriter, r *http.Request) {
  webhookurl := r.FormValue("webhookurl")
  c := appengine.NewContext(r)
  wurl := WebhookURL {
     PropName  : "webhookurl",
     PropValue : webhookurl,
  }

key := webhookURLKey(c,webhookurl)
  _, err := datastore.Put(c, key, &wurl)
  if err != nil {
          http.Error(w, err.Error(), http.StatusInternalServerError)
          return
  }

resp := APIResponse{
  ErrorCode : "200",
  ErrorMsg  : "",
}
b,e := json.Marshal(&resp);
if e != nil {}
fmt.Fprint(w,string(b))

}

func getWebhook(c context.Context) string {
  var wurl WebhookURL
  q := datastore.NewQuery("WebhookURL")

    var results []WebhookURL
    keys,err := q.GetAll(c,&results)

    err = datastore.Get(c, keys[0], &wurl)
    if err != nil {

                  return ""
          }

    return wurl.PropValue

}

//API - Slack Team URL
func getSlackTeamWebhook(w http.ResponseWriter, r *http.Request) {

  var wurl WebhookURL
  c := appengine.NewContext(r)
  q := datastore.NewQuery("WebhookURL")

    var results []WebhookURL
    keys,err := q.GetAll(c,&results)

    err = datastore.Get(c, keys[0], &wurl)
    if err != nil {
                  http.Error(w, err.Error(), http.StatusInternalServerError)
                  return
          }

    fmt.Fprint(w,wurl.PropValue)
}

//Creates a new Asset Key
func assetKey(c context.Context, assetid string) *datastore.Key {
	return datastore.NewKey(c, "Asset", assetid, 0, nil)
}

//JSON API to insert Asset
func addAsset(w http.ResponseWriter, r *http.Request) {

	assetid := r.FormValue("AssetId")
	assetname := r.FormValue("AssetName")
	beaconid := r.FormValue("BeaconId")
	beaconname := r.FormValue("BeaconName")
	apptype := r.FormValue("AppType")
	appdata := r.FormValue("AppData")
	status := r.FormValue("Status")

	//Need to either insert or update the Asset

	c := appengine.NewContext(r)
	a := Asset{
		AssetId:    assetid,
		AssetName:  assetname,
		BeaconId:   beaconid,
		BeaconName: beaconname,
		AppType:    apptype,
		AppData:    appdata,
		Status:     status,
	}

	key := assetKey(c, assetid)
	_, err := datastore.Put(c, key, &a)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := APIResponse{
		ErrorCode: "200",
		ErrorMsg:  "",
	}
	b, e := json.Marshal(&resp)
	if e != nil {
	}
	fmt.Fprint(w, string(b))
}

//Creates a new Issue Key
func issueKey(c context.Context, issueid string) *datastore.Key {
	return datastore.NewKey(c, "Issue", issueid, 0, nil)
}

//JSON API to insert Issue
func raiseIssue(w http.ResponseWriter, r *http.Request) {

	assetid := r.FormValue("AssetId")
	assetname := r.FormValue("AssetName")
	issuetype := r.FormValue("IssueType")
	issuedetails := r.FormValue("IssueDetails")
	issueraisedby := r.FormValue("IssueRaisedBy")
	t := time.Now()
	issuedatetime := t
	status := r.FormValue("Status")

	//Need to either insert Issue

	c := appengine.NewContext(r)
	i := Issue{
		AssetId:       assetid,
		AssetName:     assetname,
		IssueType:     issuetype,
		IssueDetails:  issuedetails,
		IssueDateTime: issuedatetime,
		IssueRaisedBy: issueraisedby,
		Status:        status,
	}

	key := issueKey(c, string(t.Unix()))
	_, err := datastore.Put(c, key, &i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Send Slack Notification

	m := Message{"#general", "***************************************" + "\n New Problem Report" + "\n Asset Id: " + i.AssetId + "\n Asset Name :" + i.AssetName + "\n Issue Type: " + i.IssueType + "\n Reported By: " + i.IssueRaisedBy + "\n Details: " + i.IssueDetails + "\n Reported on: " + i.IssueDateTime.Format("Mon Jan _2 15:04:05 2006") + "\n ************************", "HelpDesk Bot", ":computer:"}
	b, err := json.Marshal(m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := urlfetch.Client(c)

	v := url.Values{}
	v.Set("payload", string(b))
	_, err = client.PostForm(incoming_webhook_url, v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := APIResponse{
		ErrorCode: "200",
		ErrorMsg:  "",
	}
	b, e := json.Marshal(&resp)
	if e != nil {
	}
	fmt.Fprint(w, string(b))
}

//Asset List Template
var assetsListTemplate = template.Must(template.New("assetsListTemplate").Parse(`
<html>
  <head>
    <title>Assets List</title>
  </head>
  <body>
    <h1>Assets List</h1>
    <table>
	<tr><th>AssetId</th><th>AssetName</th><th>BeaconId</th><th>BeaconName</th><th>AppType</th><th>AppData</th><th>Status</th></tr>
    {{range .}}
      <tr>
	  <td>{{.AssetId}}</td>
	  <td>{{.AssetName}}</td>
    <td>{{.BeaconId}}</td>
	  <td>{{.BeaconName}}</td>
    <td>{{.AppType}}</td>
    <td>{{.AppData}}</td>
    <td>{{.Status}}</td>
	  </tr>
    {{end}}
  </body>
</html>
`))

//Issues List Template
var issuesListTemplate = template.Must(template.New("issuesListTemplate").Parse(`
<html>
  <head>
    <title>Issues List</title>
  </head>
  <body>
    <h1>Issues List</h1>
    <table>
	<tr><th>AssetId</th><th>AssetName</th><th>IssueType</th><th>IssueDetails</th><th>IssueDateTime</th><th>IssueRaisedBy</th><th>Status</th></tr>
    {{range .}}
      <tr>
	  <td>{{.AssetId}}</td>
	  <td>{{.AssetName}}</td>
    <td>{{.IssueType}}</td>
    <td>{{.IssueDetails}}</td>
	  <td>{{.IssueDateTime}}</td>
    <td>{{.IssueRaisedBy}}</td>
    <td>{{.Status}}</td>
	  </tr>
    {{end}}
  </body>
</html>
`))
