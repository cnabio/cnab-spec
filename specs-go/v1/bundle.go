package v1

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Bundle represent a CNAB descriptor
type Bundle struct {
	// The version of the CNAB specification. This should always be the integer 1 for this schema version.
	SchemaVersion string `json:"schemaVersion"`
	// The name of this bundle
	Name string `json:"name"`
	// A SemVer2 version for this bundle
	Version string `json:"version"`
	// A description of this bundle, intended for users
	Description string `json:"description,omitempty"`
	// A list of keywords describing the bundle, intended for users
	Keywords []string `json:"keywords,omitempty"`
	// A list of parties responsible for this bundle, with contact info
	Maintainers []Maintainer `json:"maintainers,omitempty"`
	// The SPDX license code or proprietary license name for this bundle
	License string `json:"license,omitempty"`
	// The array of invocation image definitions for this bundle
	InvocationImages []InvocationImage `json:"invocationImages"`
	// The application images installed by this bundle
	Images map[string]Image `json:"images,omitempty"`
	// Credentials to be injected into the invocation image
	Credentials map[string]Credential `json:"credentials,omitempty"`
	// Custom actions that can be triggered on this bundle
	Actions map[string]Action `json:"actions,omitempty"`
	// reserved for future usage
	Extensions map[string]interface{} `json:"extensions,omitempty"`
	// Parameters that can be injected into the invocation image
	Parameters map[string]Parameter `json:"parameters,omitempty"`
}

// Maintainer is an object that describes a maintainer
type Maintainer struct {
	// Name of party reponsible for this bundle
	Name string
	// Email address of responsible party
	Email string
	// URL of the responsible party, perhaps containing additional contact info
	URL string
}

// Action is a custom action executable on the bundle
type Action struct {
	// Must be set to true if the action can change any resource managed by this bundle
	Modifies bool
}

// Platform qualifies an image or invocation image target platform
type Platform struct {
	// The architecture of the image (i386, amd64, arm32, arm64,...)
	Architecture string `json:"architecture,omitempty"`
	// The operating system of the image (linux, windows, darwin,...)
	Os string `json:"os,omitempty"`
}

// ImageBase contains common fields between Image and Invocation image
type ImageBase struct {
	// A resolvable reference to the image. This may be interpreted differently based on imageType, but the default is to treat this as an OCI image
	Image string `json:"image"`
	// The type of image. If this is not specified, 'oci' is assumed
	ImageType string `json:"imageType,omitempty"`
	// A cryptographic hash digest that can be used to validate the image. This may be interpreted differently based on imageType
	Digest string `json:"digest,omitempty"`
	// The image size in bytes
	Size int `json:"size,omitempty"`
	// The target platform
	Platform *Platform `json:"platform,omitempty"`
	// The media type of the image
	MediaType string `json:"mediaType,omitempty"`
}

// InvocationImage is a bootstrapping image for the CNAB bundle.
type InvocationImage struct {
	ImageBase
}

// LocationReference is a reference to the image within the invocation image dile system
type LocationReference struct {
	// The path in the CNAB bundle to a file that references this image. It will be calculated from the root of the container
	Path string `json:"path,omitempty"`

	// The field to be replaced in the file specified by the 'path' property
	Field string `json:"field,omitempty"`

	// The MIME type of the file, used for determining how to process it
	MediaType string `json:"mediaType,omitempty"`
}

// Image is an application image for this CNAB bundle
type Image struct {
	ImageBase
	// The locations in the invocation image that reference this image. Used for rewriting
	Refs []LocationReference `json:"refs,omitempty"`

	// A description of the purpose of this image
	Description string `json:"description,omitempty"`
} // struct image

// Credential defines a particular credential, and where it should be placed in the invocation image
type Credential struct {
	// The path inside of the invocation image where credentials will be mounted
	Path string `json:"path,omitempty"`

	// The environment variable name, such as MY_VALUE, into which the credential will be placed
	Env string `json:"env,omitempty"`

	// A user-friendly description of this credential
	Description string `json:"description,omitempty"`
} // struct credential

// Parameter that can be passed into the invocation image
type Parameter struct {

	// Minimum integer value (ignored for non-integer parameters)
	MinValue int `json:"minValue,omitempty"`

	// Maximum integer value (ignored for non-integer parameters)
	MaxValue int `json:"maxValue,omitempty"`

	// Maximum string length (ignored for non-string parameters)
	MaxLength int `json:"maxLength,omitempty"`

	// Extra data about the parameter
	Metadata struct {

		// Description of this parameter
		Description string `json:"description,omitempty"`
	} `json:"metadata,omitempty"`

	// Minimum string length (ignored for non-string parameters)
	MinLength int `json:"minLength,omitempty"`

	Destination struct {

		// The path inside of the invocation image where parameter data is mounted
		Path string `json:"path,omitempty"`

		// The environment variable name, such as MY_VALUE, in which the parameter value is stored
		Env string `json:"env,omitempty"`

		// A user-friendly description of this parameter
		Description string `json:"description,omitempty"`
	} `json:"destination,omitempty"`

	// The data type of the parameter
	Type ParameterType `json:"type"`

	// If true, this parameter must be supplied
	Required bool `json:"required,omitempty"`

	// The default value of this parameter
	RawDefaultValue json.RawMessage `json:"defaultValue,omitempty"`

	// An optional exhaustive list of allowed values
	RawAllowedValues json.RawMessage `json:"allowedValues,omitempty"`
} // struct parameter

