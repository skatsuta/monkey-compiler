let fib = fn(x) {
   if (x == 0) {
     0;
   } else {
     if (x == 1) {
       1;
     } else {
       fib(x - 1) + fib(x - 2);
     }
   }
};
puts(fib(15));
