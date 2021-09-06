resource "alertlogic_user" "user" {
  name         = "Bob Loblaw"
  email        = "bob@bobloblawlaw.com"
  mobile_phone = "555-123-4567"
  active       = true
  role_ids     = ["F578CCE5-9574-4489-BF05-A04075838DE3"]
}
