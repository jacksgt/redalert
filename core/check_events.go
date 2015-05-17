package core

import "strings"

var MaxEventsStored = 100

func (c *Check) storeEvent(event *Event) {
	c.LastEvent = event
	c.EventHistory.PushFront(event)
	if c.EventHistory.Len() > MaxEventsStored {
		c.EventHistory.Remove(c.EventHistory.Back())
	}
}

func (c *Check) triggerAlerts(event *Event) {

	go func() {

		var err error
		for _, alert := range c.Alerts {

			if event.isRedAlert() {
				c.Log.Println(red, "Sending red alert via", alert.Name(), "for", c.Name, reset)
			} else if event.isGreenAlert() {
				c.Log.Println(green, "Sending green alert via", alert.Name(), "for", c.Name, reset)
			}

			err = alert.Notify(event)
			if err != nil {
				c.Log.Println(red, "CRITICAL: Failure triggering alert ["+alert.Name()+"]: ", err.Error())
			}
		}

	}()
}

func (c *Check) RecentMetrics(metric string) string {
	var output []string
	for e := c.EventHistory.Front(); e != nil; e = e.Next() {
		event := e.Value.(*Event)
		if event != nil {
			metricStr := event.DisplayMetric(metric)
			output = append([]string{metricStr}, output...)
		}
	}
	return strings.Join(output, ",")
}
