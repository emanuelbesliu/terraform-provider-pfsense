resource "pfsense_firewall_schedule" "business_hours" {
  name        = "BusinessHours"
  description = "Weekdays 9am to 5pm"

  time_range = [
    {
      position          = "1,2,3,4,5"
      start_time        = "9:00"
      stop_time         = "17:00"
      range_description = "Monday through Friday"
    },
  ]
}

resource "pfsense_firewall_schedule" "holidays" {
  name        = "Holidays"
  description = "Specific calendar days"

  time_range = [
    {
      month             = "12"
      day               = "25"
      start_time        = "0:00"
      stop_time         = "23:59"
      range_description = "Christmas Day"
    },
    {
      month             = "1"
      day               = "1"
      start_time        = "0:00"
      stop_time         = "23:59"
      range_description = "New Year's Day"
    },
  ]
}
