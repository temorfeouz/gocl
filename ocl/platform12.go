// +build cl12 cl20

package ocl

import (
	"fmt"
	"github.com/hmwill/gocl/cl"
)

func (this *platform) UnloadCompiler() error {
	if errCode := cl.CLUnloadPlatformCompiler(this.platform_id); errCode != cl.CL_SUCCESS {
		return fmt.Errorf("UnloadCompiler failure with errcode_ret %d: %s", errCode, cl.ERROR_CODES_STRINGS[-errCode])
	}
	return nil
}
