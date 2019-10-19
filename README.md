# VoronoiDiagrams
The goal of this project is to build Voronoi Diagrams using Fortunes Algorithm.

https://github.com/VenkoChakalov/VoronoiDiagrams/blob/master/images/demo.png

## Usage
Execute the following:

```
   go run run.go -port 8333
```
Options:
 
 -  -port = select a port for the server
 -  -proxies = select the count of proxies behind the server
 
### API
To use the API provided with the implementation:
 - The endpoint to obtain the results from diagram partitioning is:
    - /api/voronoi POST
 - The function accepts a JSON array of objects:
    -  ```json
       [
         {
          "X":160,
      
          "Y":330
         }
       ]
       ```
 - The response we get from the server is:
    - ```json
      [
         {
           "Start": {
              "X":197.1625,
              "Y":346
            },
           "End":  {
              "X":1140.1875275354098,
              "Y":-797.0606394368613
            },
           "Done": true
         }
      ]
      ```
 - Sample ajax request:
   - ```javascript
         $.ajax({
                    url: ".../api/voronoi",
                    type: "POST",
                    data: {
                        points: JSON.stringify(points)
                    },
                    success: function (response) {
                        /.../
                    },
                    error: function (jqXHR, textStatus, errorThrown) {
                       /.../
                    }
                })
     ```


## Improvement:


  - Replace the LinkedList which represents the Beach Line with a balancing search tree, to achieve O(nlogn) complexity.


  - Remove priority queue casting, it's dangerous and slow.





