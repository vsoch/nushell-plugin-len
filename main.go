// Copyright 2019 Vanessa Sochat. All rights reserved.
// Use of this source code is governed by the Polyform Strict license
// that can be found in the LICENSE file and available at
// https://polyformproject.org/licenses/noncommercial/1.0.0

package main
 
import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)


// NushellPlugin represents an interface for a Nushell Plugin. It includes 
// a configuration, along with supporting functions
type NushellPlugin struct {
	Config	ConfigParams
}

// A set of base params is string[string]
type Params map[string]string
type IntParams map[string]int

// Returning a length, must be int - must be customized for the plugin type
type ResponseParams map[string]IntParams

// A set of params is a map
type ConfigParams struct {
	Name	string	`json:"name"`
	Usage	string	`json:"usage"`
	Positionals	[]string	`json:"positionals"`
	Named	Params	`json:"named"`
	IsFilter	bool	`json:"is_filter"`
}

// StringResponse is nested
type StringResponse struct {
	Item struct {
    		Primitive struct {
			String string `json:"String"`
		} `json:"Primivite"`
	} `json:"item"`
}

// JsonResponse is a standard Json Response as defined by:
// https://www.jsonrpc.org/specification#response_object
// where the params are a single dictionary of params
type JsonResponse struct {
	Jsonrpc string	`json:"jsonrpc"`	// jsonrpc version, e.g., 2.0
	Method string	`json:"method"`		// method, e.g., response
	Params *ResponseParams	`json:"params"`	// arbitrary params
}


// configure will add the configuration to the plugin, akin to an init.
// Here we primarily provide the plugin name, usage, and arguments, but
// other functions could be performed. It takes no arguments,
// as it's assumed that a plugin provides only one config
func (plugin *NushellPlugin) configure() {
	var config = ConfigParams{
		Name: "len",
		Usage: "Return the length of a string",
		Named: Params{},
		Positionals: make([]string, 0),
		IsFilter: true}
	plugin.Config = config
}

// getLength of a string value, nested at stringValue["item"]["Primitive"]["String"]
func (plugin *NushellPlugin) getLength(stringValue interface{}) int {

	// I hope there is a more elegant way to do this
	jsonValues := stringValue.(map[string]interface{})
	item := jsonValues["item"].(map[string]interface{})
	primitive := item["Primitive"].(map[string]interface{})
	finalString := primitive["String"].(string)
	length := len(finalString)
	return length
}


// printGoodResponse will print a json response to the terminal. The 
// status would typically be "Ok" for a good response.
func (plugin *NushellPlugin) printGoodResponse(params IntParams) error {

	responseParams := &ResponseParams{"Ok": params}
	
	// Wrap params in json response
	response := &JsonResponse{
		Jsonrpc: "2.0",
		Method: "response",
		Params: responseParams}

	// Serialize the struct to json, exit out if there is an error
	jsonString, err := json.Marshal(response) 
	if err != nil {
		return err
	}

	// Write the response to stdout
	fmt.Println(string(jsonString))
	return nil
}

// printConfigResponse will print the config json response to the terminal.
func (plugin *NushellPlugin) printConfigResponse() error {

	// Serialize the struct to json, exit out if there is an error
	jsonString, err := json.Marshal(plugin.Config) 
	if err != nil {
		return err
	}

	// Write the response to stdout
	fmt.Println(string(jsonString))
	return nil
}


func main() {

	// Set up temporary logger
	f, err := os.OpenFile("/tmp/nu_plugin_len.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()	
	logger := log.New(f, "nu_plugin_len ", log.LstdFlags)

	// Instantiate new plugin interface, generate config
	plugin := NushellPlugin{}
	plugin.configure()

	// Read into json decoded
	decoder := json.NewDecoder(os.Stdin)

	line := make(map[string]interface{})

	for {
		err := decoder.Decode(&line) 
		if err != nil {
			fmt.Errorf("unable to read json: %s", err)
		} 

		// look for a method in the line
		if method, ok := line["method"]; ok {	
	
			// Case 1: method is config
			if method == "config" {
				logger.Println("Request for config", line)
			        plugin.printConfigResponse()
				break

			} else if method == "begin_filter" {
				logger.Println("Request for begin filter", line)
				emptyResponse := make([]string, 0)
				fmt.Println(emptyResponse)

			} else if method == "filter" {
				logger.Println("Request for filter", line)
				if params, ok := line["params"]; ok {
					intLength := plugin.getLength(params)
					value := IntParams{"Value": intLength}
					plugin.printGoodResponse(value)
				}

			} else if method == "end_filter" {
				logger.Println("Request for end filter", line)
				emptyResponse := make([]string, 0)
				fmt.Println(emptyResponse)
				break
			}

		} else {
			break
		}
	}
}
