// +build cl20

package ocl

import (
	"github.com/hmwill/gocl/cl"
	"fmt"
	"unsafe"
)

type context20 interface {
	GetID() cl.CL_context
	GetInfo(param_name cl.CL_context_info) (interface{}, error)
	Retain() error
	Release() error

	CreateBuffer(flags cl.CL_mem_flags,
		size cl.CL_size_t,
		host_ptr unsafe.Pointer) (Buffer, error)
	CreateEvent() (Event, error)
	CreateProgramWithSource(count cl.CL_uint,
		strings [][]byte,
		lengths []cl.CL_size_t) (Program, error)
	CreateProgramWithBinary(devices []Device,
		lengths []cl.CL_size_t,
		binaries [][]byte,
		binary_status []cl.CL_int) (Program, error)

	GetSupportedImageFormats(flags cl.CL_mem_flags,
		image_type cl.CL_mem_object_type) ([]cl.CL_image_format, error)

	//cl20
	CreateCommandQueueWithProperties(device Device,
		properties []cl.CL_command_queue_properties) (CommandQueue, error)
	CreateSamplerWithProperties(properties []cl.CL_sampler_properties) (Sampler, error)
}

func (this *context) CreateCommandQueueWithProperties(device Device,
	properties []cl.CL_command_queue_properties) (CommandQueue, error) {
	var errCode cl.CL_int

	if command_queue_id := cl.CLCreateCommandQueueWithProperties(this.context_id, device.GetID(), properties, &errCode); errCode != cl.CL_SUCCESS {
		return nil, fmt.Errorf("CreateCommandQueueWithProperties failure with errcode_ret %d: %s", errCode, cl.ERROR_CODES_STRINGS[-errCode])
	} else {
		return &command_queue{command_queue_id}, nil
	}
}

func (this *context) CreateSamplerWithProperties(properties []cl.CL_sampler_properties) (Sampler, error) {
	var errCode cl.CL_int

	if sampler_id := cl.CLCreateSamplerWithProperties(this.context_id, properties, &errCode); errCode != cl.CL_SUCCESS {
		return nil, fmt.Errorf("CreateSamplerWithProperties failure with errcode_ret %d: %s", errCode, cl.ERROR_CODES_STRINGS[-errCode])
	} else {
		return &sampler{sampler_id}, nil
	}

}
