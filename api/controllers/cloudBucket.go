package controllers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

var (
	storageClient *storage.Client
)

const jactezCredentials = ` {
	"type": "service_account",
	"project_id": "jactez-com-cal",
	"private_key_id": "4328d5cd6cfd5e81fef50b54a30e86858bdea834",
	"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCFsw18Bb0xa3dA\n+MK5YclAweSA19hkon3PZ6Lyjb9+tLeAS4YyKlW4RFJw+LFI53vkSufV6aUkCD+R\nBzu7K+q7aOJzGHeu5iXXeloPsHtEJ+tJIedRwUZSw2KyQ6m30So++9XOdSPaqWGm\nsYyM+pXupiePUHGJCnEt4b32mP4DoOQ00e5DrM4xHb9bEplWTpCo2ZUjIM9uaIAl\nHAljzWmDpWEkeVBuOQGaATCmUQmqn1JEeT74BbErzMhOU7WhmMx6VENo7589CIfv\nigOiPKzB9Ggrk6euykScR/x4EYRyv0Mr5rYYCAfnOl7WX5zZs+QRa+y9lVYtzKwp\nQcUniJLRAgMBAAECggEAA+42S0XFQMVZWtOgBeCAu3f3xmySCVitkiyzIaZWfDHG\nje2qeExo+2MmzqNrstT7U27QT1LwbpVqKP1UHYs2T2QkpCA9WdWc7y9to7XF1JnH\n1P07gQUe6CKjrXydFLEAKPHbtezDrRNQXPIt9kcM+VhxI/Qdd7CS5eJpD0LYKjXE\nc+fsX/TG99HcVU9OKtyKAlwnfh8oyv1lbFOyayng6lWc0dvhkCuLGfBhsAJq6VSr\n39vQrMoZPgyC0q5tV6vX49B8kUS8lD5vXTscP68tO/OzT//DAc5UNe6RKuqkwHql\nuGiV/6alZLg+iBRvAfmYMqBHZ7MZQP0Awcafhe898QKBgQC8SnQThH4iFnuVf/sB\n7Nvcxmj0z1okixzywjUHRaTRWKULsXDAfrNwOrpAqKvAIqM9Cy7/LBB7Jhw6jGjg\nG8C8Sf5e7Ai5mlpwIE1JiAdBZwhP3rSLHRYabCxXBvUbYjN7xkhxQGszHrQ1YYrl\nzXQPWN3UzRR6HOM/sJhc5tmR5wKBgQC1xw+wPmi/Vd9xAh2vN9ZUXO0aRF1mZQqn\n6LBBKXUeQHTXjpPjjLH0kD4hA00Gxjll8IR2/bnyJjtVWZRgI7crVxWIIVAQM42M\nEpWkphD6nkwJor3CihXeJobFqlwOS/muit7KZDGvqdWhWjeEDy5bgypvde2wHkDj\nJ8fiPl4OhwKBgEwfzc4WQuiFLnHCzDh7Cmi3zrcrHcaod4ut+MJ35aq9q/yOQIeS\nsfktxR9fEhEb7+M+IkIIDqG7Rq5lFgGFNubpA25c/yoKvYWXiaew1z4Z6cJgx512\npPkJwuNsbKwlh6sC/0bKRIzmXPU3+m/uIH4T75uZTi4Qf8/AFdl5e30BAoGAZjQ2\nLfHpEytFJlT6O7o5V9Wnuk0V9qx5AU9jSj/1Cb1T9J7Fp/tDUy6GwCkK9fQd8aL5\n161xDyVP1v0235c1NbkQ8ilIytMxksAgQyLcCQ1X01MdPnRFN1KSFDFk8OTmzaxm\n94S+KvZilwYSkL24Ytus0F6N9agM86s6R4qpd6ECgYEAk26Sy/N7JQ/nUdEZOW5n\nHClN3k244NZiD42MEpGPpifChlDb+sLCgVF5c/4Xyr59wiGG96/z3uVx9naJT9hH\n86lP5mvuFZ2MARROTbOUVujvVPQhUMQ3QEeCCKHv4qjqHArfvYNzM2AiYuI2hDde\niUvSWuP2kej51H5kF/+panM=\n-----END PRIVATE KEY-----\n",
	"client_email": "jactez-bucket-service@jactez-com-cal.iam.gserviceaccount.com",
	"client_id": "111212006094252198730",
	"auth_uri": "https://accounts.google.com/o/oauth2/auth",
	"token_uri": "https://oauth2.googleapis.com/token",
	"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	"client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/jactez-bucket-service%40jactez-com-cal.iam.gserviceaccount.com",
	"universe_domain": "googleapis.com"
 }`

func getNewGoogleStorageClient(c *gin.Context) (*storage.Client, context.Context, error) {
	var err error

	ctx := appengine.NewContext(c.Request)

	var o any
	json.Unmarshal([]byte(jactezCredentials), o)

	storageClient, err = storage.NewClient(ctx, option.WithCredentialsJSON([]byte(jactezCredentials)))
	return storageClient, ctx, err

}

// HandleFileUploadToBucket uploads file to bucket
func HandleFileUploadToBucket(c *gin.Context) {
	bucket := "jactez01" //your bucket name

	var err error
	var ctx context.Context
	storageClient, ctx, err = getNewGoogleStorageClient(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	f, uploadedFile, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	defer f.Close()

	sw := storageClient.Bucket(bucket).Object(uploadedFile.Filename).NewWriter(ctx)

	if _, err := io.Copy(sw, f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	if err := sw.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	u, err := url.Parse("/" + bucket + "/" + sw.Attrs().Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"pathname": u.EscapedPath(),
	})
}

func DeleteBucketObject(c *gin.Context) {
	bucket := "jactez01" //your bucket name
	var bucketObject string = c.Param("id")

	var err error
	var ctx context.Context
	storageClient, ctx, err = getNewGoogleStorageClient(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}
	defer storageClient.Close()

	o := storageClient.Bucket(bucket).Object(path.Base(bucketObject))
	if err := o.Delete(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

}

/*
// deleteFile removes specified object.
func deleteFile(w io.Writer, bucket, object string) error {
	// bucket := "bucket-name"
	// object := "object-name"
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	o := client.Bucket(bucket).Object(object)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to delete the file is aborted
	// if the object's generation number does not match your precondition.
	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", object, err)
	}
	fmt.Fprintf(w, "Blob %v deleted.\n", object)
	return nil
}
*/
