package main

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/temorfeouz/gocl/cl_demo/utils"

	"github.com/temorfeouz/gocl/cl"
)

const INPUT_FILE = "input.png"
const OUTPUT_FILE = "output.png"
const PROGRAM_FILE = "interp.cl"

var KERNEL_FUNC = []byte("interp")

func main() {

	/* Host/device data structures */
	var device []cl.CL_device_id
	var context cl.CL_context
	var queue cl.CL_command_queue
	var program *cl.CL_program
	var kernel cl.CL_kernel
	var err cl.CL_int

	var err1 error
	var global_size [2]cl.CL_size_t

	/* Image data */
	var input_pixelsR, input_pixelsG, input_pixelsB, input_pixelsA, output_pixelsR, output_pixelsG, output_pixelsB, output_pixelsA []uint16
	var png_format cl.CL_image_format
	var input_imageR, input_imageG, input_imageB, input_imageA, output_imageR, output_imageG, output_imageB, output_imageA cl.CL_mem
	var origin, region [3]cl.CL_size_t
	var width, height cl.CL_size_t

	/* Open input file and read image data */
	input_pixelsR, input_pixelsG, input_pixelsB, input_pixelsA, width, height, err1 = utils.Read_image_data(INPUT_FILE)
	if err1 != nil {
		fmt.Printf("err on read image %s", err1)
		return
	} else {
		fmt.Printf("width=%d, height=%d\n", width, height)
	}
	output_pixelsR = make([]uint16, width*height)
	output_pixelsG = make([]uint16, width*height)
	output_pixelsB = make([]uint16, width*height)
	output_pixelsA = make([]uint16, width*height)

	/* Create a device and context */
	device = utils.Create_device()
	context = cl.CLCreateContext(nil, 1, device[:], nil, nil, &err)
	if err < 0 {
		println("Couldn't create a context")
		return
	}

	/* Build the program and create a kernel */
	program = utils.Build_program(context, device[:], PROGRAM_FILE, nil)
	kernel = cl.CLCreateKernel(*program, KERNEL_FUNC, &err)
	if err < 0 {
		fmt.Printf("Couldn't create a kernel: %d", err)
		return
	}

	/* Create input image object */
	png_format.Image_channel_order = cl.CL_LUMINANCE
	png_format.Image_channel_data_type = cl.CL_UNORM_INT16
	input_imageR = cl.CLCreateImage2D(context,
		cl.CL_MEM_READ_ONLY|cl.CL_MEM_COPY_HOST_PTR,
		&png_format, width, height, 0, unsafe.Pointer(&input_pixelsR[0]), &err)
	if err < 0 {
		println("Couldn't create the image object")
		return
	}
	input_imageG = cl.CLCreateImage2D(context,
		cl.CL_MEM_READ_ONLY|cl.CL_MEM_COPY_HOST_PTR,
		&png_format, width, height, 0, unsafe.Pointer(&input_pixelsG[0]), &err)
	if err < 0 {
		println("Couldn't create the image object")
		return
	}
	input_imageB = cl.CLCreateImage2D(context,
		cl.CL_MEM_READ_ONLY|cl.CL_MEM_COPY_HOST_PTR,
		&png_format, width, height, 0, unsafe.Pointer(&input_pixelsB[0]), &err)
	if err < 0 {
		println("Couldn't create the image object")
		return
	}
	input_imageA = cl.CLCreateImage2D(context,
		cl.CL_MEM_READ_ONLY|cl.CL_MEM_COPY_HOST_PTR,
		&png_format, width, height, 0, unsafe.Pointer(&input_pixelsA[0]), &err)
	if err < 0 {
		println("Couldn't create the image object")
		return
	}

	/* Create output image object */
	output_imageR = cl.CLCreateImage2D(context,
		cl.CL_MEM_WRITE_ONLY, &png_format, width,
		height, 0, nil, &err)
	if err < 0 {
		println("Couldn't create the image object")
		return
	}
	output_imageG = cl.CLCreateImage2D(context,
		cl.CL_MEM_WRITE_ONLY, &png_format, width,
		height, 0, nil, &err)
	if err < 0 {
		println("Couldn't create the image object")
		return
	}
	output_imageB = cl.CLCreateImage2D(context,
		cl.CL_MEM_WRITE_ONLY, &png_format, width,
		height, 0, nil, &err)
	if err < 0 {
		println("Couldn't create the image object")
		return
	}
	output_imageA = cl.CLCreateImage2D(context,
		cl.CL_MEM_WRITE_ONLY, &png_format, width,
		height, 0, nil, &err)
	if err < 0 {
		println("Couldn't create the image object")
		return
	}

	/* Create kernel arguments */
	err = cl.CLSetKernelArg(kernel, 0, cl.CL_size_t(unsafe.Sizeof(input_imageR)), unsafe.Pointer(&input_imageR))
	err = cl.CLSetKernelArg(kernel, 1, cl.CL_size_t(unsafe.Sizeof(input_imageG)), unsafe.Pointer(&input_imageG))
	err = cl.CLSetKernelArg(kernel, 2, cl.CL_size_t(unsafe.Sizeof(input_imageB)), unsafe.Pointer(&input_imageB))
	err = cl.CLSetKernelArg(kernel, 3, cl.CL_size_t(unsafe.Sizeof(input_imageA)), unsafe.Pointer(&input_imageA))
	err |= cl.CLSetKernelArg(kernel, 4, cl.CL_size_t(unsafe.Sizeof(output_imageR)), unsafe.Pointer(&output_imageR))
	err |= cl.CLSetKernelArg(kernel, 5, cl.CL_size_t(unsafe.Sizeof(output_imageG)), unsafe.Pointer(&output_imageG))
	err |= cl.CLSetKernelArg(kernel, 6, cl.CL_size_t(unsafe.Sizeof(output_imageB)), unsafe.Pointer(&output_imageB))
	err |= cl.CLSetKernelArg(kernel, 7, cl.CL_size_t(unsafe.Sizeof(output_imageA)), unsafe.Pointer(&output_imageA))
	if err < 0 {
		println("Couldn't set a kernel argument")
		return
	}
	ts := time.Now()
	/* Create a command queue */
	queue = cl.CLCreateCommandQueue(context, device[0], 0, &err)
	if err < 0 {
		println("Couldn't create a command queue")
		return
	}

	/* Enqueue kernel */
	global_size[0] = width
	global_size[1] = height
	err = cl.CLEnqueueNDRangeKernel(queue, kernel, 2, nil, global_size[:],
		nil, 0, nil, nil)
	if err < 0 {
		println("Couldn't enqueue the kernel")
		return
	}

	/* Read the image object */
	origin[0] = 0
	origin[1] = 0
	origin[2] = 0
	region[0] = width
	region[1] = height
	region[2] = 1

	err = cl.CLEnqueueReadImage(queue, output_imageR, cl.CL_TRUE, origin,
		region, 0, 0, unsafe.Pointer(&output_pixelsR[0]), 0, nil, nil)
	if err < 0 {
		println("Couldn't read from the image object")
		return
	}
	err = cl.CLEnqueueReadImage(queue, output_imageG, cl.CL_TRUE, origin,
		region, 0, 0, unsafe.Pointer(&output_pixelsG[0]), 0, nil, nil)
	if err < 0 {
		println("Couldn't read from the image object")
		return
	}
	err = cl.CLEnqueueReadImage(queue, output_imageB, cl.CL_TRUE, origin,
		region, 0, 0, unsafe.Pointer(&output_pixelsB[0]), 0, nil, nil)
	if err < 0 {
		println("Couldn't read from the image object")
		return
	}
	err = cl.CLEnqueueReadImage(queue, output_imageA, cl.CL_TRUE, origin,
		region, 0, 0, unsafe.Pointer(&output_pixelsA[0]), 0, nil, nil)
	if err < 0 {
		println("Couldn't read from the image object")
		return
	}
	fmt.Printf("upload, work and read result took %s\n", time.Now().Sub(ts))
	ts = time.Now()
	/* Create output PNG file and write data */
	utils.Write_image_data(OUTPUT_FILE, output_pixelsR, output_pixelsG, output_pixelsB, output_pixelsA, width, height)
	fmt.Printf("saving image took %s\n", time.Now().Sub(ts))
	ts = time.Now()
	/* Deallocate resources */
	cl.CLReleaseMemObject(input_imageR)
	cl.CLReleaseMemObject(output_imageR)
	cl.CLReleaseMemObject(input_imageG)
	cl.CLReleaseMemObject(output_imageG)
	cl.CLReleaseMemObject(input_imageB)
	cl.CLReleaseMemObject(output_imageB)
	cl.CLReleaseMemObject(input_imageA)
	cl.CLReleaseMemObject(output_imageA)
	cl.CLReleaseKernel(kernel)
	cl.CLReleaseCommandQueue(queue)
	cl.CLReleaseProgram(*program)
	cl.CLReleaseContext(context)
	fmt.Printf("release took %s\n", time.Now().Sub(ts))
}
