curl -isSX POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.web1.com","alias":"000000000000"}' > test18.out 2>&1
curl -isSX POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.web1.com","alias":"000000000001"}' > test18.out 2>&1
curl -isSX POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.web1.com","alias":"000000000002"}' > test18.out 2>&1
diff test18.out test18.ref