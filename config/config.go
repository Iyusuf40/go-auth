package config

const ApiPort = "8080"
const AuthPort = "8081"
const BaseApiUrl = "http://localhost:" + ApiPort + "/api/"

const DBMS = "file"
const HOST = "localhost"
const USER = "yusuf"
const PASSWORD = "0"

const UsersDatabase = "test"
const UsersRecords = "users"

const TempStoreDb = "test"
const TempStoreType = "file"
const RedisUrl = "localhost:6379"
const RedisPassword = ""

const GmailPassword = "vcmu kvwj kjgp lxyi"
const GmailSource = "punkgts40@gmail.com"

const RequireEmailVerification = true

// you can change this message to your liking
// so far as you keep the `##link##` somewhere in your
// custom message, it will be substituted with the string
// "link" which will be a live link for users to confirm
// their emails and complete signup.
const EmailConfirmationMessage = `
<div>
	<h1>Welcome</h1>
	<p>Complete signup by clicking ##link##.</p>
	<br>
	<br>
	<br>
	<p>powered by Go-Auth https://github.com/iyusuf40/go-auth.</p>
</div>
`

// if you change this, make sure to use the new value in
// EmailConfirmationMessage
const LinkSubstitute = "##link##"
