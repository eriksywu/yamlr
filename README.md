# yamlr
## API
1. `CREATE METADATA `
* url: `POST /api/metadata`
* body: metadata taml 
* returns: guid associated with metadata or error
2. `GET METADATA (extra functionality)`
* url: `GET /api/metadata/{guid}`
* returns: metadata or error
3. `UPDATE METADATA (extra functionality)`
* url: `PUT /api/metadata/{guid}`
* body: metadata
* returns: none or error
4. `SEARCH METADATA`
* url: `POST /api/search`
* body: metadata search taml (see notes below)
* returns: list of matching metadata or error

## SEARCH
Searching on all fields except for Description and Maintainer Name are supported. For example, to search for metadata with company=Microsoft, Title=AKS, Maintainers=[erik.wu@microsoft.com]:
```
company: Microsoft
title: AKS
maintainers:
    -
        email: erik.wu@microsoft.com
```
Fuzzy logic or complex queries are not supported

## Implementation Notes
### Searching
Metadata is stored in an in-memory repository, as implemented in `memoryrepo/repo.go`. Because the stored metadata should be queryable, I decided not to implement a naive in-memory cache as querying will not be efficient and also tedious to code. The in-memory repository is built on the `github.com/hashicorp/go-memdb` package which provides an in-memory, schema-based object store with rudimentary indexing capabilities, thus allowing for easier querying. 

The `description` field was interpreted as a free-text field in this exercise so it should be free-text queryable. Thus querying against the `description` field is not supported because text searching is not trivial to implement. In a 'real-life' application, I would off-load all indexing/querying to a managed search service (i.e Azure Search).

I decided to implement the search endpoint as `POST /api/search/{guid}` instead of `GET /api/search/{query params}` mostly because it's easier to test with.

### Bugs and TODOs
1. Data denormalization between the metadata and maintainer tables. Updates/Deletes are not 100% optimized because of this.
2. The update API does not fully update maintainer data association. The above TODO kinda preludes this.
3. HTTPS not supported
4. Shutdown logic not implemented
