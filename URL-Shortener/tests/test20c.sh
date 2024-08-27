curl -s -w "\nResponse code: %{http_code}\n" -X GET http://localhost:8000/urlshortener/expand/0 >> test20.out 2>&1
curl -s -w "\nResponse code: %{http_code}\n" -X GET http://localhost:8000/urlshortener/expand/1 >> test20.out 2>&1
curl -s -w "\nResponse code: %{http_code}\n" -X POST http://localhost:8000/urlshortener/shorten -H "Content-Type: application/json" -d '{"url":"https://www.web3.com"}' >> test20.out 2>&1
curl -s -w "\nResponse code: %{http_code}\n" -X GET http://localhost:8000/urlshortener/analytics/0 >> test20.out 2>&1
curl -s -w "\nResponse code: %{http_code}\n" -X GET http://localhost:8000/urlshortener/analytics/1 >> test20.out 2>&1
curl -s -w "\nResponse code: %{http_code}\n" -X GET http://localhost:8000/urlshortener/analytics/2 >> test20.out 2>&1
diff test20.out test20.ref