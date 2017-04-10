# Server examples

tsdbb bench write graphite tcp://127.0.0.1:2003
tsdbb bench read graphite-web tcp://127.0.0.1:8080
tsdbb bench read carbon-zipper tcp://127.0.0.1:8080


tsdbb historical write 

# Cli examples

tsdbb add 1000
tsdbb upto 10000 --step 1000

tsdbb status
