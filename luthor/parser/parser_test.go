package parser

import (
	"testing"

	"github.com/getnoops/terraform-gen/luthor/ast"
)

const stringStr = `string`
const complexStr = `
object({
	allowed_methods = list(string)
	cached_methods  = list(string)
	cache_policy = object({
		cookie_behavior       = string
		cookie_items          = optional(list(string))
		header_behavior       = string
		header_items          = optional(list(string))
		query_string_behavior = string
		query_string_items    = optional(list(string))
	})
	compress                  = optional(bool)
	default_ttl               = optional(number)
	field_level_encryption_id = optional(string)
	lambda_function_association = optional(list(object({
		event_type   = string
		lambda_arn   = string
		include_body = optional(bool)
	})))
	function_association = optional(list(object({
		event_type   = string
		function_arn = string
	})))
	max_ttl                    = optional(number)
	min_ttl                    = optional(number)
	origin_request_policy_id   = optional(string)
	realtime_log_config_arn    = optional(string)
	response_headers_policy_id = optional(string)
	smooth_streaming           = optional(bool)
	target_origin_id           = string
	trusted_key_groups         = optional(list(string))
	trusted_signers            = optional(list(string))
	viewer_protocol_policy     = string
})
`

const complexStr2 = `
list(object({
	id         = string
	put_events = bool
}))
`

func Test_Kitchen(t *testing.T) {
	f, err := ParseType(&ast.Source{
		Name:  "test",
		Input: complexStr2,
	})
	if err != nil {
		t.Fatal(err)
		return
	}

	t.Log(f.Name)
}
