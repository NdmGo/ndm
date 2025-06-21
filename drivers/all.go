package drivers

import (
	_ "ndm/drivers/ftp"
	_ "ndm/drivers/local"
	_ "ndm/drivers/s3"
	_ "ndm/drivers/sftp"
	_ "ndm/drivers/tpl"
)

// All do nothing,just for import
// same as _ import
func All() {

}
