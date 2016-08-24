#ifndef JPEG_COMPRESS_H
#define JPEG_COMPRESS_H

#include <stdint.h>

//int compress_image_to_jpeg(char *vd, unsigned char *buffer, int size, int quality);
void compressYUYVtoJPEG(const uint8_t* const input, uint8_t* output, uint32_t width, uint32_t height);

#endif
