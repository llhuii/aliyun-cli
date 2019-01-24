package meta

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestApi_GetMethod(t *testing.T) {
	api := &Api{
		Method:  "post",
	}
	method := api.GetMethod()
	assert.Equal(t, method, "POST")

	api.Method = "get"
	method = api.GetMethod()
	assert.Equal(t, method, "GET")

	api.Method = "put"
	method = api.GetMethod()
	assert.Equal(t, method, "GET")
}

func TestApi_GetProtocol(t *testing.T) {
	api := &Api{
		Protocol:  "https://",
	}
	protocol := api.GetProtocol()
	assert.Equal(t, protocol, "https")

	api.Protocol = "http://"
	protocol = api.GetProtocol()
	assert.Equal(t, protocol, "http")
}

func TestApi_FindParameter(t *testing.T) {
	api := &Api{}
	api.Parameters = []Parameter{
		{
			Name:  "paramter",
	    },
	}
	parameter := api.FindParameter("paramter")
	assert.Equal(t, parameter.Name, "paramter")

	api.Parameters = []Parameter{
		{
			SubParameters: []Parameter{
				{
					Name: "subparameters",
				},
			},
			Name: "paramter",
	    },
	}
	parameter = api.FindParameter("paramter_test")
	assert.Nil(t, parameter)

	parameter = api.FindParameter("paramter.1.10")
	assert.Nil(t, parameter)

	api.Parameters = []Parameter{
		{
			Type:  "RepeatList",
			Name: "paramter",
		},
	}
	parameter = api.FindParameter("paramter.1.1")
	assert.Equal(t, parameter.Name, "paramter")

	parameter = api.FindParameter("paramter_test")
	assert.Nil(t, parameter)

	api.Parameters = nil
	parameter =api.FindParameter("paramter.1.10")
	assert.Nil(t, parameter)
}

func TestApi_ForeachParameters(t *testing.T) {
	api := &Api{}
	f := func (s string, p *Parameter) {
		p.Name = s
	}
	api.Parameters = []Parameter{
		{
			Type:  "RepeatList",
			Name:  "paramter",
		},
	}
	api.ForeachParameters(f)
	assert.Equal(t, api.Parameters[0].Name, "paramter.1")

	api.Parameters[0].Type = ""
	api.ForeachParameters(f)
	assert.Equal(t, api.Parameters[0].Name, "paramter.1")

	api.Parameters = []Parameter{
		{
			SubParameters: []Parameter{
				{
					Name: "subparameters",
				},
			},
			Name: "paramter",
		},
	}
	api.ForeachParameters(f)
	assert.Equal(t, api.Parameters[0].SubParameters[0].Name, "paramter.1.subparameters")
}

func TestApi_GetDocumentLink(t *testing.T) {
	api := &Api{
		Name: "apitest",
		Product: &Product{
			Code: "code",
		},
	}
	link := api.GetDocumentLink()
	assert.Equal(t, link, "https://help.aliyun.com/api/code/apitest.html")
}

func TestApi_CheckRequiredParameters(t *testing.T) {
	api := &Api{
		Parameters: []Parameter{
			{
				Name:"api",
				Required: true,
				Type: "Repeat",
			},
		},
	}
	checker := func(string)bool {
		return false
	}
	err := api.CheckRequiredParameters(checker)
	assert.Contains(t, err.Error(), "required parameters not assigned")

	api.Parameters = nil
	err = api.CheckRequiredParameters(checker)
	assert.Nil(t, err)
}
func TestParameterSliceLenAndSwap(t *testing.T) {
	parameters := ParameterSlice{
		{
			Name:"api",
			Required: true,
			Type: "RepeatList",
		},
	}
	len := parameters.Len()
	assert.Equal(t, len, 1)

	parameters = append(parameters, Parameter{
		Name:"api_test",
		Required: true,
		Type: "RepeatList",
	})
	parameters.Swap(0, 1)
	assert.Equal(t, parameters[0].Name, "api_test")
	assert.Equal(t, parameters[1].Name, "api")
}

func TestParameterSlice_Less(t *testing.T) {
	parameters := ParameterSlice{
		{
			Name:     "api",
			Required: true,
			Type:     "RepeatList",
		},
		{
			Name:"api_test",
			Required: true,
			Type: "RepeatList",
		},
	}
	ok := parameters.Less(0,1)
	assert.True(t, ok)

	parameters[1].Required = false
	ok = parameters.Less(0,1)
	assert.True(t, ok)

	parameters[0].Required = false
	ok = parameters.Less(0,1)
	assert.True(t, ok)

	parameters[1].Required = true
	ok = parameters.Less(0,1)
	assert.False(t, ok)
}