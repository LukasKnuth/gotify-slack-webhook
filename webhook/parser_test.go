package webhook

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const exampleBody = `{
    "text": "Danny Torrence left a 1 star review for your property.",
    "blocks": [
    	{
    		"type": "section",
    		"text": {
    			"type": "mrkdwn",
    			"text": "Danny Torrence left the following review for your property:"
    		}
    	}
    ]
}
`

func TestParse(t *testing.T) {
	body, err := Parse(exampleBody)
	assert.Nil(t, err)
	assert.Equal(t, "Danny Torrence left a 1 star review for your property.", body.Text)
	assert.Len(t, body.Blocks, 1)
}
