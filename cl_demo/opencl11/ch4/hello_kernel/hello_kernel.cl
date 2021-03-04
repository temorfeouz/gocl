#pragma OPENCL EXTENSION cl_khr_int64_base_atomics : enable

__kernel void hello_kernel(__global char16 *msg, __global int * counter) {

 constant char16 *p = "adimaslkdmaslkdmaslkdmaskdmalksdmlakmsdlkasmd";
 *msg=*p;

  //printf("%s\n", msg);

     // atomic_add(counter,1);
*counter=get_local_id(0);
   //local int sz = atomic_add(&(size),1);
   // if (sz >= capacity)
   //     return;
//
   // unsigned int i = get_global_id(0);
   // a[sz] = list[i];
}

