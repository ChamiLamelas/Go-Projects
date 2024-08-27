curl -isSX POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.google.com", "alias":"custom"}' > test10.out 2>&1
curl -isSX POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.nytimes.com", "alias":"custom"}' >> test10.out 2>&1
diff test10.out test10.ref