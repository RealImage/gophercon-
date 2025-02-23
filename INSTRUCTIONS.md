# Challenge 2015

This solution finds the shortest degree of connection between two actors based on their movie collaborations using single source shortest path breadth-first algorithm.

## Prerequisites

- Go 1.23.5 or later

## Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/c-harish/toupdate.git
    cd challenge2015
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

## Usage

1. Run the program:
    ```sh
    go run main.go
    ```

2. Enter the source actor's moviebuff URL when prompted:
    ```
    source actor : <source_actor_url>
    ```

3. Enter the target actor's moviebuff URL when prompted:
    ```
    target actor : <target_actor_url>
    ```
    
Refer to https://www.moviebuff.com/ for valid actor_url

The program will output the degrees of separation and the list movies connecting the two actors.
