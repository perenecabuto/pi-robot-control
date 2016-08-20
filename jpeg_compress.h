#ifndef JPEG_COMPRESS_H
#define JPEG_COMPRESS_H

//int compress_image_to_jpeg(char *vd, unsigned char *buffer, int size, int quality);
void compressYUYVtoJPEG(const char* input, const char* output, const int width, const int height);

#endif