// DefaultValueBool returns the default value, assuming it is a boolean
func (p *Parameter) DefaultValueBool() (*bool, error) {
	if p.Type != ParameterTypeBoolean {
		return nil, fmt.Errorf(`parameter type is %q, not "boolean"`, p.Type)
	}
	switch {
	case len(p.RawDefaultValue) == 0, strings.TrimSpace(string(p.RawDefaultValue)) == "null":
		return nil, nil
	default:
		var value bool
		if err := json.Unmarshal(p.RawDefaultValue, &value); err != nil {
			return nil, err
		}
		return &value, nil
	}
}

// AllowedValuesBool returns the allowed values, assuming it is a boolean
func (p *Parameter) AllowedValuesBool() ([]bool, error) {
	if p.Type != ParameterTypeBoolean {
		return nil, fmt.Errorf(`parameter type is %q, not "boolean"`, p.Type)
	}
	switch {
	case len(p.RawAllowedValues) == 0, strings.TrimSpace(string(p.RawAllowedValues)) == "null":
		return nil, nil
	default:
		var value []bool
		if err := json.Unmarshal(p.RawAllowedValues, &value); err != nil {
			return nil, err
		}
		return value, nil
	}
}

// DefaultValueString returns the default value, assuming it is a string
func (p *Parameter) DefaultValueString() (*string, error) {
	if p.Type != ParameterTypeString {
		return nil, fmt.Errorf(`parameter type is %q, not "string"`, p.Type)
	}
	switch {
	case len(p.RawDefaultValue) == 0, strings.TrimSpace(string(p.RawDefaultValue)) == "null":
		return nil, nil
	default:
		var value string
		if err := json.Unmarshal(p.RawDefaultValue, &value); err != nil {
			return nil, err
		}
		return &value, nil
	}
}

// AllowedValuesString returns the allowed values, assuming it is a string
func (p *Parameter) AllowedValuesString() ([]string, error) {
	if p.Type != ParameterTypeString {
		return nil, fmt.Errorf(`parameter type is %q, not "string"`, p.Type)
	}
	switch {
	case len(p.RawAllowedValues) == 0, strings.TrimSpace(string(p.RawAllowedValues)) == "null":
		return nil, nil
	default:
		var value []string
		if err := json.Unmarshal(p.RawAllowedValues, &value); err != nil {
			return nil, err
		}
		return value, nil
	}
}

// DefaultValueInt returns the default value, assuming it is an integer
func (p *Parameter) DefaultValueInt() (*int, error) {
	if p.Type != ParameterTypeInteger {
		return nil, fmt.Errorf(`parameter type is %q, not "int"`, p.Type)
	}
	switch {
	case len(p.RawDefaultValue) == 0, strings.TrimSpace(string(p.RawDefaultValue)) == "null":
		return nil, nil
	default:
		var value int
		if err := json.Unmarshal(p.RawDefaultValue, &value); err != nil {
			return nil, err
		}
		return &value, nil
	}
}

// AllowedValuesInt returns the allowed values, assuming it is an integer
func (p *Parameter) AllowedValuesInt() ([]int, error) {
	if p.Type != ParameterTypeInteger {
		return nil, fmt.Errorf(`parameter type is %q, not "int"`, p.Type)
	}
	switch {
	case len(p.RawAllowedValues) == 0, strings.TrimSpace(string(p.RawAllowedValues)) == "null":
		return nil, nil
	default:
		var value []int
		if err := json.Unmarshal(p.RawAllowedValues, &value); err != nil {
			return nil, err
		}
		return value, nil
	}
}

// DefaultValue returns the default value of the parameter (nil if no default pvalue is provided)
func (p *Parameter) DefaultValue() (interface{}, error) {
	switch p.Type {
	case ParameterTypeString:
		return p.DefaultValueString()
	case ParameterTypeBoolean:
		return p.DefaultValueBool()
	case ParameterTypeInteger:
		return p.DefaultValueInt()
	}
	return nil, fmt.Errorf("unsupported parameter type %q", p.Type)
}

// AllowedValues returns a slice of allowed values
func (p *Parameter) AllowedValues() ([]interface{}, error) {
	switch p.Type {
	case ParameterTypeString:
		values, err := p.AllowedValuesString()
		if err != nil {
			return nil, err
		}
		result := make([]interface{}, len(values))
		for ix, v := range values {
			result[ix] = v
		}
		return result, nil
	case ParameterTypeBoolean:
		values, err := p.AllowedValuesBool()
		if err != nil {
			return nil, err
		}
		result := make([]interface{}, len(values))
		for ix, v := range values {
			result[ix] = v
		}
		return result, nil
	case ParameterTypeInteger:
		values, err := p.AllowedValuesInt()
		if err != nil {
			return nil, err
		}
		result := make([]interface{}, len(values))
		for ix, v := range values {
			result[ix] = v
		}
		return result, nil
	}
	return nil, fmt.Errorf("unsupported parameter type %q", p.Type)
}

// ParameterType is the expected type of a parameter's value
type ParameterType string

const (
	// ParameterTypeString is a string
	ParameterTypeString ParameterType = "string"
	// ParameterTypeInteger is an integer
	ParameterTypeInteger ParameterType = "int"
	// ParameterTypeBoolean is a boolean
	ParameterTypeBoolean ParameterType = "boolean"
)
