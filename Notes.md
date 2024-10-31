# Key Considerations in Solving the Assessment

1. **Algorithm**  
   The solution uses the Breadth-First Search (BFS) algorithm to find the shortest degree of separation between two artists efficiently.

2. **Rate Limiting**  
   A custom HTTP client with an adjustable rate limiter is implemented to handle `http.StatusTooManyRequests` responses effectively.

3. **Optimizations**  
   - **Concurrency**: The solution leverages goroutines and channels for concurrent requests, significantly reducing search times.
   - **Caching**: A concurrency-safe `sync.Map` cache stores results from previously fetched requests, minimizing redundant API calls.
     - **Note**: A substantial number of requests returned `403 Forbidden`, which are also cached to avoid repeated requests to those URLs.

4. **Best Practices**  
   The `FetchEntityDetails()` function is designed as a generic utility, serving both "Person" and "Movie" requests for streamlined and reusable code.
