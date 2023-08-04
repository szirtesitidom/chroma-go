/*
Testing Chroma Client
*/

package test

import (
	"fmt"
	chroma "github.com/amikos-tech/chroma-go"
	openai "github.com/amikos-tech/chroma-go/openai"
	godotenv "github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func Test_chroma_client(t *testing.T) {

	client := chroma.NewClient("http://localhost:8000")

	t.Run("Test Heartbeat", func(t *testing.T) {
		resp, err := client.Heartbeat()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Truef(t, resp["nanosecond heartbeat"] > 0, "Heartbeat should be greater than 0")
	})

	t.Run("Test CreateCollection", func(t *testing.T) {
		collectionName := "test-collection"
		metadata := map[string]string{}
		err := godotenv.Load("../.env")
		if err != nil {
			assert.Failf(t, "Error loading .env file", "%s", err)
		}
		embeddingFunction := openai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
		distanceFunction := chroma.L2
		_, errRest := client.Reset()
		if errRest != nil {
			assert.Fail(t, fmt.Sprintf("Error resetting database: %s", errRest))
		}
		resp, err := client.CreateCollection(collectionName, metadata, true, embeddingFunction, distanceFunction)
		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, collectionName, resp.Name)
		fmt.Printf("resp: %v\n", resp.Metadata)
		assert.Equal(t, 2, len(resp.Metadata))
		//assert the metadata contains key embedding_function
		assert.Contains(t, chroma.GetStringTypeOfEmbeddingFunction(embeddingFunction), resp.Metadata["embedding_function"])
	})

	t.Run("Test Add Documents", func(t *testing.T) {
		collectionName := "test-collection"
		metadata := map[string]string{}
		err := godotenv.Load("../.env")
		if err != nil {
			assert.Failf(t, "Error loading .env file", "%s", err)
		}
		embeddingFunction := openai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
		distanceFunction := chroma.L2
		_, errRest := client.Reset()
		if errRest != nil {
			assert.Fail(t, fmt.Sprintf("Error resetting database: %s", errRest))
		}
		resp, err := client.CreateCollection(collectionName, metadata, true, embeddingFunction, distanceFunction)
		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, collectionName, resp.Name)
		fmt.Printf("resp: %v\n", resp.Metadata)
		assert.Equal(t, 2, len(resp.Metadata))
		//assert the metadata contains key embedding_function
		assert.Contains(t, chroma.GetStringTypeOfEmbeddingFunction(embeddingFunction), resp.Metadata["embedding_function"])
		documents := []string{
			"Document 1 content here",
			"Document 2 content here",
		}
		ids := []string{
			"ID1",
			"ID2",
		}

		metadatas := []map[string]string{
			{"key1": "value1"},
			{"key2": "value2"},
		}
		//_, _ := embeddingFunction.CreateEmbedding(documents)
		_, addError := resp.Add(nil, metadatas, documents, ids)
		require.Nil(t, addError)
	})

	t.Run("Test Get Collection Documents", func(t *testing.T) {
		collectionName := "test-collection"
		metadata := map[string]string{}
		err := godotenv.Load("../.env")
		if err != nil {
			assert.Failf(t, "Error loading .env file", "%s", err)
		}
		embeddingFunction := openai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
		distanceFunction := chroma.L2
		_, errRest := client.Reset()
		if errRest != nil {
			assert.Fail(t, fmt.Sprintf("Error resetting database: %s", errRest))
		}
		resp, err := client.CreateCollection(collectionName, metadata, true, embeddingFunction, distanceFunction)
		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, collectionName, resp.Name)
		assert.Equal(t, 2, len(resp.Metadata))
		//assert the metadata contains key embedding_function
		assert.Contains(t, chroma.GetStringTypeOfEmbeddingFunction(embeddingFunction), resp.Metadata["embedding_function"])
		documents := []string{
			"Document 1 content here",
			"Document 2 content here",
		}
		ids := []string{
			"ID1",
			"ID2",
		}

		metadatas := []map[string]string{
			{"key1": "value1"},
			{"key2": "value2"},
		}
		col, addError := resp.Add(nil, metadatas, documents, ids)
		require.Nil(t, addError)

		col, geterr := col.Get(nil, nil, nil)
		require.Nil(t, geterr)
		assert.Equal(t, 2, len(col.CollectionData.Ids))
		assert.Contains(t, col.CollectionData.Ids, "ID1")
		assert.Contains(t, col.CollectionData.Ids, "ID2")
	})

	t.Run("Test Query Collection Documents", func(t *testing.T) {
		collectionName := "test-collection"
		metadata := map[string]string{}
		err := godotenv.Load("../.env")
		if err != nil {
			assert.Failf(t, "Error loading .env file", "%s", err)
		}
		embeddingFunction := openai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
		distanceFunction := chroma.L2
		_, errRest := client.Reset()
		if errRest != nil {
			assert.Fail(t, fmt.Sprintf("Error resetting database: %s", errRest))
		}
		resp, err := client.CreateCollection(collectionName, metadata, true, embeddingFunction, distanceFunction)
		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, collectionName, resp.Name)
		assert.Equal(t, 2, len(resp.Metadata))
		//assert the metadata contains key embedding_function
		assert.Contains(t, chroma.GetStringTypeOfEmbeddingFunction(embeddingFunction), resp.Metadata["embedding_function"])
		documents := []string{
			"This is a document about cats. Cats are great.",
			"this is a document about dogs. Dogs are great.",
		}
		ids := []string{
			"ID1",
			"ID2",
		}

		metadatas := []map[string]string{
			{"key1": "value1"},
			{"key2": "value2"},
		}
		col, addError := resp.Add(nil, metadatas, documents, ids)
		require.Nil(t, addError)

		qr, qrerr := col.Query([]string{"I love dogs"}, 5, nil, nil, nil)
		require.Nil(t, qrerr)
		fmt.Printf("qr: %v\n", qr)
		assert.Equal(t, 2, len(qr.Documents[0]))
		assert.Equal(t, documents[1], qr.Documents[0][0]) //ensure that the first document is the one about dogs
	})

	t.Run("Test Count Collection Documents", func(t *testing.T) {
		collectionName := "test-collection"
		metadata := map[string]string{}
		err := godotenv.Load("../.env")
		if err != nil {
			assert.Failf(t, "Error loading .env file", "%s", err)
		}
		embeddingFunction := openai.NewOpenAIEmbeddingFunction(os.Getenv("OPENAI_API_KEY"))
		distanceFunction := chroma.L2
		_, errRest := client.Reset()
		if errRest != nil {
			assert.Fail(t, fmt.Sprintf("Error resetting database: %s", errRest))
		}
		resp, err := client.CreateCollection(collectionName, metadata, true, embeddingFunction, distanceFunction)
		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, collectionName, resp.Name)
		assert.Equal(t, 2, len(resp.Metadata))
		//assert the metadata contains key embedding_function
		assert.Contains(t, chroma.GetStringTypeOfEmbeddingFunction(embeddingFunction), resp.Metadata["embedding_function"])
		documents := []string{
			"This is a document about cats. Cats are great.",
			"this is a document about dogs. Dogs are great.",
		}
		ids := []string{
			"ID1",
			"ID2",
		}

		metadatas := []map[string]string{
			{"key1": "value1"},
			{"key2": "value2"},
		}
		col, addError := resp.Add(nil, metadatas, documents, ids)
		require.Nil(t, addError)

		countDocs, qrerr := col.Count()
		require.Nil(t, qrerr)
		fmt.Printf("qr: %v\n", countDocs)
		assert.Equal(t, int32(2), countDocs)
	})

}
