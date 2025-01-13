package cloudinitgo

import "fmt"

// This function returns a slice of strings containing cloud-init configuration arguments for QEMU
// It creates a single argument string in the format:
// "ds=nocloud-net;s=http://localhost:<port>/"
// where:
// - ds=nocloud-net: Specifies to use NoCloud network data source
// - s=http://localhost:<port>/: Specifies the metadata server URL running locally
// https://niekdeschipper.com/projects/cloud-init.html
func (c *cloudInit) QemuConfig() []string {

	cloudInitArgs := fmt.Sprintf("ds=nocloud-net;s=http://localhost:%d/", c.Port)

	return []string{cloudInitArgs}
}
