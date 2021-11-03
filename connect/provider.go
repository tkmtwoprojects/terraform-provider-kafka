package connect

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tkmtwoprojects/go-kafka/connect"

	"golang.org/x/oauth2/clientcredentials"
)

func Provider() *schema.Provider {
	log.Print("[DEBUG] Creating Provider")
	provider := schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKA_CONNECT_URL", ""),
				ValidateFunc: func(val interface{}, key string) (warns []string, errs []error) {
					v := val.(string)
					if !strings.HasSuffix(v, "/") {
						errs = append(errs, fmt.Errorf("Parameter '%q' value '%q' does not end in slash.", key, v))
					}
					return
				},
			},

			"basic_auth_username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"basic_auth_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth2_auth_clientid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth2_auth_clientsecret": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth2_auth_tokenurl": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth2_auth_params": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},

		ConfigureFunc: providerConfigure,
		ResourcesMap: map[string]*schema.Resource{
			"connect_connector": connectConnectorResource(),
		},
	}
	log.Printf("[DEBUG] Created provider: %v", provider)
	return &provider
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Print("[DEBUG] Initializing go-kafka/connect client")

	httpClient, httpClientErr := buildHttpClient(d)
	if httpClientErr != nil {
		log.Print("[ERROR] Error while building httpClient")
		return nil, httpClientErr
	} else {
		log.Print("[DEBUG] Created httpClient for use in go-kafka")
	}

	c, ncerr := connect.NewClient(httpClient, d.Get("url").(string))
	if ncerr != nil {
		log.Print("[ERROR] Error while creating new go-kafka client")
		return nil, ncerr
	} else {
		log.Print("[DEBUG] Created new go-kafka client")
	}

	return c, nil
}

//
// Doesn't feel like this should be necessary, but here we go...
//
// The paramter "oauth2_auth_params" is a TypeMap of strings.
// For each value, split on comma and pack the url.Values{}
// along the way.
//
func buildEndpointParams(d *schema.ResourceData) url.Values {
	eps := url.Values{}

	cm := d.Get("oauth2_auth_params").(map[string]interface{})

	/*
		if cm == nil {
			log.Printf("Using cm")
		} else {
			for x, y := range cm {
				log.Printf("LOOKS LIKE %q is a %T", x, y)
			}
		}
	*/

	for k, v := range cm {
		vs := strings.Split(v.(string), ",")
		for i := range vs {
			log.Printf("[DEBUG] Adding endpoint param '%q' -> '%q'", k, vs[i])
			eps.Add(k, vs[i])
		}
	}

	return eps
}

func buildHttpClient(d *schema.ResourceData) (*http.Client, error) {
	//
	//TODO: remove
	//
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	if (d.Get("oauth2_auth_clientid").(string) != "") &&
		(d.Get("oauth2_auth_clientsecret").(string) != "") &&
		(d.Get("oauth2_auth_tokenurl").(string) != "") {

		//
		// TODO: remove this from logging and move the decision
		// on what kind of client to go-kafka or another function.  This func and approach
		// assumed that we could set basic auth (or later mtls) here.
		// That assumption was probably wrong
		//
		log.Print("[DEBUG] Building oauth2 client")
		log.Print("URL   is         ", d.Get("url").(string))
		log.Print("CLIENTID is      ", d.Get("oauth2_auth_clientid").(string))
		log.Print("CLIENTSECRET is  ", d.Get("oauth2_auth_clientsecret").(string))
		log.Print("TOKENURL is      ", d.Get("oauth2_auth_tokenurl").(string))
		log.Print("PARAMS are       ", d.Get("oauth2_auth_params").(map[string]interface{}))

		cfg := clientcredentials.Config{
			ClientID:       d.Get("oauth2_auth_clientid").(string),
			ClientSecret:   d.Get("oauth2_auth_clientsecret").(string),
			TokenURL:       d.Get("oauth2_auth_tokenurl").(string),
			EndpointParams: buildEndpointParams(d),
		}
		return cfg.Client(context.Background()), nil

	}

	//return nil, nil
	return &http.Client{}, nil
}
