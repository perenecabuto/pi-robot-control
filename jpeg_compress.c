#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#include <string.h>
#include <linux/types.h>
#include <linux/videodev2.h>
#include <jpeglib.h>


/**
 * converts a YUYV raw buffer to a JPEG buffer.
 * input is in YUYV (YUV 422). output is JPEG binary.
 * from https://linuxtv.org/downloads/v4l-dvb-apis/V4L2-PIX-FMT-YUYV.html:
 *      Each four bytes is two pixels.
 *      Each four bytes is two Y's, a Cb and a Cr.
 *      Each Y goes to one of the pixels, and the Cb and Cr belong to both pixels.
 *
 * inspired by: http://stackoverflow.com/questions/17029136
 * /weird-image-while-trying-to-compress-yuv-image-to-jpeg-using-libjpeg
 */
void compressYUYVtoJPEG(unsigned char* input, unsigned char* output, const int width, const int height) {
	struct jpeg_compress_struct cinfo;
	struct jpeg_error_mgr jerr;
	JSAMPROW row_pointer[1];
	unsigned char *line_buffer, *yuyv;
	int z = 0;

	uint8_t* outbuffer = NULL;
	uint64_t outlen = 0;

	yuyv = input;

	cinfo.err = jpeg_std_error(&jerr);
	jpeg_create_compress(&cinfo);
	jpeg_mem_dest(&cinfo, &outbuffer, &outlen);

	// jrow is a libjpeg row of samples array of 1 row pointer
	cinfo.image_width = width & -1;
	cinfo.image_height = height & -1;
	cinfo.input_components = 3;
	cinfo.in_color_space = JCS_YCbCr; //libJPEG expects YUV 3bytes, 24bit
	line_buffer = calloc(cinfo.image_width * 3, 1);

	jpeg_set_defaults(&cinfo);
	jpeg_set_quality(&cinfo, 100, TRUE);
	jpeg_start_compress(&cinfo, TRUE);

	while(cinfo.next_scanline < cinfo.image_height) {
		int x;
		unsigned char *ptr = line_buffer;

		for (x = 0; x < cinfo.image_width; x++) {
			int r, g, b;
			int y, u, v;

			if (!z)
				y = yuyv[1] << 8;
			else
				y = yuyv[3] << 8;
			u = yuyv[0] - 128;
			v = yuyv[2] - 128;

			r = (y + (359 * v)) >> 8;
			g = (y - (88 * u) - (183 * v)) >> 8;
			b = (y + (454 * u)) >> 8;

			*(ptr++) = (r > 255) ? 255 : ((r < 0) ? 0 : r);
			*(ptr++) = (g > 255) ? 255 : ((g < 0) ? 0 : g);
			*(ptr++) = (b > 255) ? 255 : ((b < 0) ? 0 : b);

			if (z++) {
				z = 0;
				yuyv += 4;
			}
		}

		row_pointer[0] = line_buffer;
		jpeg_write_scanlines(&cinfo, row_pointer, 1);
	}

	jpeg_finish_compress(&cinfo);
	jpeg_destroy_compress(&cinfo);
	free(line_buffer);

	/*printf("libjpeg produced %d bytes %d\n", outlen, cinfo.image_width);*/
	strncpy(output, outbuffer, outlen);
}

