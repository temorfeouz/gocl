constant sampler_t sampler = CLK_NORMALIZED_COORDS_FALSE
   | CLK_ADDRESS_CLAMP | CLK_FILTER_NEAREST;

constant float SCALE = 3;

__kernel void interp(
read_only image2d_t src_imageR,
read_only image2d_t src_imageG,
read_only image2d_t src_imageB,
read_only image2d_t src_imageA,

write_only image2d_t dst_imageR,
write_only image2d_t dst_imageG,
write_only image2d_t dst_imageB,
write_only image2d_t dst_imageA) {

   float4 pixel;

   /* Determine input coordinate */
   float2 input_coord = (float2)
      (get_global_id(0) + (1.0f/(SCALE*2)),
       get_global_id(1) + (1.0f/(SCALE*2)));

   /* Determine output coordinate */
   int2 output_coord = (int2)
      (SCALE*get_global_id(0),
       SCALE*get_global_id(1));

   /* Compute interpolation */
   for(int i=0; i<SCALE; i++) {
      for(int j=0; j<SCALE; j++) {
         pixel = read_imagef(src_imageR, sampler,
           (float2)(input_coord +
           (float2)(1.0f*i/SCALE, 1.0f*j/SCALE)));

         write_imagef(dst_imageR, output_coord +
                      (int2)(i, j), pixel);
      }
   }

   for(int i=0; i<SCALE; i++) {
      for(int j=0; j<SCALE; j++) {
         pixel = read_imagef(src_imageG, sampler,
           (float2)(input_coord +
           (float2)(1.0f*i/SCALE, 1.0f*j/SCALE)));

         write_imagef(dst_imageG, output_coord +
                      (int2)(i, j), pixel);
      }
   }

   for(int i=0; i<SCALE; i++) {
      for(int j=0; j<SCALE; j++) {
         pixel = read_imagef(src_imageB, sampler,
           (float2)(input_coord +
           (float2)(1.0f*i/SCALE, 1.0f*j/SCALE)));

         write_imagef(dst_imageB, output_coord +
                      (int2)(i, j), pixel);
      }
   }

   for(int i=0; i<SCALE; i++) {
      for(int j=0; j<SCALE; j++) {
         pixel = read_imagef(src_imageA, sampler,
           (float2)(input_coord +
           (float2)(1.0f*i/SCALE, 1.0f*j/SCALE)));

         write_imagef(dst_imageA, output_coord +
                      (int2)(i, j), pixel);
      }
   }
}
