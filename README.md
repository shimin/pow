## Word of Wisdom tcp server

Demonstrates client-server communication based on "Hashcash" POW algorythm.
1. Client connects to the server
2. Client starts auth process and gets challenge task to solve.
3. As soon as calculated the solution has to be sent to the server for quick validation.
4. Server validates the solution and sends random word of wisdom in positive case.  


### Protection from DOS attacks with POW challenge-response protocol

Client gets random set of bytes and calculates SHA-256 hash of this set.
Client has to find a correct hash of source bytes using bruteforce.
The solution must satisfy the simple rule: first N bits has to be 0.
Amount of bits (N) describes complexity of bruteforce calculations for the Client.
The probability of getting correct answer is between 1..2^N hash calc iterations.
The variable that sets complexity of algorythm (sets N bits) named 'TargetBits'

Server receives 8-bytes SHA key that can be quickly validated.
In case of successful validation Server sends to the Client random quote from the book of quotes of wisdom.

### How to run

```sh
docker-compose build
docker-compose up -d
docker ps
docker logs <container_name> -f
```

### Encryption algorythm
Hashcash is based on SHA256
SHA-256 is popular and presented in standart go package
simple validation

