package webhook

import (
	"bytes"
	"testing"

	"github.com/lukasknuth/gotify-slack-webhook/gotify"
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
    	},
    	{
    		"type": "section",
    		"block_id": "section567",
    		"text": {
    			"type": "mrkdwn",
    			"text": "<https://example.com|Overlook Hotel> \n :star: \n Doors had too many axe holes, guest in room 237 was far too rowdy, whole place felt stuck in the 1920s."
    		},
    		"accessory": {
    			"type": "image",
    			"image_url": "https://is5-ssl.mzstatic.com/image/thumb/Purple3/v4/d3/72/5c/d3725c8f-c642-5d69-1904-aa36e4297885/source/256x256bb.jpg",
    			"alt_text": "Haunted hotel image"
    		}
    	},
    	{
    		"type": "section",
    		"block_id": "section789",
    		"text": {
    			"type": "plain_text"
    			"text": "Rating"
    		}
    		"fields": [
    			{
    				"type": "mrkdwn",
    				"text": "*1.0*"
    			},
    			{
    				"type": "plain_text",
    				"text": "Not a good time"
    			}
    		]
    	}
    ]
}
`

func TestParse(t *testing.T) {
	body, err := Parse(exampleBody)
	assert.Nil(t, err)
	assert.Equal(t, "Danny Torrence left a 1 star review for your property.", body.Text)
	assert.Len(t, body.Blocks, 3)

	buf := new(bytes.Buffer)
	err = body.Blocks[0].Render(gotify.Wrap(buf))
	assert.Nil(t, err)
	assert.Equal(t, "Danny Torrence left the following review for your property:\n", buf.String())

	buf = new(bytes.Buffer)
	err = body.Blocks[2].Render(gotify.Wrap(buf))
	assert.Nil(t, err)
	assert.Equal(t, "Rating\n*1.0*\nNot a good time\n", buf.String())
}
