[
  {
    "name": "snooper",
    "image": "${REPOSITORY_URL}:latest",
    "networkMode": "awsvpc",
    "essential": true,

	  "environment": [
     {
    	"name": "SWHURL",
      "value": "${SWHURL}"
   	 }
    ]
  }
]