package grab

import (
  "unsafe"
  "syscall"
)
type ( // <linux/fb.h>
/*      bitfield struct {
         offset, length,
              msb_right uint32
                        } */
  var_screeninfo struct { // see /usr/include/linux/vt.h
             xres, yres,
           xres_virtual,
           yres_virtual,
       xoffset, yoffset,
         bits_per_pixel,
              grayscale uint32
/*     red, green, blue,
                 transp bitfield  // 12 == 4 * 3
       nonstd, activate,
          height, width,
  accel_flags, pixclock,
            left_margin,
           right_margin,
           upper_margin,
           lower_margin,
   hsync_len, vsync_len,
    sync, vmode, rotate uint32    // 15
               reserved [5]uint32 //  5
*/                                // 32
                  dummy [32]uint32
                        }
  fix_screeninfo struct {
                     id [16]byte
             smem_start,
               smem_len uint32
/*       type, type_aux,
          visual uint32           // 3
     xpanstep, ypanstep,
              ywrapstep int    // 1.5
            line_length uint32    // 1
             mmio_start uint      // 1
        mmio_len, accel uint32    // 2
               reserved [3]int // 1.5
*/                                // 10
                  dummy [10]uint32
                        }
)

func Geometry(fd uintptr) (x, y, b int) {
//
  const FBIOGET_VSCREENINFO = 0x4600
  var v var_screeninfo
  syscall.Syscall (syscall.SYS_IOCTL, fd, FBIOGET_VSCREENINFO, uintptr(unsafe.Pointer(&v)))
  return int(v.xres), int(v.yres), int(v.bits_per_pixel)/8
}
