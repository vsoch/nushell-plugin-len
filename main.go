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

// Tag is passed from stream during filter, we parse and pass forward
type Tag struct {
	Anchor interface{}	`json:"anchor"`
	Span map[string]int	`json:"span"`
}

// LengthResponse is nested under Primitive -> Int -> value
type LengthResponse struct {
	Item struct {
    		Primitive struct {
			Int int `json:"Int"`
		} `json:"Primitive"`
	} `json:"item"`
	Tag Tag	`json:"tag"`
}

// NestedParams support inner {"Value": LengthResponse} with Item.Primitive.Int
type NestedParams map[string]LengthResponse

// ResponseParams support inner {"Ok": NestedParams}
type ResponseParams map[string]NestedParams

// ArrayResponseParams are an array of []ResponseParams for JsonResponse.Params
type ArrayResponseParams []ResponseParams

// FinalResponseParams wrap ArrayResponseParams: {"Ok": ArrayResponseParams}
type FinalResponseParams map[string]ArrayResponseParams

// Params are general string[string] map for the config
type Params map[string]string

// ConfigParams holds configuration information for the plugin
type ConfigParams struct {
	Name	string			`json:"name"`
	Usage	string			`json:"usage"`
	Positional	[]string	`json:"positional"`
	RestPositional []string		`json:"rest_positional"`
	Named	Params			`json:"named"`
	IsFilter	bool		`json:"is_filter"`
}

// ConfigResponseParams add another level of nesting for {"Ok": ConfigParams}
type ConfigResponseParams map[string]ConfigParams


// ArrayParams are intended for end_filter and begin_filter
type ArrayParams []string

// EmptyResponseParams supports {"Ok": ArrayParams} (e.g., {"Ok": []})
type EmptyResponseParams map[string]ArrayParams

// JsonResponse is a standard Json Response as defined by:
// https://www.jsonrpc.org/specification#response_object
// Returns our FinalResponseParams as JsonResponse.Params
type JsonResponse struct {
	Jsonrpc string			`json:"jsonrpc"`	// jsonrpc version, e.g., 2.0
	Method string			`json:"method"`		// method, e.g., response
	Params FinalResponseParams	`json:"params"`		// arbitrary params
}

// ConfigResponse is specifically to return Jsonrpc with Params.ConfigResponseParams
type ConfigResponse struct {
	Jsonrpc string			`json:"jsonrpc"`
	Method string			`json:"method"`
	Params *ConfigResponseParams	`json:"params"`
}

// ArrayResponse returns an ArrayParams instead
type ArrayResponse struct {
	Jsonrpc string			`json:"jsonrpc"`
	Method string			`json:"method"`
	Params *EmptyResponseParams	`json:"params"`
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
		Positional: make([]string, 0),
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


// getTag from the filter input to return in response
func (plugin *NushellPlugin) getTag(stringValue interface{}) Tag {
	
	// I hope there is a more elegant way to do this
	jsonValues := stringValue.(map[string]interface{})

	// Create the span to hold a start and end, and empty tag
	span := map[string]int{}
	tag := Tag{}

	// Real nushell invokation will include a tag
	if tagGroup, ok := jsonValues["tag"].(map[string]interface{}); ok {

		spanGroup := tagGroup["span"].(map[string]interface{})

		// Create the span, it has a start and end
		span["start"] = int(spanGroup["start"].(float64))
		span["end"] = int(spanGroup["end"].(float64))
		tag.Span = span

		// If anchor isn't nil, add it (not sure if type is correct)
		if anchor, ok := tagGroup["anchor"].(interface{}); ok {
			tag.Anchor = anchor
		}

	// Otherwise generate a dummy one for local testing
	} else {
		span["start"] = 0
		span["end"] = 0
		tag.Span = span
	}

	return tag
}


// printGoodResponse will print a json response to the terminal. The 
// generates a JsonResponse with Params.FinalResponseParams
func (plugin *NushellPlugin) printGoodResponse(length int, tag Tag) error {

	lengthResponse := LengthResponse{}
	lengthResponse.Item.Primitive.Int = length
	lengthResponse.Tag = tag

	nestedParams := NestedParams{"Value": lengthResponse}
	responseParams := ResponseParams{"Ok": nestedParams}

	arrayResponseParams := ArrayResponseParams{}
	arrayResponseParams = append(arrayResponseParams, responseParams)
	finalResponseParams := FinalResponseParams{"Ok": arrayResponseParams}
	
	// Wrap params in json response
	response := &JsonResponse{
		Jsonrpc: "2.0",
		Method: "response",
		Params: finalResponseParams}

	// Serialize the struct to json, exit out if there is an error
	jsonString, err := json.Marshal(response) 
	if err != nil {
		return err
	}

	// Write the response to stdout
	fmt.Println(string(jsonString))
	return nil
}

// printEmptyResponse will print an ArrayResponse that is empty.
// the intende use case is for an end_filter or start_filter
// generates an ArrayResponse with Params.EmptyResponseParams
func (plugin *NushellPlugin) printEmptyResponse() error {

	emptyArray := make([]string, 0)
	responseParams := &EmptyResponseParams{"Ok": emptyArray}
	
	// Wrap params in json response
	response := &ArrayResponse{
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
// generates an ConfigResponse with Params.ConfigResponseParams
func (plugin *NushellPlugin) printConfigResponse() error {

	responseParams := &ConfigResponseParams{"Ok": plugin.Config}
	
	// Wrap params in json response
	response := &ConfigResponse{
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
				plugin.printEmptyResponse()

			} else if method == "filter" {
				logger.Println("Request for filter", line)
				if params, ok := line["params"]; ok {
					intLength := plugin.getLength(params)
					tag := plugin.getTag(params)
					plugin.printGoodResponse(intLength, tag)
				}

			} else if method == "end_filter" {
				logger.Println("Request for end filter", line)
				plugin.printEmptyResponse()
				break
			}

		} else {
			break
		}
	}
}
