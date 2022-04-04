= Building Neo4j Applications with Go

> Learn how to interact with Neo4j from Go using the Neo4j Go Driver

This repository accompanies the link:https://graphacademy.neo4j.com/courses/app-go/[Building Neo4j Applications with Go course^] on link:https://graphacademy.neo4j.com/[Neo4j GraphAcademy^].

For a complete walkthrough of this repository,  link:https://graphacademy.neo4j.com/courses/app-go/[enroll now^].

== Setup

* Clone repository
* Update config.json with the connection details
[source,json]
----
{
  "APP_PORT": 3000,
  "NEO4J_URI": "neo4j://localhost:7687",
  "NEO4J_USERNAME": "neo4j",
  "NEO4J_PASSWORD": "letmein",
  "JWT_SECRET": "secret",
  "SALT_ROUNDS": 10
}
----

* Start the project

----
go run ./cmd/neoflix
----

== A Note on comments

You may spot a number of comments in this repository that look a little like this:

[source,java]
----
// tag::something[]
someCode()
// end::something[]
----


We use link:https://asciidoc-py.github.io/index.html[Asciidoc^] to author our courses.
Using these tags means that we can use a macro to include portions of code directly into the course itself.

From the point of view of the course, you can go ahead and ignore them.