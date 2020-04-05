# remember

Small app in Go to help me remember things, like watching tv series, movies, reading books etc

## Media Tracker

* replace Keep Watching
* media to support
    * games
    * books
    * movies
    * tv series
    * sport matches
    * link to source
* exportable calendar

### Model

* Consumer
* MediaItem --> interface
    * Name
    * Type
    * Source (URL)
    * MediaItems 0...* MediaItem
    * 
* MediaList
    * contains 0...* MediaItemEntry (ordered list)
    * Owner 1 Consumer
    * Co-Owners 0...* Consumer
    * Public
    * Viewers
* MediaItemEntry
    * MediaItem
    * DateAdded
    * DateChanged
    * Finished
    * Progression
* Progression
    * MediaItem
    * ...

Example of a TV Series
* MediaItem (TV Series)
    * Name: The Witcher
    * MediaItems
        * MediaItem 0:
            * Name: Season 1
            * Type: TV Series Season
            * MediaItems
                * MediaItem 0:
                    * Name Episode 1: ...
                    
## Links

* https://echo.labstack.com/middleware/logger
* https://dev.to/douglasmakey/oauth2-example-with-go-3n8a
* https://fly.io/docs/app-guides/continuous-deployment-with-github-actions/
* https://sbstjn.com/host-golang-slackbot-on-heroku-with-hanu.html