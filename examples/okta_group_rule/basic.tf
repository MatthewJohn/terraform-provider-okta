resource "okta_group" "test" {
  name = "testAcc_%[1]d"
}

resource "okta_group_rule" "test" {
  name              = "testAcc_%[1]d"
  status            = "ACTIVE"
  group_assignments = ["${okta_group.test.id}"]
  expression_type   = "urn:okta:expression:1.0"
  expression_value  = "String.startsWith(user.articulateId,\"auth0|\")"
}
