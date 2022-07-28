## Word of Wisdom tcp server

Demonstrates client-server communication based on "Hashcash" POW algorythm.
1. Client connects to the server
2. Client starts auth process and gets challenge task to solve.
3. As soon as calculated the solution has to be sent to the server for quick validation.
4. Server validates the solution and sends random word of wisdom in positive case.  


### Protection from DOS attacks with POW challenge-response protocol

Client gets random set of bytes and calculates SHA-256 hash of this set.\
Client has to find a correct hash of source bytes using bruteforce.\
The solution must satisfy the simple rule: first N bits has to be 0.\
Amount of bits (N) describes complexity of bruteforce calculations for the Client.\
The probability of getting correct answer is between 1..2^N hash calc iterations.\
The variable that sets complexity of algorythm (sets N bits) named 'TargetBits'\
\
According to benchmarks 20..24 is the most suitable values of target bits\
\
Server receives 8-bytes SHA key that can be quickly validated.\
In case of successful validation Server sends to the Client random quote from the book of quotes of wisdom.\

### How to run

```sh
docker-compose build
docker-compose up -d
docker ps
docker logs <container_name> -f
```

### Encryption algorythm
- Hashcash is based on SHA256
- SHA-256 is popular and presented in standart go package
- simple validation

### Benchmarks
cpu: Intel(R) Core(TM) i5-8250U CPU @ 1.60GHz

#### Calculate

| complexity |	  cycles    |     avg time     |    alloc bytes    |     alloc count      |
| -----------|--------------|------------------|-------------------|----------------------|
|	2          | 	 1000000	  |      1255 ns/op  |	     328 B/op    |	         7 allocs/op|
|	8          | 	   64663	  |     16080 ns/op  |	    1912 B/op    |	        40 allocs/op|
|	16         | 	     100	  |  10595463 ns/op  |	 1034778 B/op    |	     21558 allocs/op|
|	20         | 	       6	  | 183752350 ns/op  |	20055217 B/op    |	    417816 allocs/op|
|	24         | 	       1	  |9222631900 ns/op  |	917003464 B/op   |	  19104235 allocs/op|
|	25         | 	       1	  |19911452500 ns/op |	2047857208 B/op  |	  42663649 allocs/op|

#### Validation

| complexity |	  cycles    |     avg time     |    alloc bytes    |     alloc count      |
| -----------|--------------|------------------|-------------------|----------------------|
|   2        | 	 1000000	|      1003 ns/op  |     280 B/op	   |    6 allocs/op		  |
|   8        | 	 1000000	|      1139 ns/op  |     280 B/op	   |    6 allocs/op		  |
|   16       | 	 1201792	|      1013 ns/op  |     280 B/op	   |    6 allocs/op		  |
|   20       | 	 1267788	|      988.2 ns/op |     280 B/op	   |    6 allocs/op		  |
|   24       | 	 1363494	|      924.5 ns/op |     280 B/op	   |    6 allocs/op		  |
|   25       | 	 1260970	|      810.0 ns/op |     280 B/op	   |    6 allocs/op		  |
