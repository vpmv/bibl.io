Bibl.io Library API
---

No-nonsense Go implementation of developer assignment "Build a Books/Authors API". 
[Read assignment](ASSIGNMENT.md) | [Work log](WORKLOG.md)

Features
- <strike>Blazing fast API</strike> Broken AF API 
- Generics (one job to rule them all 💍)
- Concurrent job handling using basic channels
  - Sync mutexes to prevent firing the same jobs
- OpenAPI v3
  - with separate public and internal models using DTOs
- Interchangeable storage backends
  - Implemented GORM for now
- Easily manageable API with version support <br>
  Binding /v2/ routes/packages is a breeze
- Authorization modules
  - mocked for now

Known issues
- Altough it worked fine a day ago, something is now causing deadlocks. The first few queries are resolved (to an extent) and then it stops listening. The API pretends to do something, but the job queue will stall. 
- 
Todo
- Fix above issues
- Handle external API rate limiting <br> 
    The API now floods their API with requests. I implemented a QueryRegistry to try and keep somewhat limit queries, but it's very bare bones
- Unit tests 

Requirements
- docker
  - `docker network create dev_proxy` to make PHP and Go api play nice
- Go v1.22+ optional