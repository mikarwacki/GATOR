# GATOR
This is blog aggreGATOR running in the CLI, it allows logged users to fetch feeds of their followed channels that are supporting RSS feeds.

# Why?
I created gator to be able to read latest hackernews posts and cbs news directly without exiting my terminal during daily workflow.

# Quickstart
### Prerequisites: postgres installed, go 1.23 installed
```bash
brew install go
brew install postgresql@15
```
### Install gator using go install command
```bash
github.com/mikarwacki/gator@latest
```
### Set-up the database:
from the terminal open psql tui: psql postgres
Now run: CREATE DATABASE gator;
Create a config file in your home directory called ".gatorconfig.json"
In the created file input URI to your postgres database: 
```json
{"db_url":"postgres://your_username:@localhost:5432/gator?sslmode=disable"}
```
Now you can run your gator commands by typing
```bash
./gator command_name
```

# Usage
Available commands:
- register "user_name" -> registers a new user in your database, it accepts one parameter (user_name)
- login "user_name" -> logins into one of registered user, so that gator knows which feed to fetch, one parameter required
- reset -> reset is used to reset all users from your database, it doesn't require an argument
- users -> displays all registered users
- agg "time_duration" -> this function fetches and aggregates the RSS feeds that the currently loged in user follows, it accept one argument which is the frequency of feed fetches
- addFeed "feed_name" "feed_url" -> adds the feed to the database and adds it to the follows list of the current user feed_name tells gator how to later tag a feed in agg func, feed_url is location of the news site you want to add and follow
- feeds -> displays all feeds added to the database by all users
- follow "feed_name" -> allows user to follow currently existing feed, requires one argument that is name of the feed
- following -> displays all the feeds followed by the current user
- unfollow "feed_name" -> allows user to unfollow one of theirs feeds
- browse "number_of_feeds" -> displays the latest feeds followed by a user, requires one argument that tells how many feeds should be displayed

# Future improvements:
  - add BubbleTea support for a better user experience with nice TUI elements
  - host gator so that users can ssh into the server and browse their feeds
