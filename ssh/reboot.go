package ssh

import "time"

func (c *Client) PossibleRebootWait() {
	c.logger.InfoContext(c.ctx, "Possible reboot, waiting... (Press Ctrl+C to abort)")
	time.Sleep(time.Second)

	timeout := c.Options.Timeout
	c.Options.Timeout = 5 * time.Second
	var retest bool
	for {
		c.logger.DebugContext(c.ctx, "Trying to connect...")
		// todo check error code for connection refused or timeout
		if _, err := c.Connect(); err == nil {
			if retest {
				break
			}
			// let's wait a bit more to be sure
			time.Sleep(time.Second)
			c.Close()
			retest = true
			continue
		}
	}
	c.Close()
	c.logger.InfoContext(c.ctx, "Connected")
	c.Options.Timeout = timeout
}
