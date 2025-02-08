I used this assignment to try out echo instead of gorilla mux which what I am used to. 

There are a few coding decisions made that are a bit odd that were done to test out middleware and binding in echo. I have comments on those parts.


My tests also are all over the place in the way they are written. This was done on purpose because I have noticed various people write test in different ways.
The different ways are
1. new test for each part
2. using t.Run to have one test split up the parts
3. table tests

I personally prefer either 1 or 2 depending on the functionality being tested. I am not a fan of tables tests as I find them hard to read and maintain. 
