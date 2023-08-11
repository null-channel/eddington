package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sendgrid/sendgrid-go"
)

type MarketingController struct {
	apiKey string
}

// New creates a new MarketingController responsible for marketing endpoints
func New(apiKey string, router *mux.Router) MarketingController {
	fmt.Println("New MarketingController")

	mc := MarketingController{}

	mc.AddAllControllers(router)

	return MarketingController{apiKey: apiKey}
}

func (m *MarketingController) AddAllControllers(router *mux.Router) {
	router.HandleFunc("/email", m.POSTEmailSubscriber).Methods(http.MethodPost)
}

//	@BasePath	/api/v1/

// CreateUser godoc
//
//	@Summary	Create an user
//	@Schemes
//	@Description	create a user
//	@Tags			Users
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/marketing/email [post]
func (m *MarketingController) POSTEmailSubscriber(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	fmt.Println("email: " + email)
	err = m.Addrecipients(email)

	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "New e-mail added to list successfully!")
}

// GetIDfromEmail : Get ID from email
/*
{
  "result": {
    "jane_doe@example.com": {
      "contact": {
        "address_line_1": "",
        "address_line_2": "",
        "alternate_emails": [
          "janedoe@example1.com"
        ],
        "city": "",
        "country": "",
        "email": "jane_doe@example.com",
        "first_name": "Jane",
        "id": "asdf-Jkl-zxCvBNm",
        "last_name": "Doe",
        "list_ids": [],
        "segment_ids": [],
        "postal_code": "",
        "state_province_region": "",
        "phone_number": "",
        "whatsapp": "",
        "line": "",
        "facebook": "",
        "unique_name": "",
        "custom_fields": {},
        "created_at": "2021-03-02T15:25:47Z",
        "updated_at": "2021-03-30T15:26:16Z",
        "_metadata": {
          "self": "<metadata_url>"
        }
      }
    },
}
*/
func (m *MarketingController) GetIDfromEmail(email string) (string, error) {
	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest(m.apiKey, "/v3/marketing/contacts/search/emails", host)
	request.Method = http.MethodPost
	v := EmailList{Emails: []string{email}}
	b, e := json.Marshal(v)
	if e != nil {
		fmt.Println(e)
		return "", e
	}
	request.Body = b

	fmt.Println("request: " + string(request.Body))

	response, err := sendgrid.MakeRequest(request)
	//response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
		return "", err
	}
	var result map[string]interface{}

	fmt.Println(response.StatusCode)
	fmt.Println(response.Body)
	fmt.Println(response.Headers)

	err = json.Unmarshal([]byte(response.Body), &result)
	if err != nil {
		log.Println(err)
		return "", err
	}

	result = result["result"].(map[string]interface{})

	for key, result := range result {
		person := result.(map[string]interface{})
		//just the first e-mail is fine. e-mail is all we have to go off of right now anyhow.
		if key == email {
			return person["id"].(string), nil
		}
	}

	return "", fmt.Errorf("no ID found for email: %s", email)
}

func (m *MarketingController) AddRecipientToWaitingList(recipientID string) error {
	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest(m.apiKey, fmt.Sprintf("/v3/marketing/contacts/lists/bdd5bf34-a5ba-43a5-b24a-e098b2ae3b68/recipients/%s", recipientID), host)
	request.Method = http.MethodPost
	response, err := sendgrid.MakeRequest(request)
	//response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
		return err
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
		return nil
	}

}

// Addrecipients : Add recipients
// POST /contactdb/recipients
func (u *MarketingController) Addrecipients(email string) error {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest(apiKey, "/v3/marketing/contacts", host)
	request.Method = http.MethodPut
	v := RecipientList{Contacts: []Contact{{Email: email}}}
	b, e := json.Marshal(v)
	if e != nil {
		fmt.Println(e)
		return e
	}
	request.Body = b
	response, err := sendgrid.MakeRequest(request)
	if err != nil {
		fmt.Println("addrecipients request finished with error ------")
		log.Println(err)
		return err
	} else {
		fmt.Println("Addrecipients worked ------")
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
	//TODO just return the id once I figure out how to parse the response
	return nil
}

type RecipientList struct {
	Contacts []Contact `json:"contacts"`
}

type EmailList struct {
	Emails []string `json:"emails"`
}

type Contact struct {
	Email string `json:"email"`
}
