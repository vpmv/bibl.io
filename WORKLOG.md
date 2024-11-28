Work log
---

Built using an old dev laptop, which was struggling to keep up with modern day.

| Date      | Tasks                                                                                                           | comment                                                                              |
|-----------|-----------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------|
| sun 24/11 | Setup bare minimum API using old boilerplate code                                                               |                                                                                      |
|           | Investigate libraries exposing OpenAPI v3<br>Discovered Fuego and went with it                                  | Fuego is new to me but it's def up to par                                            |
| mon 25/11 | Started modelling database structs                                                                              | Should've started looking at OpenLibrary API first                                   |
|           | Wrote draft Job/Queue logic                                                                                     |                                                                                      |
| tue 26/11 | Work out OpenLibrary responses                                                                                  | O_O                                                                                  |
|           | Bang head against wall a few times                                                                              |                                                                                      |
|           | Initial API testing                                                                                             |                                                                                      |
| wed 27/11 | Job signatures grew out of hand. Decided to implement generics anyway. This proved more difficult then expected |                                                                                      |
| thu 28/11 | Finalize generics. Refactored queueing system                                                                   |                                                                                      |
|           | Test API                                                                                                        | Database TX seem broken, fix later                                                   |
|           | Re-read assignment                                                                                              | 💩                                                                                   |
|           | Create bare minimum symfony frontend to pair with it                                                            | symfony CLI FTW \o/                                                                  |
|           | Updating machine to support php 8.4 <br>Setup Sf7.1 + docker                                                    |                                                                                      |
|           | Symfony's docker can't access localhost (logically), compiling docker image for Go                              | Lost a lot of time running into a weird compilation error, caused by a faulty import | 
|           | Resolving issues with docker networking                                                                         | .. so much time                                                                      | 
|           | Check database logic for errors                                                                                 | Still not writing correctly. Need to fix later                                       | 
|           | More testing                                                                                                    |                                                                                      | 
|           | Add author routes. Fix generic support.                                                                         |                                                                                      | 
|           | Running into thread/deadlocks. Trying to debug/refactor sync mutex. It's getting too late to solve. :(          | To be continued                                                                      | 