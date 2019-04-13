let fib = fn(x) {
   if (x == 0) {
     return 0;
   }
   if (x == 1) {
     return 1;
   }
   fib(x - 1) + fib(x - 2);
};

let N = 15;
puts(fib(N));
