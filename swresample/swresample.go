package swresample

/*
	#cgo pkg-config: libswresample
	#include <libswresample/swresample.h>
*/
import "C"

import (
	"errors"
	"unsafe"

	"github.com/wang1219/go-libav/avutil"
)

var (
	ErrAllocationError      = errors.New("allocation error")
)

type Context struct {
	CSwrContext *C.SwrContext
}

func NewSwrContextWithOpts(ocl avutil.ChannelLayout, osf avutil.SampleFormat, osr int, icl avutil.ChannelLayout, isf avutil.SampleFormat, isr, lo, lc int) (*Context, error) {
	cSwrContext := C.swr_alloc_set_opts((*C.SwrContext)(nil), C.int64_t(ocl), (C.enum_AVSampleFormat)(osf), C.int(osr), C.int64_t(icl), (C.enum_AVSampleFormat)(isf), C.int(isr), C.int(lo), unsafe.Pointer(&lc))
	if cSwrContext == nil {
		return nil, ErrAllocationError
	}
	return NewSwrContextFromC(unsafe.Pointer(cSwrContext)), nil
}

func NewSwrContextFromC(cSwrContext unsafe.Pointer) *Context {
	return &Context{CSwrContext: (*C.SwrContext)(cSwrContext)}
}

//Initialize context after user parameters have been set.
func (s *Context) Init() error {
	code := C.swr_init(s.CSwrContext)
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

//Check whether an swr context has been initialized or not.
func (s *Context) IsInitialized() bool {
	code := C.swr_is_initialized(s.CSwrContext)
	return code == 0
}

//Context destructor functions. Free the given Context and set the pointer to NULL.
func (s *Context) Free() {
	C.swr_free((**C.SwrContext)(&s.CSwrContext))
}

//Closes the context so that swr_is_initialized() returns 0.
func (s *Context) Close() {
	C.swr_close(s.CSwrContext)
}

//Core conversion functions. Convert audio
func (s *Context) Convert(out **uint8, oc int, in **uint8, ic int) error {
	code := C.swr_convert(s.CSwrContext, (**C.uint8_t)(unsafe.Pointer(out)), C.int(oc), (**C.uint8_t)(unsafe.Pointer(in)), C.int(ic))
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

//Convert the next timestamp from input to output timestamps are in 1/(in_sample_rate * out_sample_rate) units.
func (s *Context) NextPts(pts int64) int64 {
	return int64(C.swr_next_pts(s.CSwrContext, C.int64_t(pts)))
}

//Low-level option setting functions
//These functons provide a means to set low-level options that is not possible with the AvOption API.
//Activate resampling compensation ("soft" compensation).
func (s *Context) SetCompensation(sd, cd int) error {
	code := C.swr_set_compensation(s.CSwrContext, C.int(sd), C.int(cd))
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

//Set a customized input channel mapping.
func (s *Context) SetChannelMapping(cm *int) error {
	code := C.swr_set_channel_mapping(s.CSwrContext, (*C.int)(unsafe.Pointer(cm)))
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

//Set a customized remix matrix.
func (s *Context) SetMatrix(m *int, t int) error {
	code := C.swr_set_matrix(s.CSwrContext, (*C.double)(unsafe.Pointer(m)), C.int(t))
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

//Sample handling functions. Drops the specified number of output samples.
func (s *Context) DropOutput(c int) error {
	code := C.swr_drop_output(s.CSwrContext, C.int(c))
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

//Injects the specified number of silence samples.
func (s *Context) InjectSilence(c int) int {
	return int(C.swr_inject_silence(s.CSwrContext, C.int(c)))
}

//Gets the delay the next input sample will experience relative to the next output sample.
func (s *Context) GetDelay(b int64) int64 {
	return int64(C.swr_get_delay(s.CSwrContext, C.int64_t(b)))
}

//Frame based API. Convert the samples in the input Frame and write them to the output Frame.
func (s *Context) ConvertFrame(o, i *avutil.Frame) error {
	code := C.swr_convert_frame(s.CSwrContext, (*C.AVFrame)(unsafe.Pointer(o.CAVFrame)), (*C.AVFrame)(unsafe.Pointer(i.CAVFrame)))
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}

//Configure or reconfigure the Context using the information provided by the AvFrames.
func (s *Context) ConfigFrame(o, i *avutil.Frame) error {
	code := C.swr_config_frame(s.CSwrContext, (*C.AVFrame)(unsafe.Pointer(o.CAVFrame)), (*C.AVFrame)(unsafe.Pointer(i.CAVFrame)))
	if code < 0 {
		return avutil.NewErrorFromCode(avutil.ErrorCode(code))
	}
	return nil
}