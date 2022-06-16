## Gets the number of [ApsimX](https://github.com/APSIMInitiative/ApsimX) issues closed by a user

Before running the script, you need to create a file called credentials.dat that looks like this:

username=GITHUBUSERNAME
token=GITHUBTOKEN

where GITHUBUSERNAME is your user name and GITHUBTOKEN is your GitHub personal token


# Running in Docker:

1. Fork and/or clone this repo
2. docker build -t apsimissues .
3. docker run -it -v $(pwd):/wd apsimissues


# Running without Docker:

1. Download and install [Go](https://golang.org/dl/)
2. Fork and/or clone this repo
3. Compile
```cmd
cd %userprofile%\go\src\ApsimIssues
go build
```

```sh
cd ~/go/src/ApsimIssues
go build
```
4. Run