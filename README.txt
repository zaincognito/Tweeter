Distributed and Parallel Systems Project Part 4


Contributors: Nabil Ahmed (na1489) and Zulkifl Arefin (za453)


How to Build: Have four terminal windows open. For the first terminal change directories to Tweeter and then change directories to Frontend. Type go build, hit enter and then ./Frontend and enter. For the second terminal change directories to Tweeter and then change directories to Backend1. Type go build, hit enter and then ./Backend1 and enter. For the third terminal change directories to Tweeter and then change directories to Backend2. Type go build, hit enter and then ./Backend2 and enter. For the fourth terminal change directories to Tweeter and then change directories to Backend4. Type go build, hit enter and then ./Backend4 and enter. Go to a browser and then go to localhost:7050/


How data is kept in sync across replicas:

There is one master server (frontend chooses who that is). It is the most updated one and it does that by checking which server has the highest transaction number (which is incremented by


FILES:

***BACKEND***

addFriend.go: takes username and friend to add to the username friends, written to friends text file.

backend.go: server, takes in a string from client, frontend.go and parses it. The first part of string is the command, with this command it goes into if-else structures to then do what the frontend needs.

bremoveAcct: removes the user from every file. Removes any relationships in the friends file, if user is someone's friend or user has friends, it deletes each of these relationships. Removes all posts that are by the user and then deletes the user from the users.txt file.

bsignin.go: takes in login input and checks if valid user, if so, sends back confirmation to frontend through backend.go.

bsignup.go: takes in signup info and if it is valid info, it puts it into the file system.

locking.go: has functions acquireLock() and releaseLock() to handle the locking system for this concurrent application. A RWMutex is used, which is a shared lock when reading and an exclusive lock when writing. 

removeFriend.go: removes a friend (unfollow) from a users friend from the friends.txt file.

retrievePosts.go: Sends back one string of all posts in a nice fashion.

writePost.go: Takes input of a potential post and writes it to the posts text file.




***FRONTEND***

flogout.go: when logout is clicked on frontend, deletes cookie.

fnewsfeed.go: sets cookie, displays newsfeed by getting different things from backend, such as name of user, can put in form of searching friends, shows posts of all friends and self.

fprofile.go: sets cookie, displays profile of user through calls to backend retrieving info just like newsfeed.

frontend.go: sets up server for frontend, has handler funcs.

fsignin.go: displays signin page, makes request to backend to check signin creds.

fsignup.go: displays signup page, makes request to backend to create user.

logout.html: dummy file in order for the user to be able to click on a link and then the home.go file processes this and logs a user out.

navbar.html: is the navigation bar on top of the newsfeed and profile where it shows links such as "Newsfeed", "Profile", "Logout" and "Delete Account".

newsfeed.html: is the html file which shows the newsfeed, it has a form for the user to post tweets as well, see their own tweets and friends tweets.

profile.html: is the html file which shows the user profile, shows all of the friends the user has, it has a form for the user to post tweets, and it shows the tweets of the user only.

removeAcct: is a form to ask the user if they are sure they want to delete their account, if it is a yes home.go deletes the account, if no then they are redirected to their newsfeed.

signin.html: is a form for anyone to login, and a link for signing up.

signup.html: is a form for anyone to signup.




***DATABASES***

friends.txt: file for all friend relationships, where each line has one relationship: user,friendTheyFollow.

postCount.txt: a storage of current count of number of posts, to create post IDs.

posts.txt: storage of all posts in format: user,postID,post.

users.txt: storage of users in format: username,password.
