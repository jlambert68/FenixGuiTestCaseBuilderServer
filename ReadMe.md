MIT License

Copyright (c) 2024 Jonas Lambert

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

***

# Fenix Inception

## BuilderServer
BuilderServer has the responsibility to server TesterGui with all components that the user is allowed to use when building TestCase. BuilderServer is also responsible to save, load and list TestCases for the user. 

![Fenix Inception - BuilderServer](./Documentation/FenixInception-Overview-NonDetailed-BuilderServer-Worker.png "Fenix Inception - BuilderServer")

The following environment variable is needed for BuilderServer to be able to run.

| Environment variable                        | Example value                            | comment                                                                        |
|---------------------------------------------|------------------------------------------|--------------------------------------------------------------------------------|
| DB_HOST                                     | 127.0.0.1                                |                                                                                |
| DB_NAME                                     | fenix_gcp_database                       |                                                                                |
| DB_PASS                                     | database password                        |                                                                                |
| DB_POOL_MAX_CONNECTIONS                     | 4                                        | The number of connections towards the database that the ExectionServer can use |
| DB_PORT                                     | 5432                                     |                                                                                |
| DB_SCHEMA                                   | Not used                                 |                                                                                |
| DB_USER                                     | postgres                                 |                                                                                |
| ExecutionLocation    | GCP                                      |                                                                                |
| FenixGuiBuilderServerAddress | 127.0.0.1                                |                                                                                |
| FenixGuiBuilderServerPort                 | 6670                                     |                                                                                |
| LogAllSQLs                    | true                                     |                                                                                |


