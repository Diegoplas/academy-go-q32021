# PokeAPI

With this API you can check the first and second generation pokemons' characteristics such as pokedex id, name, height and types.

### Requirements

* Go 1.15

### Framework

This project utilizes Gorilla web toolkit.

### 1. Install dependencies

As this project utilizes go modules, the dependencies can be easily downloaded executing the following line:
```
go mod download
```

### 2. Usage

1. The main program should be excecuted from root for the paths to match correcly. To excecute it we can use:
   ```
   go run ./main
   ```

2. First Generation Endpoint - /first-generation/{id}

   This endpoint gets the information from a specific pokemon from a CSV file. It should be from the first generation (pokedex IDs 1 to 151.)
   ```
   Eg. http://localhost:8000/first-generation/15
   ```

3. Second Generation Endpoint - /second-generation/{id}

   This endpoint gets the information from a specific pokemon from https://pokeapi.co/. It should be from the second generation (pokedex IDs 152 to 251.)
   ```
   Eg. http://localhost:8000/second-generation/251
   ```

3. Worker Pool - /worker-pool

   This endpoint gets the information from a list of pokemon from a CSV file using go routines. Three parameters must be added to this endpoint for it to work:
   - type: Only support "odd" or "even".
   - items:  Is an Int and is the number of valid items you need to display as a response.
   - items_per_workers: Is an Int and is the number of valid items the worker should append to the response.
   ```
   Eg. http://localhost:8000/worker-pool?type=Odd&items=8&items_per_workers=2
   ```