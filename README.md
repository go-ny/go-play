This should be a forum application, with surrounding functions.

The application is based on the following tech stack:

I. backend

1. Go (with gonic/gin framework: https://github.com/gin-gonic/gin)
2. MongoDB 
3. Redis (https://github.com/go-redis/redis)

II. frontend

1. React (and its bucket: using this: https://github.com/facebook/create-react-app)

III. DevOps & Others:

a. AWS: 

1. S3 bucket
2. EC2

b. Docker


Create your AWS account, AWS S3 bucket, and change all the keys sets in `.env` file.

Create your MongoDB account, create a DB called `my_blogs`, and a collection in it called `blogs`, copy connection URI 
to the file `./consts/dumTest.go`

Test:
Better first write something to your collection:
POST
`localhost:4000/post-blog`
in body
```
{
	"views": 0,
	"cover":"oo",
	"title": "dsds",
	"content":"dd"
}
```

Test Your Upload
POST
`localhost:4000/upload`
choose `multi-part`
choose `file`
and choose an image

You might need to change line 128 in `./awsOp/awsOp.go` to your bucket returning link

have fun!