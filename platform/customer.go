package platform

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Customer struct {
	ID                        string              `json:"id"`
	Version                   int                 `json:"version"`
	VersionModifiedAt         string              `json:"versionModifiedAt"`
	LastMessageSequenceNumber int                 `json:"lastMessageSequenceNumber"`
	CreatedAt                 string              `json:"createdAt"`
	LastModifiedAt            string              `json:"lastModifiedAt"`
	LastModifiedBy            *LastModifiedBy     `json:"lastModifiedBy,omitempty"`
	CreatedBy                 *CreatedBy          `json:"createdBy,omitempty"`
	Email                     string              `json:"email"`
	FirstName                 string              `json:"firstName"`
	LastName                  string              `json:"lastName"`
	Addresses                 []Address           `json:"addresses"`
	ShippingAddressIds        []string            `json:"shippingAddressIds"`
	BillingAddressIds         []string            `json:"billingAddressIds"`
	IsEmailVerified           bool                `json:"isEmailVerified"`
	Stores                    []StoreKeyReference `json:"stores"`
	AuthenticationMode        string              `json:"authenticationMode"`
}
type CustomerList struct {
	Limit   int        `json:"limit"`
	Offset  int        `json:"offset"`
	Count   int        `json:"count"`
	Total   int        `json:"total"`
	Results []Customer `json:"results"`
}

type LastModifiedBy struct {
	ClientId         *string `json:"clientId,omitempty"`
	IsPlatformClient bool    `json:"IsPlatformClient,omitempty"`
}

type CreatedBy struct {
	ClientId         *string `json:"clientId,omitempty"`
	IsPlatformClient bool    `json:"IsPlatformClient,omitempty"`
}
type Address struct {
	ID                    string `json:"id"`
	Key                   string `json:"key"`
	Title                 string `json:"title"`
	Salutation            string `json:"salutation"`
	FirstName             string `json:"firstName"`
	LastName              string `json:"lastName"`
	StreetName            string `json:"streetName"`
	StreetNumber          string `json:"streetNumber"`
	AdditionalStreetInfo  string `json:"additionalStreetInfo"`
	PostalCode            string `json:"postalCode"`
	City                  string `json:"city"`
	Region                string `json:"region"`
	State                 string `json:"state"`
	Country               string `json:"country"`
	Company               string `json:"company"`
	Department            string `json:"department"`
	Building              string `json:"building"`
	Apartment             string `json:"apartment"`
	POBox                 string `json:"pOBox"`
	Phone                 string `json:"phone"`
	Mobile                string `json:"mobile"`
	Email                 string `json:"email"`
	Fax                   string `json:"fax"`
	AdditionalAddressInfo string `json:"additionalAddressInfo"`
	ExternalId            string `json:"externalId"`
}

type StoreKeyReference struct {
	Key string `json:"key"`
}

type NewCustomer struct {
	Email              string `json:"email"`
	FirstName          string `json:"firstName"`
	LastName           string `json:"lastName"`
	AuthenticationMode string `json:"authenticationMode"`
}

type CustomerResponse struct {
	Customer Customer `json:"customer"`
}

var customerPost NewCustomer

func (m *Client) Post(customer NewCustomer) *MethodHead {
	customerPost = customer
	return &MethodHead{
		EndPoint: fmt.Sprintf("%s/%s/customers", m.ApiUrl, m.ProjectKey),
		Method:   "POST",
	}
}

func (m *Client) Where(email string) *MethodHead {
	where := fmt.Sprintf("email in (\"%s\")", email)
	encodedWhere := url.QueryEscape(where)
	return &MethodHead{
		EndPoint: fmt.Sprintf("%s/%s/customers?where=%s", m.ApiUrl, m.ProjectKey, encodedWhere),
		Method:   "GET",
	}
}

func (m *Client) Get() *MethodHead {
	return &MethodHead{
		EndPoint: fmt.Sprintf("%s/%s/customers", m.ApiUrl, m.ProjectKey),
		Method:   "GET",
	}
}

func (m *MethodHead) Execute() ([]Customer, error) {
	var req *http.Request
	var err error
	switch m.Method {
	case "POST":
		customerJson, err := json.Marshal(customerPost)
		if err != nil {
			return nil, err
		}
		body := []byte(customerJson)
		req, err = http.NewRequest(m.Method, m.EndPoint, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		resp, err := clientDo(req)
		var customerResponse CustomerResponse
		err = json.NewDecoder(resp.Body).Decode(&customerResponse)
		if err != nil {
			fmt.Println("Error al decodificar la respuesta:", err)
			return nil, err
		}
		customerArray := []Customer{customerResponse.Customer}
		return customerArray, err

	case "GET":
		req, err = http.NewRequest(m.Method, m.EndPoint, nil)
		if err != nil {
			return nil, err
		}
		resp, err := clientDo(req)
		var customerList CustomerList
		err = json.NewDecoder(resp.Body).Decode(&customerList)
		if err != nil {
			fmt.Println("Error al decodificar la respuesta:", err)
			return nil, err
		}
		return customerList.Results, nil
	default:
		return nil, errors.New("No se asigno un metodo")
	}
}

func clientDo(req *http.Request) (resp *http.Response, err error) {
	client := &http.Client{}
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println("Error al realizar la solicitud:", err)
		return nil, err
	}

	if err != nil {
		fmt.Println("Error al leer el cuerpo de la respuesta:", err)
		return nil, err
	}

	return resp, err
}
