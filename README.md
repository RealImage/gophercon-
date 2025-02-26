# 🎬 Degrees of Separation - Moviebuff API (Go)
**Find the shortest connection between two actors or filmmakers using Moviebuff data.**  

---

## 📌 Overview  
This Go program determines the **degrees of separation** between two people (actors, directors, etc.) using **Moviebuff's API**.  
It builds a **graph of connections** using movies, cast, and crew details and finds the shortest path using **Breadth-First Search (BFS)**.

---

## 📂 Project Structure  
```
qube_assignment/
│── main.go                   # Entry point (CLI)
│── go.mod                    # Go module file
│── go.sum                    # Dependencies
│── models/
│   ├── models.go             # Data models (Movie, Person, Cast, Crew)
│── services/
│   ├── fetch.go              # Fetches & parses API data
│   ├── bfs.go                # BFS algorithm to find connections between two people
│── utils/
│   ├── utils.go              # API rate limiting
│── README.md                 # Documentation
```

---

## 🚀 Installation & Setup  

### 1️⃣ Clone the Repository  
```sh
git clone https://github.com/viswamvs/qube_assignment.git
cd qube_assignment
```

### 2️⃣ Initialize Go Modules  
```sh
go mod init qube_assignment
go mod tidy
```

### 3️⃣ Run the Program  
```sh
go run main.go amitabh-bachchan robert-de-niro
```
📌 **Example Output:**  
```
Degrees of Separation: 3

1. amitabh-bachchan → the-great-gatsby
2. the-great-gatsby → leonardo-dicaprio
3. leonardo-dicaprio → robert-de-niro
```

---

## 🛠 How It Works  
1️⃣ **Fetches Data from Moviebuff API**  
   - Retrieves **actor, director, and movie details** using JSON API.  
2️⃣ **Builds a Connection Graph**  
   - Uses **cast & crew relationships** to link people via movies.  
3️⃣ **Finds Shortest Path Using BFS**  
   - Ensures **minimum degree of separation** is found efficiently.  
4️⃣ **Applies Rate Limiting & Caching**  
   - Prevents API throttling and **reduces redundant requests**.  

---

## 📌 Features  
✅ **Supports Actors, Directors, Producers, Writers, etc.**  
✅ **Handles API Rate Limiting** (Prevents 403 errors)  
✅ **Uses BFS for Efficient Search**  
✅ **Caches API Results** (Faster Performance)  

---

## 📝 API Response Format  
Example response for a **movie** (`taxi-driver`):  
```json
{
  "url": "taxi-driver",
  "type": "Movie",
  "name": "Taxi Driver",
  "cast": [
    {"url": "robert-de-niro", "name": "Robert De Niro", "role": "Actor"},
    {"url": "martin-scorsese", "name": "Martin Scorsese", "role": "Supporting Actor"}
  ],
  "crew": [
    {"url": "martin-scorsese", "name": "Martin Scorsese", "role": "Director"}
  ]
}
```

Example response for a **person** (`amitabh-bachchan`):  
```json
{
  "url": "amitabh-bachchan",
  "type": "Person",
  "name": "Amitabh Bachchan",
  "movies": [
    {"url": "the-great-gatsby", "name": "The Great Gatsby"}
  ]
}
```
